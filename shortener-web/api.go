package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HtmlContextData struct {
	Ok      bool
	Message *string
	Url     *string
}

type Api interface {
	Run()
}

func NewApi(config *Config, svc ShortenerService) Api {
	return NewHttpApi(config, svc)
}

type HttpApi struct {
	config *Config
	Router *mux.Router
	svc    ShortenerService
}

func NewHttpApi(
	config *Config,
	svc ShortenerService) *HttpApi {

	return &HttpApi{
		Router: mux.NewRouter().StrictSlash(true),
		config: config,
		svc:    svc,
	}
}

func (httpApi *HttpApi) Run() {
	httpApi.Router.HandleFunc("/health", httpApi.health).Methods("GET")
	httpApi.Router.HandleFunc("/", httpApi.homeHandler).Methods(http.MethodGet)
	httpApi.Router.HandleFunc("/", httpApi.shortHandler).Methods(http.MethodPost)
	httpApi.Router.HandleFunc("/{short}", httpApi.searchShortHandler).Methods(http.MethodGet)
	httpApi.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(httpApi.config.HTTP.Port, httpApi.Router))
}

func (httpApi *HttpApi) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UP"))
}

func (httpApi *HttpApi) render(w http.ResponseWriter, tpl string, data interface{}) {
	tmpl, err := template.ParseFiles(tpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (httpApi *HttpApi) homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (api *HttpApi) shortHandler(w http.ResponseWriter, r *http.Request) {
	data := &HtmlContextData{Ok: true, Message: nil, Url: nil}
	err := r.ParseForm()
	if err != nil {
		msg := err.Error()
		data.Ok = false
		data.Message = &msg
		api.render(w, "templates/index.html", data)
		return
	}

	url := r.PostFormValue("url")

	short, err := api.svc.Create(url, ReadUserIP(r))
	if err != nil {
		msg := err.Error()
		data.Ok = false
		data.Message = &msg
		api.render(w, "templates/index.html", data)
		return
	}

	msg := "Your short link was created successfuly"
	shortGen := fmt.Sprintf("%v/%v", api.config.App.Domain, short)
	data.Ok = true
	data.Message = &msg
	data.Url = &shortGen
	api.render(w, "templates/index.html", data)

}

func (api *HttpApi) searchShortHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["short"]
	url, err := api.svc.Search(id, ReadUserIP(r))
	if err != nil {
		msg := err.Error()
		data := &HtmlContextData{Ok: false, Message: &msg, Url: nil}
		api.render(w, "templates/index.html", data)
		return
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
