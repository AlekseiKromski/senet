package processor

import (
	"embed"
	"fmt"
	"github.com/AlekseiKromski/at-socket-server/core"
	"senet/config"
	"senet/processor/api"
	"senet/processor/lb"
	"senet/processor/storage"
	"senet/processor/storage/dbstorage"
	hs "senet/processor/ws"
)

type Processor struct {
	config   *config.Config
	lb       *lb.LoadBalancer
	storage  storage.Storage
	handlers *core.Handlers
}

func NewProcessor(config *config.Config) (*Processor, error) {
	store, err := dbstorage.NewDbStorage(config.DbConfig)
	if err != nil {
		return nil, fmt.Errorf("problem with database: %v", err)
	}

	p := &Processor{
		config:  config,
		lb:      lb.NewLoadBalancer(store),
		storage: store,
	}
	p.handlers = p.registerWebsocketHandlers()

	return p, nil
}

func (p *Processor) Start(frontend embed.FS) error {
	app, err := core.Start(p.handlers, p.config.AppConfig)
	if err != nil {
		return fmt.Errorf("cannot start core application: %v", err)
	}

	api := api.NewApi(p.lb, &api.ProcessorFs{frontend}, app.Engine, p.storage)

	api.RegisterStaticFiles()
	api.Register()

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
