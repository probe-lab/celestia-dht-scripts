package main

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	kad "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/core/protocol"
	routingdisc "github.com/libp2p/go-libp2p/p2p/discovery/routing"

	"github.com/probe-lab/celestia-dht-scripts/commons"
)

func main() {
	// "commons.Mainnet" or "commons.Arabica" or "commons.Mocha"
	celestiaNet := commons.Mocha
	// "commons.ArchivalNode" or "commons.ArchivalNode"
	recordKey := commons.ArchivalNode

	ctx := context.Background()

	h, err := libp2p.New(
		libp2p.UserAgent("celestia-celestia"),
		libp2p.Identity(commons.LoadPrivKey()),
		//libp2p.NATPortMap(), // enable upnp
		libp2p.DisableRelay(),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Peer ID:", h.ID())

	dhtOpts := []kad.Option{
		kad.Mode(kad.ModeClient),
		kad.BootstrapPeers(commons.BootstrapPeers(celestiaNet)...),
		kad.ProtocolPrefix(protocol.ID(fmt.Sprintf("/celestia/%s", celestiaNet))),
	}
	dht, err := kad.New(ctx, h, dhtOpts...)
	if err != nil {
		panic(err)
	}

	for _, bootstrapper := range commons.BootstrapPeers(celestiaNet) {
		if err := h.Connect(ctx, bootstrapper); err != nil {
			fmt.Println("couldn't connect to", bootstrapper, ":", err)
		}
	}

	dht.Bootstrap(ctx)

	time.Sleep(5 * time.Second)

	fmt.Println("Routing table size:", dht.RoutingTable().Size())

	disc := routingdisc.NewRoutingDiscovery(dht)

	findCtx, cancelFunc := context.WithTimeout(ctx, 15*time.Second)
	defer cancelFunc()

	peers, err := disc.FindPeers(findCtx, recordKey.String(), discovery.Limit(0))
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	fmt.Println("Found peers:")
	c := 0
	for p := range peers {
		fmt.Println(p.ID)
		c += 1
	}
	fmt.Println("Total peers found:", c)
}
