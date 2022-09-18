package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

//Функци обработчик ошибок
func errorHandler(err error) {
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", ProxyServer)
	errorHandler(http.ListenAndServe(os.Getenv("SERVER"), nil))
}

func ProxyServer(w http.ResponseWriter, r *http.Request) {

	req, err := prepareProxyRequest(r)

	if err != nil {
		throw500(w, err)
		return
	}

	resp, err := doRequest(req)

	if err != nil {
		throw500(w, err)
		return
	}

	sendResponse(w, resp)
}

func sendResponse(w http.ResponseWriter, resp *http.Response) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(bodyBytes)
}

func doRequest(r *http.Request) (*http.Response, error) {
	return cacheable(r, func() (*http.Response, error) {
		client := &http.Client{}
		resp, err := client.Do(r)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})
}

type RequestCache struct {
	Resp http.Response
	Err  error
}

func (c RequestCache) Format() (*http.Response, error) {
	return &c.Resp, c.Err
}

func cacheable(r *http.Request, fn func() (*http.Response, error)) (*http.Response, error) {

	if os.Getenv("CACHE") == "true" {
		key := Hash(r)
		c := cache.New(5*time.Minute, 10*time.Minute)

		itemFromCache, found := c.Get(key)
		if found {
			return itemFromCache.(RequestCache).Format()
		}

		resp, err := fn()

		cacheItem := RequestCache{Resp: *resp, Err: err}

		c.Set(key, cacheItem, cache.NoExpiration)

		return resp, err
	}

	return fn()
}

func Hash(s *http.Request) string {
	var b bytes.Buffer
	_ = gob.NewEncoder(&b).Encode(s)
	return string(b.Bytes())
}

func throw500(w http.ResponseWriter, err error) {
	errorHandler(err)
	w.WriteHeader(500)
	_, _ = w.Write([]byte("Err :" + err.Error()))
}

func prepareProxyRequest(current *http.Request) (*http.Request, error) {

	queryUrl := current.URL.Query().Get("url")

	if queryUrl == "" {
		return nil, errors.New("empty url")
	}

	decodedUrl, err := url.QueryUnescape(queryUrl)

	if err != nil {
		return nil, errors.New("cant decode url")
	}

	if err != nil {
		return nil, errors.New("cant parse url")
	}

	req, err := http.NewRequest(current.Method, decodedUrl, current.Body)

	if err != nil {
		return nil, errors.New("cant create request")
	}

	for name, values := range current.Header {
		for _, value := range values {
			req.Header.Set(name, value)
		}
	}

	return req, nil
}
