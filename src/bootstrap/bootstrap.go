package bootstrap

import (
	"context"
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/domain/impl/auth"
	"github.com/africarealty/server/src/domain/impl/communications"
	"github.com/africarealty/server/src/domain/impl/filestore"
	httpAuth "github.com/africarealty/server/src/http/auth"
	"github.com/africarealty/server/src/http/system"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/auth/impl"
	kitHttp "github.com/africarealty/server/src/kit/http"
	"github.com/africarealty/server/src/kit/queue"
	"github.com/africarealty/server/src/kit/queue/jetstream"
	"github.com/africarealty/server/src/kit/queue/listener"
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
	queue                 queue.Queue
	queueListener         listener.QueueListener
	instanceId            string
}

// New creates a new instance of the service
func New() kitService.Service {
	s := &serviceImpl{
		instanceId: kit.UUID(4),
	}

	s.queue = jetstream.New(service.LF())
	s.queueListener = listener.NewQueueListener(s.queue, service.LF())
	s.authStorage = authStrg.NewAdapter()
	s.communicationsAdapter = commStrg.NewAdapter()
	templateStorage := communications.NewTemplateGenerator(s.communicationsAdapter.GetTemplateStorage())
	fileStore := filestore.NewStoreService(nil)
	s.smtpAdapter = smtp.NewAdapter(s.communicationsAdapter.GetEmailStorage())
	s.emailService = communications.NewEmailService(templateStorage, s.queue, fileStore, s.smtpAdapter)

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
	userRegUc := authUc.NewUserUseCases(userService, s.emailService, impl.NewPasswordService(service.LF()), s.cfg)

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

	// open Queue connection
	if err := s.queue.Open(ctx, s.instanceId, s.cfg.Nats); err != nil {
		return err
	}
	// declare topics
	if err := s.queue.Declare(ctx, queue.QueueTypeAtLeastOnce, domain.EmailRequestTopic); err != nil {
		return err
	}
	// add listeners
	s.queueListener.New(domain.EmailRequestTopic).AtLeastOnce(s.GetCode()).WithLoadBalancing(s.GetCode()).WithHandler(s.emailService.RequestHandler()).Add()

	// init services
	sessionService.Init(s.cfg.Auth.Session)
	if err := s.authStorage.Init(ctx, s.cfg); err != nil {
		return err
	}
	if err := s.communicationsAdapter.Init(ctx, s.cfg); err != nil {
		return err
	}
	if err := s.smtpAdapter.Init(ctx, s.cfg); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) Start(ctx context.Context) error {
	// start listening REST
	s.http.Listen()
	// queue listener
	s.queueListener.ListenAsync()
	return nil
}

func (s *serviceImpl) Close(ctx context.Context) {
	s.queueListener.Stop()
	_ = s.authStorage.Close(ctx)
	s.http.Close()
	_ = s.authStorage.Close(ctx)
	_ = s.smtpAdapter.Close(ctx)
	_ = s.queue.Close()
}
