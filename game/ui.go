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

	// End screen
	endTitle   *widget.Label
	endMessage *widget.Label

	// Callbacks
	onStartGame     func()
	onSpeedChange   func(speed int)
	onReturnToTitle func()
}

func newUIManager(g *Game) *UIManager {
	ui := &UIManager{}
	ui.onStartGame = func() { g.initBattle() }
	ui.onSpeedChange = func(speed int) { g.gameSpeed = speed }
	ui.onReturnToTitle = func() { g.state = entity.StateTitle }

	ui.buildTitleUI()
	ui.buildBattleUI()
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
	uiColorDarkBG = color.NRGBA{0x0f, 0x16, 0x22, 0xD8}
	uiColorBorder = color.NRGBA{0x39, 0x52, 0x72, 0xFF}
	uiColorAccent = color.NRGBA{0x36, 0xd1, 0xdc, 0xFF}
	uiColorGreen  = color.NRGBA{0x29, 0xd4, 0xa1, 0xFF}
	uiColorRed    = color.NRGBA{0xff, 0x5d, 0x73, 0xFF}
	uiColorYellow = color.NRGBA{0xff, 0xd5, 0x66, 0xFF}
	uiColorOrange = color.NRGBA{0xff, 0x9d, 0x54, 0xFF}
	uiColorWhite  = color.NRGBA{0xf3, 0xf8, 0xff, 0xFF}
	uiColorGrey   = color.NRGBA{0x8f, 0xa1, 0xbb, 0xFF}
)

func labelColor(c color.NRGBA) *widget.LabelColor {
	return &widget.LabelColor{Idle: c, Disabled: uiColorGrey}
}

func buttonImage(idle, hover, pressed color.NRGBA) *widget.ButtonImage {
	return &widget.ButtonImage{
		Idle:    borderedNineSlice(idle, uiColorBorder, 1),
		Hover:   borderedNineSlice(hover, uiColorAccent, 1),
		Pressed: borderedNineSlice(pressed, color.NRGBA{0x23, 0xaa, 0xb5, 0xFF}, 1),
	}
}

// facePtr returns a *text.Face from a text.Face for ebitenui
func facePtr(f text.Face) *text.Face {
	return &f
}

// ---- Title Screen ----

func (ui *UIManager) buildTitleUI() {
	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSlice(color.NRGBA{0x00, 0x00, 0x00, 0x00})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	centerPanel := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
			Padding:            &widget.Insets{Top: 120},
		})),
		widget.ContainerOpts.BackgroundImage(borderedNineSlice(
			color.NRGBA{0x0d, 0x14, 0x20, 0xD8},
			color.NRGBA{0x41, 0x67, 0x8d, 0xF0},
			1,
		)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
			widget.RowLayoutOpts.Padding(&widget.Insets{Left: 26, Right: 26, Top: 24, Bottom: 24}),
		)),
	)

	statusLabel := widget.NewLabel(
		widget.LabelOpts.Text("STRATEGY DEFENSE", facePtr(fontSmall), labelColor(color.NRGBA{0x9f, 0xc8, 0xe8, 0xFF})),
		widget.LabelOpts.TextOpts(widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter)),
	)
	centerPanel.AddChild(statusLabel)

	subLabel := widget.NewLabel(
		widget.LabelOpts.Text("자동전투 덱빌딩 디펜스", facePtr(fontMedium), labelColor(uiColorWhite)),
		widget.LabelOpts.TextOpts(widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter)),
	)
	centerPanel.AddChild(subLabel)

	spacer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(0, 10)),
	)
	centerPanel.AddChild(spacer)

	startBtn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage(
			color.NRGBA{0x1a, 0x97, 0xa8, 0xFF},
			color.NRGBA{0x25, 0xb9, 0xcd, 0xFF},
			color.NRGBA{0x17, 0x7f, 0x8f, 0xFF},
		)),
		widget.ButtonOpts.Text("게임 시작", facePtr(fontLarge), &widget.ButtonTextColor{
			Idle: uiColorWhite,
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 54, Right: 54, Top: 12, Bottom: 12}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.onStartGame()
		}),
	)
	centerPanel.AddChild(startBtn)

	spacer2 := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(0, 8)),
	)
	centerPanel.AddChild(spacer2)

	instructions := []string{
		"카드를 선택해 전장에 배치하세요.",
		"유닛은 사거리 내 적을 자동으로 공격합니다.",
		"10웨이브를 막아내면 승리합니다.",
	}
	for _, line := range instructions {
		l := widget.NewLabel(
			widget.LabelOpts.Text(line, facePtr(fontSmall), labelColor(uiColorGrey)),
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
		widget.ContainerOpts.BackgroundImage(borderedNineSlice(uiColorDarkBG, uiColorBorder, 1)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(&widget.Insets{Left: 18, Right: 18, Top: 9, Bottom: 9}),
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
			color.NRGBA{0x1a, 0x97, 0xa8, 0xFF},
			color.NRGBA{0x25, 0xb9, 0xcd, 0xFF},
			color.NRGBA{0x17, 0x7f, 0x8f, 0xFF},
		)),
		widget.ButtonOpts.Text("1x", facePtr(fontSmall), &widget.ButtonTextColor{Idle: uiColorWhite}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 12, Right: 12, Top: 3, Bottom: 3}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ui.onSpeedChange(1)
		}),
	)
	speedContainer.AddChild(ui.speedBtn1)

	ui.speedBtn2 = widget.NewButton(
		widget.ButtonOpts.Image(buttonImage(
			color.NRGBA{0x1f, 0x2e, 0x42, 0xFF},
			color.NRGBA{0x25, 0xb9, 0xcd, 0xFF},
			color.NRGBA{0x17, 0x7f, 0x8f, 0xFF},
		)),
		widget.ButtonOpts.Text("2x", facePtr(fontSmall), &widget.ButtonTextColor{Idle: uiColorWhite}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 12, Right: 12, Top: 3, Bottom: 3}),
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
			Padding:            &widget.Insets{Top: 40, Left: 6},
		})),
		widget.ContainerOpts.BackgroundImage(borderedNineSlice(uiColorDarkBG, uiColorBorder, 1)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(6),
			widget.RowLayoutOpts.Padding(&widget.Insets{Left: 11, Right: 11, Top: 12, Bottom: 12}),
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
			Padding:            &widget.Insets{Top: 40, Right: 6},
		})),
		widget.ContainerOpts.BackgroundImage(borderedNineSlice(uiColorDarkBG, uiColorBorder, 1)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(6),
			widget.RowLayoutOpts.Padding(&widget.Insets{Left: 11, Right: 11, Top: 12, Bottom: 12}),
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

	ui.battleUI = &ebitenui.UI{Container: root}
}

// ---- End Screen (Game Over / Victory) ----

func (ui *UIManager) buildEndUI() {
	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSlice(color.NRGBA{0x04, 0x0b, 0x14, 0xE6})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	centerPanel := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.ContainerOpts.BackgroundImage(borderedNineSlice(uiColorDarkBG, uiColorBorder, 1)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(12),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(24)),
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
			color.NRGBA{0x1f, 0x2e, 0x42, 0xFF},
			color.NRGBA{0x2f, 0x4c, 0x6d, 0xFF},
			color.NRGBA{0x17, 0x24, 0x34, 0xFF},
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
