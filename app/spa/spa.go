package spa

import (
	"embed"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"net/http"
)

//go:embed files/ui/*
var UiFiles embed.FS

func NewCarbonSpa() (http.Handler, error) {
	return handlers.NewSpaHAndler(
		UiFiles,
		"files/ui",
		"/",
	)
}
