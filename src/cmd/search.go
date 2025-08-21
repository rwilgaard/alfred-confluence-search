package cmd

import (
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/rwilgaard/alfred-confluence-search/src/internal/util"
	cf "github.com/rwilgaard/confluence-go-api"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var (
	highlightReplacer = strings.NewReplacer("@@@hl@@@", "", "@@@endhl@@@", "")
	searchCmd         = &cobra.Command{
		Use:     "search",
		GroupID: groupBase,
		Run: func(cmd *cobra.Command, args []string) {
			if ok := alfredutils.HandleAuthentication(wf, keychainAccount); !ok {
				return
			}
			query := args[0]
			parsedQuery := util.ParseQuery(query)

			if a := util.Autocomplete(query); a != "" {
				if err := wf.Alfred.RunTrigger(a, query); err != nil {
					wf.FatalError(err)
				}
				return
			}

			api, err := setupAPIClient()
			if err != nil {
				wf.FatalError(err)
			}

			cql := util.BuildJQL(parsedQuery)
			params := cf.SearchQuery{
				CQL:    cql,
				Limit:  cfg.ResultsLimit,
				Expand: []string{"content.space", "content.history"},
			}

			log.Println(cql)

			pages, err := api.Search(params)
			if err != nil {
				wf.FatalError(err)
			}

			addPageItems(parsedQuery, pages)
			alfredutils.HandleFeedback(wf)
		},
	}
)

func addPageItems(query *util.ParsedQuery, pages *cf.Search) {
	if len(query.Spaces) == 1 {
		spaceId := strings.ToUpper(query.Spaces[0])
		wf.NewItem(fmt.Sprintf("Open %s Space Home", spaceId)).
			Icon(homeIcon).
			Arg("space").
			Var("item_url", fmt.Sprintf("%s/display/%s", cfg.BaseURL, spaceId)).
			Valid(true)
	}

	for _, page := range pages.Results {
		title := highlightReplacer.Replace(page.Title)
		title = html.UnescapeString(title)

		modTime := page.LastModified.Format("02-01-2006 15:04")
		space := page.Content.Space.Name
		subtitle := fmt.Sprintf("%s  •  Updated: %s", space, modTime)

		wf.NewItem(title).Subtitle(subtitle).
			UID(page.Content.ID).
			Icon(pageIcon).
			Var("item_url", cfg.BaseURL+page.URL).
			Arg("page").
			Valid(true)
	}
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
