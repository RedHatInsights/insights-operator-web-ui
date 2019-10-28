package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const API_PREFIX = "/api/v1/"

type Cluster struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ConfigurationProfile struct {
	Id            int    `json:"id"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Description   string `json:"description"`
}

type ClusterConfiguration struct {
	Id            int    `json:"id"`
	Cluster       string `json:"cluster"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Active        string `json:"active"`
	Reason        string `json:"reason"`
}

var controllerURL = ""

func performReadRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}
	body, readErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		return nil, fmt.Errorf("Unable to read response body")
	}

	return body, nil
}

func readListOfClusters(controllerUrl string, apiPrefix string) ([]Cluster, error) {
	clusters := []Cluster{}

	url := controllerUrl + apiPrefix + "client/cluster"
	body, err := performReadRequest(url)

	err = json.Unmarshal(body, &clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func readListOfConfigurationProfiles(controllerUrl string, apiPrefix string) ([]ConfigurationProfile, error) {
	profiles := []ConfigurationProfile{}

	url := controllerUrl + apiPrefix + "client/profile"
	body, err := performReadRequest(url)

	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func readListOfConfigurations(controllerUrl string, apiPrefix string) ([]ClusterConfiguration, error) {
	configurations := []ClusterConfiguration{}

	url := controllerUrl + apiPrefix + "client/configuration"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &configurations)
	if err != nil {
		return nil, err
	}
	return configurations, nil
}

func readConfigurationProfile(controllerUrl string, apiPrefix string, profileId string) (*ConfigurationProfile, error) {
	var profile ConfigurationProfile
	url := controllerUrl + apiPrefix + "client/profile/" + profileId
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func readClusterConfigurationById(controllerUrl string, apiPrefix string, configurationId string) (*string, error) {
	url := controllerUrl + apiPrefix + "client/configuration/" + configurationId
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	str := string(body)
	return &str, nil
}

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

type ListClustersDynContent struct {
	Items []Cluster
}

func listClusters(writer http.ResponseWriter, request *http.Request) {
	clusters, err := readListOfClusters(controllerURL, API_PREFIX)
	if err != nil {
		log.Println("Error reading list of clusters", err)
		return
	}

	t, err := template.ParseFiles("html/list_clusters.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	dynData := ListClustersDynContent{Items: clusters}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

type ListProfilesDynContent struct {
	Items []ConfigurationProfile
}

func listProfiles(writer http.ResponseWriter, request *http.Request) {
	profiles, err := readListOfConfigurationProfiles(controllerURL, API_PREFIX)
	if err != nil {
		log.Println("Error reading list of configuration profiles", err)
		return
	}

	t, err := template.ParseFiles("html/list_profiles.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	dynData := ListProfilesDynContent{Items: profiles}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

type ListConfigurationsDynContent struct {
	Items []ClusterConfiguration
}

func listConfigurations(writer http.ResponseWriter, request *http.Request) {
	configurations, err := readListOfConfigurations(controllerURL, API_PREFIX)
	if err != nil {
		log.Println("Error reading list of cluster configurations", err)
		return
	}

	t, err := template.ParseFiles("html/list_configurations.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	dynData := ListConfigurationsDynContent{Items: configurations}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

type DescribeConfigurationDynContent struct {
	Configuration ConfigurationProfile
}

func describeConfiguration(writer http.ResponseWriter, request *http.Request) {
	configId, ok := request.URL.Query()["configuration"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	configuration, err := readConfigurationProfile(controllerURL, API_PREFIX, configId[0])
	fmt.Println(configuration)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	t, err := template.ParseFiles("html/describe_configuration.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Error parsing template")
		return
	}

	dynData := DescribeConfigurationDynContent{Configuration: *configuration}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

func startHttpServer() {
	http.HandleFunc("/", staticPage("html/index.html"))
	http.HandleFunc("/bootstrap.min.css", staticPage("html/bootstrap.min.css"))
	http.HandleFunc("/bootstrap.min.js", staticPage("html/bootstrap.min.js"))
	http.HandleFunc("/ccx.ccs", staticPage("html/ccs.ccs"))
	http.HandleFunc("/list-clusters", listClusters)
	http.HandleFunc("/list-profiles", listProfiles)
	http.HandleFunc("/list-configurations", listConfigurations)
	http.HandleFunc("/describe-configuration", describeConfiguration)
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

	controllerURL = viper.GetString("URL")

	log.Println("Starting the service")
	startHttpServer()
}
