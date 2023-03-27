package generate

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/exec"
	"runtime"
	"time"
)

var (
	backend string
)

// generateCmd represents the auth command
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Key for your chosen backend (opens browser)",
	Long:  `Opens your browser to generate a key for your chosen backend.`,
	Run: func(cmd *cobra.Command, args []string) {

		backendType := viper.GetString("backend_type")
		if backendType == "" {
			// Set the default backend
			backend = "openai"
		}
		// override the default backend if a flag is provided
		if backend != "" {
			backendType = backend
		}
		fmt.Println("")
		color.Green("Opening: https://beta.openai.com/account/api-keys to generate a key for %s", backendType)
		color.Green("Please copy the generated key and run `k8sgpt auth` to add it to your config file")
		fmt.Println("")
		time.Sleep(5 * time.Second)
		openbrowser("https://beta.openai.com/account/api-keys")
	},
}

func init() {
	// add flag for backend
	GenerateCmd.Flags().StringVarP(&backend, "backend", "b", "openai", "Backend AI provider")
}

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
	}
}
