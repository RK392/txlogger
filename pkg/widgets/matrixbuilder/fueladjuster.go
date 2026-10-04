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

// Default axes for Trionic 7 BFuelCal.Map (18 columns x 16 rows)
var (
	DefaultAirXSP = []float64{75, 125, 150, 180, 240, 300, 360, 420, 480, 540, 600, 660, 720, 780, 840, 900, 1050, 1300}
	DefaultRpmYSP = []float64{700, 880, 1260, 1640, 2020, 2400, 2780, 3160, 3540, 3920, 4300, 4680, 5060, 5440, 5820, 6200}

	// Reference T7 closed-loop airmass boundary (LambdaCal.MaxLoadNormTab)
	DefaultClosedLoopRpm     = []float64{700, 880, 1260, 1640, 2020, 2400, 2780, 3160, 3540, 3920, 4300, 4680, 5060, 5440, 5820, 6200}
	DefaultClosedLoopMaxLoad = []float64{650, 650, 650, 650, 650, 650, 660, 660, 660, 650, 570, 510, 450, 330, 330, 330}
)

// Hardcoded default base Target AFR table (18 cols x 16 rows in AFR units, petrol base).
// Columns 0..10 are 14.7 (stoichiometric); columns 11..17 taper richer under boost.
var HardcodedBaseTargetAFR = []float64{
	// Row 0 (700 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 13.2, 13.0, 12.8, 12.6, 12.4, 12.2, 12.0,
	// Row 1 (880 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 13.1, 12.9, 12.7, 12.5, 12.3, 12.1, 11.9,
	// Row 2 (1260 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 13.0, 12.8, 12.6, 12.4, 12.2, 12.0, 11.8,
	// Row 3 (1640 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 12.9, 12.7, 12.5, 12.3, 12.1, 11.9, 11.7,
	// Row 4 (2020 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 12.8, 12.6, 12.4, 12.2, 12.0, 11.8, 11.6,
	// Row 5 (2400 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 12.7, 12.5, 12.3, 12.1, 11.9, 11.7, 11.5,
	// Row 6 (2780 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 12.6, 12.4, 12.2, 12.0, 11.9, 11.7, 11.5,
	// Row 7 (3160 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 12.7, 12.5, 12.3, 12.1, 11.9, 11.7, 11.6,
	// Row 8 (3540 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 12.8, 12.6, 12.4, 12.2, 12.0, 11.9, 11.7,
	// Row 9 (3920 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 12.9, 12.7, 12.5, 12.3, 12.1, 11.9, 11.7,
	// Row 10 (4300 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 13.0, 12.8, 12.6, 12.4, 12.2, 12.0, 11.8,
	// Row 11 (4680 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 13.1, 12.9, 12.7, 12.5, 12.3, 12.1, 11.9,
	// Row 12 (5060 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 13.2, 13.0, 12.8, 12.6, 12.4, 12.2, 12.0,
	// Row 13 (5440 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 13.3, 13.1, 12.9, 12.7, 12.5, 12.3, 12.1,
	// Row 14 (5820 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 13.4, 13.2, 13.0, 12.8, 12.6, 12.4, 12.2,
	// Row 15 (6200 rpm)
	14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 14.7, 13.4, 13.2, 13.0, 12.9, 12.7, 12.5, 12.3,
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

// GetDefaultTargetAFR returns a copy of the hardcoded Target AFR map scaled for the given ethanol percentage.
func GetDefaultTargetAFR(ethanolPct float64) []float64 {
	stoich := EthanolStoichAFR(ethanolPct)
	ratio := stoich / 14.70
	out := make([]float64, len(HardcodedBaseTargetAFR))
	for i, v := range HardcodedBaseTargetAFR {
		out[i] = math.Round(v*ratio*10) / 10
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
	ZSeries     string    // "Lambda.LambdaInt" or "Lambda.External"
	LearnedZ    []float64 // MatrixBuilder learned values (same length as CurrentFuel)
	Counts      []int     // MatrixBuilder sample counts per cell
	CurrentFuel []float64 // Base BFuelCal.Map values
	TargetAFR   []float64 // Editable Target AFR table
	ClosedLoop  []bool    // true = closed loop, false = open loop
	EthanolPct  float64   // 0 .. 85 (default 0)
	Smoothing   float64   // 0.0 .. 1.0 (e.g. 0.5 for 50%)
}

// FuelAdjustmentResult contains the computed new fuel map and adjustment statistics.
type FuelAdjustmentResult struct {
	NewFuel      []float64
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

// CalculateFuelAdjustment calculates the new BFuelCal.Map:
// - If Lambda.LambdaInt is evaluated: ONLY closed loop cells are updated.
// - If Lambda.External is evaluated: ONLY open loop cells are updated.
func CalculateFuelAdjustment(cfg FuelAdjustmentConfig) FuelAdjustmentResult {
	n := len(cfg.CurrentFuel)
	res := FuelAdjustmentResult{
		NewFuel: make([]float64, n),
	}
	copy(res.NewFuel, cfg.CurrentFuel)

	isLambdaInt := IsClosedLoopSeries(cfg.ZSeries)
	isExternal := IsOpenLoopSeries(cfg.ZSeries)

	stoich := EthanolStoichAFR(cfg.EthanolPct)
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
			tgtAFR := stoich
			if i < len(cfg.TargetAFR) && cfg.TargetAFR[i] > 0 {
				tgtAFR = cfg.TargetAFR[i]
			}
			tgtLam := tgtAFR / stoich
			if tgtLam <= 0 {
				res.SkippedCells++
				continue
			}
			correctionRatio = learned / tgtLam
		} else {
			// Other/generic series: skip
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
		newFuel := math.Round(curFuel*smoothedRatio*100) / 100

		delta := newFuel - curFuel
		res.NewFuel[i] = newFuel
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

	if res.UpdatedCells > 0 {
		res.AvgDelta = totalDelta / float64(res.UpdatedCells)
	}

	return res
}

// SerializeT7FuelString serializes map data into the T7Suite clipboard string format:
// "20:0:val:~1:0:val:~...x:y:val:~"
// Preserves the historical '20' prefix on the first coordinate for T7Suite compatibility.
func SerializeT7FuelString(zData []float64, cols, rows int) string {
	tokens := make([]string, 0, len(zData))
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := r*cols + c
			if idx >= len(zData) {
				continue
			}
			val := zData[idx]
			outX := c
			if r == 0 && c == 0 {
				outX = 20
			}
			var valStr string
			if val == math.Trunc(val) {
				valStr = strconv.Itoa(int(val))
			} else {
				valStr = strconv.FormatFloat(val, 'f', 2, 64)
			}
			tokens = append(tokens, fmt.Sprintf("%d:%d:%s", outX, r, valStr))
		}
	}
	return strings.Join(tokens, ":~") + ":~"
}

// FuelUpdater is the interface provided by MainWindow to read/update BFuelCal.Map and its closed loop region.
type FuelUpdater interface {
	GetBFuelCal() (xData, yData, zData []float64, err error)
	GetClosedLoopRegion(xData, yData []float64) []bool
	UpdateBFuelCal(zData []float64) error
	Log(msg string)
}

// FuelAdjusterWidget is the interactive UI component to tune BFuelCal.Map using MatrixBuilder results.
type FuelAdjusterWidget struct {
	widget.BaseWidget
	mb      *MatrixBuilder
	updater FuelUpdater

	ethanolPct float64
	smoothing  float64

	targetAFR []float64
	newFuel   []float64
	result    FuelAdjustmentResult

	// UI controls
	ethSlider     *widget.Slider
	ethLabel      *widget.Label
	smoothSlider  *widget.Slider
	smoothLabel   *widget.Label
	ruleBadge     *widget.Label
	statsLabel    *widget.Label
	statusLabel   *widget.Label
	targetViewer  *mapviewer.MapViewer
	previewViewer *mapviewer.MapViewer
	tabContainer  *container.AppTabs

	content fyne.CanvasObject
}

// NewFuelAdjusterWidget creates a new fuel adjuster interface.
func NewFuelAdjusterWidget(mb *MatrixBuilder, updater FuelUpdater) *FuelAdjusterWidget {
	fa := &FuelAdjusterWidget{
		mb:         mb,
		updater:    updater,
		ethanolPct: 0,
		smoothing:  0.50,
		targetAFR:  GetDefaultTargetAFR(0),
	}
	fa.ExtendBaseWidget(fa)
	fa.buildUI()
	fa.recalculate()
	return fa
}

func (fa *FuelAdjusterWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(fa.content)
}

func (fa *FuelAdjusterWidget) buildUI() {
	// Ethanol slider (0..85%)
	fa.ethLabel = widget.NewLabel("0% (Stoich: 14.70)")
	fa.ethSlider = widget.NewSlider(0, 85)
	fa.ethSlider.Step = 1
	fa.ethSlider.SetValue(fa.ethanolPct)
	fa.ethSlider.OnChanged = func(val float64) {
		fa.ethanolPct = val
		stoich := EthanolStoichAFR(val)
		fa.ethLabel.SetText(fmt.Sprintf("%d%% (Stoich: %.2f)", int(val), stoich))
		// Auto-rescale the target AFR baseline when ethanol content changes
		fa.targetAFR = GetDefaultTargetAFR(val)
		if fa.targetViewer != nil {
			_ = fa.targetViewer.SetZData(fa.targetAFR)
			fa.targetViewer.Refresh()
		}
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

	fa.statsLabel = widget.NewLabel("")
	fa.statusLabel = widget.NewLabel("")

	resetAFRBtn := widget.NewButtonWithIcon("Reset Target AFR", theme.ViewRefreshIcon(), func() {
		fa.targetAFR = GetDefaultTargetAFR(fa.ethanolPct)
		if fa.targetViewer != nil {
			_ = fa.targetViewer.SetZData(fa.targetAFR)
			fa.targetViewer.Refresh()
		}
		fa.recalculate()
		fa.statusLabel.SetText("Reset Target AFR map to default")
	})

	topControls := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Evaluated Series:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(fa.mb.zSeries),
			layout.NewSpacer(),
			fa.ruleBadge,
		),
		container.NewGridWithColumns(2,
			container.NewBorder(nil, nil, widget.NewLabel("Ethanol:"), fa.ethLabel, fa.ethSlider),
			container.NewBorder(nil, nil, widget.NewLabel("Smoothing:"), fa.smoothLabel, fa.smoothSlider),
		),
		container.NewHBox(
			resetAFRBtn,
			layout.NewSpacer(),
			fa.statsLabel,
		),
		widget.NewSeparator(),
	)

	// Create Target AFR MapViewer (Editable)
	tCfg := &mapviewer.Config{
		Name:           "Target AFR",
		XData:          DefaultAirXSP,
		YData:          DefaultRpmYSP,
		ZData:          fa.targetAFR,
		XPrecision:     0,
		YPrecision:     0,
		ZPrecision:     1,
		XLabel:         "Airmass (mg/c)",
		YLabel:         "RPM",
		ZLabel:         "Target AFR",
		Editable:       true,
		ColorblindMode: colors.ModeNormal,
		OnUpdateCell: func(idx int, data []float64) {
			if idx >= 0 && idx < len(fa.targetAFR) {
				fa.targetAFR[idx] = data[idx]
				fa.recalculate()
			}
		},
	}
	tMv, _ := mapviewer.New(tCfg)
	fa.targetViewer = tMv

	// Create Preview MapViewer
	pCfg := &mapviewer.Config{
		Name:           "Adjusted BFuelCal.Map",
		XData:          DefaultAirXSP,
		YData:          DefaultRpmYSP,
		ZData:          make([]float64, len(HardcodedBaseTargetAFR)),
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

	fa.tabContainer = container.NewAppTabs(
		container.NewTabItemWithIcon("Target AFR Map (Editable)", theme.DocumentIcon(), fa.targetViewer),
		container.NewTabItemWithIcon("Adjusted BFuelCal.Map Preview", theme.GridIcon(), fa.previewViewer),
	)

	// Bottom action buttons
	applyBtn := widget.NewButtonWithIcon("Apply to BFuelCal.Map", theme.DocumentSaveIcon(), func() {
		fa.applyToFuelMap()
	})
	applyBtn.Importance = widget.HighImportance

	copyT7Btn := widget.NewButtonWithIcon("Copy T7Suite String", theme.ContentCopyIcon(), func() {
		fa.copyT7SuiteString()
	})

	bottomBar := container.NewBorder(nil, nil,
		fa.statusLabel,
		container.NewHBox(copyT7Btn, applyBtn),
	)

	fa.content = container.NewBorder(topControls, bottomBar, nil, nil, fa.tabContainer)
}

func (fa *FuelAdjusterWidget) getClosedLoopMask(xData, yData []float64) []bool {
	if fa.updater != nil {
		if mask := fa.updater.GetClosedLoopRegion(xData, yData); len(mask) == len(xData)*len(yData) {
			return mask
		}
	}
	return GetClosedLoopMask(xData, yData, DefaultClosedLoopRpm, DefaultClosedLoopMaxLoad)
}

func (fa *FuelAdjusterWidget) getBaseFuelMap() (xData, yData, zData []float64) {
	if fa.updater != nil {
		if x, y, z, err := fa.updater.GetBFuelCal(); err == nil && len(z) == len(x)*len(y) {
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

	// Update rule badge
	if IsClosedLoopSeries(zSeries) {
		fa.ruleBadge.SetText("Mode: CLOSED-LOOP ONLY (updating cells where Airmass <= MaxLoad)")
	} else if IsOpenLoopSeries(zSeries) {
		fa.ruleBadge.SetText("Mode: OPEN-LOOP ONLY (updating cells where Airmass > MaxLoad)")
	} else {
		fa.ruleBadge.SetText("Mode: Neutral (" + zSeries + ")")
	}

	xData, yData, baseZ := fa.getBaseFuelMap()
	mask := fa.getClosedLoopMask(xData, yData)

	learnedZ := fa.mb.zData
	counts := fa.mb.counts

	// If matrix dimensions differ, create zero counts
	if len(learnedZ) != len(baseZ) {
		learnedZ = make([]float64, len(baseZ))
		counts = make([]int, len(baseZ))
	}

	cfg := FuelAdjustmentConfig{
		ZSeries:     zSeries,
		LearnedZ:    learnedZ,
		Counts:      counts,
		CurrentFuel: baseZ,
		TargetAFR:   fa.targetAFR,
		ClosedLoop:  mask,
		EthanolPct:  fa.ethanolPct,
		Smoothing:   fa.smoothing,
	}

	res := CalculateFuelAdjustment(cfg)
	fa.result = res
	fa.newFuel = res.NewFuel

	if fa.previewViewer != nil {
		_ = fa.previewViewer.SetZData(res.NewFuel)
		fa.previewViewer.Refresh()
	}

	fa.statsLabel.SetText(fmt.Sprintf("Cells Updated: %d | Skipped: %d | Avg Delta: %+.2f",
		res.UpdatedCells, res.SkippedCells, res.AvgDelta))
}

func (fa *FuelAdjusterWidget) applyToFuelMap() {
	if fa.updater == nil {
		fa.statusLabel.SetText("No active ECU binary connection to update BFuelCal.Map")
		return
	}
	if err := fa.updater.UpdateBFuelCal(fa.newFuel); err != nil {
		fa.statusLabel.SetText("Error updating BFuelCal.Map: " + err.Error())
		return
	}
	fa.statusLabel.SetText(fmt.Sprintf("Successfully updated %d cells in BFuelCal.Map!", fa.result.UpdatedCells))
}

func (fa *FuelAdjusterWidget) copyT7SuiteString() {
	xData, yData, _ := fa.getBaseFuelMap()
	str := SerializeT7FuelString(fa.newFuel, len(xData), len(yData))
	if app := fyne.CurrentApp(); app != nil {
		if drv := app.Driver(); drv != nil {
			windows := drv.AllWindows()
			if len(windows) > 0 {
				windows[0].Clipboard().SetContent(str)
				fa.statusLabel.SetText("Copied adjusted BFuelCal.Map to clipboard (T7Suite format)")
				return
			}
		}
	}
	fa.statusLabel.SetText("Failed to access system clipboard")
}
