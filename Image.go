package viewsonic

// Splash Screen
type SplashScreen int8

const (
	SplashScreenBlack     SplashScreen = 0x00
	SplashScreenBlue      SplashScreen = 0x01
	SplashScreenViewSonic SplashScreen = 0x02
	SplashScreenCapture   SplashScreen = 0x03
	SplashScreenOff       SplashScreen = 0x04
)

func (conn *ViewSonic) SetSplashScreen(screen SplashScreen) error {
	return conn.Write(0x110A, int8(screen))
}

func (conn *ViewSonic) GetSplashScreen() (SplashScreen, error) {
	value, err := conn.Read(0x110A) // PDF #12
	if err != nil {
		return 0, err
	}
	return SplashScreen(value), nil
}

// Projector Position
type ProjectorPosition int8

const (
	ProjectorPositionFrontTable   ProjectorPosition = 0x00
	ProjectorPositionRearTable    ProjectorPosition = 0x01
	ProjectorPositionRearCeiling  ProjectorPosition = 0x02
	ProjectorPositionFrontCeiling ProjectorPosition = 0x03
)

func (conn *ViewSonic) SetProjectorPosition(pos ProjectorPosition) error {
	return conn.Write(0x1200, int8(pos)) // PDF #27-30
}

func (conn *ViewSonic) GetProjectorPosition() (ProjectorPosition, error) {
	value, err := conn.Read(0x1200) // PDF #31
	if err != nil {
		return 0, err
	}
	return ProjectorPosition(value), nil
}

// Contrast
func (conn *ViewSonic) IncreaseContrast() error {
	return conn.Write(0x1202, 0x01) // PDF #43
}

func (conn *ViewSonic) DecreaseContrast() error {
	return conn.Write(0x1202, 0x00) // PDF #42
}

func (conn *ViewSonic) GetContrast() (int16, error) {
	return conn.Read2Bytes(0x1202) // PDF #44
}

// Brightness
func (conn *ViewSonic) IncreaseBrightness() error {
	return conn.Write(0x1203, 0x01) // PDF #46
}

func (conn *ViewSonic) DecreaseBrightness() error {
	return conn.Write(0x1203, 0x00) // PDF #45
}

func (conn *ViewSonic) GetBrightness() (int16, error) {
	return conn.Read2Bytes(0x1203) // PDF #47
}

// Aspect ratio
type AspectRatio int8

const (
	AspectRatioAuto       AspectRatio = 0x00
	AspectRatio4To3       AspectRatio = 0x02
	AspectRatio16To9      AspectRatio = 0x03
	AspectRatio16To10     AspectRatio = 0x04
	AspectRatioAnamorphic AspectRatio = 0x05
	AspectRatioWide       AspectRatio = 0x06
	AspectRatio235To1     AspectRatio = 0x07
	AspectRatioPanorama   AspectRatio = 0x08
	AspectRatioNative     AspectRatio = 0x09
)

func (conn *ViewSonic) SetAspectRatio(ratio AspectRatio) error {
	return conn.Write(0x1204, int8(ratio)) // PDF #48-56
}

func (conn *ViewSonic) GetAspectRatio() (AspectRatio, error) {
	value, err := conn.Read(0x1204) // PDF #58
	if err != nil {
		return 0, err
	}
	return AspectRatio(value), nil
}

func (conn *ViewSonic) CycleAspectRatio() error {
	return conn.Write(0x1331, 0x00) // PDF #57
}

// Auto Adjust
func (conn *ViewSonic) AutoAdjust() error {
	return conn.Write(0x1205, 0x00) // PDF #59
}

// Blank
func (conn *ViewSonic) SetBlank(blank bool) error {
	value := int8(0x00) // Off
	if blank {
		value = 0x01 // On
	}
	return conn.Write(0x1209, value) // PDF #71, 72
}

func (conn *ViewSonic) GetBlank() (bool, error) {
	val, err := conn.Read(0x1209) // PDF #73
	return val == 0x01, err
}

// Freeze
func (conn *ViewSonic) SetFreeze(freeze bool) error {
	value := int8(0x00)
	if freeze {
		value = 0x01
	}
	return conn.Write(0x1300, value) // PDF #112, 113
}

func (conn *ViewSonic) GetFreeze() (bool, error) {
	val, err := conn.Read(0x1300) // PDF #114
	return val == 0x01, err
}

// Over Scan (0-5)
func (conn *ViewSonic) SetOverScan(value int8) error {
	if value > 5 {
		value = 5
	}
	return conn.Write(0x1133, value) // PDF #205-210
}
func (conn *ViewSonic) GetOverScan() (int8, error) {
	return conn.Read(0x1133) // PDF #211
}

// 3D Sync Mode
type ThreeDSyncMode int8

const (
	ThreeDSyncOff             ThreeDSyncMode = 0x00
	ThreeDSyncAuto            ThreeDSyncMode = 0x01
	ThreeDSyncFrameSequential ThreeDSyncMode = 0x02
	ThreeDSyncFramePacking    ThreeDSyncMode = 0x03
	ThreeDSyncTopBottom       ThreeDSyncMode = 0x04
	ThreeDSyncSideBySide      ThreeDSyncMode = 0x05
)

func (conn *ViewSonic) SetThreeDSyncMode(mode ThreeDSyncMode) error {
	return conn.Write(0x1220, int8(mode)) // PDF #32-37
}

func (conn *ViewSonic) GetThreeDSyncMode() (ThreeDSyncMode, error) {
	value, err := conn.Read(0x1220) // PDF #38
	if err != nil {
		return 0, err
	}
	return ThreeDSyncMode(value), nil
}

func (conn *ViewSonic) SetThreeDSyncInvert(enable bool) error {
	value := int8(0x00)
	if enable {
		value = 0x01
	}
	return conn.Write(0x1221, value) // PDF #39, 40
}

func (conn *ViewSonic) GetThreeDSyncInvert() (bool, error) {
	value, err := conn.Read(0x1221) // PDF #41
	return value == 0x01, err
}
