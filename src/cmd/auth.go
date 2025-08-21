package cmd

import (
	"fmt"
	"log"

	"github.com/ncruces/zenity"
	cf "github.com/rwilgaard/confluence-go-api"
	"github.com/spf13/cobra"
)

var (
	authCmd = &cobra.Command{
		Use:     "auth",
		Short:   "Authenticate",
		Args:    cobra.NoArgs,
		GroupID: groupBase,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, token, err := zenity.Password(zenity.Title("Enter API Token"))
			if err != nil {
				return err
			}

			if err := validateAPIToken(token); err != nil {
				zerr := zenity.Error(
					fmt.Sprintf("Error authenticating: %s", err),
					zenity.ErrorIcon,
				)
				if zerr != nil {
					log.Printf("Zenity error dialog failed: %v", zerr)
				}
				return err
			}

			if err := wf.Keychain.Set(keychainAccount, token); err != nil {
				zerr := zenity.Error(
					fmt.Sprintf("Failed to store token in Keychain: %s", err),
					zenity.ErrorIcon,
				)
				if zerr != nil {
					log.Printf("Zenity error dialog failed: %v", zerr)
				}
				return err
			}

			_ = zenity.Notify("Authentication successful!", zenity.Title("Alfred Confluence Search"))
			return nil
		},
	}
)

func validateAPIToken(token string) error {
	api, err := cf.NewAPI(cfg.BaseURL, cfg.Username, token)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	if _, err = api.CurrentUser(); err != nil {
		return fmt.Errorf("invalid token or API error: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(authCmd)
}
