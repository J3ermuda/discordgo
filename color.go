package discordgo

import (
	"strconv"
	"strings"
)

const (
	ColorGreen     = Color(0x2ecc71)
	ColorBlue      = Color(0x3498db)
	ColorPurple    = Color(0x9b59b6)
	ColorGold      = Color(0xf1c40f)
	ColorRed       = Color(0xe74c3c)
	ColorOrange    = Color(0xe67e22)
	ColorMagenta   = Color(0xe91e63)
	ColorTeal      = Color(0x1abc9c)
	ColorLightGrey = Color(0x979c9f)
	ColorDarkGrey  = Color(0x607d8b)
	ColorBlurple   = Color(0x7289da)
	ColorGreyple   = Color(0x99aab5)
)

// Color is a type around the int value of a discord color
type Color int

// DefaultColor returns the Color object with a value of 0,
// the default color of discord roles
func DefaultColor() Color {
	return Color(0)
}

func (c *Color) getByte(b uint) int {
	return int((*c >> (8 * b)) & 0xFF)
}

// String formats the Color into a string of the hexcode representing it
func (c Color) String() string {
	return strconv.FormatInt(int64(c), 16)
}

// B returns the blue component of the Color
func (c *Color) B() int {
	return c.getByte(0)
}

// G returns the green component of the Color
func (c *Color) G() int {
	return c.getByte(1)
}

// R returns the red component of the Color
func (c *Color) R() int {
	return c.getByte(2)
}

// SetRGB sets the value of the color to the given RGB values
// r :  the red component
// g :  the green component
// b :  the blue component
func (c *Color) SetRGB(r, g, b int) {
	*c = Color((r << 16) + (g << 8) + b)
}

// SetHex sets the value of the color to the given hexcode string
// hex :  the string containing the hex code of the color
func (c *Color) SetHex(hex string) (err error) {
	hex = strings.Trim(hex, "#")
	i, err := strconv.ParseInt(hex, 16, 0)
	if err != nil {
		return
	}
	*c = Color(i)
	return
}
