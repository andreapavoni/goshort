package main

import (
	"regexp"
  "fmt"
  "errors"
  "github.com/vmihailenco/redis"

	"math/rand"
	"time"

	"strconv"
  "strings"
)

type RedisModel struct {
  redisPrefix string
  client *redis.Client
}

type ShortUrl struct {
  Id, Url string
}

// Initialize a new RedisModel instance
func NewRedisModel(host, password string, db int64) *RedisModel {
	client := redis.NewTCPClient(host, password, db)
	defer client.Close()

  return &RedisModel{"goshort", client}
}

// Create a new record on db. If url was already shortened, it uses its id
func (m *RedisModel) Create(url string) (*ShortUrl, error) {
  if !validateUrlFormat(url) {
    return nil, errors.New("data: Invalid url format")
  }

  if res, err := m.FindBy("url", url); err == nil {
    return res, nil
  }

  multi, _ := m.client.MultiClient()
  defer multi.Close()

  id := m.generateId()

  _, err := multi.Exec(func() {
    multi.Set(m.redisKey("id", id), url)
    multi.Set(m.redisKey("url", url), id)
  })

  if err != nil {
    return nil, err
  }

  return &ShortUrl{id, url}, nil
}

// Find a record by id or url
func (d *RedisModel) FindBy(key, value string) (*ShortUrl, error) {
  val, err := d.getKey(d.redisKey(key, value))

  if err != nil {
    return nil, err
  }

  var result ShortUrl

  switch key {
  case "id":
    result = ShortUrl{value, val}
  case "url":
    result = ShortUrl{val, value}
  }

  return &result, nil
}

func (d *RedisModel) redisKey(key, value string) string {
  return fmt.Sprintf("%s.%s|%s", d.redisPrefix, key, value)
}

func (d *RedisModel) getKey(key string) (string, error) {
  get := d.client.Get(key)

  if get.Err() != nil {
    return "", get.Err()
  }

  val := get.Val()

  if val == "" {
    err := fmt.Sprintf("data: Key not found %s", key)
    return "", errors.New(err)
  }

  return val, nil
}

// Misc utils

func (m *RedisModel) generateId() string {
  for  {
    id := base62(randInt())

    if  _, err := m.FindBy("id", id); err != nil {
      return id
    }
  }

  return ""
}

func randInt() int {
	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)

	val := int(seed)
	if val < 0 {
		val = -val
	}

	return (1 + rand.Intn(val))
}

func base62(num int) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	enc := make([]uint8, len(string(strconv.Itoa(num))))

	if num == 0 {
		return "0"
	}

	for i := 0; num > 0; i++ {
		n := num % 62

		enc[i] = chars[n]
		num /= 62
	}

  var n uint8
  return strings.Trim(string(enc), string(n))

}

func validateUrlFormat(url string) bool {
	const pattern = "(^$)|(^(http|https)://[a-z0-9]+([-.]{1}[a-z0-9]+)*.[a-z]{2,5}(([0-9]{1,5})?/.*)?$)"

  if url == "" {
    return false
  }

	result, _ := regexp.MatchString(pattern, url)

	return result
}
