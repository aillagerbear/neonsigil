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
	vector.FillCircle(screen, fx, fy+r*0.14, r*0.96, color.RGBA{0x01, 0x05, 0x0b, 0x88}, false)

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
	shadowOp.ColorScale.ScaleWithColor(color.RGBA{0x03, 0x08, 0x12, 0xB8})
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
			Outer: color.RGBA{0x2F, 0x5F, 0x89, 0xFF},
			Inner: color.RGBA{0x75, 0xB7, 0xE4, 0xFF},
			Ring:  color.RGBA{0x9B, 0xD6, 0xFF, 0xFF},
			Glow:  color.RGBA{0x7F, 0xC7, 0xF2, 0x7F},
			Text:  color.RGBA{0xF3, 0xF8, 0xFF, 0xFF},
			Label: "S",
		}
	case entity.CardArcher:
		return avatarStyle{
			Outer: color.RGBA{0x2C, 0x6F, 0x62, 0xFF},
			Inner: color.RGBA{0x6F, 0xD8, 0xBC, 0xFF},
			Ring:  color.RGBA{0x98, 0xF2, 0xD4, 0xFF},
			Glow:  color.RGBA{0x72, 0xE9, 0xCA, 0x7F},
			Text:  color.RGBA{0xF0, 0xFF, 0xFA, 0xFF},
			Label: "A",
		}
	case entity.CardSpearman:
		return avatarStyle{
			Outer: color.RGBA{0x6E, 0x52, 0x2A, 0xFF},
			Inner: color.RGBA{0xEC, 0xBE, 0x6A, 0xFF},
			Ring:  color.RGBA{0xFF, 0xDC, 0x95, 0xFF},
			Glow:  color.RGBA{0xFF, 0xCE, 0x72, 0x7A},
			Text:  color.RGBA{0xFF, 0xFB, 0xEE, 0xFF},
			Label: "P",
		}
	case entity.CardMage:
		return avatarStyle{
			Outer: color.RGBA{0x4A, 0x3F, 0x7E, 0xFF},
			Inner: color.RGBA{0xA7, 0x92, 0xEE, 0xFF},
			Ring:  color.RGBA{0xC3, 0xB5, 0xFF, 0xFF},
			Glow:  color.RGBA{0xB6, 0xA5, 0xFF, 0x7A},
			Text:  color.RGBA{0xF5, 0xF0, 0xFF, 0xFF},
			Label: "M",
		}
	default:
		return avatarStyle{
			Outer: color.RGBA{0x7C, 0x3E, 0x2A, 0xFF},
			Inner: color.RGBA{0xFF, 0x8D, 0x62, 0xFF},
			Ring:  color.RGBA{0xFF, 0xB1, 0x87, 0xFF},
			Glow:  color.RGBA{0xFF, 0x9D, 0x6E, 0x86},
			Text:  color.RGBA{0xFF, 0xF4, 0xED, 0xFF},
			Label: "F",
		}
	}
}

func enemyAvatarStyle(enemyType entity.EnemyType) avatarStyle {
	switch enemyType {
	case entity.EnemyGoblin:
		return avatarStyle{
			Outer: color.RGBA{0x3A, 0x72, 0x53, 0xFF},
			Inner: color.RGBA{0x7F, 0xD8, 0x96, 0xFF},
			Ring:  color.RGBA{0x9F, 0xEB, 0xB8, 0xFF},
			Glow:  color.RGBA{0x8C, 0xE6, 0xAA, 0x7A},
			Text:  color.RGBA{0xF5, 0xFF, 0xF8, 0xFF},
			Label: "G",
		}
	case entity.EnemyOrc:
		return avatarStyle{
			Outer: color.RGBA{0x7A, 0x4C, 0x3A, 0xFF},
			Inner: color.RGBA{0xD7, 0x94, 0x74, 0xFF},
			Ring:  color.RGBA{0xE9, 0xB1, 0x95, 0xFF},
			Glow:  color.RGBA{0xD6, 0x9E, 0x82, 0x74},
			Text:  color.RGBA{0xFF, 0xF6, 0xEF, 0xFF},
			Label: "O",
		}
	case entity.EnemyBossOrc:
		return avatarStyle{
			Outer: color.RGBA{0x8A, 0x3E, 0x34, 0xFF},
			Inner: color.RGBA{0xFF, 0x86, 0x6C, 0xFF},
			Ring:  color.RGBA{0xFF, 0xAB, 0x92, 0xFF},
			Glow:  color.RGBA{0xFF, 0x8E, 0x78, 0x88},
			Text:  color.RGBA{0xFF, 0xF4, 0xEF, 0xFF},
			Label: "B",
		}
	default:
		return avatarStyle{
			Outer: color.RGBA{0x5E, 0x2F, 0x6A, 0xFF},
			Inner: color.RGBA{0xC9, 0x84, 0xDB, 0xFF},
			Ring:  color.RGBA{0xE0, 0xAC, 0xF0, 0xFF},
			Glow:  color.RGBA{0xCF, 0x93, 0xE1, 0x8E},
			Text:  color.RGBA{0xFF, 0xF2, 0xFF, 0xFF},
			Label: "X",
		}
	}
}
