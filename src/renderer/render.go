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
  "strconv"
)

var GlobalTemplateFS embed.FS

type PageData struct {
  Data any
}

/***
Global template functions available to all templates processed by Render.

Global Level Configuration
When you define your functions globally, you typically attach them to a central variable or map once when the server first
boots up (usually in main() or an init() block).
* How it works: Go compiles and saves the helper logic into memory a single time during initialization.
* The Big Advantage: Any HTML template page processed by your web application can call get_pagination_pages at any time.
  You do not have to copy and paste registration code for every individual page router endpoint.
* Performance: It is faster because Go doesn't waste time rebuilding the mathematical function layout rules every time a
  user clicks a button or reloads a page.

Function Level Configuration
When you define functions inside a local handler (like our renderTemplate example), you attach them to the template instance
directly within that specific HTTP request path block.
* How it works: Every single time a user requests that page, Go allocates fresh system memory to re-register the helper
  function logic, processes the layout, and then destroys it when the request completes.
* The Big Advantage: It allows you to build unique helpers that require direct access to a specific request's data values
  (like reading a user's logged-in session cookies or parsing URL queries).
* The Catch: No other pages on your site can see that helper function. If you need pagination on an "Accounts" page and a
  "Users" page, you have to write out the configuration code inside both function routers.
***/
var funcMap = template.FuncMap{
  "alphabetRange": func(val string, labels []string) string {
    //If the labels array is empty, fall back safely.
    if len(labels) == 0 {
      return ""
    }
    //Convert the slider's string value (e.g., "2") into an integer index; use Atoi because the slider values come
    //through as strings.
    index, err := strconv.Atoi(val)
    if err != nil {
      return labels[0]  //Fallback to the first item if parsing fails.
    }
    //The slider starts at index 1, but Go slices start at index 0; subtract 1.
    sliceIndex := index - 1
    //Prevent index out of bounds.
    if -1 < sliceIndex && sliceIndex < len(labels) {
      return labels[sliceIndex]
    }
    return labels[0]
  },
  /***
  Go templates do not naturally let you iterate over a numeric loop (like 1 to totalPages) out of the box. To make an
  HTML loop code work seamlessly, simply attach a quick sequence helper function to your template functions right before
  parsing the file.
  ***/
  "get_pagination_pages": func(current, total int) []int {
    var pages []int
    if total <= 7 {
      for i := 1; i <= total; i++ {
        pages = append(pages, i)
      }
      return pages
    }
    pages = append(pages, 1)  //Always show page 1
    if current <= 4 {
      for i := 2; i <= 5; i++ {
        pages = append(pages, i)
      }
      pages = append(pages, -1)  //-1 triggers the ellipsis (...) branch
    } else if current >= (total - 3) {
      pages = append(pages, -1)
      for i := total - 4; i <= (total - 1); i++ {
        pages = append(pages, i)
      }
    } else {
      pages = append(pages, -1)
      pages = append(pages, current - 1, current, current + 1)
      pages = append(pages, -1)
    }
    pages = append(pages, total)  //Always show the last page
    return pages
  },
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
