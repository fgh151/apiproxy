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
	errorHandler(http.ListenAndServe(os.Getenv("PORT"), nil))
}

func ProxyServer(w http.ResponseWriter, r *http.Request) {

	req, err := prepareProxyRequest(r)

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
	decodedUrl, err := url.QueryUnescape(queryUrl)

	if err != nil {
		return nil, errors.New("cant decode url")
	}

	newRequest := current
	newRequest.URL, err = url.Parse(decodedUrl)

	if err != nil {
		return nil, errors.New("cant parse url")
	}

	return newRequest, nil
}
