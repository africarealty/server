package bootstrap

import (
	"context"
	"github.com/africarealty/server/src/domain/impl/auth"
	httpAuth "github.com/africarealty/server/src/http/auth"
	"github.com/africarealty/server/src/http/system"
	"github.com/africarealty/server/src/kit/auth/impl"
	kitHttp "github.com/africarealty/server/src/kit/http"
	kitService "github.com/africarealty/server/src/kit/service"
	"github.com/africarealty/server/src/repository/storage"
	"github.com/africarealty/server/src/service"
)

type serviceImpl struct {
	cfg            *service.Config
	http           *kitHttp.Server
	storageAdapter storage.Adapter
}

// New creates a new instance of the service
func New() kitService.Service {
	s := &serviceImpl{}

	s.storageAdapter = storage.NewAdapter()

	return s
}

func (s *serviceImpl) GetCode() string {
	return "africarealty"
}

// Init does all initializations
func (s *serviceImpl) Init(ctx context.Context) error {

	// load config
	var err error
	s.cfg, err = service.LoadConfig()
	if err != nil {
		return err
	}

	// set log config
	service.Logger.Init(s.cfg.Log)

	// create HTTP server
	s.http = kitHttp.NewHttpServer(s.cfg.Http, service.LF())

	// create resource policy manager
	resourcePolicyManager := impl.NewResourcePolicyManager(service.LF())
	authorizeSession := auth.NewAuthorizeService(s.storageAdapter)
	sessionService := impl.NewSessionsService(service.LF(), s.storageAdapter, s.storageAdapter, authorizeSession)
	userService := auth.NewUserService(s.storageAdapter)

	// create and set middlewares
	mdw := kitHttp.NewMiddleware(service.LF(), sessionService, authorizeSession, resourcePolicyManager)
	s.http.RootRouter.Use(mdw.SetContextMiddleware)

	// set up routing
	routeBuilder := kitHttp.NewRouteBuilder(s.http, resourcePolicyManager, mdw)

	// setup routes & controllers
	routers := []kitHttp.RouteSetter{
		httpAuth.NewRouter(httpAuth.NewController(sessionService, userService), routeBuilder),
		system.NewRouter(system.NewController(), routeBuilder),
	}
	for _, r := range routers {
		if err := r.Set(); err != nil {
			return err
		}
	}

	// init services
	sessionService.Init(s.cfg.Auth)

	if err := s.storageAdapter.Init(ctx, s.cfg); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) Start(ctx context.Context) error {

	// start listening REST
	s.http.Listen()

	return nil
}

func (s *serviceImpl) Close(ctx context.Context) {
	_ = s.storageAdapter.Close(ctx)
	s.http.Close()
}
