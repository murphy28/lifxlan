package lifxlan

const (
	LifxPort   = 56700       // Default UDP port for LIFX LAN protocol
	HeaderSize = 8 + 16 + 12 // Frame Header + Frame Address + Protocol Header Size
	Protocol   = 1024        // LIFX Protocol Number
)
