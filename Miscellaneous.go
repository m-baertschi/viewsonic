package viewsonic

import (
	"encoding/binary"
	"fmt"
)

// High Altitude Mode
func (conn *ViewSonic) SetHighAltitudeMode(enable bool) error {
	value := uint8(0x00)
	if enable {
		value = 0x01
	}
	return conn.Write(0x110C, value) // PDF #16, 17
}

func (conn *ViewSonic) GetHighAltitudeMode() (bool, error) {
	value, err := conn.Read(0x110C) // PDF #18
	return value == 0x01, err
}

// Message
func (conn *ViewSonic) SetMessageDisplay(enable bool) error {
	value := uint8(0x00)
	if enable {
		value = 0x01
	}
	return conn.Write(0x1127, value) // PDF #24, 25
}

func (conn *ViewSonic) GetMessageDisplay() (bool, error) {
	value, err := conn.Read(0x1127) // PDF #26
	return value == 0x01, err
}

// Language
type Language uint8

const (
	LanguageEnglish     Language = 0x00
	LanguageFrench      Language = 0x01
	LanguageGerman      Language = 0x02
	LanguageItalian     Language = 0x03
	LanguageSpanish     Language = 0x04
	LanguageRussian     Language = 0x05
	LanguageTradChinese Language = 0x06
	LanguageSimpChinese Language = 0x07
	LanguageJapanese    Language = 0x08
	LanguageKorean      Language = 0x09
	LanguageSwedish     Language = 0x0A
	LanguageDutch       Language = 0x0B
	LanguageTurkish     Language = 0x0C
	LanguageCzech       Language = 0x0D
	LanguagePortuguese  Language = 0x0E
	LanguageThai        Language = 0x0F
	LanguagePolish      Language = 0x10
	LanguageFinnish     Language = 0x11
	LanguageArabic      Language = 0x12
	LanguageIndonesian  Language = 0x13
	LanguageHindi       Language = 0x14
	LanguageVietnamese  Language = 0x15
)

func (conn *ViewSonic) SetLanguage(lang Language) error {
	return conn.Write(0x1500, uint8(lang)) // PDF #141-162
}

func (conn *ViewSonic) GetLanguage() (Language, error) {
	val, err := conn.Read(0x1500) // PDF #163
	return Language(val), err
}

// Remote Control Code (1-8)
func (conn *ViewSonic) SetRemoteControlCode(code uint8) error {
	if code < 1 || code > 8 {
		return fmt.Errorf("remote code must be between 1 and 8")
	}
	return conn.Write(0x0C48, code-1) // PDF #190-197: code 1 is value 0
}
func (conn *ViewSonic) GetRemoteControlCode() (uint8, error) {
	val, err := conn.Read(0x0C48) // PDF #198
	return val + 1, err
}

// Remote Key
type RemoteKey uint8

const (
	RemoteKeyMenu     RemoteKey = 0x0F
	RemoteKeyExit     RemoteKey = 0x13
	RemoteKeyTop      RemoteKey = 0x0B
	RemoteKeyBottom   RemoteKey = 0x0C
	RemoteKeyLeft     RemoteKey = 0x0D
	RemoteKeyRight    RemoteKey = 0x0E
	RemoteKeySource   RemoteKey = 0x04
	RemoteKeyEnter    RemoteKey = 0x15
	RemoteKeyAuto     RemoteKey = 0x08
	RemoteKeyMyButton RemoteKey = 0x11
)

// 0x02 0x14 0x00 0x04 0x00 | 0x34 | 0x02 0x04 0x0F 0x61
func (conn *ViewSonic) SendRemoteKey(key RemoteKey) error {
	return conn.WriteKey(0x0204, uint8(key))
}

// Light Source
func (conn *ViewSonic) ResetLightSourceUsageTime() error {
	return conn.Write(0x1501, 0x00) // PDF #164
}

func (conn *ViewSonic) GetLightSourceUsageTime() (uint32, error) {
	data, err := conn.ReadNBytes(0x1501) // PDF #165
	if err != nil {
		return 0, err
	}
	if len(data) < 4 {
		return 0, fmt.Errorf("not enough data for usage time")
	}
	// Note 4: HEX2DEC(ddccbbaa)
	return binary.LittleEndian.Uint32(data), nil
}

type LightSourceMode uint8

const (
	LightSourceModeNormal     LightSourceMode = 0x00
	LightSourceModeEco        LightSourceMode = 0x01
	LightSourceModeDynamicEco LightSourceMode = 0x02
	LightSourceModeSuperEco   LightSourceMode = 0x03
)

func (conn *ViewSonic) SetLightSourceMode(mode LightSourceMode) error {
	return conn.Write(0x1110, uint8(mode)) // PDF #19-22
}

func (conn *ViewSonic) GetLightSourceMode() (LightSourceMode, error) {
	value, err := conn.Read(0x1110) // PDF #23
	if err != nil {
		return 0, err
	}
	return LightSourceMode(value), nil
}

func (conn *ViewSonic) CycleLampMode() error {
	return conn.Write(0x1336, 0x00) // PDF #224
}
