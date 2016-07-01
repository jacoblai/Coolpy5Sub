package Redico

import (
	"github.com/garyburd/redigo/redis"
	"testing"
	"fmt"
)

func TestRedico(t *testing.T) {
	s, err := Run()
	if err != nil {
		t.Error(err)
	}
	defer s.Close()

	// Configure you application to connect to redis at s.Addr()
	// Any redis client should work, as long as you use redis commands which
	c, err := redis.Dial("tcp", s.Addr())
	if err != nil {
		t.Error(err)
	}

	_, err = c.Do("AUTH", "foo", "bar")
	fmt.Println(err != nil, "no password set")

	s.RequireAuth("nocomment")
	_, err = c.Do("PING", "foo", "bar")
	fmt.Println(err != nil, "need AUTH")

	_, err = c.Do("AUTH", "wrongpasswd")
	fmt.Println(err != nil, "wrong password")

	_, err = c.Do("AUTH", "nocomment")
	fmt.Println(err)

	_, err = c.Do("PING")
	fmt.Println(err)

	_, err = c.Do("SET", "foo", "bar")
	if err != nil {
		t.Error(err)
	}

	_, err = c.Do("SET", "joo", "bar")
	if v, err := redis.Strings(c.Do("KEYS", "j*"));err !=nil || v[0] != "joo"{
		t.Error("Keys not fire *")
	}

	// You can ask about keys directly, without going over the network.
	if got, err := s.Get("foo"); err != nil || got != "bar" {
		t.Error("Didn't get 'bar' back")
	}

	// Or with a DB id
	if _, err := s.DB(42).Get("foo"); err != ErrKeyNotFound {
		t.Error("Didn't use a different DB")
	}
	// Or use a Check* function which Fail()s if the key is not what we expect
	// (checks for existence, key type and the value)
	// s.CheckGet(t, "foo", "bar")

	// Check if there really was only one connection.
	if s.TotalConnectionCount() != 1 {
		t.Error("Too many connections made")
	}
}