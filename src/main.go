package main

import (
    "log"
    "os"
    "os/exec"
    "time"

    aw "github.com/deanishe/awgo"
    "github.com/deanishe/awgo/update"
    cf "github.com/rwilgaard/confluence-go-api"
    "go.deanishe.net/fuzzy"
)

type workflowConfig struct {
    URL      string `env:"confluence_url"`
    Username string `env:"username"`
    APIToken string
}

const (
    repo            = "rwilgaard/alfred-confluence-search"
    updateJobName   = "checkForUpdates"
    keychainAccount = "alfred-confluence-search"
)

var (
    wf          *aw.Workflow
    cfg         *workflowConfig
    cacheName   = "spaces.json"
    maxCacheAge = 24 * time.Hour
    spaceCache  []Space
)

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
        aw.AddMagic(magicAuth{wf}),
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

    cfg = &workflowConfig{}
    if err := wf.Config.To(cfg); err != nil {
        panic(err)
    }

    if opts.Auth {
        runAuth()
    }

    token, err := wf.Keychain.Get(keychainAccount)
    if err != nil {
        wf.NewItem("You're not logged in.").
            Subtitle("Press ⏎ to authenticate").
            Icon(aw.IconInfo).
            Arg("auth").
            Valid(true)
        wf.SendFeedback()
        return
    }

    cfg.APIToken = token

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

    runSearch(api)

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
