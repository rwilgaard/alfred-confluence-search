package main

import "flag"

var (
    opts = &options{}
    cli  = flag.NewFlagSet("alfred-confluence-search", flag.ContinueOnError)
)

type options struct {
    // Arguments
    Query string

    // Commands
    GetIcons bool
    Update   bool
    Cache    bool
    Spaces   bool
}

func init() {
    cli.BoolVar(&opts.Update, "update", false, "check for updates")
    cli.BoolVar(&opts.Cache, "cache", false, "cache spaces")
    cli.BoolVar(&opts.GetIcons, "geticons", false, "get all space icons")
    cli.BoolVar(&opts.Spaces, "spaces", false, "list spaces")
}
