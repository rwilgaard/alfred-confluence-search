package cmd

import (
	"fmt"
	"log"

	"github.com/rwilgaard/alfred-confluence-search/src/internal/models"
	cf "github.com/rwilgaard/confluence-go-api"
	"github.com/spf13/cobra"
)

var (
	cacheCmd = &cobra.Command{
		Use:   "cache",
		Short: "refresh cache of spaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("[cache] fetching spaces...")

			api, err := setupAPIClient()
			if err != nil {
				return err
			}

			if err := fetchAndCacheSpaces(api); err != nil {
				return err
			}

			log.Printf("[cache] spaces fetched")
			return nil
		},
	}
)

func fetchAndCacheSpaces(api *cf.API) error {
	params := cf.AllSpacesQuery{
		Limit: 9999,
		Type:  "global",
	}

	result, err := api.GetAllSpaces(params)
	if err != nil {
		return fmt.Errorf("failed to get all spaces from API: %w", err)
	}

	spaces := make([]models.Space, 0, len(result.Results))
	for _, space := range result.Results {
		spaces = append(spaces, models.Space{Key: space.Key, Name: space.Name})
	}

	if err := wf.Cache.StoreJSON(spaceCacheName, spaces); err != nil {
		return fmt.Errorf("failed to store spaces in cache: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}
