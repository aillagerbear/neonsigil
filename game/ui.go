package game

import (
	"fmt"
	"image/color"

	"ebitengine-testing/config"
	"ebitengine-testing/entity"

	"github.com/ebitenui/ebitenui"
	euiimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// UIManager manages all ebitenui screens
type UIManager struct {
	titleUI  *ebitenui.UI
	battleUI *ebitenui.UI
	rewardUI *ebitenui.UI
	endUI    *ebitenui.UI

	// Battle HUD labels (updated each frame)
	waveLabel     *widget.Label
	enemyLabel    *widget.Label
	fpsLabel      *widget.Label
	hpLabel       *widget.Label
	manaLabel     *widget.Label
	synHumanLabel *widget.Label
	synElfLabel   *widget.Label
	synHumanDesc  *widget.Label
	synElfDesc    *widget.Label
	deckLabel     *widget.Label
	graveLabel    *widget.Label

	// Speed buttons
	speedBtn1 *widget.Button
	speedBtn2 *widget.Button

	// Card buttons (bottom hand)
	cardButtons   []*widget.Button
	cardContainer *widget.Container

	// State tracking to avoid per-frame rebuild
	prevHandLen      int
	prevHandNames    string
	prevSelectedCard int
	prevManaFloor    int

	// Reward card buttons
	rewardButtons   []*widget.Button
	rewardTitle     *widget.Label
	rewardContainer *widget.Container

	// End screen
	endTitle   *widget.Label
	endMessage *widget.Label

	// Callbacks
	onStartGame     func()
	onSpeedChange   func(speed int)
	onCardSelect    func(index int)
	onRewardSelect  func(index int)
	onReturnToTitle func()
}

func newUIManager(g *Game) *UIManager {
	ui := &UIManager{}
	ui.onStartGame = func() { g.initBattle() }
	ui.onSpeedChange = func(speed int) { g.gameSpeed = speed }
	ui.onCardSelect = func(index int) { g.handleCardSelect(index) }
	ui.onRewardSelect = func(index int) { g.handleRewardSelect(index) }
	ui.onReturnToTitle = func() { g.state = entity.StateTitle }

	ui.prevHandLen = -1
	ui.prevSelectedCard = -99
	ui.prevManaFloor = -1

	ui.buildTitleUI()
	ui.buildBattleUI()
	ui.buildRewardUI()
	ui.buildEndUI()
	return ui
}

// ---- Color helpers ----

func nineSlice(c color.NRGBA) *euiimage.NineSlice {
	return euiimage.NewNineSliceColor(c)
}

func borderedNineSlice(bg, border color.NRGBA, borderWidth int) *euiimage.NineSlice {
	return euiimage.NewBorderedNineSliceColor(bg, border, borderWidth)
}

var (
	uiColorDarkBG    = color.NRGBA{0x12, 0x12, 0x20, 0xF0}
	uiColorBorder    = color.NRGBA{0x3a, 0x3a, 0x5a, 0xFF}
	uiColorAccent    = color.NRGBA{0x29, 0xAD, 0xFF, 0xFF}
	uiColorGreen     = color.NRGBA{0x00, 0xE4, 0x36, 0xFF}
	uiColorRed       = color.NRGBA{0xFF, 0x00, 0x4D, 0xFF}
	uiColorYellow    = color.NRGBA{0xFF, 0xEC, 0x27, 0xFF}
	uiColorOrange    = color.NRGBA{0xFF, 0xA3, 0x00, 0xFF}
	uiColorWhite     = color.NRGBA{0xFF, 0xF1, 0xE8, 0xFF}
	uiColorGrey      = color.NRGBA{0x80, 0x80, 0x90, 0xFF}
	uiColorCardBG    = color.NRGBA{0x25, 0x25, 0x3a, 0xF0}
	uiColorCardHover = color.NRGBA{0x35, 0x40, 0x55, 0xFF}
	uiColorOverlay   = color.NRGBA{0x00, 0x00, 0x00, 0xCC}
)

func labelColor(c color.NRGBA) *widget.LabelColor {
	return &widget.LabelColor{Idle: c, Disabled: uiColorGrey}
}

func buttonImage(idle, hover, pressed color.NRGBA) *widget.ButtonImage {
	return &widget.ButtonImage{
		Idle:    borderedNineSlice(idle, uiColorBorder, 2),
		Hover:   borderedNineSlice(hover, uiColorAccent, 2),
		Pressed: borderedNineSlice(pressed, uiColorAccent, 2),
	}
}

// facePtr returns a *text.Face from a text.Face for ebitenui
func facePtr(f text.Face) *text.Face {
	return &f
}

// ---- Title Screen ----

func (ui *UIManager) buildTitleUI() {
	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSlice(color.NRGBA{0x08, 0x08, 0x12, 0xFF})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	centerPanel := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(12),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
		)),
	)

	// Title
	titleLabel := widget.NewLabel(
		widget.LabelOpts.Text("소환사의 수호", facePtr(fontTitle), labelColor(uiColorAccent)),
		widget.LabelOpts.TextOpts(widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter)),
	)
	centerPanel.AddChild(titleLabel)

	// Subtitle
	subLabel := widget.NewLabel(
		widget.LabelOpts.Text("자동전투 덱빌딩 디펜스", facePtr(fontLarge), labelColor(uiColorWhite)),
		widget.LabelOpts.TextOpts(widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter)),
	)
	centerPanel.AddChild(subLabel)

	// Spacer
	spacer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(0, 20)),
	)
	centerPanel.AddChild(spacer)

	// Start button
	startBtn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage(
			color.NRGBA{0x00, 0x60, 0x30, 0xFF},
			color.NRGBA{0x00, 0x87, 0x51, 0xFF},
			color.NRGBA{0x00, 0x50, 0x28, 0xFF},
		)),
		widget.ButtonOpts.Text("게임 시작", facePtr(fontLarge), &widget.ButtonTextColor{
			Idle: uiColorWhite,
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 40, Right: 40, Top: 12, Bottom: 12}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.onStartGame()
		}),
	)
	centerPanel.AddChild(startBtn)

	// Spacer
	spacer2 := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(0, 15)),
	)
	centerPanel.AddChild(spacer2)

	// Instructions
	instructions := []string{
		"조작법:",
		"• 카드를 클릭하여 선택, 그리드에 배치",
		"• 유닛이 주변 적을 자동으로 공격합니다",
		"• 10 웨이브를 버텨서 승리하세요!",
	}
	for _, line := range instructions {
		l := widget.NewLabel(
			widget.LabelOpts.Text(line, facePtr(fontMedium), labelColor(uiColorGrey)),
			widget.LabelOpts.TextOpts(widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter)),
		)
		centerPanel.AddChild(l)
	}

	root.AddChild(centerPanel)
	ui.titleUI = &ebitenui.UI{Container: root}
}

// ---- Battle HUD ----

func (ui *UIManager) buildBattleUI() {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// === Top bar ===
	topBar := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionStart,
			StretchHorizontal:  true,
		})),
		widget.ContainerOpts.BackgroundImage(nineSlice(uiColorDarkBG)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(&widget.Insets{Left: 15, Right: 15, Top: 6, Bottom: 6}),
		)),
	)

	ui.waveLabel = widget.NewLabel(
		widget.LabelOpts.Text("웨이브 1/10", facePtr(fontMedium), labelColor(uiColorAccent)),
	)
	topBar.AddChild(ui.waveLabel)

	ui.enemyLabel = widget.NewLabel(
		widget.LabelOpts.Text("적: 0", facePtr(fontMedium), labelColor(uiColorOrange)),
	)
	topBar.AddChild(ui.enemyLabel)

	ui.fpsLabel = widget.NewLabel(
		widget.LabelOpts.Text("FPS: 60", facePtr(fontSmall), labelColor(uiColorGrey)),
	)
	topBar.AddChild(ui.fpsLabel)

	// Speed buttons
	speedContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(4),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionEnd,
		})),
	)

	ui.speedBtn1 = widget.NewButton(
		widget.ButtonOpts.Image(buttonImage(
			color.NRGBA{0x00, 0x87, 0x51, 0xFF},
			color.NRGBA{0x00, 0xA0, 0x60, 0xFF},
			color.NRGBA{0x00, 0x60, 0x38, 0xFF},
		)),
		widget.ButtonOpts.Text("1x", facePtr(fontSmall), &widget.ButtonTextColor{Idle: uiColorWhite}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 8, Right: 8, Top: 2, Bottom: 2}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.onSpeedChange(1)
		}),
	)
	speedContainer.AddChild(ui.speedBtn1)

	ui.speedBtn2 = widget.NewButton(
		widget.ButtonOpts.Image(buttonImage(
			color.NRGBA{0x30, 0x30, 0x45, 0xFF},
			color.NRGBA{0x00, 0xA0, 0x60, 0xFF},
			color.NRGBA{0x00, 0x60, 0x38, 0xFF},
		)),
		widget.ButtonOpts.Text("2x", facePtr(fontSmall), &widget.ButtonTextColor{Idle: uiColorWhite}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 8, Right: 8, Top: 2, Bottom: 2}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.onSpeedChange(2)
		}),
	)
	speedContainer.AddChild(ui.speedBtn2)
	topBar.AddChild(speedContainer)
	root.AddChild(topBar)

	// === Left sidebar ===
	leftSide := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
			Padding:            &widget.Insets{Top: 40},
		})),
		widget.ContainerOpts.BackgroundImage(nineSlice(uiColorDarkBG)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(6),
			widget.RowLayoutOpts.Padding(&widget.Insets{Left: 8, Right: 8, Top: 10, Bottom: 10}),
		)),
	)

	hpTitle := widget.NewLabel(
		widget.LabelOpts.Text("체력", facePtr(fontSmall), labelColor(uiColorRed)),
	)
	leftSide.AddChild(hpTitle)

	ui.hpLabel = widget.NewLabel(
		widget.LabelOpts.Text("100/100", facePtr(fontMedium), labelColor(uiColorWhite)),
	)
	leftSide.AddChild(ui.hpLabel)

	manaTitle := widget.NewLabel(
		widget.LabelOpts.Text("마나", facePtr(fontSmall), labelColor(uiColorAccent)),
	)
	leftSide.AddChild(manaTitle)

	ui.manaLabel = widget.NewLabel(
		widget.LabelOpts.Text("5/10", facePtr(fontMedium), labelColor(uiColorWhite)),
	)
	leftSide.AddChild(ui.manaLabel)

	ui.deckLabel = widget.NewLabel(
		widget.LabelOpts.Text("덱: 0", facePtr(fontSmall), labelColor(uiColorGrey)),
	)
	leftSide.AddChild(ui.deckLabel)

	ui.graveLabel = widget.NewLabel(
		widget.LabelOpts.Text("묘지: 0", facePtr(fontSmall), labelColor(uiColorGrey)),
	)
	leftSide.AddChild(ui.graveLabel)

	root.AddChild(leftSide)

	// === Right sidebar (synergies) ===
	rightSide := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionEnd,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
			Padding:            &widget.Insets{Top: 40},
		})),
		widget.ContainerOpts.BackgroundImage(nineSlice(uiColorDarkBG)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(6),
			widget.RowLayoutOpts.Padding(&widget.Insets{Left: 8, Right: 8, Top: 10, Bottom: 10}),
		)),
	)

	synTitle := widget.NewLabel(
		widget.LabelOpts.Text("시너지", facePtr(fontSmall), labelColor(uiColorYellow)),
	)
	rightSide.AddChild(synTitle)

	ui.synHumanLabel = widget.NewLabel(
		widget.LabelOpts.Text("인간 0/2", facePtr(fontMedium), labelColor(uiColorAccent)),
	)
	rightSide.AddChild(ui.synHumanLabel)

	ui.synHumanDesc = widget.NewLabel(
		widget.LabelOpts.Text("", facePtr(fontSmall), labelColor(uiColorGrey)),
	)
	rightSide.AddChild(ui.synHumanDesc)

	ui.synElfLabel = widget.NewLabel(
		widget.LabelOpts.Text("엘프 0/2", facePtr(fontMedium), labelColor(uiColorGreen)),
	)
	rightSide.AddChild(ui.synElfLabel)

	ui.synElfDesc = widget.NewLabel(
		widget.LabelOpts.Text("", facePtr(fontSmall), labelColor(uiColorGrey)),
	)
	rightSide.AddChild(ui.synElfDesc)

	root.AddChild(rightSide)

	// === Bottom card hand ===
	bottomBar := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
			StretchHorizontal:  true,
		})),
		widget.ContainerOpts.BackgroundImage(nineSlice(uiColorDarkBG)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(8),
			widget.RowLayoutOpts.Padding(&widget.Insets{Left: 15, Right: 15, Top: 8, Bottom: 8}),
		)),
	)

	ui.cardContainer = bottomBar
	root.AddChild(bottomBar)

	ui.battleUI = &ebitenui.UI{Container: root}
}

// ---- Reward Screen ----

func (ui *UIManager) buildRewardUI() {
	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSlice(uiColorOverlay)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	centerPanel := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(15),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
		)),
	)

	ui.rewardTitle = widget.NewLabel(
		widget.LabelOpts.Text("웨이브 클리어! 카드를 선택하세요:", facePtr(fontLarge), labelColor(uiColorYellow)),
		widget.LabelOpts.TextOpts(widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter)),
	)
	centerPanel.AddChild(ui.rewardTitle)

	ui.rewardContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(15),
		)),
	)
	centerPanel.AddChild(ui.rewardContainer)

	root.AddChild(centerPanel)
	ui.rewardUI = &ebitenui.UI{Container: root}
}

// ---- End Screen (Game Over / Victory) ----

func (ui *UIManager) buildEndUI() {
	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSlice(color.NRGBA{0x08, 0x08, 0x12, 0xF0})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	centerPanel := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(15),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
		)),
	)

	ui.endTitle = widget.NewLabel(
		widget.LabelOpts.Text("게임 오버", facePtr(fontTitle), labelColor(uiColorRed)),
		widget.LabelOpts.TextOpts(widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter)),
	)
	centerPanel.AddChild(ui.endTitle)

	ui.endMessage = widget.NewLabel(
		widget.LabelOpts.Text("", facePtr(fontMedium), labelColor(uiColorWhite)),
		widget.LabelOpts.TextOpts(widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter)),
	)
	centerPanel.AddChild(ui.endMessage)

	spacer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(0, 15)),
	)
	centerPanel.AddChild(spacer)

	restartBtn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage(
			color.NRGBA{0x50, 0x30, 0x10, 0xFF},
			color.NRGBA{0x70, 0x50, 0x20, 0xFF},
			color.NRGBA{0x40, 0x20, 0x08, 0xFF},
		)),
		widget.ButtonOpts.Text("타이틀로 돌아가기", facePtr(fontMedium), &widget.ButtonTextColor{
			Idle: uiColorWhite,
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 30, Right: 30, Top: 10, Bottom: 10}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.onReturnToTitle()
		}),
	)
	centerPanel.AddChild(restartBtn)

	root.AddChild(centerPanel)
	ui.endUI = &ebitenui.UI{Container: root}
}

// ---- Update methods ----

func (ui *UIManager) updateBattleHUD(g *Game) {
	// Wave info
	ui.waveLabel.Label = fmt.Sprintf("웨이브 %d/%d", g.wave+1, g.maxWave)

	// Enemy count
	aliveCount := 0
	for _, e := range g.enemies {
		if !e.Dead && !e.Reached {
			aliveCount++
		}
	}
	aliveCount += len(g.spawnQueue)
	ui.enemyLabel.Label = fmt.Sprintf("적: %d", aliveCount)

	// FPS
	ui.fpsLabel.Label = fmt.Sprintf("FPS: %.0f", ebiten.ActualFPS())

	// HP
	ui.hpLabel.Label = fmt.Sprintf("%d/%d", g.summonerHP, config.SummonerMaxHP)

	// Mana
	ui.manaLabel.Label = fmt.Sprintf("%.0f/%d", g.mana, g.maxMana)

	// Deck/Graveyard
	ui.deckLabel.Label = fmt.Sprintf("덱: %d", len(g.deck))
	ui.graveLabel.Label = fmt.Sprintf("묘지: %d", len(g.graveyard))

	// Speed button text indicator
	if g.gameSpeed == 1 {
		ui.speedBtn1.Text().Label = "1x ●"
		ui.speedBtn2.Text().Label = "2x"
	} else {
		ui.speedBtn1.Text().Label = "1x"
		ui.speedBtn2.Text().Label = "2x ●"
	}

	// Synergies
	humanText := fmt.Sprintf("인간 %d/2", g.humanCount)
	if g.humanSynergy {
		humanText += " 활성"
	}
	ui.synHumanLabel.Label = humanText
	if g.humanSynergy {
		ui.synHumanDesc.Label = "체력 +20%"
	} else {
		ui.synHumanDesc.Label = ""
	}

	elfText := fmt.Sprintf("엘프 %d/2", g.elfCount)
	if g.elfSynergy {
		elfText += " 활성"
	}
	ui.synElfLabel.Label = elfText
	if g.elfSynergy {
		ui.synElfDesc.Label = "공속 +20%"
	} else {
		ui.synElfDesc.Label = ""
	}

	// Only rebuild card buttons when hand state actually changes
	handNames := ""
	for _, c := range g.hand {
		handNames += c.Data.Name
	}
	manaFloor := int(g.mana)
	if len(g.hand) != ui.prevHandLen || handNames != ui.prevHandNames ||
		g.selectedCard != ui.prevSelectedCard || manaFloor != ui.prevManaFloor {
		ui.rebuildCardButtons(g)
		ui.prevHandLen = len(g.hand)
		ui.prevHandNames = handNames
		ui.prevSelectedCard = g.selectedCard
		ui.prevManaFloor = manaFloor
	}
}

func (ui *UIManager) rebuildCardButtons(g *Game) {
	ui.cardContainer.RemoveChildren()
	ui.cardButtons = nil

	for i, card := range g.hand {
		idx := i
		canAfford := g.mana >= float64(card.Data.Cost)

		bgColor := uiColorCardBG
		if !canAfford {
			bgColor = color.NRGBA{0x30, 0x18, 0x18, 0xF0}
		}
		if g.selectedCard == idx {
			bgColor = color.NRGBA{0x20, 0x50, 0x30, 0xFF}
		}

		borderColor := uiColorBorder
		if g.selectedCard == idx {
			borderColor = uiColorGreen
		}

		// Build card text
		nameStr := card.Data.Name
		costStr := fmt.Sprintf("비용: %d", card.Data.Cost)
		if !canAfford {
			costStr += " (부족)"
		}

		var statsStr string
		var rangeStr string
		if card.Data.Type != entity.CardFireball {
			statsStr = fmt.Sprintf("체력:%d 공격:%d", card.Data.HP, card.Data.Atk)
			if card.Data.Range <= 1 {
				rangeStr = "근접"
			} else {
				rangeStr = fmt.Sprintf("사거리:%d", card.Data.Range)
			}
		} else {
			statsStr = "피해:20 광역"
			rangeStr = "타겟 클릭"
		}

		raceName := ""
		switch card.Data.Race {
		case entity.RaceHuman:
			raceName = "[인간]"
		case entity.RaceElf:
			raceName = "[엘프]"
		}

		fullText := nameStr + "\n" + costStr + "\n" + statsStr + "\n" + rangeStr
		if raceName != "" {
			fullText += "\n" + raceName
		}

		cardBtn := widget.NewButton(
			widget.ButtonOpts.Image(&widget.ButtonImage{
				Idle:    borderedNineSlice(bgColor, borderColor, 2),
				Hover:   borderedNineSlice(uiColorCardHover, uiColorAccent, 2),
				Pressed: borderedNineSlice(bgColor, uiColorAccent, 2),
			}),
			widget.ButtonOpts.Text(fullText, facePtr(fontCardInfo), &widget.ButtonTextColor{
				Idle:    uiColorWhite,
				Hover:   uiColorYellow,
				Pressed: uiColorAccent,
			}),
			widget.ButtonOpts.TextPadding(&widget.Insets{Left: 8, Right: 8, Top: 4, Bottom: 4}),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(125, 95)),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				ui.onCardSelect(idx)
			}),
		)

		ui.cardButtons = append(ui.cardButtons, cardBtn)
		ui.cardContainer.AddChild(cardBtn)
	}
}

func (ui *UIManager) updateReward(g *Game) {
	ui.rewardTitle.Label = fmt.Sprintf("웨이브 %d 클리어! 카드를 선택하세요:", g.wave+1)

	ui.rewardContainer.RemoveChildren()
	ui.rewardButtons = nil

	for i, card := range g.rewardCards {
		idx := i

		var statsText string
		if card.Data.Type != entity.CardFireball {
			speedText := "보통"
			if card.Data.AtkSpeed < 50 {
				speedText = "빠름"
			} else if card.Data.AtkSpeed > 70 {
				speedText = "느림"
			}
			rangeText := "근접"
			if card.Data.Range > 1 {
				rangeText = fmt.Sprintf("%d칸", card.Data.Range)
			}
			statsText = fmt.Sprintf("%s\n비용: %d\n체력: %d\n공격: %d\n속도: %s\n사거리: %s",
				card.Data.Name, card.Data.Cost, card.Data.HP, card.Data.Atk, speedText, rangeText)
		} else {
			statsText = fmt.Sprintf("%s\n비용: %d\n마법 카드\n피해: 20 (광역)\n주변 모든 적에게\n피해를 줍니다",
				card.Data.Name, card.Data.Cost)
		}

		raceName := ""
		switch card.Data.Race {
		case entity.RaceHuman:
			raceName = "[인간]"
		case entity.RaceElf:
			raceName = "[엘프]"
		}
		if raceName != "" {
			statsText += "\n" + raceName
		}

		rewardBtn := widget.NewButton(
			widget.ButtonOpts.Image(buttonImage(
				uiColorCardBG,
				uiColorCardHover,
				color.NRGBA{0x20, 0x40, 0x30, 0xFF},
			)),
			widget.ButtonOpts.Text(statsText+"\n\n클릭하여 추가", facePtr(fontCardInfo), &widget.ButtonTextColor{
				Idle:    uiColorWhite,
				Hover:   uiColorYellow,
				Pressed: uiColorAccent,
			}),
			widget.ButtonOpts.TextPadding(&widget.Insets{Left: 15, Right: 15, Top: 10, Bottom: 10}),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(170, 230)),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				ui.onRewardSelect(idx)
			}),
		)

		ui.rewardButtons = append(ui.rewardButtons, rewardBtn)
		ui.rewardContainer.AddChild(rewardBtn)
	}
}

func (ui *UIManager) updateEndScreen(g *Game) {
	if g.state == entity.StateGameOver {
		ui.endTitle.Label = "게임 오버"
		ui.endMessage.Label = fmt.Sprintf("웨이브 %d/%d에서 패배했습니다", g.wave+1, g.maxWave)
	} else {
		ui.endTitle.Label = "승리!"
		ui.endMessage.Label = "모든 10 웨이브를 방어했습니다!\n소환사가 안전합니다!"
	}
}
