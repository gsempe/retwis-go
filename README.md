retwis-go
=========
## What is it
retwis-go is the port to Go of the redis tutorial [Twitter clone](http://redis.io/topics/twitter-clone)  
You can see it in action at <a href="http://retwis.sempe.net" target="_blank">http://retwis.sempe.net</a>

## How to contribute
retwis-go is a direct port without almost no improvements done on the go. I have done it as a way to practice Golang and redis.
If you have the same goals and search projects to practice, just fork it and open pull requests.
There is lot of things that can be done:
- protect users passwords. They are not encrypted in the database
- Many errors messages should not be shown to the retwis user
- User profile is very poor. For instance there is no way to know who follows who
- There is no reply feature, no retweet feature, no favorite feature, etc... 

## Usage
Set the `GOPATH` env, like described in the [Golang documentation](http://golang.org/doc/code.html#GOPATH)

Install retwis-go:
```
go get github.com/gsempe/retwis-go
```

Run retwis-go:
```
go build
./retwis-go
```

Note: A redis database must run on the same machine

## How it's done
To get it done faster and safer the port uses :
- [redigo](https://github.com/garyburd/redigo) Go client for Redis
- [negroni](https://github.com/codegangsta/negroni) Idiomatic HTTP Middleware for Golang
- [httprouter](https://github.com/julienschmidt/httprouter) A high performance HTTP request router that scales well
- [render](https://github.com/unrolled/render) Go package for easily rendering JSON, XML, and HTML template responses.
- [securecookie](https://github.com/gorilla/securecookie) Gorilla package that encodes and decodes authenticated and optionally encrypted cookie values.
