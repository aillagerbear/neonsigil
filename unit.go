package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func NewUnit(def *UnitDef) *Unit {
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

func (u *Unit) Place(gx, gy int) {
	u.GridX = gx
	u.GridY = gy
	u.BenchSlot = -1
	u.Deployed = true
}

func (u *Unit) PlaceBench(slot int) {
	u.GridX = -1
	u.GridY = -1
	u.BenchSlot = slot
	u.Deployed = false
}

func (u *Unit) FindTarget(enemies []*Enemy, board *Board) *Enemy {
	if !u.Deployed {
		return nil
	}

	unitPx := float64(BoardOffsetX+u.GridX*TileSize) + float64(TileSize)/2
	unitPy := float64(BoardOffsetY+u.GridY*TileSize) + float64(TileSize)/2
	rangePixels := float64(u.Range) * float64(TileSize)

	var best *Enemy
	var bestScore float64

	for _, e := range enemies {
		if !e.Alive || e.Reached {
			continue
		}
		if !e.Visible && e.Def.Type == EnemyStalker {
			continue
		}

		dx := e.Pos.X - unitPx
		dy := e.Pos.Y - unitPy
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > rangePixels+float64(TileSize)/2 {
			continue
		}

		var score float64
		switch u.Def.Targeting {
		case TargetFrontmost:
			score = e.GetProgress(board)
		case TargetLowHP:
			score = 1.0 - (e.HP / e.MaxHP)
		case TargetNearest:
			score = 1.0 - (dist / rangePixels)
		}

		if best == nil || score > bestScore {
			best = e
			bestScore = score
		}
	}

	return best
}

func (u *Unit) Update(enemies []*Enemy, board *Board, projectiles *[]*Projectile) {
	if !u.Deployed {
		return
	}

	u.AtkCooldown -= 1.0 / 60.0
	if u.AtkCooldown > 0 {
		return
	}

	target := u.FindTarget(enemies, board)
	if target == nil {
		return
	}

	u.AtkCooldown = 1.0 / u.AtkSpeed

	unitPx := float64(BoardOffsetX+u.GridX*TileSize) + float64(TileSize)/2
	unitPy := float64(BoardOffsetY+u.GridY*TileSize) + float64(TileSize)/2

	if u.Def.AtkType == AttackMelee {
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

func (u *Unit) Draw(screen *ebiten.Image, tick int) {
	if !u.Deployed {
		return
	}

	sx := float32(BoardOffsetX+u.GridX*TileSize) + float32(TileSize)/2
	sy := float32(BoardOffsetY+u.GridY*TileSize) + float32(TileSize)/2
	s := float32(TileSize/2 - 6)

	// Faction color base
	fc := FactionColors[u.Def.Faction]

	// Unit body (rounded square)
	vector.DrawFilledRect(screen, sx-s, sy-s, s*2, s*2, color.RGBA{fc.R / 3, fc.G / 3, fc.B / 3, 240}, false)
	vector.StrokeRect(screen, sx-s, sy-s, s*2, s*2, 2, fc, false)

	// Class indicator (inner shape)
	cc := ClassColors[u.Def.Class]
	innerS := float32(8)
	switch u.Def.Class {
	case ClassVanguard:
		// Shield shape
		vector.DrawFilledRect(screen, sx-innerS, sy-innerS, innerS*2, innerS*2, cc, false)
	case ClassMarksman:
		// Cross
		vector.DrawFilledRect(screen, sx-1, sy-innerS, 2, innerS*2, cc, false)
		vector.DrawFilledRect(screen, sx-innerS, sy-1, innerS*2, 2, cc, false)
	case ClassCaster:
		// Circle
		vector.DrawFilledCircle(screen, sx, sy, innerS, cc, false)
	case ClassEngineer:
		// Gear (hexagon-ish)
		vector.StrokeCircle(screen, sx, sy, innerS, 2, cc, false)
		vector.DrawFilledCircle(screen, sx, sy, innerS*0.5, cc, false)
	case ClassSupport:
		// Plus
		vector.DrawFilledRect(screen, sx-innerS, sy-2, innerS*2, 4, cc, false)
		vector.DrawFilledRect(screen, sx-2, sy-innerS, 4, innerS*2, cc, false)
	}

	// Star indicator
	starY := sy + s + 6
	for i := 0; i < u.Star; i++ {
		starX := sx - float32(u.Star-1)*5 + float32(i)*10
		vector.DrawFilledCircle(screen, starX, starY, 3, ColorNeonYellow, false)
	}

	// Range indicator when selected (could add later)

	// Attack cooldown indicator (small bar at bottom)
	if u.AtkCooldown > 0 {
		barW := s * 2
		ratio := float32(u.AtkCooldown / (1.0 / u.AtkSpeed))
		if ratio > 1 {
			ratio = 1
		}
		vector.DrawFilledRect(screen, sx-s, sy+s-2, barW*(1-ratio), 2, ColorNeonCyan, false)
	}
}

func (u *Unit) DrawOnBench(screen *ebiten.Image, slot int, tick int) {
	bx := float32(benchSlotX(slot))
	by := float32(benchSlotY())
	s := float32(25)

	fc := FactionColors[u.Def.Faction]

	// Background
	vector.DrawFilledRect(screen, bx, by, s*2, s*2, color.RGBA{fc.R / 4, fc.G / 4, fc.B / 4, 220}, false)
	vector.StrokeRect(screen, bx, by, s*2, s*2, 1.5, fc, false)

	// Class indicator
	cc := ClassColors[u.Def.Class]
	cx := bx + s
	cy := by + s
	vector.DrawFilledCircle(screen, cx, cy, 8, cc, false)

	// Star
	for i := 0; i < u.Star; i++ {
		starX := cx - float32(u.Star-1)*4 + float32(i)*8
		vector.DrawFilledCircle(screen, starX, by+s*2+6, 2.5, ColorNeonYellow, false)
	}
}

func benchSlotX(slot int) int {
	return BoardOffsetX + slot*66
}

func benchSlotY() int {
	return BoardOffsetY + BoardRows*TileSize + 16
}

// UpdateProjectiles updates all projectiles
func UpdateProjectiles(projectiles []*Projectile, enemies []*Enemy) {
	for _, p := range projectiles {
		if !p.Alive {
			continue
		}

		if p.TargetID < 0 || p.TargetID >= len(enemies) {
			p.Alive = false
			continue
		}

		target := enemies[p.TargetID]
		if !target.Alive || target.Reached {
			p.Alive = false
			continue
		}

		dx := target.Pos.X - p.X
		dy := target.Pos.Y - p.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < 8 {
			target.TakeDamage(p.Damage, DamagePhys)
			p.Alive = false
			continue
		}

		speed := p.Speed / 60.0
		p.X += (dx / dist) * speed
		p.Y += (dy / dist) * speed
	}
}

// DrawProjectiles draws all projectiles
func DrawProjectiles(screen *ebiten.Image, projectiles []*Projectile) {
	for _, p := range projectiles {
		if !p.Alive {
			continue
		}
		vector.DrawFilledCircle(screen, float32(p.X), float32(p.Y), 3, ColorNeonCyan, false)
		// Trail
		vector.DrawFilledCircle(screen, float32(p.X-2), float32(p.Y-1), 2, color.RGBA{0, 200, 255, 100}, false)
	}
}
