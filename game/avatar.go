package game

import (
	"image/color"

	"ebitengine-testing/entity"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type avatarStyle struct {
	Outer color.RGBA
	Inner color.RGBA
	Ring  color.RGBA
	Glow  color.RGBA
	Text  color.RGBA
	Label string
}

func avatarFaceForRadius(radius float64) text.Face {
	if radius >= 30 {
		return fontLarge
	}
	if radius >= 22 {
		return fontCardCost
	}
	return fontCardName
}

func drawAvatarBadge(screen *ebiten.Image, x, y, radius float64, style avatarStyle) {
	if radius <= 0 {
		return
	}
	r := float32(radius)
	fx := float32(x)
	fy := float32(y)

	// Soft drop shadow
	vector.FillCircle(screen, fx, fy+r*0.14, r*0.96, color.RGBA{0xA1, 0x8F, 0xC1, 0x72}, false)

	// Main body
	vector.FillCircle(screen, fx, fy, r, style.Outer, false)
	vector.FillCircle(screen, fx, fy-r*0.15, r*0.72, style.Inner, false)
	vector.FillCircle(screen, fx-r*0.30, fy-r*0.34, r*0.20, color.RGBA{0xFF, 0xFF, 0xFF, 0x2E}, false)

	// Ring and glow
	vector.StrokeCircle(screen, fx, fy, r, 2, style.Ring, false)
	if style.Glow.A > 0 {
		vector.StrokeCircle(screen, fx, fy, r+2, 1, style.Glow, false)
	}

	// Center label
	face := avatarFaceForRadius(radius)
	lw, lh := text.Measure(style.Label, face, 0)
	shadowOp := &text.DrawOptions{}
	shadowOp.GeoM.Translate(x-lw/2+1, y-lh/2+1)
	shadowOp.ColorScale.ScaleWithColor(color.RGBA{0x75, 0x66, 0x96, 0xA8})
	text.Draw(screen, style.Label, face, shadowOp)

	op := &text.DrawOptions{}
	op.GeoM.Translate(x-lw/2, y-lh/2)
	op.ColorScale.ScaleWithColor(style.Text)
	text.Draw(screen, style.Label, face, op)
}

func cardAvatarStyle(cardType entity.CardType) avatarStyle {
	switch cardType {
	case entity.CardSoldier:
		return avatarStyle{
			Outer: color.RGBA{0xA9, 0xC5, 0xF5, 0xFF},
			Inner: color.RGBA{0xE8, 0xF3, 0xFF, 0xFF},
			Ring:  color.RGBA{0x8F, 0xAF, 0xE5, 0xFF},
			Glow:  color.RGBA{0xBF, 0xD9, 0xFF, 0x8A},
			Text:  color.RGBA{0x4B, 0x5F, 0x86, 0xFF},
			Label: "S",
		}
	case entity.CardArcher:
		return avatarStyle{
			Outer: color.RGBA{0x9D, 0xE8, 0xC7, 0xFF},
			Inner: color.RGBA{0xE8, 0xFF, 0xF4, 0xFF},
			Ring:  color.RGBA{0x72, 0xC8, 0xA8, 0xFF},
			Glow:  color.RGBA{0xB7, 0xF7, 0xDE, 0x8A},
			Text:  color.RGBA{0x3C, 0x7A, 0x67, 0xFF},
			Label: "A",
		}
	case entity.CardSpearman:
		return avatarStyle{
			Outer: color.RGBA{0xF7, 0xD4, 0xA3, 0xFF},
			Inner: color.RGBA{0xFF, 0xF3, 0xD9, 0xFF},
			Ring:  color.RGBA{0xE8, 0xB7, 0x6E, 0xFF},
			Glow:  color.RGBA{0xFF, 0xE1, 0xA9, 0x8A},
			Text:  color.RGBA{0x8D, 0x68, 0x37, 0xFF},
			Label: "P",
		}
	case entity.CardMage:
		return avatarStyle{
			Outer: color.RGBA{0xC9, 0xB7, 0xF9, 0xFF},
			Inner: color.RGBA{0xF2, 0xEC, 0xFF, 0xFF},
			Ring:  color.RGBA{0xA7, 0x90, 0xE8, 0xFF},
			Glow:  color.RGBA{0xDE, 0xD1, 0xFF, 0x86},
			Text:  color.RGBA{0x62, 0x53, 0x9B, 0xFF},
			Label: "M",
		}
	default:
		return avatarStyle{
			Outer: color.RGBA{0xF7, 0xBC, 0xA8, 0xFF},
			Inner: color.RGBA{0xFF, 0xEA, 0xE0, 0xFF},
			Ring:  color.RGBA{0xE8, 0x9D, 0x84, 0xFF},
			Glow:  color.RGBA{0xFF, 0xD2, 0xC2, 0x90},
			Text:  color.RGBA{0x8E, 0x55, 0x46, 0xFF},
			Label: "F",
		}
	}
}

func enemyAvatarStyle(enemyType entity.EnemyType) avatarStyle {
	switch enemyType {
	case entity.EnemyGoblin:
		return avatarStyle{
			Outer: color.RGBA{0xA1, 0xDA, 0xB1, 0xFF},
			Inner: color.RGBA{0xEB, 0xFB, 0xEF, 0xFF},
			Ring:  color.RGBA{0x78, 0xB9, 0x8A, 0xFF},
			Glow:  color.RGBA{0xC5, 0xF0, 0xCF, 0x86},
			Text:  color.RGBA{0x3C, 0x7A, 0x54, 0xFF},
			Label: "G",
		}
	case entity.EnemyOrc:
		return avatarStyle{
			Outer: color.RGBA{0xE4, 0xB1, 0xA0, 0xFF},
			Inner: color.RGBA{0xFF, 0xE8, 0xDF, 0xFF},
			Ring:  color.RGBA{0xCA, 0x8C, 0x79, 0xFF},
			Glow:  color.RGBA{0xF6, 0xD1, 0xC5, 0x84},
			Text:  color.RGBA{0x7E, 0x4E, 0x43, 0xFF},
			Label: "O",
		}
	case entity.EnemyBossOrc:
		return avatarStyle{
			Outer: color.RGBA{0xF2, 0x95, 0xAE, 0xFF},
			Inner: color.RGBA{0xFF, 0xDF, 0xE8, 0xFF},
			Ring:  color.RGBA{0xD7, 0x6B, 0x93, 0xFF},
			Glow:  color.RGBA{0xFF, 0xBF, 0xD4, 0x8A},
			Text:  color.RGBA{0x7F, 0x3A, 0x5D, 0xFF},
			Label: "B",
		}
	default:
		return avatarStyle{
			Outer: color.RGBA{0xC2, 0x9C, 0xE5, 0xFF},
			Inner: color.RGBA{0xF1, 0xE4, 0xFF, 0xFF},
			Ring:  color.RGBA{0x9F, 0x76, 0xC8, 0xFF},
			Glow:  color.RGBA{0xDC, 0xC3, 0xF5, 0x8E},
			Text:  color.RGBA{0x5B, 0x3B, 0x7C, 0xFF},
			Label: "X",
		}
	}
}
