package lifxlan

import (
	"encoding/binary"
	"time"
)

// BuildDiscoveryPacket creates a discovery packet
func BuildDiscoveryPacket(source uint32) []byte {
	header := DefaultHeader(source, make([]byte, 6), GetService, 0)
	header.SetTagged(true)

	return header[:]
}

// BuildSetPowerPacket creates a packet to set the power state of a device
func BuildSetPowerPacket(source uint32, target []byte, on bool) []byte {
	payload := make([]byte, 2)
	if on {
		binary.LittleEndian.PutUint16(payload, 65535) // Power on
	} else {
		binary.LittleEndian.PutUint16(payload, 0) // Power off
	}

	header := DefaultHeader(source, target, SetPower, uint16(len(payload)))

	// add the payload to the header to create the packet
	packet := make([]byte, HeaderSize+len(payload))
	copy(packet, header[:])
	copy(packet[HeaderSize:], payload)
	return packet
}

// BuildSetColorPacket creates a packet to set the color of a device
func BuildSetColorPacket(source uint32, target []byte, color LIFXColor, duration time.Duration) []byte {
	payload := make([]byte, 13)
	binary.LittleEndian.PutUint16(payload[1:3], color.hue)
	binary.LittleEndian.PutUint16(payload[3:5], color.saturation)
	binary.LittleEndian.PutUint16(payload[5:7], color.brightness)
	binary.LittleEndian.PutUint16(payload[7:9], color.kelvin)

	// Set the duration in milliseconds
	durationMs := uint32(duration.Milliseconds())
	binary.LittleEndian.PutUint32(payload[9:13], durationMs)

	header := DefaultHeader(source, target, SetColor, uint16(len(payload)))

	// add the payload to the header to create the packet
	packet := make([]byte, HeaderSize+len(payload))
	copy(packet, header[:])
	copy(packet[HeaderSize:], payload)
	return packet
}

// BuildGetLabelPacket creates a packet to get the label of a device
func BuildGetLabelPacket(source uint32, target []byte) []byte {
	header := DefaultHeader(source, target, GetLabel, 0)

	return header[:]
}

// BuildSetLabelPacket creates a packet to set the label of a device
func BuildSetLabelPacket(source uint32, target []byte, label string) []byte {
	payload := make([]byte, 32)
	copy(payload, label)

	header := DefaultHeader(source, target, SetLabel, uint16(len(payload)))

	// add the payload to the header to create the packet
	packet := make([]byte, HeaderSize+len(payload))
	copy(packet, header[:])
	copy(packet[HeaderSize:], payload)
	return packet
}

// BuildGetVersionPacket creates a packet to get the version of a device
func BuildGetVersionPacket(source uint32, target []byte) []byte {
	header := DefaultHeader(source, target, GetVersion, 0)

	return header[:]
}

func BuildEchoRequestPacket(source uint32, target []byte, echo []byte) []byte {
	payload := echo

	header := DefaultHeader(source, target, EchoRequest, uint16(len(payload)))

	// add the payload to the header to create the packet
	packet := make([]byte, HeaderSize+len(payload))
	copy(packet, header[:])
	copy(packet[HeaderSize:], payload)
	return packet
}
