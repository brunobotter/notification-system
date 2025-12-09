package main

import (
	"github.com/brunobotter/notification-system/main/app"
	"github.com/brunobotter/notification-system/main/providers"
)

func main() {
	app.NewApplication(providers.List()).Bootstrap()
}
