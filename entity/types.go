package entity

// GameState represents the current screen/phase of the game.
type GameState int

const (
	StateTitle GameState = iota
	StateBattle
	StateReward
	StateGameOver
	StateVictory
)

// CardType identifies the kind of card/unit.
type CardType int

const (
	CardSoldier CardType = iota
	CardArcher
	CardSpearman
	CardMage
	CardFireball
)

// Race identifies the racial faction for synergy calculation.
type Race int

const (
	RaceNone Race = iota
	RaceHuman
	RaceElf
)

// EnemyType identifies enemy variants.
type EnemyType int

const (
	EnemyGoblin EnemyType = iota
	EnemyOrc
	EnemyBossOrc
	EnemyFinalBoss
)

// CardData holds the template stats for a card type.
type CardData struct {
	Name     string
	Type     CardType
	Race     Race
	Cost     int
	HP       int
	Atk      int
	AtkSpeed int // 쿨다운 틱 (낮을수록 빠름)
	Range    int // 0=근접, 2-3=원거리 (그리드 칸)
}

// WaveData defines the composition of enemies in a wave.
type WaveData struct {
	Goblins   int
	Orcs      int
	BossOrc   int // 체력 100 보스
	FinalBoss int // 체력 200 최종 보스
}

// Point is a 2D coordinate used for waypoints.
type Point struct {
	X, Y float64
}

// Card represents a card instance in the deck, hand, or graveyard.
type Card struct {
	Data CardData
}

// Summoner is a placed unit on the grid.
type Summoner struct {
	Card      CardData
	GridX     int
	GridY     int
	CurrentHP int
	MaxHP     int
	Atk       int
	AtkSpeed  int
	Range     int
	AtkTimer  int
	ScreenX   float64
	ScreenY   float64
}

// Enemy is a hostile unit traversing the path.
type Enemy struct {
	Type      EnemyType
	HP        int
	MaxHP     int
	Atk       int
	Speed     float64
	X, Y      float64
	PathIndex int
	PathT     float64 // 현재 세그먼트에서의 진행도 (0~1)
	Dead      bool
	Reached   bool // 소환사에게 도달했는지
	HasAura   bool // 최종 보스 오라
}

// Projectile is a ranged attack in flight.
type Projectile struct {
	X, Y       float64
	TargetX    float64
	TargetY    float64
	Damage     int
	Target     *Enemy
	Speed      float64
	IsFireball bool
	AOERange   float64
}
