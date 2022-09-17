package main

import (
	"errors"
	"fmt"
	"github.com/tkanos/gonfig"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

//Функци обработчик ошибок
func errorHandler(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}

//Структура файла конфигурации
type Configuration struct {
	Port string
}

func main() {
	configuration := Configuration{}

	errorHandler(gonfig.GetConf(os.Args[1], &configuration))

	http.HandleFunc("/", ProxyServer)
	fmt.Println("Listen port " + configuration.Port)

	errorHandler(http.ListenAndServe(configuration.Port, nil))
}

func ProxyServer(w http.ResponseWriter, r *http.Request) {

	req, err := prepareProxyRequest(r)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(bodyBytes)

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
