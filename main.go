package main

import (
  "flag"
  "fmt"
  "github.com/pilu/traffic"
  "log"
  "net/http"
  "net/url"
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
  router.Get("/:id", handleShow)

  model, _ = initRedis(*redisconf)

  address := fmt.Sprintf("%s:%s", *host, *port)
  err := http.ListenAndServe(address, router)

  if err != nil {
    log.Fatal(err)
  }
}

type ResponseData struct {
  Message string
}

func initRedis(config string) (*RedisModel, error) {
  var password string

  u, _ := url.Parse(config)

  host := u.Host

  if auth := u.User; auth != nil {
    password, _ = u.User.Password()
  } else {
    password = ""
  }

  return NewRedisModel(host, password, int64(-1)), nil
}

func handleIndex(w traffic.ResponseWriter, r *http.Request) {
  traffic.Render(w, "index", ResponseData{""})
}

func handleCreate(w traffic.ResponseWriter, r *http.Request) {
  msg := ResponseData{""}
  param_url := r.FormValue("url")
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
    msg = ResponseData{fmt.Sprintf("Shortened URL from %s to %s://%s/%s", short.Url, scheme, r.Host, short.Id)}
  } else {
    msg = ResponseData{fmt.Sprintf("Error: %s", err)}
  }

  traffic.Render(w, "index", msg)
}

func handleShow(w traffic.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  id := params.Get("id")
  log.Printf("Request ID: %s", id)
  
  res, err := model.FindBy("id", id)


  if err == nil {
    http.Redirect(w, r, res.Url, 302)
    return
  } else {
    traffic.Render(w, "index", ResponseData{fmt.Sprintf("Not Found Id %s", id)})
  }
}
