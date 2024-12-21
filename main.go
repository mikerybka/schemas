package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"encoding/json"

	"bytes"

	"github.com/mikerybka/types"
	"github.com/oklog/ulid/v2"
)

var dataURL string

func init() {
	dataURL = os.Getenv("DATA_URL")
	if dataURL == "" {
		panic("DATA_URL not set")
	}
}

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			index(w)
		} else {
			show(w, r)
		}
	})
	http.HandleFunc("POST /", create)
	http.HandleFunc("PUT /", update)
	http.HandleFunc("DELETE /", del)
	panic(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter) {
	resp, err := http.Get(dataURL)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		panic(err)
	}
}

func show(w http.ResponseWriter, r *http.Request) {
	url := filepath.Join(dataURL, r.URL.Path)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		panic(err)
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	id := ulid.Make().String()
	s := types.Schema{
		ID: id,
	}

	// Read request body
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write to data store
	url := filepath.Join(dataURL, id)
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.StatusCode != http.StatusOK {
		w.WriteHeader(res.StatusCode)
		_, err = io.Copy(w, res.Body)
		if err != nil {
			panic(err)
		}
		return
	}

	// Return schema object
	err = json.NewEncoder(w).Encode(s)
	if err != nil {
		panic(err)
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	// Read request body
	s := types.Schema{}
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write to data store
	url := filepath.Join(dataURL, r.URL.Path)
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.StatusCode != http.StatusOK {
		w.WriteHeader(res.StatusCode)
		_, err = io.Copy(w, res.Body)
		if err != nil {
			panic(err)
		}
		return
	}

	// Return schema object
	err = json.NewEncoder(w).Encode(s)
	if err != nil {
		panic(err)
	}
}

func del(w http.ResponseWriter, r *http.Request) {}
