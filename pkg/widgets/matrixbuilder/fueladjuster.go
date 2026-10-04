package matrixbuilder

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/roffe/txlogger/pkg/colors"
	"github.com/roffe/txlogger/pkg/widgets/mapviewer"
)

// Default axes for Trionic 7 BFuelCal.Map and BFuelCal.StartMap (18 columns x 16 rows)
var (
	DefaultAirXSP = []float64{75, 125, 150, 180, 240, 300, 360, 420, 480, 540, 600, 660, 720, 780, 840, 900, 1050, 1300}
	DefaultRpmYSP = []float64{700, 880, 1260, 1640, 2020, 2400, 2780, 3160, 3540, 3920, 4300, 4680, 5060, 5440, 5820, 6200}

	// Reference T7 petrol closed-loop airmass boundary (LambdaCal.MaxLoadNormTab)
	DefaultClosedLoopRpm     = []float64{700, 880, 1260, 1640, 2020, 2400, 2780, 3160, 3540, 3920, 4300, 4680, 5060, 5440, 5820, 6200}
	DefaultClosedLoopMaxLoad = []float64{650, 650, 650, 650, 650, 650, 660, 660, 660, 650, 570, 510, 450, 330, 330, 330}

	// Reference T7 E85 closed-loop airmass boundary (LambdaCal.MaxLoadE85Tab)
	// Higher airmass is permitted in closed-loop on E85 due to lower knock sensitivity and charge cooling.
	DefaultClosedLoopMaxLoadE85 = []float64{750, 750, 750, 750, 750, 750, 780, 780, 780, 750, 680, 620, 550, 450, 450, 450}
)

// Hardcoded default base Target AFR table (18 cols x 16 rows in AFR units, petrol base 14.70).
// Row 0 is 700 RPM, Row 15 is 6200 RPM, aligned with DefaultRpmYSP and T7Suite.
var HardcodedBaseTargetAFR = []float64{
	// Row 0 (700 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 13.40, 13.20, 13.00, 12.90, 12.70, 12.50, 12.30,
	// Row 1 (880 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 13.40, 13.20, 13.00, 12.80, 12.60, 12.40, 12.20,
	// Row 2 (1260 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 13.30, 13.10, 12.90, 12.70, 12.50, 12.30, 12.10,
	// Row 3 (1640 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 13.20, 13.00, 12.80, 12.60, 12.40, 12.20, 12.00,
	// Row 4 (2020 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 13.10, 12.90, 12.70, 12.50, 12.30, 12.10, 11.90,
	// Row 5 (2400 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 13.00, 12.80, 12.60, 12.40, 12.20, 12.00, 11.80,
	// Row 6 (2780 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 12.90, 12.70, 12.50, 12.30, 12.10, 11.90, 11.70,
	// Row 7 (3160 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 12.80, 12.60, 12.40, 12.20, 12.00, 11.90, 11.70,
	// Row 8 (3540 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 12.70, 12.50, 12.30, 12.10, 11.90, 11.70, 11.60,
	// Row 9 (3920 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 12.60, 12.40, 12.20, 12.00, 11.90, 11.70, 11.50,
	// Row 10 (4300 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 12.70, 12.50, 12.30, 12.10, 11.90, 11.70, 11.50,
	// Row 11 (4680 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 12.80, 12.60, 12.40, 12.20, 12.00, 11.80, 11.60,
	// Row 12 (5060 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 12.90, 12.70, 12.50, 12.30, 12.10, 11.90, 11.70,
	// Row 13 (5440 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 13.00, 12.80, 12.60, 12.40, 12.20, 12.00, 11.80,
	// Row 14 (5820 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 13.10, 12.90, 12.70, 12.50, 12.30, 12.10, 11.90,
	// Row 15 (6200 rpm)
	14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 14.70, 13.20, 13.00, 12.80, 12.60, 12.40, 12.20, 12.00,
}

// EthanolStoichAFR returns the stoichiometric air-fuel ratio for a given ethanol percentage (0..85).
// Pure petrol (E0) = 14.70, E85 = 9.76.
func EthanolStoichAFR(ethanolPct float64) float64 {
	if ethanolPct < 0 {
		ethanolPct = 0
	}
	if ethanolPct > 85 {
		ethanolPct = 85
	}
	return 14.70 - (ethanolPct/85.0)*(14.70-9.76)
}

// GetDefaultTargetAFR returns a copy of the hardcoded Target AFR map scaled for the given ethanol percentage (2 decimal points).
func GetDefaultTargetAFR(ethanolPct float64) []float64 {
	stoich := EthanolStoichAFR(ethanolPct)
	ratio := stoich / 14.70
	out := make([]float64, len(HardcodedBaseTargetAFR))
	for i, v := range HardcodedBaseTargetAFR {
		out[i] = math.Round(v*ratio*100) / 100
	}
	return out
}

// GetDefaultTargetLambda returns the hardcoded baseline in Lambda units (AFR / 14.70) (3 decimal points).
func GetDefaultTargetLambda() []float64 {
	out := make([]float64, len(HardcodedBaseTargetAFR))
	for i, v := range HardcodedBaseTargetAFR {
		out[i] = math.Round((v/14.70)*1000) / 1000
	}
	return out
}

// Lookup1D performs clamped linear interpolation of ys over xs at point x.
func Lookup1D(xs, ys []float64, x float64) float64 {
	n := len(xs)
	if n == 0 || len(ys) != n {
		return 0
	}
	if x <= xs[0] {
		return ys[0]
	}
	if x >= xs[n-1] {
		return ys[n-1]
	}
	for i := 0; i < n-1; i++ {
		if x >= xs[i] && x <= xs[i+1] {
			span := xs[i+1] - xs[i]
			if span == 0 {
				return ys[i]
			}
			frac := (x - xs[i]) / span
			return ys[i] + frac*(ys[i+1]-ys[i])
		}
	}
	return ys[n-1]
}

// GetClosedLoopMask returns a bool slice (row-major: r*cols + c) flagging whether each cell
// is in closed loop (airmass <= maxLoad at that row's RPM).
func GetClosedLoopMask(xData, yData, limitRpm, limitMaxLoad []float64) []bool {
	mask := make([]bool, len(xData)*len(yData))
	for r, rpm := range yData {
		maxAir := Lookup1D(limitRpm, limitMaxLoad, rpm)
		for c, air := range xData {
			mask[r*len(xData)+c] = air <= maxAir
		}
	}
	return mask
}

// FuelAdjustmentConfig holds all inputs needed to calculate the updated fuel map.
type FuelAdjustmentConfig struct {
	ZSeries      string    // "Lambda.LambdaInt" or "Lambda.External"
	TargetSymbol string    // "BFuelCal.Map" or "BFuelCal.StartMap" (optional, defaults by EthanolPct)
	LearnedZ     []float64 // MatrixBuilder learned values (same length as CurrentFuel)
	Counts       []int     // MatrixBuilder sample counts per cell
	CurrentFuel  []float64 // Base fuel map values (BFuelCal.Map or BFuelCal.StartMap)
	OtherFuel    []float64 // Complementary fuel map (StartMap when tuning Map, or Map when tuning StartMap)
	EthanolPct   float64   // Flex fuel ethanol percentage (0..85)
	TargetLambda []float64 // Target lambda values (per cell)
	ClosedLoop   []bool    // true = closed loop, false = open loop
	Smoothing    float64   // 0.0 .. 1.0 (e.g. 0.5 for 50%)
}

// FuelAdjustmentResult contains the computed new fuel map and adjustment statistics.
type FuelAdjustmentResult struct {
	NewFuel      []float64
	Difference   []float64
	UpdatedCells int
	SkippedCells int
	MinDelta     float64
	MaxDelta     float64
	AvgDelta     float64
}

// IsClosedLoopSeries returns true if series is a closed-loop trim series (Lambda.LambdaInt).
func IsClosedLoopSeries(s string) bool {
	lower := strings.ToLower(s)
	return strings.Contains(lower, "lambdaint")
}

// IsOpenLoopSeries returns true if series is an external wideband sensor (Lambda.External).
func IsOpenLoopSeries(s string) bool {
	lower := strings.ToLower(s)
	return strings.Contains(lower, "external") || strings.Contains(lower, "wbl")
}

// CalculateFuelAdjustment calculates the new fuel map:
// - If Lambda.LambdaInt is evaluated: ONLY closed loop cells are updated.
// - If Lambda.External is evaluated: ONLY open loop cells are updated.
// - Takes into account the Biopower flexfuel blend ((1-w)*Map + w*StartMap) to de-blend corrections.
func CalculateFuelAdjustment(cfg FuelAdjustmentConfig) FuelAdjustmentResult {
	n := len(cfg.CurrentFuel)
	res := FuelAdjustmentResult{
		NewFuel:    make([]float64, n),
		Difference: make([]float64, n),
	}
	copy(res.NewFuel, cfg.CurrentFuel)

	isLambdaInt := IsClosedLoopSeries(cfg.ZSeries)
	isExternal := IsOpenLoopSeries(cfg.ZSeries)

	var totalDelta float64
	first := true

	for i := 0; i < n; i++ {
		// Only cells with at least 1 hit in the log can be updated
		if i >= len(cfg.Counts) || cfg.Counts[i] == 0 {
			res.SkippedCells++
			continue
		}

		isClosed := false
		if i < len(cfg.ClosedLoop) {
			isClosed = cfg.ClosedLoop[i]
		}

		// Core rule:
		// Lambda.LambdaInt -> update ONLY closed loop cells
		// Lambda.External -> update ONLY open loop cells
		if isLambdaInt && !isClosed {
			res.SkippedCells++
			continue
		}
		if isExternal && isClosed {
			res.SkippedCells++
			continue
		}

		curFuel := cfg.CurrentFuel[i]
		learned := cfg.LearnedZ[i]
		var correctionRatio float64

		if isLambdaInt {
			// LambdaInt is short-term trim in % (e.g. +3.5 means +3.5%)
			trim := learned
			correctionRatio = 1.0 + (trim / 100.0)
		} else if isExternal {
			// External is actual measured lambda (e.g. 0.85)
			if learned <= 0 {
				res.SkippedCells++
				continue
			}
			tgtLam := 1.00
			if i < len(cfg.TargetLambda) && cfg.TargetLambda[i] > 0 {
				tgtLam = cfg.TargetLambda[i]
			}
			correctionRatio = learned / tgtLam
		} else {
			res.SkippedCells++
			continue
		}

		// Apply smoothing
		smoothing := cfg.Smoothing
		if smoothing < 0 {
			smoothing = 0
		} else if smoothing > 1 {
			smoothing = 1
		}
		smoothedRatio := 1.0 + (correctionRatio-1.0)*smoothing

		// Biopower Flexfuel blend weighting:
		// w = ethanolPct / 85.0
		// In Trionic 7, injected fuel is a blend: (1 - w)*BFuelCal.Map + w*BFuelCal.StartMap
		w := cfg.EthanolPct / 85.0
		if w < 0 {
			w = 0
		} else if w > 1 {
			w = 1
		}

		otherVal := curFuel
		if i < len(cfg.OtherFuel) && cfg.OtherFuel[i] > 0 {
			otherVal = cfg.OtherFuel[i]
		}

		isStartMap := (cfg.TargetSymbol == "BFuelCal.StartMap") || (cfg.TargetSymbol == "" && cfg.EthanolPct >= 55.0)

		var newFuel float64
		if isStartMap {
			// Tuning BFuelCal.StartMap (E85 table). Weight on StartMap is w.
			if w < 0.05 {
				newFuel = curFuel * smoothedRatio
			} else {
				currentBlended := (1.0-w)*otherVal + w*curFuel
				targetBlended := currentBlended * smoothedRatio
				newFuel = (targetBlended - (1.0-w)*otherVal) / w
			}
		} else {
			// Tuning BFuelCal.Map (petrol table). Weight on Map is (1 - w).
			if (1.0 - w) < 0.05 {
				newFuel = curFuel * smoothedRatio
			} else {
				currentBlended := (1.0-w)*curFuel + w*otherVal
				targetBlended := currentBlended * smoothedRatio
				newFuel = (targetBlended - w*otherVal) / (1.0 - w)
			}
		}

		if newFuel < 0.10 {
			newFuel = 0.10
		}
		newFuel = math.Round(newFuel*100) / 100

		delta := newFuel - curFuel
		res.NewFuel[i] = newFuel
		res.Difference[i] = math.Round(delta*100) / 100
		res.UpdatedCells++
		totalDelta += math.Abs(delta)

		if first {
			res.MinDelta = delta
			res.MaxDelta = delta
			first = false
		} else {
			if delta < res.MinDelta {
				res.MinDelta = delta
			}
			if delta > res.MaxDelta {
				res.MaxDelta = delta
			}
		}
	}

	for i := range res.Difference {
		res.Difference[i] = math.Round((res.NewFuel[i]-cfg.CurrentFuel[i])*100) / 100
	}

	if res.UpdatedCells > 0 {
		res.AvgDelta = totalDelta / float64(res.UpdatedCells)
	}

	return res
}

// SerializeT7FuelString serializes map data into the T7Suite clipboard string format:
// "20:0:val:~1:0:val:~...x:y:val:~"
// In T7Suite, Y-axis coordinates start at 0 for the highest RPM (6200) down to rows-1 for lowest RPM (700).
// In txlogger, row 0 is 700 RPM and row rows-1 is 6200 RPM.
// Therefore, T7Suite y = (rows - 1) - r.
func SerializeT7FuelString(zData []float64, cols, rows int) string {
	tokens := make([]string, 0, len(zData))
	for y := 0; y < rows; y++ {
		r := (rows - 1) - y // map T7Suite y (top=0) to txlogger row r (bottom=0)
		for c := 0; c < cols; c++ {
			idx := r*cols + c
			if idx >= len(zData) {
				continue
			}
			val := zData[idx]
			outX := c
			if y == 0 && c == 0 {
				outX = 20
			}
			var valStr string
			if val == math.Trunc(val) {
				valStr = strconv.Itoa(int(val))
			} else {
				valStr = strconv.FormatFloat(val, 'f', 2, 64)
			}
			tokens = append(tokens, fmt.Sprintf("%d:%d:%s", outX, y, valStr))
		}
	}
	return strings.Join(tokens, ":~") + ":~"
}

// FuelUpdater is the interface provided by MainWindow to read/update BFuelCal maps and closed loop regions.
type FuelUpdater interface {
	GetFuelMap(symbolName string) (xData, yData, zData []float64, err error)
	GetClosedLoopRegion(regionSymbol string, xData, yData []float64) []bool
	UpdateFuelMap(symbolName string, zData []float64) error
	Log(msg string)
}

// FuelAdjusterWidget is the interactive UI component to tune BFuelCal maps using MatrixBuilder results.
type FuelAdjusterWidget struct {
	widget.BaseWidget
	mb      *MatrixBuilder
	updater FuelUpdater

	ethanolPct    float64
	smoothing     float64
	isLambda      bool   // false = Target AFR mode, true = Target Lambda mode
	targetMapMode string // "Auto", "BFuelCal.Map", "BFuelCal.StartMap"

	targetAFR    []float64 // base target AFR values (petrol 14.70 base)
	targetLambda []float64 // normalized target lambda values
	displayedZ   []float64 // values currently shown in targetViewer (AFR or Lambda)
	newFuel      []float64
	difference   []float64
	result       FuelAdjustmentResult

	// UI controls
	modeSelect      *widget.Select
	targetMapSelect *widget.Select
	targetMapLbl    *widget.Label
	ethSlider       *widget.Slider
	ethLabel        *widget.Label
	smoothSlider    *widget.Slider
	smoothLabel     *widget.Label
	ruleBadge       *widget.Label
	warningBadge    *widget.Label
	statsLabel      *widget.Label
	statusLabel     *widget.Label
	applyBtn        *widget.Button

	targetViewer  *mapviewer.MapViewer
	previewViewer *mapviewer.MapViewer
	diffViewer    *mapviewer.MapViewer
	tabContainer  *container.AppTabs

	content fyne.CanvasObject
}

// NewFuelAdjusterWidget creates a new fuel adjuster interface.
func NewFuelAdjusterWidget(mb *MatrixBuilder, updater FuelUpdater) *FuelAdjusterWidget {
	fa := &FuelAdjusterWidget{
		mb:            mb,
		updater:       updater,
		ethanolPct:    0,
		smoothing:     0.50,
		isLambda:      false,
		targetMapMode: "Auto",
		targetAFR:     GetDefaultTargetAFR(0),
		targetLambda:  GetDefaultTargetLambda(),
	}
	fa.ExtendBaseWidget(fa)
	fa.buildUI()
	fa.recalculate()
	return fa
}

func (fa *FuelAdjusterWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(fa.content)
}

func (fa *FuelAdjusterWidget) targetMapNames() (fuelSymbol, regionSymbol string) {
	if fa.targetMapMode == "BFuelCal.StartMap" {
		return "BFuelCal.StartMap", "LambdaCal.MaxLoadE85Tab"
	}
	if fa.targetMapMode == "BFuelCal.Map" {
		return "BFuelCal.Map", "LambdaCal.MaxLoadNormTab"
	}
	// "Auto" mode: switch to BFuelCal.StartMap at 55% or higher
	if fa.ethanolPct >= 55.0 {
		return "BFuelCal.StartMap", "LambdaCal.MaxLoadE85Tab"
	}
	return "BFuelCal.Map", "LambdaCal.MaxLoadNormTab"
}

func (fa *FuelAdjusterWidget) buildUI() {
	// Mode selector: Target AFR vs. Target Lambda
	fa.modeSelect = widget.NewSelect([]string{"Target AFR", "Target Lambda"}, func(s string) {
		fa.isLambda = (s == "Target Lambda")
		fa.updateTargetViewerDisplay()
		fa.recalculate()
	})
	fa.modeSelect.Selected = "Target AFR"

	// Target Map dropdown: Auto vs explicit BFuelCal.Map vs BFuelCal.StartMap
	fa.targetMapSelect = widget.NewSelect([]string{"Auto", "BFuelCal.Map (Petrol)", "BFuelCal.StartMap (E85)"}, func(s string) {
		switch s {
		case "BFuelCal.Map (Petrol)":
			fa.targetMapMode = "BFuelCal.Map"
		case "BFuelCal.StartMap (E85)":
			fa.targetMapMode = "BFuelCal.StartMap"
		default:
			fa.targetMapMode = "Auto"
		}
		fuelSymbol, _ := fa.targetMapNames()
		fa.targetMapLbl.SetText("Target: " + fuelSymbol)
		if fa.applyBtn != nil {
			fa.applyBtn.SetText("Apply to " + fuelSymbol)
		}
		fa.recalculate()
	})
	fa.targetMapSelect.Selected = "Auto"

	// Target Map label (e.g. BFuelCal.Map vs BFuelCal.StartMap)
	fa.targetMapLbl = widget.NewLabelWithStyle("Target: BFuelCal.Map", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Ethanol slider (0..85%) with 0.1% precision
	fa.ethLabel = widget.NewLabel("0.0% (Stoich: 14.70)")
	fa.ethSlider = widget.NewSlider(0, 85)
	fa.ethSlider.Step = 0.1
	fa.ethSlider.SetValue(fa.ethanolPct)
	fa.ethSlider.OnChanged = func(val float64) {
		fa.ethanolPct = math.Round(val*10) / 10
		stoich := EthanolStoichAFR(fa.ethanolPct)
		fa.ethLabel.SetText(fmt.Sprintf("%.1f%% (Stoich: %.2f)", fa.ethanolPct, stoich))

		fuelSymbol, _ := fa.targetMapNames()
		fa.targetMapLbl.SetText("Target: " + fuelSymbol)
		if fa.applyBtn != nil {
			fa.applyBtn.SetText("Apply to " + fuelSymbol)
		}

		fa.updateTargetViewerDisplay()
		fa.recalculate()
	}

	// Smoothing slider (0..100%)
	fa.smoothLabel = widget.NewLabel("50%")
	fa.smoothSlider = widget.NewSlider(0, 100)
	fa.smoothSlider.Step = 5
	fa.smoothSlider.SetValue(fa.smoothing * 100)
	fa.smoothSlider.OnChanged = func(val float64) {
		fa.smoothing = val / 100.0
		fa.smoothLabel.SetText(fmt.Sprintf("%d%%", int(val)))
		fa.recalculate()
	}

	// Rule explanation badge
	fa.ruleBadge = widget.NewLabel("")
	fa.ruleBadge.TextStyle = fyne.TextStyle{Bold: true}

	// Mid-blend warning badge
	fa.warningBadge = widget.NewLabel("")
	fa.warningBadge.TextStyle = fyne.TextStyle{Bold: true}

	fa.statsLabel = widget.NewLabel("")
	fa.statusLabel = widget.NewLabel("")

	resetAFRBtn := widget.NewButtonWithIcon("Reset Target", theme.ViewRefreshIcon(), func() {
		fa.targetAFR = GetDefaultTargetAFR(0)
		fa.targetLambda = GetDefaultTargetLambda()
		fa.updateTargetViewerDisplay()
		fa.recalculate()
		fa.statusLabel.SetText("Reset Target table to defaults")
	})

	topControls := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Series:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(fa.mb.zSeries),
			widget.NewSeparator(),
			fa.targetMapLbl,
			layout.NewSpacer(),
			fa.ruleBadge,
		),
		container.NewGridWithColumns(4,
			container.NewBorder(nil, nil, widget.NewLabel("Mode:"), nil, fa.modeSelect),
			container.NewBorder(nil, nil, widget.NewLabel("Target:"), nil, fa.targetMapSelect),
			container.NewBorder(nil, nil, widget.NewLabel("Ethanol:"), fa.ethLabel, fa.ethSlider),
			container.NewBorder(nil, nil, widget.NewLabel("Smoothing:"), fa.smoothLabel, fa.smoothSlider),
		),
		container.NewHBox(
			resetAFRBtn,
			layout.NewSpacer(),
			fa.warningBadge,
			layout.NewSpacer(),
			fa.statsLabel,
		),
		widget.NewSeparator(),
	)

	initialMask := fa.getClosedLoopMask("LambdaCal.MaxLoadNormTab", DefaultAirXSP, DefaultRpmYSP)
	fa.displayedZ = make([]float64, len(HardcodedBaseTargetAFR))
	copy(fa.displayedZ, HardcodedBaseTargetAFR)

	// Tab 1: Target Map (Editable)
	tCfg := &mapviewer.Config{
		Name:           "Target Map",
		XData:          DefaultAirXSP,
		YData:          DefaultRpmYSP,
		ZData:          fa.displayedZ,
		RegionBorder:   initialMask,
		XPrecision:     0,
		YPrecision:     0,
		ZPrecision:     2,
		XLabel:         "Airmass (mg/c)",
		YLabel:         "RPM",
		ZLabel:         "Target AFR",
		Editable:       true,
		ColorblindMode: colors.ModeNormal,
		OnUpdateCell: func(idx int, data []float64) {
			if idx >= 0 && idx < len(fa.targetAFR) {
				stoich := EthanolStoichAFR(fa.ethanolPct)
				val := data[idx]
				if fa.isLambda {
					fa.targetLambda[idx] = val
					fa.targetAFR[idx] = val * 14.70
				} else {
					fa.targetAFR[idx] = val / (stoich / 14.70)
					fa.targetLambda[idx] = val / stoich
				}
				fa.recalculate()
			}
		},
	}
	tMv, _ := mapviewer.New(tCfg)
	fa.targetViewer = tMv

	// Tab 2: Adjusted Fuel Map Preview
	pCfg := &mapviewer.Config{
		Name:           "Adjusted Fuel Map Preview",
		XData:          DefaultAirXSP,
		YData:          DefaultRpmYSP,
		ZData:          make([]float64, len(HardcodedBaseTargetAFR)),
		RegionBorder:   initialMask,
		XPrecision:     0,
		YPrecision:     0,
		ZPrecision:     2,
		XLabel:         "Airmass (mg/c)",
		YLabel:         "RPM",
		ZLabel:         "Fuel Factor",
		Editable:       false,
		ColorblindMode: colors.ModeNormal,
	}
	pMv, _ := mapviewer.New(pCfg)
	fa.previewViewer = pMv

	// Tab 3: Difference Map (Δ New - Current)
	dCfg := &mapviewer.Config{
		Name:           "Difference (Δ New - Current)",
		XData:          DefaultAirXSP,
		YData:          DefaultRpmYSP,
		ZData:          make([]float64, len(HardcodedBaseTargetAFR)),
		RegionBorder:   initialMask,
		XPrecision:     0,
		YPrecision:     0,
		ZPrecision:     2,
		XLabel:         "Airmass (mg/c)",
		YLabel:         "RPM",
		ZLabel:         "Δ Fuel Factor",
		Editable:       false,
		ColorblindMode: colors.ModeNormal,
	}
	dMv, _ := mapviewer.New(dCfg)
	fa.diffViewer = dMv

	fuelSymbol, _ := fa.targetMapNames()
	fa.tabContainer = container.NewAppTabs(
		container.NewTabItemWithIcon("Target Map (Editable)", theme.DocumentIcon(), fa.targetViewer),
		container.NewTabItemWithIcon("Preview "+fuelSymbol, theme.GridIcon(), fa.previewViewer),
		container.NewTabItemWithIcon("Δ Difference (New - Current)", theme.ViewRefreshIcon(), fa.diffViewer),
	)

	// Bottom action buttons
	fa.applyBtn = widget.NewButtonWithIcon("Apply to "+fuelSymbol, theme.DocumentSaveIcon(), func() {
		fa.applyToFuelMap()
	})
	fa.applyBtn.Importance = widget.HighImportance

	copyT7Btn := widget.NewButtonWithIcon("Copy T7Suite String", theme.ContentCopyIcon(), func() {
		fa.copyT7SuiteString()
	})

	bottomBar := container.NewBorder(nil, nil,
		fa.statusLabel,
		container.NewHBox(copyT7Btn, fa.applyBtn),
	)

	fa.content = container.NewBorder(topControls, bottomBar, nil, nil, fa.tabContainer)
	fa.updateTargetViewerDisplay()
}

func (fa *FuelAdjusterWidget) updateTargetViewerDisplay() {
	if fa.targetViewer == nil {
		return
	}
	stoich := EthanolStoichAFR(fa.ethanolPct)
	ratio := stoich / 14.70
	displayData := make([]float64, len(fa.targetAFR))
	if fa.isLambda {
		fa.targetViewer.SetZPrecision(3)
		for i, afr := range fa.targetAFR {
			lam := afr / 14.70
			displayData[i] = math.Round(lam*1000) / 1000
		}
	} else {
		fa.targetViewer.SetZPrecision(2)
		for i, afr := range fa.targetAFR {
			displayData[i] = math.Round(afr*ratio*100) / 100
		}
	}
	_ = fa.targetViewer.SetZData(displayData)
	fa.targetViewer.Refresh()

	if len(fa.tabContainer.Items) > 0 {
		label := "Target Map (AFR) (Editable)"
		if fa.isLambda {
			label = "Target Map (Lambda) (Editable)"
		}
		fa.tabContainer.Items[0].Text = label
		fa.tabContainer.Refresh()
	}
}

func (fa *FuelAdjusterWidget) getClosedLoopMask(regionSymbol string, xData, yData []float64) []bool {
	if fa.updater != nil {
		if mask := fa.updater.GetClosedLoopRegion(regionSymbol, xData, yData); len(mask) == len(xData)*len(yData) {
			return mask
		}
	}
	if regionSymbol == "LambdaCal.MaxLoadE85Tab" {
		return GetClosedLoopMask(xData, yData, DefaultClosedLoopRpm, DefaultClosedLoopMaxLoadE85)
	}
	return GetClosedLoopMask(xData, yData, DefaultClosedLoopRpm, DefaultClosedLoopMaxLoad)
}

func (fa *FuelAdjusterWidget) getBaseFuelMap(fuelSymbol string) (xData, yData, zData []float64) {
	if fa.updater != nil {
		if x, y, z, err := fa.updater.GetFuelMap(fuelSymbol); err == nil && len(z) == len(x)*len(y) {
			return x, y, z
		}
	}
	// Fallback reference map (1.00 base)
	size := len(DefaultAirXSP) * len(DefaultRpmYSP)
	fallbackZ := make([]float64, size)
	for i := range fallbackZ {
		fallbackZ[i] = 1.00
	}
	return DefaultAirXSP, DefaultRpmYSP, fallbackZ
}

func (fa *FuelAdjusterWidget) recalculate() {
	zSeries := fa.mb.zSeries
	if zSeries == "" {
		zSeries = "Lambda.External"
	}

	fuelSymbol, regionSymbol := fa.targetMapNames()
	otherSymbol := "BFuelCal.StartMap"
	if fuelSymbol == "BFuelCal.StartMap" {
		otherSymbol = "BFuelCal.Map"
	}

	// Update warning badge for mid-blends
	if fa.ethanolPct >= 20.0 && fa.ethanolPct <= 65.0 {
		fa.warningBadge.SetText("⚠ Warning: Mid-blend (20-65% E) - high cross-table sensitivity! Tune near E0 or E85 recommended.")
	} else {
		fa.warningBadge.SetText("")
	}

	// Update rule badge
	if IsClosedLoopSeries(zSeries) {
		fa.ruleBadge.SetText("Mode: CLOSED-LOOP ONLY (updating cells where Airmass <= MaxLoad)")
	} else if IsOpenLoopSeries(zSeries) {
		fa.ruleBadge.SetText("Mode: OPEN-LOOP ONLY (updating cells where Airmass > MaxLoad)")
	} else {
		fa.ruleBadge.SetText("Mode: Neutral (" + zSeries + ")")
	}

	xData, yData, baseZ := fa.getBaseFuelMap(fuelSymbol)
	_, _, otherZ := fa.getBaseFuelMap(otherSymbol)
	mask := fa.getClosedLoopMask(regionSymbol, xData, yData)

	learnedZ := fa.mb.zData
	counts := fa.mb.counts

	// If matrix dimensions differ, create zero counts
	if len(learnedZ) != len(baseZ) {
		learnedZ = make([]float64, len(baseZ))
		counts = make([]int, len(baseZ))
	}

	// Keep targetLambda synchronized with targetAFR
	if fa.targetLambda == nil || len(fa.targetLambda) != len(fa.targetAFR) {
		fa.targetLambda = make([]float64, len(fa.targetAFR))
	}
	for i, afr := range fa.targetAFR {
		fa.targetLambda[i] = afr / 14.70
	}

	cfg := FuelAdjustmentConfig{
		ZSeries:      zSeries,
		TargetSymbol: fuelSymbol,
		LearnedZ:     learnedZ,
		Counts:       counts,
		CurrentFuel:  baseZ,
		OtherFuel:    otherZ,
		EthanolPct:   fa.ethanolPct,
		TargetLambda: fa.targetLambda,
		ClosedLoop:   mask,
		Smoothing:    fa.smoothing,
	}

	res := CalculateFuelAdjustment(cfg)
	fa.result = res
	fa.newFuel = res.NewFuel
	fa.difference = res.Difference

	// Update all 3 viewers with the grey border and new data
	if fa.targetViewer != nil {
		fa.targetViewer.SetRegionBorder(mask)
	}
	if fa.previewViewer != nil {
		fa.previewViewer.SetRegionBorder(mask)
		_ = fa.previewViewer.SetZData(res.NewFuel)
		fa.previewViewer.Refresh()
	}
	if fa.diffViewer != nil {
		fa.diffViewer.SetRegionBorder(mask)
		_ = fa.diffViewer.SetZData(res.Difference)
		fa.diffViewer.Refresh()
	}

	if len(fa.tabContainer.Items) > 1 {
		fa.tabContainer.Items[1].Text = "Preview " + fuelSymbol
		fa.tabContainer.Refresh()
	}

	fa.statsLabel.SetText(fmt.Sprintf("%s | Updated: %d | Skipped: %d | Avg Δ: %+.2f",
		fuelSymbol, res.UpdatedCells, res.SkippedCells, res.AvgDelta))
}

func (fa *FuelAdjusterWidget) applyToFuelMap() {
	if fa.updater == nil {
		fa.statusLabel.SetText("No active ECU binary connection to update fuel map")
		return
	}
	fuelSymbol, _ := fa.targetMapNames()
	if err := fa.updater.UpdateFuelMap(fuelSymbol, fa.newFuel); err != nil {
		fa.statusLabel.SetText("Error updating " + fuelSymbol + ": " + err.Error())
		return
	}
	fa.statusLabel.SetText(fmt.Sprintf("Successfully updated %d cells in %s!", fa.result.UpdatedCells, fuelSymbol))
}

func (fa *FuelAdjusterWidget) copyT7SuiteString() {
	fuelSymbol, _ := fa.targetMapNames()
	xData, yData, _ := fa.getBaseFuelMap(fuelSymbol)
	str := SerializeT7FuelString(fa.newFuel, len(xData), len(yData))
	if app := fyne.CurrentApp(); app != nil {
		if drv := app.Driver(); drv != nil {
			windows := drv.AllWindows()
			if len(windows) > 0 {
				windows[0].Clipboard().SetContent(str)
				fa.statusLabel.SetText(fmt.Sprintf("Copied %s to clipboard (T7Suite format)", fuelSymbol))
				return
			}
		}
	}
	fa.statusLabel.SetText("Failed to access system clipboard")
}
