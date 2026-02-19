package battle

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"neonsigil/internal/config"
	"neonsigil/internal/entity"
	"neonsigil/internal/ui"
)

// Draw renders the entire battle screen
func (b *BattleState) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(config.ColorBG)

	// Board (with placement highlights when selecting a bench unit)
	highlightPlaceable := b.SelectedUnit != nil && !b.SelectedUnit.Deployed
	var occupiedTiles map[config.Pos]bool
	if highlightPlaceable {
		occupiedTiles = make(map[config.Pos]bool)
		for _, u := range b.Units {
			if u.Deployed {
				occupiedTiles[config.Pos{X: u.GridX, Y: u.GridY}] = true
			}
		}
	}
	b.Board.DrawWithHighlight(screen, b.Tick, highlightPlaceable, occupiedTiles)

	// Node connections
	DrawNodeIndicator(screen, b, b.Tick)

	// Barrier effect visual
	if b.BarrierActive > 0 {
		drawBarrierEffect(screen, b, b.Tick)
	}

	// Selected unit range indicator
	if b.SelectedUnit != nil && b.SelectedUnit.Deployed {
		drawRangeIndicator(screen, b.SelectedUnit)
	}

	// Enemies
	for _, e := range b.Enemies {
		e.Draw(screen, b.Tick)
	}

	// Units on board
	for _, u := range b.Units {
		u.Draw(screen, b.Tick)
	}

	// Selected unit highlight
	if b.SelectedUnit != nil && b.SelectedUnit.Deployed {
		sx := float32(config.BoardOffsetX+b.SelectedUnit.GridX*config.TileSize) + float32(config.TileSize)/2
		sy := float32(config.BoardOffsetY+b.SelectedUnit.GridY*config.TileSize) + float32(config.TileSize)/2
		s := float32(config.TileSize/2 + 2)
		vector.StrokeRect(screen, sx-s, sy-s, s*2, s*2, 2, config.ColorNeonCyan, false)
	}

	// Projectiles
	entity.DrawProjectiles(screen, b.Projectiles)

	// UI
	DrawHUD(screen, b, b.Tick)
	DrawBenchUI(screen, b, b.Tick)
	DrawShopUI(screen, b, b.Tick)
	DrawInfoPanel(screen, b, b.Tick)

	// Game over overlay
	if b.GameOver {
		drawGameOverlay(screen, b, b.Tick)
	}
}

func drawRangeIndicator(screen *ebiten.Image, u *entity.Unit) {
	cx := float32(config.BoardOffsetX+u.GridX*config.TileSize) + float32(config.TileSize)/2
	cy := float32(config.BoardOffsetY+u.GridY*config.TileSize) + float32(config.TileSize)/2
	r := float32(u.Range*config.TileSize) + float32(config.TileSize)/2
	vector.StrokeCircle(screen, cx, cy, r, 1, config.WithAlpha(config.ColorNeonCyan, 60), false)
}

func drawBarrierEffect(screen *ebiten.Image, battle *BattleState, tick int) {
	// Full screen overlay flash
	alpha := uint8(battle.BarrierActive / 3.0 * 30)
	vector.DrawFilledRect(screen, float32(config.BoardOffsetX), float32(config.BoardOffsetY),
		float32(config.BoardCols*config.TileSize), float32(config.BoardRows*config.TileSize),
		config.WithAlpha(config.ColorNeonBlue, alpha), false)
}

func drawGameOverlay(screen *ebiten.Image, battle *BattleState, tick int) {
	// Dark overlay
	vector.DrawFilledRect(screen, 0, 0, config.ScreenWidth, config.ScreenHeight, config.WithAlpha(config.ColorBG, 180), false)

	if battle.Victory {
		ui.DrawTextGlowCentered(screen, "STAGE CLEAR", ui.FontBold(36), config.ScreenWidth/2, config.ScreenHeight/2-40, config.ColorNeonCyan)
		ui.DrawTextCentered(screen, "Press ENTER to continue", ui.FontRegular(14), config.ScreenWidth/2, config.ScreenHeight/2+30, config.ColorWhiteDim)
	} else {
		ui.DrawTextGlowCentered(screen, "BREACH DETECTED", ui.FontBold(36), config.ScreenWidth/2, config.ScreenHeight/2-40, config.ColorNeonRed)
		ui.DrawTextCentered(screen, "Press ENTER to retry", ui.FontRegular(14), config.ScreenWidth/2, config.ScreenHeight/2+30, config.ColorWhiteDim)
	}

	// Stats
	statsY := float64(config.ScreenHeight/2 + 70)
	ui.DrawTextCentered(screen, fmt.Sprintf("Kills: %d   Waves: %d/%d",
		battle.KillCount, battle.WaveMgr.CurrentWave, battle.WaveMgr.TotalWaves()),
		ui.FontRegular(11), config.ScreenWidth/2, statsY, config.ColorWhiteDim)
}
