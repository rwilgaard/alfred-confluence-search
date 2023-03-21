package main

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
            Icon(getSpaceIcon(s.Key)).
            Subtitle(s.Name).
            Arg(prevQuery + s.Key).
            Valid(true)
    }
}
