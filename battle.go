package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// BattleState manages the entire battle screen
type BattleState struct {
	Stage        *StageDef
	Board        *Board
	Shop         *Shop
	WaveMgr      *WaveManager
	Enemies      []*Enemy
	Units        []*Unit
	Projectiles  []*Projectile
	Integrity    int
	MaxIntegrity int
	Phase        BattlePhase
	Tick         int
	Rng          *rand.Rand

	// UI state
	SelectedUnit  *Unit
	DraggingUnit  *Unit
	DragFromBench bool
	DragFromGrid  bool
	DragOrigX     int
	DragOrigY     int

	// Buttons
	BtnStartWave Button
	BtnReroll    Button
	BtnLevelUp   Button
	BtnSell      Button

	// Result
	Victory bool
	GameOver bool

	// Barrier
	BarrierCooldown float64
	BarrierActive   float64

	// Stats
	KillCount   int
	WaveTime    float64
}

func NewBattleState(stage *StageDef) *BattleState {
	rng := rand.New(rand.NewSource(rand.Int63()))
	board := NewBoard(stage)
	shop := NewShop(stage, rng)
	waveMgr := NewWaveManager(stage, board)

	return &BattleState{
		Stage:        stage,
		Board:        board,
		Shop:         shop,
		WaveMgr:      waveMgr,
		Enemies:      make([]*Enemy, 0),
		Units:        make([]*Unit, 0),
		Projectiles:  make([]*Projectile, 0),
		Integrity:    stage.Integrity,
		MaxIntegrity: stage.Integrity,
		Phase:        PhasePrepare,
		Rng:          rng,
	}
}

func (b *BattleState) Update() {
	b.Tick++

	// Handle input
	b.handleInput()

	if b.GameOver {
		return
	}

	// Update wave spawning
	if b.WaveMgr.WaveActive {
		b.WaveTime += 1.0 / 60.0
		newEnemies := b.WaveMgr.Update()
		b.Enemies = append(b.Enemies, newEnemies...)
	}

	// Update enemies
	for _, e := range b.Enemies {
		e.Update(b.Board)
		if e.Reached && e.Alive {
			b.Integrity -= e.Def.LeakDamage
			e.Alive = false
			if b.Integrity <= 0 {
				b.Integrity = 0
				b.GameOver = true
				b.Victory = false
			}
		}
	}

	// Update units (combat)
	for _, u := range b.Units {
		u.Update(b.Enemies, b.Board, &b.Projectiles)
	}

	// Update projectiles
	UpdateProjectiles(b.Projectiles, b.Enemies)

	// Clean up dead projectiles
	alive := make([]*Projectile, 0, len(b.Projectiles))
	for _, p := range b.Projectiles {
		if p.Alive {
			alive = append(alive, p)
		}
	}
	b.Projectiles = alive

	// Count kills and award gold
	for _, e := range b.Enemies {
		if !e.Alive && !e.Reached && e.HP <= 0 {
			// Mark as processed
			e.Reached = true // reuse flag to prevent double-counting
			b.KillCount++
			b.Shop.AddGold(1) // 1 gold per kill
		}
	}

	// Update barrier
	if b.BarrierActive > 0 {
		b.BarrierActive -= 1.0 / 60.0
	}
	if b.BarrierCooldown > 0 {
		b.BarrierCooldown -= 1.0 / 60.0
	}

	// Check barrier activation
	if b.Stage.NodesEnabled && b.BarrierCooldown <= 0 {
		occupied := b.GetOccupiedNodes()
		if len(occupied) >= 3 {
			b.ActivateBarrier()
		}
	}

	// Check wave end — give bonus gold and refresh shop
	if !b.WaveMgr.WaveActive && b.Phase == PhaseWave {
		b.Phase = PhasePrepare
		b.Shop.AddGold(3 + b.WaveMgr.CurrentWave) // wave bonus
		b.Shop.Refresh()
	}

	// Check victory
	if b.WaveMgr.AllDone && !b.GameOver {
		allDead := true
		for _, e := range b.Enemies {
			if e.Alive && !e.Reached {
				allDead = false
				break
			}
		}
		if allDead {
			b.GameOver = true
			b.Victory = true
		}
	}
}

func (b *BattleState) handleInput() {
	mx, my := ebiten.CursorPosition()

	// Update button hover states
	b.BtnStartWave.Hovered = b.BtnStartWave.Contains(mx, my)
	b.BtnReroll.Hovered = b.BtnReroll.Contains(mx, my)
	b.BtnLevelUp.Hovered = b.BtnLevelUp.Contains(mx, my)
	b.BtnSell.Hovered = b.BtnSell.Contains(mx, my)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		b.handleClick(mx, my)
	}

	// Right click to deselect
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		b.SelectedUnit = nil
		b.DraggingUnit = nil
	}
}

func (b *BattleState) handleClick(mx, my int) {
	// Check buttons first
	if b.BtnStartWave.Contains(mx, my) && !b.BtnStartWave.Disabled && !b.WaveMgr.WaveActive {
		b.WaveMgr.StartWave()
		b.Phase = PhaseWave
		return
	}

	if b.BtnReroll.Contains(mx, my) && !b.BtnReroll.Disabled {
		b.Shop.Reroll()
		return
	}

	if b.BtnLevelUp.Contains(mx, my) && !b.BtnLevelUp.Disabled {
		b.Shop.LevelUp()
		return
	}

	if b.BtnSell.Contains(mx, my) && !b.BtnSell.Disabled && b.SelectedUnit != nil {
		b.SellSelectedUnit()
		return
	}

	// Check shop slot clicks
	shopY := float64(ScreenHeight - 110)
	if float64(my) >= shopY && float64(my) < shopY+62 {
		for i := 0; i < ShopSlots; i++ {
			slotX := 80 + i*150
			if mx >= slotX && mx < slotX+140 {
				if b.Shop.CanBuy(i) {
					b.BuyUnit(i)
				}
				return
			}
		}
	}

	// Check bench clicks
	bY := benchSlotY()
	if my >= bY && my < bY+50 {
		for i := 0; i < BenchSlots; i++ {
			bx := benchSlotX(i)
			if mx >= bx && mx < bx+50 {
				// Select bench unit
				for _, u := range b.Units {
					if u.BenchSlot == i {
						if b.SelectedUnit == u {
							// Double click = deselect
							b.SelectedUnit = nil
						} else {
							b.SelectedUnit = u
						}
						return
					}
				}
				// Empty bench slot - if we have a selected deployed unit, move to bench
				if b.SelectedUnit != nil && b.SelectedUnit.Deployed {
					b.SelectedUnit.PlaceBench(i)
					b.SelectedUnit = nil
				}
				return
			}
		}
	}

	// Check board clicks
	gx, gy := b.Board.ScreenToGrid(mx, my)
	if gx >= 0 && gx < BoardCols && gy >= 0 && gy < BoardRows {
		// Check if clicking on a deployed unit
		for _, u := range b.Units {
			if u.Deployed && u.GridX == gx && u.GridY == gy {
				if b.SelectedUnit == u {
					b.SelectedUnit = nil
				} else {
					b.SelectedUnit = u
				}
				return
			}
		}

		// If we have a selected bench unit, deploy it
		if b.SelectedUnit != nil && !b.SelectedUnit.Deployed && b.Board.CanPlace(gx, gy) {
			if b.DeployedCount() < b.Shop.DeployCap {
				// Check no other unit is there
				occupied := false
				for _, u := range b.Units {
					if u.Deployed && u.GridX == gx && u.GridY == gy {
						occupied = true
						break
					}
				}
				if !occupied {
					b.SelectedUnit.Place(gx, gy)
					b.SelectedUnit = nil
				}
			}
			return
		}

		// If we have a selected deployed unit, move it
		if b.SelectedUnit != nil && b.SelectedUnit.Deployed && b.Board.CanPlace(gx, gy) {
			occupied := false
			for _, u := range b.Units {
				if u.Deployed && u.GridX == gx && u.GridY == gy {
					occupied = true
					break
				}
			}
			if !occupied {
				b.SelectedUnit.Place(gx, gy)
				b.SelectedUnit = nil
			}
			return
		}

		// Click on empty space deselects
		b.SelectedUnit = nil
	}
}

func (b *BattleState) BuyUnit(slot int) {
	// Find free bench slot
	benchSlot := -1
	for i := 0; i < BenchSlots; i++ {
		occupied := false
		for _, u := range b.Units {
			if u.BenchSlot == i {
				occupied = true
				break
			}
		}
		if !occupied {
			benchSlot = i
			break
		}
	}

	if benchSlot == -1 {
		return // Bench full
	}

	unit := b.Shop.Buy(slot)
	if unit == nil {
		return
	}

	unit.PlaceBench(benchSlot)
	b.Units = append(b.Units, unit)

	// Check for TRI-FUSE (3 same units = upgrade)
	if b.Stage.TriFuseEnabled {
		b.CheckTriFuse(unit.Def.ID)
	}
}

func (b *BattleState) SellSelectedUnit() {
	if b.SelectedUnit == nil {
		return
	}
	b.Shop.SellUnit(b.SelectedUnit)

	// Remove from units list
	for i, u := range b.Units {
		if u == b.SelectedUnit {
			b.Units = append(b.Units[:i], b.Units[i+1:]...)
			break
		}
	}
	b.SelectedUnit = nil
}

func (b *BattleState) DeployedCount() int {
	count := 0
	for _, u := range b.Units {
		if u.Deployed {
			count++
		}
	}
	return count
}

func (b *BattleState) CountSynergies() (map[Faction]int, map[UnitClass]int) {
	factions := make(map[Faction]int)
	classes := make(map[UnitClass]int)
	for _, u := range b.Units {
		if u.Deployed {
			factions[u.Def.Faction]++
			classes[u.Def.Class]++
		}
	}
	return factions, classes
}

func (b *BattleState) GetOccupiedNodes() []Pos {
	var occupied []Pos
	for node := range b.Board.NodeSet {
		for _, u := range b.Units {
			if u.Deployed && u.GridX == node.X && u.GridY == node.Y {
				occupied = append(occupied, node)
				break
			}
		}
	}
	return occupied
}

func (b *BattleState) ActivateBarrier() {
	b.BarrierCooldown = 20.0 // 20 second cooldown
	b.BarrierActive = 3.0    // 3 second duration

	switch b.Stage.BarrierEffect {
	case "BARRIER_SLOW":
		for _, e := range b.Enemies {
			if e.Alive && !e.Reached {
				e.SlowTimer = 3.0
			}
		}
	case "BARRIER_MARK":
		// Increase damage taken (simplified: reduce HP slightly)
		for _, e := range b.Enemies {
			if e.Alive && !e.Reached {
				e.TakeDamage(e.MaxHP*0.05, DamageMagic)
			}
		}
	case "BARRIER_REVEAL":
		for _, e := range b.Enemies {
			if e.Alive && !e.Reached && e.Def.Type == EnemyStalker {
				e.Visible = true
			}
		}
	}
}

func (b *BattleState) CheckTriFuse(unitID string) {
	// Find all units with same ID and star level
	for star := 1; star <= 2; star++ {
		var matching []*Unit
		for _, u := range b.Units {
			if u.Def.ID == unitID && u.Star == star {
				matching = append(matching, u)
			}
		}

		if len(matching) >= 3 {
			// Fuse! Keep the first one, remove the other two
			keeper := matching[0]
			keeper.Star++
			keeper.MaxHP *= 1.6
			keeper.HP = keeper.MaxHP
			keeper.ATK *= 1.35
			// Remove the other 2
			for _, rm := range matching[1:3] {
				for i, u := range b.Units {
					if u == rm {
						b.Units = append(b.Units[:i], b.Units[i+1:]...)
						break
					}
				}
			}
			break
		}
	}
}

func (b *BattleState) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(ColorBG)

	// Board (with placement highlights when selecting a bench unit)
	highlightPlaceable := b.SelectedUnit != nil && !b.SelectedUnit.Deployed
	var occupiedTiles map[Pos]bool
	if highlightPlaceable {
		occupiedTiles = make(map[Pos]bool)
		for _, u := range b.Units {
			if u.Deployed {
				occupiedTiles[Pos{u.GridX, u.GridY}] = true
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
		sx := float32(BoardOffsetX+b.SelectedUnit.GridX*TileSize) + float32(TileSize)/2
		sy := float32(BoardOffsetY+b.SelectedUnit.GridY*TileSize) + float32(TileSize)/2
		s := float32(TileSize/2 + 2)
		vector.StrokeRect(screen, sx-s, sy-s, s*2, s*2, 2, ColorNeonCyan, false)
	}

	// Projectiles
	DrawProjectiles(screen, b.Projectiles)

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

func drawRangeIndicator(screen *ebiten.Image, u *Unit) {
	cx := float32(BoardOffsetX+u.GridX*TileSize) + float32(TileSize)/2
	cy := float32(BoardOffsetY+u.GridY*TileSize) + float32(TileSize)/2
	r := float32(u.Range*TileSize) + float32(TileSize)/2
	vector.StrokeCircle(screen, cx, cy, r, 1, withAlpha(ColorNeonCyan, 60), false)
}

func drawBarrierEffect(screen *ebiten.Image, battle *BattleState, tick int) {
	// Full screen overlay flash
	alpha := uint8(battle.BarrierActive / 3.0 * 30)
	vector.DrawFilledRect(screen, float32(BoardOffsetX), float32(BoardOffsetY),
		float32(BoardCols*TileSize), float32(BoardRows*TileSize),
		withAlpha(ColorNeonBlue, alpha), false)
}

func drawGameOverlay(screen *ebiten.Image, battle *BattleState, tick int) {
	// Dark overlay
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, withAlpha(ColorBG, 180), false)

	if battle.Victory {
		DrawTextGlowCentered(screen, "STAGE CLEAR", fontBold(36), ScreenWidth/2, ScreenHeight/2-40, ColorNeonCyan)
		DrawTextCentered(screen, "Press ENTER to continue", fontRegular(14), ScreenWidth/2, ScreenHeight/2+30, ColorWhiteDim)
	} else {
		DrawTextGlowCentered(screen, "BREACH DETECTED", fontBold(36), ScreenWidth/2, ScreenHeight/2-40, ColorNeonRed)
		DrawTextCentered(screen, "Press ENTER to retry", fontRegular(14), ScreenWidth/2, ScreenHeight/2+30, ColorWhiteDim)
	}

	// Stats
	statsY := float64(ScreenHeight/2 + 70)
	DrawTextCentered(screen, fmt.Sprintf("Kills: %d   Waves: %d/%d",
		battle.KillCount, battle.WaveMgr.CurrentWave, battle.WaveMgr.TotalWaves()),
		fontRegular(11), ScreenWidth/2, statsY, ColorWhiteDim)
}

func withAlpha(c color.RGBA, a uint8) color.RGBA {
	return color.RGBA{c.R, c.G, c.B, a}
}
