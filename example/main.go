package main

import (
	"fmt"
	"log/slog"

	"github.com/kr/pretty"
	cs2gsi "github.com/nescabir/go-cs2-gsi"
	"github.com/nescabir/go-cs2-gsi/models"
)

func main() {
	gsi := cs2gsi.New(cs2gsi.Config{
		ServerAddr: ":3000",
		LogLevel:   slog.LevelWarn,
	})

	cs2gsi.Subscribe(cs2gsi.Data, func(event cs2gsi.Event[*models.State]) {
		var activeWeapon *models.Weapon
		for _, weapon := range event.Data.Player.Weapons {
			if weapon.State == "active" {
				activeWeapon = weapon
				break
			}
		}
		pretty.Printf("Data: %+v\n", activeWeapon)
	})

	cs2gsi.Subscribe(cs2gsi.Mvp, func(event cs2gsi.Event[*models.Player]) {
		fmt.Printf("MVP: %s with %d kills (%d HS)\n",
			event.Data.Name, event.Data.State.Round_kills, event.Data.State.Round_killhs)
	})

	gsi.Listen()
}
