package cmd

import (
	"fmt"
	"log"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/ncruces/zenity"
	cf "github.com/rwilgaard/confluence-go-api"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
	"go.deanishe.net/fuzzy"
)

type workflowConfig struct {
	BaseURL          string `env:"confluence_url"`
	Username         string `env:"username"`
	APIToken         string
	ResultsLimit     int    `env:"results_limit"`
	CacheAge         int    `env:"cache_age"`
	RestrictToSpaces string `env:"restrict_to_spaces"`
}

const (
	repo            = "rwilgaard/alfred-confluence-search"
	keychainAccount = "alfred-confluence-search"
	groupBase       = "group-base"
	groupPages      = "group-pages"
	groupSpaces     = "group-spaces"
	spaceCacheName  = "spaces.json"
)

var (
	wf       *aw.Workflow
	cfg      = &workflowConfig{}
	pageIcon = &aw.Icon{Value: "icons/page.png"}
	homeIcon = &aw.Icon{Value: "icons/home.png"}
	rootCmd  = &cobra.Command{
		Use:           "confluence",
		Short:         "confluence is a CLI to be used by Alfred for ...",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)

func Execute() {
	wf.Run(run)
}

func run() {
	rootCmd.AddGroup(
		&cobra.Group{
			ID:    groupBase,
			Title: "Basic Commands: ",
		},
	)
	alfredutils.AddClearAuthMagic(wf, keychainAccount)

	if err := alfredutils.InitWorkflow(wf, cfg); err != nil {
		wf.FatalError(err)
	}

	if err := alfredutils.CheckForUpdates(wf); err != nil {
		wf.FatalError(err)
	}

	if err := rootCmd.Execute(); err != nil {
		wf.FatalError(err)
	}
}

func setupAPIClient() (*cf.API, error) {
	token, err := wf.Keychain.Get(keychainAccount)
	if err != nil {
		zerr := zenity.Error(
			fmt.Sprintf("Error retrieving credentials from Keychain: %s", err),
			zenity.ErrorIcon,
		)
		if zerr != nil {
			log.Printf("Zenity error dialog failed: %v", zerr)
		}
		return nil, err
	}

	api, err := cf.NewAPI(cfg.BaseURL, cfg.Username, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Confluence API client: %w", err)
	}

	return api, nil
}

func init() {
	sopts := []fuzzy.Option{
		fuzzy.AdjacencyBonus(10.0),
		fuzzy.LeadingLetterPenalty(-0.1),
		fuzzy.MaxLeadingLetterPenalty(-3.0),
		fuzzy.UnmatchedLetterPenalty(-0.5),
	}
	wf = aw.New(
		aw.SortOptions(sopts...),
		update.GitHub(repo),
	)
}
