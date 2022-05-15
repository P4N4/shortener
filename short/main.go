package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/wzshiming/ssdb"
)

var DB, err = ssdb.Connect(
	ssdb.URL("ssdb://127.0.0.1:8888"),
)

const rnd = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rnd[rand.Intn(len(rnd))]
	}
	return string(b)
}
func urlShortener(w http.ResponseWriter, r *http.Request) {

	url := string(r.URL.Query().Get("url"))

	rand := randStringBytes(8)

	x, err := DB.SetNX(url, ssdb.Value(rand))
	if err != nil {
		fmt.Println(err)
	}

	s, err := DB.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	if x == true {
		DB.Set(rand, ssdb.Value(url))
		if err != nil {
			fmt.Println(err)
		}
	}

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(s))

	return
}

func urlRedirect(w http.ResponseWriter, r *http.Request) {
	url := string(r.URL.RequestURI())
	value := url[3:]

	s, err := DB.Get(value)
	if err != nil {
		fmt.Println(err)
	}

	newStringValue := string(s)

	http.Redirect(w, r, newStringValue, http.StatusMovedPermanently)
	return
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/a/", urlShortener)
	mux.HandleFunc("/s/", urlRedirect)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error creating http server: ", err)
	}
}
