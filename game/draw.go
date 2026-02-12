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
	shadowOp.ColorScale.ScaleWithColor(color.RGBA{0x00, 0x00, 0x00, 0xAA})
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

	// Draw HP/Mana bars (pixel art style, over the ebitenui panels)
	g.drawHPManaBar(screen)
}

// ---- Path ----

func (g *Game) drawPath(screen *ebiten.Image) {
	pathColor := color.RGBA{0x25, 0x25, 0x35, 0xFF}
	pathBorder := color.RGBA{0x30, 0x30, 0x45, 0x88}
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
		vector.StrokeRect(screen, rx, ry, rw, rh, 1, pathBorder, false)
	}

	// Animated energy dots along path
	drawPathEnergyDots(screen, g.animTick)
}

// ---- Grid ----

func (g *Game) drawGrid(screen *ebiten.Image) {
	mx, my := ebiten.CursorPosition()

	for row := 0; row < config.GridRows; row++ {
		for col := 0; col < config.GridCols; col++ {
			x := float32(config.GridStartX + col*config.TileSize)
			y := float32(config.GridStartY + row*config.TileSize)

			var tile *ebiten.Image
			if g.grid[row][col] != nil {
				tile = g.sprites.OccupiedTile
			} else {
				tile = g.sprites.GrassTile
			}

			// Hover effect
			if g.selectedCard >= 0 && mx >= int(x) && mx < int(x)+config.TileSize && my >= int(y) && my < int(y)+config.TileSize {
				if g.grid[row][col] == nil {
					tile = g.sprites.HoverTile
				} else {
					tile = g.sprites.BlockedTile
				}
			}

			// Draw tile texture
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(tile, op)
		}
	}
}

// ---- Summoners (placed units) ----

func (g *Game) drawSummoners(screen *ebiten.Image) {
	for _, s := range g.summoners {
		sx := s.ScreenX
		sy := s.ScreenY

		// Get sprite for unit type
		var sprite *ebiten.Image
		switch s.Card.Type {
		case entity.CardSoldier:
			sprite = g.sprites.Soldier
		case entity.CardArcher:
			sprite = g.sprites.Archer
		case entity.CardSpearman:
			sprite = g.sprites.Spearman
		case entity.CardMage:
			sprite = g.sprites.Mage
		}

		// Attack flash effect
		scale := 4.0
		if s.AtkTimer == 0 {
			scale = 4.5
		}

		drawSpriteAt(screen, sprite, sx, sy, scale)

		// Health bar (pixel style)
		hpRatio := float64(s.CurrentHP) / float64(s.MaxHP)
		barW := float32(32)
		barH := float32(4)
		barX := float32(sx) - barW/2
		barY := float32(sy) - 24

		hpColor := pixelPalette['b'] // green
		if hpRatio < 0.3 {
			hpColor = pixelPalette['8'] // red
		} else if hpRatio < 0.6 {
			hpColor = pixelPalette['a'] // yellow
		}
		drawPixelBar(screen, barX, barY, barW, barH, hpRatio, hpColor, color.RGBA{0x20, 0x20, 0x20, 0xFF})

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
					vector.FillRect(screen, float32(dotX), float32(dotY), 1, 1, color.RGBA{0x50, 0x50, 0x70, 0x60}, false)
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

		var sprite *ebiten.Image
		var scale float64
		switch e.Type {
		case entity.EnemyGoblin:
			sprite = g.sprites.Goblin
			scale = 3.0
		case entity.EnemyOrc:
			sprite = g.sprites.Orc
			scale = 3.0
		case entity.EnemyBossOrc:
			sprite = g.sprites.BossOrc
			scale = 3.0
		case entity.EnemyFinalBoss:
			sprite = g.sprites.FinalBoss
			scale = 3.0
		}

		// Boss glow effect (animated)
		if e.Type == entity.EnemyBossOrc || e.Type == entity.EnemyFinalBoss {
			glowAlpha := byte(60 + int(40*math.Sin(float64(g.animTick)*0.08)))
			glowSize := float32(scale*6) + float32(4*math.Sin(float64(g.animTick)*0.05))
			vector.FillCircle(screen, float32(e.X), float32(e.Y), glowSize, color.RGBA{0xFF, 0xA0, 0x00, glowAlpha}, false)
		}

		// Final boss aura ring
		if e.HasAura {
			auraSize := float32(50 + 5*math.Sin(float64(g.animTick)*0.06))
			numDots := 16
			for i := 0; i < numDots; i++ {
				angle := float64(i)/float64(numDots)*math.Pi*2 + float64(g.animTick)*0.03
				dotX := float32(e.X) + float32(math.Cos(angle))*auraSize
				dotY := float32(e.Y) + float32(math.Sin(angle))*auraSize
				vector.FillRect(screen, dotX-1, dotY-1, 3, 3, color.RGBA{0xFF, 0x80, 0x00, 0x99}, false)
			}
		}

		drawSpriteAt(screen, sprite, e.X, e.Y, scale)

		// Health bar
		hpRatio := float64(e.HP) / float64(e.MaxHP)
		spriteH := float64(sprite.Bounds().Dy()) * scale
		barW := float32(scale * float64(sprite.Bounds().Dx()))
		barH := float32(3)
		barX := float32(e.X) - barW/2
		barY := float32(e.Y) - float32(spriteH)/2 - 6

		drawPixelBar(screen, barX, barY, barW, barH, hpRatio, pixelPalette['8'], color.RGBA{0x20, 0x20, 0x20, 0xFF})
	}
}

// ---- Projectiles ----

func (g *Game) drawProjectiles(screen *ebiten.Image) {
	for _, p := range g.projectiles {
		if p.IsFireball {
			scale := 3.0 + 0.5*math.Sin(float64(g.animTick)*0.3)
			drawSpriteAt(screen, g.sprites.Fireball, p.X, p.Y, scale)
		} else {
			drawSpriteAt(screen, g.sprites.Arrow, p.X, p.Y, 2.0)
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
	vector.FillCircle(screen, float32(bx), float32(by), glowSize, color.RGBA{0x29, 0xAD, 0xFF, glowAlpha}, false)

	// Crystal sprite
	pulse := 3.0 + 0.2*math.Sin(float64(g.animTick)*0.05)
	drawSpriteAt(screen, g.sprites.Base, bx, by, pulse)

	// Label
	drawKoreanTextWithShadow(screen, "기지", fontSmall, bx-10, by+18, color.RGBA{0xFF, 0xF1, 0xE8, 0xFF})
}

// ---- HP/Mana bars (pixel art style) ----

func (g *Game) drawHPManaBar(screen *ebiten.Image) {
	// HP bar below left sidebar label area
	hpRatio := float64(g.summonerHP) / float64(config.SummonerMaxHP)
	hpColor := pixelPalette['b'] // green
	if hpRatio < 0.3 {
		hpColor = pixelPalette['8'] // red
	} else if hpRatio < 0.6 {
		hpColor = pixelPalette['a'] // yellow
	}
	drawPixelBar(screen, 10, 70, 96, 10, hpRatio, hpColor, color.RGBA{0x20, 0x20, 0x20, 0xFF})

	// Mana bar
	manaRatio := g.mana / float64(g.maxMana)
	manaColor := pixelPalette['c'] // blue
	drawPixelBar(screen, 10, 115, 96, 10, manaRatio, manaColor, color.RGBA{0x20, 0x20, 0x20, 0xFF})

	// Mana regen indicator
	regenProgress := float64(g.manaTimer) / float64(config.ManaRegenTicks)
	if g.mana < float64(g.maxMana) {
		regenBarW := float32(96 * regenProgress)
		vector.FillRect(screen, 10, 126, regenBarW, 2, color.RGBA{0x29, 0xAD, 0xFF, 0x60}, false)
	}
}

// ---- Fireball mode indicator ----

func (g *Game) drawFireballIndicator(screen *ebiten.Image) {
	flash := byte(150 + int(105*math.Sin(float64(g.animTick)*0.15)))
	drawKoreanTextWithShadow(screen, ">> 화염구: 타겟 지역을 클릭하세요! <<", fontMedium,
		float64(config.ScreenWidth)/2-140, float64(config.ScreenHeight-config.HandHeight-30),
		color.RGBA{0xFF, flash, 0x00, 0xFF})
	// Animated border flash
	vector.StrokeRect(screen, 0, 0, config.ScreenWidth, config.ScreenHeight, 2, color.RGBA{0xFF, flash, 0x00, flash}, false)
}

// ---- Unused drawTextWithShadow kept for compatibility ----

func drawTextWithShadow(screen *ebiten.Image, s string, x, y int) {
	drawKoreanTextWithShadow(screen, s, fontSmall, float64(x), float64(y), color.RGBA{0xFF, 0xF1, 0xE8, 0xFF})
}
