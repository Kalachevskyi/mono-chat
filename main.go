package main

import (
	"log"
	"os"

	"github.com/Kalachevskyi/mono-chat/cmd"
)

//main - start mono-chat app...
func main() {
	app := cmd.RootCMD{}
	if err := app.Init().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
