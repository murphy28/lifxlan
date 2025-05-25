# LIFX LAN
`lifxlan` is simple Go library for interacting with LIFX devices via the [LIFX LAN Protocol](https://lan.developer.lifx.com/docs/introduction).

### Installation
``go get github.com/murphy28/lifxlan``

### Features
 - Discover LIFX devices over the LAN
 - Ping devices
 - Rename Devices
 - View device product info
 - Turn devices on and off
 - Set device color using RGB, HSV, Hex Colors, or by passing in a color from `go-colorful`

### Example
In the example below, a new LIFX client is created, device discovery runs for 5 seconds, all discovered devices are turned on, and the device named "Nightstand" is set to a purple color.
```go
package main

import (
	"fmt"
	"time"

	"github.com/murphy28/lifxlan"
)

func main() {

	// Create a new LIFX client
	lifx, err := lifxlan.NewClient()
	if err != nil {
		fmt.Println("Error creating LIFX client:", err)
		return
	}
	defer lifx.Close()

	// Discover LIFX devices on the network
	err = lifx.Discover(5 * time.Second)
	if err != nil {
		fmt.Println("Error discovering LIFX devices:", err)
		return
	}

	// Get the list of devices
	devices := lifx.GetDevices()

	// Turn on all devices
	for _, device := range devices {
		device.TurnOn()
	}

	// Set an individual device to a specific color
	lifx.SetColor("Nightstand", lifxlan.NewColorFromRGB(144, 48, 240), 5*time.Second)

}
```

### Supported Devices
This library was built with standard LIFX bulbs in mind. I do not own any of the fancier LIFX devices, but support will come if I get the chance. All code was tested on the LIFX Color.

While there are no concrete plans to support other LIFX devices, there are other libraries that may fit your needs.

 - [fishy/lifxlan](https://github.com/fishy/lifxlan)
 - [2tvenom/golifx](https://github.com/2tvenom/golifx)