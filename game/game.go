package game

import (
	"ebitengine-testing/config"
	"ebitengine-testing/entity"

	"github.com/hajimehoshi/ebiten/v2"
)

// CardRect represents the screen area of a drawn card for click detection
type CardRect struct {
	X, Y, W, H float64
	Index       int
}

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

	// Visual effects
	particles []Particle
	sprites   *SpriteCache
	bgStars   []Star
	animTick  int

	// Custom card UI
	hoverCard       int        // -1 = none
	cardRects       []CardRect // current frame card positions
	rewardCardRects []CardRect // reward card positions
	hoverReward     int        // -1 = none

	// ebitenui
	ui             *UIManager
	endScreenReady bool
}

// New creates a new Game in the title state.
func New() *Game {
	initFonts()

	g := &Game{
		state:        entity.StateTitle,
		maxWave:      10,
		maxMana:      config.MaxMana,
		gameSpeed:    1,
		selectedCard: -1,
		rewardHover:  -1,
		hoverCard:    -1,
		hoverReward:  -1,
	}
	g.sprites = initSprites()
	g.bgStars = initStars(80)
	g.ui = newUIManager(g)
	return g
}

// Update implements ebiten.Game.
func (g *Game) Update() error {
	g.animTick++
	g.updateParticles()

	switch g.state {
	case entity.StateTitle:
		g.ui.titleUI.Update()
	case entity.StateBattle:
		g.updateCardHover()
		g.ui.updateBattleHUD(g)
		g.ui.battleUI.Update()
		for i := 0; i < g.gameSpeed; i++ {
			g.updateBattle()
		}
	case entity.StateReward:
		g.updateRewardHover()
		g.handleRewardInput()
	case entity.StateGameOver, entity.StateVictory:
		if !g.endScreenReady {
			g.ui.updateEndScreen(g)
			g.endScreenReady = true
		}
		g.ui.endUI.Update()
	}
	return nil
}

// Draw implements ebiten.Game.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colorBG)

	switch g.state {
	case entity.StateTitle:
		drawStars(screen, g.bgStars, g.animTick)
		g.drawTitleDecorations(screen)
		g.ui.titleUI.Draw(screen)
	case entity.StateBattle:
		g.drawBattle(screen)
		g.drawCustomHUD(screen)
		g.ui.battleUI.Draw(screen)
		g.drawHandCards(screen)
		if g.fireballMode {
			g.drawFireballIndicator(screen)
		}
	case entity.StateReward:
		g.drawBattle(screen)
		g.drawRewardScreen(screen)
	case entity.StateGameOver, entity.StateVictory:
		drawStars(screen, g.bgStars, g.animTick)
		g.drawEndScreenEffects(screen)
		g.ui.endUI.Draw(screen)
	}
}

// Layout implements ebiten.Game.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}

// handleCardSelect is called by ebitenui when a card button is clicked
func (g *Game) handleCardSelect(index int) {
	if index < 0 || index >= len(g.hand) {
		return
	}

	if g.selectedCard == index {
		g.selectedCard = -1
		return
	}

	g.selectedCard = index
	if g.hand[index].Data.Type == entity.CardFireball {
		if g.mana >= float64(g.hand[index].Data.Cost) {
			g.mana -= float64(g.hand[index].Data.Cost)
			g.graveyard = append(g.graveyard, g.hand[index])
			g.hand = append(g.hand[:index], g.hand[index+1:]...)
			g.drawCard()
			g.fireballMode = true
			g.selectedCard = -1
		} else {
			g.selectedCard = -1
		}
	}
}

// handleRewardSelect is called by ebitenui when a reward card is clicked
func (g *Game) handleRewardSelect(index int) {
	if index < 0 || index >= len(g.rewardCards) {
		return
	}

	g.deck = append(g.deck, g.rewardCards[index])
	g.shuffleDeck()

	for len(g.hand) < config.MaxHandSize {
		if len(g.deck) == 0 && len(g.graveyard) == 0 {
			break
		}
		g.drawCard()
	}

	g.wave++
	g.state = entity.StateBattle
	g.startWave()
}
