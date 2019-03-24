// +build !windows

package promptui

var (
	bold       = Styler(FGBold)
	faint      = Styler(FGFaint)
	underlined = Styler(FGUnderline)
	blue       = Styler(FGBlue)
)

// Icons used for displaying prompts or status
var (
	IconInitial = Styler(FGBlue)("?")
	IconGood    = Styler(FGGreen)("✔")
	IconQuest   = Styler(FGBlue)("✔")
	IconWarn    = Styler(FGYellow)("⚠")
	IconBad     = Styler(FGRed)("✗")
)

var red = Styler(FGBold, FGRed)
