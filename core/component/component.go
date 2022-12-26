package component

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/smallnest/rpcx/server"
	"golang.org/x/exp/slog"
)

// type definitions
type Description struct {
	Name     string   `json:"name"`
	Version  int32    `json:"version"`
	Protocol string   `json:"protocol"`
	Addr     string   `json:"address"`
	Tags     []string `json:"tags"`
}

type Config struct {
	DaemonPort int `json:"daemon_port" msg:"daemon_port"`
}

// 用于发现组件
type Discoverer struct {
}
type EventHandler interface {
	SetDiscoverer(Discoverer)
	OnStart()
	OnShutdown()
}

type Component struct {
	server *server.Server
	logger *slog.Logger

	eventHandler EventHandler

	description Description
}

type Option func(c *Component)

type StopReply struct {
	Code   int    `json:"code"`
	Remark string `json:"remark"`
}

//go:generate stringer -type ShutdownReason -linecomment
type ShutdownReason int

// constants
const (
	ShutdownReasonUninstall ShutdownReason = iota // 卸载组件, 永久下线
	ShutdownReasonRemaining                       // 维护组件,暂时下线
)

const (
	kDefaultProtocol = "tcp"
	kDefaultAddress  = "localhost:4096"
)

// global functions
func New(eventHandler EventHandler, opts ...Option) *Component {
	c := Component{
		server:       server.NewServer(),
		eventHandler: eventHandler,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func WithNameAndVersion(name string, version int32) Option {
	return func(c *Component) {
		c.description.Name = name
		c.description.Version = version
	}
}

func WithProtocolAndAddr(protocol string, addr string) Option {
	return func(c *Component) {
		c.description.Protocol = protocol
		c.description.Addr = addr
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *Component) {
		c.logger = logger
	}
}

func WithTags(tags []string) Option {
	return func(c *Component) {
		c.description.Tags = tags
	}
}

func Serve(eventHandler EventHandler, opts ...Option) error {
	c := New(eventHandler, opts...)

	return c.Serve()
}

// methods

// Component public methods
func (c *Component) Serve() error {
	c.init()

	c.Logger().Info("component started", c.logWithDescription()...)
	if c.eventHandler != nil {
		c.eventHandler.OnStart()
	}
	err := c.server.Serve(c.description.Protocol, c.description.Addr)
	if err != nil && err != server.ErrServerClosed {

		c.Logger().Error("failed to serve component", err)
		return fmt.Errorf("failed to serve component. %s", err)
	}
	c.Logger().Info("component finished")
	return nil
}

func (c *Component) Logger() *slog.Logger {
	if c.logger == nil {
		c.logger = slog.New(slog.NewJSONHandler(os.Stdout))
	}

	return c.logger
}

// Description - rpc functions
func (c *Component) Description(ctx context.Context, reserve int, reply *Description) error {
	*reply = c.description
	return nil
}

func (c *Component) Shutdown(ctx context.Context, reason ShutdownReason, reply *StopReply) error {
	c.Logger().Info("component shutdown", "reason", reason.String())
	reply.Code = int(reason)
	reply.Remark = "ok"
	if c.eventHandler != nil {
		c.eventHandler.OnShutdown()
	}
	time.AfterFunc(3*time.Second, func() {
		c.server.Shutdown(context.TODO())
	})
	return nil
}

// Component private methods
func (c *Component) init() {
	if c.description.Protocol == "" {
		c.description.Protocol = kDefaultProtocol
	}
	if c.description.Addr == "" {
		c.description.Addr = kDefaultAddress
	}

	c.server.DisableJSONRPC = true
	c.server.DisableHTTPGateway = true
	c.server.RegisterName("Component", c, "")

	// init 最后添加日志上下文(组件描述信息)
	// example:
	// logger := c.Logger()
	// c.logger = logger.With(c.logWithDescription()...)
}

func (c *Component) logWithDescription(keyValues ...any) []any {
	metaData := []any{
		"name", c.description.Name,
		"version", c.description.Version,
		"protocol", c.description.Protocol,
		"address", c.description.Addr,
	}
	if len(keyValues) == 0 {
		return metaData
	}

	return append(metaData, keyValues...)
}

// Discoverer functions
func (d *Discoverer) FindComponent() {

}
