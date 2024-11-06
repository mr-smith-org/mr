package style

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/mr-smith/mr/internal/debug"
)

func PrintStyles() {
	fmt.Println(TitleWithBgStyle.Render("Title") + "\n")
	fmt.Println(ErrorStyle.Render("Error") + "\n")
	fmt.Println(LogStyle.Render("Log") + "\n")
	fmt.Println(FocusedStyle.Render("Focused") + "\n")
	fmt.Println(SelectedItemStyle.Render("Selected") + "\n")
	fmt.Println(CheckStyle.Render("Check") + "\n")
	fmt.Println(CrossMarkStyle.Render("Cross") + "\n")
	fmt.Println(TagsStyle.Render("Tags") + "\n")
	fmt.Println(DescriptionStyle.Render("Description") + "\n")
}

func TitlePrint(text string, withBG bool) {
	if withBG {
		fmt.Println(TitleWithBgStyle.Render(text) + "\n")
		return
	}
	fmt.Println(TitleStyle.Render(text) + "\n")
}

func LogPrint(text string) {
	fmt.Println(LogStyle.Render(text) + "\n")
}

func CheckMarkPrint(text string) {
	fmt.Println(CheckStyle.Render("✔") + text)
}

func CrossMarkPrint(text string) {
	fmt.Println(CrossMarkStyle.Render("✖") + text)
}

func ErrorPrint(text string) {
	fmt.Println(ErrorStyle.Render(text) + "\n")
}

func DebugPrint(header, text string) {
	fmt.Println()
	if debug.Debug {
		color.New(color.FgBlack, color.BgYellow).Println(" - " + header + " - ")
		color.Gray.Println(text)
	}
}
