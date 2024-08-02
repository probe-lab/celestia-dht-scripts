package base_crawler

import (
	"context"
	"time"

	pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-msgio/protoio"
)

// MessageSender handles sending wire protocol messages to a given peer
type MessageSender struct {
	H         host.Host
	Protocols []protocol.ID
	Timeout   time.Duration
}

// SendRequest sends a peer a message and waits for its response
func (ms *MessageSender) SendRequest(ctx context.Context, p peer.ID, pmes *pb.Message) (*pb.Message, error) {
	s, err := ms.H.NewStream(ctx, p, ms.Protocols...)
	if err != nil {
		return nil, err
	}

	w := protoio.NewDelimitedWriter(s)
	if err := w.WriteMsg(pmes); err != nil {
		return nil, err
	}

	r := protoio.NewDelimitedReader(s, network.MessageSizeMax)
	tctx, cancel := context.WithTimeout(ctx, ms.Timeout)
	defer cancel()
	defer func() { _ = s.Close() }()

	msg := new(pb.Message)
	if err := ctxReadMsg(tctx, r, msg); err != nil {
		_ = s.Reset()
		return nil, err
	}

	return msg, nil
}

func ctxReadMsg(ctx context.Context, rc protoio.ReadCloser, mes *pb.Message) error {
	errc := make(chan error, 1)
	go func(r protoio.ReadCloser) {
		defer close(errc)
		err := r.ReadMsg(mes)
		errc <- err
	}(rc)

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// SendMessage sends a peer a message without waiting on a response
func (ms *MessageSender) SendMessage(ctx context.Context, p peer.ID, pmes *pb.Message) error {
	s, err := ms.H.NewStream(ctx, p, ms.Protocols...)
	if err != nil {
		return err
	}
	defer func() { _ = s.Close() }()

	w := protoio.NewDelimitedWriter(s)
	return w.WriteMsg(pmes)
}
