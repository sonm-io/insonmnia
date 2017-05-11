package daemon

import (
	"context"
	"net"

	gcontext "golang.org/x/net/context"

	"google.golang.org/grpc"

	pb "github.com/sonm-io/insonmnia/daemon/miner"
	"github.com/sonm-io/insonmnia/insonmnia"
)

// Daemon accepts incoming requests to control apps
type Daemon struct {
	ctx    context.Context
	cancel context.CancelFunc

	ovr insonmnia.Overseer

	grpcServer *grpc.Server
}

var _ pb.MinerServer = &Daemon{}

// New creates new daemon
func New(ctx context.Context) (*Daemon, error) {
	// NOTE: set actual options
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	ctx, cancel := context.WithCancel(ctx)
	ovr, err := insonmnia.NewOverseer(ctx)
	if err != nil {
		cancel()
		return nil, err
	}

	dmn := Daemon{
		ctx:    ctx,
		cancel: cancel,

		grpcServer: grpcServer,
		ovr:        ovr,
	}

	pb.RegisterMinerServer(grpcServer, &dmn)

	return &dmn, nil
}

// Serve starts request handling
func (d *Daemon) Serve(ln net.Listener) error {
	return d.grpcServer.Serve(ln)
}

// Spawn starts new task
func (d *Daemon) Spawn(ctx gcontext.Context, req *pb.SpawnRequest) (*pb.SpawnReply, error) {
	var (
		description insonmnia.Description
		reply       = new(pb.SpawnReply)
	)

	// NOTE: wrap into status.Errorf
	if err := d.ovr.Spawn(ctx, description); err != nil {
		return nil, err
	}

	return reply, nil
}

// Close releases all resources associated with Daemon
func (d *Daemon) Close() {
	// TODO: implement graceful shutdown via context & SIGNAL handling
	d.grpcServer.Stop()
	d.ovr.Close()
}
