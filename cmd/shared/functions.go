package shared

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh/spinner"
	"github.com/mr-smith/mr/pkg/style"
)

func RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ReadFileFromURL(url string) (string, error) {
	var bodyBytes []byte
	err := spinner.New().
		Title("Downloading variables file").
		Action(func() {
			resp, err := http.Get(url)
			if err != nil {
				style.ErrorPrint("downloading variables file error: " + err.Error())
				os.Exit(1)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				style.ErrorPrint(fmt.Sprintf("bad status: %s", resp.Status))
				os.Exit(1)
			}

			bodyBytes, err = io.ReadAll(resp.Body)
			if err != nil {
				style.ErrorPrint("reading file error: " + err.Error())
				os.Exit(1)
			}
		}).Run()
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
