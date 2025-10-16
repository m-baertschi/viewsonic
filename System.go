package viewsonic

import (
	"encoding/binary"
	"fmt"
)

// Power and Status
type PowerState uint8

const (
	PowerStateOn  PowerState = 0x01
	PowerStateOff PowerState = 0x00
)

type ProjectorStatusValue uint8

const (
	ProjectorStatusPowerOff ProjectorStatusValue = 0x00
	ProjectorStatusWarmUp   ProjectorStatusValue = 0x01
	ProjectorStatusPowerOn  ProjectorStatusValue = 0x02
	ProjectorStatusCoolDown ProjectorStatusValue = 0x03
)

func (conn *ViewSonic) SetPower(state PowerState) error {
	if state == PowerStateOn {
		return conn.Write(0x1100, 0x00) // PDF #1
	}
	return conn.Write(0x1101, 0x00) // PDF #2
}

func (conn *ViewSonic) GetPower() (PowerState, error) {
	value, err := conn.Read(0x1100) // PDF #3
	if err != nil {
		return 0, err
	}
	return PowerState(value), nil
}

func (conn *ViewSonic) GetProjectorStatus() (ProjectorStatusValue, error) {
	value, err := conn.Read(0x1126) // PDF #4
	if err != nil {
		return 0, err
	}
	return ProjectorStatusValue(value), nil
}

// Reset All Settings
func (conn *ViewSonic) ResetAllSettings() error {
	return conn.Write(0x1102, 0x00) // PDF #5
}

func (conn *ViewSonic) ResetCurrentColorSettings() error {
	return conn.Write(0x112A, 0x00) // PDF #6
}

// Quick Power Off
func (conn *ViewSonic) SetQuickPowerOff(enable bool) error {
	value := uint8(0x00)
	if enable {
		value = 0x01
	}
	return conn.Write(0x110B, value) // PDF #13, 14
}

func (conn *ViewSonic) GetQuickPowerOff() (bool, error) {
	value, err := conn.Read(0x110B) // PDF #15
	return value == 0x01, err
}

// Error status
type LampModeStatus uint8

const (
	LampModeStatusStandby                    LampModeStatus = 0x00
	LampModeStatusIgnition                   LampModeStatus = 0x01 // 0x01, 0x02, 0x03
	LampModeStatusLampRunUp                  LampModeStatus = 0x04
	LampModeStatusCoolDown                   LampModeStatus = 0x05
	LampModeStatusNormalLampOperation        LampModeStatus = 0x06
	LampModeStatusShutdownUnrecoverableError LampModeStatus = 0x08
	LampModeStatusPreHeatingPhase            LampModeStatus = 0x09 // 0x09, 0x0C
)

type LampModeErrorStatus uint16

const (
	LampModeErrorStatusNoError                         LampModeErrorStatus = 0x00
	LampModeErrorStatusTemperatureShutdown             LampModeErrorStatus = 0x01
	LampModeErrorStatusShortCircuit                    LampModeErrorStatus = 0x02
	LampModeErrorStatusEndOfLampLife                   LampModeErrorStatus = 0x03
	LampModeErrorStatusLampDidNotIgnite                LampModeErrorStatus = 0x04
	LampModeErrorStatusLampExtinguishedDuringOperation LampModeErrorStatus = 0x05
	LampModeErrorStatusLampExtinguishedDuringRunUp     LampModeErrorStatus = 0x06
	LampModeErrorStatusEEPROMWriteError                LampModeErrorStatus = 0x07
	LampModeErrorStatusEEPROMWriteBufferOverflow       LampModeErrorStatus = 0x08
	LampModeErrorStatusUARTBufferOverflow              LampModeErrorStatus = 0x09
	LampModeErrorStatusLampCurrentCalculationError     LampModeErrorStatus = 0x0A
	LampModeErrorStatusCorruptedSoftwareConfiguration  LampModeErrorStatus = 0x0B
	LampModeErrorStatusLampVoltageTooLow               LampModeErrorStatus = 0x0C
	LampModeErrorStatusEEPROMConfigMismatch            LampModeErrorStatus = 0x0F
	LampModeErrorStatusMaxPreHeatingTimeElapsed        LampModeErrorStatus = 0x10
)

type ErrorStatus struct {
	LampFailCount               uint8
	LampLitErrorCount           uint8
	Fan1ErrorCount              uint8
	Fan2ErrorCount              uint8
	Fan3ErrorCount              uint8
	Fan4ErrorCount              uint8
	Diode1OpenErrorCount        uint8
	Diode2OpenErrorCount        uint8
	Diode1ShortErrorCount       uint8
	Diode2ShortErrorCount       uint8
	TemperatureErrorCount       uint8
	Temperature2ErrorCount      uint8
	FanIC1ErrorCount            uint8
	ColorWheelErrorCount        uint8
	ColorWheelStartupErrorCount uint8
	UART1ErrorCount             uint8
	AbnormalPowerdown           uint8
	FirstBurnInErrorMinute      uint32
	LampStatus                  LampModeStatus
	LampErrorStatus             LampModeErrorStatus
}

// GetErrorStatus returns the decoded error status.
func (conn *ViewSonic) GetErrorStatus() (*ErrorStatus, error) {
	data, err := conn.ReadNBytes(0x0C0D) // PDF #177
	if err != nil {
		return nil, err
	}

	if len(data) < 24 {
		return nil, fmt.Errorf("not enough data for error status, expected 24 bytes, got %d", len(data))
	}

	lampModeStatus := LampModeStatus(data[21])
	switch lampModeStatus {
	case 0x01, 0x02, 0x03:
		lampModeStatus = LampModeStatusIgnition
	case 0x09, 0x0C:
		lampModeStatus = LampModeStatusPreHeatingPhase
	}

	status := &ErrorStatus{
		LampFailCount:               data[0],
		LampLitErrorCount:           data[1],
		Fan1ErrorCount:              data[2],
		Fan2ErrorCount:              data[3],
		Fan3ErrorCount:              data[4],
		Fan4ErrorCount:              data[5],
		Diode1OpenErrorCount:        data[6],
		Diode2OpenErrorCount:        data[7],
		Diode1ShortErrorCount:       data[8],
		Diode2ShortErrorCount:       data[9],
		TemperatureErrorCount:       data[10],
		Temperature2ErrorCount:      data[11],
		FanIC1ErrorCount:            data[12],
		ColorWheelErrorCount:        data[13],
		ColorWheelStartupErrorCount: data[14],
		UART1ErrorCount:             data[15],
		AbnormalPowerdown:           data[16],
		FirstBurnInErrorMinute:      binary.LittleEndian.Uint32(data[17:21]),
		LampStatus:                  lampModeStatus,
		LampErrorStatus:             LampModeErrorStatus(binary.LittleEndian.Uint16(data[22:24])),
	}

	return status, nil
}

// Temperature status
func (conn *ViewSonic) GetOperatingTemperature() (float32, float32, error) {
	data, err := conn.ReadNBytes(0x1503) // PDF #223
	if err != nil {
		return 0, 0, err
	}
	if len(data) < 4 {
		return 0, 0, fmt.Errorf("not enough data for temperature")
	}
	// Note 1: HEX2DEC(ddccbbaa)/10
	val := binary.LittleEndian.Uint32(data[2:])
	val2 := binary.LittleEndian.Uint32(data[6:])
	return float32(val) / 10.0, float32(val2) / 10.0, nil
}
