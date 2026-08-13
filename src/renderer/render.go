package renderer

import (
  "bytes"
  "embed"
  /***
  Package template (html/template) implements data-driven templates for generating HTML output safe against code injection. It
  provides the same interface as text/template and should be used instead of text/template whenever the output is HTML.
  ***/
  "html/template"
  "net/http"
)

var GlobalTemplateFS embed.FS

type PageData struct {
  Data any
}

//Global template functions available to all templates processed by Render.
var funcMap = template.FuncMap{
  "alphabetRange": func(val string) string {
    switch val {
    case "1": return "A - G"
    case "2": return "H - N"
    case "3": return "O - T"
    case "4": return "U - Z"
    default: return "A - G"
    }
  },
  // Add your new functions right here:
}

func InitTemplates(templateFS embed.FS) {
  GlobalTemplateFS = templateFS
}

func Render(res http.ResponseWriter, layoutName string, templatePaths []string, data PageData) {
  //Create an empty template root, attach functions, then parse the FS files
  tmpl, err := template.New(layoutName).Funcs(funcMap).ParseFS(GlobalTemplateFS, templatePaths...)
  if err != nil {
    http.Error(res, "Template parsing failure: " + err.Error(), http.StatusInternalServerError)
    return
  }
  //Render to a temporary buffer first to catch runtime template errors safely.
  var buf bytes.Buffer
  err = tmpl.ExecuteTemplate(&buf, layoutName, data)
  if err != nil {
    http.Error(res, "Template rendering error: " + err.Error(), http.StatusInternalServerError)
    return
  }
  res.Header().Set("Content-Type", "text/html; charset=utf-8")
  buf.WriteTo(res)
}
