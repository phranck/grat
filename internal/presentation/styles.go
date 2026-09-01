package presentation

import "charm.land/lipgloss/v2"

// The palette. Every style in this package is built from these values, so a colour
// is decided once and read everywhere.
const (
	colorAccent  = "#2ABEF6"
	colorProject = "#F5A524"
	colorSuccess = "#38D17A"
	colorFailure = "#F05D5E"
	colorBright  = "#D8E0EA"
	colorBody    = "#9AA6B5"

	// One grey for everything subordinate, whichever view prints it. There were
	// two, a few percent apart, and a heading of the line-based output sitting
	// beside the footer of the lifecycle view read as an inconsistency rather
	// than as a distinction.
	colorDetail = "#8793A2"
)

// Styles of the line-based output.
var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorAccent))
	// One style, read by both views. Two identical definitions are how the two
	// greys drifted apart in the first place, so there is only one to change.
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
