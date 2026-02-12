package game

import (
	"bytes"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	fontSource *text.GoTextFaceSource

	// Stored as text.Face so we can take &fontXxx for ebitenui (which wants *text.Face)
	fontSmall    text.Face
	fontMedium   text.Face
	fontLarge    text.Face
	fontTitle    text.Face
	fontCardInfo text.Face

	// Also keep GoTextFace versions for direct text.Draw calls
	goFontSmall  *text.GoTextFace
	goFontMedium *text.GoTextFace
	goFontLarge  *text.GoTextFace
)

func initFonts() {
	// Try embedded font first, then system fonts as fallback
	fontData, err := os.ReadFile("assets/fonts/NotoSansKR.ttf")
	if err != nil {
		// Fallback: try macOS system Korean font
		fontData, err = os.ReadFile("/System/Library/Fonts/Supplemental/AppleGothic.ttf")
		if err != nil {
			panic("한국어 폰트를 찾을 수 없습니다. assets/fonts/NotoSansKR.ttf 파일이 필요합니다.")
		}
	}

	fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		panic("폰트 로드 실패: " + err.Error())
	}

	goFontSmall = &text.GoTextFace{Source: fontSource, Size: 11}
	goFontMedium = &text.GoTextFace{Source: fontSource, Size: 14}
	goFontLarge = &text.GoTextFace{Source: fontSource, Size: 18}

	fontSmall = goFontSmall
	fontCardInfo = &text.GoTextFace{Source: fontSource, Size: 12}
	fontMedium = goFontMedium
	fontLarge = goFontLarge
	fontTitle = &text.GoTextFace{Source: fontSource, Size: 28}
}
