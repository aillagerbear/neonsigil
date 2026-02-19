package data

import "neonsigil/internal/config"

// UnitDefs contains all unit definitions
var UnitDefs = []*UnitDef{
	// 1-cost (8)
	{ID: "MOTH", Name: "MOTH", Cost: 1, Faction: config.FactionStreet, Class: config.ClassVanguard,
		HP: 520, ATK: 38, AtkSpeed: 1.0, Range: 1, Armor: 10,
		AtkType: config.AttackMelee, DmgType: config.DamagePhys, Targeting: config.TargetFrontmost,
		SkillDesc: "Taunt + DMG Reduction"},
	{ID: "VICE", Name: "VICE", Cost: 1, Faction: config.FactionStreet, Class: config.ClassMarksman,
		HP: 280, ATK: 55, AtkSpeed: 1.2, Range: 3, Armor: 3,
		AtkType: config.AttackRanged, DmgType: config.DamagePhys, Targeting: config.TargetLowHP,
		SkillDesc: "Backline Strike"},
	{ID: "KNOT", Name: "KNOT", Cost: 1, Faction: config.FactionCoven, Class: config.ClassCaster,
		HP: 300, ATK: 48, AtkSpeed: 0.9, Range: 3, Armor: 4,
		AtkType: config.AttackRanged, DmgType: config.DamageMagic, Targeting: config.TargetFrontmost,
		SkillDesc: "Curse Stack"},
	{ID: "TAR", Name: "TAR", Cost: 1, Faction: config.FactionCoven, Class: config.ClassSupport,
		HP: 350, ATK: 30, AtkSpeed: 0.8, Range: 2, Armor: 5,
		AtkType: config.AttackRanged, DmgType: config.DamageMagic, Targeting: config.TargetFrontmost,
		SkillDesc: "Slow Charm"},
	{ID: "SPARK", Name: "SPARK", Cost: 1, Faction: config.FactionArcTech, Class: config.ClassEngineer,
		HP: 300, ATK: 42, AtkSpeed: 1.0, Range: 2, Armor: 5,
		AtkType: config.AttackRanged, DmgType: config.DamagePhys, Targeting: config.TargetNearest,
		SkillDesc: "Deploy Drone"},
	{ID: "GLINT", Name: "GLINT", Cost: 1, Faction: config.FactionArcTech, Class: config.ClassMarksman,
		HP: 260, ATK: 58, AtkSpeed: 1.3, Range: 4, Armor: 2,
		AtkType: config.AttackRanged, DmgType: config.DamagePhys, Targeting: config.TargetFrontmost,
		SkillDesc: "Single Shot"},
	{ID: "HALO", Name: "HALO", Cost: 1, Faction: config.FactionExorcist, Class: config.ClassSupport,
		HP: 380, ATK: 25, AtkSpeed: 0.7, Range: 2, Armor: 6,
		AtkType: config.AttackRanged, DmgType: config.DamageMagic, Targeting: config.TargetNearest,
		SkillDesc: "Purify Shield"},
	{ID: "IRON", Name: "IRON", Cost: 1, Faction: config.FactionExorcist, Class: config.ClassVanguard,
		HP: 580, ATK: 32, AtkSpeed: 0.9, Range: 1, Armor: 14,
		AtkType: config.AttackMelee, DmgType: config.DamagePhys, Targeting: config.TargetFrontmost,
		SkillDesc: "DMG Reduction + Block"},

	// 2-cost (6)
	{ID: "GLASS", Name: "GLASS", Cost: 2, Faction: config.FactionStreet, Class: config.ClassMarksman,
		HP: 320, ATK: 72, AtkSpeed: 1.1, Range: 4, Armor: 3,
		AtkType: config.AttackRanged, DmgType: config.DamagePhys, Targeting: config.TargetLowHP,
		SkillDesc: "Mark Snipe"},
	{ID: "INK", Name: "INK", Cost: 2, Faction: config.FactionCoven, Class: config.ClassCaster,
		HP: 360, ATK: 60, AtkSpeed: 0.85, Range: 3, Armor: 5,
		AtkType: config.AttackRanged, DmgType: config.DamageMagic, Targeting: config.TargetFrontmost,
		SkillDesc: "Lingering Zone"},
	{ID: "PATCH", Name: "PATCH", Cost: 2, Faction: config.FactionArcTech, Class: config.ClassEngineer,
		HP: 400, ATK: 35, AtkSpeed: 0.8, Range: 2, Armor: 8,
		AtkType: config.AttackRanged, DmgType: config.DamagePhys, Targeting: config.TargetNearest,
		SkillDesc: "Repair Module"},
	{ID: "VOLT", Name: "VOLT", Cost: 2, Faction: config.FactionArcTech, Class: config.ClassCaster,
		HP: 340, ATK: 65, AtkSpeed: 0.9, Range: 3, Armor: 4,
		AtkType: config.AttackRanged, DmgType: config.DamageMagic, Targeting: config.TargetFrontmost,
		SkillDesc: "Shock AoE"},
	{ID: "LAMP", Name: "LAMP", Cost: 2, Faction: config.FactionExorcist, Class: config.ClassMarksman,
		HP: 300, ATK: 68, AtkSpeed: 1.2, Range: 4, Armor: 3,
		AtkType: config.AttackRanged, DmgType: config.DamageMagic, Targeting: config.TargetFrontmost,
		SkillDesc: "Undead Bane"},
	{ID: "LITANY", Name: "LITANY", Cost: 2, Faction: config.FactionExorcist, Class: config.ClassCaster,
		HP: 380, ATK: 58, AtkSpeed: 0.8, Range: 3, Armor: 6,
		AtkType: config.AttackRanged, DmgType: config.DamageMagic, Targeting: config.TargetFrontmost,
		SkillDesc: "Purify AoE"},

	// 3-cost (3)
	{ID: "COIN", Name: "COIN", Cost: 3, Faction: config.FactionStreet, Class: config.ClassSupport,
		HP: 420, ATK: 40, AtkSpeed: 0.7, Range: 2, Armor: 6,
		AtkType: config.AttackRanged, DmgType: config.DamagePhys, Targeting: config.TargetNearest,
		SkillDesc: "Kill Gold + Buff"},
	{ID: "DOLL", Name: "DOLL", Cost: 3, Faction: config.FactionCoven, Class: config.ClassCaster,
		HP: 450, ATK: 55, AtkSpeed: 0.75, Range: 3, Armor: 7,
		AtkType: config.AttackRanged, DmgType: config.DamageMagic, Targeting: config.TargetFrontmost,
		SkillDesc: "Summon Block"},
	{ID: "NODE", Name: "NODE", Cost: 3, Faction: config.FactionArcTech, Class: config.ClassSupport,
		HP: 480, ATK: 30, AtkSpeed: 0.6, Range: 2, Armor: 9,
		AtkType: config.AttackRanged, DmgType: config.DamagePhys, Targeting: config.TargetNearest,
		SkillDesc: "Deploy + CDR"},

	// 4-cost (1)
	{ID: "ORISON", Name: "ORISON", Cost: 4, Faction: config.FactionExorcist, Class: config.ClassCaster,
		HP: 500, ATK: 85, AtkSpeed: 0.6, Range: 3, Armor: 8,
		AtkType: config.AttackRanged, DmgType: config.DamageMagic, Targeting: config.TargetFrontmost,
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
