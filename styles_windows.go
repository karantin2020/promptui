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
	IconGood    = Styler(FGGreen)("v")
	IconQuest   = Styler(FGBlue)("v")
	IconWarn    = Styler(FGYellow)("!")
	IconBad     = Styler(FGRed)("x")
)

var red = Styler(FGBold, FGRed)
