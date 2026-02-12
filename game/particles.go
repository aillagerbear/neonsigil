package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Particle represents a single visual effect particle
type Particle struct {
	X, Y    float64
	VX, VY  float64
	Life    int
	MaxLife int
	Color   color.RGBA
	Size    float64
}

func (g *Game) updateParticles() {
	alive := g.particles[:0]
	for i := range g.particles {
		p := &g.particles[i]
		p.X += p.VX
		p.Y += p.VY
		p.VY += 0.05 // slight gravity
		p.Life--
		if p.Life > 0 {
			alive = append(alive, *p)
		}
	}
	g.particles = alive
}

func (g *Game) drawParticles(screen *ebiten.Image) {
	for i := range g.particles {
		p := &g.particles[i]
		alpha := float64(p.Life) / float64(p.MaxLife)
		a := byte(float64(p.Color.A) * alpha)
		if a < 10 {
			continue
		}
		c := color.RGBA{p.Color.R, p.Color.G, p.Color.B, a}
		sz := float32(p.Size * alpha)
		if sz < 1 {
			sz = 1
		}
		vector.FillRect(screen, float32(p.X)-sz/2, float32(p.Y)-sz/2, sz, sz, c, false)
	}
}

// spawnDeathParticles creates a pixel explosion when an enemy dies
func (g *Game) spawnDeathParticles(x, y float64, baseColor color.RGBA) {
	count := 12 + rand.Intn(8)
	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2
		speed := 1.0 + rand.Float64()*3.0
		life := 20 + rand.Intn(25)

		// Vary the color slightly
		c := baseColor
		if rand.Intn(3) == 0 {
			c = pixelPalette['a'] // yellow flash
		}

		g.particles = append(g.particles, Particle{
			X:       x + rand.Float64()*8 - 4,
			Y:       y + rand.Float64()*8 - 4,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle)*speed - 1.5,
			Life:    life,
			MaxLife: life,
			Color:   c,
			Size:    2 + rand.Float64()*3,
		})
	}
}

// spawnHitParticles creates small sparkle effect on attack hit
func (g *Game) spawnHitParticles(x, y float64) {
	count := 4 + rand.Intn(4)
	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2
		speed := 0.5 + rand.Float64()*2.0
		life := 10 + rand.Intn(10)

		g.particles = append(g.particles, Particle{
			X:       x + rand.Float64()*6 - 3,
			Y:       y + rand.Float64()*6 - 3,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle)*speed - 0.5,
			Life:    life,
			MaxLife: life,
			Color:   pixelPalette['7'], // white sparkle
			Size:    1 + rand.Float64()*2,
		})
	}
}

// spawnFireballExplosion creates a large fiery explosion
func (g *Game) spawnFireballExplosion(x, y float64) {
	colors := []color.RGBA{
		pixelPalette['8'], // red
		pixelPalette['9'], // orange
		pixelPalette['a'], // yellow
	}
	count := 20 + rand.Intn(10)
	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2
		speed := 1.0 + rand.Float64()*4.0
		life := 25 + rand.Intn(20)

		g.particles = append(g.particles, Particle{
			X:       x + rand.Float64()*10 - 5,
			Y:       y + rand.Float64()*10 - 5,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle)*speed - 2.0,
			Life:    life,
			MaxLife: life,
			Color:   colors[rand.Intn(len(colors))],
			Size:    3 + rand.Float64()*4,
		})
	}
}

// spawnTrailParticle creates a single trailing particle for projectiles
func (g *Game) spawnTrailParticle(x, y float64, isFireball bool) {
	var c color.RGBA
	var sz float64
	if isFireball {
		colors := []color.RGBA{pixelPalette['9'], pixelPalette['8']}
		c = colors[rand.Intn(len(colors))]
		sz = 2 + rand.Float64()*2
	} else {
		c = pixelPalette['a'] // yellow trail
		sz = 1 + rand.Float64()
	}

	g.particles = append(g.particles, Particle{
		X:       x + rand.Float64()*4 - 2,
		Y:       y + rand.Float64()*4 - 2,
		VX:      rand.Float64()*0.5 - 0.25,
		VY:      rand.Float64()*0.5 - 0.25,
		Life:    8 + rand.Intn(6),
		MaxLife: 14,
		Color:   c,
		Size:    sz,
	})
}

// spawnPlaceParticles creates particles when placing a unit
func (g *Game) spawnPlaceParticles(x, y float64) {
	for i := 0; i < 8; i++ {
		angle := float64(i) / 8.0 * math.Pi * 2
		speed := 1.5
		g.particles = append(g.particles, Particle{
			X:       x,
			Y:       y,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle) * speed,
			Life:    15,
			MaxLife: 15,
			Color:   pixelPalette['c'], // blue sparkle
			Size:    2,
		})
	}
}
