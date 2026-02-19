package entity

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"neonsigil/internal/board"
	"neonsigil/internal/config"
	"neonsigil/internal/data"
)

// Unit is a live unit on the board or bench
type Unit struct {
	Def         *data.UnitDef
	Star        int
	GridX, GridY int // -1 if on bench
	BenchSlot   int  // -1 if on board
	HP          float64
	MaxHP       float64
	ATK         float64
	AtkSpeed    float64
	Range       int
	AtkCooldown float64
	Deployed    bool
}

// NewUnit creates a new unit from a definition
func NewUnit(def *data.UnitDef) *Unit {
	return &Unit{
		Def:       def,
		Star:      1,
		GridX:     -1,
		GridY:     -1,
		BenchSlot: -1,
		HP:        def.HP,
		MaxHP:     def.HP,
		ATK:       def.ATK,
		AtkSpeed:  def.AtkSpeed,
		Range:     def.Range,
		Deployed:  false,
	}
}

// Place deploys the unit to the board at the given grid position
func (u *Unit) Place(gx, gy int) {
	u.GridX = gx
	u.GridY = gy
	u.BenchSlot = -1
	u.Deployed = true
}

// PlaceBench moves the unit to the bench
func (u *Unit) PlaceBench(slot int) {
	u.GridX = -1
	u.GridY = -1
	u.BenchSlot = slot
	u.Deployed = false
}

// FindTarget finds the best target enemy based on the unit's targeting mode
func (u *Unit) FindTarget(enemies []*Enemy, b *board.Board) *Enemy {
	if !u.Deployed {
		return nil
	}

	unitPx := float64(config.BoardOffsetX+u.GridX*config.TileSize) + float64(config.TileSize)/2
	unitPy := float64(config.BoardOffsetY+u.GridY*config.TileSize) + float64(config.TileSize)/2
	rangePixels := float64(u.Range) * float64(config.TileSize)

	var best *Enemy
	var bestScore float64

	for _, e := range enemies {
		if !e.Alive || e.Reached {
			continue
		}
		if !e.Visible && e.Def.Type == config.EnemyStalker {
			continue
		}

		dx := e.Pos.X - unitPx
		dy := e.Pos.Y - unitPy
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > rangePixels+float64(config.TileSize)/2 {
			continue
		}

		var score float64
		switch u.Def.Targeting {
		case config.TargetFrontmost:
			score = e.GetProgress(b)
		case config.TargetLowHP:
			score = 1.0 - (e.HP / e.MaxHP)
		case config.TargetNearest:
			score = 1.0 - (dist / rangePixels)
		}

		if best == nil || score > bestScore {
			best = e
			bestScore = score
		}
	}

	return best
}

// Update runs the unit's combat logic
func (u *Unit) Update(enemies []*Enemy, b *board.Board, projectiles *[]*Projectile) {
	if !u.Deployed {
		return
	}

	u.AtkCooldown -= 1.0 / 60.0
	if u.AtkCooldown > 0 {
		return
	}

	target := u.FindTarget(enemies, b)
	if target == nil {
		return
	}

	u.AtkCooldown = 1.0 / u.AtkSpeed

	unitPx := float64(config.BoardOffsetX+u.GridX*config.TileSize) + float64(config.TileSize)/2
	unitPy := float64(config.BoardOffsetY+u.GridY*config.TileSize) + float64(config.TileSize)/2

	if u.Def.AtkType == config.AttackMelee {
		// Instant damage
		target.TakeDamage(u.ATK, u.Def.DmgType)
	} else {
		// Spawn projectile
		for i, e := range enemies {
			if e == target {
				proj := &Projectile{
					X:        unitPx,
					Y:        unitPy,
					TargetID: i,
					Damage:   u.ATK,
					Speed:    400.0,
					Alive:    true,
				}
				*projectiles = append(*projectiles, proj)
				break
			}
		}
	}
}

// Draw renders the unit on the board
func (u *Unit) Draw(screen *ebiten.Image, tick int) {
	if !u.Deployed {
		return
	}

	sx := float32(config.BoardOffsetX+u.GridX*config.TileSize) + float32(config.TileSize)/2
	sy := float32(config.BoardOffsetY+u.GridY*config.TileSize) + float32(config.TileSize)/2
	s := float32(config.TileSize/2 - 6)

	// Faction color base
	fc := config.FactionColors[u.Def.Faction]

	// Unit body (rounded square)
	vector.DrawFilledRect(screen, sx-s, sy-s, s*2, s*2, color.RGBA{fc.R / 3, fc.G / 3, fc.B / 3, 240}, false)
	vector.StrokeRect(screen, sx-s, sy-s, s*2, s*2, 2, fc, false)

	// Class indicator (inner shape)
	cc := config.ClassColors[u.Def.Class]
	innerS := float32(8)
	switch u.Def.Class {
	case config.ClassVanguard:
		// Shield shape
		vector.DrawFilledRect(screen, sx-innerS, sy-innerS, innerS*2, innerS*2, cc, false)
	case config.ClassMarksman:
		// Cross
		vector.DrawFilledRect(screen, sx-1, sy-innerS, 2, innerS*2, cc, false)
		vector.DrawFilledRect(screen, sx-innerS, sy-1, innerS*2, 2, cc, false)
	case config.ClassCaster:
		// Circle
		vector.DrawFilledCircle(screen, sx, sy, innerS, cc, false)
	case config.ClassEngineer:
		// Gear (hexagon-ish)
		vector.StrokeCircle(screen, sx, sy, innerS, 2, cc, false)
		vector.DrawFilledCircle(screen, sx, sy, innerS*0.5, cc, false)
	case config.ClassSupport:
		// Plus
		vector.DrawFilledRect(screen, sx-innerS, sy-2, innerS*2, 4, cc, false)
		vector.DrawFilledRect(screen, sx-2, sy-innerS, 4, innerS*2, cc, false)
	}

	// Star indicator
	starY := sy + s + 6
	for i := 0; i < u.Star; i++ {
		starX := sx - float32(u.Star-1)*5 + float32(i)*10
		vector.DrawFilledCircle(screen, starX, starY, 3, config.ColorNeonYellow, false)
	}

	// Attack cooldown indicator (small bar at bottom)
	if u.AtkCooldown > 0 {
		barW := s * 2
		ratio := float32(u.AtkCooldown / (1.0 / u.AtkSpeed))
		if ratio > 1 {
			ratio = 1
		}
		vector.DrawFilledRect(screen, sx-s, sy+s-2, barW*(1-ratio), 2, config.ColorNeonCyan, false)
	}
}

// DrawOnBench renders the unit in a bench slot
func (u *Unit) DrawOnBench(screen *ebiten.Image, slot int, tick int) {
	bx := float32(config.BenchSlotX(slot))
	by := float32(config.BenchSlotY())
	s := float32(25)

	fc := config.FactionColors[u.Def.Faction]

	// Background
	vector.DrawFilledRect(screen, bx, by, s*2, s*2, color.RGBA{fc.R / 4, fc.G / 4, fc.B / 4, 220}, false)
	vector.StrokeRect(screen, bx, by, s*2, s*2, 1.5, fc, false)

	// Class indicator
	cc := config.ClassColors[u.Def.Class]
	cx := bx + s
	cy := by + s
	vector.DrawFilledCircle(screen, cx, cy, 8, cc, false)

	// Star
	for i := 0; i < u.Star; i++ {
		starX := cx - float32(u.Star-1)*4 + float32(i)*8
		vector.DrawFilledCircle(screen, starX, by+s*2+6, 2.5, config.ColorNeonYellow, false)
	}
}
