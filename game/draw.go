package game

import (
	"image/color"
	"math"

	"ebitengine-testing/config"
	"ebitengine-testing/entity"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ---- Korean text drawing ----

func drawKoreanText(screen *ebiten.Image, s string, face text.Face, x, y float64, clr color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, s, face, op)
}

func drawKoreanTextWithShadow(screen *ebiten.Image, s string, face text.Face, x, y float64, clr color.Color) {
	// Shadow
	shadowOp := &text.DrawOptions{}
	shadowOp.GeoM.Translate(x+1, y+1)
	shadowOp.ColorScale.ScaleWithColor(color.RGBA{0x73, 0x64, 0x93, 0x8A})
	text.Draw(screen, s, face, shadowOp)
	// Main text
	drawKoreanText(screen, s, face, x, y, clr)
}

// ---- Battle Screen ----

func (g *Game) drawBattle(screen *ebiten.Image) {
	drawStars(screen, g.bgStars, g.animTick)
	g.drawPath(screen)
	g.drawGrid(screen)
	g.drawSummonerBase(screen)
	g.drawSummoners(screen)
	g.drawEnemies(screen)
	g.drawProjectiles(screen)
	g.drawParticles(screen)
}

// ---- Path ----

func (g *Game) drawPath(screen *ebiten.Image) {
	pathColor := color.RGBA{0xF2, 0xE0, 0xF8, 0xE8}
	pathBorder := color.RGBA{0xD2, 0xAE, 0xDE, 0xB8}
	pathWidth := float32(32)
	halfW := pathWidth / 2

	for i := 0; i < len(entity.EnemyPath)-1; i++ {
		x1, y1 := float32(entity.EnemyPath[i].X), float32(entity.EnemyPath[i].Y)
		x2, y2 := float32(entity.EnemyPath[i+1].X), float32(entity.EnemyPath[i+1].Y)

		var rx, ry, rw, rh float32
		if x1 == x2 { // vertical
			rx = x1 - halfW
			ry = float32(math.Min(float64(y1), float64(y2)))
			rw = pathWidth
			rh = float32(math.Abs(float64(y2-y1))) + pathWidth
		} else { // horizontal
			rx = float32(math.Min(float64(x1), float64(x2))) - halfW
			ry = y1 - halfW
			rw = float32(math.Abs(float64(x2-x1))) + pathWidth
			rh = pathWidth
		}

		vector.FillRect(screen, rx, ry, rw, rh, pathColor, false)
		vector.FillRect(screen, rx+1, ry+1, rw-2, rh*0.34, color.RGBA{0xFF, 0xFF, 0xFF, 0x66}, false)
		vector.StrokeRect(screen, rx, ry, rw, rh, 1, pathBorder, false)
	}

	// Animated energy dots along path
	drawPathEnergyDots(screen, g.animTick)
}

// ---- Grid ----

func (g *Game) drawGrid(screen *ebiten.Image) {
	mx, my := ebiten.CursorPosition()
	showPlacement := g.selectedCard >= 0

	for row := 0; row < config.GridRows; row++ {
		for col := 0; col < config.GridCols; col++ {
			x := float32(config.GridStartX + col*config.TileSize)
			y := float32(config.GridStartY + row*config.TileSize)
			cx := x + float32(config.TileSize)/2
			cy := y + float32(config.TileSize)/2

			if !showPlacement {
				// Keep only a tiny anchor point so the map is visible while no card is selected.
				vector.FillCircle(screen, cx, cy, 1.6, color.RGBA{0xCC, 0xB5, 0xE8, 0x70}, false)
				continue
			}

			occupied := g.grid[row][col] != nil
			tile := g.sprites.GrassTile
			if occupied {
				tile = g.sprites.OccupiedTile
			}

			// Hover effect
			isHover := mx >= int(x) && mx < int(x)+config.TileSize && my >= int(y) && my < int(y)+config.TileSize
			if isHover {
				if !occupied {
					tile = g.sprites.HoverTile
				} else {
					tile = g.sprites.BlockedTile
				}
			}

			// Draw tile texture only while in placement mode.
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			alpha := float32(0.56)
			if occupied {
				alpha = 0.48
			}
			if isHover {
				alpha = 0.8
			}
			op.ColorScale.ScaleAlpha(alpha)
			screen.DrawImage(tile, op)
			vector.StrokeRect(screen, x, y, config.TileSize, config.TileSize, 1, color.RGBA{0xD7, 0xB8, 0xE9, 0x82}, false)

			if isHover {
				hoverBorder := color.RGBA{0x86, 0xE0, 0xCC, 0xD8}
				if occupied {
					hoverBorder = color.RGBA{0xF7, 0xA0, 0xC6, 0xD8}
				}
				vector.StrokeRect(screen, x+1, y+1, config.TileSize-2, config.TileSize-2, 2, hoverBorder, false)
			}
		}
	}
}

// ---- Summoners (placed units) ----

func (g *Game) drawSummoners(screen *ebiten.Image) {
	for _, s := range g.summoners {
		sx := s.ScreenX
		sy := s.ScreenY

		radius := 14.0
		if s.AtkTimer == 0 {
			radius = 16.0
		}
		drawAvatarBadge(screen, sx, sy, radius, cardAvatarStyle(s.Card.Type))

		// Health bar (pixel style)
		hpRatio := float64(s.CurrentHP) / float64(s.MaxHP)
		barW := float32(radius * 2.3)
		barH := float32(4)
		barX := float32(sx) - barW/2
		barY := float32(sy) - float32(radius) - 10

		hpColor := color.RGBA{0x8E, 0xE5, 0xA8, 0xFF}
		if hpRatio < 0.3 {
			hpColor = color.RGBA{0xFF, 0x95, 0xB3, 0xFF}
		} else if hpRatio < 0.6 {
			hpColor = color.RGBA{0xFF, 0xE0, 0x8D, 0xFF}
		}
		drawPixelBar(screen, barX, barY, barW, barH, hpRatio, hpColor, color.RGBA{0xEF, 0xE0, 0xF7, 0xFF})

		// Range indicator (subtle dashed circle)
		if g.selectedCard < 0 {
			rangePixels := float64(s.Range) * float64(config.TileSize)
			if rangePixels > float64(config.TileSize) {
				numDots := 24
				for i := 0; i < numDots; i++ {
					if i%3 != 0 {
						continue
					}
					angle := float64(i) / float64(numDots) * math.Pi * 2
					dotX := sx + math.Cos(angle)*rangePixels
					dotY := sy + math.Sin(angle)*rangePixels
					vector.FillRect(screen, float32(dotX), float32(dotY), 1, 1, color.RGBA{0xD9, 0xBF, 0xF0, 0x90}, false)
				}
			}
		}
	}
}

// ---- Enemies ----

func (g *Game) drawEnemies(screen *ebiten.Image) {
	for _, e := range g.enemies {
		if e.Dead || e.Reached {
			continue
		}

		style := enemyAvatarStyle(e.Type)
		radius := 13.0
		switch e.Type {
		case entity.EnemyGoblin:
			radius = 12.0
		case entity.EnemyOrc:
			radius = 14.0
		case entity.EnemyBossOrc:
			radius = 18.0
		case entity.EnemyFinalBoss:
			radius = 22.0
		}

		// Boss glow effect (animated)
		if e.Type == entity.EnemyBossOrc || e.Type == entity.EnemyFinalBoss {
			glowAlpha := byte(60 + int(40*math.Sin(float64(g.animTick)*0.08)))
			glowSize := float32(radius*1.8) + float32(4*math.Sin(float64(g.animTick)*0.05))
			vector.FillCircle(screen, float32(e.X), float32(e.Y), glowSize, color.RGBA{0xFF, 0xC6, 0xD8, glowAlpha}, false)
		}

		// Final boss aura ring
		if e.HasAura {
			auraSize := float32(50 + 5*math.Sin(float64(g.animTick)*0.06))
			numDots := 16
			for i := 0; i < numDots; i++ {
				angle := float64(i)/float64(numDots)*math.Pi*2 + float64(g.animTick)*0.03
				dotX := float32(e.X) + float32(math.Cos(angle))*auraSize
				dotY := float32(e.Y) + float32(math.Sin(angle))*auraSize
				vector.FillRect(screen, dotX-1, dotY-1, 3, 3, color.RGBA{0xC8, 0xB5, 0xF2, 0x99}, false)
			}
		}

		drawAvatarBadge(screen, e.X, e.Y, radius, style)

		// Health bar
		hpRatio := float64(e.HP) / float64(e.MaxHP)
		barW := float32(radius * 2.2)
		barH := float32(3)
		barX := float32(e.X) - barW/2
		barY := float32(e.Y) - float32(radius) - 8

		drawPixelBar(screen, barX, barY, barW, barH, hpRatio, color.RGBA{0xFF, 0xA5, 0xC4, 0xFF}, color.RGBA{0xEF, 0xE0, 0xF7, 0xFF})
	}
}

// ---- Projectiles ----

func (g *Game) drawProjectiles(screen *ebiten.Image) {
	for _, p := range g.projectiles {
		if p.IsFireball {
			pulse := float32(4.0 + 0.8*math.Sin(float64(g.animTick)*0.3))
			vector.FillCircle(screen, float32(p.X), float32(p.Y), pulse*1.25, color.RGBA{0xFF, 0xB3, 0xC9, 0x68}, false)
			vector.FillCircle(screen, float32(p.X), float32(p.Y), pulse, color.RGBA{0xFF, 0xC4, 0x93, 0xD8}, false)
			vector.FillCircle(screen, float32(p.X), float32(p.Y), pulse*0.45, color.RGBA{0xFF, 0xF0, 0xBE, 0xF0}, false)
		} else {
			vector.FillCircle(screen, float32(p.X), float32(p.Y), 2.4, color.RGBA{0xA2, 0xE2, 0xF2, 0xE0}, false)
		}
	}
}

// ---- Summoner Base ----

func (g *Game) drawSummonerBase(screen *ebiten.Image) {
	lastP := entity.EnemyPath[len(entity.EnemyPath)-1]
	bx, by := lastP.X, lastP.Y

	// Animated glow
	glowAlpha := byte(40 + int(30*math.Sin(float64(g.animTick)*0.04)))
	glowSize := float32(28 + 4*math.Sin(float64(g.animTick)*0.03))
	vector.FillCircle(screen, float32(bx), float32(by), glowSize, color.RGBA{0x9A, 0xE6, 0xD4, glowAlpha}, false)

	pulse := 18.0 + 1.5*math.Sin(float64(g.animTick)*0.05)
	drawAvatarBadge(screen, bx, by, pulse, avatarStyle{
		Outer: color.RGBA{0xD2, 0xC2, 0xF7, 0xFF},
		Inner: color.RGBA{0xF2, 0xE9, 0xFF, 0xFF},
		Ring:  color.RGBA{0xA6, 0xE9, 0xD5, 0xFF},
		Glow:  color.RGBA{0xB8, 0xF2, 0xE5, 0x88},
		Text:  color.RGBA{0x5A, 0x4D, 0x7A, 0xFF},
		Label: "C",
	})

	// Label
	drawKoreanTextWithShadow(screen, "코어", fontSmall, bx-10, by+18, color.RGBA{0x5A, 0x4D, 0x7A, 0xFF})
}

// ---- Fireball mode indicator ----

func (g *Game) drawFireballIndicator(screen *ebiten.Image) {
	flash := byte(140 + int(95*math.Sin(float64(g.animTick)*0.15)))
	drawKoreanTextWithShadow(screen, "스킬 준비! 타겟 위치를 눌러주세요", fontMedium,
		float64(config.ScreenWidth)/2-140, float64(config.ScreenHeight-config.HandHeight-30),
		color.RGBA{0xFF, flash, 0xB2, 0xFF})
	// Animated border flash
	vector.StrokeRect(screen, 0, 0, config.ScreenWidth, config.ScreenHeight, 2, color.RGBA{0xFF, flash, 0xB2, flash}, false)
}
