package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"sidus.io/boogrocha/internal/cli/commands"
)

const ApplicationName = "BooGroCha"
const Version = "0.1"

func init() {

	BgcCmd.AddCommand(commands.BookCmd(getBookingService, getRankingService))
	BgcCmd.AddCommand(commands.ConfigCmd(getSavePassword))
	BgcCmd.AddCommand(commands.DeleteCmd(getBookingService))
	BgcCmd.AddCommand(commands.ListCmd(getBookingService))
	BgcCmd.AddCommand(commands.VersionCmd(ApplicationName, Version))
	BgcCmd.AddCommand(commands.CompletionCmd)

	loadFlags()

	err := loadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %s", err)
		os.Exit(1)
	}

	err = bindFlags()
	if err != nil {
		fmt.Printf("Failed to bind flags: %v \n", err)
		os.Exit(1)
	}
}

var BgcCmd = &cobra.Command{
	Use:   "bgc",
	Short: "Manage your group room bookings at Chalmers",
	Long: `A lightweight, easy to use application for managing your
		   group room bookings at Chalmers University of Technology`,
	Run: nil,
}
