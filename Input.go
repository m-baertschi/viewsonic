package viewsonic

// Source input
type SourceInput uint8

const (
	SourceInputDSub1      SourceInput = 0x00
	SourceInputDSub2      SourceInput = 0x08
	SourceInputHDMI1      SourceInput = 0x03
	SourceInputHDMI2      SourceInput = 0x07
	SourceInputHDMI3      SourceInput = 0x09
	SourceInputHDMIMHL4   SourceInput = 0x0E
	SourceInputComposite  SourceInput = 0x05
	SourceInputSVideo     SourceInput = 0x06
	SourceInputDVI        SourceInput = 0x0A
	SourceInputComponent  SourceInput = 0x0B
	SourceInputHDBaseT    SourceInput = 0x0C
	SourceInputUSBC       SourceInput = 0x0F
	SourceInputUSBReader  SourceInput = 0x1A
	SourceInputLANWiFi    SourceInput = 0x1B
	SourceInputUSBDisplay SourceInput = 0x1C
)

func (conn *ViewSonic) SetSourceInput(input SourceInput) error {
	return conn.Write(0x1301, uint8(input)) // PDF #115-129
}

func (conn *ViewSonic) GetSourceInput() (SourceInput, error) {
	value, err := conn.Read(0x1301) // PDF #130
	if err != nil {
		return 0, err
	}
	return SourceInput(value), nil
}

// Quick Auto Search
func (conn *ViewSonic) SetQuickAutoSearch(enable bool) error {
	value := uint8(0x00)
	if enable {
		value = 0x01
	}
	return conn.Write(0x1302, value) // PDF #131, 132
}

func (conn *ViewSonic) GetQuickAutoSearch() (bool, error) {
	val, err := conn.Read(0x1302) // PDF #133
	return val == 0x01, err
}

// HDMI Format
type HdmiFormat uint8

const (
	HdmiFormatRGB  HdmiFormat = 0x00
	HdmiFormatYUV  HdmiFormat = 0x01
	HdmiFormatAuto HdmiFormat = 0x02
)

func (conn *ViewSonic) SetHdmiFormat(format HdmiFormat) error {
	return conn.Write(0x1128, uint8(format)) // PDF #166-168
}
func (conn *ViewSonic) GetHdmiFormat() (HdmiFormat, error) {
	val, err := conn.Read(0x1128) // PDF #169
	return HdmiFormat(val), err
}

// Hdmi Range
type HdmiRange uint8

const (
	HdmiRangeEnhanced HdmiRange = 0x00 // 0-255
	HdmiRangeNormal   HdmiRange = 0x01 // 16-235
	HdmiRangeAuto     HdmiRange = 0x02
)

func (conn *ViewSonic) SetHdmiRange(r HdmiRange) error {
	return conn.Write(0x1129, uint8(r)) // PDF #170-172
}
func (conn *ViewSonic) GetHdmiRange() (HdmiRange, error) {
	val, err := conn.Read(0x1129) // PDF #173
	return HdmiRange(val), err
}

// HDMI CEC
func (conn *ViewSonic) SetCEC(enable bool) error {
	value := uint8(0x00)
	if enable {
		value = 0x01
	}
	return conn.Write(0x112B, value) // PDF #174, 175
}
func (conn *ViewSonic) GetCEC() (bool, error) {
	val, err := conn.Read(0x112B) // PDF #176
	return val == 0x01, err
}

// Image Position
func (conn *ViewSonic) ShiftHorizontalPosition(right bool) error {
	value := uint8(0x00) // Left
	if right {
		value = 0x01
	}
	return conn.Write(0x1206, value) // PDF #60, 61
}

func (conn *ViewSonic) GetHorizontalPosition() (uint8, error) {
	return conn.Read(0x1206) // PDF #62
}

func (conn *ViewSonic) ShiftVerticalPosition(up bool) error {
	value := uint8(0x00) // Up
	if !up {             // down
		value = 0x01
	}
	return conn.Write(0x1207, value) // PDF #63, 64
}

func (conn *ViewSonic) GetVerticalPosition() (uint8, error) {
	return conn.Read(0x1207) // PDF #65
}

// Keystone
func (conn *ViewSonic) IncreaseKeystoneVertical() error {
	return conn.Write(0x120A, 0x01) // PDF #75
}

func (conn *ViewSonic) DecreaseKeystoneVertical() error {
	return conn.Write(0x120A, 0x00) // PDF #74
}

func (conn *ViewSonic) GetKeystoneVertical() (uint8, error) {
	return conn.Read(0x120A) // PDF #76
}

func (conn *ViewSonic) IncreaseKeystoneHorizontal() error {
	return conn.Write(0x1131, 0x01) // PDF #78
}

func (conn *ViewSonic) DecreaseKeystoneHorizontal() error {
	return conn.Write(0x1131, 0x00) // PDF #77
}

func (conn *ViewSonic) GetKeystoneHorizontal() (uint8, error) {
	return conn.Read(0x1131) // PDF #79
}
