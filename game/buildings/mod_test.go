package buildings

import "testing"

func TestModuloOfNegativeNumbers(t *testing.T) {
	if m := -6 % 4; m != -2 {
		t.Error("-6 % 4 =", m)
	}
	if m := -5 % 4; m != -1 {
		t.Error("-6 % 5 =", m)
	}
	if m := -4 % 4; m != 0 {
		t.Error("-6 % 4 =", m)
	}
}

func TestNormalizedRotationIsInRange0To359(t *testing.T) {
	checkInt(t, normalizeRotation(0), 0)
	checkInt(t, normalizeRotation(359), 359)
	checkInt(t, normalizeRotation(360), 0)
	checkInt(t, normalizeRotation(-1), 359)
	checkInt(t, normalizeRotation(-360), 0)
	checkInt(t, normalizeRotation(-359), 1)
	checkInt(t, normalizeRotation(-359-360), 1)
	checkInt(t, normalizeRotation(15+3*360), 15)
}

func checkInt(t *testing.T, actual, expected int) {
	if actual != expected {
		t.Error("expected", expected, "but was", actual)
	}
}
