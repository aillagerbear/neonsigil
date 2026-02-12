package entity

// CardTemplates defines the base stats for each card type.
var CardTemplates = map[CardType]CardData{
	CardSoldier:  {Name: "전사", Type: CardSoldier, Race: RaceHuman, Cost: 2, HP: 50, Atk: 8, AtkSpeed: 60, Range: 1},
	CardArcher:   {Name: "궁수", Type: CardArcher, Race: RaceHuman, Cost: 3, HP: 25, Atk: 10, AtkSpeed: 60, Range: 3},
	CardSpearman: {Name: "창병", Type: CardSpearman, Race: RaceElf, Cost: 3, HP: 40, Atk: 10, AtkSpeed: 45, Range: 2},
	CardMage:     {Name: "마법사", Type: CardMage, Race: RaceElf, Cost: 4, HP: 20, Atk: 15, AtkSpeed: 90, Range: 3},
	CardFireball: {Name: "화염구", Type: CardFireball, Race: RaceNone, Cost: 3, HP: 0, Atk: 20, AtkSpeed: 0, Range: 0},
}

// Waves defines the enemy composition for each of the 10 waves.
var Waves = []WaveData{
	{Goblins: 5},
	{Goblins: 8},
	{Goblins: 10, Orcs: 1},
	{Goblins: 12, Orcs: 2},
	{BossOrc: 1},
	{Goblins: 15, Orcs: 3},
	{Orcs: 5},
	{Goblins: 20},
	{Goblins: 10, Orcs: 5},
	{FinalBoss: 1},
}

// EnemyPath defines the S-shaped waypoints enemies follow.
var EnemyPath = []Point{
	{X: 520, Y: -20},  // 시작 (화면 밖)
	{X: 520, Y: 80},   // 상단
	{X: 750, Y: 80},   // 오른쪽으로
	{X: 750, Y: 230},  // 아래로
	{X: 280, Y: 230},  // 왼쪽으로
	{X: 280, Y: 380},  // 아래로
	{X: 750, Y: 380},  // 오른쪽으로
	{X: 750, Y: 530},  // 아래로
	{X: 520, Y: 530},  // 가운데로
	{X: 520, Y: 700},  // 소환사 위치
}
