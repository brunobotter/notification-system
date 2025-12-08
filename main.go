package main

import (
	"github.com/brunobotter/notification-system/app/main/app"
	"github.com/brunobotter/notification-system/app/main/providers"
)

func main() {
	app.NewApplication(providers.List()).Bootstrap()
}
