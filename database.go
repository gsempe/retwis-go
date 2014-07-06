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
	// registerScript is used to register a new user atomically.
	// Keys and arguments:
	// KEYS[1] == "users"
	// KEYS[2] == "user:"+userId
	// KEYS[3] == "auths"
	// KEYS[4] == "users_by_time"
	// ARGV[1] == the userId
	// ARGV[2] == the username wanted
	// ARGV[3] == the password used
	// ARGV[4] == the auth token to use
	// ARGV[5] == current timestamp
	// Return the authentification key or an error
	registerScript = redis.NewScript(4, `
local username = redis.call("HGET", KEYS[1], ARGV[2])
if username then
	return redis.error_reply("Sorry the selected username is already in use.")
end
redis.call("HSET", KEYS[1], ARGV[2], ARGV[1])
redis.call("HMSET", KEYS[2], "username", ARGV[2], "password", ARGV[3], "auth", ARGV[4])
redis.call("HSET", KEYS[3], ARGV[4], ARGV[1])
redis.call("ZADD", KEYS[4], ARGV[5], ARGV[2])
return ARGV[4]
		`)
)

func init() {

	var err error
	conn, err = redis.Dial("tcp", address)
	if nil != err {
		log.Fatalln("Error: Connection to redis:", err)
	}
	registerScript.Load(conn)
}
