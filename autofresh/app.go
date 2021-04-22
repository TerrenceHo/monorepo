package autofresh

import "fmt"

const logo = `
    ___         __        ______               __  
   /   | __  __/ /_____  / ____/_______  _____/ /_ 
  / /| |/ / / / __/ __ \/ /_  / ___/ _ \/ ___/ __ \
 / ___ / /_/ / /_/ /_/ / __/ / /  /  __(__  ) / / /
/_/  |_\__,_/\__/\____/_/   /_/   \___/____/_/ /_/ 
`

func Start(c Config) {
	if !c.HideBanner {
		fmt.Println(logo)
	}
	fmt.Println(c.Cmd)
	fmt.Println(c.Watch)
	fmt.Println(c.Extensions)
}
