package main

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

const (
	address = "127.0.0.1:6379"
)

var (
	conn redis.Conn
)

func init() {

	var err error
	conn, err = redis.Dial("tcp", address)
	if nil != err {
		log.Fatalln("Error: Connection to redis:", err)
	}
}
