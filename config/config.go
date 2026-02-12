package config

const (
	ScreenWidth  = 1024
	ScreenHeight = 768

	// 그리드 설정
	GridCols   = 5
	GridRows   = 5
	TileSize   = 80
	GridStartX = 280
	GridStartY = 150

	// HUD 영역
	HUDHeight    = 40
	HandHeight   = 120
	SidebarWidth = 120

	// 게임 설정
	MaxMana               = 10
	StartMana             = 5
	ManaRegenTicks        = 120 // 2초 (60TPS * 2)
	MaxHandSize           = 5
	SummonerMaxHP         = 100
	EnemyDamageToSummoner = 10

	// 적 스폰 간격 (틱)
	EnemySpawnInterval = 40
)
