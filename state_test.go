package zorya

import (
	"math"
	"runtime/debug"
	"testing"
)

func TestRoundㆍNormal(t *testing.T) {
	before := []float64{-2.7, -2.5, -2.3, -1, 0, 0.1, 0.5, 0.8, 1, 1.5, 1.7, 2}
	after := []float64{-3, -3, -2, -1, 0, 0, 1, 1, 1, 2, 2, 2}
	for i, x := range before {
		r := round(x)
		if r != after[i] {
			t.Errorf("round(%f) = %f; != %f", x, round(x), after[i])
		}
	}
}

func TestRoundㆍSpecial(t *testing.T) {
	r := round(math.NaN())
	if !math.IsNaN(r) {
		t.Errorf("round(NaN) != NaN (%v)", r)
	}
	r = round(math.Inf(+1))
	if !math.IsInf(r, +1) {
		t.Errorf("round(+Inf) != +Inf (%v)", r)
	}
	r = round(math.Inf(-1))
	if !math.IsInf(r, -1) {
		t.Errorf("round(-Inf) != -Inf (%v)", r)
	}
}

func expectPanic(t *testing.T, code int) {
	if err := recover(); err != nil {
		if e, ok := err.(*Error); ok && e.Code == code {
			t.Logf("Received expected Zorya error: [%d] %v", e.Code, e.Msg)
			return
		} else if ok {
			t.Errorf("Unexpected Zorya error: [%d] %v", e.Code, e.Msg)
			return
		}
		debug.PrintStack()
		t.Error("Unexpected panic:", err)
	} else {
		t.Error("Expected error, got none.")
	}
}

func TestThreadㆍBadOpcode(t *testing.T) {
	th := Thread{reg: make([]float64, 256)}
	err := th.exec(-5, nil, 0)
	if !IsBadOpcode(err) {
		t.Error("Expected bad opcode error, got:", err)
	}
}

func TestThreadㆍBadAccess(t *testing.T) {
	th := Thread{reg: make([]float64, 256)}
	err := th.exec(OpAdd, []float64{-1, -2, -3}, 0xFFFFFFFF)
	if !IsBadAccess(err) {
		t.Error("Expected bad access error, got:", err)
	}

	defer expectPanic(t, CodeBadAccess)
	err = th.exec(OpAdd, []float64{0, -2, -3}, 0xFFFFFFFF)
	if err != nil {
		t.Error("Received error instead of panic:", err)
	}
}
