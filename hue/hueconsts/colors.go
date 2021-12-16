package hueconsts

type Color []float32

func makeColor(x, y float32) Color {
	return []float32{x, y}
}

var (
	Daylight   = makeColor(0.312, 0.3281)
	Nightlight = makeColor(0.5612, 0.4042)
)