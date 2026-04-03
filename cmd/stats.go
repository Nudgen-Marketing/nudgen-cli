package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nudgen/nudgen-cli/internal/api"
	"github.com/nudgen/nudgen-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/charmbracelet/lipgloss"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View Nudgen metrics and stats",
	Long: `Display a high-level analytics overview for your active team. 
Includes counts for Sent, Delivered, Opens, Clicks, and Engagement Rate (CTR).`,
	Example: `  nudgen stats
  nudgen stats --json`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		stats, err := client.GetStats()
		if err != nil {
			log.Fatalf("Error fetching stats: %v", err)
		}

		if UseJSON {
			json.NewEncoder(os.Stdout).Encode(stats)
			return
		}

		renderDashboard(stats)
	},
}

func renderDashboard(stats *api.Stats) {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D1FF")).
		MarginBottom(1)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#333333")).
		Padding(1, 2).
		MarginRight(2).
		Width(24)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		MarginTop(1)

	createCard := func(label string, value string) string {
		return cardStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				labelStyle.Render(label),
				valueStyle.Render(value),
			),
		)
	}

	fmt.Println(titleStyle.Render("NUDGEN ANALYTICS OVERVIEW"))

	row1 := lipgloss.JoinHorizontal(lipgloss.Top,
		createCard("TOTAL SENT", fmt.Sprintf("%d", stats.TotalSent)),
		createCard("DELIVERED", fmt.Sprintf("%d", stats.TotalDelivered)),
		createCard("FAILED", fmt.Sprintf("%d", stats.TotalFailed)),
	)

	row2 := lipgloss.JoinHorizontal(lipgloss.Top,
		createCard("OPENS", fmt.Sprintf("%d", stats.TotalOpens)),
		createCard("CLICKS", fmt.Sprintf("%d", stats.TotalClicks)),
		createCard("CTR", fmt.Sprintf("%.2f%%", stats.CTR)),
	)

	fmt.Println(row1)
	fmt.Println()
	fmt.Println(row2)
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
