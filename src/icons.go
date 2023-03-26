package main

import (
	"fmt"
	"os"

	aw "github.com/deanishe/awgo"
)

var (
    pageIcon = &aw.Icon{Value: "icons/page.png"}
    homeIcon = &aw.Icon{Value: "icons/home.png"}
)

func getSpaceIcon(spaceKey string) *aw.Icon {
    iconPath := fmt.Sprintf("icons/%s.png", spaceKey)
    icon := aw.IconWorkflow

    if _, err := os.Stat(iconPath); err == nil {
        icon = &aw.Icon{Value: iconPath}
    }

    return icon
}
