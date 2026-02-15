package game

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed fonts/NotoSansKR.ttf
var fontTTF []byte

var (
	fontSource *text.GoTextFaceSource

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
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontTTF))
	if err != nil {
		panic("폰트 로드 실패: " + err.Error())
	}

	goFontSmall = &text.GoTextFace{Source: fontSource, Size: 11}
	goFontMedium = &text.GoTextFace{Source: fontSource, Size: 14}
	goFontLarge = &text.GoTextFace{Source: fontSource, Size: 18}

	fontSmall = goFontSmall
	fontCardInfo = &text.GoTextFace{Source: fontSource, Size: 12}
	fontCardName = &text.GoTextFace{Source: fontSource, Size: 13}
	fontCardCost = &text.GoTextFace{Source: fontSource, Size: 16}
	fontBarValue = &text.GoTextFace{Source: fontSource, Size: 10}
	fontMedium = goFontMedium
	fontLarge = goFontLarge
	fontTitle = &text.GoTextFace{Source: fontSource, Size: 28}
}
