package battle

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"neonsigil/internal/board"
	"neonsigil/internal/config"
	"neonsigil/internal/data"
	"neonsigil/internal/entity"
	"neonsigil/internal/shop"
	"neonsigil/internal/ui"
	"neonsigil/internal/wave"
)

// BattleState manages the entire battle screen
type BattleState struct {
	Stage        *data.StageDef
	Board        *board.Board
	Shop         *shop.Shop
	WaveMgr      *wave.WaveManager
	Enemies      []*entity.Enemy
	Units        []*entity.Unit
	Projectiles  []*entity.Projectile
	Integrity    int
	MaxIntegrity int
	Phase        config.BattlePhase
	Tick         int
	Rng          *rand.Rand

	// UI state
	SelectedUnit  *entity.Unit
	DraggingUnit  *entity.Unit
	DragFromBench bool
	DragFromGrid  bool
	DragOrigX     int
	DragOrigY     int

	// Buttons
	BtnStartWave ui.Button
	BtnReroll    ui.Button
	BtnLevelUp   ui.Button
	BtnSell      ui.Button

	// Result
	Victory  bool
	GameOver bool

	// Barrier
	BarrierCooldown float64
	BarrierActive   float64

	// Stats
	KillCount int
	WaveTime  float64
}

// NewBattleState creates a new battle state for the given stage
func NewBattleState(stage *data.StageDef) *BattleState {
	rng := rand.New(rand.NewSource(rand.Int63()))
	b := board.NewBoard(stage)
	s := shop.NewShop(stage, rng)
	waveMgr := wave.NewWaveManager(stage, b)

	return &BattleState{
		Stage:        stage,
		Board:        b,
		Shop:         s,
		WaveMgr:      waveMgr,
		Enemies:      make([]*entity.Enemy, 0),
		Units:        make([]*entity.Unit, 0),
		Projectiles:  make([]*entity.Projectile, 0),
		Integrity:    stage.Integrity,
		MaxIntegrity: stage.Integrity,
		Phase:        config.PhasePrepare,
		Rng:          rng,
	}
}

// Update runs one frame of battle logic
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
	entity.UpdateProjectiles(b.Projectiles, b.Enemies)

	// Clean up dead projectiles
	alive := make([]*entity.Projectile, 0, len(b.Projectiles))
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
	if !b.WaveMgr.WaveActive && b.Phase == config.PhaseWave {
		b.Phase = config.PhasePrepare
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
		b.Phase = config.PhaseWave
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
	shopY := float64(config.ScreenHeight - 110)
	if float64(my) >= shopY && float64(my) < shopY+62 {
		for i := 0; i < config.ShopSlots; i++ {
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
	bY := config.BenchSlotY()
	if my >= bY && my < bY+50 {
		for i := 0; i < config.BenchSlots; i++ {
			bx := config.BenchSlotX(i)
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
	if gx >= 0 && gx < config.BoardCols && gy >= 0 && gy < config.BoardRows {
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

// BuyUnit purchases a unit from the shop and places it on the bench
func (b *BattleState) BuyUnit(slot int) {
	// Find free bench slot
	benchSlot := -1
	for i := 0; i < config.BenchSlots; i++ {
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

// SellSelectedUnit sells the currently selected unit
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

// DeployedCount returns the number of deployed units
func (b *BattleState) DeployedCount() int {
	count := 0
	for _, u := range b.Units {
		if u.Deployed {
			count++
		}
	}
	return count
}

// CountSynergies counts deployed faction and class synergies
func (b *BattleState) CountSynergies() (map[config.Faction]int, map[config.UnitClass]int) {
	factions := make(map[config.Faction]int)
	classes := make(map[config.UnitClass]int)
	for _, u := range b.Units {
		if u.Deployed {
			factions[u.Def.Faction]++
			classes[u.Def.Class]++
		}
	}
	return factions, classes
}

// GetOccupiedNodes returns node positions that have units on them
func (b *BattleState) GetOccupiedNodes() []config.Pos {
	var occupied []config.Pos
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

// ActivateBarrier activates the barrier effect
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
				e.TakeDamage(e.MaxHP*0.05, config.DamageMagic)
			}
		}
	case "BARRIER_REVEAL":
		for _, e := range b.Enemies {
			if e.Alive && !e.Reached && e.Def.Type == config.EnemyStalker {
				e.Visible = true
			}
		}
	}
}

// CheckTriFuse checks and performs TRI-FUSE combination
func (b *BattleState) CheckTriFuse(unitID string) {
	// Find all units with same ID and star level
	for star := 1; star <= 2; star++ {
		var matching []*entity.Unit
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
