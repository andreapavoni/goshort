# GoShort

An experiment to make a simple web app using [Go](http://golang.org) and [Redis](http://redis.io)

## Installation

You'll need a working installation of Redis and Go, then:

* ```go get github.com/pilu/traffic```
* ```go get github.com/apeacox/goshort```
* ```cd $GOPATH/src/github.com/apeacox/goshort```
* ```go build```
* start a Redis server
* ```./goshort```
* point your browser to http://0.0.0.0:8080

### Command line options:

GoShort accepts the following command line options:

* ```-host <hostname>```: hostname to listen (default ```0.0.0.0```)
* ```-p <portNumber>```: port to listen (default ```8080```)
* ```-redis <redis://[user:pass@]host:port>``` specify Redis connection URL (default: ```redis://localhost:6379/```)

These are especially useful for deploying in production, check the [Procfile](https://github.com/apeacox/goshort/blob/master/Procfile) to see how it's used on Heroku.

## Demo

There's a working demo on Heroku: http://goshort.herokuapp.com

## Contributing

1. Fork it!
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

### Testing ?

I'm still figuring out how to write them :-P

## License

Copyright (c) 2013 Andrea Pavoni http://andreapavoni.com
