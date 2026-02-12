package game

import (
	"image/color"

	"ebitengine-testing/config"
	"ebitengine-testing/entity"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game implements the ebiten.Game interface and holds all game state.
type Game struct {
	state   entity.GameState
	wave    int
	maxWave int

	summonerHP int
	mana       float64
	maxMana    int
	manaTimer  int

	deck      []entity.Card
	hand      []entity.Card
	graveyard []entity.Card

	grid        [config.GridRows][config.GridCols]*entity.Summoner
	summoners   []*entity.Summoner
	enemies     []*entity.Enemy
	projectiles []*entity.Projectile

	selectedCard int // -1 = 선택 안 됨, 0~4 = 핸드 인덱스
	gameSpeed    int // 1 또는 2

	// 웨이브 스포닝
	spawnQueue   []entity.EnemyType
	spawnTimer   int
	waveComplete bool
	allSpawned   bool

	// 보상 선택
	rewardCards []entity.Card
	rewardHover int

	// 파이어볼 타겟팅
	fireballMode bool

	// 시너지
	humanCount   int
	elfCount     int
	humanSynergy bool
	elfSynergy   bool

	// 틱 카운터
	ticks int
}

// New creates a new Game in the title state.
func New() *Game {
	return &Game{
		state:        entity.StateTitle,
		maxWave:      10,
		maxMana:      config.MaxMana,
		gameSpeed:    1,
		selectedCard: -1,
		rewardHover:  -1,
	}
}

// Update implements ebiten.Game.
func (g *Game) Update() error {
	switch g.state {
	case entity.StateTitle:
		g.updateTitle()
	case entity.StateBattle:
		for i := 0; i < g.gameSpeed; i++ {
			g.updateBattle()
		}
	case entity.StateReward:
		g.updateReward()
	case entity.StateGameOver, entity.StateVictory:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.state = entity.StateTitle
		}
	}
	return nil
}

// Draw implements ebiten.Game.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x1a, 0x1a, 0x2e, 0xff})

	switch g.state {
	case entity.StateTitle:
		g.drawTitle(screen)
	case entity.StateBattle:
		g.drawBattle(screen)
	case entity.StateReward:
		g.drawBattle(screen)
		g.drawReward(screen)
	case entity.StateGameOver:
		g.drawGameOver(screen)
	case entity.StateVictory:
		g.drawVictory(screen)
	}
}

// Layout implements ebiten.Game.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
