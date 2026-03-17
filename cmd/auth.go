package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/misonikomipan/homebox-cli/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"os"
)

// newLoginCmd is also exposed as top-level `hb login`.
func newLoginCmd() *cobra.Command {
	var email, password, endpoint string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login and store authentication token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if endpoint != "" {
				if err := config.SetEndpoint(endpoint); err != nil {
					return err
				}
			}
			if email == "" {
				fmt.Print("Email: ")
				fmt.Scanln(&email)
			}
			if password == "" {
				fmt.Print("Password: ")
				b, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()
				if err != nil {
					return err
				}
				password = string(b)
			}
			c, err := client.New(false)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/users/login", map[string]string{
				"username": email,
				"password": password,
			})
			if err != nil {
				return err
			}
			var resp map[string]any
			if err := json.Unmarshal(data, &resp); err != nil {
				return err
			}
			token, ok := resp["token"].(string)
			if !ok || token == "" {
				return fmt.Errorf("no token in response")
			}
			// Remove "Bearer " prefix if it exists in the token string from API
			token = strings.TrimPrefix(token, "Bearer ")
			if err := config.SetToken(token); err != nil {
				return err
			}
			out, _ := json.MarshalIndent(map[string]string{
				"message":  "Login successful",
				"endpoint": config.GetEndpoint(),
			}, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
	cmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	cmd.Flags().StringVar(&endpoint, "endpoint", "", "Override API endpoint URL")
	return cmd
}

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout and clear stored token",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err == nil {
				_, _ = c.Post("/v1/users/logout", nil)
			}
			if err := config.ClearToken(); err != nil {
				return err
			}
			fmt.Println(`{"message": "Logged out"}`)
			return nil
		},
	}
}

func newAuthCmd() *cobra.Command {
	auth := &cobra.Command{
		Use:   "auth",
		Short: "Authentication and user account",
	}

	auth.AddCommand(newLoginCmd())
	auth.AddCommand(newLogoutCmd())

	auth.AddCommand(&cobra.Command{
		Use:   "refresh",
		Short: "Refresh authentication token",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/users/refresh", nil)
			if err != nil {
				return err
			}
			var resp map[string]any
			if err := json.Unmarshal(data, &resp); err == nil {
				if token, ok := resp["token"].(string); ok && token != "" {
					_ = config.SetToken(token)
				}
			}
			client.PrintJSON(data)
			return nil
		},
	})

	auth.AddCommand(&cobra.Command{
		Use:   "me",
		Short: "Get current user information",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/users/self", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	var updateName, updateEmail string
	updateMe := &cobra.Command{
		Use:   "update-me",
		Short: "Update current user profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			payload := map[string]string{}
			if updateName != "" {
				payload["name"] = updateName
			}
			if updateEmail != "" {
				payload["email"] = updateEmail
			}
			data, err := c.Put("/v1/users/self", payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	updateMe.Flags().StringVar(&updateName, "name", "", "New name")
	updateMe.Flags().StringVar(&updateEmail, "email", "", "New email")
	auth.AddCommand(updateMe)

	auth.AddCommand(&cobra.Command{
		Use:   "change-password",
		Short: "Change user password",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("Current password: ")
			cur, _ := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			fmt.Print("New password: ")
			nw, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				return err
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Put("/v1/users/change-password", map[string]string{
				"current": string(cur),
				"new":     string(nw),
			})
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	return auth
}
