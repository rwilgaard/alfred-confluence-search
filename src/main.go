package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	cf "github.com/rwilgaard/confluence-go-api"
)

type workflowConfig struct {
    URL      string `env:"confluence_url"`
    Username string `env:"username"`
    APIToken string `env:"apitoken"`
}

const (
    repo          = "rwilgaard/alfred-confluence-search"
    updateJobName = "checkForUpdates"
)

var (
    wf          *aw.Workflow
    cacheName   = "spaces.json"
    maxCacheAge = 24 * time.Hour
    spaceCache  []Space
)

func init() {
    wf = aw.New(
        update.GitHub(repo),
    )
}

func run() {
    if err := cli.Parse(wf.Args()); err != nil {
        wf.FatalError(err)
    }
    opts.Query = cli.Arg(0)

    if opts.Update {
        wf.Configure(aw.TextErrors(true))
        log.Println("Checking for updates...")
        if err := wf.CheckForUpdate(); err != nil {
            wf.FatalError(err)
        }
        return
    }

    if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
        log.Println("Running update check in background...")
        cmd := exec.Command(os.Args[0], "-update")
        if err := wf.RunInBackground(updateJobName, cmd); err != nil {
            log.Printf("Error starting update check: %s", err)
        }
    }

    if wf.UpdateAvailable() {
        wf.Configure(aw.SuppressUIDs(true))
        wf.NewItem("Update Available!").
            Subtitle("Press ⏎ to install").
            Autocomplete("workflow:update").
            Valid(false).
            Icon(aw.IconInfo)
    }

    cfg := &workflowConfig{}
    if err := wf.Config.To(cfg); err != nil {
        panic(err)
    }

    api, err := cf.NewAPI(cfg.URL, cfg.Username, cfg.APIToken)
    if err != nil {
        panic(err)
    }

    if opts.Cache {
        wf.Configure(aw.TextErrors(true))
        log.Println("[main] fetching spaces...")
        spaces := getSpaces(*api)
        if err := wf.Cache.StoreJSON(cacheName, spaces); err != nil {
            wf.FatalError(err)
        }
        log.Println("[main] cached spaces")
        return
    }

    if wf.Cache.Expired(cacheName, maxCacheAge) {
        wf.Rerun(0.3)
        if !wf.IsRunning("cache") {
            log.Println("[main] fetching spaces...")
            cmd := exec.Command(os.Args[0], "-cache")
            if err := wf.RunInBackground("cache", cmd); err != nil {
                wf.FatalError(err)
            } else {
                log.Printf("cache job already running.")
            }
        }
    }

    if opts.Spaces {
        runSpaces()
        if len(opts.Query) > 0 {
            wf.Filter(opts.Query)
        }
        wf.SendFeedback()
        return
    }

    if autocompleteSpaces(opts.Query) {
        if err := wf.Cache.StoreJSON("prev_query", opts.Query); err != nil {
            wf.FatalError(err)
        }
        if err := wf.Alfred.RunTrigger("spaces", ""); err != nil {
            wf.FatalError(err)
        }
        return
    }

    cql, spaceList := parseQuery(opts.Query)
    pages := getPages(*api, cql)

    if len(spaceList) == 1 {
        homeIcon := aw.Icon{Value: "icons/home.png"}
        spaceId := strings.ToUpper(spaceList[0])
        wf.NewItem(fmt.Sprintf("Open %s Space Home", spaceId)).
            Icon(&homeIcon).
            Arg("space").
            Var("item_url", spaceId).
            Valid(true)
    }

    for _, page := range pages.Results {
        title := strings.ReplaceAll(page.Title, "@@@hl@@@", "")
        title = strings.ReplaceAll(title, "@@@endhl@@@", "")
        modTime := page.LastModified.Time.Format("02-01-2006 15:04")
        sub := fmt.Sprintf("%s  |  Updated: %s", page.Content.Space.Name, modTime)
        wf.NewItem(html.UnescapeString(title)).Subtitle(sub).
            Var("item_url", page.URL).
            Arg("page").
            Icon(getSpaceIcon(page.Content.Space.Key)).
            Valid(true)
    }

    getSpaces(*api)
    if wf.IsEmpty() {
        wf.NewItem("No results found...").
            Subtitle("Try a different query?").
            Icon(aw.IconInfo)
    }
    wf.SendFeedback()
}

func main() {
    wf.Run(run)
}
