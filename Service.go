package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	fmt.Println("Welcome to Go Status checker:")
	muxRouter := mux.NewRouter()

	//start status check for url
	muxRouter.HandleFunc("/websites", StatusCheck).Methods(http.MethodPost)
	// fetch all with status
	muxRouter.HandleFunc("/websites", FetchAll).Methods(http.MethodGet)
	//find  by url
	// muxRouter.HandleFunc("/websites/{name}", FindByUrl).Methods(http.MethodGet)

	server := negroni.Classic()
	server.UseHandler(muxRouter)
	server.Run(":3000")
}

type StatusChecker struct {
	Url    string `json:"url"`
	Status string `json:"status"`
}

var list = []StatusChecker{}

func StatusCheck(w http.ResponseWriter, r *http.Request) {
	var arr = []string{}
	err := json.NewDecoder(r.Body).Decode(&arr)
	if err == nil {
		go func() {
			for {
				for _, url := range arr {
					go func(url string) {
						s := StatusChecker{url, ""}
						list = append(list, s)

						resp, err := http.Get(url)

						if err == nil {
							fmt.Printf("Site Url: %v ,Status %v", url, resp.StatusCode)
							if resp.StatusCode == 200 {
								updateStatus(url, "UP")
							}
						} else {
							updateStatus(url, "DOWN")
							fmt.Printf("Exception:%v", err)
						}
					}(url)
				}
				time.Sleep(time.Minute * 1)
			}
		}()

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
	}
}

func updateStatus(url, status string) {
	for index, value := range list {
		fmt.Println(value)
		if value.Url == url {
			val := value
			val.Status = status
			list[index] = val
			break
		}
	}

}

// if query param present return specific else all
func FetchAll(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("name")
	fmt.Println(url)
	if url != "" {
		FindByUrl(w, r, url)
	} else {
		data, err := json.Marshal(list)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	}
}
func FindByUrl(w http.ResponseWriter, r *http.Request, str string) {
	// param := mux.Vars(r)
	// str := param["name"]
	fmt.Println("name:", str)
	var user StatusChecker
	for _, usr := range list {
		if usr.Url == str {
			user = usr
			break
		}
	}
	data, err := json.Marshal(user)

	if err == nil {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(data))
	}
}
