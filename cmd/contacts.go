package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"strings"
	"github.com/charmbracelet/huh"
	"github.com/nudgen/nudgen-cli/internal/api"
	"github.com/nudgen/nudgen-cli/internal/config"
	"github.com/spf13/cobra"
)

var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Manage contacts in Nudgen",
	Long: `Lead management and segmentation for your team.
Use this command to list all contacts, delete outdated ones,
import new leads from CSV, or create them manually.`,
	Example: `  nudgen contacts list
  nudgen contacts create
  nudgen contacts delete c...id`,
}

var contactsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your contacts",
	Long:  `Retrieve and display managed leads from the active team's database.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		contacts, err := client.GetContacts()
		if err != nil {
			log.Fatalf("Error fetching contacts: %v", err)
		}

		if UseJSON {
			json.NewEncoder(os.Stdout).Encode(contacts)
			return
		}

		if len(contacts) == 0 {
			fmt.Println("No contacts found.")
			return
		}

		fmt.Println("--- Active Contacts ---")
		for _, c := range contacts {
			phone := "-"
			if c.Phone != "" {
				phone = c.Phone
			}
			tags := "-"
			if len(c.Tags) > 0 {
				tags = strings.Join(c.Tags, ", ")
			}
			desc := ""
			if c.Description != "" {
				desc = fmt.Sprintf(" | Desc: %s", c.Description)
			}
			fmt.Printf("[%s] %s <%s> | Status: %s | Created: %s | Phone: %s | Tags: [%s]%s\n",
				c.ID, c.Name, c.Email, c.Status, c.CreatedAt.Format("2006-01-02"), phone, tags, desc)
		}
	},
}

var contactsDeleteCmd = &cobra.Command{
	Use:   "delete [contact-id]",
	Short: "Delete a contact",
	Long:  `Permanently remove a lead from your team's contact list.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		id := args[0]
		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		err := client.DeleteContact(id)
		if err != nil {
			log.Fatalf("Error deleting contact: %v", err)
		}

		fmt.Printf("Successfully deleted contact: %s\n", id)
	},
}

var contactsImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Interactively import a CSV mapping of contacts",
	Long:  `Launch an interactive importer to map columns and import leads from a local CSV file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[Not Implemented] Interactive file picker and CSV mapping UI using Bubbletea.")
	},
}

var contactsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new contact manually",
	Long:  `Adds a lone contact using an interactive wizard for quick entry.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		var (
			email       string
			name        string
			phone       string
			description string
			tagsRaw     string
		)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Email Address").
					Value(&email).
					Validate(func(s string) error {
						if len(s) < 3 {
							return fmt.Errorf("Invalid email")
						}
						return nil
					}),
				huh.NewInput().
					Title("Full Name").
					Value(&name),
				huh.NewInput().
					Title("Phone Number (Optional)").
					Value(&phone),
				huh.NewInput().
					Title("Description (Optional)").
					Value(&description),
				huh.NewInput().
					Title("Tags (comma separated)").
					Value(&tagsRaw),
			),
		)

		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		var tags []string
		if tagsRaw != "" {
			tags = strings.Split(tagsRaw, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		body := map[string]interface{}{
			"email":       email,
			"name":        name,
			"phone":       phone,
			"description": description,
			"tags":        tags,
		}

		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		contactID, err := client.CreateContact(body)
		if err != nil {
			log.Fatalf("Error creating contact: %v", err)
		}

		fmt.Printf("Successfully created contact: %s <%s> (ID: %s)\n", name, email, contactID)
	},
}

var contactsUpdateCmd = &cobra.Command{
	Use:   "update [contact-id]",
	Short: "Update a contact settings",
	Long:  `Launches a wizard to update details or tags for an existing contact.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		id := args[0]
		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		// Pre-fetch data or just use empty if not found
		// For now we'll just show the update wizard.
		var (
			name        string
			phone       string
			description string
			tagsRaw     string
			status      string
		)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Full Name").
					Value(&name),
				huh.NewInput().
					Title("Phone Number").
					Value(&phone),
				huh.NewInput().
					Title("Description").
					Value(&description),
				huh.NewInput().
					Title("Tags (comma separated)").
					Value(&tagsRaw),
				huh.NewSelect[string]().
					Title("Status").
					Options(
						huh.NewOption("Active", "active"),
						huh.NewOption("Unsubscribed", "unsubscribed"),
						huh.NewOption("Bounced", "bounced"),
					).
					Value(&status),
			).Title("Update Contact: " + id),
		)

		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		body := map[string]interface{}{}
		if name != "" {
			body["name"] = name
		}
		if phone != "" {
			body["phone"] = phone
		}
		if description != "" {
			body["description"] = description
		}
		if status != "" {
			body["status"] = status
		}
		if tagsRaw != "" {
			parts := strings.Split(tagsRaw, ",")
			tags := []string{}
			for _, p := range parts {
				tags = append(tags, strings.TrimSpace(p))
			}
			body["tags"] = tags
		}

		contactID, err := client.UpdateContact(id, body)
		if err != nil {
			log.Fatalf("Error updating contact: %v", err)
		}

		fmt.Printf("Successfully updated contact: %s\n", contactID)
	},
}

func init() {
	rootCmd.AddCommand(contactsCmd)
	contactsCmd.AddCommand(contactsListCmd)
	contactsCmd.AddCommand(contactsImportCmd)
	contactsCmd.AddCommand(contactsCreateCmd)
	contactsCmd.AddCommand(contactsUpdateCmd)
	contactsCmd.AddCommand(contactsDeleteCmd)
}
