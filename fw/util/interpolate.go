package util

import (
	"golang.org/x/exp/constraints"
)

type LutPoint[K any, V any] struct {
	Key   K
	Value V
}

var Divider_47k_ntc3950 = [...]LutPoint[uint16, int16]{
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

func Interpolate[K constraints.Integer, V constraints.Integer | constraints.Float](key K, lut []LutPoint[K, V]) V {

	if key >= lut[0].Key {
		return lut[0].Value
	}
	last := lut[len(lut)-1]
	if key <= last.Key {
		return last.Value
	}

	lo, hi := 0, len(lut)-2
	for lo <= hi {
		mid := (lo + hi) / 2
		a0 := lut[mid].Key
		a1 := lut[mid+1].Key

		if key <= a0 && key >= a1 {
			t0 := lut[mid].Value
			t1 := lut[mid+1].Value
			da := a0 - a1
			x := a0 - key
			return t0 + V(x)*(t1-t0)/V(da)
		}
		if key > a0 {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return last.Value
}

func Compensate(min25, min75, max25, max75, temp, measured int32) int32 {

	pt := float32(temp-25) / float32(75-25)
	mint := float32(min25) + pt*float32(min75-min25)
	maxt := float32(max25) + pt*float32(max75-max25)
	pv := (float32(measured) - mint) / (maxt - mint)

	return int32(float32(min25) + pv*float32(max25-min25))
}
