package config

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

// Shop/bench constants
const (
	ShopSlots  = 5
	BenchSlots = 8
	RerollCost = 2
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

// Pos is a grid coordinate
type Pos struct {
	X, Y int
}

// FPos is a float position for smooth movement
type FPos struct {
	X, Y float64
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

// BenchSlotX returns the screen X coordinate for a bench slot
func BenchSlotX(slot int) int {
	return BoardOffsetX + slot*66
}

// BenchSlotY returns the screen Y coordinate for bench slots
func BenchSlotY() int {
	return BoardOffsetY + BoardRows*TileSize + 16
}

// WithAlpha returns a color with the given alpha value
func WithAlpha(c color.RGBA, a uint8) color.RGBA {
	return color.RGBA{c.R, c.G, c.B, a}
}
