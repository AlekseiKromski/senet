package processor

import (
	"embed"
	"fmt"
	"github.com/AlekseiKromski/at-socket-server/core"
	"net/http"
	"senet/processor/api"
	hs "senet/processor/ws"
)

type Processor struct {
	config   *core.Config
	handlers *core.Handlers
}

func NewProcessor(config *core.Config) *Processor {
	p := &Processor{
		config: config,
	}
	p.handlers = p.registerWebsocketHandlers()

	return p
}

func (p *Processor) Start(frontend embed.FS) error {
	app, err := core.Start(p.handlers, p.config)
	if err != nil {
		return fmt.Errorf("cannot start core application: %v", err)
	}

	//reassign fs system for complete /webclient path
	pfs := &processorFs{frontend}

	p.registerStaticFiles(app, pfs)
	p.registerApiHandlers(app, pfs)

	//Start server in additional gorutine
	go func() {
		app.Engine.Run(app.Config.GetServerString())
	}()

	//Block main thread for showing log information
	for {
		hook := <-app.Hooks
		switch hook.HookType {
		case core.CLIENT_ADDED:
			fmt.Printf("Client added: %s\n", hook.Data)
		case core.CLIENT_CLOSED_CONNECTION:
			fmt.Printf("Client closed connection: %s\n", hook.Data)
		case core.ERROR:
			fmt.Printf("Error: %s\n", hook.Data)
		}
	}
}

func (p *Processor) registerWebsocketHandlers() *core.Handlers {
	handlers := make(core.Handlers)
	handlers[hs.SEND_MESSAGE] = &hs.SenderHandler{}

	return &handlers
}

func (p *Processor) registerApiHandlers(app *core.App, pfs *processorFs) {
	app.Engine.GET("/healthz", api.Healthz)
	app.Engine.GET("/", api.Webclient(pfs.content))
}

func (p *Processor) registerStaticFiles(app *core.App, pfs *processorFs) {
	app.Engine.StaticFS("/static", http.FS(pfs))
}
