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
var nextDelTime string

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/data", dataHandle)
	http.HandleFunc("/delTime", delTimeHandle)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		provide fake data ！！
		router:[/data]
		method:[POST] {
			header:[key] or url parameter[key]
			body:what you want to get
		}

		method:[GET] {
			header:[key] or url parameter[key]
			reponse:what you post by key
		}

		method:[delete] {
			header:[key] or url parameter[key]
			reponse:delete data by [key]
		}

		======================================
		router:[/delTime]
		method:[GET] {
			reponse:next delete time
		}

		`)
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
		now := time.Now().Add(24 * time.Hour)
		nextDelTime = now.Format("2006/01/02/ 15:04")
		time.Sleep(24 * time.Hour)
	}
}
func delTimeHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte("next delete time : " + nextDelTime))
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
