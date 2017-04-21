package main

import (
	"./app"
)

func main() {
	app := app.New()
	app.LoadConfig("config.yaml")
	app.Run()
}
