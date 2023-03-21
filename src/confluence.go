package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	cf "github.com/rwilgaard/confluence-go-api"
)

type Space struct {
    Key  string
    Name string
}

func getSpaces(api cf.API) []Space {
    var spaces []Space
    params := cf.AllSpacesQuery{
        Limit: 9999,
        Type:  "global",
    }

    result, err := api.GetAllSpaces(params)
    if err != nil {
        panic(err)
    }

    log.Println(result.Size)

    for _, s := range result.Results {
        space := Space{Key: s.Key, Name: s.Name}
        spaces = append(spaces, space)
    }

    return spaces
}

func parseQuery(query string) (string, []string) {
    if wf.Cache.Exists(cacheName) {
        if err := wf.Cache.LoadJSON(cacheName, &spaceCache); err != nil {
            wf.FatalError(err)
        }
    }

    cql := "siteSearch ~ '%s' AND type = page"
    matchParam := regexp.MustCompile(`^@\w+`)
    var text string
    var spaceList []string
    for _, w := range strings.Split(query, " ") {
        switch {
        case matchParam.MatchString(w):
            spaceKey := w[1:]
            if spaceExists(spaceKey, spaceCache) {
                spaceList = append(spaceList, spaceKey)
            } else {
                title := fmt.Sprintf("%s space not found...", strings.ToUpper(spaceKey))
                s := fuzzy.Find(spaceKey, spaceCacheToSlice(spaceCache))
                sub := fmt.Sprintf("Did you mean %s?", strings.Join(s, ", "))
                wf.NewItem(title).Subtitle(sub).Icon(aw.IconInfo)
            }
        default:
            text = text + fmt.Sprintf("%s ", w)
        }
    }

    cql = fmt.Sprintf(cql, strings.TrimSpace(text))
    if len(spaceList) > 0 {
        cql = cql + fmt.Sprintf(" AND space.key in (%s)", strings.Join(spaceList, ","))
    }

    return cql, spaceList
}

func spaceExists(key string, spaces []Space) bool {
    for _, s := range spaces{
        if strings.EqualFold(s.Key, key) {
            return true
        }
    }
    return false
}

func spaceCacheToSlice(spaces []Space) []string {
    var spaceList []string
    for _, s := range spaces {
        spaceList = append(spaceList, s.Key)
    }
    return spaceList
}

func autocompleteSpaces(query string) bool {
    for _, w := range strings.Split(query, " ") {
        if w == "@" {
            return true
        }
    }
    return false
}

func getPages(api cf.API, cql string) *cf.Search {
    params := cf.SearchQuery{
        CQL:    cql,
        Limit:  15,
        Expand: []string{"content.space", "content.history"},
    }

    result, err := api.Search(params)
    if err != nil {
        panic(err)
    }

    return result
}
