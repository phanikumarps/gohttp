package apiserver

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var defaultStopTimeout = time.Second * 30

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) (*APIServer, error) {
	if addr == "" {
		return nil, errors.New("addr cannot be blank")
	}

	return &APIServer{
		addr: addr,
	}, nil
}

// Start starts a server with a stop channel
func (s *APIServer) Start(stop <-chan struct{}) error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.router(),
	}

	go func() {
		logrus.WithField("addr", srv.Addr).Info("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), defaultStopTimeout)
	defer cancel()

	logrus.WithField("timeout", defaultStopTimeout).Info("stopping server")
	return srv.Shutdown(ctx)
}

func (s *APIServer) router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", s.defaultRoute)
	router.HandleFunc("/hello", s.hello)
	router.HandleFunc("/account", s.account)
	router.HandleFunc("/metadata", s.metadata)
	return router
}

func (s *APIServer) defaultRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}
func (s *APIServer) hello(w http.ResponseWriter, r *http.Request) {

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL()),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	request, err := http.NewRequest("GET", getDestUrl().String(), nil)
	if err != nil {
		log.Println(err)
	}
	h := map[string][]string{
		"Content-Type":                     {"application/json"},
		"SAP-Connectivity-SCC-Location_ID": {""},
	}
	request.Header = h

	//calling the URL
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	//printing the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
func (s *APIServer) account(w http.ResponseWriter, r *http.Request) {

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL()),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	opurl := getDestUrl().String() + "/j"
	request, err := http.NewRequest("GET", opurl, nil)
	if err != nil {
		log.Println(err)
	}
	h := map[string][]string{
		"Content-Type":                     {"application/json"},
		"SAP-Connectivity-SCC-Location_ID": {""},
	}
	request.Header = h

	//calling the URL
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	//printing the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
func (s *APIServer) metadata(w http.ResponseWriter, r *http.Request) {

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL()),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	opurl := getDestUrl().String() + "/sap/opu/odata/sap/ERP_ISU_UMC/$metadata"
	request, err := http.NewRequest("GET", opurl, nil)
	if err != nil {
		log.Println(err)
	}
	h := map[string][]string{
		"Content-Type":                     {"application/json"},
		"SAP-Connectivity-SCC-Location_ID": {""},
	}
	request.Header = h

	//calling the URL
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	//printing the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
func getDestUrl() *url.URL {

	opUrl, err := url.Parse("http://http-host:8001/")
	if err != nil {
		log.Println(err)
	}
	return opUrl
}

func proxyURL() *url.URL {
	Url, err := url.Parse("http://connectivity-proxy.kyma-system.svc.cluster.local:20003")
	if err != nil {
		log.Println(err)
	}
	return Url
}
