package data

import (
	"image/color"

	"neonsigil/internal/config"
)

// EnemyDef defines enemy base stats
type EnemyDef struct {
	Type       config.EnemyType
	Name       string
	BaseHP     float64
	Speed      float64
	LeakDamage int
	Color      color.RGBA
	ShieldPct  float64 // ranged damage reduction (SHIELD type)
}

// UnitDef defines a unit template
type UnitDef struct {
	ID        string
	Name      string
	Cost      int
	Faction   config.Faction
	Class     config.UnitClass
	HP        float64
	ATK       float64
	AtkSpeed  float64
	Range     int
	Armor     float64
	AtkType   config.AttackType
	DmgType   config.DamageType
	Targeting config.TargetMode
	SkillDesc string
}

// SpecialTileDef defines a special tile on the map
type SpecialTileDef struct {
	Pos  config.Pos
	Type config.SpecialType
}

// WaveGroup defines a spawn group within a wave
type WaveGroup struct {
	Enemy    config.EnemyType
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
	Waypoints []config.Pos
}

// ShopRules for a stage
type ShopRules struct {
	RerollEnabled  bool
	LevelUpEnabled bool
	AllowedCosts   []int
}

// StageDef defines a complete stage
type StageDef struct {
	ID             string
	Name           string
	Integrity      int
	StartingGold   int
	StartingLv     int
	DeployCapBase  int
	ShopRules      ShopRules
	TriFuseEnabled bool
	NodesEnabled   bool
	Blocks         []config.Pos
	Nodes          []config.Pos
	Specials       []SpecialTileDef
	Paths          []PathDef
	Waves          []WaveDef
	EnemyHPMul     float64
	EnemySpdMul    float64
	BarrierEffect  string
}
