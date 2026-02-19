package entity

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"neonsigil/internal/config"
)

// Projectile represents an in-flight projectile
type Projectile struct {
	X, Y     float64
	TargetID int
	Damage   float64
	Speed    float64
	Alive    bool
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
			target.TakeDamage(p.Damage, config.DamagePhys)
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
		vector.DrawFilledCircle(screen, float32(p.X), float32(p.Y), 3, config.ColorNeonCyan, false)
		// Trail
		vector.DrawFilledCircle(screen, float32(p.X-2), float32(p.Y-1), 2, color.RGBA{0, 200, 255, 100}, false)
	}
}
