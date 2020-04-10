package main


import (
	"fmt"
	"github.com/tkanos/gonfig"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

//Функци обработчик ошибок
func errorHandler(err error)  {
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}

//Структура файла конфигурации
type Configuration struct {
	Port         string
}

func main() {
	configuration := Configuration{}

	errorHandler(gonfig.GetConf(os.Args[1], &configuration))

	http.HandleFunc("/", ProxyServer)
	fmt.Println("Listen port "+configuration.Port)

	errorHandler(http.ListenAndServe(configuration.Port, nil))
}

func ProxyServer(w http.ResponseWriter, r *http.Request) {


	keys, ok := r.URL.Query()["url"]
//Получаем запрашиваемуй адрес
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'url' is missing")
		return
	}
	//Валидация адреса
	hostUrl, err := url.ParseRequestURI(keys[0])
	if err != nil {
		log.Println("Url Param 'url' is invalid")
		return
	}
	//Шлем запрос
	resp, err := http.Get(hostUrl.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	//Возвращаем ответ
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		fmt.Fprintf(w, bodyString)
	}
}
