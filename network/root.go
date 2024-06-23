package network

import (
	"diary/types"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// api
func MakeHandler() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/diary/{id:[0-9]+}", GetDiary).Methods("GET")
	mux.HandleFunc("/diary", GetDiaryList).Methods("GET")
	mux.HandleFunc("/diary", PostDiary).Methods("POST")
	types.MemoMap = make(map[int]types.Memo)
	types.MemoMap[0] = types.Memo{Id: 0, Title: "title0", Content: "content0"}
	types.MemoMap[1] = types.Memo{Id: 1, Title: "title1", Content: "content1"}
	types.DiaryId = 2
	return mux
}

func GetDiary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	v, ok := types.MemoMap[id]
	if ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(v)

	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.Memo{})
	}

}

func GetDiaryList(w http.ResponseWriter, r *http.Request) {
	var memos types.Memos = make([]types.Memo, 0)
	for _, v := range types.MemoMap {
		memos = append(memos, v)
	}

	sort.Sort(memos)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memos)
}

func PostDiary(w http.ResponseWriter, r *http.Request) {
	var diary types.Memo
	err := json.NewDecoder(r.Body).Decode(&diary)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	diary.Id = types.DiaryId // id 부여
	types.DiaryId++
	diary.Year, diary.Month, diary.Day = time.Now().Date() // 등록시간 부여

	types.MemoMap[diary.Id] = diary // 등록

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(diary)
}
