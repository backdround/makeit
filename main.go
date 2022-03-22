package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
  MakeItFile = "makeit.yaml"
)

type Task struct {
  Script interface{}
}

type MakeItData struct {
  Tasks map[string]Task
}


func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func main() {
  file, err := os.Open(MakeItFile)
  check(err)

  fileData, err := ioutil.ReadAll(file)
  check(err)

  makeItData := MakeItData{
    Tasks: make(map[string]Task),
  }

  err = yaml.Unmarshal(fileData, makeItData)
  check(err)

  fmt.Println(makeItData)
}
