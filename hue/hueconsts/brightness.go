package hueconsts

const (
	MaxBrightness           uint8 = 254
	QuarterBrightness             = MaxBrightness / 4
	HalfBrightness                = MaxBrightness / 2
	ThreeQuartersBrightness       = QuarterBrightness + HalfBrightness
	MinBrightness           uint8 = 1
)

func CalculateMaxBrightnessDelta(delta uint8) uint8 {
	return MaxBrightness / delta
}

func CalculateMaxBrightnessDeltaMul(delta uint8, mul uint8) uint8 {
	return (MaxBrightness / delta) * mul
}

func CalculateBrightness(b1 uint8, offset uint8) uint8 {
	return b1 + offset
}