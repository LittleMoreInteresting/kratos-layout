package taskserver

import (
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/robfig/cron/v3"
	"golang.org/x/net/context"
)

var (
	_ transport.Server = (*Server)(nil)
)

type Server struct {
	sync.RWMutex
	Cron       *cron.Cron
	log        *log.Helper
	baseCtx    context.Context
	MapEntryId map[string]cron.EntryID
	started    bool
	err        error
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		Cron: cron.New(cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		)),
		MapEntryId: make(map[string]cron.EntryID),
		started:    false,
	}
	srv.init(opts...)
	return srv
}

func (s *Server) init(opts ...ServerOption) {
	for _, o := range opts {
		o(s)
	}
}

func (s *Server) Name() string {
	return "task-cron"
}

func (s *Server) RegisterMapEntry(id cron.EntryID, taskName string) {
	s.Lock()
	defer s.Unlock()
	if len(taskName) > 0 {
		s.MapEntryId[taskName] = id
		s.log.WithContext(s.baseCtx).Infow("[task-cron] register task:", taskName, "id:", id)
	}
}

func (s *Server) GetMapEntry() map[string]cron.EntryID {
	return s.MapEntryId
}

func (s *Server) AddFunc(spec string, fun func()) (cron.EntryID, error) {
	return s.Cron.AddFunc(spec, fun)
}
func (s *Server) AddJob(spec string, job cron.Job) (cron.EntryID, error) {
	return s.Cron.AddJob(spec, job)
}

func (s *Server) Start(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}
	if s.started {
		return nil
	}
	s.Cron.Start()
	s.log.WithContext(ctx).Info("[task-corn] server starting")
	s.baseCtx = ctx
	s.started = true
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.log.WithContext(s.baseCtx).Info("[task-cron] server stopping")
	s.started = false
	s.Cron.Stop()
	return nil
}

type TaskServer interface {
	Schedule() string
	TaskName() string
	TaskFunc(ctx context.Context) func()
}

func RegisterServer(srv *Server, s TaskServer) {
	RegisterFunc(srv, s.Schedule(), s.TaskName(), s.TaskFunc(srv.baseCtx))
}
func RegisterFunc(srv *Server, spec string, taskName string, fun func()) {
	entryId, err := srv.AddFunc(spec, fun)
	if err != nil {
		srv.log.WithContext(srv.baseCtx).Errorw("register func is error")
	}
	srv.RegisterMapEntry(entryId, taskName)
}
