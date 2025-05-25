package lifxlan

type PacketType uint16

// Query Packet types (https://lan.developer.lifx.com/docs/querying-the-device-for-data)

const (
	GetService PacketType = 2 // Discovery Request

	GetHostFirmware PacketType = 14
	GetWifiInfo     PacketType = 16
	GetWifiFirmware PacketType = 18
	GetPower        PacketType = 20
	GetLabel        PacketType = 23
	GetVersion      PacketType = 32
	GetInfo         PacketType = 34
	GetLocation     PacketType = 48
	GetGroup        PacketType = 51
	EchoRequest     PacketType = 58

	GetColor                 PacketType = 101
	GetLightPower            PacketType = 116
	GetInfrared              PacketType = 120
	GetHevCycle              PacketType = 142
	GetHevCycleConfiguration PacketType = 145
	GetLastHevCycleResult    PacketType = 148

	GetColorZones         PacketType = 502
	GetMultiZoneEffect    PacketType = 507
	GetExtendedColorZones PacketType = 511

	GetRPower PacketType = 816

	GetDeviceChain        PacketType = 701
	Get64                 PacketType = 707
	GetTileEffect         PacketType = 718
	SensorGetAmbientLight PacketType = 401
)

// Set Packet types (https://lan.developer.lifx.com/docs/changing-a-device)

const (
	SetPower    PacketType = 21
	SetLabel    PacketType = 24
	SetReboot   PacketType = 38
	SetLocation PacketType = 49
	SetGroup    PacketType = 52

	SetColor                 PacketType = 102
	SetWaveform              PacketType = 103
	SetLightPower            PacketType = 117
	SetWafeformOptional      PacketType = 119
	SetInfrared              PacketType = 122
	SetHevCycle              PacketType = 143
	SetHevCycleConfiguration PacketType = 146

	SetColorZones         PacketType = 501
	SetMultiZoneEffect    PacketType = 508
	SetExtendedColorZones PacketType = 510

	SetRPower PacketType = 817

	SetUserPosition PacketType = 703
	Set64           PacketType = 715
	SetTileEffect   PacketType = 719
)

// Information Packet types (https://lan.developer.lifx.com/docs/information-messages)

const (
	Acknowledgement PacketType = 45 // Returned when `ack_required=1` is specified.

	StateService PacketType = 3 // Discovery Response

	StateHostFirmware PacketType = 15
	StateWifiInfo     PacketType = 17
	StateWifiFirmware PacketType = 19
	StatePower        PacketType = 22
	StateLabel        PacketType = 25
	StateVersion      PacketType = 33
	StateInfo         PacketType = 35
	StateLocation     PacketType = 50
	StateGroup        PacketType = 53
	EchoResponse      PacketType = 59
	StateUnhandled    PacketType = 223

	LightState                 PacketType = 107
	StateLightPower            PacketType = 118
	StateInfrared              PacketType = 121
	StateHevCycle              PacketType = 144
	StateHevCycleConfiguration PacketType = 147
	StateLastHevCycleResult    PacketType = 149

	StateZone               PacketType = 503
	StateMultiZone          PacketType = 506
	StateMultiZoneEffect    PacketType = 509
	StateExtendedColorZones PacketType = 512

	StateRPower PacketType = 818

	StateDeviceChain  PacketType = 702
	State64           PacketType = 711
	StateTileEffect   PacketType = 720
	StateAmbientLight PacketType = 402
)
