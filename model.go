package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type User struct {
	Id       string
	Username string `redis:"username"`
	Auth     string `redis:"auth"`
}

type Post struct {
	UserId   string `redis:"user_id"`
	Username string
	Body     string `redis:"body"`
	Elapsed  string `redis:"time"`
}

func (user *User) Is(aUser *User) bool {
	if user == nil || aUser == nil {
		return false
	}
	if user.Id == aUser.Id {
		return true
	}
	return false
}

func (u *User) IsFollowing(p *User) bool {

	v, err := redis.Int(conn.Do("ZSCORE", "following:"+u.Id, p.Id))
	if err != nil {
		return false
	}
	if v > 0 {
		return true
	} else {
		return false
	}
}

func (u *User) Followers() int {
	nbFollowers, err := redis.Int(conn.Do("ZCARD", "followers:"+u.Id))
	if err != nil {
		return 0
	} else {
		return nbFollowers
	}
}

func (u *User) Following() int {
	nbFollowing, err := redis.Int(conn.Do("ZCARD", "following:"+u.Id))
	if err != nil {
		return 0
	} else {
		return nbFollowing
	}
}

func (u *User) Follow(p *User) {
	conn.Do("ZADD", "followers:"+p.Id, time.Now().Unix(), u.Id)
	conn.Do("ZADD", "following:"+u.Id, time.Now().Unix(), p.Id)
}

func (u *User) Unfollow(p *User) {
	conn.Do("ZREM", "followers:"+p.Id, u.Id)
	conn.Do("ZREM", "following:"+u.Id, p.Id)
}
