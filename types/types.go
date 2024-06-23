package types

import (
	"time"
)

type Memo struct {
	Id      int    `json:"Id,omitempty"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
	Date    `json:"Date,omitempty"`
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

var DiaryId int

var MemoMap map[int]Memo

type Memos []Memo

func (m Memos) Len() int {
	return len(m)
}

func (m Memos) Less(i, j int) bool {
	return m[i].Id > m[j].Id
}

func (m Memos) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
