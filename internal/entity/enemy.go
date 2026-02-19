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

// Enemy is a live enemy on the board
type Enemy struct {
	Def         *data.EnemyDef
	HP          float64
	MaxHP       float64
	Pos         config.FPos // pixel position
	PathID      string
	WaypointIdx int
	Speed       float64
	Alive       bool
	Reached     bool // reached the end
	SlowTimer   float64
	StunTimer   float64
	Visible     bool // for STALKER
}

// NewEnemy creates a new enemy from a definition
func NewEnemy(def *data.EnemyDef, pathID string, hpMul, spdMul float64, b *board.Board) *Enemy {
	// Find starting position from path
	var startPos config.FPos
	for _, pd := range b.PathDefs {
		if pd.ID == pathID && len(pd.Waypoints) > 0 {
			wp := pd.Waypoints[0]
			startPos = config.FPos{
				X: float64(config.BoardOffsetX+wp.X*config.TileSize) + float64(config.TileSize)/2,
				Y: float64(config.BoardOffsetY+wp.Y*config.TileSize) + float64(config.TileSize)/2,
			}
			break
		}
	}

	hp := def.BaseHP * hpMul
	visible := true
	if def.Type == config.EnemyStalker {
		visible = false
	}

	return &Enemy{
		Def:         def,
		HP:          hp,
		MaxHP:       hp,
		Pos:         startPos,
		PathID:      pathID,
		WaypointIdx: 0,
		Speed:       def.Speed * spdMul,
		Alive:       true,
		Visible:     visible,
	}
}

// Update moves the enemy along its path
func (e *Enemy) Update(b *board.Board) {
	if !e.Alive || e.Reached {
		return
	}

	// Stun check
	if e.StunTimer > 0 {
		e.StunTimer -= 1.0 / 60.0
		return
	}

	// Find current path
	var path *data.PathDef
	for i := range b.PathDefs {
		if b.PathDefs[i].ID == e.PathID {
			path = &b.PathDefs[i]
			break
		}
	}
	if path == nil {
		return
	}

	if e.WaypointIdx >= len(path.Waypoints)-1 {
		e.Reached = true
		return
	}

	// Move toward next waypoint
	nextWP := path.Waypoints[e.WaypointIdx+1]
	targetX := float64(config.BoardOffsetX+nextWP.X*config.TileSize) + float64(config.TileSize)/2
	targetY := float64(config.BoardOffsetY+nextWP.Y*config.TileSize) + float64(config.TileSize)/2

	dx := targetX - e.Pos.X
	dy := targetY - e.Pos.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	// Apply slow
	speed := e.Speed
	if e.SlowTimer > 0 {
		speed *= 0.6
		e.SlowTimer -= 1.0 / 60.0
	}

	moveSpeed := speed * 60.0 // pixels per second at speed 1.0 = 60px/s
	movePerFrame := moveSpeed / 60.0

	if dist <= movePerFrame {
		e.Pos.X = targetX
		e.Pos.Y = targetY
		e.WaypointIdx++
	} else {
		e.Pos.X += (dx / dist) * movePerFrame
		e.Pos.Y += (dy / dist) * movePerFrame
	}
}

// TakeDamage applies damage to the enemy
func (e *Enemy) TakeDamage(dmg float64, dmgType config.DamageType) {
	if !e.Alive {
		return
	}
	actualDmg := dmg
	// Shield enemies reduce ranged/phys damage
	if e.Def.ShieldPct > 0 && dmgType == config.DamagePhys {
		actualDmg *= (1.0 - e.Def.ShieldPct)
	}
	e.HP -= actualDmg
	if e.HP <= 0 {
		e.HP = 0
		e.Alive = false
	}
}

// GetProgress returns how far along the path this enemy is (0.0 to 1.0)
func (e *Enemy) GetProgress(b *board.Board) float64 {
	var path *data.PathDef
	for i := range b.PathDefs {
		if b.PathDefs[i].ID == e.PathID {
			path = &b.PathDefs[i]
			break
		}
	}
	if path == nil {
		return 0
	}
	totalWP := len(path.Waypoints)
	if totalWP <= 1 {
		return 0
	}
	return float64(e.WaypointIdx) / float64(totalWP-1)
}

// Draw renders the enemy on screen
func (e *Enemy) Draw(screen *ebiten.Image, tick int) {
	if !e.Alive || e.Reached {
		return
	}

	x := float32(e.Pos.X)
	y := float32(e.Pos.Y)
	r := float32(10)

	// Draw body based on enemy type
	c := e.Def.Color
	if !e.Visible {
		// Stalker invisible: show faint shimmer
		pulse := math.Sin(float64(tick%30)/30.0*math.Pi*2)*0.3 + 0.3
		c = color.RGBA{c.R, c.G, c.B, uint8(float64(60) * pulse)}
	}

	switch e.Def.Type {
	case config.EnemyBruiser, config.EnemyBoss:
		// Square for tanky
		r2 := r * 1.3
		vector.DrawFilledRect(screen, x-r2, y-r2, r2*2, r2*2, c, false)
		vector.StrokeRect(screen, x-r2, y-r2, r2*2, r2*2, 1.5, brighten(c, 0.5), false)
	case config.EnemyShield:
		// Hexagon-ish
		vector.DrawFilledCircle(screen, x, y, r*1.2, c, false)
		vector.StrokeCircle(screen, x, y, r*1.2, 2, brighten(c, 0.4), false)
		// Shield indicator
		vector.StrokeCircle(screen, x, y, r*0.6, 1.5, color.RGBA{200, 230, 255, 200}, false)
	case config.EnemyFlyer:
		// Diamond for flyer
		var path vector.Path
		path.MoveTo(x, y-r*1.3)
		path.LineTo(x+r, y)
		path.LineTo(x, y+r*1.3)
		path.LineTo(x-r, y)
		path.Close()
		vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range vs {
			vs[i].ColorR = float32(c.R) / 255
			vs[i].ColorG = float32(c.G) / 255
			vs[i].ColorB = float32(c.B) / 255
			vs[i].ColorA = float32(c.A) / 255
		}
		screen.DrawTriangles(vs, is, board.EmptyImage, nil)
	case config.EnemySplitter:
		// Two small circles
		vector.DrawFilledCircle(screen, x-4, y, r*0.8, c, false)
		vector.DrawFilledCircle(screen, x+4, y, r*0.8, c, false)
	default:
		// Circle for basic
		vector.DrawFilledCircle(screen, x, y, r, c, false)
		if e.Def.Type != config.EnemyStalker || e.Visible {
			vector.StrokeCircle(screen, x, y, r, 1, brighten(c, 0.3), false)
		}
	}

	// HP bar
	if e.Visible || e.Def.Type != config.EnemyStalker {
		barW := float32(24)
		barH := float32(3)
		barX := x - barW/2
		barY := y - r - 8

		// Background
		vector.DrawFilledRect(screen, barX, barY, barW, barH, color.RGBA{40, 40, 40, 200}, false)
		// Fill
		ratio := float32(e.HP / e.MaxHP)
		hpColor := config.ColorHP
		if ratio < 0.3 {
			hpColor = config.ColorHPLow
		}
		vector.DrawFilledRect(screen, barX, barY, barW*ratio, barH, hpColor, false)
	}

	// Slow indicator
	if e.SlowTimer > 0 {
		vector.StrokeCircle(screen, x, y, r+3, 1, color.RGBA{0, 200, 255, 150}, false)
	}
}

func brighten(c color.RGBA, amount float64) color.RGBA {
	r := math.Min(float64(c.R)+255*amount, 255)
	g := math.Min(float64(c.G)+255*amount, 255)
	b := math.Min(float64(c.B)+255*amount, 255)
	return color.RGBA{uint8(r), uint8(g), uint8(b), c.A}
}
