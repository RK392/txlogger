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
	if e0[0] != 14.70 {
		t.Errorf("e0[0] = %v, want 14.70", e0[0])
	}
	// Row 0 (700 rpm) at 1300 load (col 17) is 12.30
	if e0[17] != 12.30 {
		t.Errorf("e0[17] = %v, want 12.30", e0[17])
	}
	// Row 15 (6200 rpm) at 1300 load (col 17) is 12.00
	if e0[15*18+17] != 12.00 {
		t.Errorf("e0[15*18+17] = %v, want 12.00", e0[15*18+17])
	}

	e85 := GetDefaultTargetAFR(85)
	expectedStoich := 9.76 // 14.70 * (9.76 / 14.70) = 9.76
	if math.Abs(e85[0]-expectedStoich) > 0.05 {
		t.Errorf("e85[0] = %v, want ~%v", e85[0], expectedStoich)
	}

	lambdas := GetDefaultTargetLambda()
	if len(lambdas) != 18*16 {
		t.Fatalf("expected len 288, got %d", len(lambdas))
	}
	if lambdas[0] != 1.00 {
		t.Errorf("lambdas[0] = %v, want 1.00", lambdas[0])
	}
	// Row 15 (6200 rpm) at 1300 load is 0.816 (12.00 / 14.70)
	if lambdas[15*18+17] != 0.816 {
		t.Errorf("lambdas[15*18+17] = %v, want 0.816", lambdas[15*18+17])
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

	// Test E85 higher boundary vs Petrol
	e85Mask := GetClosedLoopMask([]float64{700}, []float64{2000}, DefaultClosedLoopRpm, DefaultClosedLoopMaxLoadE85)
	normMask := GetClosedLoopMask([]float64{700}, []float64{2000}, DefaultClosedLoopRpm, DefaultClosedLoopMaxLoad)
	// At 2000 RPM: E85 max is 750 (700 <= 750 -> true), Petrol max is 650 (700 <= 650 -> false)
	if !e85Mask[0] {
		t.Errorf("expected 700 mg/c to be closed loop on E85 at 2000 RPM")
	}
	if normMask[0] {
		t.Errorf("expected 700 mg/c to be open loop on petrol at 2000 RPM")
	}
}

func TestCalculateFuelAdjustment_LambdaInt_ClosedLoopOnly(t *testing.T) {
	// 2 cells: cell 0 is closed-loop, cell 1 is open-loop
	cfg := FuelAdjustmentConfig{
		ZSeries:      "Lambda.LambdaInt",
		LearnedZ:     []float64{4.0, 4.0}, // +4% trim
		Counts:       []int{10, 10},
		CurrentFuel:  []float64{1.00, 1.00},
		TargetLambda: []float64{1.00, 0.85},
		ClosedLoop:   []bool{true, false}, // cell 0 is closed, cell 1 is open
		Smoothing:    1.0,                 // 100% smoothing (direct application)
	}

	res := CalculateFuelAdjustment(cfg)
	if res.UpdatedCells != 1 {
		t.Errorf("expected 1 cell updated, got %d", res.UpdatedCells)
	}
	// Cell 0 (closed-loop) should be updated: 1.00 * (1 + 4/100) = 1.04
	if res.NewFuel[0] != 1.04 {
		t.Errorf("cell 0 new fuel = %v, want 1.04", res.NewFuel[0])
	}
	// Difference map should show +0.04
	if res.Difference[0] != 0.04 {
		t.Errorf("cell 0 delta = %v, want 0.04", res.Difference[0])
	}
	// Cell 1 (open-loop) should remain UNTOUCHED
	if res.NewFuel[1] != 1.00 {
		t.Errorf("cell 1 new fuel = %v, want 1.00 (untouched)", res.NewFuel[1])
	}
	if res.Difference[1] != 0.00 {
		t.Errorf("cell 1 delta = %v, want 0.00", res.Difference[1])
	}
}

func TestCalculateFuelAdjustment_LambdaExternal_OpenLoopOnly(t *testing.T) {
	// 2 cells: cell 0 is closed-loop, cell 1 is open-loop
	// Target Lambda: cell 1 target is 0.80
	// Measured wideband lambda is 0.88 (+10% lean -> ratio 0.88 / 0.80 = 1.10)
	cfg := FuelAdjustmentConfig{
		ZSeries:      "Lambda.External",
		LearnedZ:     []float64{0.88, 0.88},
		Counts:       []int{15, 15},
		CurrentFuel:  []float64{1.00, 1.00},
		TargetLambda: []float64{1.00, 0.80},
		ClosedLoop:   []bool{true, false}, // cell 0 closed, cell 1 open
		Smoothing:    1.0,
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
	if res.Difference[1] != 0.10 {
		t.Errorf("cell 1 delta = %v, want 0.10", res.Difference[1])
	}
}

func TestCalculateFuelAdjustment_Smoothing(t *testing.T) {
	// With +10% correction ratio and 50% smoothing, expected ratio is 1.05
	cfg := FuelAdjustmentConfig{
		ZSeries:      "Lambda.LambdaInt",
		LearnedZ:     []float64{10.0},
		Counts:       []int{5},
		CurrentFuel:  []float64{1.00},
		TargetLambda: []float64{1.00},
		ClosedLoop:   []bool{true},
		Smoothing:    0.50, // 50%
	}

	res := CalculateFuelAdjustment(cfg)
	if res.NewFuel[0] != 1.05 {
		t.Errorf("smoothed fuel = %v, want 1.05", res.NewFuel[0])
	}
	if res.Difference[0] != 0.05 {
		t.Errorf("smoothed delta = %v, want 0.05", res.Difference[0])
	}
}

func TestCalculateFuelAdjustment_UnhitCellsSkipped(t *testing.T) {
	cfg := FuelAdjustmentConfig{
		ZSeries:      "Lambda.LambdaInt",
		LearnedZ:     []float64{10.0},
		Counts:       []int{0}, // 0 hits!
		CurrentFuel:  []float64{1.00},
		TargetLambda: []float64{1.00},
		ClosedLoop:   []bool{true},
		Smoothing:    1.0,
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
	// In T7Suite, y=0 is the top row (r=1 in txlogger), starting with '20:0:'
	if !strings.HasPrefix(str, "20:0:0.85:~") {
		t.Errorf("expected string to start with '20:0:0.85:~', got %q", str)
	}
	// Cell (1,0) - top right cell
	if !strings.Contains(str, "1:0:1.10:~") {
		t.Errorf("missing cell (1,0), got %q", str)
	}
	// Cell (0,1) - bottom left cell (r=0 in txlogger)
	if !strings.Contains(str, "0:1:1:~") {
		t.Errorf("missing cell (0,1), got %q", str)
	}
	// Cell (1,1) - bottom right cell
	if !strings.Contains(str, "1:1:1.25:~") {
		t.Errorf("missing cell (1,1), got %q", str)
	}
}

func TestTargetMapNames_EthanolSwitching(t *testing.T) {
	fa := &FuelAdjusterWidget{ethanolPct: 0, targetMapMode: "Auto"}
	fuelSym, regSym := fa.targetMapNames()
	if fuelSym != "BFuelCal.Map" || regSym != "LambdaCal.MaxLoadNormTab" {
		t.Errorf("at 0%% eth: got %s / %s, want BFuelCal.Map / LambdaCal.MaxLoadNormTab", fuelSym, regSym)
	}

	fa.ethanolPct = 50
	fuelSym, regSym = fa.targetMapNames()
	if fuelSym != "BFuelCal.Map" || regSym != "LambdaCal.MaxLoadNormTab" {
		t.Errorf("at 50%% eth: got %s / %s, want BFuelCal.Map / LambdaCal.MaxLoadNormTab", fuelSym, regSym)
	}

	fa.ethanolPct = 55
	fuelSym, regSym = fa.targetMapNames()
	if fuelSym != "BFuelCal.StartMap" || regSym != "LambdaCal.MaxLoadE85Tab" {
		t.Errorf("at 55%% eth: got %s / %s, want BFuelCal.StartMap / LambdaCal.MaxLoadE85Tab", fuelSym, regSym)
	}

	fa.ethanolPct = 85
	fuelSym, regSym = fa.targetMapNames()
	if fuelSym != "BFuelCal.StartMap" || regSym != "LambdaCal.MaxLoadE85Tab" {
		t.Errorf("at 85%% eth: got %s / %s, want BFuelCal.StartMap / LambdaCal.MaxLoadE85Tab", fuelSym, regSym)
	}

	// Explicit override to BFuelCal.Map even at 85%
	fa.targetMapMode = "BFuelCal.Map"
	fuelSym, regSym = fa.targetMapNames()
	if fuelSym != "BFuelCal.Map" || regSym != "LambdaCal.MaxLoadNormTab" {
		t.Errorf("explicit BFuelCal.Map at 85%% eth: got %s / %s, want BFuelCal.Map / LambdaCal.MaxLoadNormTab", fuelSym, regSym)
	}

	// Explicit override to BFuelCal.StartMap even at 0%
	fa.targetMapMode = "BFuelCal.StartMap"
	fa.ethanolPct = 0
	fuelSym, regSym = fa.targetMapNames()
	if fuelSym != "BFuelCal.StartMap" || regSym != "LambdaCal.MaxLoadE85Tab" {
		t.Errorf("explicit BFuelCal.StartMap at 0%% eth: got %s / %s, want BFuelCal.StartMap / LambdaCal.MaxLoadE85Tab", fuelSym, regSym)
	}
}

func TestCalculateFuelAdjustment_FlexFuelDeblend(t *testing.T) {
	// At E = 8.5% (w = 0.10):
	// BFuelCal.Map = 1.00, BFuelCal.StartMap = 1.20
	// Current blended fuel = 0.90 * 1.00 + 0.10 * 1.20 = 1.02
	// Target Lambda = 0.85, Learned Lambda = 0.8925 (+5% lean -> correctionRatio = 1.05)
	// Target blended fuel = 1.02 * 1.05 = 1.071
	// BFuelCal.Map_new = (1.071 - 0.10 * 1.20) / 0.90 = (1.071 - 0.120) / 0.90 = 0.951 / 0.90 = 1.05667 -> 1.06
	cfg := FuelAdjustmentConfig{
		ZSeries:      "Lambda.External",
		TargetSymbol: "BFuelCal.Map",
		LearnedZ:     []float64{0.8925},
		Counts:       []int{20},
		CurrentFuel:  []float64{1.00},
		OtherFuel:    []float64{1.20},
		EthanolPct:   8.5, // w = 0.10
		TargetLambda: []float64{0.85},
		ClosedLoop:   []bool{false}, // open-loop cell
		Smoothing:    1.0,
	}

	res := CalculateFuelAdjustment(cfg)
	if res.UpdatedCells != 1 {
		t.Fatalf("expected 1 cell updated, got %d", res.UpdatedCells)
	}
	if res.NewFuel[0] != 1.06 {
		t.Errorf("expected de-blended fuel = 1.06, got %v", res.NewFuel[0])
	}

	// Test tuning BFuelCal.StartMap at E = 76.5% (w = 0.90):
	// Current BFuelCal.Map (OtherFuel) = 1.00, BFuelCal.StartMap (CurrentFuel) = 1.30
	// Current blended fuel = (1 - 0.90)*1.00 + 0.90*1.30 = 0.10 + 1.17 = 1.27
	// Target Lambda = 0.85, Learned Lambda = 0.8925 (+5% lean -> correctionRatio = 1.05)
	// Target blended fuel = 1.27 * 1.05 = 1.3335
	// BFuelCal.StartMap_new = (1.3335 - 0.10 * 1.00) / 0.90 = 1.2335 / 0.90 = 1.3705 -> 1.37
	cfgStartMap := FuelAdjustmentConfig{
		ZSeries:      "Lambda.External",
		TargetSymbol: "BFuelCal.StartMap",
		LearnedZ:     []float64{0.8925},
		Counts:       []int{20},
		CurrentFuel:  []float64{1.30},
		OtherFuel:    []float64{1.00},
		EthanolPct:   76.5, // w = 0.90
		TargetLambda: []float64{0.85},
		ClosedLoop:   []bool{false},
		Smoothing:    1.0,
	}
	resStartMap := CalculateFuelAdjustment(cfgStartMap)
	if resStartMap.UpdatedCells != 1 {
		t.Fatalf("expected 1 cell updated, got %d", resStartMap.UpdatedCells)
	}
	if resStartMap.NewFuel[0] != 1.37 {
		t.Errorf("expected de-blended StartMap fuel = 1.37, got %v", resStartMap.NewFuel[0])
	}
}

