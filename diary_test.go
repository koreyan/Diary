package main

import (
	"bytes"
	"diary/network"
	"diary/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDiary(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/diary/1", nil)

	mux := network.MakeHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	var diary1 types.Memo
	err := json.NewDecoder(res.Body).Decode(&diary1)
	assert.Nil(err)
	assert.Equal(1, diary1.Id)
	assert.Equal("title1", diary1.Title)
	assert.Equal("content1", diary1.Content)

	res = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/diary/10", nil)
	mux = network.MakeHandler()
	mux.ServeHTTP(res, req)
	assert.Equal(http.StatusNotFound, res.Code)
	var diary types.Memo
	err = json.NewDecoder(res.Body).Decode(&diary)
	assert.Nil(err)
	assert.Equal(0, diary.Id)
	assert.Equal("", diary.Title)
	assert.Equal("", diary.Content)
}

func TestGetDiaryList(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/diary", nil)

	mux := network.MakeHandler()
	mux.ServeHTTP(res, req)

	var memos types.Memos
	err := json.NewDecoder(res.Body).Decode(&memos)
	assert.Nil(err)
	assert.Equal(0, memos[1].Id)
	assert.Equal("title0", memos[1].Title)
	assert.Equal("content0", memos[1].Content)

	assert.Equal(1, memos[0].Id)
	assert.Equal("title1", memos[0].Title)
	assert.Equal("content1", memos[0].Content)

}

func TestPostDiary(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	var data struct {
		Title   string
		Content string
	} = struct {
		Title   string
		Content string
	}{Title: "오늘은 멋진 날이다.", Content: "왜냐면 오늘 살아있기 때문이다."}

	jsonData, err := json.Marshal(data)
	requestBody := bytes.NewBuffer(jsonData)

	assert.Nil(err)

	req := httptest.NewRequest("POST", "/diary", requestBody)

	mux := network.MakeHandler()
	mux.ServeHTTP(res, req)
	assert.Equal(http.StatusCreated, res.Code)
	var memo types.Memo
	err = json.NewDecoder(res.Body).Decode(&memo)
	assert.Nil(err)
	assert.Equal(2, memo.Id)
	assert.Equal("오늘은 멋진 날이다.", memo.Title)
	assert.Equal("왜냐면 오늘 살아있기 때문이다.", memo.Content)
}
