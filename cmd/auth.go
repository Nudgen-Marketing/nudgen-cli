package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nudgen/nudgen-cli/internal/api"
	"github.com/nudgen/nudgen-cli/internal/config"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with a existing Nudgen Personal Access Token",
	Long: `Authorization is required to access your Nudgen data. 
This command opens a local web server to receive a Personal Access Token (PAT) 
from the Nudgen web dashboard.`,
	Example: `  nudgen login`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.IsLoggedIn() {
			fmt.Println("You are already logged in.")
			return
		}

		port := "3456"
		callbackURL := fmt.Sprintf("http://localhost:%s/callback", port)
		loginURL := fmt.Sprintf("https://app.nudgen.net/settings/cli/login?callback=%s", callbackURL)

		fmt.Println("Opening your browser to authenticate with Nudgen.")
		fmt.Printf("If the browser does not open automatically, visit:\n%s\n\n", loginURL)
		fmt.Println("Waiting for authentication callback...")

		// Channel to receive token from the local web server
		tokenChan := make(chan string)
		errChan := make(chan error)

		mux := http.NewServeMux()
		srv := &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		}

		mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")
			if token == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error: No token provided in the callback."))
				errChan <- fmt.Errorf("no token provided in callback")
				return
			}
			
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<html><body><h2>Authentication successful!</h2><p>You can safely close this browser window and return to your terminal.</p></body></html>"))
			
			tokenChan <- token
		})

		// Start server in background
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start local callback server: %v", err)
			}
		}()

		// Open browser
		if err := browser.OpenURL(loginURL); err != nil {
			fmt.Printf("Failed to open browser automatically. Please manually navigate to the URL above.\n")
		}

		// Wait for token
		var pat string
		select {
		case pat = <-tokenChan:
			// Ensure server is shut down properly
			go srv.Shutdown(context.Background())
		case err := <-errChan:
			go srv.Shutdown(context.Background())
			log.Fatalf("Authentication failed: %v", err)
		}

		if pat == "" {
			fmt.Println("Login cancelled or token missing.")
			os.Exit(0)
		}

		fmt.Println("\nToken received! Verifying...")

		// Verify PAT against API
		client := api.NewClientWithToken(pat)
		user, err := client.GetMe()
		if err != nil {
			log.Fatalf("Authentication failed. Is the PAT correct? Error: %v", err)
		}

		// Success
		if err := config.SaveToken(pat); err != nil {
			log.Fatalf("Failed to save token to keyring: %v", err)
		}

		fmt.Printf("Successfully logged in as %s (%s)!\n", user.Name, user.Email)
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of the CLI by removing the PAT",
	Long:  `Permanently removes the stored Personal Access Token from your local keyring.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.ClearToken(); err != nil {
			fmt.Println("You were not logged in, or there was an error logging you out.")
			return
		}
		fmt.Println("Logged out successfully.")
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "See who is currently logged in",
	Long:  `Display your Nudgen account identity and the current session's user email.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.IsLoggedIn() {
			fmt.Println("Not logged in. Run 'nudgen login'.")
			return
		}

		token, err := config.GetToken()
		if err != nil {
			log.Fatalf("Error reading token: %v", err)
		}

		client := api.NewClientWithToken(token)
		user, err := client.GetMe()
		if err != nil {
			log.Fatalf("Error connecting to API: %v", err)
		}

		fmt.Printf("Logged in as: %s (%s)\n", user.Name, user.Email)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
}
