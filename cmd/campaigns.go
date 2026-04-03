package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nudgen/nudgen-cli/internal/api"
	"github.com/nudgen/nudgen-cli/internal/config"
	"github.com/spf13/cobra"
)

var campaignsCmd = &cobra.Command{
	Use:   "campaigns",
	Short: "Manage and run campaigns",
	Long: `Campaigns are outreach efforts using Email or LinkedIn.
From here you can list your active campaigns, track their status,
create new ones, or delete those no longer needed.`,
	Example: `  nudgen campaigns list
  nudgen campaigns delete cm...id`,
}

var campaignsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your active campaigns",
	Long:  `Fetch and display all campaigns created by the active team.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		campaigns, err := client.GetCampaigns()
		if err != nil {
			log.Fatalf("Error fetching campaigns: %v", err)
		}

		if UseJSON {
			json.NewEncoder(os.Stdout).Encode(campaigns)
			return
		}

		if len(campaigns) == 0 {
			fmt.Println("No active campaigns found.")
			return
		}

		fmt.Println("--- Active Campaigns ---")
		for _, c := range campaigns {
			fmt.Printf("[%s] %s (%s)\n", c.ID, c.Name, c.Status)
		}
	},
}

var campaignsDeleteCmd = &cobra.Command{
	Use:   "delete [campaign-id]",
	Short: "Delete a campaign",
	Long:  `Permanently remove a campaign and its associated data. A campaign cannot be deleted while it is currently sending.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		id := args[0]
		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		err := client.DeleteCampaign(id)
		if err != nil {
			log.Fatalf("Error deleting campaign: %v", err)
		}

		fmt.Printf("Successfully deleted campaign: %s\n", id)
	},
}

var campaignsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new campaign",
	Long:  `Launches a wizard to create a new campaign from scratch.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[Not Implemented] Interactive Form for Campaign Generation using Bubbletea.")
	},
}

var campaignsUpdateCmd = &cobra.Command{
	Use:   "update [campaign-id]",
	Short: "Update a campaign settings",
	Long:  `Launches a wizard to update an existing campaign's settings or content.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[Not Implemented] Interactive Form for Campaign Update using Bubbletea.")
	},
}

func init() {
	rootCmd.AddCommand(campaignsCmd)
	campaignsCmd.AddCommand(campaignsListCmd)
	campaignsCmd.AddCommand(campaignsCreateCmd)
	campaignsCmd.AddCommand(campaignsUpdateCmd)
	campaignsCmd.AddCommand(campaignsDeleteCmd)
}
