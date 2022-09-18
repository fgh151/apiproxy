package main

import (
	"errors"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Err :" + err.Error()))
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
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
