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

// CORS 미들웨어: 모든 요청에 CORS 헤더 추가
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 모든 도메인에서의 접근을 허용
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 허용되는 메소드
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// 허용되는 헤더
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// preflight 요청에 대한 처리
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 다음 핸들러로 요청 전달
		next.ServeHTTP(w, r)
	})
}

// api
func MakeHandler() http.Handler {
	mux := mux.NewRouter()
	// 미들웨어 추가 => 모든 HandleFunc가 실행되기 전 미들웨어부터 실행됨
	mux.Use(corsMiddleware)
	mux.HandleFunc("/diary/{id:[0-9]+}", GetDiary).Methods("GET")
	mux.HandleFunc("/diary", GetDiaryList).Methods("GET")
	mux.HandleFunc("/diary", PostDiary).Methods("POST")
	mux.HandleFunc("/diary/{id:[0-9]+}", DeleteDiary).Methods("DELETE")

	types.MemoMap = make(map[int]types.Memo)
	// id는 1부터 시작해야함, 0으로 시작하면 웹에서 id key를 삭제함
	types.MemoMap[1] = types.Memo{Id: 1, Title: "title1", Content: "content1"}
	types.MemoMap[2] = types.Memo{Id: 2, Title: "title2", Content: "content2"}
	types.DiaryId = 3
	return mux
}

func GetDiary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	v, ok := types.MemoMap[id]
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(v)

	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(diary)
}

func DeleteDiary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 조회
	if _, ok := types.MemoMap[id]; ok {
		// 삭제
		delete(types.MemoMap, id)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
