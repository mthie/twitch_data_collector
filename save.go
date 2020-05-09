package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

func saveContent(kind, filename, content string) {
	p := path.Join(kind, filename+".txt")
	if err := os.MkdirAll(kind, 0777); err != nil {
		log.WithField("path", p).WithError(err).Error("Unable to create directory")
		return
	}

	if err := ioutil.WriteFile(p, []byte(content), 0644); err != nil {
		log.WithError(err).Error("Unable to write content")
	}
}

func saveJSON(kind, filename string, data interface{}) {
	p := path.Join(kind, filename+".json")
	if err := os.MkdirAll(kind, 0777); err != nil {
		log.WithField("path", p).WithError(err).Error("Unable to create directory")
		return
	}

	f, err := os.Create(p)
	if err != nil {
		log.WithField("path", p).WithError(err).Error("Unable to create file")
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(data); err != nil {
		log.WithField("path", p).WithError(err).Error("Unable to encode json")
		return
	}
}
