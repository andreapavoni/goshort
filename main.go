package main

import (
  "net/http"
  "github.com/pilu/traffic"
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
  flag.Parse()

  router := traffic.New()

  router.Get("/", handleIndex)
  router.Post("/", handleCreate)
  router.Get("/(.*)", handleShow)

  http.Handle("/", router)

  model, _ = initRedis(*redisconf)
  listen := fmt.Sprintf("%s:%s", *host, *port)

  http.ListenAndServe(listen, nil)
}

type Msg struct {
	Text string
}

func renderTemplate(w http.ResponseWriter, tmpl string, msg Msg) {
	t, _ := template.ParseFiles("views/" + tmpl + ".html")
  t.Execute(w, msg)
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

func handleIndex(w http.ResponseWriter, r *http.Request) {
  renderTemplate(w, "index", Msg{""})
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
  msg := Msg{""}
  param_url := r.FormValue( "url" )
  short, err := model.Create(param_url)
  var scheme string

  /*
    A dirty hack to check wether the app is hosted on https or not.
    This is the only way to guess, because URI.Scheme is empty.
    Read here for details:
    http://stackoverflow.com/questions/6899069/why-are-request-url-host-and-scheme-blank-in-the-development-server
  */
  if r.TLS != nil {
    scheme = "https"
  } else {
    scheme = "http"
  }

  if err == nil {
    msg = Msg{fmt.Sprintf("Shortened URL from %s to %s://%s/%s", short.Url, scheme, r.Host, short.Id)}
  } else {
    msg = Msg{fmt.Sprintf("Error: %s", err)}
  }

  renderTemplate(w, "index", msg)
}

func handleShow(w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  id := params.Get("id")
  res, err := model.FindBy("id", id)

  if err == nil {
    http.Redirect(w, r, res.Url, 302)
    return
  } else {
    renderTemplate(w, "index", Msg{fmt.Sprintf("Not Found Id %s", id)})
  }
}
