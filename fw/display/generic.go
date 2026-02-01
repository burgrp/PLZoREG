package display

// Initializes the display with given number of digits and intensity 0-255
func Init(digits, intensity int) {
	configurePeripherals(digits, intensity)
}

func DigitAt(digit int, dp bool, position int) {
	var numberGlyphs = []int{0x3F, 0x06, 0x5B, 0x4F, 0x66, 0x6D, 0x7D, 0x07, 0x7F, 0x6F}
	if digit >= 0 && digit <= 10 {
		g := numberGlyphs[digit]
		if dp {
			g = g | (1 << 7)
		}
		GlyphAt(g, position)
	}
}

func NumberAt(number int, leadingZero bool, dotPosition, basePosition, length int) {
	for p := basePosition; p >= basePosition-length+1; p-- {
		if !leadingZero && number == 0 && p < basePosition {
			GlyphAt(0, p)
		} else {
			DigitAt(number%10, p == dotPosition, p)
			number = number / 10
		}
	}
}
