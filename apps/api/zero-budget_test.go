package zerobudget

import (
	"testing"
)

func TestZeroBudget(t *testing.T) {
	result := ZeroBudget("works")
	if result != "ZeroBudget works" {
		t.Error("Expected ZeroBudget to append 'works'")
	}
}
