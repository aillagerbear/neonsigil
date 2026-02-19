package game

import (
	"image/color"
)

// Screen dimensions
const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	TileSize     = 60
	BoardCols    = 8
	BoardRows    = 8
	BoardOffsetX = 40
	BoardOffsetY = 50
)

// Game states
type GameState int

const (
	StateTitle GameState = iota
	StateStageSelect
	StateBattle
	StateResult
)

// Battle phases
type BattlePhase int

const (
	PhasePrepare BattlePhase = iota
	PhaseWave
	PhaseWaveEnd
)

// Tile types
type TileType int

const (
	TileBuild TileType = iota
	TilePath
	TileBlock
	TileNode
	TileSpecial
)

// Special tile types
type SpecialType string

const (
	SpecialSeal      SpecialType = "SEAL"
	SpecialAntenna   SpecialType = "ANTENNA"
	SpecialWorkbench SpecialType = "WORKBENCH"
	SpecialGround    SpecialType = "GROUND"
)

// Enemy types
type EnemyType string

const (
	EnemyRunner   EnemyType = "RUNNER"
	EnemyBruiser  EnemyType = "BRUISER"
	EnemyShield   EnemyType = "SHIELD"
	EnemySplitter EnemyType = "SPLITTER"
	EnemyFlyer    EnemyType = "FLYER"
	EnemyStalker  EnemyType = "STALKER"
	EnemyHacker   EnemyType = "HACKER"
	EnemyCharger  EnemyType = "CHARGER"
	EnemyTotem    EnemyType = "TOTEM"
	EnemyBoss     EnemyType = "BOSS_GATE"
)

// Faction types
type Faction string

const (
	FactionStreet   Faction = "STREET"
	FactionCoven    Faction = "COVEN"
	FactionArcTech  Faction = "ARC_TECH"
	FactionExorcist Faction = "EXORCIST"
)

// Class types
type UnitClass string

const (
	ClassVanguard UnitClass = "VANGUARD"
	ClassMarksman UnitClass = "MARKSMAN"
	ClassCaster   UnitClass = "CASTER"
	ClassEngineer UnitClass = "ENGINEER"
	ClassSupport  UnitClass = "SUPPORT"
)

// Targeting types
type TargetMode string

const (
	TargetFrontmost TargetMode = "FRONTMOST"
	TargetLowHP     TargetMode = "LOW_HP"
	TargetNearest   TargetMode = "NEAREST"
)

// Attack types
type AttackType string

const (
	AttackMelee  AttackType = "MELEE"
	AttackRanged AttackType = "RANGED"
)

// Damage types
type DamageType string

const (
	DamagePhys  DamageType = "PHYS"
	DamageMagic DamageType = "MAGIC"
)

// EnemyDef defines enemy base stats
type EnemyDef struct {
	Type       EnemyType
	Name       string
	BaseHP     float64
	Speed      float64
	LeakDamage int
	Color      color.RGBA
	ShieldPct  float64 // ranged damage reduction (SHIELD type)
}

// UnitDef defines a unit template
type UnitDef struct {
	ID       string
	Name     string
	Cost     int
	Faction  Faction
	Class    UnitClass
	HP       float64
	ATK      float64
	AtkSpeed float64
	Range    int
	Armor    float64
	AtkType  AttackType
	DmgType  DamageType
	Targeting TargetMode
	SkillDesc string
}

// Pos is a grid coordinate
type Pos struct {
	X, Y int
}

// FPos is a float position for smooth movement
type FPos struct {
	X, Y float64
}

// SpecialTile defines a special tile on the map
type SpecialTileDef struct {
	Pos  Pos
	Type SpecialType
}

// WaveGroup defines a spawn group within a wave
type WaveGroup struct {
	Enemy    EnemyType
	Count    int
	Interval float64 // seconds between spawns
	PathID   string
}

// WaveDef defines a complete wave
type WaveDef struct {
	ID     string
	Groups []WaveGroup
}

// PathDef defines an enemy path
type PathDef struct {
	ID        string
	Waypoints []Pos
}

// ShopRules for a stage
type ShopRules struct {
	RerollEnabled  bool
	LevelUpEnabled bool
	AllowedCosts   []int
}

// StageDef defines a complete stage
type StageDef struct {
	ID           string
	Name         string
	Integrity    int
	StartingGold int
	StartingLv   int
	DeployCapBase int
	ShopRules    ShopRules
	TriFuseEnabled bool
	NodesEnabled   bool
	Blocks       []Pos
	Nodes        []Pos
	Specials     []SpecialTileDef
	Paths        []PathDef
	Waves        []WaveDef
	EnemyHPMul   float64
	EnemySpdMul  float64
	BarrierEffect string
}

// --- Runtime structs ---

// Enemy is a live enemy on the board
type Enemy struct {
	Def        *EnemyDef
	HP         float64
	MaxHP      float64
	Pos        FPos   // pixel position
	PathID     string
	WaypointIdx int
	Speed      float64
	Alive      bool
	Reached    bool  // reached the end
	SlowTimer  float64
	StunTimer  float64
	Visible    bool  // for STALKER
}

// Projectile represents an in-flight projectile
type Projectile struct {
	X, Y     float64
	TargetID int
	Damage   float64
	Speed    float64
	Alive    bool
}

// Unit is a live unit on the board or bench
type Unit struct {
	Def        *UnitDef
	Star       int
	GridX, GridY int  // -1 if on bench
	BenchSlot  int    // -1 if on board
	HP         float64
	MaxHP      float64
	ATK        float64
	AtkSpeed   float64
	Range      int
	AtkCooldown float64
	Deployed   bool
}

// Neon color palette
var (
	ColorBG         = color.RGBA{10, 10, 26, 255}
	ColorBGLight    = color.RGBA{18, 18, 42, 255}
	ColorNeonCyan   = color.RGBA{0, 255, 255, 255}
	ColorNeonMagenta = color.RGBA{255, 0, 255, 255}
	ColorNeonYellow = color.RGBA{255, 255, 0, 255}
	ColorNeonGreen  = color.RGBA{0, 255, 136, 255}
	ColorNeonRed    = color.RGBA{255, 51, 102, 255}
	ColorNeonBlue   = color.RGBA{0, 120, 255, 255}
	ColorNeonPurple = color.RGBA{180, 60, 255, 255}
	ColorNeonOrange = color.RGBA{255, 160, 0, 255}
	ColorGridLine   = color.RGBA{30, 30, 70, 255}
	ColorPathTile   = color.RGBA{15, 15, 35, 255}
	ColorBuildTile  = color.RGBA{22, 22, 50, 255}
	ColorBlockTile  = color.RGBA{35, 35, 55, 255}
	ColorNodeTile   = color.RGBA{10, 30, 70, 255}
	ColorWhite      = color.RGBA{255, 255, 255, 255}
	ColorWhiteDim   = color.RGBA{180, 180, 200, 255}
	ColorGold       = color.RGBA{255, 215, 0, 255}
	ColorHP         = color.RGBA{50, 255, 80, 255}
	ColorHPLow      = color.RGBA{255, 80, 50, 255}
	ColorShopBG     = color.RGBA{15, 15, 35, 240}
	ColorBtnBG      = color.RGBA{25, 25, 60, 255}
	ColorBtnHover   = color.RGBA{40, 40, 90, 255}
)

// Faction/class colors
var FactionColors = map[Faction]color.RGBA{
	FactionStreet:   {255, 120, 50, 255},
	FactionCoven:    {180, 60, 255, 255},
	FactionArcTech:  {0, 200, 255, 255},
	FactionExorcist: {255, 220, 100, 255},
}

var ClassColors = map[UnitClass]color.RGBA{
	ClassVanguard: {100, 180, 255, 255},
	ClassMarksman: {255, 100, 100, 255},
	ClassCaster:   {200, 100, 255, 255},
	ClassEngineer: {100, 255, 200, 255},
	ClassSupport:  {255, 255, 100, 255},
}
