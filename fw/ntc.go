package main

type ntcPoint struct {
	adc  uint16 // ADC code
	tC10 int16  // temperature in 0.1°C
}

var ntc3950_47k_LUT = [...]ntcPoint{
	{3666, -400},
	{3509, -350},
	{3316, -300},
	{3089, -250},
	{2832, -200},
	{2554, -150},
	{2266, -100},
	{1981, -50},
	{1708, 0},
	{1456, 50},
	{1230, 100},
	{1032, 150},
	{862, 200},
	{718, 250},
	{598, 300},
	{498, 350},
	{415, 400},
	{347, 450},
	{290, 500},
	{244, 550},
	{206, 600},
	{174, 650},
	{148, 700},
	{126, 750},
	{108, 800},
	{93, 850},
	{80, 900},
	{69, 950},
	{60, 1000},
	{52, 1050},
	{46, 1100},
	{40, 1150},
	{35, 1200},
	{31, 1250},
}

// adc: 0..4095
// return: temperature in 0.1°C
func NTC3950TempC10(adc uint32) int32 {
	lut := ntc3950_47k_LUT[:]

	if adc >= uint32(lut[0].adc) {
		return int32(lut[0].tC10)
	}
	last := lut[len(lut)-1]
	if adc <= uint32(last.adc) {
		return int32(last.tC10)
	}

	lo, hi := 0, len(lut)-2
	for lo <= hi {
		mid := (lo + hi) / 2
		a0 := uint32(lut[mid].adc)
		a1 := uint32(lut[mid+1].adc)

		if adc <= a0 && adc >= a1 {
			t0 := int32(lut[mid].tC10)
			t1 := int32(lut[mid+1].tC10)
			da := int32(a0 - a1)
			x := int32(a0 - adc)
			return t0 + (t1-t0)*x/da
		}
		if adc > a0 {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return int32(last.tC10)
}
