package web

import (
	"encoding/json"
	"github.com/oleksiy-os/porto-events/configs"
	"github.com/oleksiy-os/porto-events/internal/model"
	telegramApi "github.com/oleksiy-os/porto-events/internal/model/client/telegram"
	"github.com/oleksiy-os/porto-events/internal/model/event"
	"github.com/oleksiy-os/porto-events/internal/store"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"os"
)

const (
	pathWeb  = "internal/web"
	htmlPath = "internal/web/templates/"
)

var templates = template.Must(template.ParseFiles(htmlPath + "home.html"))

type (
	Server struct {
		store  store.StoreInterface
		config *configs.Config
	}
)

func New(config *configs.Config, store *store.StoreInterface) *Server {
	s := &Server{
		store:  *store,
		config: config,
	}

	s.configureRouter()

	return s
}

func (s *Server) ListenAndServe() {
	log.Fatal(http.ListenAndServe(s.config.Server.BindAddr, nil))
}

func (s *Server) configureRouter() {
	http.HandleFunc("/", s.homeHandler)
	http.HandleFunc("/move/", s.changeCategoryHandler)
	http.HandleFunc("/save/", s.saveHandler)
	http.HandleFunc("/delete/", s.deleteHandler)
	http.HandleFunc("/get/", s.getHandler)
	http.HandleFunc("/publish/", s.publishHandler)

	http.HandleFunc("/assets/", s.staticHandler)
	http.HandleFunc("/templates/", s.staticHandler)
}

func (s *Server) homeHandler(w http.ResponseWriter, _ *http.Request) {
	if !s.config.ProductionMode { // for live changes in html during develop
		templates = template.Must(template.ParseFiles(htmlPath + "home.html"))
	}

	if err := templates.ExecuteTemplate(w, "home.html", *s.store.Event().Get()); err != nil {
		log.Error("exec template|", err)
	}
}

func (s *Server) changeCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer closeBody(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("read body|", err)
		http.Error(w, "wrong data", http.StatusBadRequest)
		return
	}

	var data store.ChangeCategoryData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Error("unmarshal|", err)
		http.Error(w, "wrong data", http.StatusBadRequest)
		return
	}

	if ok := s.store.Event().ChangeCategory(data); !ok {
		http.Error(w, "failed save data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer closeBody(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("read body|", err)
		http.Error(w, "wrong data", http.StatusBadRequest)
		return
	}

	var ev model.Event
	if err = json.Unmarshal(body, &ev); err != nil {
		log.Error("unmarshal|", err)
		http.Error(w, "wrong data", http.StatusBadRequest)
		return
	}

	if ok := s.store.Event().Save(&ev); !ok {
		http.Error(w, "failed save data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer closeBody(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		log.Error("read body|", err)
		http.Error(w, "wrong data", http.StatusBadRequest)
		return
	}

	if ok := s.store.Event().Delete(string(body)); !ok {
		http.Error(w, "failed delete data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sources, err := model.GetSources(s.config.SourcesListPath)
	if err != nil {
		log.Error("parse toml| ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	events := event.Collect(sources)

	for _, e := range *events {
		s.store.Event().Add(&e)
	}

	evs, err := json.Marshal(s.store.Event().Get())
	if err != nil {
		log.Error("get events, json marshal|", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(evs)
	if err != nil {
		log.Error("write data to response|", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) staticHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(pathWeb + r.URL.Path); err != nil {
		log.Error("file path|", err)
	}

	http.ServeFile(w, r, pathWeb+r.URL.Path)
}

func (s *Server) publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bot := telegramApi.New(s.config.Telegram)
	for _, ev := range *s.store.Event().GetCategoryPublish() {
		if err := bot.Publish(&ev); err != nil {
			continue
		}
		s.store.Event().ChangeCategory(store.ChangeCategoryData{
			Id:       ev.ID,
			Category: store.CategoryPublished,
		})
	}

	w.WriteHeader(http.StatusOK)
}

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Error("close body|", err)
	}
}
