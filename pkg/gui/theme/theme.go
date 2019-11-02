package theme

import (
	"github.com/fatih/color"
)

func DefaultTheme() *Theme {
	return &Theme{
		Colors:     DefaultColorSettings,
		ListWidths: DefaultListWidths,
	}
}

type Theme struct {
	Colors
	ListWidths
}

func (t *Theme) ListsPerPage(maxW int) int {
	return DefaultListWidths.Match(maxW)
}

func (t *Theme) Color(class string, isSelected bool) *color.Color {
	setting, found := t.Colors[class]
	if !found {
		setting = DefaultColorSettings[DefaultClass]
	}
	if isSelected {
		return setting.selected
	}
	return setting.normal
}

type Colors map[string]ColorSetting

type ListWidths [][2]int

const (
	DefaultClass              = ""
	BoardDescriptionClass     = "board-description"
	CardTitleClass            = "selected-card-title"
	PopupCardTitle            = "popup-title"
	PopupCardDescriptionClass = "popup-description"
)

func (l *ListWidths) Match(w int) int {
	for _, v := range DefaultListWidths {
		if w <= v[0] {
			return v[1]
		}
	}
	return 1
}

type ColorSetting struct {
	normal   *color.Color
	selected *color.Color
}

var DefaultColorSettings = Colors{
	DefaultClass: {
		normal:   color.New(color.FgWhite, color.BgBlack),
		selected: color.New(color.FgGreen, color.BgBlack),
	},

	CardTitleClass: {
		normal:   color.New(color.FgWhite, color.BgBlack),
		selected: color.New(color.FgBlack, color.BgYellow),
	},

	PopupCardDescriptionClass: {
		normal: color.New(color.FgBlue, color.BgBlack),
	},

	PopupCardTitle: {
		normal:   color.New(color.FgGreen, color.BgBlack),
		selected: color.New(color.FgGreen, color.BgBlack),
	},
}

var DefaultListWidths = ListWidths{
	{50, 1},
	{100, 2},
	{150, 3},
	{200, 4},
	{250, 5},
	{300, 6},
	{450, 7},
	{100000000, 8},
}
