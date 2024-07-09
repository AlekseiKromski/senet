package gin_server

import (
	"alekseikromski.com/senet/core"
	v1 "alekseikromski.com/senet/modules/gin_server/v1"
	"alekseikromski.com/senet/modules/gin_server/ws"
	"context"
	"embed"
	"net"
	"net/http"
)

type ServerConfig struct {
	Address      string
	Secret       []byte
	CookieDomain string
}

func NewServerConfig(secret string, address string, cookieDomain string) *ServerConfig {
	return &ServerConfig{
		Address:      address,
		Secret:       []byte(secret),
		CookieDomain: cookieDomain,
	}
}

type Server struct {
	config    *ServerConfig
	server    *http.Server
	ws        *ws.WebSocket
	api       Api
	busEvent  chan core.BusEvent
	resources embed.FS
}

func NewServer(conf *ServerConfig, resources embed.FS) *Server {
	return &Server{
		config:    conf,
		resources: resources,
	}
}

func (s *Server) Start(notifyChannel chan struct{}, busEventChannel chan core.BusEvent, requirements map[string]core.Module) {
	s.Log("init http server")

	s.busEvent = busEventChannel
	go s.listenEventBus()

	storage, err := s.getStorageFromRequirement(requirements)
	if err != nil {
		s.Log("cannot get storage requirement", err.Error())
		return
	}
	serverKeyStorage, err := s.getServerKeyStorageRequirement(requirements)
	if err != nil {
		s.Log("cannot get server key storage requirement", err.Error())
		return
	}

	s.api = v1.NewV1Api(storage, s.config.Secret, s.config.CookieDomain, serverKeyStorage, s.Log)

	if err := s.api.RegisterRoutes(s.resources); err != nil {
		s.Log("cannot register routes")
		return
	}

	s.ws, err = ws.NewWebSocket(s.api.GetEngine(), s.config.Secret, s.api.GetGuard(), storage, serverKeyStorage, s.Log)
	if err != nil {
		s.Log("cannot start websocket server", err.Error())
		return
	}

	// Create tcp listener and server
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		s.Log("cannot create tcp listener and server")
		return
	}
	s.server = &http.Server{
		Handler: s.api.GetEngine(),
	}

	// Notify core, that we started listener
	notifyChannel <- struct{}{}

	// Start server
	if err := s.server.Serve(listener); err != nil {
		s.Log("cannot serve", err.Error())
		return
	}
}

func (s *Server) Stop() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		s.Log("cannot stop server", err.Error())

		return
	}
}

func (s *Server) Signature() string {
	return "gin_server"
}

func (s *Server) listenEventBus() {
	for {
		event := <-s.busEvent // listen event bus
		if event.Receiver != s.Signature() {
			continue
		}

		s.Log("New bus event received, send to all clients")

		if err := s.ws.SendDatapointsToAllClients(event.Payload); err != nil {
			s.Log("cannot send datapoints to ws clients", err.Error())
		}
	}
}
