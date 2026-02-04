package util

import "testing"

func TestCompensate(t *testing.T) {
	min25 := int32(100)
	min75 := int32(75)
	max25 := int32(250)
	max75 := int32(193)

	if compensated := Compensate(min25, min75, max25, max75, 25, min25); compensated != min25 {
		t.Errorf("compensate at 25C: got %d, want %d", compensated, min25)
	}
	if compensated := Compensate(min25, min75, max25, max75, 75, min75); compensated != min25 {
		t.Errorf("compensate at 75C: got %d, want %d", compensated, min25)
	}
	if compensated := Compensate(min25, min75, max25, max75, 25, max25); compensated != max25 {
		t.Errorf("compensate at 25C: got %d, want %d", compensated, max25)
	}
	if compensated := Compensate(min25, min75, max25, max75, 75, max75); compensated != max25 {
		t.Errorf("compensate at 75C: got %d, want %d", compensated, max25)
	}

}
