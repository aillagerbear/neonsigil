package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TitleScreen draws the title screen
func DrawTitleScreen(screen *ebiten.Image, tick int) {
	screen.Fill(ColorBG)

	// Animated scanlines
	for y := 0; y < ScreenHeight; y += 4 {
		alpha := uint8(8 + int(4*math.Sin(float64(y+tick)/20.0)))
		vector.DrawFilledRect(screen, 0, float32(y), ScreenWidth, 1, color.RGBA{0, 255, 255, alpha}, false)
	}

	// Title glow background
	pulse := math.Sin(float64(tick%120)/120.0*math.Pi*2)*0.3 + 0.7
	titleY := float64(ScreenHeight/2 - 80)

	// Large glow behind title
	glowAlpha := uint8(float64(20) * pulse)
	vector.DrawFilledRect(screen, float32(ScreenWidth/2-300), float32(titleY-30),
		600, 80, color.RGBA{0, 255, 255, glowAlpha}, false)

	// Main title
	DrawTextGlowCentered(screen, "NEON SIGIL", fontBold(48), ScreenWidth/2, titleY, ColorNeonCyan)

	// Subtitle
	DrawTextGlowCentered(screen, "TRI-FUSE DEFENSE", fontBold(20), ScreenWidth/2, titleY+60, ColorNeonMagenta)

	// Animated border lines
	lineAlpha := uint8(float64(100) * pulse)
	lineColor := color.RGBA{0, 255, 255, lineAlpha}
	lineW := float32(400)
	vector.DrawFilledRect(screen, float32(ScreenWidth/2)-lineW/2, float32(titleY+90), lineW, 2, lineColor, false)

	// Start prompt
	if tick%80 < 60 {
		DrawTextCentered(screen, "PRESS ENTER TO START", fontRegular(14), ScreenWidth/2, ScreenHeight/2+100, ColorWhite)
	}

	// Version
	DrawText(screen, "v0.1 MVP", fontRegular(9), 20, float64(ScreenHeight-30), ColorWhiteDim)

	// Decorative corner elements
	drawCornerDecor(screen, tick)
}

func drawCornerDecor(screen *ebiten.Image, tick int) {
	pulse := math.Sin(float64(tick%90)/90.0*math.Pi*2)*0.3 + 0.7
	alpha := uint8(float64(80) * pulse)
	c := color.RGBA{0, 255, 255, alpha}

	// Top-left
	vector.StrokeLine(screen, 20, 20, 80, 20, 2, c, false)
	vector.StrokeLine(screen, 20, 20, 20, 80, 2, c, false)

	// Top-right
	vector.StrokeLine(screen, ScreenWidth-80, 20, ScreenWidth-20, 20, 2, c, false)
	vector.StrokeLine(screen, ScreenWidth-20, 20, ScreenWidth-20, 80, 2, c, false)

	// Bottom-left
	vector.StrokeLine(screen, 20, ScreenHeight-20, 80, ScreenHeight-20, 2, c, false)
	vector.StrokeLine(screen, 20, ScreenHeight-80, 20, ScreenHeight-20, 2, c, false)

	// Bottom-right
	vector.StrokeLine(screen, ScreenWidth-80, ScreenHeight-20, ScreenWidth-20, ScreenHeight-20, 2, c, false)
	vector.StrokeLine(screen, ScreenWidth-20, ScreenHeight-80, ScreenWidth-20, ScreenHeight-20, 2, c, false)
}

// StageSelectScreen
type StageSelectState struct {
	Selected int
	Tick     int
}

func NewStageSelectState() *StageSelectState {
	return &StageSelectState{Selected: 0}
}

func (s *StageSelectState) Update() int {
	s.Tick++

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.Selected--
		if s.Selected < 0 {
			s.Selected = len(Stages) - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.Selected++
		if s.Selected >= len(Stages) {
			s.Selected = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return s.Selected
	}

	// Mouse selection
	mx, my := ebiten.CursorPosition()
	if mx >= 340 && mx <= 940 {
		for i := range Stages {
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

func (s *StageSelectState) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBG)

	// Scanlines
	for y := 0; y < ScreenHeight; y += 6 {
		alpha := uint8(5 + int(3*math.Sin(float64(y+s.Tick)/25.0)))
		vector.DrawFilledRect(screen, 0, float32(y), ScreenWidth, 1, color.RGBA{255, 0, 255, alpha}, false)
	}

	// Header
	DrawTextGlowCentered(screen, "SELECT STAGE", fontBold(28), ScreenWidth/2, 60, ColorNeonCyan)

	// Stage list
	for i, stage := range Stages {
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
			borderColor = ColorNeonCyan
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
		numColor := ColorNeonMagenta
		if isSelected {
			numColor = ColorNeonCyan
		}
		DrawText(screen, numStr, fontBold(16), x+12, y+12, numColor)

		// Stage name
		nameColor := ColorWhite
		if isSelected {
			nameColor = ColorNeonCyan
		}
		DrawText(screen, stage.Name, fontBold(13), x+60, y+13, nameColor)

		// Stage ID
		DrawText(screen, stage.ID, fontRegular(9), x+w-80, y+16, ColorWhiteDim)

		// Difficulty indicator (integrity)
		intStr := fmt.Sprintf("HP:%d", stage.Integrity)
		DrawText(screen, intStr, fontRegular(8), x+60, y+30, ColorWhiteDim)

		// Waves
		waveStr := fmt.Sprintf("W:%d", len(stage.Waves))
		DrawText(screen, waveStr, fontRegular(8), x+140, y+30, ColorWhiteDim)
	}

	// Navigation hint
	DrawTextCentered(screen, "UP/DOWN to select, ENTER to start, ESC to go back",
		fontRegular(10), ScreenWidth/2, float64(ScreenHeight-40), ColorWhiteDim)

	drawCornerDecor(screen, s.Tick)
}
