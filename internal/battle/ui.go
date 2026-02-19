package battle

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"neonsigil/internal/config"
	"neonsigil/internal/data"
	"neonsigil/internal/ui"
)

// DrawHUD draws the top HUD bar
func DrawHUD(screen *ebiten.Image, battle *BattleState, tick int) {
	// Top bar background
	vector.DrawFilledRect(screen, 0, 0, config.ScreenWidth, 42, color.RGBA{8, 8, 20, 240}, false)
	vector.DrawFilledRect(screen, 0, 42, config.ScreenWidth, 1, color.RGBA{0, 255, 255, 60}, false)

	// Stage name
	ui.DrawTextGlow(screen, battle.Stage.ID+" - "+battle.Stage.Name, ui.FontBold(13), 15, 14, config.ColorNeonCyan)

	// Wave info
	waveStr := fmt.Sprintf("WAVE %d/%d", battle.WaveMgr.CurrentWave+1, battle.WaveMgr.TotalWaves())
	ui.DrawText(screen, waveStr, ui.FontBold(12), 360, 14, config.ColorWhite)

	// Gold
	goldStr := fmt.Sprintf("%d", battle.Shop.Gold)
	ui.DrawText(screen, "GOLD", ui.FontRegular(9), 530, 10, config.ColorGold)
	ui.DrawTextGlow(screen, goldStr, ui.FontBold(16), 580, 10, config.ColorGold)

	// Level & deploy cap
	deployed := battle.DeployedCount()
	lvStr := fmt.Sprintf("LV.%d [%d/%d]", battle.Shop.Level, deployed, battle.Shop.DeployCap)
	ui.DrawText(screen, lvStr, ui.FontBold(11), 680, 14, config.ColorNeonGreen)

	// Integrity
	intStr := fmt.Sprintf("INTEGRITY %d", battle.Integrity)
	intColor := config.ColorNeonCyan
	if battle.Integrity <= 5 {
		intColor = config.ColorNeonRed
		// Pulse when low
		if tick%30 < 15 {
			intColor = config.ColorNeonYellow
		}
	}
	ui.DrawText(screen, intStr, ui.FontBold(12), 870, 14, intColor)

	// Integrity bar
	barX := float32(1010)
	barY := float32(14)
	barW := float32(120)
	barH := float32(14)
	vector.DrawFilledRect(screen, barX, barY, barW, barH, color.RGBA{30, 30, 50, 200}, false)
	ratio := float32(battle.Integrity) / float32(battle.MaxIntegrity)
	if ratio < 0 {
		ratio = 0
	}
	barColor := config.ColorNeonCyan
	if ratio < 0.3 {
		barColor = config.ColorNeonRed
	}
	vector.DrawFilledRect(screen, barX, barY, barW*ratio, barH, barColor, false)
	vector.StrokeRect(screen, barX, barY, barW, barH, 1, color.RGBA{0, 200, 255, 100}, false)
}

// DrawShopUI draws the shop panel at the bottom
func DrawShopUI(screen *ebiten.Image, battle *BattleState, tick int) {
	shopY := float32(config.ScreenHeight - 110)

	// Shop background
	vector.DrawFilledRect(screen, 0, shopY, config.ScreenWidth, 110, color.RGBA{10, 10, 25, 240}, false)
	vector.DrawFilledRect(screen, 0, shopY, config.ScreenWidth, 1, color.RGBA{255, 0, 255, 40}, false)

	// "SHOP" label
	ui.DrawText(screen, "SHOP", ui.FontBold(11), 15, float64(shopY)+8, config.ColorNeonMagenta)

	// Shop slots
	for i := 0; i < config.ShopSlots; i++ {
		slotX := float32(80 + i*150)
		slotY := shopY + 6
		slotW := float32(140)
		slotH := float32(56)

		def := battle.Shop.Slots[i]
		if def == nil {
			// Empty slot
			vector.DrawFilledRect(screen, slotX, slotY, slotW, slotH, color.RGBA{20, 20, 35, 200}, false)
			vector.StrokeRect(screen, slotX, slotY, slotW, slotH, 1, color.RGBA{40, 40, 60, 150}, false)
			continue
		}

		// Background with faction color tint
		fc := config.FactionColors[def.Faction]
		bg := color.RGBA{fc.R / 6, fc.G / 6, fc.B / 6, 220}
		vector.DrawFilledRect(screen, slotX, slotY, slotW, slotH, bg, false)

		canBuy := battle.Shop.Gold >= def.Cost
		borderColor := fc
		if !canBuy {
			borderColor = color.RGBA{60, 60, 80, 150}
		}
		vector.StrokeRect(screen, slotX, slotY, slotW, slotH, 1.5, borderColor, false)

		// Unit name
		nameColor := config.ColorWhite
		if !canBuy {
			nameColor = config.ColorWhiteDim
		}
		ui.DrawText(screen, def.Name, ui.FontBold(11), float64(slotX)+6, float64(slotY)+6, nameColor)

		// Faction/class
		tagStr := fmt.Sprintf("%s/%s", def.Faction, def.Class)
		ui.DrawText(screen, tagStr, ui.FontRegular(7), float64(slotX)+6, float64(slotY)+22, fc)

		// Cost
		costStr := fmt.Sprintf("$%d", def.Cost)
		ui.DrawText(screen, costStr, ui.FontBold(12), float64(slotX)+float64(slotW)-35, float64(slotY)+6, config.ColorGold)

		// Stats mini
		statStr := fmt.Sprintf("ATK:%d RNG:%d", int(def.ATK), def.Range)
		ui.DrawText(screen, statStr, ui.FontRegular(7), float64(slotX)+6, float64(slotY)+38, config.ColorWhiteDim)
	}

	// Buttons area
	btnY := float64(shopY) + 68
	btnH := 32.0

	// Reroll button
	battle.BtnReroll.X = 80
	battle.BtnReroll.Y = btnY
	battle.BtnReroll.W = 120
	battle.BtnReroll.H = btnH
	battle.BtnReroll.Label = fmt.Sprintf("REROLL $%d", config.RerollCost)
	battle.BtnReroll.Color = config.ColorNeonMagenta
	battle.BtnReroll.Disabled = !battle.Shop.CanReroll()
	battle.BtnReroll.Draw(screen, tick)

	// Level Up button
	battle.BtnLevelUp.X = 210
	battle.BtnLevelUp.Y = btnY
	battle.BtnLevelUp.W = 140
	battle.BtnLevelUp.H = btnH
	battle.BtnLevelUp.Label = fmt.Sprintf("LEVEL UP $%d", battle.Shop.LevelUpCost())
	battle.BtnLevelUp.Color = config.ColorNeonGreen
	battle.BtnLevelUp.Disabled = !battle.Shop.CanLevelUp()
	battle.BtnLevelUp.Draw(screen, tick)

	// Start Wave button
	if !battle.WaveMgr.WaveActive {
		battle.BtnStartWave.X = 1080
		battle.BtnStartWave.Y = btnY
		battle.BtnStartWave.W = 160
		battle.BtnStartWave.H = btnH
		battle.BtnStartWave.Label = "START WAVE"
		battle.BtnStartWave.Color = config.ColorNeonCyan
		battle.BtnStartWave.Disabled = false
		battle.BtnStartWave.Draw(screen, tick)
	}

	// Sell button (when unit selected)
	if battle.SelectedUnit != nil {
		battle.BtnSell.X = 400
		battle.BtnSell.Y = btnY
		battle.BtnSell.W = 120
		battle.BtnSell.H = btnH
		refund := max(1, battle.SelectedUnit.Def.Cost/2)
		battle.BtnSell.Label = fmt.Sprintf("SELL $%d", refund)
		battle.BtnSell.Color = config.ColorNeonRed
		battle.BtnSell.Disabled = false
		battle.BtnSell.Draw(screen, tick)
	}
}

// DrawBenchUI draws the bench area
func DrawBenchUI(screen *ebiten.Image, battle *BattleState, tick int) {
	benchY := float32(config.BenchSlotY())

	// Bench label
	ui.DrawText(screen, "BENCH", ui.FontRegular(9), float64(config.BoardOffsetX), float64(benchY)-12, config.ColorWhiteDim)

	// Bench slots
	for i := 0; i < config.BenchSlots; i++ {
		bx := float32(config.BenchSlotX(i))
		by := benchY
		s := float32(50)

		// Empty slot
		vector.DrawFilledRect(screen, bx, by, s, s, color.RGBA{18, 18, 35, 200}, false)
		vector.StrokeRect(screen, bx, by, s, s, 1, color.RGBA{40, 40, 70, 150}, false)
	}

	// Draw units on bench
	for _, u := range battle.Units {
		if u.BenchSlot >= 0 && u.BenchSlot < config.BenchSlots {
			u.DrawOnBench(screen, u.BenchSlot, tick)
		}
	}
}

// DrawInfoPanel draws the right-side info panel
func DrawInfoPanel(screen *ebiten.Image, battle *BattleState, tick int) {
	panelX := float64(config.BoardOffsetX + config.BoardCols*config.TileSize + 20)
	panelY := float64(config.BoardOffsetY)
	panelW := float64(config.ScreenWidth) - panelX - 20

	// Panel background
	vector.DrawFilledRect(screen, float32(panelX), float32(panelY), float32(panelW), float32(config.ScreenHeight-230),
		color.RGBA{12, 12, 28, 220}, false)
	vector.StrokeRect(screen, float32(panelX), float32(panelY), float32(panelW), float32(config.ScreenHeight-230),
		1, color.RGBA{60, 60, 100, 100}, false)

	y := panelY + 15

	// Synergy display
	ui.DrawText(screen, "SYNERGIES", ui.FontBold(11), panelX+12, y, config.ColorNeonMagenta)
	y += 28

	factionCounts, classCounts := battle.CountSynergies()

	for _, f := range []config.Faction{config.FactionStreet, config.FactionCoven, config.FactionArcTech, config.FactionExorcist} {
		count := factionCounts[f]
		fc := config.FactionColors[f]
		nameClr := fc
		if count == 0 {
			nameClr = color.RGBA{60, 60, 80, 200}
		}
		ui.DrawText(screen, string(f), ui.FontRegular(9), panelX+12, y, nameClr)

		// Dots
		for d := 0; d < 4; d++ {
			dotX := float32(panelX) + 140 + float32(d)*14
			dotY := float32(y) + 5
			if d < count {
				vector.DrawFilledCircle(screen, dotX, dotY, 4, fc, false)
			} else {
				vector.StrokeCircle(screen, dotX, dotY, 4, 1, color.RGBA{60, 60, 80, 150}, false)
			}
		}

		// Active indicator
		if count >= 2 {
			ui.DrawText(screen, "ACTIVE", ui.FontRegular(7), panelX+200, y+1, config.ColorNeonGreen)
		}
		y += 20
	}

	y += 10
	for _, c := range []config.UnitClass{config.ClassVanguard, config.ClassMarksman, config.ClassCaster, config.ClassEngineer, config.ClassSupport} {
		count := classCounts[c]
		cc := config.ClassColors[c]
		nameClr := cc
		if count == 0 {
			nameClr = color.RGBA{60, 60, 80, 200}
		}
		ui.DrawText(screen, string(c), ui.FontRegular(8), panelX+12, y, nameClr)

		for d := 0; d < 4; d++ {
			dotX := float32(panelX) + 140 + float32(d)*14
			dotY := float32(y) + 5
			if d < count {
				vector.DrawFilledCircle(screen, dotX, dotY, 4, cc, false)
			} else {
				vector.StrokeCircle(screen, dotX, dotY, 4, 1, color.RGBA{60, 60, 80, 150}, false)
			}
		}

		if count >= 2 {
			ui.DrawText(screen, "ACTIVE", ui.FontRegular(7), panelX+200, y+1, config.ColorNeonGreen)
		}
		y += 18
	}

	// Selected unit info
	if battle.SelectedUnit != nil {
		y += 20
		vector.DrawFilledRect(screen, float32(panelX+8), float32(y), float32(panelW-16), 1,
			color.RGBA{60, 60, 100, 100}, false)
		y += 12

		u := battle.SelectedUnit
		fc := config.FactionColors[u.Def.Faction]
		ui.DrawText(screen, u.Def.Name, ui.FontBold(14), panelX+12, y, fc)
		y += 22

		starStr := ""
		for i := 0; i < u.Star; i++ {
			starStr += "* "
		}
		ui.DrawText(screen, starStr, ui.FontBold(12), panelX+12, y, config.ColorNeonYellow)
		y += 18

		tagStr := fmt.Sprintf("%s / %s", u.Def.Faction, u.Def.Class)
		ui.DrawText(screen, tagStr, ui.FontRegular(9), panelX+12, y, config.ColorWhiteDim)
		y += 22

		ui.DrawText(screen, fmt.Sprintf("ATK  %d", int(u.ATK)), ui.FontRegular(10), panelX+12, y, config.ColorNeonRed)
		ui.DrawText(screen, fmt.Sprintf("SPD  %.1f", u.AtkSpeed), ui.FontRegular(10), panelX+120, y, config.ColorNeonCyan)
		y += 18
		ui.DrawText(screen, fmt.Sprintf("RNG  %d", u.Range), ui.FontRegular(10), panelX+12, y, config.ColorNeonGreen)
		ui.DrawText(screen, fmt.Sprintf("ARM  %d", int(u.Def.Armor)), ui.FontRegular(10), panelX+120, y, config.ColorNeonYellow)
		y += 22

		ui.DrawText(screen, u.Def.SkillDesc, ui.FontRegular(8), panelX+12, y, config.ColorNeonMagenta)
	}

	// Wave preview
	y = float64(config.ScreenHeight-230) - 90 + panelY
	if !battle.WaveMgr.AllDone && battle.WaveMgr.CurrentWave < len(battle.Stage.Waves) {
		vector.DrawFilledRect(screen, float32(panelX+8), float32(y), float32(panelW-16), 1,
			color.RGBA{60, 60, 100, 100}, false)
		y += 12
		ui.DrawText(screen, "NEXT WAVE", ui.FontBold(10), panelX+12, y, config.ColorNeonCyan)
		y += 20

		wv := battle.Stage.Waves[battle.WaveMgr.CurrentWave]
		for _, g := range wv.Groups {
			def := data.EnemyDefs[g.Enemy]
			if def != nil {
				eStr := fmt.Sprintf("%s x%d", def.Name, g.Count)
				ui.DrawText(screen, eStr, ui.FontRegular(9), panelX+16, y, def.Color)
				y += 16
			}
		}
	}
}

// DrawNodeIndicator draws node connection lines when 3 nodes occupied
func DrawNodeIndicator(screen *ebiten.Image, battle *BattleState, tick int) {
	if !battle.Stage.NodesEnabled {
		return
	}

	occupiedNodes := battle.GetOccupiedNodes()
	if len(occupiedNodes) < 3 {
		return
	}

	// Draw triangle between first 3 occupied nodes
	pulse := math.Sin(float64(tick%60)/60.0*math.Pi*2)*0.4 + 0.6
	alpha := uint8(float64(80) * pulse)

	for i := 0; i < len(occupiedNodes) && i < 3; i++ {
		next := (i + 1) % min(len(occupiedNodes), 3)
		x1 := float32(config.BoardOffsetX+occupiedNodes[i].X*config.TileSize) + float32(config.TileSize)/2
		y1 := float32(config.BoardOffsetY+occupiedNodes[i].Y*config.TileSize) + float32(config.TileSize)/2
		x2 := float32(config.BoardOffsetX+occupiedNodes[next].X*config.TileSize) + float32(config.TileSize)/2
		y2 := float32(config.BoardOffsetY+occupiedNodes[next].Y*config.TileSize) + float32(config.TileSize)/2
		vector.StrokeLine(screen, x1, y1, x2, y2, 3, color.RGBA{0, 180, 255, alpha}, false)
	}
}
