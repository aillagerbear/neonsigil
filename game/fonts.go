package game

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed fonts/Pretendard-Regular.otf
var fontRegularData []byte

//go:embed fonts/Pretendard-Bold.otf
var fontBoldData []byte

var (
	fontRegularSource *text.GoTextFaceSource
	fontBoldSource    *text.GoTextFaceSource

	// Stored as text.Face so we can take &fontXxx for ebitenui (which wants *text.Face)
	fontSmall    text.Face
	fontMedium   text.Face
	fontLarge    text.Face
	fontTitle    text.Face
	fontCardInfo text.Face
	fontCardName text.Face
	fontCardCost text.Face
	fontBarValue text.Face

	// Also keep GoTextFace versions for direct text.Draw calls
	goFontSmall  *text.GoTextFace
	goFontMedium *text.GoTextFace
	goFontLarge  *text.GoTextFace
)

func initFonts() {
	fontRegularSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontRegularData))
	if err != nil {
		panic("폰트 로드 실패: " + err.Error())
	}
	fontBoldSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontBoldData))
	if err != nil {
		panic("폰트 로드 실패: " + err.Error())
	}

	goFontSmall = &text.GoTextFace{Source: fontRegularSource, Size: 12}
	goFontMedium = &text.GoTextFace{Source: fontRegularSource, Size: 15}
	goFontLarge = &text.GoTextFace{Source: fontBoldSource, Size: 20}

	fontSmall = goFontSmall
	fontCardInfo = &text.GoTextFace{Source: fontRegularSource, Size: 12}
	fontCardName = &text.GoTextFace{Source: fontBoldSource, Size: 14}
	fontCardCost = &text.GoTextFace{Source: fontBoldSource, Size: 16}
	fontBarValue = &text.GoTextFace{Source: fontBoldSource, Size: 11}
	fontMedium = goFontMedium
	fontLarge = goFontLarge
	fontTitle = &text.GoTextFace{Source: fontBoldSource, Size: 34}
}
