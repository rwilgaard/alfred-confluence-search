package cmd

import (
	"fmt"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/rwilgaard/alfred-confluence-search/src/internal/models"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var (
	getSpacesCmd = &cobra.Command{
		Use:     "spaces",
		Short:   "Get list of spaces",
		GroupID: groupSpaces,
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if ok := alfredutils.HandleAuthentication(wf, keychainAccount); !ok {
				return
			}

			var query string
			if len(args) > 0 {
				query = args[0]
			}

			if err := addSpaceItems(); err != nil {
				wf.FatalError(err)
			}

			if len(query) > 0 {
				wf.Filter(query)
			}
			alfredutils.HandleFeedback(wf)
		},
	}
)

func addSpaceItems() error {
	var spaceCache []models.Space

	maxCacheAge := time.Duration(cfg.CacheAge * int(time.Minute))
	if err := alfredutils.RefreshCache(wf, spaceCacheName, maxCacheAge, []string{"cache"}); err != nil {
		wf.FatalError(err)
	}

	err := alfredutils.LoadCache(wf, spaceCacheName, &spaceCache)
	if err != nil {
		return err
	}

	for _, space := range spaceCache {
		i := wf.NewItem(space.Key).
			UID(space.Key).
			Match(fmt.Sprintf("%s %s", space.Key, space.Name)).
			Icon(aw.IconWorkflow).
			Subtitle(space.Name).
			Arg(space.Key).
			Valid(true)

        i.NewModifier(aw.ModCmd).
            Subtitle("Cancel").
            Arg("cancel")
	}

	return nil
}

func init() {
	getCmd.AddCommand(getSpacesCmd)
}
