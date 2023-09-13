package entity

import "time"

type Statistic struct {
	Id     int
	Postid int
	Word   string
	Count  int
	Time   time.Time
}
