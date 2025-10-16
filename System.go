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
// GetErrorStatus returns the raw error status bytes. See PDF Note 3 for decoding.
func (conn *ViewSonic) GetErrorStatus() ([]byte, error) {
	return conn.ReadNBytes(0x0D0D) // PDF #177
}

// Temperature status
func (conn *ViewSonic) GetOperatingTemperature() (float32, error) {
	data, err := conn.ReadNBytes(0x1503) // PDF #223
	if err != nil {
		return 0, err
	}
	if len(data) < 4 {
		return 0, fmt.Errorf("not enough data for temperature")
	}
	// Note 1: HEX2DEC(ddccbbaa)/10
	val := binary.LittleEndian.Uint32(data)
	return float32(val) / 10.0, nil
}
