package main

import (
  "github.com/hoisie/web"
	"html/template"
  "fmt"
)

var model = NewRedisModel("localhost:6379", "", int64(-1))

func main() {
  web.Get("/", handleIndex)
  web.Post("/", handleCreate)
  web.Get("/(.*)", handleShow)

  web.Run("localhost:9999")
}

type Msg struct {
	Text string
}

func renderTemplate(ctx *web.Context, tmpl string, msg Msg) {
	t, _ := template.ParseFiles("views/" + tmpl + ".html")
  t.Execute(ctx.ResponseWriter, msg)
}

func handleIndex(ctx *web.Context) {
  renderTemplate(ctx, "index", Msg{""})
}

func handleCreate(ctx *web.Context) {
  msg := Msg{""}
  url := ctx.Params["url"]
  short, err := model.Create(url)

  if err == nil {
    msg = Msg{fmt.Sprintf("Saved URL %s with Id %s", short.Url, short.Id)}
  } else {
    msg = Msg{fmt.Sprintf("Error: %s", err)}
  }

  renderTemplate(ctx, "index", msg)
}

func handleShow(ctx *web.Context, id string) {
  res, err := model.FindBy("id", id)

  if err == nil {
    ctx.Redirect(302, res.Url)
    return
  } else {
    renderTemplate(ctx, "index", Msg{fmt.Sprintf("Not Found Id %s", id)})
  }
}
