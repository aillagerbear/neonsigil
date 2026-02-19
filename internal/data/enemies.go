package data

import (
	"image/color"

	"neonsigil/internal/config"
)

// EnemyDefs contains all enemy type definitions
var EnemyDefs = map[config.EnemyType]*EnemyDef{
	config.EnemyRunner: {
		Type: config.EnemyRunner, Name: "RUNNER", BaseHP: 80, Speed: 1.3,
		LeakDamage: 1, Color: color.RGBA{100, 255, 100, 255},
	},
	config.EnemyBruiser: {
		Type: config.EnemyBruiser, Name: "BRUISER", BaseHP: 220, Speed: 0.8,
		LeakDamage: 2, Color: color.RGBA{200, 100, 50, 255},
	},
	config.EnemyShield: {
		Type: config.EnemyShield, Name: "SHIELD", BaseHP: 180, Speed: 0.9,
		LeakDamage: 2, Color: color.RGBA{100, 150, 255, 255}, ShieldPct: 0.3,
	},
	config.EnemySplitter: {
		Type: config.EnemySplitter, Name: "SPLITTER", BaseHP: 140, Speed: 1.0,
		LeakDamage: 1, Color: color.RGBA{255, 200, 50, 255},
	},
	config.EnemyFlyer: {
		Type: config.EnemyFlyer, Name: "FLYER", BaseHP: 100, Speed: 1.2,
		LeakDamage: 1, Color: color.RGBA{200, 100, 255, 255},
	},
	config.EnemyStalker: {
		Type: config.EnemyStalker, Name: "STALKER", BaseHP: 110, Speed: 1.1,
		LeakDamage: 1, Color: color.RGBA{80, 80, 80, 255},
	},
	config.EnemyHacker: {
		Type: config.EnemyHacker, Name: "HACKER", BaseHP: 160, Speed: 0.95,
		LeakDamage: 2, Color: color.RGBA{0, 255, 200, 255},
	},
	config.EnemyCharger: {
		Type: config.EnemyCharger, Name: "CHARGER", BaseHP: 240, Speed: 1.15,
		LeakDamage: 3, Color: color.RGBA{255, 80, 80, 255},
	},
	config.EnemyTotem: {
		Type: config.EnemyTotem, Name: "TOTEM", BaseHP: 300, Speed: 0.7,
		LeakDamage: 3, Color: color.RGBA{255, 255, 100, 255},
	},
	config.EnemyBoss: {
		Type: config.EnemyBoss, Name: "GATEKEEPER", BaseHP: 2000, Speed: 0.5,
		LeakDamage: 99, Color: color.RGBA{255, 50, 50, 255},
	},
}
