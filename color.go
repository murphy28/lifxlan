package lifxlan

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

type LIFXColor struct {
	hue        uint16
	saturation uint16
	brightness uint16
	kelvin     uint16
}

func NewColor(hue, saturation, brightness, kelvin uint16) LIFXColor {
	return LIFXColor{
		hue:        hue,
		saturation: saturation,
		brightness: brightness,
		kelvin:     kelvin,
	}
}

func (c *LIFXColor) SetHue(hue float64) {
	c.hue = uint16(int((math.Round(0x10000*hue) / 360)) % 0x10000)
}

func (c *LIFXColor) GetHue() float64 {
	return float64(c.hue) * 360 / 0x10000
}

func (c *LIFXColor) SetSaturation(saturation float64) {
	c.saturation = uint16(math.Round(0xFFFF * saturation))
}

func (c *LIFXColor) GetSaturation() float64 {
	return float64(c.saturation) / 0xFFFF
}

func (c *LIFXColor) SetBrightness(brightness float64) {
	c.brightness = uint16(math.Round(0xFFFF * brightness))
}

func (c *LIFXColor) GetBrightness() float64 {
	return float64(c.brightness) / 0xFFFF
}

func (c *LIFXColor) SetKelvin(kelvin float64) {
	c.kelvin = uint16(math.Round(kelvin))
}

func (c *LIFXColor) GetKelvin() float64 {
	return float64(c.kelvin)
}

func ColorfulToLifxColor(color colorful.Color) LIFXColor {
	h, s, v := color.Hsv()

	lifxColor := NewColor(0, 0, 0, 3750)

	lifxColor.SetHue(h)
	lifxColor.SetSaturation(s)
	lifxColor.SetBrightness(v)

	return lifxColor
}

func NewColorFromHex(hex string) (LIFXColor, error) {
	c, err := colorful.Hex(hex)
	if err != nil {
		return NewColor(0, 0, 0, 3750), err
	}

	return ColorfulToLifxColor(c), nil
}

func NewColorFromRGB(r, g, b uint8) LIFXColor {
	c := colorful.Color{R: float64(r) / 255, G: float64(g) / 255, B: float64(b) / 255}

	return ColorfulToLifxColor(c)
}

func NewColorFromHSV(h, s, v float64) LIFXColor {
	c := colorful.Hsv(h, s, v)

	return ColorfulToLifxColor(c)
}
