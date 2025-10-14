package viewsonic

// Color temperature
type ColorTemperature uint8

const (
	ColorTemperatureWarm    ColorTemperature = 0x00
	ColorTemperatureNormal  ColorTemperature = 0x01
	ColorTemperatureNeutral ColorTemperature = 0x02
	ColorTemperatureCool    ColorTemperature = 0x03
)

func (conn *ViewSonic) SetColorTemperature(temp ColorTemperature) error {
	return conn.Write(0x1208, uint8(temp)) // PDF #66-69
}

func (conn *ViewSonic) GetColorTemperature() (ColorTemperature, error) {
	val, err := conn.Read(0x1208) // PDF #70
	return ColorTemperature(val), err
}

// Color mode
type ColorMode uint8

const (
	ColorModeBrightest     ColorMode = 0x00
	ColorModeMovie         ColorMode = 0x01
	ColorModeStandard      ColorMode = 0x04
	ColorModeSRGBViewMatch ColorMode = 0x05
	ColorModeDynamic       ColorMode = 0x08
	ColorModeRec709        ColorMode = 0x09
	ColorModeDICOMSIM      ColorMode = 0x0A
	ColorModeSports        ColorMode = 0x11
	ColorModePhoto         ColorMode = 0x13
	ColorModePresentation  ColorMode = 0x14
	ColorModeGaming        ColorMode = 0x12
	ColorModeVivid         ColorMode = 0x15
	ColorModeISFDay        ColorMode = 0x16
	ColorModeISFNight      ColorMode = 0x17
)

func (conn *ViewSonic) SetColorMode(mode ColorMode) error {
	return conn.Write(0x120B, uint8(mode)) // PDF #80-90
}

func (conn *ViewSonic) GetColorMode() (ColorMode, error) {
	val, err := conn.Read(0x120B) // PDF #92
	return ColorMode(val), err
}

func (conn *ViewSonic) CycleColorMode() error {
	return conn.Write(0x1333, 0x00) // PDF #91
}

// Primary Color
type PrimaryColor uint8

const (
	PrimaryColorR PrimaryColor = 0x00
	PrimaryColorG PrimaryColor = 0x01
	PrimaryColorB PrimaryColor = 0x02
	PrimaryColorC PrimaryColor = 0x03
	PrimaryColorM PrimaryColor = 0x04
	PrimaryColorY PrimaryColor = 0x05
)

func (conn *ViewSonic) SelectPrimaryColor(color PrimaryColor) error {
	return conn.Write(0x1210, uint8(color)) // PDF #93-98
}

func (conn *ViewSonic) GetSelectedPrimaryColor() (PrimaryColor, error) {
	val, err := conn.Read(0x1210) // PDF #99
	return PrimaryColor(val), err
}

// Hue/Tint
func (conn *ViewSonic) IncreaseHue() error {
	return conn.Write(0x1211, 0x01) // PDF #101
}
func (conn *ViewSonic) DecreaseHue() error {
	return conn.Write(0x1211, 0x00) // PDF #100
}
func (conn *ViewSonic) GetHue() (int16, error) {
	return conn.Read2Bytes(0x1211) // PDF #102
}

// Saturation
func (conn *ViewSonic) IncreaseSaturation() error {
	return conn.Write(0x1212, 0x01) // PDF #104
}
func (conn *ViewSonic) DecreaseSaturation() error {
	return conn.Write(0x1212, 0x00) // PDF #103
}
func (conn *ViewSonic) GetSaturation() (int16, error) {
	return conn.Read2Bytes(0x1212) // PDF #105
}

// Sharpness
func (conn *ViewSonic) IncreaseSharpness() error {
	return conn.Write(0x120E, 0x01) // PDF #110
}
func (conn *ViewSonic) DecreaseSharpness() error {
	return conn.Write(0x120E, 0x00) // PDF #109
}
func (conn *ViewSonic) GetSharpness() (int16, error) {
	return conn.Read2Bytes(0x120E) // PDF #111
}

// Gain
func (conn *ViewSonic) IncreaseGain() error {
	return conn.Write(0x1213, 0x01) // PDF #107
}
func (conn *ViewSonic) DecreaseGain() error {
	return conn.Write(0x1213, 0x00) // PDF #106
}
func (conn *ViewSonic) GetGain() (int16, error) {
	return conn.Read2Bytes(0x1213) // PDF #108
}

// Brilliant Color
// SetBrilliantColor sets the brilliant color level (0-10)
func (conn *ViewSonic) SetBrilliantColor(level uint8) error {
	if level > 10 {
		level = 10
	}
	// PDF #178-188: Brilliant Color OFF is value 0, Color 1 is value 1, etc.
	return conn.Write(0x120F, level)
}
func (conn *ViewSonic) GetBrilliantColor() (uint8, error) {
	return conn.Read(0x120F) // PDF #189
}

// Screen Color
type ScreenColor uint8

const (
	ScreenColorOff        ScreenColor = 0x00
	ScreenColorBlackboard ScreenColor = 0x01
	ScreenColorGreenboard ScreenColor = 0x02
	ScreenColorWhiteboard ScreenColor = 0x03
	ScreenColorBlueboard  ScreenColor = 0x04
)

func (conn *ViewSonic) SetScreenColor(color ScreenColor) error {
	return conn.Write(0x1132, uint8(color)) // PDF #199-203
}
func (conn *ViewSonic) GetScreenColor() (ScreenColor, error) {
	val, err := conn.Read(0x1132) // PDF #204
	return ScreenColor(val), err
}
