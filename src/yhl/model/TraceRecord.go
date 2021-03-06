package model

import (
	"time"
)

type TraceRecord struct {
	Ip        string //`json:"ip"`
	Domain    string
	Uri       string //`json:"uri"`
	Datetime  string
	Refer     string
	UserAgent string
	User      interface{}
	Time      time.Time //`json:"time"`
}
