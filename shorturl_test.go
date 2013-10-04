package main

import (
  "errors"
  "fmt"
  "reflect"
  "testing"
)

var redisModel = NewRedisModel("localhost:6379", "", int64(-1))

func TestNewRedisModel(t *testing.T) {
  expected := "*main.RedisModel"
  result := fmt.Sprintf("%v", reflect.TypeOf(redisModel))

  if result != expected {
    t.Errorf("NewRedisModel() returned %s, expected %s", result, expected)
  }
}

func TestCreate(t *testing.T) {
  validUrl := "http://www.example.com"
  wrongUrl := "www.example"
  expected := "*main.ShortUrl"

  shorted, err := redisModel.Create(validUrl)
  result := fmt.Sprintf("%v", reflect.TypeOf(shorted))
  if result != expected || err != nil {
    t.Errorf("Create('%s') expected (%v, %v), returned (%v, %v)", validUrl, expected, nil, result, err)
  }

  invalidErr := errors.New("data: Invalid url format")
  shorted, err = redisModel.Create(wrongUrl)
  if err == nil {
    t.Errorf("Create('%s') expected '%s' error, returned %s", wrongUrl, invalidErr, err)
  }
}

func TestFindBy(t *testing.T) {
  shorted, err := redisModel.Create("http://www.google.com")
  
  result, err := redisModel.FindBy("id", shorted.Id)
  if (result.Id != shorted.Id && result.Url != shorted.Url) || err != nil {
    t.Errorf("FindBy('id', '%s') expected (%v, %v), returned (%v, %v)", shorted.Id, shorted, nil, result, err)
  }

  result, err = redisModel.FindBy("url", shorted.Url)
  if (result.Id != shorted.Id && result.Url != shorted.Url) || err != nil {
    t.Errorf("FindBy('url', '%s') expected (%v, %v), returned (%v, %v)", shorted.Url, shorted, nil, result, err)
  }

  badId := "0xDEADBEEF"
  expectedErr := errors.New("data: Key not found")

  result, err = redisModel.FindBy("id", badId)
  if result != nil || err == nil {
    t.Errorf("FindBy('id', '%s') expected (%v, %v), returned (%v, %v)", badId, nil, expectedErr, result, err)
  }
}
