package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var dataMap sync.Map

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/data", dataHandle)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "provide fake data ！！")
	})
	if port == "" {
		port = "8888"
	}
	go clearMap()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func clearMap() {
	for {
		dataMap.Range(func(key interface{}, value interface{}) bool {
			dataMap.Delete(key)
			return true
		})
		time.Sleep(24 * time.Hour)
	}
}

func dataHandle(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	key := keys.Get("key")
	if key == "" {
		key = r.Header.Get("key")
	}
	if r.Method == "POST" {
		if byteArr, readErr := ioutil.ReadAll(r.Body); readErr != nil {
			w.Write([]byte("read error : " + readErr.Error()))
		} else {
			dataMap.Store(key, string(byteArr))
		}
	} else if r.Method == "GET" {

		if v, ok := dataMap.Load(key); !ok {
			w.Write([]byte("no data with:`" + key + "`"))
		} else {
			val := v.(string)
			w.Write([]byte(val))
		}
	} else if r.Method == "DELETE" {
		dataMap.Delete(key)
		w.Write([]byte("delete `" + key + "` succeed"))
	}
}
