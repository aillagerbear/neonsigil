package game

import (
	"fmt"
	"image/color"

	"ebitengine-testing/config"
	"ebitengine-testing/entity"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawTitle(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x0a, 0x0a, 0x1e, 0xff})

	ebitenutil.DebugPrintAt(screen, "========================================", 312, 200)
	ebitenutil.DebugPrintAt(screen, "    SUMMONER'S DEFENSE", 370, 230)
	ebitenutil.DebugPrintAt(screen, "========================================", 312, 260)
	ebitenutil.DebugPrintAt(screen, "    Auto-Battle Deck Building Defense", 330, 310)
	ebitenutil.DebugPrintAt(screen, "    Click anywhere to start!", 370, 400)

	ebitenutil.DebugPrintAt(screen, "Controls:", 400, 480)
	ebitenutil.DebugPrintAt(screen, "- Click card in hand to select", 370, 510)
	ebitenutil.DebugPrintAt(screen, "- Click grid tile to place unit", 370, 530)
	ebitenutil.DebugPrintAt(screen, "- Units auto-attack nearby enemies", 370, 550)
	ebitenutil.DebugPrintAt(screen, "- Survive 10 waves to win!", 370, 570)
}

func (g *Game) drawBattle(screen *ebiten.Image) {
	g.drawPath(screen)
	g.drawGrid(screen)
	g.drawSummoners(screen)
	g.drawEnemies(screen)
	g.drawProjectiles(screen)
	g.drawSummonerBase(screen)
	g.drawHUD(screen)
	g.drawHand(screen)
	g.drawSynergies(screen)
}

func (g *Game) drawPath(screen *ebiten.Image) {
	pathColor := color.RGBA{0x3a, 0x3a, 0x50, 0xff}
	pathWidth := float32(30)

	for i := 0; i < len(entity.EnemyPath)-1; i++ {
		x1, y1 := float32(entity.EnemyPath[i].X), float32(entity.EnemyPath[i].Y)
		x2, y2 := float32(entity.EnemyPath[i+1].X), float32(entity.EnemyPath[i+1].Y)
		vector.StrokeLine(screen, x1, y1, x2, y2, pathWidth, pathColor, true)
	}
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	mx, my := ebiten.CursorPosition()

	for row := 0; row < config.GridRows; row++ {
		for col := 0; col < config.GridCols; col++ {
			x := float32(config.GridStartX + col*config.TileSize)
			y := float32(config.GridStartY + row*config.TileSize)

			tileColor := color.RGBA{0x2a, 0x4a, 0x2a, 0xaa}
			if g.grid[row][col] != nil {
				tileColor = color.RGBA{0x4a, 0x4a, 0x2a, 0xaa}
			}

			// 호버 강조
			if g.selectedCard >= 0 && mx >= int(x) && mx < int(x)+config.TileSize && my >= int(y) && my < int(y)+config.TileSize {
				if g.grid[row][col] == nil {
					tileColor = color.RGBA{0x3a, 0x7a, 0x3a, 0xcc}
				} else {
					tileColor = color.RGBA{0x7a, 0x3a, 0x3a, 0xcc}
				}
			}

			vector.FillRect(screen, x+1, y+1, float32(config.TileSize-2), float32(config.TileSize-2), tileColor, true)
			vector.StrokeRect(screen, x, y, float32(config.TileSize), float32(config.TileSize), 1, color.RGBA{0x55, 0x88, 0x55, 0xff}, true)
		}
	}
}

func (g *Game) drawSummoners(screen *ebiten.Image) {
	for _, s := range g.summoners {
		sx := float32(s.ScreenX)
		sy := float32(s.ScreenY)

		var unitColor color.RGBA
		switch s.Card.Type {
		case entity.CardSoldier:
			unitColor = color.RGBA{0x40, 0x80, 0xff, 0xff} // 파랑 - 보병
		case entity.CardArcher:
			unitColor = color.RGBA{0x80, 0xff, 0x40, 0xff} // 녹색 - 궁수
		case entity.CardSpearman:
			unitColor = color.RGBA{0xff, 0xc0, 0x40, 0xff} // 주황 - 창병
		case entity.CardMage:
			unitColor = color.RGBA{0xc0, 0x40, 0xff, 0xff} // 보라 - 마법사
		}

		size := float32(30)
		vector.FillRect(screen, sx-size/2, sy-size/2, size, size, unitColor, true)

		label := ""
		switch s.Card.Type {
		case entity.CardSoldier:
			label = "S"
		case entity.CardArcher:
			label = "A"
		case entity.CardSpearman:
			label = "P"
		case entity.CardMage:
			label = "M"
		}
		ebitenutil.DebugPrintAt(screen, label, int(sx)-3, int(sy)-6)

		// 체력바
		hpRatio := float32(s.CurrentHP) / float32(s.MaxHP)
		barW := float32(30)
		barH := float32(4)
		vector.FillRect(screen, sx-barW/2, sy-size/2-6, barW, barH, color.RGBA{0x40, 0x40, 0x40, 0xff}, true)
		vector.FillRect(screen, sx-barW/2, sy-size/2-6, barW*hpRatio, barH, color.RGBA{0x40, 0xff, 0x40, 0xff}, true)
	}
}

func (g *Game) drawEnemies(screen *ebiten.Image) {
	for _, e := range g.enemies {
		if e.Dead || e.Reached {
			continue
		}

		var enemyColor color.RGBA
		var size float32
		switch e.Type {
		case entity.EnemyGoblin:
			enemyColor = color.RGBA{0xff, 0x60, 0x60, 0xff}
			size = 12
		case entity.EnemyOrc:
			enemyColor = color.RGBA{0xff, 0x30, 0x30, 0xff}
			size = 16
		case entity.EnemyBossOrc:
			enemyColor = color.RGBA{0xff, 0x00, 0x00, 0xff}
			size = 22
		case entity.EnemyFinalBoss:
			enemyColor = color.RGBA{0xff, 0x00, 0x80, 0xff}
			size = 28
		}

		vector.FillCircle(screen, float32(e.X), float32(e.Y), size, enemyColor, true)

		// 보스 표시
		if e.Type == entity.EnemyBossOrc || e.Type == entity.EnemyFinalBoss {
			vector.StrokeCircle(screen, float32(e.X), float32(e.Y), size+3, 2, color.RGBA{0xff, 0xff, 0x00, 0xff}, true)
		}

		// 오라 표시
		if e.HasAura {
			vector.StrokeCircle(screen, float32(e.X), float32(e.Y), size+8, 1, color.RGBA{0xff, 0x80, 0x00, 0x80}, true)
		}

		// 체력바
		hpRatio := float32(e.HP) / float32(e.MaxHP)
		barW := size * 2
		barH := float32(3)
		vector.FillRect(screen, float32(e.X)-barW/2, float32(e.Y)-size-6, barW, barH, color.RGBA{0x40, 0x40, 0x40, 0xff}, true)
		vector.FillRect(screen, float32(e.X)-barW/2, float32(e.Y)-size-6, barW*hpRatio, barH, color.RGBA{0xff, 0x40, 0x40, 0xff}, true)
	}
}

func (g *Game) drawProjectiles(screen *ebiten.Image) {
	for _, p := range g.projectiles {
		if p.IsFireball {
			vector.FillCircle(screen, float32(p.X), float32(p.Y), 8, color.RGBA{0xff, 0x80, 0x00, 0xff}, true)
		} else {
			vector.FillCircle(screen, float32(p.X), float32(p.Y), 4, color.RGBA{0xff, 0xff, 0x80, 0xff}, true)
		}
	}
}

func (g *Game) drawSummonerBase(screen *ebiten.Image) {
	lastP := entity.EnemyPath[len(entity.EnemyPath)-1]
	baseColor := color.RGBA{0x00, 0xaa, 0xff, 0xff}
	vector.FillCircle(screen, float32(lastP.X), float32(lastP.Y), 20, baseColor, true)
	vector.StrokeCircle(screen, float32(lastP.X), float32(lastP.Y), 22, 2, color.RGBA{0x00, 0xff, 0xff, 0xff}, true)
	ebitenutil.DebugPrintAt(screen, "BASE", int(lastP.X)-14, int(lastP.Y)-6)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	// 상단 바 배경
	vector.FillRect(screen, 0, 0, config.ScreenWidth, float32(config.HUDHeight), color.RGBA{0x10, 0x10, 0x20, 0xee}, true)

	// 웨이브 정보
	waveText := fmt.Sprintf("Wave %d/%d", g.wave+1, g.maxWave)
	ebitenutil.DebugPrintAt(screen, waveText, 20, 12)

	// 남은 적
	aliveCount := 0
	for _, e := range g.enemies {
		if !e.Dead && !e.Reached {
			aliveCount++
		}
	}
	aliveCount += len(g.spawnQueue)
	enemyText := fmt.Sprintf("Enemies: %d", aliveCount)
	ebitenutil.DebugPrintAt(screen, enemyText, 180, 12)

	// FPS
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.0f", ebiten.ActualFPS()), 400, 12)

	// 속도 버튼
	btn1Color := color.RGBA{0x40, 0x40, 0x60, 0xff}
	btn2Color := color.RGBA{0x40, 0x40, 0x60, 0xff}
	if g.gameSpeed == 1 {
		btn1Color = color.RGBA{0x40, 0x80, 0x40, 0xff}
	} else {
		btn2Color = color.RGBA{0x40, 0x80, 0x40, 0xff}
	}
	vector.FillRect(screen, 880, 5, 40, 30, btn1Color, true)
	vector.FillRect(screen, 930, 5, 40, 30, btn2Color, true)
	ebitenutil.DebugPrintAt(screen, "1x", 892, 12)
	ebitenutil.DebugPrintAt(screen, "2x", 942, 12)

	// 좌측 사이드바
	vector.FillRect(screen, 0, float32(config.HUDHeight), float32(config.SidebarWidth), float32(config.ScreenHeight-config.HUDHeight-config.HandHeight), color.RGBA{0x15, 0x15, 0x25, 0xdd}, true)

	// HP 바
	ebitenutil.DebugPrintAt(screen, "HP", 10, 60)
	hpRatio := float32(g.summonerHP) / float32(config.SummonerMaxHP)
	vector.FillRect(screen, 10, 80, 100, 20, color.RGBA{0x40, 0x40, 0x40, 0xff}, true)
	hpBarColor := color.RGBA{0x40, 0xff, 0x40, 0xff}
	if hpRatio < 0.3 {
		hpBarColor = color.RGBA{0xff, 0x40, 0x40, 0xff}
	} else if hpRatio < 0.6 {
		hpBarColor = color.RGBA{0xff, 0xff, 0x40, 0xff}
	}
	vector.FillRect(screen, 10, 80, 100*hpRatio, 20, hpBarColor, true)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d/%d", g.summonerHP, config.SummonerMaxHP), 25, 83)

	// 마나 바
	ebitenutil.DebugPrintAt(screen, "MANA", 10, 120)
	manaRatio := float32(g.mana) / float32(g.maxMana)
	vector.FillRect(screen, 10, 140, 100, 20, color.RGBA{0x40, 0x40, 0x40, 0xff}, true)
	vector.FillRect(screen, 10, 140, 100*manaRatio, 20, color.RGBA{0x40, 0x80, 0xff, 0xff}, true)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f/%d", g.mana, g.maxMana), 30, 143)
}

func (g *Game) drawHand(screen *ebiten.Image) {
	handY := config.ScreenHeight - config.HandHeight

	// 핸드 배경
	vector.FillRect(screen, 0, float32(handY), config.ScreenWidth, float32(config.HandHeight), color.RGBA{0x15, 0x15, 0x25, 0xee}, true)

	// 덱/묘지 카운트
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Deck:%d", len(g.deck)), 15, handY+50)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Grav:%d", len(g.graveyard)), config.ScreenWidth-80, handY+50)

	// 카드 렌더링
	cardWidth := 130
	cardHeight := 90
	cardSpacing := 10
	totalWidth := len(g.hand)*(cardWidth+cardSpacing) - cardSpacing
	startX := (config.ScreenWidth - totalWidth) / 2

	mx, my := ebiten.CursorPosition()

	for i, card := range g.hand {
		cx := startX + i*(cardWidth+cardSpacing)
		cy := handY + 10

		bgColor := color.RGBA{0x30, 0x30, 0x50, 0xff}
		canAfford := g.mana >= float64(card.Data.Cost)

		if !canAfford {
			bgColor = color.RGBA{0x30, 0x20, 0x20, 0xff}
		}

		// 호버 효과
		isHover := mx >= cx && mx <= cx+cardWidth && my >= cy && my <= cy+cardHeight
		if isHover && canAfford {
			bgColor = color.RGBA{0x40, 0x50, 0x70, 0xff}
		}

		// 선택 효과
		if g.selectedCard == i {
			bgColor = color.RGBA{0x50, 0x70, 0x50, 0xff}
		}

		vector.FillRect(screen, float32(cx), float32(cy), float32(cardWidth), float32(cardHeight), bgColor, true)

		// 카드 테두리
		borderColor := color.RGBA{0x80, 0x80, 0xa0, 0xff}
		if g.selectedCard == i {
			borderColor = color.RGBA{0x80, 0xff, 0x80, 0xff}
		}
		vector.StrokeRect(screen, float32(cx), float32(cy), float32(cardWidth), float32(cardHeight), 2, borderColor, true)

		// 종족 색상 표시
		var raceColor color.RGBA
		switch card.Data.Race {
		case entity.RaceHuman:
			raceColor = color.RGBA{0x40, 0x80, 0xff, 0xff}
		case entity.RaceElf:
			raceColor = color.RGBA{0x40, 0xff, 0x80, 0xff}
		default:
			raceColor = color.RGBA{0xff, 0x80, 0x40, 0xff}
		}
		vector.FillRect(screen, float32(cx), float32(cy), 5, float32(cardHeight), raceColor, true)

		// 카드 정보
		ebitenutil.DebugPrintAt(screen, card.Data.Name, cx+10, cy+5)

		costColor := ""
		if !canAfford {
			costColor = "(X)"
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cost:%d%s", card.Data.Cost, costColor), cx+10, cy+22)

		if card.Data.Type != entity.CardFireball {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP:%d ATK:%d", card.Data.HP, card.Data.Atk), cx+10, cy+39)
			rangeText := "Melee"
			if card.Data.Range > 1 {
				rangeText = fmt.Sprintf("Range:%d", card.Data.Range)
			}
			ebitenutil.DebugPrintAt(screen, rangeText, cx+10, cy+56)
		} else {
			ebitenutil.DebugPrintAt(screen, "DMG:20 AOE", cx+10, cy+39)
			ebitenutil.DebugPrintAt(screen, "Click target", cx+10, cy+56)
		}

		// 종족 표시
		raceName := ""
		switch card.Data.Race {
		case entity.RaceHuman:
			raceName = "[Human]"
		case entity.RaceElf:
			raceName = "[Elf]"
		}
		if raceName != "" {
			ebitenutil.DebugPrintAt(screen, raceName, cx+10, cy+73)
		}
	}

	// 파이어볼 모드 안내
	if g.fireballMode {
		ebitenutil.DebugPrintAt(screen, ">> FIREBALL: Click target area! <<", 380, handY-20)
	}
}

func (g *Game) drawSynergies(screen *ebiten.Image) {
	// 우측 사이드바
	sx := float32(config.ScreenWidth - config.SidebarWidth)
	vector.FillRect(screen, sx, float32(config.HUDHeight), float32(config.SidebarWidth), float32(config.ScreenHeight-config.HUDHeight-config.HandHeight), color.RGBA{0x15, 0x15, 0x25, 0xdd}, true)

	ebitenutil.DebugPrintAt(screen, "Synergies", config.ScreenWidth-config.SidebarWidth+10, 55)

	// 인간 시너지
	humanText := fmt.Sprintf("Human %d/2", g.humanCount)
	humanColor := color.RGBA{0x60, 0x60, 0x60, 0xff}
	if g.humanSynergy {
		humanColor = color.RGBA{0x40, 0x80, 0xff, 0xff}
		humanText += " ON"
	}
	vector.FillRect(screen, sx+5, 80, 110, 40, humanColor, true)
	ebitenutil.DebugPrintAt(screen, humanText, config.ScreenWidth-config.SidebarWidth+10, 85)
	if g.humanSynergy {
		ebitenutil.DebugPrintAt(screen, "HP +20%", config.ScreenWidth-config.SidebarWidth+10, 102)
	}

	// 엘프 시너지
	elfText := fmt.Sprintf("Elf %d/2", g.elfCount)
	elfColor := color.RGBA{0x60, 0x60, 0x60, 0xff}
	if g.elfSynergy {
		elfColor = color.RGBA{0x40, 0xff, 0x80, 0xff}
		elfText += " ON"
	}
	vector.FillRect(screen, sx+5, 130, 110, 40, elfColor, true)
	ebitenutil.DebugPrintAt(screen, elfText, config.ScreenWidth-config.SidebarWidth+10, 135)
	if g.elfSynergy {
		ebitenutil.DebugPrintAt(screen, "SPD +20%", config.ScreenWidth-config.SidebarWidth+10, 152)
	}
}

func (g *Game) drawReward(screen *ebiten.Image) {
	// 반투명 오버레이
	vector.FillRect(screen, 0, 0, config.ScreenWidth, config.ScreenHeight, color.RGBA{0x00, 0x00, 0x00, 0xbb}, true)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Wave %d Clear! Choose a card:", g.wave+1), 380, 120)

	cardWidth := 180
	cardHeight := 250
	spacing := 30
	totalWidth := 3*(cardWidth+spacing) - spacing
	startX := (config.ScreenWidth - totalWidth) / 2
	startY := (config.ScreenHeight - cardHeight) / 2

	for i, card := range g.rewardCards {
		cx := startX + i*(cardWidth+spacing)
		cy := startY

		bgColor := color.RGBA{0x30, 0x30, 0x50, 0xff}
		if g.rewardHover == i {
			bgColor = color.RGBA{0x50, 0x60, 0x80, 0xff}
		}

		vector.FillRect(screen, float32(cx), float32(cy), float32(cardWidth), float32(cardHeight), bgColor, true)

		borderColor := color.RGBA{0x80, 0x80, 0xa0, 0xff}
		if g.rewardHover == i {
			borderColor = color.RGBA{0xff, 0xff, 0x40, 0xff}
		}
		vector.StrokeRect(screen, float32(cx), float32(cy), float32(cardWidth), float32(cardHeight), 2, borderColor, true)

		// 종족 색상 바
		var raceColor color.RGBA
		switch card.Data.Race {
		case entity.RaceHuman:
			raceColor = color.RGBA{0x40, 0x80, 0xff, 0xff}
		case entity.RaceElf:
			raceColor = color.RGBA{0x40, 0xff, 0x80, 0xff}
		default:
			raceColor = color.RGBA{0xff, 0x80, 0x40, 0xff}
		}
		vector.FillRect(screen, float32(cx), float32(cy), 8, float32(cardHeight), raceColor, true)

		// 카드 정보
		ebitenutil.DebugPrintAt(screen, card.Data.Name, cx+20, cy+20)

		raceName := ""
		switch card.Data.Race {
		case entity.RaceHuman:
			raceName = "Human"
		case entity.RaceElf:
			raceName = "Elf"
		}
		if raceName != "" {
			ebitenutil.DebugPrintAt(screen, raceName, cx+20, cy+45)
		}

		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cost: %d", card.Data.Cost), cx+20, cy+70)

		if card.Data.Type != entity.CardFireball {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP: %d", card.Data.HP), cx+20, cy+95)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ATK: %d", card.Data.Atk), cx+20, cy+120)

			speedText := "Normal"
			if card.Data.AtkSpeed < 50 {
				speedText = "Fast"
			} else if card.Data.AtkSpeed > 70 {
				speedText = "Slow"
			}
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Speed: %s", speedText), cx+20, cy+145)

			rangeText := "Melee"
			if card.Data.Range > 1 {
				rangeText = fmt.Sprintf("%d tiles", card.Data.Range)
			}
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Range: %s", rangeText), cx+20, cy+170)
		} else {
			ebitenutil.DebugPrintAt(screen, "Spell Card", cx+20, cy+95)
			ebitenutil.DebugPrintAt(screen, "DMG: 20 (AOE)", cx+20, cy+120)
			ebitenutil.DebugPrintAt(screen, "Damages all", cx+20, cy+145)
			ebitenutil.DebugPrintAt(screen, "enemies nearby", cx+20, cy+170)
		}

		ebitenutil.DebugPrintAt(screen, "Click to add", cx+45, cy+220)
	}
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x20, 0x00, 0x00, 0xff})
	ebitenutil.DebugPrintAt(screen, "========================================", 312, 280)
	ebitenutil.DebugPrintAt(screen, "          GAME OVER", 390, 310)
	ebitenutil.DebugPrintAt(screen, "========================================", 312, 340)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("   You reached Wave %d/%d", g.wave+1, g.maxWave), 370, 400)
	ebitenutil.DebugPrintAt(screen, "   Click to return to title", 370, 460)
}

func (g *Game) drawVictory(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x10, 0x20, 0xff})
	ebitenutil.DebugPrintAt(screen, "========================================", 312, 250)
	ebitenutil.DebugPrintAt(screen, "        VICTORY!", 400, 280)
	ebitenutil.DebugPrintAt(screen, "========================================", 312, 310)
	ebitenutil.DebugPrintAt(screen, "   You defended against all 10 waves!", 340, 370)
	ebitenutil.DebugPrintAt(screen, "   The summoner is safe!", 370, 400)
	ebitenutil.DebugPrintAt(screen, "   Click to return to title", 370, 460)
}
