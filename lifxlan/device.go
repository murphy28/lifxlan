package lifxlan

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Device struct {
	MAC     []byte `json:"mac"`
	IP      net.IP `json:"ip"`
	client  *Client
	Label   string  `json:"label"`
	Product Product `json:"product"`
}

func NewDevice(mac []byte, ip net.IP, c *Client) *Device {
	device := &Device{
		MAC:    mac,
		IP:     ip,
		client: c,
	}

	return device
}

// ExportJSON exports the device information as a JSON string using marshall
func (d *Device) ExportDeviceJSON() string {
	jsonBytes, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return "{}"
	}

	return string(jsonBytes)
}

// RefreshInfo refreshes the device information by getting the label and product details
func (d *Device) RefreshInfo() error {
	// Ping the device to ensure it's reachable
	if !d.Ping() {
		return fmt.Errorf("device %s is not reachable", d.GetMACAddress())
	}

	// Obtain current device label
	_, err := d.GetLabel()
	if err != nil {
		return fmt.Errorf("failed to get label for device %s: %w", d.GetMACAddress(), err)
	}

	// Obtain current product information
	_, err = d.GetProduct()
	if err != nil {
		return fmt.Errorf("failed to get product for device %s: %w", d.GetMACAddress(), err)
	}

	return nil
}

func (d *Device) Send(packet []byte) error {
	return d.client.Send(packet, &net.UDPAddr{
		IP:   d.IP,
		Port: LifxPort,
	})
}

func (d *Device) SendAndWait(packet []byte, pktType PacketType, duration time.Duration) ([]byte, error) {
	return d.client.SendAndWait(packet, d.UDPAddr(), pktType, duration)
}

func (d *Device) TurnOn() error {
	packet := BuildSetPowerPacket(d.client.identifier, d.MAC, true)
	return d.Send(packet)
}

func (d *Device) TurnOff() error {
	packet := BuildSetPowerPacket(d.client.identifier, d.MAC, false)
	return d.Send(packet)
}

func (d *Device) SetColor(color LIFXColor, duration time.Duration) error {
	packet := BuildSetColorPacket(d.client.identifier, d.MAC, color, duration)
	return d.Send(packet)
}

func (d *Device) GetLabel() (string, error) {
	packet := BuildGetLabelPacket(d.client.identifier, d.MAC)
	response, err := d.SendAndWait(packet, StateLabel, 2*time.Second)
	if err != nil {
		return "", err
	}

	// The label is in bytes 36-52 of the response
	labelBytes := bytes.TrimRight(response[HeaderSize:HeaderSize+32], "\x00") // Get the raw bytes
	label := string(labelBytes)
	d.Label = label // Update the device's label field
	return label, nil
}

func (d *Device) SetLabel(label string) error {
	// Ensure the label is 32 bytes long
	if len(label) > 32 {
		label = label[:32]
	}

	packet := BuildSetLabelPacket(d.client.identifier, d.MAC, label)
	return d.Send(packet)
}

func (d *Device) GetProduct() (Product, error) {
	packet := BuildGetVersionPacket(d.client.identifier, d.MAC)
	response, err := d.SendAndWait(packet, StateVersion, 2*time.Second)
	if err != nil {
		return Product{}, err
	}

	vendorID := binary.LittleEndian.Uint32(response[HeaderSize : HeaderSize+4])
	productID := binary.LittleEndian.Uint32(response[HeaderSize+4 : HeaderSize+8])

	product, err := GetProduct(int(vendorID), int(productID))

	if err != nil {
		return Product{}, err
	}

	d.Product = product

	return product, nil
}

func (d *Device) Ping() bool {
	uniquePayload := make([]byte, 64)
	rand.Read(uniquePayload)

	packet := BuildEchoRequestPacket(d.client.identifier, d.MAC, uniquePayload)

	response, err := d.SendAndWait(packet, EchoResponse, 2*time.Second)
	if err != nil {
		fmt.Printf("Error sending echo request to device %s: %v\n", d.GetMACAddress(), err)
		return false
	}

	echoing := response[HeaderSize : HeaderSize+64]

	return bytes.Equal(echoing, uniquePayload)
}

func (d *Device) GetMACAddress() string {
	return net.HardwareAddr(d.MAC).String()
}

func (d *Device) UDPAddr() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   d.IP,
		Port: LifxPort,
	}
}
