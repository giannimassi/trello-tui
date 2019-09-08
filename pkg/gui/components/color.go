package components

import (
	"github.com/fatih/color"
)

type ElementClass int

const (
	DefaultClass ElementClass = iota
	BoardDescriptionClass
	CardTitleClass
)

type ColorSettings map[ElementClass]ColorSetting

type ColorSetting struct {
	normal   *color.Color
	selected *color.Color
}

var DefaultColorSettings = ColorSettings{
	DefaultClass: {
		normal:   color.New(color.FgWhite, color.BgBlack),
		selected: color.New(color.FgGreen, color.BgBlack),
	},

	CardTitleClass: {
		normal:   color.New(color.FgWhite, color.BgBlack),
		selected: color.New(color.FgBlack, color.BgYellow),
	},
}
