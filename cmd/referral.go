package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/nudgen/nudgen-cli/internal/api"
	"github.com/nudgen/nudgen-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/charmbracelet/lipgloss"
)

var referralCmd = &cobra.Command{
	Use:   "referral",
	Short: "Manage your Nudgen referral account",
	Long: `Track your referrals, commissions, and get your unique referral link 
to share with your network and earn commissions.`,
	Example: `  nudgen referral link
  nudgen referral stats
  nudgen referral activity
  nudgen referral payouts`,
}

var referralLinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Get your unique referral link",
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		info, err := client.GetReferralInfo()
		if err != nil {
			log.Fatalf("Error fetching referral info: %v", err)
		}

		if info.Account == nil {
			fmt.Println("You haven't signed up for the referral program yet.")
			fmt.Println("Visit https://nudgen.cc/affiliates to join!")
			return
		}

		link := "https://nudgen.cc/r/" + info.Account.ReferralCode
		fmt.Printf("\nYour Referral Link: %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#00D1FF")).Bold(true).Render(link))
		fmt.Println("Share this link with your audience to earn 30% recurring commission!")
	},
}

var referralStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View your referral performance stats",
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		info, err := client.GetReferralInfo()
		if err != nil {
			log.Fatalf("Error fetching referral info: %v", err)
		}

		if info.Account == nil {
			fmt.Println("You haven't signed up for the referral program yet.")
			return
		}

		if UseJSON {
			json.NewEncoder(os.Stdout).Encode(info)
			return
		}

		renderReferralDashboard(info)
	},
}

func renderReferralDashboard(info *api.ReferralInfo) {
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

	fmt.Println(titleStyle.Render("NUDGEN REFERRAL OVERVIEW"))

	row1 := lipgloss.JoinHorizontal(lipgloss.Top,
		createCard("CLICKS", fmt.Sprintf("%d", info.Account.Count.Clicks)),
		createCard("REFERRALS", fmt.Sprintf("%d", info.Account.Count.Referrals)),
		createCard("COMMISSIONS", fmt.Sprintf("%d", info.Account.Count.Commissions)),
	)

	totalEarnings := (info.Balances.Approved + info.Balances.Paid) / 100
	row2 := lipgloss.JoinHorizontal(lipgloss.Top,
		createCard("PENDING", fmt.Sprintf("$%.2f", float64(info.Balances.Pending)/100)),
		createCard("APPROVED", fmt.Sprintf("$%.2f", float64(info.Balances.Approved)/100)),
		createCard("TOTAL EARNED", fmt.Sprintf("$%.d", totalEarnings)),
	)

	fmt.Println(row1)
	fmt.Println()
	fmt.Println(row2)
}

var referralActivityCmd = &cobra.Command{
	Use:   "activity",
	Short: "View recent referral activity and click history",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		analytics, err := client.GetAffiliateAnalytics()
		if err != nil {
			log.Fatalf("Error fetching activity: %v", err)
		}

		if UseJSON {
			json.NewEncoder(os.Stdout).Encode(analytics)
			return
		}

		headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00D1FF")).Underline(true)
		fmt.Println(headerStyle.Render("\nRECENT CLICK HISTORY"))
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "DATE\tIP ADDRESS\tUA\tOUTCOME")
		for _, click := range analytics.RecentClicks {
			outcome := "Visit"
			if click.Converted {
				outcome = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("Signed Up")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", 
				click.CreatedAt.Format("2006-01-02 15:04"), 
				click.IP, 
				click.UA,
				outcome,
			)
		}
		w.Flush()
	},
}

var referralListCmd = &cobra.Command{
	Use:   "referrals",
	Short: "List all teams referred by you",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		analytics, err := client.GetAffiliateAnalytics()
		if err != nil {
			log.Fatalf("Error fetching referrals: %v", err)
		}

		if UseJSON {
			json.NewEncoder(os.Stdout).Encode(analytics.Referrals)
			return
		}

		headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00D1FF")).Underline(true)
		fmt.Println(headerStyle.Render("\nYOUR REFERRALS"))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "TEAM NAME\tCREATED AT\tSTATUS")
		for _, ref := range analytics.Referrals {
			status := ref.Status
			if status == "active" {
				status = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("Paying Plan")
			} else {
				status = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Free Trial")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", 
				ref.TeamName, 
				ref.CreatedAt.Format("2006-01-02"),
				status,
			)
		}
		w.Flush()
	},
}

var referralPayoutsCmd = &cobra.Command{
	Use:   "payouts",
	Short: "View your payout request history",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := config.GetToken()
		client := api.NewClientWithToken(token)

		payouts, err := client.GetAffiliatePayouts()
		if err != nil {
			log.Fatalf("Error fetching payouts: %v", err)
		}

		if UseJSON {
			json.NewEncoder(os.Stdout).Encode(payouts)
			return
		}

		headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00D1FF")).Underline(true)
		fmt.Println(headerStyle.Render("\nPAYOUT HISTORY"))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "DATE\tAMOUNT\tMETHOD\tSTATUS")
		for _, p := range payouts {
			statusStyle := lipgloss.NewStyle()
			switch p.Status {
			case "paid":
				statusStyle = statusStyle.Foreground(lipgloss.Color("#00FF00"))
			case "pending":
				statusStyle = statusStyle.Foreground(lipgloss.Color("#FFA500"))
			case "failed":
				statusStyle = statusStyle.Foreground(lipgloss.Color("#FF0000"))
			}

			fmt.Fprintf(w, "%s\t$%.2f\t%s\t%s\n", 
				p.CreatedAt.Format("2006-01-02"), 
				float64(p.Amount)/100,
				p.Method,
				statusStyle.Render(p.Status),
			)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(referralCmd)
	referralCmd.AddCommand(referralLinkCmd)
	referralCmd.AddCommand(referralStatsCmd)
	referralCmd.AddCommand(referralActivityCmd)
	referralCmd.AddCommand(referralListCmd)
	referralCmd.AddCommand(referralPayoutsCmd)
}
