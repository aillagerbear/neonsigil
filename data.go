package main

import "image/color"

// ---- Enemy Definitions ----

var EnemyDefs = map[EnemyType]*EnemyDef{
	EnemyRunner: {
		Type: EnemyRunner, Name: "RUNNER", BaseHP: 80, Speed: 1.3,
		LeakDamage: 1, Color: color.RGBA{100, 255, 100, 255},
	},
	EnemyBruiser: {
		Type: EnemyBruiser, Name: "BRUISER", BaseHP: 220, Speed: 0.8,
		LeakDamage: 2, Color: color.RGBA{200, 100, 50, 255},
	},
	EnemyShield: {
		Type: EnemyShield, Name: "SHIELD", BaseHP: 180, Speed: 0.9,
		LeakDamage: 2, Color: color.RGBA{100, 150, 255, 255}, ShieldPct: 0.3,
	},
	EnemySplitter: {
		Type: EnemySplitter, Name: "SPLITTER", BaseHP: 140, Speed: 1.0,
		LeakDamage: 1, Color: color.RGBA{255, 200, 50, 255},
	},
	EnemyFlyer: {
		Type: EnemyFlyer, Name: "FLYER", BaseHP: 100, Speed: 1.2,
		LeakDamage: 1, Color: color.RGBA{200, 100, 255, 255},
	},
	EnemyStalker: {
		Type: EnemyStalker, Name: "STALKER", BaseHP: 110, Speed: 1.1,
		LeakDamage: 1, Color: color.RGBA{80, 80, 80, 255},
	},
	EnemyHacker: {
		Type: EnemyHacker, Name: "HACKER", BaseHP: 160, Speed: 0.95,
		LeakDamage: 2, Color: color.RGBA{0, 255, 200, 255},
	},
	EnemyCharger: {
		Type: EnemyCharger, Name: "CHARGER", BaseHP: 240, Speed: 1.15,
		LeakDamage: 3, Color: color.RGBA{255, 80, 80, 255},
	},
	EnemyTotem: {
		Type: EnemyTotem, Name: "TOTEM", BaseHP: 300, Speed: 0.7,
		LeakDamage: 3, Color: color.RGBA{255, 255, 100, 255},
	},
	EnemyBoss: {
		Type: EnemyBoss, Name: "GATEKEEPER", BaseHP: 2000, Speed: 0.5,
		LeakDamage: 99, Color: color.RGBA{255, 50, 50, 255},
	},
}

// ---- Unit Definitions ----

var UnitDefs = []*UnitDef{
	// 1-cost (8)
	{ID: "MOTH", Name: "MOTH", Cost: 1, Faction: FactionStreet, Class: ClassVanguard,
		HP: 520, ATK: 38, AtkSpeed: 1.0, Range: 1, Armor: 10,
		AtkType: AttackMelee, DmgType: DamagePhys, Targeting: TargetFrontmost,
		SkillDesc: "Taunt + DMG Reduction"},
	{ID: "VICE", Name: "VICE", Cost: 1, Faction: FactionStreet, Class: ClassMarksman,
		HP: 280, ATK: 55, AtkSpeed: 1.2, Range: 3, Armor: 3,
		AtkType: AttackRanged, DmgType: DamagePhys, Targeting: TargetLowHP,
		SkillDesc: "Backline Strike"},
	{ID: "KNOT", Name: "KNOT", Cost: 1, Faction: FactionCoven, Class: ClassCaster,
		HP: 300, ATK: 48, AtkSpeed: 0.9, Range: 3, Armor: 4,
		AtkType: AttackRanged, DmgType: DamageMagic, Targeting: TargetFrontmost,
		SkillDesc: "Curse Stack"},
	{ID: "TAR", Name: "TAR", Cost: 1, Faction: FactionCoven, Class: ClassSupport,
		HP: 350, ATK: 30, AtkSpeed: 0.8, Range: 2, Armor: 5,
		AtkType: AttackRanged, DmgType: DamageMagic, Targeting: TargetFrontmost,
		SkillDesc: "Slow Charm"},
	{ID: "SPARK", Name: "SPARK", Cost: 1, Faction: FactionArcTech, Class: ClassEngineer,
		HP: 300, ATK: 42, AtkSpeed: 1.0, Range: 2, Armor: 5,
		AtkType: AttackRanged, DmgType: DamagePhys, Targeting: TargetNearest,
		SkillDesc: "Deploy Drone"},
	{ID: "GLINT", Name: "GLINT", Cost: 1, Faction: FactionArcTech, Class: ClassMarksman,
		HP: 260, ATK: 58, AtkSpeed: 1.3, Range: 4, Armor: 2,
		AtkType: AttackRanged, DmgType: DamagePhys, Targeting: TargetFrontmost,
		SkillDesc: "Single Shot"},
	{ID: "HALO", Name: "HALO", Cost: 1, Faction: FactionExorcist, Class: ClassSupport,
		HP: 380, ATK: 25, AtkSpeed: 0.7, Range: 2, Armor: 6,
		AtkType: AttackRanged, DmgType: DamageMagic, Targeting: TargetNearest,
		SkillDesc: "Purify Shield"},
	{ID: "IRON", Name: "IRON", Cost: 1, Faction: FactionExorcist, Class: ClassVanguard,
		HP: 580, ATK: 32, AtkSpeed: 0.9, Range: 1, Armor: 14,
		AtkType: AttackMelee, DmgType: DamagePhys, Targeting: TargetFrontmost,
		SkillDesc: "DMG Reduction + Block"},

	// 2-cost (6)
	{ID: "GLASS", Name: "GLASS", Cost: 2, Faction: FactionStreet, Class: ClassMarksman,
		HP: 320, ATK: 72, AtkSpeed: 1.1, Range: 4, Armor: 3,
		AtkType: AttackRanged, DmgType: DamagePhys, Targeting: TargetLowHP,
		SkillDesc: "Mark Snipe"},
	{ID: "INK", Name: "INK", Cost: 2, Faction: FactionCoven, Class: ClassCaster,
		HP: 360, ATK: 60, AtkSpeed: 0.85, Range: 3, Armor: 5,
		AtkType: AttackRanged, DmgType: DamageMagic, Targeting: TargetFrontmost,
		SkillDesc: "Lingering Zone"},
	{ID: "PATCH", Name: "PATCH", Cost: 2, Faction: FactionArcTech, Class: ClassEngineer,
		HP: 400, ATK: 35, AtkSpeed: 0.8, Range: 2, Armor: 8,
		AtkType: AttackRanged, DmgType: DamagePhys, Targeting: TargetNearest,
		SkillDesc: "Repair Module"},
	{ID: "VOLT", Name: "VOLT", Cost: 2, Faction: FactionArcTech, Class: ClassCaster,
		HP: 340, ATK: 65, AtkSpeed: 0.9, Range: 3, Armor: 4,
		AtkType: AttackRanged, DmgType: DamageMagic, Targeting: TargetFrontmost,
		SkillDesc: "Shock AoE"},
	{ID: "LAMP", Name: "LAMP", Cost: 2, Faction: FactionExorcist, Class: ClassMarksman,
		HP: 300, ATK: 68, AtkSpeed: 1.2, Range: 4, Armor: 3,
		AtkType: AttackRanged, DmgType: DamageMagic, Targeting: TargetFrontmost,
		SkillDesc: "Undead Bane"},
	{ID: "LITANY", Name: "LITANY", Cost: 2, Faction: FactionExorcist, Class: ClassCaster,
		HP: 380, ATK: 58, AtkSpeed: 0.8, Range: 3, Armor: 6,
		AtkType: AttackRanged, DmgType: DamageMagic, Targeting: TargetFrontmost,
		SkillDesc: "Purify AoE"},

	// 3-cost (3)
	{ID: "COIN", Name: "COIN", Cost: 3, Faction: FactionStreet, Class: ClassSupport,
		HP: 420, ATK: 40, AtkSpeed: 0.7, Range: 2, Armor: 6,
		AtkType: AttackRanged, DmgType: DamagePhys, Targeting: TargetNearest,
		SkillDesc: "Kill Gold + Buff"},
	{ID: "DOLL", Name: "DOLL", Cost: 3, Faction: FactionCoven, Class: ClassCaster,
		HP: 450, ATK: 55, AtkSpeed: 0.75, Range: 3, Armor: 7,
		AtkType: AttackRanged, DmgType: DamageMagic, Targeting: TargetFrontmost,
		SkillDesc: "Summon Block"},
	{ID: "NODE", Name: "NODE", Cost: 3, Faction: FactionArcTech, Class: ClassSupport,
		HP: 480, ATK: 30, AtkSpeed: 0.6, Range: 2, Armor: 9,
		AtkType: AttackRanged, DmgType: DamagePhys, Targeting: TargetNearest,
		SkillDesc: "Deploy + CDR"},

	// 4-cost (1)
	{ID: "ORISON", Name: "ORISON", Cost: 4, Faction: FactionExorcist, Class: ClassCaster,
		HP: 500, ATK: 85, AtkSpeed: 0.6, Range: 3, Armor: 8,
		AtkType: AttackRanged, DmgType: DamageMagic, Targeting: TargetFrontmost,
		SkillDesc: "Grand Purify"},
}

// UnitDefByID provides quick lookup
var UnitDefByID = func() map[string]*UnitDef {
	m := make(map[string]*UnitDef)
	for _, u := range UnitDefs {
		m[u.ID] = u
	}
	return m
}()

// ---- Stage Definitions (CH1-01 ~ CH1-10) ----

var Stages = []*StageDef{
	{
		ID: "CH1-01", Name: "Tape on the Asphalt",
		Integrity: 15, StartingGold: 8, StartingLv: 1, DeployCapBase: 3,
		ShopRules: ShopRules{RerollEnabled: false, LevelUpEnabled: false, AllowedCosts: []int{1}},
		TriFuseEnabled: false, NodesEnabled: false,
		Blocks: []Pos{{4, 3}, {4, 4}},
		Paths: []PathDef{{ID: "P0", Waypoints: []Pos{
			{0, 3}, {1, 3}, {2, 3}, {2, 2}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 5}, {6, 5}, {7, 5},
		}}},
		Waves: []WaveDef{
			{ID: "W1", Groups: []WaveGroup{{EnemyRunner, 14, 0.60, "P0"}}},
			{ID: "W2", Groups: []WaveGroup{{EnemyRunner, 10, 0.55, "P0"}, {EnemyBruiser, 2, 1.20, "P0"}}},
			{ID: "W3", Groups: []WaveGroup{{EnemyRunner, 16, 0.55, "P0"}}},
			{ID: "W4", Groups: []WaveGroup{{EnemyBruiser, 4, 1.10, "P0"}}},
			{ID: "W5", Groups: []WaveGroup{{EnemyRunner, 10, 0.55, "P0"}, {EnemyBruiser, 3, 1.10, "P0"}}},
			{ID: "W6", Groups: []WaveGroup{{EnemyShield, 3, 1.00, "P0"}, {EnemyRunner, 10, 0.55, "P0"}}},
		},
		EnemyHPMul: 0.95, EnemySpdMul: 1.0,
	},
	{
		ID: "CH1-02", Name: "Cheap Tricks",
		Integrity: 15, StartingGold: 10, StartingLv: 1, DeployCapBase: 3,
		ShopRules: ShopRules{RerollEnabled: true, LevelUpEnabled: false, AllowedCosts: []int{1}},
		TriFuseEnabled: false, NodesEnabled: false,
		Blocks: []Pos{{3, 2}, {3, 3}, {3, 4}},
		Paths: []PathDef{{ID: "P0", Waypoints: []Pos{
			{0, 4}, {1, 4}, {2, 4}, {2, 3}, {2, 2}, {3, 2}, {4, 2}, {5, 2}, {5, 3}, {5, 4}, {6, 4}, {7, 4},
		}}},
		Waves: []WaveDef{
			{ID: "W1", Groups: []WaveGroup{{EnemyRunner, 16, 0.55, "P0"}}},
			{ID: "W2", Groups: []WaveGroup{{EnemyRunner, 10, 0.55, "P0"}, {EnemyBruiser, 3, 1.10, "P0"}}},
			{ID: "W3", Groups: []WaveGroup{{EnemyBruiser, 5, 1.10, "P0"}}},
			{ID: "W4", Groups: []WaveGroup{{EnemyRunner, 18, 0.50, "P0"}}},
			{ID: "W5", Groups: []WaveGroup{{EnemyShield, 4, 1.00, "P0"}, {EnemyRunner, 10, 0.55, "P0"}}},
			{ID: "W6", Groups: []WaveGroup{{EnemyBruiser, 4, 1.10, "P0"}, {EnemyShield, 4, 1.00, "P0"}}},
		},
		EnemyHPMul: 1.0, EnemySpdMul: 1.0,
	},
	{
		ID: "CH1-03", Name: "Breakpoint",
		Integrity: 15, StartingGold: 10, StartingLv: 1, DeployCapBase: 3,
		ShopRules: ShopRules{RerollEnabled: true, LevelUpEnabled: true, AllowedCosts: []int{1, 2}},
		TriFuseEnabled: false, NodesEnabled: false,
		Blocks: []Pos{{4, 2}, {4, 3}, {4, 4}, {4, 5}},
		Paths: []PathDef{{ID: "P0", Waypoints: []Pos{
			{0, 2}, {1, 2}, {2, 2}, {3, 2}, {3, 3}, {3, 4}, {3, 5}, {4, 5}, {5, 5}, {6, 5}, {7, 5},
		}}},
		Waves: []WaveDef{
			{ID: "W1", Groups: []WaveGroup{{EnemyRunner, 18, 0.55, "P0"}}},
			{ID: "W2", Groups: []WaveGroup{{EnemyBruiser, 4, 1.10, "P0"}, {EnemyRunner, 8, 0.55, "P0"}}},
			{ID: "W3", Groups: []WaveGroup{{EnemyShield, 6, 1.00, "P0"}}},
			{ID: "W4", Groups: []WaveGroup{{EnemyRunner, 22, 0.50, "P0"}}},
			{ID: "W5", Groups: []WaveGroup{{EnemyBruiser, 5, 1.10, "P0"}, {EnemyShield, 4, 1.00, "P0"}}},
			{ID: "W6", Groups: []WaveGroup{{EnemyShield, 6, 1.00, "P0"}, {EnemyRunner, 12, 0.55, "P0"}}},
		},
		EnemyHPMul: 1.05, EnemySpdMul: 1.0,
	},
	{
		ID: "CH1-04", Name: "TRI-FUSE",
		Integrity: 15, StartingGold: 10, StartingLv: 1, DeployCapBase: 3,
		ShopRules: ShopRules{RerollEnabled: true, LevelUpEnabled: true, AllowedCosts: []int{1, 2, 3}},
		TriFuseEnabled: true, NodesEnabled: false,
		Blocks: []Pos{{2, 4}, {3, 4}, {4, 4}, {5, 4}},
		Paths: []PathDef{{ID: "P0", Waypoints: []Pos{
			{0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, {5, 2}, {5, 1}, {6, 1}, {7, 1},
		}}},
		Waves: []WaveDef{
			{ID: "W1", Groups: []WaveGroup{{EnemyRunner, 16, 0.55, "P0"}}},
			{ID: "W2", Groups: []WaveGroup{{EnemySplitter, 6, 0.90, "P0"}}},
			{ID: "W3", Groups: []WaveGroup{{EnemyRunner, 12, 0.55, "P0"}, {EnemyShield, 4, 1.00, "P0"}}},
			{ID: "W4", Groups: []WaveGroup{{EnemySplitter, 8, 0.90, "P0"}, {EnemyRunner, 6, 0.55, "P0"}}},
			{ID: "W5", Groups: []WaveGroup{{EnemyBruiser, 4, 1.10, "P0"}, {EnemySplitter, 6, 0.90, "P0"}}},
			{ID: "W6", Groups: []WaveGroup{{EnemyShield, 6, 1.00, "P0"}, {EnemyRunner, 10, 0.55, "P0"}}},
		},
		EnemyHPMul: 1.05, EnemySpdMul: 1.0,
	},
	{
		ID: "CH1-05", Name: "Triangle of Salt",
		Integrity: 15, StartingGold: 10, StartingLv: 1, DeployCapBase: 3,
		ShopRules: ShopRules{RerollEnabled: true, LevelUpEnabled: true, AllowedCosts: []int{1, 2, 3}},
		TriFuseEnabled: true, NodesEnabled: true,
		Nodes:  []Pos{{2, 2}, {4, 5}, {6, 2}},
		Blocks: []Pos{{3, 3}, {4, 3}},
		Paths: []PathDef{{ID: "P0", Waypoints: []Pos{
			{0, 5}, {1, 5}, {2, 5}, {2, 4}, {2, 3}, {2, 2}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {7, 1},
		}}},
		Waves: []WaveDef{
			{ID: "W1", Groups: []WaveGroup{{EnemyRunner, 18, 0.55, "P0"}}},
			{ID: "W2", Groups: []WaveGroup{{EnemyBruiser, 4, 1.10, "P0"}}},
			{ID: "W3", Groups: []WaveGroup{{EnemySplitter, 6, 0.90, "P0"}, {EnemyRunner, 8, 0.55, "P0"}}},
			{ID: "W4", Groups: []WaveGroup{{EnemyShield, 6, 1.00, "P0"}}},
			{ID: "W5", Groups: []WaveGroup{{EnemyBruiser, 3, 1.10, "P0"}, {EnemyShield, 4, 1.00, "P0"}, {EnemyRunner, 8, 0.55, "P0"}}},
			{ID: "W6", Groups: []WaveGroup{{EnemySplitter, 10, 0.80, "P0"}}},
		},
		EnemyHPMul: 1.10, EnemySpdMul: 1.0, BarrierEffect: "BARRIER_SLOW",
	},
	{
		ID: "CH1-06", Name: "Seal the Doorway",
		Integrity: 18, StartingGold: 10, StartingLv: 1, DeployCapBase: 3,
		ShopRules: ShopRules{RerollEnabled: true, LevelUpEnabled: true, AllowedCosts: []int{1, 2, 3}},
		TriFuseEnabled: true, NodesEnabled: true,
		Nodes:    []Pos{{2, 5}, {5, 5}, {6, 2}},
		Specials: []SpecialTileDef{{Pos: Pos{3, 2}, Type: SpecialSeal}, {Pos: Pos{4, 2}, Type: SpecialSeal}},
		Blocks:   []Pos{{4, 4}},
		Paths: []PathDef{{ID: "P0", Waypoints: []Pos{
			{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {5, 2}, {5, 3}, {5, 4}, {5, 5}, {6, 5}, {7, 5},
		}}},
		Waves: []WaveDef{
			{ID: "W1", Groups: []WaveGroup{{EnemyRunner, 18, 0.55, "P0"}}},
			{ID: "W2", Groups: []WaveGroup{{EnemyShield, 5, 1.00, "P0"}}},
			{ID: "W3", Groups: []WaveGroup{{EnemySplitter, 6, 0.90, "P0"}, {EnemyRunner, 10, 0.55, "P0"}}},
			{ID: "W4", Groups: []WaveGroup{{EnemyBruiser, 5, 1.10, "P0"}}},
			{ID: "W5", Groups: []WaveGroup{{EnemyShield, 6, 1.00, "P0"}, {EnemyBruiser, 3, 1.10, "P0"}}},
			{ID: "W6", Groups: []WaveGroup{{EnemySplitter, 8, 0.90, "P0"}, {EnemyShield, 4, 1.00, "P0"}, {EnemyRunner, 8, 0.55, "P0"}}},
		},
		EnemyHPMul: 1.10, EnemySpdMul: 1.0, BarrierEffect: "BARRIER_SLOW",
	},
	{
		ID: "CH1-07", Name: "Air Above Neon",
		Integrity: 18, StartingGold: 12, StartingLv: 1, DeployCapBase: 3,
		ShopRules: ShopRules{RerollEnabled: true, LevelUpEnabled: true, AllowedCosts: []int{1, 2, 3}},
		TriFuseEnabled: true, NodesEnabled: true,
		Nodes:    []Pos{{2, 3}, {4, 5}, {6, 3}},
		Specials: []SpecialTileDef{{Pos: Pos{1, 1}, Type: SpecialAntenna}, {Pos: Pos{6, 6}, Type: SpecialAntenna}},
		Blocks:   []Pos{{3, 4}, {4, 4}},
		Paths: []PathDef{
			{ID: "P0", Waypoints: []Pos{{0, 3}, {1, 3}, {2, 3}, {3, 3}, {3, 2}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {7, 1}}},
			{ID: "P1", Waypoints: []Pos{{0, 6}, {2, 6}, {4, 6}, {6, 6}, {7, 5}}},
		},
		Waves: []WaveDef{
			{ID: "W1", Groups: []WaveGroup{{EnemyRunner, 16, 0.55, "P0"}}},
			{ID: "W2", Groups: []WaveGroup{{EnemyFlyer, 8, 0.80, "P1"}}},
			{ID: "W3", Groups: []WaveGroup{{EnemyShield, 4, 1.00, "P0"}, {EnemyRunner, 10, 0.55, "P0"}}},
			{ID: "W4", Groups: []WaveGroup{{EnemyFlyer, 10, 0.70, "P1"}, {EnemyRunner, 8, 0.55, "P0"}}},
			{ID: "W5", Groups: []WaveGroup{{EnemyBruiser, 4, 1.10, "P0"}, {EnemyFlyer, 8, 0.80, "P1"}}},
			{ID: "W6", Groups: []WaveGroup{{EnemyFlyer, 12, 0.65, "P1"}, {EnemyShield, 4, 1.00, "P0"}}},
		},
		EnemyHPMul: 1.12, EnemySpdMul: 1.0, BarrierEffect: "BARRIER_SLOW",
	},
	{
		ID: "CH1-08", Name: "The Unseen Lane",
		Integrity: 18, StartingGold: 12, StartingLv: 1, DeployCapBase: 3,
		ShopRules: ShopRules{RerollEnabled: true, LevelUpEnabled: true, AllowedCosts: []int{1, 2, 3}},
		TriFuseEnabled: true, NodesEnabled: true,
		Nodes:    []Pos{{1, 4}, {4, 1}, {6, 4}},
		Specials: []SpecialTileDef{{Pos: Pos{3, 6}, Type: SpecialAntenna}},
		Blocks:   []Pos{{2, 2}, {5, 5}},
		Paths: []PathDef{{ID: "P0", Waypoints: []Pos{
			{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 3}, {4, 2}, {4, 1}, {5, 1}, {6, 1}, {7, 1},
		}}},
		Waves: []WaveDef{
			{ID: "W1", Groups: []WaveGroup{{EnemyRunner, 18, 0.55, "P0"}}},
			{ID: "W2", Groups: []WaveGroup{{EnemyStalker, 8, 0.80, "P0"}}},
			{ID: "W3", Groups: []WaveGroup{{EnemySplitter, 6, 0.90, "P0"}, {EnemyRunner, 8, 0.55, "P0"}}},
			{ID: "W4", Groups: []WaveGroup{{EnemyStalker, 10, 0.70, "P0"}, {EnemyRunner, 10, 0.55, "P0"}}},
			{ID: "W5", Groups: []WaveGroup{{EnemyShield, 5, 1.00, "P0"}, {EnemyStalker, 6, 0.80, "P0"}}},
			{ID: "W6", Groups: []WaveGroup{{EnemyStalker, 12, 0.65, "P0"}, {EnemySplitter, 6, 0.90, "P0"}}},
		},
		EnemyHPMul: 1.13, EnemySpdMul: 1.0, BarrierEffect: "BARRIER_REVEAL",
	},
	{
		ID: "CH1-09", Name: "Static on the Line",
		Integrity: 18, StartingGold: 12, StartingLv: 1, DeployCapBase: 3,
		ShopRules: ShopRules{RerollEnabled: true, LevelUpEnabled: true, AllowedCosts: []int{1, 2, 3}},
		TriFuseEnabled: true, NodesEnabled: true,
		Nodes:    []Pos{{2, 2}, {4, 4}, {6, 2}},
		Specials: []SpecialTileDef{{Pos: Pos{2, 6}, Type: SpecialWorkbench}, {Pos: Pos{5, 6}, Type: SpecialWorkbench}},
		Blocks:   []Pos{{3, 3}},
		Paths: []PathDef{{ID: "P0", Waypoints: []Pos{
			{0, 6}, {1, 6}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {5, 5}, {5, 4}, {5, 3}, {6, 3}, {7, 3},
		}}},
		Waves: []WaveDef{
			{ID: "W1", Groups: []WaveGroup{{EnemyRunner, 18, 0.55, "P0"}}},
			{ID: "W2", Groups: []WaveGroup{{EnemyHacker, 5, 0.90, "P0"}, {EnemyRunner, 8, 0.55, "P0"}}},
			{ID: "W3", Groups: []WaveGroup{{EnemyShield, 5, 1.00, "P0"}, {EnemyHacker, 4, 0.90, "P0"}}},
			{ID: "W4", Groups: []WaveGroup{{EnemySplitter, 8, 0.90, "P0"}}},
			{ID: "W5", Groups: []WaveGroup{{EnemyBruiser, 4, 1.10, "P0"}, {EnemyHacker, 6, 0.90, "P0"}}},
			{ID: "W6", Groups: []WaveGroup{{EnemyHacker, 8, 0.80, "P0"}, {EnemyShield, 6, 1.00, "P0"}}},
		},
		EnemyHPMul: 1.15, EnemySpdMul: 1.0, BarrierEffect: "BARRIER_SLOW",
	},
	{
		ID: "CH1-10", Name: "Backroom Gatekeeper",
		Integrity: 18, StartingGold: 14, StartingLv: 1, DeployCapBase: 3,
		ShopRules: ShopRules{RerollEnabled: true, LevelUpEnabled: true, AllowedCosts: []int{1, 2, 3, 4}},
		TriFuseEnabled: true, NodesEnabled: true,
		Nodes:    []Pos{{2, 3}, {4, 6}, {6, 3}},
		Specials: []SpecialTileDef{{Pos: Pos{4, 2}, Type: SpecialSeal}, {Pos: Pos{4, 5}, Type: SpecialWorkbench}},
		Blocks:   []Pos{{3, 4}, {4, 4}, {5, 4}},
		Paths: []PathDef{{ID: "P0", Waypoints: []Pos{
			{0, 3}, {1, 3}, {2, 3}, {3, 3}, {3, 2}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {6, 2}, {6, 3}, {6, 4}, {6, 5}, {7, 5},
		}}},
		Waves: []WaveDef{
			{ID: "W1", Groups: []WaveGroup{{EnemyRunner, 20, 0.50, "P0"}}},
			{ID: "W2", Groups: []WaveGroup{{EnemyStalker, 8, 0.80, "P0"}, {EnemyRunner, 8, 0.55, "P0"}}},
			{ID: "W3", Groups: []WaveGroup{{EnemyHacker, 6, 0.90, "P0"}, {EnemyShield, 4, 1.00, "P0"}}},
			{ID: "W4", Groups: []WaveGroup{{EnemyBruiser, 5, 1.10, "P0"}, {EnemySplitter, 6, 0.90, "P0"}}},
			{ID: "W5", Groups: []WaveGroup{{EnemyRunner, 12, 0.55, "P0"}, {EnemyFlyer, 8, 0.80, "P0"}}},
			{ID: "W6", Groups: []WaveGroup{{EnemyBoss, 1, 0, "P0"}, {EnemyRunner, 12, 0.55, "P0"}}},
		},
		EnemyHPMul: 1.18, EnemySpdMul: 1.0, BarrierEffect: "BARRIER_MARK",
	},
}

// Shop pool weights by level
var ShopWeights = map[int][]float64{
	1: {1.0, 0, 0, 0},     // level 1: only 1-cost
	2: {0.75, 0.25, 0, 0}, // level 2
	3: {0.55, 0.30, 0.15, 0},
	4: {0.40, 0.30, 0.20, 0.10},
	5: {0.30, 0.30, 0.25, 0.15},
	6: {0.20, 0.25, 0.30, 0.25},
}

func GetUnitsForCost(cost int) []*UnitDef {
	var result []*UnitDef
	for _, u := range UnitDefs {
		if u.Cost == cost {
			result = append(result, u)
		}
	}
	return result
}
