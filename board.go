package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Board manages the 8x8 grid
type Board struct {
	Tiles    [BoardCols][BoardRows]TileType
	Specials map[Pos]SpecialType
	Paths    map[Pos]bool // which tiles are path tiles
	NodeSet  map[Pos]bool
	PathDefs []PathDef
}

func NewBoard(stage *StageDef) *Board {
	b := &Board{
		Specials: make(map[Pos]SpecialType),
		Paths:    make(map[Pos]bool),
		NodeSet:  make(map[Pos]bool),
		PathDefs: stage.Paths,
	}

	// Default all to BUILD
	for x := 0; x < BoardCols; x++ {
		for y := 0; y < BoardRows; y++ {
			b.Tiles[x][y] = TileBuild
		}
	}

	// Mark path tiles
	for _, p := range stage.Paths {
		for _, wp := range p.Waypoints {
			if wp.X >= 0 && wp.X < BoardCols && wp.Y >= 0 && wp.Y < BoardRows {
				b.Tiles[wp.X][wp.Y] = TilePath
				b.Paths[wp] = true
			}
		}
	}

	// Mark blocks
	for _, bl := range stage.Blocks {
		if bl.X >= 0 && bl.X < BoardCols && bl.Y >= 0 && bl.Y < BoardRows {
			b.Tiles[bl.X][bl.Y] = TileBlock
		}
	}

	// Mark nodes
	for _, nd := range stage.Nodes {
		if nd.X >= 0 && nd.X < BoardCols && nd.Y >= 0 && nd.Y < BoardRows {
			b.Tiles[nd.X][nd.Y] = TileNode
			b.NodeSet[nd] = true
		}
	}

	// Mark specials
	for _, sp := range stage.Specials {
		p := sp.Pos
		if p.X >= 0 && p.X < BoardCols && p.Y >= 0 && p.Y < BoardRows {
			b.Tiles[p.X][p.Y] = TileSpecial
			b.Specials[p] = sp.Type
		}
	}

	return b
}

func (b *Board) CanPlace(x, y int) bool {
	if x < 0 || x >= BoardCols || y < 0 || y >= BoardRows {
		return false
	}
	t := b.Tiles[x][y]
	return t == TileBuild || t == TileNode || t == TileSpecial
}

func (b *Board) TileScreenPos(x, y int) (float64, float64) {
	return float64(BoardOffsetX + x*TileSize), float64(BoardOffsetY + y*TileSize)
}

func (b *Board) ScreenToGrid(sx, sy int) (int, int) {
	gx := (sx - BoardOffsetX) / TileSize
	gy := (sy - BoardOffsetY) / TileSize
	return gx, gy
}

func (b *Board) Draw(screen *ebiten.Image, tick int) {
	b.DrawWithHighlight(screen, tick, false, nil)
}

func (b *Board) DrawWithHighlight(screen *ebiten.Image, tick int, highlightPlaceable bool, occupiedTiles map[Pos]bool) {
	// Draw tiles
	for x := 0; x < BoardCols; x++ {
		for y := 0; y < BoardRows; y++ {
			sx := float32(BoardOffsetX + x*TileSize)
			sy := float32(BoardOffsetY + y*TileSize)
			ts := float32(TileSize)

			var tileColor color.RGBA
			switch b.Tiles[x][y] {
			case TileBuild:
				tileColor = ColorBuildTile
			case TilePath:
				tileColor = ColorPathTile
			case TileBlock:
				tileColor = ColorBlockTile
			case TileNode:
				tileColor = ColorNodeTile
			case TileSpecial:
				tileColor = ColorBuildTile
			}

			// Fill tile
			vector.DrawFilledRect(screen, sx+1, sy+1, ts-2, ts-2, tileColor, false)

			// Highlight placeable tiles when selecting a unit
			if highlightPlaceable && b.CanPlace(x, y) {
				isOccupied := occupiedTiles != nil && occupiedTiles[Pos{x, y}]
				if !isOccupied {
					pulse := math.Sin(float64(tick%50)/50.0*math.Pi*2)*0.3 + 0.5
					alpha := uint8(float64(35) * pulse)
					vector.DrawFilledRect(screen, sx+1, sy+1, ts-2, ts-2, color.RGBA{0, 255, 200, alpha}, false)
				}
			}

			// Grid border
			vector.StrokeRect(screen, sx, sy, ts, ts, 1, ColorGridLine, false)
		}
	}

	// Draw path lines (neon glow)
	for _, pd := range b.PathDefs {
		b.drawPath(screen, pd, tick)
	}

	// Draw node glow
	for _, nd := range b.nodeList() {
		b.drawNodeGlow(screen, nd, tick)
	}

	// Draw special tile icons
	for pos, sp := range b.Specials {
		b.drawSpecialIcon(screen, pos, sp, tick)
	}
}

func (b *Board) drawPath(screen *ebiten.Image, pd PathDef, tick int) {
	if len(pd.Waypoints) < 2 {
		return
	}
	for i := 0; i < len(pd.Waypoints)-1; i++ {
		x1 := float32(BoardOffsetX+pd.Waypoints[i].X*TileSize) + float32(TileSize)/2
		y1 := float32(BoardOffsetY+pd.Waypoints[i].Y*TileSize) + float32(TileSize)/2
		x2 := float32(BoardOffsetX+pd.Waypoints[i+1].X*TileSize) + float32(TileSize)/2
		y2 := float32(BoardOffsetY+pd.Waypoints[i+1].Y*TileSize) + float32(TileSize)/2

		// Glow layer
		pulse := float64(tick%60) / 60.0
		alpha := uint8(40 + int(20*math.Sin(pulse*math.Pi*2)))
		glowColor := color.RGBA{0, 255, 255, alpha}
		vector.StrokeLine(screen, x1, y1, x2, y2, 6, glowColor, false)

		// Core line
		lineColor := color.RGBA{0, 200, 255, 120}
		vector.StrokeLine(screen, x1, y1, x2, y2, 2, lineColor, false)
	}

	// Draw direction dots moving along path
	dotPhase := float64(tick%120) / 120.0
	totalSegs := len(pd.Waypoints) - 1
	for d := 0; d < 3; d++ {
		phase := math.Mod(dotPhase+float64(d)*0.33, 1.0)
		segF := phase * float64(totalSegs)
		seg := int(segF)
		if seg >= totalSegs {
			seg = totalSegs - 1
		}
		t := segF - float64(seg)
		wx1 := float64(BoardOffsetX+pd.Waypoints[seg].X*TileSize) + float64(TileSize)/2
		wy1 := float64(BoardOffsetY+pd.Waypoints[seg].Y*TileSize) + float64(TileSize)/2
		wx2 := float64(BoardOffsetX+pd.Waypoints[seg+1].X*TileSize) + float64(TileSize)/2
		wy2 := float64(BoardOffsetY+pd.Waypoints[seg+1].Y*TileSize) + float64(TileSize)/2
		dx := wx1 + (wx2-wx1)*t
		dy := wy1 + (wy2-wy1)*t
		vector.DrawFilledCircle(screen, float32(dx), float32(dy), 3, color.RGBA{0, 255, 255, 180}, false)
	}
}

func (b *Board) drawNodeGlow(screen *ebiten.Image, nd Pos, tick int) {
	cx := float32(BoardOffsetX+nd.X*TileSize) + float32(TileSize)/2
	cy := float32(BoardOffsetY+nd.Y*TileSize) + float32(TileSize)/2

	pulse := math.Sin(float64(tick%90)/90.0*math.Pi*2)*0.3 + 0.7
	r := float32(TileSize/2 - 4)

	// Outer glow
	glowAlpha := uint8(float64(60) * pulse)
	vector.DrawFilledCircle(screen, cx, cy, r+4, color.RGBA{0, 120, 255, glowAlpha}, false)
	// Inner
	vector.StrokeCircle(screen, cx, cy, r, 2, color.RGBA{0, 180, 255, uint8(float64(200) * pulse)}, false)
	// Diamond shape
	s := float32(8)
	var path vector.Path
	path.MoveTo(cx, cy-s)
	path.LineTo(cx+s, cy)
	path.LineTo(cx, cy+s)
	path.LineTo(cx-s, cy)
	path.Close()
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].ColorR = 0
		vs[i].ColorG = 0.7
		vs[i].ColorB = 1
		vs[i].ColorA = float32(pulse) * 0.8
	}
	screen.DrawTriangles(vs, is, emptyImage, nil)
}

func (b *Board) drawSpecialIcon(screen *ebiten.Image, pos Pos, sp SpecialType, tick int) {
	cx := float32(BoardOffsetX+pos.X*TileSize) + float32(TileSize)/2
	cy := float32(BoardOffsetY+pos.Y*TileSize) + float32(TileSize)/2

	var c color.RGBA
	var label string
	switch sp {
	case SpecialSeal:
		c = color.RGBA{180, 60, 255, 200}
		label = "S"
	case SpecialAntenna:
		c = color.RGBA{0, 255, 136, 200}
		label = "A"
	case SpecialWorkbench:
		c = color.RGBA{255, 200, 50, 200}
		label = "W"
	case SpecialGround:
		c = color.RGBA{255, 120, 0, 200}
		label = "G"
	}

	// Border
	s := float32(TileSize/2 - 6)
	vector.StrokeRect(screen, cx-s, cy-s, s*2, s*2, 1.5, c, false)

	// Label
	_ = label
	// Small indicator dot
	vector.DrawFilledCircle(screen, cx, cy, 4, c, false)
}

func (b *Board) nodeList() []Pos {
	nodes := make([]Pos, 0, len(b.NodeSet))
	for p := range b.NodeSet {
		nodes = append(nodes, p)
	}
	return nodes
}

// emptyImage is used for DrawTriangles
var emptyImage *ebiten.Image

func init() {
	emptyImage = ebiten.NewImage(3, 3)
	emptyImage.Fill(color.White)
}
