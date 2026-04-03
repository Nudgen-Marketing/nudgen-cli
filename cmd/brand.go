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

var brandCmd = &cobra.Command{
	Use:   "brand",
	Short: "Manage and update your brand identity",
	Long: `Brand identities define how the AI Agent represents your company.
You can create multiple voices (e.g., Professional, Friendly) and set a default one
to ensure consistent communication across all campaigns.`,
	Example: `  nudgen brand list
  nudgen brand create
  nudgen brand update cmm...id
  nudgen brand delete cmm...id`,
}

var brandListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all brand identities for the active team",
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		configs, err := client.GetBrandConfigs()
		if err != nil {
			log.Fatalf("Error fetching brand identities: %v", err)
		}

		if UseJSON {
			json.NewEncoder(os.Stdout).Encode(configs)
			return
		}

		if len(configs) == 0 {
			fmt.Println("No brand identities found for the current team.")
			return
		}

		fmt.Println("--- Brand Identities ---")
		for _, c := range configs {
			defaultMark := ""
			if c.IsDefault {
				defaultMark = " (DEFAULT)"
			}
			fmt.Printf("[%s] %s - Tone: %s%s\n", c.ID, c.Name, c.Tone, defaultMark)
		}
	},
}

var brandDeleteCmd = &cobra.Command{
	Use:   "delete [brand-id]",
	Short: "Delete a brand identity",
	Long:  `Permanently remove a brand identity. You cannot delete the default identity.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		id := args[0]
		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		err := client.DeleteBrandConfig(id)
		if err != nil {
			log.Fatalf("Error deleting brand: %v", err)
		}

		fmt.Printf("Successfully deleted brand identity: %s\n", id)
	},
}

var brandCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new brand identity interactively",
	Long:  `Launch an interactive wizard to define name, tone, custom instructions, and default status.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		var (
			name      string
			tone      string = "professional"
			isDefault bool   = false
			prompt    string
		)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Brand Name").
					Placeholder("e.g. My SaaS Company").
					Value(&name).
					Validate(func(s string) error {
						if len(s) < 2 {
							return fmt.Errorf("Name too short")
						}
						return nil
					}),
				huh.NewSelect[string]().
					Title("Tone of Voice").
					Options(
						huh.NewOption("Professional", "professional"),
						huh.NewOption("Friendly", "friendly"),
						huh.NewOption("Playful", "playful"),
						huh.NewOption("Urgent", "urgent"),
					).
					Value(&tone),
				huh.NewText().
					Title("Custom Instructions (AI Prompt)").
					Placeholder("Tell the agent how to write...").
					Value(&prompt),
				huh.NewConfirm().
					Title("Make this the default brand identity?").
					Value(&isDefault),
			),
		)

		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		body := map[string]interface{}{
			"name":         name,
			"tone":         tone,
			"customPrompt": prompt,
			"isDefault":    isDefault,
		}

		newConfig, err := client.CreateBrandConfig(body)
		if err != nil {
			log.Fatalf("Error creating brand: %v", err)
		}

		fmt.Printf("Successfully created brand identity: %s (ID: %s)\n", newConfig.Name, newConfig.ID)
	},
}

var brandUpdateCmd = &cobra.Command{
	Use:   "update [brand-id]",
	Short: "Update a brand identity settings",
	Long:  `Launch an interactive wizard to update settings for an existing brand identity.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		id := args[0]
		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		var (
			name      string
			tone      string
			isDefault bool
			prompt    string
		)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Brand Name (Leave empty to keep current)").
					Value(&name),
				huh.NewSelect[string]().
					Title("Tone of Voice").
					Options(
						huh.NewOption("Professional", "professional"),
						huh.NewOption("Friendly", "friendly"),
						huh.NewOption("Playful", "playful"),
						huh.NewOption("Urgent", "urgent"),
					).
					Value(&tone),
				huh.NewText().
					Title("Custom Instructions (AI Prompt)").
					Value(&prompt),
				huh.NewConfirm().
					Title("Make this the default?").
					Value(&isDefault),
			),
		)

		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		body := map[string]interface{}{}
		if name != "" {
			body["name"] = name
		}
		if tone != "" {
			body["tone"] = tone
		}
		if prompt != "" {
			body["customPrompt"] = prompt
		}
		body["isDefault"] = isDefault

		updated, err := client.UpdateBrandConfig(id, body)
		if err != nil {
			log.Fatalf("Error updating brand: %v", err)
		}

		fmt.Printf("Successfully updated brand identity: %s\n", updated.ID)
	},
}

func init() {
	rootCmd.AddCommand(brandCmd)
	brandCmd.AddCommand(brandListCmd)
	brandCmd.AddCommand(brandCreateCmd)
	brandCmd.AddCommand(brandUpdateCmd)
	brandCmd.AddCommand(brandDeleteCmd)
}
