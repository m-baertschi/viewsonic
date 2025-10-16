package viewsonic

// Mute
func (conn *ViewSonic) SetMute(mute bool) error {
	value := int8(0x00) // OFF
	if mute {
		value = 0x01 // ON
	}
	return conn.Write(0x1400, value) // PDF #134, 135
}

func (conn *ViewSonic) GetMute() (bool, error) {
	value, err := conn.Read(0x1400) // PDF #136
	if err != nil {
		return false, err
	}
	return value == 0x01, nil
}

// Volume
func (conn *ViewSonic) IncreaseVolume() error {
	return conn.Write(0x1401, 0x00) // PDF #137
}

func (conn *ViewSonic) DecreaseVolume() error {
	return conn.Write(0x1402, 0x00) // PDF #138
}

func (conn *ViewSonic) SetVolume(level int8) error {
	return conn.Write(0x132A, level) // PDF #139
}

func (conn *ViewSonic) GetVolume() (int8, error) {
	return conn.Read(0x1403) // PDF #140
}

// Audio mode cycle
func (conn *ViewSonic) CycleAudioMode() error {
	return conn.Write(0x1335, 0x00) // PDF #225
}
