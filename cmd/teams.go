package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/nudgen/nudgen-cli/internal/api"
	"github.com/nudgen/nudgen-cli/internal/config"
	"github.com/spf13/cobra"
)

var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Manage your Nudgen teams",
	Long: `Multi-team support for shared lead ops and campaigns.
From here you can create new teams, switch your active work context,
or manage your existing team memberships.`,
	Example: `  nudgen teams list
  nudgen teams switch cm..id
  nudgen teams create`,
}

var teamListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all teams you belong to",
	Long:  `Retrieve and display your roles and memberships in all Nudgen teams.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		teams, err := client.GetTeams()
		if err != nil {
			log.Fatalf("Error fetching teams: %v", err)
		}

		if UseJSON {
			json.NewEncoder(os.Stdout).Encode(teams)
			return
		}

		if len(teams) == 0 {
			fmt.Println("You don't belong to any teams yet.")
			return
		}

		fmt.Println("--- Your Teams ---")
		for _, t := range teams {
			fmt.Printf("[%s] %s (Role: %s)\n", t.ID, t.Name, t.Role)
		}
	},
}

var teamCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the currently active team context",
	Run: func(cmd *cobra.Command, args []string) {
		config.PrintActiveTeam()
	},
}

var teamSwitchCmd = &cobra.Command{
	Use:   "switch [team-id]",
	Short: "Switch the active team context",
	Long: `Updates the default team for the CLI and backend session. 
All subsequent commands (e.g., stats, campaigns, contacts) will be relative to this team.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		teamID := args[0]
		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		err := client.SwitchTeam(teamID)
		if err != nil {
			log.Fatalf("Error switching team: %v", err)
		}

		// Sync local config
		teams, _ := client.GetTeams()
		teamName := teamID
		for _, t := range teams {
			if t.ID == teamID {
				teamName = t.Name
				break
			}
		}
		config.SaveActiveTeam(teamID, teamName)

		fmt.Printf("Successfully switched active team context to: %s (%s)\n", teamName, teamID)
	},
}

var teamDeleteCmd = &cobra.Command{
	Use:   "delete [team-id]",
	Short: "Delete a team",
	Long:  `Permanently removes a team and all its data. Must be called by the team owner.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		id := args[0]
		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		err := client.DeleteTeam(id)
		if err != nil {
			log.Fatalf("Error deleting team: %v", err)
		}

		fmt.Printf("Successfully deleted team: %s\n", id)
	},
}

var teamCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new team interactively",
	Long:  `Launch an interactive wizard to name and initialize a new Nudgen team.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		var name string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Team Name").
					Placeholder("e.g. My New Team").
					Value(&name).
					Validate(func(s string) error {
						if len(s) < 2 {
							return fmt.Errorf("Name too short")
						}
						return nil
					}),
			),
		)

		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		newTeam, err := client.CreateTeam(name)
		if err != nil {
			log.Fatalf("Error creating team: %v", err)
		}

		fmt.Printf("Successfully created team: %s (ID: %s)\n", newTeam.Name, newTeam.ID)
	},
}

var teamUpdateCmd = &cobra.Command{
	Use:   "update [team-id]",
	Short: "Update a team's name",
	Long:  `Launches a wizard to rename an existing Nudgen team.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Team update via CLI is matching name only for now.")
	},
}

func init() {
	rootCmd.AddCommand(teamsCmd)
	teamsCmd.AddCommand(teamListCmd)
	teamsCmd.AddCommand(teamSwitchCmd)
	teamsCmd.AddCommand(teamCurrentCmd)
	teamsCmd.AddCommand(teamCreateCmd)
	teamsCmd.AddCommand(teamUpdateCmd)
	teamsCmd.AddCommand(teamDeleteCmd)
}
