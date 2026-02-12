package game

import (
	"math"
	"math/rand"

	"ebitengine-testing/config"
	"ebitengine-testing/entity"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) initBattle() {
	g.state = entity.StateBattle
	g.wave = 0
	g.summonerHP = config.SummonerMaxHP
	g.mana = float64(config.StartMana)
	g.manaTimer = 0
	g.selectedCard = -1
	g.fireballMode = false
	g.summoners = nil
	g.enemies = nil
	g.projectiles = nil
	g.grid = [config.GridRows][config.GridCols]*entity.Summoner{}
	g.graveyard = nil
	g.ticks = 0

	// 초기 덱 생성: 보병 4장 + 궁수 4장
	g.deck = nil
	for i := 0; i < 4; i++ {
		g.deck = append(g.deck, entity.Card{Data: entity.CardTemplates[entity.CardSoldier]})
	}
	for i := 0; i < 4; i++ {
		g.deck = append(g.deck, entity.Card{Data: entity.CardTemplates[entity.CardArcher]})
	}
	g.shuffleDeck()

	// 초기 핸드 드로우
	g.hand = nil
	for i := 0; i < config.MaxHandSize && len(g.deck) > 0; i++ {
		g.hand = append(g.hand, g.deck[0])
		g.deck = g.deck[1:]
	}

	g.startWave()
}

func (g *Game) shuffleDeck() {
	for i := len(g.deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		g.deck[i], g.deck[j] = g.deck[j], g.deck[i]
	}
}

func (g *Game) drawCard() {
	if len(g.hand) >= config.MaxHandSize {
		return
	}
	if len(g.deck) == 0 {
		if len(g.graveyard) == 0 {
			return
		}
		g.deck = append(g.deck, g.graveyard...)
		g.graveyard = nil
		g.shuffleDeck()
	}
	if len(g.deck) > 0 {
		g.hand = append(g.hand, g.deck[0])
		g.deck = g.deck[1:]
	}
}

func (g *Game) startWave() {
	g.spawnQueue = nil
	g.spawnTimer = 0
	g.waveComplete = false
	g.allSpawned = false
	g.enemies = nil
	g.projectiles = nil

	if g.wave >= len(entity.Waves) {
		return
	}

	w := entity.Waves[g.wave]
	for i := 0; i < w.Goblins; i++ {
		g.spawnQueue = append(g.spawnQueue, entity.EnemyGoblin)
	}
	for i := 0; i < w.Orcs; i++ {
		g.spawnQueue = append(g.spawnQueue, entity.EnemyOrc)
	}
	for i := 0; i < w.BossOrc; i++ {
		g.spawnQueue = append(g.spawnQueue, entity.EnemyBossOrc)
	}
	for i := 0; i < w.FinalBoss; i++ {
		g.spawnQueue = append(g.spawnQueue, entity.EnemyFinalBoss)
	}

	// 셔플 스폰 큐
	for i := len(g.spawnQueue) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		g.spawnQueue[i], g.spawnQueue[j] = g.spawnQueue[j], g.spawnQueue[i]
	}
}

func (g *Game) spawnEnemy(et entity.EnemyType) {
	e := &entity.Enemy{
		Type:      et,
		PathIndex: 0,
		PathT:     0,
		X:         entity.EnemyPath[0].X,
		Y:         entity.EnemyPath[0].Y,
	}
	switch et {
	case entity.EnemyGoblin:
		e.HP = 10
		e.MaxHP = 10
		e.Atk = 5
		e.Speed = 1.5
	case entity.EnemyOrc:
		e.HP = 30
		e.MaxHP = 30
		e.Atk = 10
		e.Speed = 0.8
	case entity.EnemyBossOrc:
		e.HP = 100
		e.MaxHP = 100
		e.Atk = 15
		e.Speed = 0.6
	case entity.EnemyFinalBoss:
		e.HP = 200
		e.MaxHP = 200
		e.Atk = 20
		e.Speed = 0.5
		e.HasAura = true
	}
	g.enemies = append(g.enemies, e)
}

func (g *Game) updateSynergies() {
	g.humanCount = 0
	g.elfCount = 0
	for _, s := range g.summoners {
		switch s.Card.Race {
		case entity.RaceHuman:
			g.humanCount++
		case entity.RaceElf:
			g.elfCount++
		}
	}
	g.humanSynergy = g.humanCount >= 2
	g.elfSynergy = g.elfCount >= 2

	// 시너지 스탯 적용
	for _, s := range g.summoners {
		baseData := entity.CardTemplates[s.Card.Type]
		s.MaxHP = baseData.HP
		s.Atk = baseData.Atk
		s.AtkSpeed = baseData.AtkSpeed

		if g.humanSynergy && s.Card.Race == entity.RaceHuman {
			s.MaxHP = int(float64(s.MaxHP) * 1.2)
			if s.CurrentHP > s.MaxHP {
				s.CurrentHP = s.MaxHP
			}
		}
		if g.elfSynergy && s.Card.Race == entity.RaceElf {
			s.AtkSpeed = int(float64(s.AtkSpeed) * 0.8) // 빠르게 = 쿨다운 감소
		}
	}
}

func (g *Game) updateTitle() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.initBattle()
	}
}

func (g *Game) updateBattle() {
	g.ticks++

	// 마나 회복
	g.manaTimer++
	if g.manaTimer >= config.ManaRegenTicks {
		g.manaTimer = 0
		if g.mana < float64(g.maxMana) {
			g.mana++
		}
	}

	// 적 스폰
	if !g.allSpawned && len(g.spawnQueue) > 0 {
		g.spawnTimer++
		if g.spawnTimer >= config.EnemySpawnInterval {
			g.spawnTimer = 0
			g.spawnEnemy(g.spawnQueue[0])
			g.spawnQueue = g.spawnQueue[1:]
			if len(g.spawnQueue) == 0 {
				g.allSpawned = true
			}
		}
	}

	// 적 이동
	for _, e := range g.enemies {
		if e.Dead || e.Reached {
			continue
		}
		g.moveEnemy(e)
	}

	// 최종 보스 오라 적용 (주변 적 공격력 증가)
	for _, e := range g.enemies {
		if e.HasAura && !e.Dead && !e.Reached {
			for _, other := range g.enemies {
				if other != e && !other.Dead && !other.Reached {
					dx := e.X - other.X
					dy := e.Y - other.Y
					if math.Sqrt(dx*dx+dy*dy) < 150 {
						// 오라 효과는 기본 스탯 기반으로 적용됨
					}
				}
			}
		}
	}

	// 소환수 공격
	for _, s := range g.summoners {
		s.AtkTimer++
		if s.AtkTimer >= s.AtkSpeed {
			target := g.findTarget(s)
			if target != nil {
				s.AtkTimer = 0
				if s.Range <= 1 {
					// 근접 공격: 즉시 피해
					target.HP -= s.Atk
				} else {
					// 원거리 공격: 투사체 생성
					p := &entity.Projectile{
						X:       s.ScreenX,
						Y:       s.ScreenY,
						TargetX: target.X,
						TargetY: target.Y,
						Damage:  s.Atk,
						Target:  target,
						Speed:   5,
					}
					g.projectiles = append(g.projectiles, p)
				}
			}
		}
	}

	// 투사체 이동
	var aliveProjectiles []*entity.Projectile
	for _, p := range g.projectiles {
		if p.Target != nil && (p.Target.Dead || p.Target.Reached) {
			continue
		}
		// 타겟 추적
		if p.Target != nil {
			p.TargetX = p.Target.X
			p.TargetY = p.Target.Y
		}
		dx := p.TargetX - p.X
		dy := p.TargetY - p.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < p.Speed+5 {
			// 명중
			if p.IsFireball {
				// AOE 피해
				for _, e := range g.enemies {
					if e.Dead || e.Reached {
						continue
					}
					edx := p.TargetX - e.X
					edy := p.TargetY - e.Y
					if math.Sqrt(edx*edx+edy*edy) < p.AOERange {
						e.HP -= p.Damage
					}
				}
			} else if p.Target != nil {
				p.Target.HP -= p.Damage
			}
			continue
		}
		p.X += dx / dist * p.Speed
		p.Y += dy / dist * p.Speed
		aliveProjectiles = append(aliveProjectiles, p)
	}
	g.projectiles = aliveProjectiles

	// 적 사망 / 도달 처리
	for _, e := range g.enemies {
		if e.Dead || e.Reached {
			continue
		}
		if e.HP <= 0 {
			e.Dead = true
		}
	}

	// 소환사 HP 체크
	if g.summonerHP <= 0 {
		g.state = entity.StateGameOver
		return
	}

	// 웨이브 완료 체크
	if g.allSpawned && len(g.spawnQueue) == 0 {
		allDead := true
		for _, e := range g.enemies {
			if !e.Dead && !e.Reached {
				allDead = false
				break
			}
		}
		if allDead {
			g.waveComplete = true
			if g.wave >= g.maxWave-1 {
				g.state = entity.StateVictory
			} else {
				g.prepareReward()
			}
		}
	}

	// 마우스 입력 처리 (카드 선택/배치)
	g.handleBattleInput()
}

func (g *Game) moveEnemy(e *entity.Enemy) {
	if e.PathIndex >= len(entity.EnemyPath)-1 {
		e.Reached = true
		g.summonerHP -= config.EnemyDamageToSummoner
		return
	}

	target := entity.EnemyPath[e.PathIndex+1]
	dx := target.X - e.X
	dy := target.Y - e.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < e.Speed+2 {
		e.PathIndex++
		if e.PathIndex >= len(entity.EnemyPath)-1 {
			e.Reached = true
			g.summonerHP -= config.EnemyDamageToSummoner
		}
		return
	}

	e.X += dx / dist * e.Speed
	e.Y += dy / dist * e.Speed
}

func (g *Game) findTarget(s *entity.Summoner) *entity.Enemy {
	var closest *entity.Enemy
	minDist := math.MaxFloat64
	rangePixels := float64(s.Range) * float64(config.TileSize)

	for _, e := range g.enemies {
		if e.Dead || e.Reached {
			continue
		}
		dx := s.ScreenX - e.X
		dy := s.ScreenY - e.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= rangePixels && dist < minDist {
			minDist = dist
			closest = e
		}
	}
	return closest
}

func (g *Game) handleBattleInput() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	mx, my := ebiten.CursorPosition()

	// 속도 조절 버튼
	if my < config.HUDHeight {
		if mx >= 880 && mx <= 920 {
			g.gameSpeed = 1
		} else if mx >= 930 && mx <= 970 {
			g.gameSpeed = 2
		}
		return
	}

	// 파이어볼 모드
	if g.fireballMode {
		fx := float64(mx)
		fy := float64(my)
		p := &entity.Projectile{
			X:          fx,
			Y:          fy - 100,
			TargetX:    fx,
			TargetY:    fy,
			Damage:     20,
			Speed:      8,
			IsFireball: true,
			AOERange:   float64(config.TileSize) * 1.5,
		}
		g.projectiles = append(g.projectiles, p)
		g.fireballMode = false
		g.selectedCard = -1
		return
	}

	// 핸드 카드 클릭 (하단)
	handY := config.ScreenHeight - config.HandHeight
	if my >= handY {
		cardWidth := 130
		cardSpacing := 10
		totalWidth := len(g.hand)*(cardWidth+cardSpacing) - cardSpacing
		startX := (config.ScreenWidth - totalWidth) / 2

		for i := range g.hand {
			cx := startX + i*(cardWidth+cardSpacing)
			if mx >= cx && mx <= cx+cardWidth && my >= handY+10 && my <= handY+config.HandHeight-10 {
				if g.selectedCard == i {
					g.selectedCard = -1
				} else {
					g.selectedCard = i
					if g.hand[i].Data.Type == entity.CardFireball {
						if g.mana >= float64(g.hand[i].Data.Cost) {
							g.mana -= float64(g.hand[i].Data.Cost)
							g.graveyard = append(g.graveyard, g.hand[i])
							g.hand = append(g.hand[:i], g.hand[i+1:]...)
							g.drawCard()
							g.fireballMode = true
						} else {
							g.selectedCard = -1
						}
					}
				}
				return
			}
		}
		return
	}

	// 그리드 타일 클릭 (배치)
	if g.selectedCard >= 0 && g.selectedCard < len(g.hand) {
		card := g.hand[g.selectedCard]
		if card.Data.Type == entity.CardFireball {
			return // 파이어볼은 위에서 처리
		}
		gx := (mx - config.GridStartX) / config.TileSize
		gy := (my - config.GridStartY) / config.TileSize
		if gx >= 0 && gx < config.GridCols && gy >= 0 && gy < config.GridRows {
			if g.grid[gy][gx] == nil && g.mana >= float64(card.Data.Cost) {
				g.mana -= float64(card.Data.Cost)
				s := &entity.Summoner{
					Card:      card.Data,
					GridX:     gx,
					GridY:     gy,
					CurrentHP: card.Data.HP,
					MaxHP:     card.Data.HP,
					Atk:       card.Data.Atk,
					AtkSpeed:  card.Data.AtkSpeed,
					Range:     card.Data.Range,
					ScreenX:   float64(config.GridStartX + gx*config.TileSize + config.TileSize/2),
					ScreenY:   float64(config.GridStartY + gy*config.TileSize + config.TileSize/2),
				}
				g.grid[gy][gx] = s
				g.summoners = append(g.summoners, s)
				g.graveyard = append(g.graveyard, g.hand[g.selectedCard])
				g.hand = append(g.hand[:g.selectedCard], g.hand[g.selectedCard+1:]...)
				g.selectedCard = -1
				g.drawCard()
				g.updateSynergies()
			}
		}
	}
}

func (g *Game) prepareReward() {
	g.state = entity.StateReward
	g.rewardCards = nil
	g.rewardHover = -1

	allTypes := []entity.CardType{entity.CardSoldier, entity.CardArcher, entity.CardSpearman, entity.CardMage, entity.CardFireball}
	// 3장의 랜덤 카드 제시
	for i := 0; i < 3; i++ {
		ct := allTypes[rand.Intn(len(allTypes))]
		g.rewardCards = append(g.rewardCards, entity.Card{Data: entity.CardTemplates[ct]})
	}
}

func (g *Game) updateReward() {
	mx, my := ebiten.CursorPosition()

	// 보상 카드 호버/클릭
	cardWidth := 180
	cardHeight := 250
	spacing := 30
	totalWidth := 3*(cardWidth+spacing) - spacing
	startX := (config.ScreenWidth - totalWidth) / 2
	startY := (config.ScreenHeight - cardHeight) / 2

	g.rewardHover = -1
	for i := 0; i < 3 && i < len(g.rewardCards); i++ {
		cx := startX + i*(cardWidth+spacing)
		if mx >= cx && mx <= cx+cardWidth && my >= startY && my <= startY+cardHeight {
			g.rewardHover = i
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && g.rewardHover >= 0 {
		// 선택한 카드를 덱에 추가
		g.deck = append(g.deck, g.rewardCards[g.rewardHover])
		g.shuffleDeck()

		// 핸드 보충
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
}
