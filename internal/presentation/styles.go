package presentation

import "charm.land/lipgloss/v2"

// The palette. Every style in this package is built from these values, so a colour
// is decided once and read everywhere.
//
// colorDetail and colorLifecycleDetail are two separate greys that sit next to each
// other in the same output. Whether they should become one is a design question and
// is tracked separately; the values here are the ones that ship today.
const (
	colorAccent          = "#2ABEF6"
	colorProject         = "#F5A524"
	colorSuccess         = "#38D17A"
	colorFailure         = "#F05D5E"
	colorBright          = "#D8E0EA"
	colorBody            = "#9AA6B5"
	colorDetail          = "#7D8794"
	colorLifecycleDetail = "#8793A2"
)

// Styles of the line-based output.
var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorAccent))
	detailStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(colorDetail))
	projectStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorProject))
	errorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorFailure))
	serviceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorBright))
	bodyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(colorBody))
)

// Styles of the step labels. These use the terminal's own first eight colours rather
// than a hex value, so a step keeps the colour the user configured for their shell.
var (
	stepInfoStyle    = lipgloss.NewStyle().Faint(true)
	stepWorkingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	stepSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	stepWarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	stepFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

// Styles of the inline lifecycle view and of the spinner, which shows the same
// accent whilst something is in progress.
var (
	lifecycleHeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorBright))
	lifecycleDetailStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(colorLifecycleDetail))
	lifecycleWorkingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
	lifecycleSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSuccess))
	lifecycleFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorFailure))
	spinnerStyle          = lifecycleWorkingStyle
)

// Geometry of the inline lifecycle view, in terminal cells.
const (
	lifecycleHorizontalInset      = 2
	lifecycleNameColumnWidth      = 18
	lifecycleStateColumnWidth     = 24
	lifecycleColumnGap            = 2
	lifecycleWideFixedColumns     = lifecycleNameColumnWidth + lifecycleColumnGap + lifecycleStateColumnWidth + lifecycleColumnGap
	lifecycleGroupedIndent        = 4
	lifecycleGroupedServiceWidth  = 26
	lifecycleReassignDividerWidth = 39
)
