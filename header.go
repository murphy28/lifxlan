package lifxlan

import (
	"encoding/binary"
	"errors"
)

// Header represents a 36-byte LIFX protocol header (https://lan.developer.lifx.com/docs/packet-contents#header)
type Header [HeaderSize]byte

// ParseHeader parses the first 36 bytes into a Header struct
func ParseHeader(data []byte) (*Header, error) {
	// Check if the data length is less than the header size
	if len(data) < HeaderSize {
		// Return an error indicating insufficient data
		return nil, errors.New("insufficient data to unpack header")
	}
	// Create a new Header instance and copy the first 36 bytes into it
	var h Header
	copy(h[:], data[:HeaderSize])

	// Return the pointer to the Header instance and nil error
	return &h, nil
}

// Size returns the total message size
func (h *Header) Size() uint16 {
	// Read the first 2 bytes as a little-endian uint16
	return binary.LittleEndian.Uint16(h[0:2])
}

// Protocol returns the protocol number
func (h *Header) Protocol() uint16 {
	// Read the first 2 bytes as a little-endian uint16
	v := binary.LittleEndian.Uint16(h[2:4])
	// Return the lower 12 bits of the value
	return v & 0x0FFF
}

// Addressable returns true if the addressable bit is set
func (h *Header) Addressable() bool {
	// Check if bit 12 of byte 2 is set
	return h[3]&0x10 != 0
}

// Tagged returns true if the tagged bit is set
func (h *Header) Tagged() bool {
	// Check if bit 13 of byte 2 is set
	return h[3]&0x20 != 0
}

// Source returns the client source identifier
func (h *Header) Source() uint32 {
	// Read bytes 4-7 as a little-endian uint32
	return binary.LittleEndian.Uint32(h[4:8])
}

// Target returns the 6-byte MAC address as a byte slice
func (h *Header) Target() []byte {
	// Extract bytes 8-13 (the MAC address)
	mac := h[8:14]
	// Convert the MAC address to a hex string and return it
	return mac
}

// ResponseRequired returns true if the response required bit (bit 0 of byte 15) is set
func (h *Header) ResponseRequired() bool {
	// Check if bit 0 of byte 15 is set
	return h[15]&0x01 != 0
}

// AckRequired returns true if the ack required bit (bit 1 of byte 15) is set
func (h *Header) AckRequired() bool {
	// Check if bit 1 of byte 15 is set
	return h[15]&0x02 != 0
}

// Sequence returns the sequence number
func (h *Header) Sequence() uint8 {
	// Shift byte 15 to the right by 2 bits and return it
	return h[15] >> 2
}

// Type returns the payload type
func (h *Header) Type() PacketType {
	// Read bytes 32-33 as a little-endian uint16 and return it as a PacketType
	return PacketType(binary.LittleEndian.Uint16(h[32:34]))
}

func NewHeader(source uint32) *Header {
	// Create a new Header instance
	h := &Header{}

	// Set the protocol number to the LIFX protocol
	protoField := uint16(Protocol)
	binary.LittleEndian.PutUint16(h[2:4], protoField)

	// Set default values for the header
	h.SetAddressable(true)
	h.SetSource(source)

	// Return the pointer to the new Header instance
	return h
}

// SetSize sets the total message size
func (h *Header) SetSize(size uint16) {
	binary.LittleEndian.PutUint16(h[0:2], size)
}

// SetAddressable sets the addressable bit to a boolean value
func (h *Header) SetAddressable(addressable bool) {
	// If addressable is true, set bit 12 of byte 2
	if addressable {
		h[3] |= 0x10
	} else {
		// Otherwise, clear bit 12 of byte 2
		h[3] &^= 0x10
	}
}

// SetTagged sets the tagged bit to a boolean value
func (h *Header) SetTagged(tagged bool) {
	// If tagged is true, set bit 13 of byte 2
	if tagged {
		h[3] |= 0x20
	} else {
		// Otherwise, clear bit 13 of byte 2
		h[3] &^= 0x20
	}
}

// SetSource sets the client source identifier
func (h *Header) SetSource(source uint32) {
	binary.LittleEndian.PutUint32(h[4:8], source)
}

// SetTarget sets the target MAC address from bytes
func (h *Header) SetTarget(mac []byte) {
	// Check if the length of the MAC address is 6 bytes
	if len(mac) != 6 {
		panic("MAC address must be 6 bytes")
	}
	// Copy the MAC address into bytes 8-13 of the header
	copy(h[8:14], mac)
}

// SetResponseRequired sets the response required bit to a boolean value
func (h *Header) SetResponseRequired(responseRequired bool) {
	// If responseRequired is true, set bit 0 of byte 15
	if responseRequired {
		h[15] |= 0x01
	} else {
		// Otherwise, clear bit 0 of byte 15
		h[15] &^= 0x01
	}
}

// SetAckRequired sets the ack required bit to a boolean value
func (h *Header) SetAckRequired(ackRequired bool) {
	// If ackRequired is true, set bit 1 of byte 15
	if ackRequired {
		h[15] |= 0x02
	} else {
		// Otherwise, clear bit 1 of byte 15
		h[15] &^= 0x02
	}
}

// SetSequence sets the sequence number
func (h *Header) SetSequence(sequence uint8) {
	// Shift the sequence number to the left by 2 bits and set it in byte 15
	h[15] = (h[15] & 0x03) | (sequence << 2)
}

// SetType sets the packet type
func (h *Header) SetType(packetType PacketType) {
	// Set the packet type in bytes 32-33
	binary.LittleEndian.PutUint16(h[32:34], uint16(packetType))
}

// PackHeader converts the header to a byte slice
func (h *Header) PackHeader() []byte {
	// Create a byte slice of the header size
	head := make([]byte, HeaderSize)
	// Copy the header bytes into the byte slice
	copy(head, h[:])
	// Return the byte slice
	return head
}

func DefaultHeader(source uint32, target []byte, packetType PacketType, payloadSize uint16) *Header {
	header := NewHeader(source)
	header.SetSize(uint16(HeaderSize + payloadSize))
	header.SetType(packetType)
	header.SetTarget(target)

	return header
}
