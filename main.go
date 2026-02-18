package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game implements ebiten.Game
type Game struct {
	State       GameState
	Tick        int
	Battle      *BattleState
	StageSelect *StageSelectState
}

func NewGame() *Game {
	InitFonts()
	return &Game{
		State: StateTitle,
	}
}

func (g *Game) Update() error {
	g.Tick++

	switch g.State {
	case StateTitle:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.StageSelect = NewStageSelectState()
			g.State = StateStageSelect
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return ebiten.Termination
		}

	case StateStageSelect:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State = StateTitle
			return nil
		}
		selected := g.StageSelect.Update()
		if selected >= 0 && selected < len(Stages) {
			g.Battle = NewBattleState(Stages[selected])
			g.State = StateBattle
		}

	case StateBattle:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			if g.Battle.GameOver {
				g.StageSelect = NewStageSelectState()
				g.State = StateStageSelect
			} else {
				// Confirm quit? For now, just go back
				g.StageSelect = NewStageSelectState()
				g.State = StateStageSelect
			}
			return nil
		}
		if g.Battle.GameOver && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.StageSelect = NewStageSelectState()
			g.State = StateStageSelect
			return nil
		}
		g.Battle.Update()

	case StateResult:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.StageSelect = NewStageSelectState()
			g.State = StateStageSelect
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.State {
	case StateTitle:
		DrawTitleScreen(screen, g.Tick)
	case StateStageSelect:
		g.StageSelect.Draw(screen)
	case StateBattle:
		g.Battle.Draw(screen)
	case StateResult:
		screen.Fill(color.RGBA{10, 10, 26, 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("NEON SIGIL: TRI-FUSE DEFENSE")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println(err)
	}
}
