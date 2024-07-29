package spa

import (
	"embed"
	"github.com/andresbott/go-carbon/libs/http/handlers"

	"net/http"
)

//go:embed files/ui/*
var UiFiles embed.FS

func NewCarbonSpa(path string) (http.Handler, error) {
	return handlers.NewSpaHAndler(
		UiFiles,
		"files/ui",
		path,
	)
}
