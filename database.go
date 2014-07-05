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
	// Warning: retwis-go is not usable in a redis cluster because of
	// of this script and the key `users .. userId` built in the LUA script
	// Keys and arguments:
	// KEYS[1] == "next_user_id"
	// ARGV[1] == the username wanted
	// ARGV[2] == the password used
	// ARGV[3] == the auth token to use
	// ARGV[4] == current timestamp
	// Return the authentification key or an error
	registerScript = redis.NewScript(1, `
local username = redis.call("HGET", "users", ARGV[1])
print("valeur de username")
print(username)
if username then
	return redis.error_reply("Sorry the selected username is already in use.")
end
local userId = redis.call("INCR",KEYS[1])
redis.call("HSET", "users", ARGV[1], userId)
redis.call("HMSET", "user:" .. userId, "username", ARGV[1], "password", ARGV[2], "auth", ARGV[3])
redis.call("HSET", "auths", ARGV[3], userId)
redis.call("ZADD", "users_by_time", ARGV[4], ARGV[1])
return ARGV[3]
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
