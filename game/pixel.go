package game

import (
	"image/color"
	"math"
	"math/rand"

	"ebitengine-testing/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PICO-8 inspired color palette
var pixelPalette = map[byte]color.RGBA{
	'0': {0x00, 0x00, 0x00, 0xFF}, // black
	'1': {0x1D, 0x2B, 0x53, 0xFF}, // dark blue
	'2': {0x7E, 0x25, 0x53, 0xFF}, // dark purple
	'3': {0x00, 0x87, 0x51, 0xFF}, // dark green
	'4': {0xAB, 0x52, 0x36, 0xFF}, // brown
	'5': {0x5F, 0x57, 0x4F, 0xFF}, // dark grey
	'6': {0xC2, 0xC3, 0xC7, 0xFF}, // light grey
	'7': {0xFF, 0xF1, 0xE8, 0xFF}, // white
	'8': {0xFF, 0x00, 0x4D, 0xFF}, // red
	'9': {0xFF, 0xA3, 0x00, 0xFF}, // orange
	'a': {0xFF, 0xEC, 0x27, 0xFF}, // yellow
	'b': {0x00, 0xE4, 0x36, 0xFF}, // green
	'c': {0x29, 0xAD, 0xFF, 0xFF}, // blue
	'd': {0x83, 0x76, 0x9C, 0xFF}, // indigo/lavender
	'e': {0xFF, 0x77, 0xA8, 0xFF}, // pink
	'f': {0xFF, 0xCC, 0xAA, 0xFF}, // peach/skin
}

// Theme colors
var (
	colorBG       = color.RGBA{0x10, 0x10, 0x1a, 0xFF}
	colorUIBG     = color.RGBA{0x18, 0x18, 0x28, 0xEE}
	colorUIBorder = color.RGBA{0x3a, 0x3a, 0x5a, 0xFF}
)

// ---- Sprite data ----

// Soldier: Human Knight (8 wide x 10 tall, blue armor)
var soldierSpriteData = []string{
	"..0000..",
	".0c6c60.",
	".06cc60.",
	"..0ff0..",
	".0cccc0.",
	"0c6cc6c0",
	".0cccc0.",
	"..0cc0..",
	"..0550..",
	"..00.00.",
}

// Archer: Human Ranger (8x10, green with brown bow)
var archerSpriteData = []string{
	"..0000..",
	".0b3b0..",
	".03bb0..",
	"..0ff0..",
	".03bb30.",
	".03bb340",
	"03bbb340",
	"..03b0..",
	"..0440..",
	"..00.00.",
}

// Spearman: Elf Warrior (8x10, orange/gold with spear tip)
var spearmanSpriteData = []string{
	"...0a0..",
	"..0000..",
	".099f90.",
	".09f90..",
	"0.0ff0.0",
	".09990..",
	"099999a0",
	".09990..",
	"..0440..",
	"..00.00.",
}

// Mage: Elf Wizard (8x10, purple robe with hat)
var mageSpriteData = []string{
	"...0a0..",
	"..0220..",
	".022220.",
	"..02d0..",
	"..0ff0..",
	".02dd20.",
	"02dddd20",
	".02dd20.",
	"..0220..",
	"..00.00.",
}

// Goblin: small green enemy (6x8)
var goblinSpriteData = []string{
	"0b..b0",
	"0bbbb0",
	"08bb80",
	".0bb0.",
	".0bb0.",
	"0b33b0",
	".0bb0.",
	".0..0.",
}

// Orc: bulky brown enemy (8x10)
var orcSpriteData = []string{
	".000000.",
	"04555540",
	"08545480",
	"04555540",
	".04540..",
	".04540..",
	"04545450",
	".04540..",
	".04.40..",
	".00.00..",
}

// Boss Orc: larger enemy with crown (10x12)
var bossOrcSpriteData = []string{
	".0a0.0a0..",
	"..000000..",
	".04555540.",
	".08545480.",
	".04555540.",
	"..045540..",
	"..045540..",
	".0454540..",
	"045454540.",
	".0454540..",
	"..04.40...",
	"..00.00...",
}

// Final Boss: large imposing enemy (12x14)
var finalBossSpriteData = []string{
	"0a0.....0a0.",
	".000000000..",
	"025555555200",
	"085555555800",
	"025555555200",
	".0255555520.",
	".0255555520.",
	"025555555520",
	".0255555520.",
	".0255555520.",
	"..02555520..",
	"...025520...",
	"...02..20...",
	"...00..00...",
}

// Base crystal (8x10)
var baseCrystalSpriteData = []string{
	"...0a0..",
	"..0aca0.",
	".0accca0",
	"0acccca0",
	"0acccca0",
	".0accca0",
	"..0aca0.",
	"...0a0..",
	"..0ccc0.",
	"..00000.",
}

// Fireball projectile (5x5)
var fireballSpriteData = []string{
	".090.",
	"09890",
	"98a89",
	"09890",
	".090.",
}

// Arrow projectile (3x5)
var arrowSpriteData = []string{
	".a.",
	".a.",
	"0a0",
	".0.",
	".0.",
}

// ---- SpriteCache ----

type SpriteCache struct {
	Soldier   *ebiten.Image
	Archer    *ebiten.Image
	Spearman  *ebiten.Image
	Mage      *ebiten.Image
	Goblin    *ebiten.Image
	Orc       *ebiten.Image
	BossOrc   *ebiten.Image
	FinalBoss *ebiten.Image
	Base      *ebiten.Image
	Fireball  *ebiten.Image
	Arrow     *ebiten.Image

	GrassTile    *ebiten.Image
	OccupiedTile *ebiten.Image
	HoverTile    *ebiten.Image
	BlockedTile  *ebiten.Image
}

// ---- Star background ----

type Star struct {
	X, Y    float64
	Bright  byte
	Twinkle float64
}

// ---- Bayer 4x4 dithering matrix ----

var bayerMatrix = [4][4]float64{
	{0.0 / 16, 8.0 / 16, 2.0 / 16, 10.0 / 16},
	{12.0 / 16, 4.0 / 16, 14.0 / 16, 6.0 / 16},
	{3.0 / 16, 11.0 / 16, 1.0 / 16, 9.0 / 16},
	{15.0 / 16, 7.0 / 16, 13.0 / 16, 5.0 / 16},
}

// ---- Sprite creation ----

func createSprite(data []string) *ebiten.Image {
	if len(data) == 0 {
		return ebiten.NewImage(1, 1)
	}
	h := len(data)
	w := len(data[0])
	img := ebiten.NewImage(w, h)
	pix := make([]byte, w*h*4)

	for y, row := range data {
		for x := 0; x < len(row) && x < w; x++ {
			ch := row[x]
			if ch == '.' {
				continue
			}
			if c, ok := pixelPalette[ch]; ok {
				offset := (y*w + x) * 4
				pix[offset] = c.R
				pix[offset+1] = c.G
				pix[offset+2] = c.B
				pix[offset+3] = c.A
			}
		}
	}

	img.WritePixels(pix)
	return img
}

func createDitheredTile(size int, c1, c2 color.RGBA, density float64, borderColor color.RGBA) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	pix := make([]byte, size*size*4)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			var c color.RGBA

			// 1px border
			if x == 0 || y == 0 || x == size-1 || y == size-1 {
				c = borderColor
			} else {
				threshold := bayerMatrix[y%4][x%4]
				if threshold < density {
					c = c1
				} else {
					c = c2
				}
				// Add some random grass-like detail
				if (x*7+y*13)%23 == 0 {
					c = lerpColor(c, c1, 0.3)
				}
			}

			offset := (y*size + x) * 4
			pix[offset] = c.R
			pix[offset+1] = c.G
			pix[offset+2] = c.B
			pix[offset+3] = c.A
		}
	}

	img.WritePixels(pix)
	return img
}

func lerpColor(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: byte(float64(a.R)*(1-t) + float64(b.R)*t),
		G: byte(float64(a.G)*(1-t) + float64(b.G)*t),
		B: byte(float64(a.B)*(1-t) + float64(b.B)*t),
		A: byte(float64(a.A)*(1-t) + float64(b.A)*t),
	}
}

func initSprites() *SpriteCache {
	cache := &SpriteCache{}

	// Unit sprites
	cache.Soldier = createSprite(soldierSpriteData)
	cache.Archer = createSprite(archerSpriteData)
	cache.Spearman = createSprite(spearmanSpriteData)
	cache.Mage = createSprite(mageSpriteData)

	// Enemy sprites
	cache.Goblin = createSprite(goblinSpriteData)
	cache.Orc = createSprite(orcSpriteData)
	cache.BossOrc = createSprite(bossOrcSpriteData)
	cache.FinalBoss = createSprite(finalBossSpriteData)

	// Other sprites
	cache.Base = createSprite(baseCrystalSpriteData)
	cache.Fireball = createSprite(fireballSpriteData)
	cache.Arrow = createSprite(arrowSpriteData)

	// Tile textures (dithered algorithmic patterns)
	grassC1 := color.RGBA{0x1a, 0x4a, 0x2a, 0xFF}
	grassC2 := color.RGBA{0x25, 0x5a, 0x35, 0xFF}
	grassBorder := color.RGBA{0x30, 0x6a, 0x40, 0x99}
	cache.GrassTile = createDitheredTile(config.TileSize, grassC1, grassC2, 0.45, grassBorder)

	occC1 := color.RGBA{0x3a, 0x35, 0x20, 0xFF}
	occC2 := color.RGBA{0x45, 0x40, 0x2a, 0xFF}
	occBorder := color.RGBA{0x55, 0x50, 0x35, 0xBB}
	cache.OccupiedTile = createDitheredTile(config.TileSize, occC1, occC2, 0.40, occBorder)

	hoverC1 := color.RGBA{0x25, 0x60, 0x35, 0xFF}
	hoverC2 := color.RGBA{0x30, 0x75, 0x45, 0xFF}
	hoverBorder := color.RGBA{0x50, 0xAA, 0x60, 0xCC}
	cache.HoverTile = createDitheredTile(config.TileSize, hoverC1, hoverC2, 0.50, hoverBorder)

	blockC1 := color.RGBA{0x55, 0x25, 0x25, 0xFF}
	blockC2 := color.RGBA{0x65, 0x30, 0x30, 0xFF}
	blockBorder := color.RGBA{0x80, 0x40, 0x40, 0xCC}
	cache.BlockedTile = createDitheredTile(config.TileSize, blockC1, blockC2, 0.45, blockBorder)

	return cache
}

// ---- Drawing utilities ----

func drawSpriteAt(screen *ebiten.Image, sprite *ebiten.Image, x, y, scale float64) {
	if sprite == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	w := float64(sprite.Bounds().Dx())
	h := float64(sprite.Bounds().Dy())
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)
	screen.DrawImage(sprite, op)
}

func drawSpriteAtWithColor(screen *ebiten.Image, sprite *ebiten.Image, x, y, scale float64, r, g, b, a float32) {
	if sprite == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	w := float64(sprite.Bounds().Dx())
	h := float64(sprite.Bounds().Dy())
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)
	op.ColorScale.Scale(r, g, b, a)
	screen.DrawImage(sprite, op)
}


// drawPixelBar draws a segmented pixel-art health/mana bar
func drawPixelBar(screen *ebiten.Image, x, y, w, h float32, ratio float64, fillColor, bgColor color.RGBA) {
	// Background
	vector.FillRect(screen, x, y, w, h, bgColor, false)
	// Fill
	fillW := float32(ratio) * w
	if fillW > 0 {
		vector.FillRect(screen, x, y, fillW, h, fillColor, false)
	}
	// 1px border
	vector.StrokeRect(screen, x, y, w, h, 1, color.RGBA{0x00, 0x00, 0x00, 0xAA}, false)
	// Pixel highlight on top edge
	if h > 2 {
		highlightColor := lerpColor(fillColor, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, 0.3)
		vector.FillRect(screen, x+1, y+1, fillW-2, 1, highlightColor, false)
	}
}

// ---- Stars ----

func initStars(count int) []Star {
	stars := make([]Star, count)
	for i := range stars {
		stars[i] = Star{
			X:       rand.Float64() * config.ScreenWidth,
			Y:       rand.Float64() * config.ScreenHeight,
			Bright:  byte(100 + rand.Intn(156)),
			Twinkle: 0.02 + rand.Float64()*0.05,
		}
	}
	return stars
}

func drawStars(screen *ebiten.Image, stars []Star, tick int) {
	for _, s := range stars {
		alpha := float64(s.Bright) * (0.3 + 0.7*math.Abs(math.Sin(float64(tick)*s.Twinkle)))
		a := byte(math.Min(alpha, 255))
		if a < 30 {
			continue
		}
		clr := color.RGBA{s.Bright, s.Bright, s.Bright, a}
		vector.FillRect(screen, float32(s.X), float32(s.Y), 1, 1, clr, false)
	}
}

// ---- Path energy dots (algorithmic animation) ----

func drawPathEnergyDots(screen *ebiten.Image, tick int) {
	dotColor := color.RGBA{0x40, 0x60, 0x80, 0xAA}

	numDots := 12
	for i := 0; i < numDots; i++ {
		t := float64(tick)/120.0 + float64(i)/float64(numDots)
		t = t - math.Floor(t) // normalize to 0-1

		// Interpolate along the full path
		totalSegments := len(entity_path_points) - 1
		if totalSegments <= 0 {
			return
		}
		segF := t * float64(totalSegments)
		seg := int(segF)
		if seg >= totalSegments {
			seg = totalSegments - 1
		}
		frac := segF - float64(seg)

		x := entity_path_points[seg][0]*(1-frac) + entity_path_points[seg+1][0]*frac
		y := entity_path_points[seg][1]*(1-frac) + entity_path_points[seg+1][1]*frac

		// Pulsing size
		pulse := 1.0 + 0.5*math.Sin(float64(tick)*0.1+float64(i)*0.8)
		sz := float32(2 * pulse)
		vector.FillRect(screen, float32(x)-sz/2, float32(y)-sz/2, sz, sz, dotColor, false)
	}
}

// Pre-computed path points for energy dots (from entity.EnemyPath)
var entity_path_points = [][2]float64{
	{520, -20}, {520, 80}, {750, 80}, {750, 230},
	{280, 230}, {280, 380}, {750, 380}, {750, 530},
	{520, 530}, {520, 700},
}
