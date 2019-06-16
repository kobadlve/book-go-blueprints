package main

import (
	"log"

	"gopkg.in/mgo.v2"
)

func main() {
}

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("Dialing MongoDB: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closedb() {
	db.Close()
	log.Println("Close connection")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}
