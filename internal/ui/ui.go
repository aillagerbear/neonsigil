package ui

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//go:embed assets/fonts/Orbitron-Bold.ttf
var orbitronBoldTTF []byte

//go:embed assets/fonts/Orbitron-Regular.ttf
var orbitronRegularTTF []byte

var (
	fontSourceBold    *text.GoTextFaceSource
	fontSourceRegular *text.GoTextFaceSource
)

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

func fontBold(size float64) *text.GoTextFace {
	return &text.GoTextFace{Source: fontSourceBold, Size: size}
}

func fontRegular(size float64) *text.GoTextFace {
	return &text.GoTextFace{Source: fontSourceRegular, Size: size}
}

// DrawText draws text with neon glow effect
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

// DrawTextGlow draws text with neon glow
func DrawTextGlow(screen *ebiten.Image, str string, face *text.GoTextFace, x, y float64, clr color.RGBA) {
	// Glow layers
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

// Button represents a clickable UI button
type Button struct {
	X, Y, W, H float64
	Label       string
	Color       color.RGBA
	Disabled    bool
	Hovered     bool
}

func (b *Button) Contains(mx, my int) bool {
	return float64(mx) >= b.X && float64(mx) <= b.X+b.W &&
		float64(my) >= b.Y && float64(my) <= b.Y+b.H
}

func (b *Button) Draw(screen *ebiten.Image, tick int) {
	bgColor := ColorBtnBG
	borderColor := b.Color
	if b.Disabled {
		bgColor = color.RGBA{20, 20, 30, 200}
		borderColor = color.RGBA{60, 60, 80, 150}
	} else if b.Hovered {
		bgColor = ColorBtnHover
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
	DrawTextCentered(screen, b.Label, fontBold(11), b.X+b.W/2, b.Y+b.H/2, textColor)
}

// DrawHUD draws the top HUD bar
func DrawHUD(screen *ebiten.Image, battle *BattleState, tick int) {
	// Top bar background
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, 42, color.RGBA{8, 8, 20, 240}, false)
	vector.DrawFilledRect(screen, 0, 42, ScreenWidth, 1, color.RGBA{0, 255, 255, 60}, false)

	// Stage name
	DrawTextGlow(screen, battle.Stage.ID+" - "+battle.Stage.Name, fontBold(13), 15, 14, ColorNeonCyan)

	// Wave info
	waveStr := fmt.Sprintf("WAVE %d/%d", battle.WaveMgr.CurrentWave+1, battle.WaveMgr.TotalWaves())
	if battle.WaveMgr.WaveActive {
		waveStr = fmt.Sprintf("WAVE %d/%d", battle.WaveMgr.CurrentWave+1, battle.WaveMgr.TotalWaves())
	}
	DrawText(screen, waveStr, fontBold(12), 360, 14, ColorWhite)

	// Gold
	goldStr := fmt.Sprintf("%d", battle.Shop.Gold)
	DrawText(screen, "GOLD", fontRegular(9), 530, 10, ColorGold)
	DrawTextGlow(screen, goldStr, fontBold(16), 580, 10, ColorGold)

	// Level & deploy cap
	deployed := battle.DeployedCount()
	lvStr := fmt.Sprintf("LV.%d [%d/%d]", battle.Shop.Level, deployed, battle.Shop.DeployCap)
	DrawText(screen, lvStr, fontBold(11), 680, 14, ColorNeonGreen)

	// Integrity
	intStr := fmt.Sprintf("INTEGRITY %d", battle.Integrity)
	intColor := ColorNeonCyan
	if battle.Integrity <= 5 {
		intColor = ColorNeonRed
		// Pulse when low
		if tick%30 < 15 {
			intColor = ColorNeonYellow
		}
	}
	DrawText(screen, intStr, fontBold(12), 870, 14, intColor)

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
	barColor := ColorNeonCyan
	if ratio < 0.3 {
		barColor = ColorNeonRed
	}
	vector.DrawFilledRect(screen, barX, barY, barW*ratio, barH, barColor, false)
	vector.StrokeRect(screen, barX, barY, barW, barH, 1, color.RGBA{0, 200, 255, 100}, false)
}

// DrawShopUI draws the shop panel at the bottom
func DrawShopUI(screen *ebiten.Image, battle *BattleState, tick int) {
	shopY := float32(ScreenHeight - 110)

	// Shop background
	vector.DrawFilledRect(screen, 0, shopY, ScreenWidth, 110, color.RGBA{10, 10, 25, 240}, false)
	vector.DrawFilledRect(screen, 0, shopY, ScreenWidth, 1, color.RGBA{255, 0, 255, 40}, false)

	// "SHOP" label
	DrawText(screen, "SHOP", fontBold(11), 15, float64(shopY)+8, ColorNeonMagenta)

	// Shop slots
	for i := 0; i < ShopSlots; i++ {
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
		fc := FactionColors[def.Faction]
		bg := color.RGBA{fc.R / 6, fc.G / 6, fc.B / 6, 220}
		vector.DrawFilledRect(screen, slotX, slotY, slotW, slotH, bg, false)

		canBuy := battle.Shop.Gold >= def.Cost
		borderColor := fc
		if !canBuy {
			borderColor = color.RGBA{60, 60, 80, 150}
		}
		vector.StrokeRect(screen, slotX, slotY, slotW, slotH, 1.5, borderColor, false)

		// Unit name
		nameColor := ColorWhite
		if !canBuy {
			nameColor = ColorWhiteDim
		}
		DrawText(screen, def.Name, fontBold(11), float64(slotX)+6, float64(slotY)+6, nameColor)

		// Faction/class
		tagStr := fmt.Sprintf("%s/%s", def.Faction, def.Class)
		DrawText(screen, tagStr, fontRegular(7), float64(slotX)+6, float64(slotY)+22, fc)

		// Cost
		costStr := fmt.Sprintf("$%d", def.Cost)
		DrawText(screen, costStr, fontBold(12), float64(slotX)+float64(slotW)-35, float64(slotY)+6, ColorGold)

		// Stats mini
		statStr := fmt.Sprintf("ATK:%d RNG:%d", int(def.ATK), def.Range)
		DrawText(screen, statStr, fontRegular(7), float64(slotX)+6, float64(slotY)+38, ColorWhiteDim)
	}

	// Buttons area
	btnY := float64(shopY) + 68
	btnH := 32.0

	// Reroll button
	battle.BtnReroll.X = 80
	battle.BtnReroll.Y = btnY
	battle.BtnReroll.W = 120
	battle.BtnReroll.H = btnH
	battle.BtnReroll.Label = fmt.Sprintf("REROLL $%d", RerollCost)
	battle.BtnReroll.Color = ColorNeonMagenta
	battle.BtnReroll.Disabled = !battle.Shop.CanReroll()
	battle.BtnReroll.Draw(screen, tick)

	// Level Up button
	battle.BtnLevelUp.X = 210
	battle.BtnLevelUp.Y = btnY
	battle.BtnLevelUp.W = 140
	battle.BtnLevelUp.H = btnH
	battle.BtnLevelUp.Label = fmt.Sprintf("LEVEL UP $%d", battle.Shop.LevelUpCost())
	battle.BtnLevelUp.Color = ColorNeonGreen
	battle.BtnLevelUp.Disabled = !battle.Shop.CanLevelUp()
	battle.BtnLevelUp.Draw(screen, tick)

	// Start Wave button
	if !battle.WaveMgr.WaveActive {
		battle.BtnStartWave.X = 1080
		battle.BtnStartWave.Y = btnY
		battle.BtnStartWave.W = 160
		battle.BtnStartWave.H = btnH
		battle.BtnStartWave.Label = "START WAVE"
		battle.BtnStartWave.Color = ColorNeonCyan
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
		battle.BtnSell.Color = ColorNeonRed
		battle.BtnSell.Disabled = false
		battle.BtnSell.Draw(screen, tick)
	}
}

// DrawBenchUI draws the bench area
func DrawBenchUI(screen *ebiten.Image, battle *BattleState, tick int) {
	benchY := float32(benchSlotY())

	// Bench label
	DrawText(screen, "BENCH", fontRegular(9), float64(BoardOffsetX), float64(benchY)-12, ColorWhiteDim)

	// Bench slots
	for i := 0; i < BenchSlots; i++ {
		bx := float32(benchSlotX(i))
		by := benchY
		s := float32(50)

		// Empty slot
		vector.DrawFilledRect(screen, bx, by, s, s, color.RGBA{18, 18, 35, 200}, false)
		vector.StrokeRect(screen, bx, by, s, s, 1, color.RGBA{40, 40, 70, 150}, false)
	}

	// Draw units on bench
	for _, u := range battle.Units {
		if u.BenchSlot >= 0 && u.BenchSlot < BenchSlots {
			u.DrawOnBench(screen, u.BenchSlot, tick)
		}
	}
}

// DrawInfoPanel draws the right-side info panel
func DrawInfoPanel(screen *ebiten.Image, battle *BattleState, tick int) {
	panelX := float64(BoardOffsetX + BoardCols*TileSize + 20)
	panelY := float64(BoardOffsetY)
	panelW := float64(ScreenWidth) - panelX - 20

	// Panel background
	vector.DrawFilledRect(screen, float32(panelX), float32(panelY), float32(panelW), float32(ScreenHeight-230),
		color.RGBA{12, 12, 28, 220}, false)
	vector.StrokeRect(screen, float32(panelX), float32(panelY), float32(panelW), float32(ScreenHeight-230),
		1, color.RGBA{60, 60, 100, 100}, false)

	y := panelY + 15

	// Synergy display
	DrawText(screen, "SYNERGIES", fontBold(11), panelX+12, y, ColorNeonMagenta)
	y += 28

	factionCounts, classCounts := battle.CountSynergies()

	for _, f := range []Faction{FactionStreet, FactionCoven, FactionArcTech, FactionExorcist} {
		count := factionCounts[f]
		fc := FactionColors[f]
		nameClr := fc
		if count == 0 {
			nameClr = color.RGBA{60, 60, 80, 200}
		}
		DrawText(screen, string(f), fontRegular(9), panelX+12, y, nameClr)

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
			DrawText(screen, "ACTIVE", fontRegular(7), panelX+200, y+1, ColorNeonGreen)
		}
		y += 20
	}

	y += 10
	for _, c := range []UnitClass{ClassVanguard, ClassMarksman, ClassCaster, ClassEngineer, ClassSupport} {
		count := classCounts[c]
		cc := ClassColors[c]
		nameClr := cc
		if count == 0 {
			nameClr = color.RGBA{60, 60, 80, 200}
		}
		DrawText(screen, string(c), fontRegular(8), panelX+12, y, nameClr)

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
			DrawText(screen, "ACTIVE", fontRegular(7), panelX+200, y+1, ColorNeonGreen)
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
		fc := FactionColors[u.Def.Faction]
		DrawText(screen, u.Def.Name, fontBold(14), panelX+12, y, fc)
		y += 22

		starStr := ""
		for i := 0; i < u.Star; i++ {
			starStr += "* "
		}
		DrawText(screen, starStr, fontBold(12), panelX+12, y, ColorNeonYellow)
		y += 18

		tagStr := fmt.Sprintf("%s / %s", u.Def.Faction, u.Def.Class)
		DrawText(screen, tagStr, fontRegular(9), panelX+12, y, ColorWhiteDim)
		y += 22

		DrawText(screen, fmt.Sprintf("ATK  %d", int(u.ATK)), fontRegular(10), panelX+12, y, ColorNeonRed)
		DrawText(screen, fmt.Sprintf("SPD  %.1f", u.AtkSpeed), fontRegular(10), panelX+120, y, ColorNeonCyan)
		y += 18
		DrawText(screen, fmt.Sprintf("RNG  %d", u.Range), fontRegular(10), panelX+12, y, ColorNeonGreen)
		DrawText(screen, fmt.Sprintf("ARM  %d", int(u.Def.Armor)), fontRegular(10), panelX+120, y, ColorNeonYellow)
		y += 22

		DrawText(screen, u.Def.SkillDesc, fontRegular(8), panelX+12, y, ColorNeonMagenta)
	}

	// Wave preview
	y = float64(ScreenHeight-230) - 90 + panelY
	if !battle.WaveMgr.AllDone && battle.WaveMgr.CurrentWave < len(battle.Stage.Waves) {
		vector.DrawFilledRect(screen, float32(panelX+8), float32(y), float32(panelW-16), 1,
			color.RGBA{60, 60, 100, 100}, false)
		y += 12
		DrawText(screen, "NEXT WAVE", fontBold(10), panelX+12, y, ColorNeonCyan)
		y += 20

		wave := battle.Stage.Waves[battle.WaveMgr.CurrentWave]
		for _, g := range wave.Groups {
			def := EnemyDefs[g.Enemy]
			if def != nil {
				eStr := fmt.Sprintf("%s x%d", def.Name, g.Count)
				DrawText(screen, eStr, fontRegular(9), panelX+16, y, def.Color)
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
		next := (i + 1) % min3(len(occupiedNodes), 3)
		x1 := float32(BoardOffsetX+occupiedNodes[i].X*TileSize) + float32(TileSize)/2
		y1 := float32(BoardOffsetY+occupiedNodes[i].Y*TileSize) + float32(TileSize)/2
		x2 := float32(BoardOffsetX+occupiedNodes[next].X*TileSize) + float32(TileSize)/2
		y2 := float32(BoardOffsetY+occupiedNodes[next].Y*TileSize) + float32(TileSize)/2
		vector.StrokeLine(screen, x1, y1, x2, y2, 3, color.RGBA{0, 180, 255, alpha}, false)
	}
}

func min3(a, b int) int {
	if a < b {
		return a
	}
	return b
}
