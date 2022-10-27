package bootstrap

import (
	"context"
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/domain/impl/auth"
	"github.com/africarealty/server/src/domain/impl/communications"
	"github.com/africarealty/server/src/domain/impl/filestore"
	httpAuth "github.com/africarealty/server/src/http/auth"
	"github.com/africarealty/server/src/http/system"
	"github.com/africarealty/server/src/kit/auth/impl"
	kitHttp "github.com/africarealty/server/src/kit/http"
	kitService "github.com/africarealty/server/src/kit/service"
	"github.com/africarealty/server/src/repository/adapters/smtp"
	authStrg "github.com/africarealty/server/src/repository/storage/auth"
	commStrg "github.com/africarealty/server/src/repository/storage/communications"
	"github.com/africarealty/server/src/service"
	authUc "github.com/africarealty/server/src/usecase/impl/auth"
)

type serviceImpl struct {
	cfg                   *service.Config
	http                  *kitHttp.Server
	authStorage           authStrg.Adapter
	communicationsAdapter commStrg.Adapter
	smtpAdapter           smtp.Adapter
	emailService          domain.EmailService
}

// New creates a new instance of the service
func New() kitService.Service {
	s := &serviceImpl{}

	s.authStorage = authStrg.NewAdapter()
	s.communicationsAdapter = commStrg.NewAdapter()
	templateStorage := communications.NewTemplateGenerator(s.communicationsAdapter.GetTemplateStorage())
	fileStore := filestore.NewStoreService(nil)
	s.smtpAdapter = smtp.NewAdapter()
	s.emailService = communications.NewEmailService(templateStorage, nil, fileStore, s.smtpAdapter)

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
	authorizeSession := auth.NewAuthorizeService(s.authStorage)
	sessionService := impl.NewSessionsService(service.LF(), s.authStorage, s.authStorage, authorizeSession)
	userService := auth.NewUserService(s.authStorage)
	userRegUc := authUc.NewUserRegistrationImpl(userService, s.emailService, impl.NewPasswordService(service.LF()))

	// create and set middlewares
	mdw := kitHttp.NewMiddleware(service.LF(), sessionService, authorizeSession, resourcePolicyManager)
	s.http.RootRouter.Use(mdw.SetContextMiddleware)

	// set up routing
	routeBuilder := kitHttp.NewRouteBuilder(s.http, resourcePolicyManager, mdw)

	// setup routes & controllers
	routers := []kitHttp.RouteSetter{
		httpAuth.NewRouter(httpAuth.NewController(sessionService, userService, userRegUc), routeBuilder),
		system.NewRouter(system.NewController(), routeBuilder),
	}
	for _, r := range routers {
		if err := r.Set(); err != nil {
			return err
		}
	}

	// init services
	sessionService.Init(&s.cfg.Auth.Config)
	if err := s.authStorage.Init(ctx, s.cfg.Auth); err != nil {
		return err
	}
	if err := s.smtpAdapter.Init(ctx, s.cfg.Communications.Email); err != nil {
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
	_ = s.authStorage.Close(ctx)
	s.http.Close()
	_ = s.authStorage.Close(ctx)
	_ = s.smtpAdapter.Close(ctx)
}
