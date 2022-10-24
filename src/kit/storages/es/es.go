package es

import (
	"context"
	"fmt"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/goroutine"
	"github.com/africarealty/server/src/kit/log"
	"github.com/olivere/elastic/v7"
)

// Config - model of ES configuration
type Config struct {
	Host     string // Host - ES host
	Port     string // Port - ES port
	Trace    bool   // Trace enables tracing mode
	Sniff    bool   // Sniff - read https://github.com/olivere/elastic/issues/387
	Shards   int    // Shards - how many shards to be created for index
	Replicas int    // Replicas - how many replicas to eb created for index
	Username string // Username - ES basic auth (if not set, no auth applied)
	Password string // Password - ES basic auth
	Ssl      bool   // Ssl - Ssl used
	Refresh  bool   // Refresh - enforces refresh after each change. It helpful for tests but MUST NOT BE USED ON PROD
}

// Search allows indexing and searching with ES
type Search interface {
	// Index indexes a document
	Index(index string, id string, data interface{}) error
	// IndexAsync indexes a document async
	IndexAsync(index string, id string, data interface{})
	// IndexBulk allows indexing bulk of documents in one hit
	IndexBulk(index string, docs map[string]interface{}) error
	// IndexBulkAsync allows indexing bulk of documents in one hit
	IndexBulkAsync(index string, docs map[string]interface{})
	// GetClient provides an access to ES client
	GetClient() *elastic.Client
	// Close closes client
	Close()
	//Ping pings server
	Ping() bool
	// Exists checks if a document exists in the index
	Exists(index, id string) (bool, error)
	// Delete deletes a document
	Delete(index string, id string) error
	// DeleteBulk deletes bulk of documents
	DeleteBulk(index string, ids []string) error
	// NewBuilder creates a new builder object
	NewBuilder() IndexBuilder
	// Refresh refreshes data in index (don't use in production code)
	Refresh(index string) error
}

type esImpl struct {
	client *elastic.Client
	logger log.CLoggerFunc
	cfg    *Config
	url    string
}

func (s *esImpl) l() log.CLogger {
	return s.logger().Cmp("es")
}

func NewEs(cfg *Config, logger log.CLoggerFunc) (Search, error) {

	s := &esImpl{
		logger: logger,
		cfg:    cfg,
	}
	l := s.l().Mth("new").F(log.FF{"host": cfg.Host, "sniff": cfg.Sniff})

	if cfg.Ssl {
		s.url = fmt.Sprintf("https://%s:%s", cfg.Host, cfg.Port)
	} else {
		s.url = fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)
	}

	opts := []elastic.ClientOptionFunc{elastic.SetURL(s.url), elastic.SetSniff(cfg.Sniff)}
	if cfg.Trace {
		opts = append(opts, elastic.SetTraceLog(s.l().Mth("es-trace")))
	}

	// basic auth
	if cfg.Username != "" {
		if cfg.Password == "" {
			return nil, ErrEsBasicAuthInvalid()
		}
		opts = append(opts, elastic.SetBasicAuth(cfg.Username, cfg.Password))
	}
	l.F(log.FF{"auth": cfg.Username != ""})

	cl, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, ErrEsNewClient(err)
	}
	s.client = cl
	l.Inf("ok")
	return s, nil
}

func (s *esImpl) NewBuilder() IndexBuilder {
	return &esIndexBuilder{
		client: s.client,
		logger: s.logger,
		cfg:    s.cfg,
	}
}

func (s *esImpl) Ping() bool {
	s.l().Mth("ping").Dbg()
	_, code, err := s.client.Ping(s.url).Do(context.Background())
	return err == nil && code == 200
}

func (s *esImpl) Index(index string, id string, doc interface{}) error {
	s.l().Mth("indexation").F(log.FF{"index": index, "id": id}).Dbg().Trc(kit.Json(doc))
	svc := s.client.Index().
		Index(index).
		Id(id).
		BodyJson(doc)
	_, err := svc.Do(context.Background())
	if err != nil {
		return ErrEsIdx(err, index, id)
	}
	// refresh
	if s.cfg.Refresh {
		return s.Refresh(index)
	}
	return nil
}

func (s *esImpl) IndexAsync(index string, id string, doc interface{}) {
	goroutine.New().
		WithLogger(s.l().Mth("index-async")).
		Go(context.Background(), func() {
			l := s.l().Mth("index-async").F(log.FF{"index": index, "id": id}).Dbg().Trc(kit.Json(doc))
			err := s.Index(index, id, doc)
			if err != nil {
				l.E(err).Err()
			}
		})
}

func (s *esImpl) IndexBulk(index string, docs map[string]interface{}) error {
	s.l().Mth("bulk-indexation").F(log.FF{"index": index, "docs": len(docs)}).Dbg()
	bulk := s.client.Bulk().Index(index)
	for id, doc := range docs {
		bulk.Add(elastic.NewBulkIndexRequest().Id(id).Doc(doc))
	}
	_, err := bulk.Do(context.Background())
	if err != nil {
		return ErrEsBulkIdx(err, index)
	}
	// refresh
	if s.cfg.Refresh {
		return s.Refresh(index)
	}
	return nil
}

func (s *esImpl) IndexBulkAsync(index string, docs map[string]interface{}) {
	goroutine.New().
		WithLogger(s.l().Mth("index-bulk-async")).
		Go(context.Background(), func() {
			l := s.l().Mth("bulk-indexation-async").F(log.FF{"index": index, "docs": len(docs)}).Dbg()
			err := s.IndexBulk(index, docs)
			if err != nil {
				l.E(err).Err()
			}
		})
}

// Exists checks if a document exists in the index
func (s *esImpl) Exists(index, id string) (bool, error) {
	l := s.l().Mth("exists").F(log.FF{"index": index, "id": id})
	res, err := s.client.Exists().Index(index).Id(id).Do(context.Background())
	if err != nil {
		return false, ErrEsExists(err, index, id)
	}
	l.DbgF("res: %v", res)
	return res, nil
}

func (s *esImpl) Delete(index string, id string) error {
	s.l().Mth("delete").F(log.FF{"index": index, "id": id}).Dbg()
	svc := s.client.
		Delete().
		Index(index).
		Id(id)
	_, err := svc.Do(context.Background())
	if err != nil {
		return ErrEsDel(err, index, id)
	}
	// refresh
	if s.cfg.Refresh {
		return s.Refresh(index)
	}
	return nil
}

func (s *esImpl) DeleteBulk(index string, ids []string) error {
	s.l().Mth("bulk-deletion").F(log.FF{"index": index, "ids": len(ids)}).Dbg()
	bulk := s.client.Bulk().Index(index)
	for _, id := range ids {
		bulk.Add(elastic.NewBulkDeleteRequest().Id(id))
	}
	_, err := bulk.Do(context.Background())
	if err != nil {
		return ErrEsBulkDel(err, index)
	}
	// refresh
	if s.cfg.Refresh {
		return s.Refresh(index)
	}
	return nil
}

func (s *esImpl) GetClient() *elastic.Client {
	return s.client
}

func (s *esImpl) Close() {
	s.client.Stop()
}

func (s *esImpl) Refresh(index string) error {
	_, err := s.client.Refresh(index).Do(context.Background())
	if err != nil {
		return ErrEsRefresh(err, index)
	}
	return nil
}
