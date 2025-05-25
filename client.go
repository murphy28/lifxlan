package lifxlan

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

var BroadcastAddress = net.UDPAddr{IP: net.IPv4bcast, Port: LifxPort}

type Client struct {
	conn       *net.UDPConn
	identifier uint32
	mu         sync.Mutex
	devices    []Device
}

func NewClient() (*Client, error) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: LifxPort})
	if err != nil {
		return nil, fmt.Errorf("cannot bind to LIFX port: %w", err)
	}

	InitializeProducts()

	client := &Client{
		conn:       conn,
		identifier: rand.Uint32(),
	}

	return client, nil
}

// Close closes the UDP connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// BroadcastPacket sends a packet to the broadcast address
func (c *Client) BroadcastPacket(packet []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Send the packet to the broadcast address
	_, err := c.conn.WriteToUDP(packet, &BroadcastAddress)
	if err != nil {
		return fmt.Errorf("failed to send broadcast packet: %w", err)
	}

	return nil
}

// Send sends a packet to a specific address
func (c *Client) Send(packet []byte, addr *net.UDPAddr) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.conn.WriteToUDP(packet, addr)

	return err
}

// SendAndWait sends a packet and waits for a response
func (c *Client) SendAndWait(packet []byte, addr *net.UDPAddr, expectedType PacketType, timeout time.Duration) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// fmt.Printf("Sending packet: %s\n", hex.EncodeToString(packet))
	// fmt.Printf("Sending packet to %s\n", addr.String())
	_, err := c.conn.WriteToUDP(packet, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to send packet: %w", err)
	}

	c.conn.SetReadDeadline(time.Now().Add(timeout))
	buf := make([]byte, 1500)
	for {
		len, _, err := c.conn.ReadFromUDP(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				return nil, fmt.Errorf("timeout waiting for response: %w", err)
			}
			return nil, fmt.Errorf("failed to read from UDP: %w", err)
		}

		h, err := ParseHeader(buf[:len])
		if err != nil {
			return nil, fmt.Errorf("failed to parse header: %w", err)
		}

		if h.Source() == c.identifier && h.Type() == expectedType {
			return buf[:len], nil
		}
	}
}

// Discover sends a discovery packet and listens for responses
func (c *Client) Discover(timeout time.Duration) error {
	packet := BuildDiscoveryPacket(c.identifier)

	// Send the discovery packet
	if err := c.BroadcastPacket(packet); err != nil {
		return fmt.Errorf("failed to send discovery packet: %w", err)
	}

	// Set read deadline for responses
	deadline := time.Now().Add(timeout)
	c.conn.SetReadDeadline(deadline)

	// Listen for responses
	buf := make([]byte, 1500)
	for {
		len, remote, err := c.conn.ReadFromUDP(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				break // Timeout reached
			}
			return fmt.Errorf("failed to read from UDP: %w", err)
		}

		// Parse the header from the received packet
		h, err := ParseHeader(buf[:len])
		if err != nil {
			return fmt.Errorf("failed to parse header: %w", err)
		}

		// Check if the packet is a response to our discovery request
		if h.Source() != c.identifier || h.Type() != StateService {
			continue // Ignore other packets
		}

		// Add the device to the list of discovered devices
		device := NewDevice(h.Target(), remote.IP, c)

		seen := false

		// Check if the device is already discovered
		for _, device := range c.devices {
			if bytes.Equal(device.MAC, h.Target()) {
				seen = true
				break
			}
		}

		// Add the new device to the list if new
		if !seen {
			c.devices = append(c.devices, *device)
			fmt.Printf("Discovered device: %s at %s\n", hex.EncodeToString(device.MAC), device.IP.String())
		}
	}

	// Retrieve device information
	c.RefreshDeviceInfo()

	return nil
}

// LoadDevices loads devices from a JSON string
func (c *Client) LoadDevices(jsonData string) error {
	// Unmarshal the JSON data into a slice of Device structs
	var devices []Device
	if err := json.Unmarshal([]byte(jsonData), &devices); err != nil {
		return fmt.Errorf("failed to unmarshal devices JSON: %w", err)
	}

	// Clear existing devices
	c.mu.Lock()
	c.devices = []Device{}

	// Add each device to the client's device list
	for _, device := range devices {
		newDevice := NewDevice(device.MAC, device.IP, c)
		newDevice.Label = device.Label
		newDevice.Product = device.Product

		c.devices = append(c.devices, *newDevice)
		fmt.Printf("Loaded device: %s at %s\n", hex.EncodeToString(newDevice.MAC), newDevice.IP.String())
	}
	c.mu.Unlock()

	// Refresh device information
	c.RefreshDeviceInfo()

	return nil
}

func (c *Client) RefreshDeviceInfo() {
	for i := range c.devices {
		device := &c.devices[i]
		if err := device.RefreshInfo(); err != nil {
			fmt.Printf("Error refreshing device info for %s: %v\n", hex.EncodeToString(device.MAC), err)
			continue
		}
	}
}

// ExportJSON exports the discovered devices as a JSON string
func (c *Client) ExportJSON() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	jsonBytes, err := json.MarshalIndent(c.devices, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return "[]"
	}

	return string(jsonBytes)
}

// Devices returns the list of discovered devices
func (c *Client) GetDevices() []Device {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.devices
}

// ClearDevices clears the list of discovered devices
func (c *Client) ClearDevices() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.devices = []Device{}
}

func (c *Client) GetDeviceByLabel(label string) (*Device, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, device := range c.devices {
		// Compare lowercase and trimmed labels
		if SanitizeLabel(device.Label) == SanitizeLabel(label) {
			return &device, nil
		}
	}

	return nil, fmt.Errorf("device with label %s not found", label)
}

func SanitizeLabel(label string) string {
	// Remove leading and trailing whitespace
	label = strings.TrimSpace(label)

	// Convert to lowercase
	label = strings.ToLower(label)

	// Limit label length to 32 characters
	if len(label) > 32 {
		label = label[:32]
	}

	return string(label)
}

// GetDeviceLabels returns a list of labels for all discovered devices
func (c *Client) GetDeviceLabels() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	labels := make([]string, len(c.devices))
	for i, device := range c.devices {
		labels[i] = device.Label
	}
	return labels
}

// TurnOn turns on the device that matches the given label
func (c *Client) TurnOn(label string) error {
	device, err := c.GetDeviceByLabel(label)
	if err != nil {
		return err
	}
	return device.TurnOn()
}

// TurnOff turns on the device that matches the given label
func (c *Client) TurnOff(label string) error {
	device, err := c.GetDeviceByLabel(label)
	if err != nil {
		return err
	}

	fmt.Printf("Turning off device: %s\n", device.Label)
	return device.TurnOff()
}

// SetColor sets the color of the device that matches the given label
func (c *Client) SetColor(label string, color LIFXColor, duration time.Duration) error {
	device, err := c.GetDeviceByLabel(label)
	if err != nil {
		return err
	}
	return device.SetColor(color, duration)
}

// SetLabel sets the label of the device that matches the given label
func (c *Client) SetLabel(oldLabel, newLabel string) error {
	device, err := c.GetDeviceByLabel(oldLabel)
	if err != nil {
		return err
	}
	return device.SetLabel(newLabel)
}
