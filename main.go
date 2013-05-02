package main

import (
  "github.com/hoisie/web"
	"html/template"
  "fmt"
  "regexp"
  "os"
)


var model, _ = initRedis()

func main() {
  web.Get("/", handleIndex)
  web.Post("/", handleCreate)
  web.Get("/(.*)", handleShow)

  listen := ":" + os.Getenv("PORT")

  if listen == "" {
    listen = ":9999"
  }

  web.Run(listen)
}

type Msg struct {
	Text string
}

func renderTemplate(ctx *web.Context, tmpl string, msg Msg) {
	t, _ := template.ParseFiles("views/" + tmpl + ".html")
  t.Execute(ctx.ResponseWriter, msg)
}

func initRedis() (*RedisModel, error) {

  if creds := os.Getenv("REDISTOGO_URL") ; creds != "" {
    re, err := regexp.Compile(`redis://(\w+:\w+)?@(.*:\d+)/`)

    if err != nil {
      return nil, err
    }

    c := re.FindStringSubmatch(creds)[1:]
    return NewRedisModel(c[1], c[0], int64(-1)), nil
  }

  return NewRedisModel("localhost:6379", "", int64(-1)), nil
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
