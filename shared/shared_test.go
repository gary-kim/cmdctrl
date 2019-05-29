package shared

import (
	"testing"
)

func TestCompatable(t *testing.T) {
	// Test same version
	got := Compatible("v1.1.1", "v1.1.1")
	if !got {
		t.Errorf(`Compatable("v1.1.1", "v1.1.1") returned %t, was expecting true`, got)
	}

	// Test different patch version
	got = Compatible("v1.1.5", "v1.1.1")
	if !got {
		t.Errorf(`Compatable("v1.1.5", "v1.1.1") returned %t, was expecting true`, got)
	}

	// Test client minor ahead of server minor
	got = Compatible("v1.2.1", "v1.1.1")
	if !got {
		t.Errorf(`Compatable("v1.2.1", "v1.1.1") returned %t, was expecting true`, got)
	}

	// Test client minor behind of server minor
	got = Compatible("v1.1.1", "v1.2.1")
	if got {
		t.Errorf(`Compatable("v1.1.1", "v1.2.1") returned %t, was expecting false`, got)
	}

	// Test different major versions
	got = Compatible("v2.0.0", "v1.0.0")
	if got {
		t.Errorf(`Compatable("v2.0.0", "v1.0.0") returned %t, was expecting false`, got)
	}
}

func TestCompare(t *testing.T) {
	// Test with higher priority
	pa0 := PendingAction{
		Priority: 1,
	}
	pa1 := PendingAction{
		Priority: 5,
	}
	got := pa0.Compare(pa1)
	if !(got < 0) {
		t.Errorf(`pa0.Compare(pa1) returned %v, was expecting smaller than 0`, got)
	}

	// Test with lower priority
	got = pa1.Compare(pa0)
	if !(got > 0) {
		t.Errorf(`pa1.Compare(pa0) returned %v, was expecting larger than 0`, got)
	}

	// Test with same priority
	pa2 := PendingAction{
		Priority: 5,
	}
	got = pa1.Compare(pa2)
	if got != 0 {
		t.Errorf(`pa1.Compare(pa2) returned %v, was expecting 0`, got)
	}
}
