package main

import (
  "github.com/hoisie/web"
	"html/template"
  "fmt"
  "net/url"
  "flag"
)


var host *string = flag.String("host", "0.0.0.0", "host and to listen (eg: localhost)")
var port *string = flag.String("port", "8080", "port to listen")
var redisconf *string = flag.String("redis", "redis://localhost:6379/", "URL for redis Server (eg: redis://[user:pass@]localhost:6379)")


var model *RedisModel

func main() {
  web.Get("/", handleIndex)
  web.Post("/", handleCreate)
  web.Get("/(.*)", handleShow)

  flag.Parse()

  model, _ = initRedis(*redisconf)
  listen := fmt.Sprintf("%s:%s", *host, *port)

  web.Run(listen)
}

type Msg struct {
	Text string
}

func renderTemplate(ctx *web.Context, tmpl string, msg Msg) {
	t, _ := template.ParseFiles("views/" + tmpl + ".html")
  t.Execute(ctx.ResponseWriter, msg)
}

func initRedis(config string) (*RedisModel, error) {
  var password string

  u, _ := url.Parse(config)

  host := u.Host

  if auth := u.User ; auth != nil {
    password, _ = u.User.Password()
  } else {
    password = ""
  }

  return NewRedisModel(host, password, int64(-1)), nil
}

func handleIndex(ctx *web.Context) {
  renderTemplate(ctx, "index", Msg{""})
}

func handleCreate(ctx *web.Context) {
  msg := Msg{""}
  url := ctx.Params["url"]
  short, err := model.Create(url)

  if err == nil {
    msg = Msg{fmt.Sprintf("Shortened URL from %s to %s/%s", short.Url, ctx.Request.Host, short.Id)}
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
