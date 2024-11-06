package style

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	NormalBg  = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
	Primary   = lipgloss.Color("#D4C985")
	Cream     = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
	Highlight = lipgloss.Color("#C5B141")
	Info      = lipgloss.Color("#89CFF0")
	Success   = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	Error     = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
)

var (
	TitleWithBgStyle  = lipgloss.NewStyle().Background(Primary).Foreground(Cream).Bold(true).Padding(0, 1, 0)
	TitleStyle        = lipgloss.NewStyle().Foreground(Primary).Bold(true)
	ErrorStyle        = lipgloss.NewStyle().Foreground(Error).Bold(true).Padding(0, 1, 0)
	LogStyle          = lipgloss.NewStyle().Foreground(Cream)
	FocusedStyle      = lipgloss.NewStyle().Foreground(Primary).Bold(true)
	SelectedItemStyle = lipgloss.NewStyle().Foreground(Highlight).Bold(true)
	CheckStyle        = lipgloss.NewStyle().Foreground(Success).Bold(true).Padding(0, 1, 0)
	CrossMarkStyle    = lipgloss.NewStyle().Foreground(Error).Bold(true).Padding(0, 1, 0)
	TagsStyle         = lipgloss.NewStyle().Foreground(Info)
	DescriptionStyle  = lipgloss.NewStyle().Foreground(Cream)
)

func Theme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
	t.Focused.Title = t.Focused.Title.Foreground(Primary).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(Primary).Bold(true).MarginBottom(1)
	t.Focused.Directory = t.Focused.Directory.Foreground(Primary)
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"})
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(Error)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(Error)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(Highlight)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(Highlight)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(Highlight)
	t.Focused.Option = t.Focused.Option.Foreground(NormalBg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(Highlight)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(Info)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#02CF92", Dark: "#02A877"}).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"}).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(NormalBg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(Cream).Background(Highlight)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(NormalBg).Background(lipgloss.AdaptiveColor{Light: "252", Dark: "237"})

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(Info)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "238"})
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(Highlight)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return t
}
