package model

type Record struct {
	name    string
	age     int
	message string
}

func NewRecord(name string, age int, message string) *Record {
	return &Record{name, age, message}
}
