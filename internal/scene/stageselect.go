package scene

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"neonsigil/internal/config"
	"neonsigil/internal/data"
	"neonsigil/internal/ui"
)

// StageSelectState manages the stage selection screen
type StageSelectState struct {
	Selected int
	Tick     int
}

// NewStageSelectState creates a new stage select state
func NewStageSelectState() *StageSelectState {
	return &StageSelectState{Selected: 0}
}

// Update handles input and returns the selected stage index, or -1 if none
func (s *StageSelectState) Update() int {
	s.Tick++

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.Selected--
		if s.Selected < 0 {
			s.Selected = len(data.Stages) - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.Selected++
		if s.Selected >= len(data.Stages) {
			s.Selected = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return s.Selected
	}

	// Mouse selection
	mx, my := ebiten.CursorPosition()
	if mx >= 340 && mx <= 940 {
		for i := range data.Stages {
			itemY := 140 + i*50
			if my >= itemY && my < itemY+44 {
				s.Selected = i
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					return s.Selected
				}
				break
			}
		}
	}

	return -1 // no selection
}

// Draw renders the stage selection screen
func (s *StageSelectState) Draw(screen *ebiten.Image) {
	screen.Fill(config.ColorBG)

	// Scanlines
	for y := 0; y < config.ScreenHeight; y += 6 {
		alpha := uint8(5 + int(3*math.Sin(float64(y+s.Tick)/25.0)))
		vector.DrawFilledRect(screen, 0, float32(y), config.ScreenWidth, 1, color.RGBA{255, 0, 255, alpha}, false)
	}

	// Header
	ui.DrawTextGlowCentered(screen, "SELECT STAGE", ui.FontBold(28), config.ScreenWidth/2, 60, config.ColorNeonCyan)

	// Stage list
	for i, stage := range data.Stages {
		y := float64(140 + i*50)
		x := 340.0
		w := 600.0
		h := 44.0

		isSelected := i == s.Selected

		// Background
		bgColor := color.RGBA{15, 15, 35, 200}
		borderColor := color.RGBA{40, 40, 70, 150}
		if isSelected {
			bgColor = color.RGBA{20, 25, 50, 240}
			borderColor = config.ColorNeonCyan
			// Glow
			pulse := math.Sin(float64(s.Tick%40)/40.0*math.Pi*2)*0.3 + 0.7
			glowAlpha := uint8(float64(30) * pulse)
			vector.DrawFilledRect(screen, float32(x-3), float32(y-3), float32(w+6), float32(h+6),
				color.RGBA{0, 255, 255, glowAlpha}, false)
		}

		vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), bgColor, false)
		vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 1.5, borderColor, false)

		// Stage number
		numStr := fmt.Sprintf("%02d", i+1)
		numColor := config.ColorNeonMagenta
		if isSelected {
			numColor = config.ColorNeonCyan
		}
		ui.DrawText(screen, numStr, ui.FontBold(16), x+12, y+12, numColor)

		// Stage name
		nameColor := config.ColorWhite
		if isSelected {
			nameColor = config.ColorNeonCyan
		}
		ui.DrawText(screen, stage.Name, ui.FontBold(13), x+60, y+13, nameColor)

		// Stage ID
		ui.DrawText(screen, stage.ID, ui.FontRegular(9), x+w-80, y+16, config.ColorWhiteDim)

		// Difficulty indicator (integrity)
		intStr := fmt.Sprintf("HP:%d", stage.Integrity)
		ui.DrawText(screen, intStr, ui.FontRegular(8), x+60, y+30, config.ColorWhiteDim)

		// Waves
		waveStr := fmt.Sprintf("W:%d", len(stage.Waves))
		ui.DrawText(screen, waveStr, ui.FontRegular(8), x+140, y+30, config.ColorWhiteDim)
	}

	// Navigation hint
	ui.DrawTextCentered(screen, "UP/DOWN to select, ENTER to start, ESC to go back",
		ui.FontRegular(10), config.ScreenWidth/2, float64(config.ScreenHeight-40), config.ColorWhiteDim)

	drawCornerDecor(screen, s.Tick)
}
