package main

import (
	"fmt"
	"html"
	"strings"

	cf "github.com/rwilgaard/confluence-go-api"
)

func runSpaces() {
    if wf.Cache.Exists(cacheName) {
        if err := wf.Cache.LoadJSON(cacheName, &spaceCache); err != nil {
            wf.FatalError(err)
        }
    }

    var prevQuery string
    if err := wf.Cache.LoadJSON("prev_query", &prevQuery); err != nil {
        wf.FatalError(err)
    }

    for _, s := range spaceCache {
        wf.NewItem(s.Key).
            UID(s.Key).
            Match(fmt.Sprintf("%s %s", s.Key, s.Name)).
            Icon(getSpaceIcon(s.Key)).
            Subtitle(s.Name).
            Arg(prevQuery + s.Key).
            Valid(true)
    }
}

func runSearch(api *cf.API) {
    cql, spaceList := parseQuery(opts.Query)
    pages := getPages(*api, cql)

    if len(spaceList) == 1 {
        spaceId := strings.ToUpper(spaceList[0])
        wf.NewItem(fmt.Sprintf("Open %s Space Home", spaceId)).
            Icon(homeIcon).
            Arg("space").
            Var("item_url", fmt.Sprintf("%s/display/%s", cfg.URL, spaceId)).
            Valid(true)
    }

    for _, page := range pages.Results {
        title := strings.ReplaceAll(page.Title, "@@@hl@@@", "")
        title = strings.ReplaceAll(title, "@@@endhl@@@", "")
        modTime := page.LastModified.Time.Format("02-01-2006 15:04")
        sub := fmt.Sprintf("%s  •  Updated: %s", page.Content.Space.Name, modTime)
        wf.NewItem(html.UnescapeString(title)).Subtitle(sub).
            Icon(pageIcon).
            Var("item_url", cfg.URL+page.URL).
            Arg("page").
            Valid(true)
    }
}
