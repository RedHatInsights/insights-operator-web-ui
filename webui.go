package main

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func getContentType(filename string) string {
	// TODO: to map
	if strings.HasSuffix(filename, ".html") {
		return "text/html"
	} else if strings.HasSuffix(filename, ".js") {
		return "application/javascript"
	} else if strings.HasSuffix(filename, ".css") {
		return "text/css"
	}
	return "text/html"
}

func sendStaticPage(writer http.ResponseWriter, filename string) {
	body, err := ioutil.ReadFile(filename)
	if err == nil {
		writer.Header().Set("Server", "A Go Web Server")
		writer.Header().Set("Content-Type", getContentType(filename))
		fmt.Fprint(writer, string(body))
	} else {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
	}
}

func staticPage(filename string) func(writer http.ResponseWriter, request *http.Request) {
	log.Println("Serving static file", filename)
	return func(writer http.ResponseWriter, request *http.Request) {
		sendStaticPage(writer, filename)
	}
}

func startHttpServer() {
	http.HandleFunc("/", staticPage("html/index.html"))
	http.HandleFunc("/bootstrap.min.css", staticPage("html/bootstrap.min.css"))
	http.HandleFunc("/bootstrap.min.js", staticPage("html/bootstrap.min.js"))
	http.HandleFunc("/ccx.ccs", staticPage("html/ccs.ccs"))
	http.ListenAndServe(":8888", nil)
}

func main() {
	log.Println("Reading configuration")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	log.Println("Starting the service")
	startHttpServer()
}
