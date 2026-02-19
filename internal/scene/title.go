package scene

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"neonsigil/internal/config"
	"neonsigil/internal/ui"
)

// DrawTitleScreen draws the title screen
func DrawTitleScreen(screen *ebiten.Image, tick int) {
	screen.Fill(config.ColorBG)

	// Animated scanlines
	for y := 0; y < config.ScreenHeight; y += 4 {
		alpha := uint8(8 + int(4*math.Sin(float64(y+tick)/20.0)))
		vector.DrawFilledRect(screen, 0, float32(y), config.ScreenWidth, 1, color.RGBA{0, 255, 255, alpha}, false)
	}

	// Title glow background
	pulse := math.Sin(float64(tick%120)/120.0*math.Pi*2)*0.3 + 0.7
	titleY := float64(config.ScreenHeight/2 - 80)

	// Large glow behind title
	glowAlpha := uint8(float64(20) * pulse)
	vector.DrawFilledRect(screen, float32(config.ScreenWidth/2-300), float32(titleY-30),
		600, 80, color.RGBA{0, 255, 255, glowAlpha}, false)

	// Main title
	ui.DrawTextGlowCentered(screen, "NEON SIGIL", ui.FontBold(48), config.ScreenWidth/2, titleY, config.ColorNeonCyan)

	// Subtitle
	ui.DrawTextGlowCentered(screen, "TRI-FUSE DEFENSE", ui.FontBold(20), config.ScreenWidth/2, titleY+60, config.ColorNeonMagenta)

	// Animated border lines
	lineAlpha := uint8(float64(100) * pulse)
	lineColor := color.RGBA{0, 255, 255, lineAlpha}
	lineW := float32(400)
	vector.DrawFilledRect(screen, float32(config.ScreenWidth/2)-lineW/2, float32(titleY+90), lineW, 2, lineColor, false)

	// Start prompt
	if tick%80 < 60 {
		ui.DrawTextCentered(screen, "PRESS ENTER TO START", ui.FontRegular(14), config.ScreenWidth/2, config.ScreenHeight/2+100, config.ColorWhite)
	}

	// Version
	ui.DrawText(screen, "v0.1 MVP", ui.FontRegular(9), 20, float64(config.ScreenHeight-30), config.ColorWhiteDim)

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
	vector.StrokeLine(screen, config.ScreenWidth-80, 20, config.ScreenWidth-20, 20, 2, c, false)
	vector.StrokeLine(screen, config.ScreenWidth-20, 20, config.ScreenWidth-20, 80, 2, c, false)

	// Bottom-left
	vector.StrokeLine(screen, 20, config.ScreenHeight-20, 80, config.ScreenHeight-20, 2, c, false)
	vector.StrokeLine(screen, 20, config.ScreenHeight-80, 20, config.ScreenHeight-20, 2, c, false)

	// Bottom-right
	vector.StrokeLine(screen, config.ScreenWidth-80, config.ScreenHeight-20, config.ScreenWidth-20, config.ScreenHeight-20, 2, c, false)
	vector.StrokeLine(screen, config.ScreenWidth-20, config.ScreenHeight-80, config.ScreenWidth-20, config.ScreenHeight-20, 2, c, false)
}
