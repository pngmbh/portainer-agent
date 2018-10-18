package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/portainer/agent"
	httpagenthandler "github.com/portainer/agent/http/handler/agent"
	"github.com/portainer/agent/http/handler/browse"
	"github.com/portainer/agent/http/handler/docker"
	"github.com/portainer/agent/http/handler/host"
	"github.com/portainer/agent/http/handler/websocket"
	"github.com/portainer/agent/http/proxy"
)

// Handler is the main handler of the application.
// Redirection to sub handlers is done in the ServeHTTP function.
type Handler struct {
	agentHandler       *httpagenthandler.Handler
	browseHandler      *browse.Handler
	browseHandlerV1    *browse.Handler
	dockerProxyHandler *docker.Handler
	webSocketHandler   *websocket.Handler
	hostHandler        *host.Handler
}

const (
	errInvalidQueryParameters = agent.Error("Invalid query parameters")
)

var dockerAPIVersionRegexp = regexp.MustCompile(`(/v[0-9]\.[0-9]*)?`)

// NewHandler returns a pointer to a Handler.
func NewHandler(systemService agent.SystemService, cs agent.ClusterService, notaryService agent.NotaryService, agentTags map[string]string) *Handler {
	agentProxy := proxy.NewAgentProxy(cs, agentTags)
	return &Handler{
		agentHandler:       httpagenthandler.NewHandler(cs, notaryService),
		browseHandler:      browse.NewHandler(agentProxy, notaryService),
		browseHandlerV1:    browse.NewHandlerV1(agentProxy, notaryService),
		dockerProxyHandler: docker.NewHandler(cs, agentTags, notaryService),
		webSocketHandler:   websocket.NewHandler(cs, agentTags, notaryService),
		hostHandler:        host.NewHandler(systemService, agentProxy, notaryService),
	}
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	request.URL.Path = dockerAPIVersionRegexp.ReplaceAllString(request.URL.Path, "")

	switch {
	case strings.HasPrefix(request.URL.Path, "/v1"):
		h.ServeHTTPV1(rw, request)
	case strings.HasPrefix(request.URL.Path, "/v2"):
		h.ServeHTTPV2(rw, request)
	case strings.HasPrefix(request.URL.Path, "/agents"):
		h.agentHandler.ServeHTTP(rw, request)
	case strings.HasPrefix(request.URL.Path, "/host"):
		h.hostHandler.ServeHTTP(rw, request)
	case strings.HasPrefix(request.URL.Path, "/browse"):
		h.browseHandler.ServeHTTP(rw, request)
	case strings.HasPrefix(request.URL.Path, "/websocket"):
		h.webSocketHandler.ServeHTTP(rw, request)
	case strings.HasPrefix(request.URL.Path, "/"):
		h.dockerProxyHandler.ServeHTTP(rw, request)
	}
}
