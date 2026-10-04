package matrixbuilder

import (
	"math"
	"strings"
	"testing"
)

func TestEthanolStoichAFR(t *testing.T) {
	cases := []struct {
		eth  float64
		want float64
	}{
		{-10, 14.70},
		{0, 14.70},
		{85, 9.76},
		{100, 9.76},
		{42.5, 12.23}, // halfway between 14.70 and 9.76
	}

	for _, c := range cases {
		got := EthanolStoichAFR(c.eth)
		if math.Abs(got-c.want) > 0.01 {
			t.Errorf("EthanolStoichAFR(%v) = %v, want %v", c.eth, got, c.want)
		}
	}
}

func TestGetDefaultTargetAFR(t *testing.T) {
	e0 := GetDefaultTargetAFR(0)
	if len(e0) != 18*16 {
		t.Fatalf("expected len 288, got %d", len(e0))
	}
	if e0[0] != 14.7 {
		t.Errorf("e0[0] = %v, want 14.7", e0[0])
	}
	if e0[17] != 12.0 {
		t.Errorf("e0[17] = %v, want 12.0", e0[17])
	}

	e85 := GetDefaultTargetAFR(85)
	expectedStoich := 9.8 // 14.7 * (9.76 / 14.7) = 9.76 -> round 1 decimal = 9.8
	if math.Abs(e85[0]-expectedStoich) > 0.1 {
		t.Errorf("e85[0] = %v, want ~%v", e85[0], expectedStoich)
	}
}

func TestGetClosedLoopMask(t *testing.T) {
	xData := []float64{100, 400, 700}
	yData := []float64{1000, 5000}
	limitRpm := []float64{1000, 5000}
	limitMax := []float64{650, 350}

	mask := GetClosedLoopMask(xData, yData, limitRpm, limitMax)
	// Row 0 (1000 RPM, max 650):
	// x=100 <= 650 -> true
	// x=400 <= 650 -> true
	// x=700 <= 650 -> false
	if !mask[0] || !mask[1] || mask[2] {
		t.Errorf("row 0 mask mismatch: %v", mask[0:3])
	}

	// Row 1 (5000 RPM, max 350):
	// x=100 <= 350 -> true
	// x=400 <= 350 -> false
	// x=700 <= 350 -> false
	if !mask[3] || mask[4] || mask[5] {
		t.Errorf("row 1 mask mismatch: %v", mask[3:6])
	}
}

func TestCalculateFuelAdjustment_LambdaInt_ClosedLoopOnly(t *testing.T) {
	// 2 cells: cell 0 is closed-loop, cell 1 is open-loop
	cfg := FuelAdjustmentConfig{
		ZSeries:     "Lambda.LambdaInt",
		LearnedZ:    []float64{4.0, 4.0}, // +4% trim
		Counts:      []int{10, 10},
		CurrentFuel: []float64{1.00, 1.00},
		TargetAFR:   []float64{14.7, 12.0},
		ClosedLoop:  []bool{true, false}, // cell 0 is closed, cell 1 is open
		EthanolPct:  0,
		Smoothing:   1.0, // 100% smoothing (direct application)
	}

	res := CalculateFuelAdjustment(cfg)
	if res.UpdatedCells != 1 {
		t.Errorf("expected 1 cell updated, got %d", res.UpdatedCells)
	}
	// Cell 0 (closed-loop) should be updated: 1.00 * (1 + 4/100) = 1.04
	if res.NewFuel[0] != 1.04 {
		t.Errorf("cell 0 new fuel = %v, want 1.04", res.NewFuel[0])
	}
	// Cell 1 (open-loop) should remain UNTOUCHED
	if res.NewFuel[1] != 1.00 {
		t.Errorf("cell 1 new fuel = %v, want 1.00 (untouched)", res.NewFuel[1])
	}
}

func TestCalculateFuelAdjustment_LambdaExternal_OpenLoopOnly(t *testing.T) {
	// 2 cells: cell 0 is closed-loop, cell 1 is open-loop
	// Target AFR: cell 1 target is 11.76 AFR -> Target Lambda = 11.76 / 14.7 = 0.80
	// Measured wideband lambda is 0.88 (+10% lean -> ratio 0.88 / 0.80 = 1.10)
	cfg := FuelAdjustmentConfig{
		ZSeries:     "Lambda.External",
		LearnedZ:    []float64{0.88, 0.88},
		Counts:      []int{15, 15},
		CurrentFuel: []float64{1.00, 1.00},
		TargetAFR:   []float64{14.7, 11.76},
		ClosedLoop:  []bool{true, false}, // cell 0 closed, cell 1 open
		EthanolPct:  0,
		Smoothing:   1.0,
	}

	res := CalculateFuelAdjustment(cfg)
	if res.UpdatedCells != 1 {
		t.Errorf("expected 1 cell updated, got %d", res.UpdatedCells)
	}
	// Cell 0 (closed-loop) must be SKIPPED because Lambda.External only updates open-loop cells
	if res.NewFuel[0] != 1.00 {
		t.Errorf("cell 0 (closed-loop) should be untouched with Lambda.External, got %v", res.NewFuel[0])
	}
	// Cell 1 (open-loop) should be updated: 1.00 * (0.88 / 0.80) = 1.10
	if res.NewFuel[1] != 1.10 {
		t.Errorf("cell 1 (open-loop) new fuel = %v, want 1.10", res.NewFuel[1])
	}
}

func TestCalculateFuelAdjustment_Smoothing(t *testing.T) {
	// With +10% correction ratio and 50% smoothing, expected ratio is 1.05
	cfg := FuelAdjustmentConfig{
		ZSeries:     "Lambda.LambdaInt",
		LearnedZ:    []float64{10.0},
		Counts:      []int{5},
		CurrentFuel: []float64{1.00},
		TargetAFR:   []float64{14.7},
		ClosedLoop:  []bool{true},
		EthanolPct:  0,
		Smoothing:   0.50, // 50%
	}

	res := CalculateFuelAdjustment(cfg)
	if res.NewFuel[0] != 1.05 {
		t.Errorf("smoothed fuel = %v, want 1.05", res.NewFuel[0])
	}
}

func TestCalculateFuelAdjustment_UnhitCellsSkipped(t *testing.T) {
	cfg := FuelAdjustmentConfig{
		ZSeries:     "Lambda.LambdaInt",
		LearnedZ:    []float64{10.0},
		Counts:      []int{0}, // 0 hits!
		CurrentFuel: []float64{1.00},
		TargetAFR:   []float64{14.7},
		ClosedLoop:  []bool{true},
		EthanolPct:  0,
		Smoothing:   1.0,
	}

	res := CalculateFuelAdjustment(cfg)
	if res.UpdatedCells != 0 {
		t.Errorf("unhit cell was updated, wanted 0 updates")
	}
	if res.NewFuel[0] != 1.00 {
		t.Errorf("unhit cell changed to %v, want 1.00", res.NewFuel[0])
	}
}

func TestSerializeT7FuelString(t *testing.T) {
	zData := []float64{1.0, 1.25, 0.85, 1.10}
	str := SerializeT7FuelString(zData, 2, 2)
	// Cell (0,0) must start with 20:0:
	if !strings.HasPrefix(str, "20:0:1:~") {
		t.Errorf("expected string to start with '20:0:1:~', got %q", str)
	}
	// Cell (1,0)
	if !strings.Contains(str, "1:0:1.25:~") {
		t.Errorf("missing cell (1,0), got %q", str)
	}
	// Cell (0,1)
	if !strings.Contains(str, "0:1:0.85:~") {
		t.Errorf("missing cell (0,1), got %q", str)
	}
}
