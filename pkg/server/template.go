package server

import (
	"embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"
)

const (
	BasePath   = "templates"
	LayoutFile = "layout.html"
)

//go:embed templates/*
var templateFS embed.FS

func RenderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
	var layoutPath = filepath.Join(BasePath, LayoutFile)
	var templatePath = filepath.Join(BasePath, tmpl+".html")

	t, err := template.New("layout").ParseFS(templateFS, layoutPath, templatePath)
	if err != nil {
		slog.Error("Failed to parse template:" + err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if data == nil {
		data = make(map[string]interface{})
	}

	data["Year"] = time.Now().Year()

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
