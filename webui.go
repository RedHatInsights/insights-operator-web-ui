/*
Copyright © 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// API_PREFIX represents part of URL that is appended before the actual endpoint address
const API_PREFIX = "/api/v1/"

// Cluster represents cluster record in the controller service.
//     ID: unique key
//     Name: cluster GUID in the following format:
//         c8590f31-e97e-4b85-b506-c45ce1911a12
type Cluster struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// ConfigurationProfile represents configuration profile record in the controller service.
//     ID: unique key
//     Configuration: a JSON structure stored in a string
//     ChangeAt: username of admin that created or updated the configuration
//     ChangeBy: timestamp of the last configuration change
//     Description: a string with any comment(s) about the configuration
type ConfigurationProfile struct {
	Id            int    `json:"id"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Description   string `json:"description"`
}

// ClusterConfiguration represents cluster configuration record in the controller service.
//     ID: unique key
//     Cluster: cluster ID (not name)
//     Configuration: a JSON structure stored in a string
//     ChangeAt: timestamp of the last configuration change
//     ChangeBy: username of admin that created or updated the configuration
//     Active: flag indicating whether the configuration is active or not
//     Reason: a string with any comment(s) about the cluster configuration
type ClusterConfiguration struct {
	Id            int    `json:"id"`
	Cluster       string `json:"cluster"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Active        string `json:"active"`
	Reason        string `json:"reason"`
}

// Trigger represents trigger record in the controller service
//     ID: unique key
//     Type: ID of trigger type
//     Cluster: cluster ID (not name)
//     Reason: a string with any comment(s) about the trigger
//     Link: link to any document with customer ACK with the trigger
//     TriggeredAt: timestamp of the last configuration change
//     TriggeredBy: username of admin that created or updated the trigger
//     AckedAt: timestamp where the insights operator acked the trigger
//     Parameters: parameters that needs to be pass to trigger code
//     Active: flag indicating whether the trigger is still active or not
type Trigger struct {
	Id          int    `json:"id"`
	Type        string `json:"type"`
	Cluster     string `json:"cluster"`
	Reason      string `json:"reason"`
	Link        string `json:"link"`
	TriggeredAt string `json:"triggered_at"`
	TriggeredBy string `json:"triggered_by"`
	AckedAt     string `json:"acked_at"`
	Parameters  string `json:"parameters"`
	Active      int    `json:"active"`
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

func performWriteRequest(url string, method string, payload io.Reader) error {
	var client http.Client

	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return fmt.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Expected HTTP status 200 OK, 201 Created or 202 Accepted, got %d", response.StatusCode)
	}
	return nil
}

func readListOfClusters(controllerUrl string, apiPrefix string) ([]Cluster, error) {
	clusters := []Cluster{}

	url := controllerUrl + apiPrefix + "client/cluster"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}

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

func readListOfTriggers(controllerUrl string, apiPrefix string, clusterName string) ([]Trigger, error) {
	var triggers []Trigger
	url := controllerUrl + apiPrefix + "client/cluster/" + clusterName + "/trigger"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &triggers)
	if err != nil {
		return nil, err
	}
	return triggers, nil
}

func readListOfAllTriggers(controllerUrl string, apiPrefix string) ([]Trigger, error) {
	var triggers []Trigger
	url := controllerUrl + apiPrefix + "client/trigger"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &triggers)
	if err != nil {
		return nil, err
	}
	return triggers, nil
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

// ListClustersDynContent represents dynamic part of HTML page with list of clusters
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

// ListProfilesDynContent represents dynamic part of HTML page with list of configuration profiles
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

// ListConfigurationsDynContent represents dynamic part of HTML page with list of configurations
type ListConfigurationsDynContent struct {
	Items []ClusterConfiguration
}

// ListTriggersDynContent represents dynamic part of HTML page with list of triggers
type ListTriggersDynContent struct {
	Items []Trigger
}

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

func listConfigurations(writer http.ResponseWriter, request *http.Request) {
	configurations, err := readListOfConfigurations(controllerURL, API_PREFIX)
	// NoCache headers
	for k, v := range noCacheHeaders {
		writer.Header().Set(k, v)
	}

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

func listTriggers(writer http.ResponseWriter, request *http.Request) {
	clusterName, ok := request.URL.Query()["clusterName"]
	var triggers []Trigger
	var err error

	if !ok {
		triggers, err = readListOfAllTriggers(controllerURL, API_PREFIX)
	} else {
		triggers, err = readListOfTriggers(controllerURL, API_PREFIX, clusterName[0])
	}

	// NoCache headers
	for k, v := range noCacheHeaders {
		writer.Header().Set(k, v)
	}

	if err != nil {
		log.Println("Error reading list of triggers", err)
		return
	}

	log.Println(triggers)
	t, err := template.ParseFiles("html/list_triggers.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	dynData := ListTriggersDynContent{Items: triggers}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

// DescribeConfigurationDynContent represents dynamic part of HTML page with configuration description
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

func storeProfile(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	form := request.Form

	username := form.Get("username")
	description := form.Get("description")
	configuration := form.Get("configuration")

	log.Println("username", username)
	log.Println("description", description)
	log.Println("configuration", configuration)

	query := "username=" + url.QueryEscape(username) + "&description=" + url.QueryEscape(description)
	url := controllerURL + API_PREFIX + "client/profile?" + query

	err := performWriteRequest(url, "POST", strings.NewReader(configuration))
	if err != nil {
		log.Println("Error communicating with the service", err)
		http.Redirect(writer, request, "/profile-not-created", 301)
	} else {
		log.Println("Configuration profile has been created")
		http.Redirect(writer, request, "/profile-created", 301)
	}
}

func storeConfiguration(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	form := request.Form

	username := form.Get("username")
	cluster := form.Get("cluster")
	reason := form.Get("reason")
	description := form.Get("description")
	configuration := form.Get("configuration")

	log.Println("username", username)
	log.Println("cluster", cluster)
	log.Println("reason", reason)
	log.Println("description", description)
	log.Println("configuration", configuration)

	query := "username=" + url.QueryEscape(username) + "&reason=" + url.QueryEscape(reason) + "&description=" + url.QueryEscape(description)
	url := controllerURL + API_PREFIX + "client/cluster/" + url.PathEscape(cluster) + "/configuration?" + query

	err := performWriteRequest(url, "POST", strings.NewReader(configuration))
	if err != nil {
		log.Println("Error communicating with the service", err)
		http.Redirect(writer, request, "/configuration-not-created", 301)
	} else {
		log.Println("Configuration has been created")
		http.Redirect(writer, request, "/configuration-created", 301)
	}
}

func enableConfiguration(writer http.ResponseWriter, request *http.Request) {
	configurationId, ok := request.URL.Query()["id"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}
	url := controllerURL + API_PREFIX + "client/configuration/" + configurationId[0] + "/enable"
	err := performWriteRequest(url, "PUT", nil)
	if err != nil {
		fmt.Println("Error communicating with the service", err)
		return
	}

	// everything is ok, configuration has been enabled
	fmt.Println("Configuration " + configurationId[0] + " has been enabled")
	http.Redirect(writer, request, "/list-configurations", 307)
}

func disableConfiguration(writer http.ResponseWriter, request *http.Request) {
	configurationId, ok := request.URL.Query()["id"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}
	url := controllerURL + API_PREFIX + "client/configuration/" + configurationId[0] + "/disable"
	err := performWriteRequest(url, "PUT", nil)
	if err != nil {
		fmt.Println("Error communicating with the service", err)
		return
	}

	// everything is ok, configuration has been disabled
	fmt.Println("Configuration " + configurationId[0] + " has been disabled")
	http.Redirect(writer, request, "/list-configurations", 307)
}

func activateTrigger(writer http.ResponseWriter, request *http.Request) {
	triggerId, ok := request.URL.Query()["id"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}
	url := controllerURL + API_PREFIX + "client/trigger/" + triggerId[0] + "/activate"

	err := performWriteRequest(url, "PUT", nil)
	if err != nil {
		fmt.Println("Error communicating with the service", err)
		return
	}

	// everything is ok, trigger has been activated
	fmt.Println("Trigger " + triggerId[0] + " has been activated")
	http.Redirect(writer, request, "/list-triggers", 307)
}

func deactivateTrigger(writer http.ResponseWriter, request *http.Request) {
	triggerId, ok := request.URL.Query()["id"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}
	url := controllerURL + API_PREFIX + "client/trigger/" + triggerId[0] + "/deactivate"

	err := performWriteRequest(url, "PUT", nil)
	if err != nil {
		fmt.Println("Error communicating with the service", err)
		return
	}

	// everything is ok, trigger has been deactivated
	fmt.Println("Trigger " + triggerId[0] + " has been deactivated")
	http.Redirect(writer, request, "/list-triggers", 307)
}

func triggerMustGatherConfiguration(writer http.ResponseWriter, request *http.Request) {
	clusterId, ok := request.URL.Query()["clusterId"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}
	id, err := strconv.Atoi(clusterId[0])
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	clusterName, ok := request.URL.Query()["clusterName"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	t, err := template.ParseFiles("html/trigger_must_gather.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Error parsing template")
		return
	}
	dynData := Cluster{Id: id, Name: clusterName[0]}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

// POST must-gather to REST API
func triggerMustGather(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	form := request.Form

	clusterId := form.Get("clusterid")
	clusterName := form.Get("clustername")
	username := form.Get("username")
	reason := form.Get("reason")
	link := form.Get("link")

	log.Println("clusterId", clusterId)
	log.Println("clusterName", clusterName)
	log.Println("username", username)
	log.Println("reason", reason)
	log.Println("link", link)

	query := "username=" + url.QueryEscape(username) + "&reason=" + url.QueryEscape(reason) + "&link=" + url.QueryEscape(link)
	log.Println(query)
	url := controllerURL + API_PREFIX + "client/cluster/" + url.PathEscape(clusterName) + "/trigger/must-gather?" + query
	log.Println(url)

	err := performWriteRequest(url, "POST", nil)
	if err != nil {
		log.Println("Error communicating with the service", err)
		http.Redirect(writer, request, "/trigger-not-created", 301)
	} else {
		log.Println("Trigger has been created")
		http.Redirect(writer, request, "/trigger-created", 301)
	}
}

func startHttpServer(address string) {
	http.HandleFunc("/", staticPage("html/index.html"))
	http.HandleFunc("/bootstrap.min.css", staticPage("html/bootstrap.min.css"))
	http.HandleFunc("/bootstrap.min.js", staticPage("html/bootstrap.min.js"))
	http.HandleFunc("/ccx.css", staticPage("html/ccx.css"))
	http.HandleFunc("/configuration-created", staticPage("html/configuration_created.html"))
	http.HandleFunc("/configuration-not-created", staticPage("html/configuration_not_created.html"))
	http.HandleFunc("/profile-created", staticPage("html/profile_created.html"))
	http.HandleFunc("/profile-not-created", staticPage("html/profile_not_created.html"))
	http.HandleFunc("/list-clusters", listClusters)
	http.HandleFunc("/list-profiles", listProfiles)
	http.HandleFunc("/list-configurations", listConfigurations)
	http.HandleFunc("/list-all-triggers", listTriggers)
	http.HandleFunc("/list-triggers", listTriggers)
	http.HandleFunc("/describe-configuration", describeConfiguration)
	http.HandleFunc("/new-profile", staticPage("html/new_profile.html"))
	http.HandleFunc("/new-configuration", staticPage("html/new_configuration.html"))
	http.HandleFunc("/store-profile", storeProfile)
	http.HandleFunc("/store-configuration", storeConfiguration)
	http.HandleFunc("/enable-configuration", enableConfiguration)
	http.HandleFunc("/disable-configuration", disableConfiguration)
	http.HandleFunc("/activate-trigger", activateTrigger)
	http.HandleFunc("/deactivate-trigger", deactivateTrigger)
	http.HandleFunc("/trigger-must-gather-configuration", triggerMustGatherConfiguration)
	http.HandleFunc("/trigger-must-gather", triggerMustGather)
	http.HandleFunc("/trigger-created", staticPage("html/trigger_created.html"))
	http.HandleFunc("/trigger-not-created", staticPage("html/trigger_not_created.html"))
	http.ListenAndServe(address, nil)
}

func main() {
	log.Println("Reading configuration")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	controllerURL = viper.GetString("controller_url")
	address := viper.GetString("address")

	log.Println("Starting the service at address: " + address)
	startHttpServer(address)
}
