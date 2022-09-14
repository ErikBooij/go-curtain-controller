package drivers

type DeviceList struct {
	AqaraShutters AqaraShutters
	SlideCurtains SlideCurtains
}

type AqaraShutters map[string]AqaraShutter
type SlideCurtains map[string]SlideCurtain
