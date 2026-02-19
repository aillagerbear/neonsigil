package ui

import (
	"bytes"
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets/fonts/Orbitron-Bold.ttf
var orbitronBoldTTF []byte

//go:embed assets/fonts/Orbitron-Regular.ttf
var orbitronRegularTTF []byte

var (
	fontSourceBold    *text.GoTextFaceSource
	fontSourceRegular *text.GoTextFaceSource
)

// InitFonts loads the embedded font files
func InitFonts() {
	var err error
	fontSourceBold, err = text.NewGoTextFaceSource(bytes.NewReader(orbitronBoldTTF))
	if err != nil {
		panic(err)
	}
	fontSourceRegular, err = text.NewGoTextFaceSource(bytes.NewReader(orbitronRegularTTF))
	if err != nil {
		panic(err)
	}
}

// FontBold returns a bold font face at the given size
func FontBold(size float64) *text.GoTextFace {
	return &text.GoTextFace{Source: fontSourceBold, Size: size}
}

// FontRegular returns a regular font face at the given size
func FontRegular(size float64) *text.GoTextFace {
	return &text.GoTextFace{Source: fontSourceRegular, Size: size}
}

// DrawText draws text at the given position
func DrawText(screen *ebiten.Image, str string, face *text.GoTextFace, x, y float64, clr color.RGBA) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, str, face, op)
}

// DrawTextCentered draws centered text
func DrawTextCentered(screen *ebiten.Image, str string, face *text.GoTextFace, cx, cy float64, clr color.RGBA) {
	w, h := text.Measure(str, face, 0)
	DrawText(screen, str, face, cx-w/2, cy-h/2, clr)
}

// DrawTextGlow draws text with neon glow effect
func DrawTextGlow(screen *ebiten.Image, str string, face *text.GoTextFace, x, y float64, clr color.RGBA) {
	glowClr := color.RGBA{clr.R, clr.G, clr.B, 40}
	for _, off := range []float64{-2, 2} {
		DrawText(screen, str, face, x+off, y, glowClr)
		DrawText(screen, str, face, x, y+off, glowClr)
	}
	DrawText(screen, str, face, x, y, clr)
}

// DrawTextGlowCentered draws centered text with glow
func DrawTextGlowCentered(screen *ebiten.Image, str string, face *text.GoTextFace, cx, cy float64, clr color.RGBA) {
	w, h := text.Measure(str, face, 0)
	DrawTextGlow(screen, str, face, cx-w/2, cy-h/2, clr)
}
