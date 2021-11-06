package main

import (
	_ "embed"

	"github.com/jairojj/docker-gui/internal"
	"github.com/wailsapp/wails"
)

//go:embed frontend/public/build/bundle.js
var js string

//go:embed frontend/public/build/bundle.css
var css string

func main() {
	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "docker-gui",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})

	api := &internal.Api{}

	app.Bind(api)
	app.Run()
}
