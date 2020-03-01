package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

var (
	infoMessage = "Thanks for using GoScrum"
)

type (
	postActionHandler func(*model.PostActionIntegrationRequest) (*model.Post, error)
)

func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", p.handleInfo).Methods(http.MethodGet)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	//apiV1.Use(checkAuthenticity)
	apiV1.HandleFunc("/configuration", p.handlePluginConfiguration).Methods(http.MethodGet)

	apiV1.HandleFunc("/user/action", p.handlePostActionIntegrationRequest(p.handleVote)).Methods(http.MethodPost)
	return r
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.API.LogDebug("New request:", "Host", r.Host, "RequestURI", r.RequestURI, "Method", r.Method)
	p.router.ServeHTTP(w, r)
}

func checkAuthenticity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Mattermost-User-ID") == "" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (p *Plugin) handleInfo(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, infoMessage)
}

func (p *Plugin) handlePluginConfiguration(w http.ResponseWriter, r *http.Request) {
	configuration := p.getConfiguration()
	b, err := json.Marshal(configuration)
	if err != nil {
		p.API.LogWarn("failed to decode configuration object.", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(b); err != nil {
		p.API.LogWarn("failed to write response.", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) handlePostActionIntegrationRequest(handler postActionHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := model.PostActionIntegrationRequestFromJson(r.Body)
		if request == nil {
			p.API.LogWarn("failed to decode PostActionIntegrationRequest")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		update, err := handler(request)
		if err != nil {
			p.API.LogWarn("failed to handle PostActionIntegrationRequest", "error", err.Error())
		}

		message := "message"

		p.SendEphemeralPost(request.ChannelId, request.UserId, message)

		response := &model.PostActionIntegrationResponse{}
		if update != nil {
			response.Update = update
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(response.ToJson()); err != nil {
			p.API.LogWarn("failed to write PostActionIntegrationResponse", "error", err.Error())
		}
	}
}

func (p *Plugin) handleVote(request *model.PostActionIntegrationRequest) (*model.Post, error) {
	spew.Dump(request)
	return nil, nil
}

func (p *Plugin) SendEphemeralPost(channelID, userID, message string) {
	ephemeralPost := &model.Post{
		ChannelId: channelID,
		UserId:    p.botUserID,
		Message:   message,
	}
	_ = p.API.SendEphemeralPost(userID, ephemeralPost)
}
