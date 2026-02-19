package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"neonsigil/internal/config"
)

// Button represents a clickable UI button
type Button struct {
	X, Y, W, H float64
	Label       string
	Color       color.RGBA
	Disabled    bool
	Hovered     bool
}

// Contains checks if the given screen coordinates are within the button
func (b *Button) Contains(mx, my int) bool {
	return float64(mx) >= b.X && float64(mx) <= b.X+b.W &&
		float64(my) >= b.Y && float64(my) <= b.Y+b.H
}

// Draw renders the button
func (b *Button) Draw(screen *ebiten.Image, tick int) {
	bgColor := config.ColorBtnBG
	borderColor := b.Color
	if b.Disabled {
		bgColor = color.RGBA{20, 20, 30, 200}
		borderColor = color.RGBA{60, 60, 80, 150}
	} else if b.Hovered {
		bgColor = config.ColorBtnHover
		// Glow effect
		pulse := math.Sin(float64(tick%40)/40.0*math.Pi*2)*0.3 + 0.7
		alpha := uint8(float64(40) * pulse)
		vector.DrawFilledRect(screen, float32(b.X-2), float32(b.Y-2), float32(b.W+4), float32(b.H+4),
			color.RGBA{b.Color.R, b.Color.G, b.Color.B, alpha}, false)
	}

	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), bgColor, false)
	vector.StrokeRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 1.5, borderColor, false)

	textColor := borderColor
	if b.Disabled {
		textColor = color.RGBA{80, 80, 100, 200}
	}
	DrawTextCentered(screen, b.Label, FontBold(11), b.X+b.W/2, b.Y+b.H/2, textColor)
}
