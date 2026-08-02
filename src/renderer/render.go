package renderer

import (
  "embed"
  "html/template"
  "net/http"
)

var GlobalTemplateFS embed.FS

type PageData struct {
  Data any
}

func InitTemplates(templateFS embed.FS) {
  GlobalTemplateFS = templateFS
}

func Render(res http.ResponseWriter, layoutName string, templatePaths []string, data PageData) {
  res.Header().Set("Content-Type", "text/html; charset=utf-8")
  //Compile the precise array of files requested by the handler.
  tmpl, err := template.ParseFS(GlobalTemplateFS, templatePaths...)
  if err != nil {
    http.Error(res, "Template parsing failure: " + err.Error(), http.StatusInternalServerError)
  } else {
    err = tmpl.ExecuteTemplate(res, layoutName, data)
    if err != nil {
      http.Error(res, "Template rendering error: " + err.Error(), http.StatusInternalServerError)
    }
  }
}
