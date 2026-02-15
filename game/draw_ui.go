package game

import (
	"fmt"
	"image/color"
	"math"

	"ebitengine-testing/config"
	"ebitengine-testing/entity"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ---- Card UI Colors ----

var (
	uiColorCardBorder   = color.NRGBA{0x4a, 0x4a, 0x6a, 0xFF}
	uiColorHumanGold    = color.NRGBA{0xFF, 0xCC, 0x44, 0xFF}
	uiColorElfEmerald   = color.NRGBA{0x44, 0xCC, 0x66, 0xFF}
	uiColorCostBadge    = color.NRGBA{0x29, 0x6B, 0xCC, 0xFF}
	uiColorCostBadgeRed = color.NRGBA{0xCC, 0x33, 0x33, 0xFF}
	uiColorGlowBlue     = color.NRGBA{0x29, 0xAD, 0xFF, 0x40}
	uiColorGlowGreen    = color.NRGBA{0x00, 0xE4, 0x36, 0x50}
	uiColorGlowGold     = color.NRGBA{0xFF, 0xCC, 0x44, 0x40}
	uiColorNameBanner   = color.NRGBA{0x10, 0x10, 0x20, 0xCC}
	uiColorDisabledOvr  = color.NRGBA{0x00, 0x00, 0x00, 0x80}
)

// ---- Helper: draw glow rectangle ----

func drawGlowRect(screen *ebiten.Image, x, y, w, h float32, glowColor color.NRGBA, glowSize float32) {
	for i := float32(1); i <= glowSize; i++ {
		alpha := byte(float32(glowColor.A) * (1.0 - i/glowSize))
		c := color.NRGBA{glowColor.R, glowColor.G, glowColor.B, alpha}
		vector.StrokeRect(screen, x-i, y-i, w+i*2, h+i*2, 1, c, false)
	}
}

// ---- Helper: get race border color ----

func raceBorderColor(race entity.Race) color.NRGBA {
	switch race {
	case entity.RaceHuman:
		return uiColorHumanGold
	case entity.RaceElf:
		return uiColorElfEmerald
	default:
		return uiColorCardBorder
	}
}

// ---- Helper: get unit sprite ----

func (g *Game) getCardSprite(cardType entity.CardType) *ebiten.Image {
	switch cardType {
	case entity.CardSoldier:
		return g.sprites.Soldier
	case entity.CardArcher:
		return g.sprites.Archer
	case entity.CardSpearman:
		return g.sprites.Spearman
	case entity.CardMage:
		return g.sprites.Mage
	case entity.CardFireball:
		return g.sprites.Fireball
	default:
		return nil
	}
}

// ---- Draw Cost Badge ----

func drawCostBadge(screen *ebiten.Image, cost int, x, y float64, affordable bool) {
	badgeColor := uiColorCostBadge
	if !affordable {
		badgeColor = uiColorCostBadgeRed
	}

	// Circle background
	vector.FillCircle(screen, float32(x), float32(y), 13, color.NRGBA{0x08, 0x08, 0x15, 0xFF}, false)
	vector.FillCircle(screen, float32(x), float32(y), 11, badgeColor, false)

	// Cost number
	costStr := fmt.Sprintf("%d", cost)
	op := &text.DrawOptions{}
	tw, th := text.Measure(costStr, fontCardCost, 0)
	op.GeoM.Translate(x-tw/2, y-th/2)
	op.ColorScale.ScaleWithColor(color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	text.Draw(screen, costStr, fontCardCost, op)
}

// ---- Draw Single Card ----

func (g *Game) drawSingleCard(screen *ebiten.Image, card entity.CardData, x, y float64, w, h float64, isHover, isSelected, canAfford bool) {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)

	// Vertical offset for hover/select
	offsetY := float32(0)
	if isSelected {
		offsetY = -10
	} else if isHover {
		offsetY = -5
	}
	fy += offsetY

	// Outer glow
	if isSelected {
		pulse := float32(0.6 + 0.4*math.Sin(float64(g.animTick)*0.1))
		glowC := uiColorGlowGreen
		glowC.A = byte(float32(glowC.A) * pulse)
		drawGlowRect(screen, fx, fy, fw, fh, glowC, 6)
	} else if isHover && canAfford {
		drawGlowRect(screen, fx, fy, fw, fh, uiColorGlowBlue, 4)
	}

	// Card background (gradient image)
	var bgImg *ebiten.Image
	if isHover && canAfford {
		bgImg = g.sprites.CardBGHover
	} else {
		bgImg = g.sprites.CardBG
	}

	op := &ebiten.DrawImageOptions{}
	// Scale to target size
	scaleX := w / float64(bgImg.Bounds().Dx())
	scaleY := h / float64(bgImg.Bounds().Dy())
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(x, float64(fy))
	screen.DrawImage(bgImg, op)

	// Border
	borderColor := raceBorderColor(card.Race)
	if isSelected {
		borderColor = color.NRGBA{0x00, 0xE4, 0x36, 0xFF}
	} else if isHover && canAfford {
		borderColor = uiColorAccent
	}
	vector.StrokeRect(screen, fx, fy, fw, fh, 2, borderColor, false)

	// Inner content area positions
	innerX := float64(fx) + 4
	innerW := float64(fw) - 8
	curY := float64(fy) + 6

	// Cost badge (top-left)
	drawCostBadge(screen, card.Cost, float64(fx)+16, float64(fy)+16, canAfford)

	// Unit sprite (centered in upper portion)
	sprite := g.getCardSprite(card.Type)
	if sprite != nil {
		spriteScale := 5.0
		if card.Type == entity.CardFireball {
			spriteScale = 6.0
		}
		spriteCenterX := float64(fx) + float64(fw)/2
		spriteCenterY := float64(fy) + 50
		drawSpriteAt(screen, sprite, spriteCenterX, spriteCenterY, spriteScale)
	}

	curY += 75

	// Name banner background
	vector.FillRect(screen, fx+2, float32(curY), fw-4, 22, uiColorNameBanner, false)

	// Name text (centered)
	nameOp := &text.DrawOptions{}
	nw, _ := text.Measure(card.Name, fontCardName, 0)
	nameOp.GeoM.Translate(innerX+(innerW-nw)/2, curY+3)
	nameOp.ColorScale.ScaleWithColor(color.RGBA{0xFF, 0xF1, 0xE8, 0xFF})
	text.Draw(screen, card.Name, fontCardName, nameOp)
	curY += 26

	// Stats section
	if card.Type != entity.CardFireball {
		// HP with heart icon
		g.drawIconStat(screen, g.sprites.HeartIcon, fmt.Sprintf("%d", card.HP),
			innerX+4, curY, color.RGBA{0xFF, 0x60, 0x60, 0xFF})

		// ATK with sword icon
		g.drawIconStat(screen, g.sprites.SwordIcon, fmt.Sprintf("%d", card.Atk),
			innerX+float64(fw)/2-4, curY, color.RGBA{0xFF, 0xCC, 0x44, 0xFF})
		curY += 18

		// Range
		rangeText := "근접"
		if card.Range > 1 {
			rangeText = fmt.Sprintf("사거리 %d", card.Range)
		}
		g.drawIconStat(screen, g.sprites.RangeIcon, rangeText,
			innerX+4, curY, color.RGBA{0xAA, 0xCC, 0xFF, 0xFF})
		curY += 18
	} else {
		// Fireball stats
		drawKoreanText(screen, "피해 20 광역", fontCardInfo, innerX+4, curY+2,
			color.RGBA{0xFF, 0x99, 0x33, 0xFF})
		curY += 18
		drawKoreanText(screen, "타겟 클릭", fontCardInfo, innerX+4, curY+2,
			color.RGBA{0xCC, 0xCC, 0xCC, 0xFF})
		curY += 18
	}

	// Race tag bar at bottom
	if card.Race != entity.RaceNone {
		tagY := float32(y+float64(offsetY)) + float32(h) - 22
		tagColor := raceBorderColor(card.Race)
		tagBG := color.NRGBA{tagColor.R, tagColor.G, tagColor.B, 0x40}
		vector.FillRect(screen, fx+2, tagY, fw-4, 18, tagBG, false)

		raceName := ""
		switch card.Race {
		case entity.RaceHuman:
			raceName = "인간"
		case entity.RaceElf:
			raceName = "엘프"
		}
		rw, _ := text.Measure(raceName, fontCardInfo, 0)
		raceOp := &text.DrawOptions{}
		raceOp.GeoM.Translate(innerX+(innerW-rw)/2, float64(tagY)+2)
		raceOp.ColorScale.ScaleWithColor(tagColor)
		text.Draw(screen, raceName, fontCardInfo, raceOp)
	}

	// Unaffordable overlay
	if !canAfford && !isSelected {
		vector.FillRect(screen, fx+1, fy+1, fw-2, fh-2, uiColorDisabledOvr, false)
	}
}

// drawIconStat draws a small icon + text stat
func (g *Game) drawIconStat(screen *ebiten.Image, icon *ebiten.Image, txt string, x, y float64, clr color.RGBA) {
	if icon != nil {
		iconOp := &ebiten.DrawImageOptions{}
		iconOp.GeoM.Scale(2, 2)
		iconOp.GeoM.Translate(x, y+1)
		screen.DrawImage(icon, iconOp)
	}
	drawKoreanText(screen, txt, fontCardInfo, x+14, y+2, clr)
}

// ---- Draw Hand Cards ----

func (g *Game) drawHandCards(screen *ebiten.Image) {
	if len(g.hand) == 0 {
		return
	}

	cardW := float64(config.CardWidth)
	cardH := float64(config.CardHeight)
	spacing := float64(config.CardSpacing)

	totalW := float64(len(g.hand))*cardW + float64(len(g.hand)-1)*spacing
	startX := (float64(config.ScreenWidth) - totalW) / 2
	startY := float64(config.ScreenHeight) - cardH - 12

	// Hand area background panel
	panelX := float32(startX) - 12
	panelY := float32(startY) - 8
	panelW := float32(totalW) + 24
	panelH := float32(cardH) + 20
	vector.FillRect(screen, panelX, panelY, panelW, panelH, color.NRGBA{0x12, 0x12, 0x20, 0xE0}, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 1, color.NRGBA{0x3a, 0x3a, 0x5a, 0x80}, false)

	// Reset card rects
	g.cardRects = g.cardRects[:0]

	for i, card := range g.hand {
		cx := startX + float64(i)*(cardW+spacing)
		cy := startY
		canAfford := g.mana >= float64(card.Data.Cost)
		isHover := g.hoverCard == i
		isSelected := g.selectedCard == i

		g.drawSingleCard(screen, card.Data, cx, cy, cardW, cardH, isHover, isSelected, canAfford)

		g.cardRects = append(g.cardRects, CardRect{
			X: cx, Y: cy - 10, // include hover offset area
			W: cardW, H: cardH + 10,
			Index: i,
		})
	}
}

// ---- Draw Enhanced HP/Mana Bar ----

func drawEnhancedBar(screen *ebiten.Image, x, y, w, h float32, ratio float64, fillColor color.RGBA, label string) {
	// Background
	vector.FillRect(screen, x, y, w, h, color.RGBA{0x1a, 0x1a, 0x1a, 0xFF}, false)

	// Fill
	fillW := float32(ratio) * w
	if fillW > 0 {
		vector.FillRect(screen, x, y, fillW, h, fillColor, false)

		// Brighter top portion (gradient effect)
		brightColor := lerpColor(fillColor, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, 0.25)
		topH := h * 0.35
		vector.FillRect(screen, x, y, fillW, float32(topH), brightColor, false)

		// Highlight edge
		hlColor := lerpColor(fillColor, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, 0.45)
		vector.FillRect(screen, x+1, y+1, fillW-2, 1, hlColor, false)
	}

	// Border
	vector.StrokeRect(screen, x, y, w, h, 1, color.RGBA{0x00, 0x00, 0x00, 0xCC}, false)

	// Text label centered on bar
	if label != "" {
		lw, lh := text.Measure(label, fontBarValue, 0)
		lx := float64(x) + float64(w)/2 - lw/2
		ly := float64(y) + float64(h)/2 - lh/2
		// Shadow
		drawKoreanText(screen, label, fontBarValue, lx+1, ly+1, color.RGBA{0x00, 0x00, 0x00, 0xCC})
		// Text
		drawKoreanText(screen, label, fontBarValue, lx, ly, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	}
}

// ---- Draw Wave Progress ----

func drawWaveProgress(screen *ebiten.Image, currentWave, maxWave, animTick int) {
	pipSize := float32(10)
	pipSpacing := float32(4)
	totalW := float32(maxWave)*pipSize + float32(maxWave-1)*pipSpacing
	startX := float32(config.ScreenWidth)/2 - totalW/2
	y := float32(30)

	for i := 0; i < maxWave; i++ {
		px := startX + float32(i)*(pipSize+pipSpacing) + pipSize/2
		py := y + pipSize/2

		if i < currentWave {
			// Completed - filled bright
			vector.FillCircle(screen, px, py, pipSize/2, color.NRGBA{0x29, 0xAD, 0xFF, 0xFF}, false)
			vector.StrokeCircle(screen, px, py, pipSize/2, 1, color.NRGBA{0x60, 0xCF, 0xFF, 0xFF}, false)
		} else if i == currentWave {
			// Current - pulsing
			pulse := byte(160 + int(95*math.Sin(float64(animTick)*0.08)))
			vector.FillCircle(screen, px, py, pipSize/2, color.NRGBA{0x29, 0xAD, 0xFF, pulse}, false)
			vector.StrokeCircle(screen, px, py, pipSize/2+1, 1, color.NRGBA{0x29, 0xAD, 0xFF, pulse / 2}, false)
		} else {
			// Future - dim outline
			vector.StrokeCircle(screen, px, py, pipSize/2, 1, color.NRGBA{0x3a, 0x3a, 0x5a, 0xFF}, false)
		}
	}
}

// ---- Draw Custom HUD (replaces old drawHPManaBar) ----

func (g *Game) drawCustomHUD(screen *ebiten.Image) {
	// HP bar
	hpRatio := float64(g.summonerHP) / float64(config.SummonerMaxHP)
	hpColor := color.RGBA{0x00, 0xCC, 0x44, 0xFF} // green
	if hpRatio < 0.3 {
		hpColor = color.RGBA{0xFF, 0x33, 0x44, 0xFF} // red
	} else if hpRatio < 0.6 {
		hpColor = color.RGBA{0xFF, 0xCC, 0x33, 0xFF} // yellow
	}
	hpLabel := fmt.Sprintf("%d / %d", g.summonerHP, config.SummonerMaxHP)
	drawEnhancedBar(screen, 8, 56, float32(config.BarWidth), float32(config.BarHeight), hpRatio, hpColor, hpLabel)

	// Mana bar
	manaRatio := g.mana / float64(g.maxMana)
	manaColor := color.RGBA{0x29, 0x99, 0xFF, 0xFF}
	manaLabel := fmt.Sprintf("%.0f / %d", g.mana, g.maxMana)
	drawEnhancedBar(screen, 8, 100, float32(config.BarWidth), float32(config.BarHeight), manaRatio, manaColor, manaLabel)

	// Mana regen indicator
	if g.mana < float64(g.maxMana) {
		regenProgress := float64(g.manaTimer) / float64(config.ManaRegenTicks)
		regenW := float32(float64(config.BarWidth) * regenProgress)
		vector.FillRect(screen, 8, float32(100+config.BarHeight+1), regenW, 3, color.RGBA{0x29, 0xAD, 0xFF, 0x60}, false)
	}

	// Wave progress pips
	drawWaveProgress(screen, g.wave, g.maxWave, g.animTick)
}

// ---- Reward Screen ----

func (g *Game) drawRewardScreen(screen *ebiten.Image) {
	// Dark overlay
	vector.FillRect(screen, 0, 0, config.ScreenWidth, config.ScreenHeight, color.NRGBA{0x00, 0x00, 0x00, 0xCC}, false)

	// Title with glow
	titleText := fmt.Sprintf("웨이브 %d 클리어!", g.wave+1)
	tw, _ := text.Measure(titleText, fontTitle, 0)
	titleX := float64(config.ScreenWidth)/2 - tw/2
	titleY := float64(80)

	// Glow effect (multi-pass at low alpha)
	glowColor := color.RGBA{0xFF, 0xCC, 0x44, 0x30}
	for _, off := range [][2]float64{{-2, 0}, {2, 0}, {0, -2}, {0, 2}, {-1, -1}, {1, 1}, {-1, 1}, {1, -1}} {
		drawKoreanText(screen, titleText, fontTitle, titleX+off[0], titleY+off[1], glowColor)
	}
	drawKoreanText(screen, titleText, fontTitle, titleX, titleY, color.RGBA{0xFF, 0xDD, 0x55, 0xFF})

	// Subtitle
	subText := "카드를 선택하세요"
	sw, _ := text.Measure(subText, fontLarge, 0)
	drawKoreanTextWithShadow(screen, subText, fontLarge,
		float64(config.ScreenWidth)/2-sw/2, titleY+40,
		color.RGBA{0xCC, 0xCC, 0xCC, 0xFF})

	// Reward cards
	cardW := float64(config.RewardCardWidth)
	cardH := float64(config.RewardCardHeight)
	spacing := float64(config.RewardCardSpacing)
	numCards := len(g.rewardCards)
	totalW := float64(numCards)*cardW + float64(numCards-1)*spacing
	startX := (float64(config.ScreenWidth) - totalW) / 2
	startY := float64(160)

	g.rewardCardRects = g.rewardCardRects[:0]

	for i, card := range g.rewardCards {
		cx := startX + float64(i)*(cardW+spacing)
		cy := startY
		isHover := g.hoverReward == i

		g.drawRewardCard(screen, card.Data, cx, cy, cardW, cardH, isHover)

		g.rewardCardRects = append(g.rewardCardRects, CardRect{
			X: cx, Y: cy - 10,
			W: cardW, H: cardH + 10,
			Index: i,
		})
	}
}

// drawRewardCard draws a single reward card (larger version)
func (g *Game) drawRewardCard(screen *ebiten.Image, card entity.CardData, x, y, w, h float64, isHover bool) {
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)

	offsetY := float32(0)
	if isHover {
		offsetY = -8
	}
	fy += offsetY

	// Glow
	if isHover {
		drawGlowRect(screen, fx, fy, fw, fh, uiColorGlowGold, 6)
	}

	// Background
	bgOp := &ebiten.DrawImageOptions{}
	scaleX := w / float64(g.sprites.RewardCardBG.Bounds().Dx())
	scaleY := h / float64(g.sprites.RewardCardBG.Bounds().Dy())
	bgOp.GeoM.Scale(scaleX, scaleY)
	bgOp.GeoM.Translate(x, float64(fy))
	screen.DrawImage(g.sprites.RewardCardBG, bgOp)

	// Border
	borderColor := raceBorderColor(card.Race)
	if isHover {
		borderColor = uiColorGlowGold
		borderColor.A = 0xFF
	}
	vector.StrokeRect(screen, fx, fy, fw, fh, 2, borderColor, false)

	innerX := float64(fx) + 6
	innerW := float64(fw) - 12

	// Cost badge
	drawCostBadge(screen, card.Cost, float64(fx)+18, float64(fy)+18, true)

	// Unit sprite
	sprite := g.getCardSprite(card.Type)
	if sprite != nil {
		spriteScale := 6.0
		if card.Type == entity.CardFireball {
			spriteScale = 7.0
		}
		drawSpriteAt(screen, sprite, float64(fx)+float64(fw)/2, float64(fy)+70, spriteScale)
	}

	curY := float64(fy) + 110

	// Name banner
	vector.FillRect(screen, fx+2, float32(curY), fw-4, 24, uiColorNameBanner, false)
	nw, _ := text.Measure(card.Name, fontLarge, 0)
	nameOp := &text.DrawOptions{}
	nameOp.GeoM.Translate(innerX+(innerW-nw)/2, curY+3)
	nameOp.ColorScale.ScaleWithColor(color.RGBA{0xFF, 0xF1, 0xE8, 0xFF})
	text.Draw(screen, card.Name, fontLarge, nameOp)
	curY += 30

	// Stats
	if card.Type != entity.CardFireball {
		g.drawIconStat(screen, g.sprites.HeartIcon, fmt.Sprintf("체력: %d", card.HP),
			innerX+4, curY, color.RGBA{0xFF, 0x60, 0x60, 0xFF})
		curY += 20
		g.drawIconStat(screen, g.sprites.SwordIcon, fmt.Sprintf("공격: %d", card.Atk),
			innerX+4, curY, color.RGBA{0xFF, 0xCC, 0x44, 0xFF})
		curY += 20

		speedText := "보통"
		if card.AtkSpeed < 50 {
			speedText = "빠름"
		} else if card.AtkSpeed > 70 {
			speedText = "느림"
		}
		drawKoreanText(screen, fmt.Sprintf("속도: %s", speedText), fontCardInfo,
			innerX+4, curY+2, color.RGBA{0xCC, 0xCC, 0xCC, 0xFF})
		curY += 20

		rangeText := "근접"
		if card.Range > 1 {
			rangeText = fmt.Sprintf("%d칸", card.Range)
		}
		g.drawIconStat(screen, g.sprites.RangeIcon, fmt.Sprintf("사거리: %s", rangeText),
			innerX+4, curY, color.RGBA{0xAA, 0xCC, 0xFF, 0xFF})
		curY += 24
	} else {
		drawKoreanText(screen, "마법 카드", fontCardInfo, innerX+4, curY+2,
			color.RGBA{0xFF, 0x99, 0x33, 0xFF})
		curY += 20
		drawKoreanText(screen, "피해: 20 (광역)", fontCardInfo, innerX+4, curY+2,
			color.RGBA{0xFF, 0xCC, 0x44, 0xFF})
		curY += 20
		drawKoreanText(screen, "주변 적 모두에게", fontCardInfo, innerX+4, curY+2,
			color.RGBA{0xCC, 0xCC, 0xCC, 0xFF})
		curY += 20
		drawKoreanText(screen, "피해를 줍니다", fontCardInfo, innerX+4, curY+2,
			color.RGBA{0xCC, 0xCC, 0xCC, 0xFF})
		curY += 24
	}

	// Race tag
	if card.Race != entity.RaceNone {
		tagY := float32(float64(fy) + h - float64(offsetY) - 26)
		tagColor := raceBorderColor(card.Race)
		tagBG := color.NRGBA{tagColor.R, tagColor.G, tagColor.B, 0x40}
		vector.FillRect(screen, fx+2, tagY, fw-4, 20, tagBG, false)

		raceName := ""
		var emblem *ebiten.Image
		switch card.Race {
		case entity.RaceHuman:
			raceName = "인간"
			emblem = g.sprites.HumanEmblem
		case entity.RaceElf:
			raceName = "엘프"
			emblem = g.sprites.ElfEmblem
		}

		// Emblem + name
		rw, _ := text.Measure(raceName, fontCardInfo, 0)
		totalTagW := 16 + rw // emblem + spacing + text
		tagStartX := float64(fx) + float64(fw)/2 - totalTagW/2

		if emblem != nil {
			eOp := &ebiten.DrawImageOptions{}
			eOp.GeoM.Scale(2, 2)
			eOp.GeoM.Translate(tagStartX, float64(tagY)+2)
			screen.DrawImage(emblem, eOp)
		}

		raceOp := &text.DrawOptions{}
		raceOp.GeoM.Translate(tagStartX+16, float64(tagY)+3)
		raceOp.ColorScale.ScaleWithColor(tagColor)
		text.Draw(screen, raceName, fontCardInfo, raceOp)
	}

	// "클릭하여 추가" prompt at bottom
	promptText := "클릭하여 추가"
	pw, _ := text.Measure(promptText, fontSmall, 0)
	promptColor := color.RGBA{0x88, 0x88, 0x99, 0xFF}
	if isHover {
		promptColor = color.RGBA{0xFF, 0xDD, 0x55, 0xFF}
	}
	drawKoreanText(screen, promptText, fontSmall,
		float64(fx)+float64(fw)/2-pw/2, float64(fy)+h-float64(offsetY)-6, promptColor)
}

// ---- Title Screen Decorations ----

func (g *Game) drawTitleDecorations(screen *ebiten.Image) {
	// Animated floating sprites in background
	t := float64(g.animTick)

	// Soldier drifting
	sx := 150 + 30*math.Sin(t*0.01)
	sy := 300 + 20*math.Cos(t*0.015)
	drawSpriteAtWithColor(screen, g.sprites.Soldier, sx, sy, 4.0, 0.3, 0.3, 0.5, 0.3)

	// Archer drifting
	ax := 870 + 25*math.Cos(t*0.012)
	ay := 350 + 15*math.Sin(t*0.018)
	drawSpriteAtWithColor(screen, g.sprites.Archer, ax, ay, 4.0, 0.3, 0.5, 0.3, 0.3)

	// Goblin drifting
	gx := 200 + 20*math.Sin(t*0.008+1.5)
	gy := 500 + 15*math.Cos(t*0.013)
	drawSpriteAtWithColor(screen, g.sprites.Goblin, gx, gy, 3.0, 0.5, 0.3, 0.3, 0.25)

	// Mage drifting
	mx := 820 + 25*math.Cos(t*0.011+2.0)
	my := 200 + 20*math.Sin(t*0.016)
	drawSpriteAtWithColor(screen, g.sprites.Mage, mx, my, 4.0, 0.4, 0.3, 0.5, 0.25)

	// Decorative line separator
	lineY := float32(config.ScreenHeight/2 + 40)
	lineW := float32(300)
	lineX := float32(config.ScreenWidth)/2 - lineW/2
	for i := float32(0); i < lineW; i++ {
		alpha := byte(80 * (1 - math.Abs(float64(i-lineW/2))/float64(lineW/2)))
		vector.FillRect(screen, lineX+i, lineY, 1, 1, color.RGBA{0x29, 0xAD, 0xFF, alpha}, false)
	}
}

// ---- End Screen Effects ----

func (g *Game) drawEndScreenEffects(screen *ebiten.Image) {
	t := float64(g.animTick)

	if g.state == entity.StateVictory {
		// Golden particle-like stars
		for i := 0; i < 20; i++ {
			fi := float64(i)
			px := float64(config.ScreenWidth)/2 + 200*math.Cos(fi*0.7+t*0.02)
			py := float64(config.ScreenHeight)/2 + 150*math.Sin(fi*0.5+t*0.025)
			alpha := byte(100 + int(80*math.Sin(t*0.05+fi*1.2)))
			sz := float32(2 + math.Sin(t*0.06+fi)*1.5)
			vector.FillRect(screen, float32(px), float32(py), sz, sz,
				color.RGBA{0xFF, 0xCC, 0x44, alpha}, false)
		}
	} else {
		// Dim red vignette for game over
		vignetteAlpha := byte(30 + int(20*math.Sin(t*0.04)))
		// Top/bottom edge darkening
		for y := 0; y < 60; y++ {
			a := byte(float64(vignetteAlpha) * (1 - float64(y)/60.0))
			vector.FillRect(screen, 0, float32(y), config.ScreenWidth, 1,
				color.RGBA{0xFF, 0x00, 0x20, a}, false)
			vector.FillRect(screen, 0, float32(config.ScreenHeight-y), config.ScreenWidth, 1,
				color.RGBA{0xFF, 0x00, 0x20, a}, false)
		}
	}
}
