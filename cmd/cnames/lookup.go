package main

import (
	"context"
	"time"

	kad "github.com/libp2p/go-libp2p-kad-dht"
	routingdisc "github.com/libp2p/go-libp2p/p2p/discovery/routing"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/probe-lab/celestia-dht-scripts/dht"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v3"
)

var lookupConfig = dht.LookupCmdConfig{
	Network:           dht.DefaultNetwork.String(),
	IsCustomNamespace: dht.DefaultIsNamespace,
	Namespace:         dht.DefaultNamespace.String(),
}

var cmdLookup = &cli.Command{
	Name:   "lookup",
	Usage:  "makes a DHT lookup for the given namespace",
	Flags:  cmdLookupFlags,
	Action: cmdLookupAction,
}

var cmdLookupFlags = []cli.Flag{
	&cli.StringFlag{
		Name: "network",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_NETWORK")},
		},
		Usage:       "celestia network where the cname will run",
		Value:       lookupConfig.Network,
		Destination: &lookupConfig.Network,
	},
	&cli.BoolFlag{
		Name: "is-custom",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_IS_CUSTOM")},
		},
		Usage:       "is the namespace custom?",
		Value:       lookupConfig.IsCustomNamespace,
		Destination: &lookupConfig.IsCustomNamespace,
	},
	&cli.StringFlag{
		Name: "namespace",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_NAMESPACE")},
		},
		Usage:       "DHT key or namespace the will be searched",
		Value:       lookupConfig.Namespace,
		Destination: &lookupConfig.Namespace,
	},
}

func cmdLookupAction(ctx context.Context, cmd *cli.Command) error {
	log.WithFields(log.Fields{
		"network":      lookupConfig.Network,
		"is-custom-ns": lookupConfig.IsCustomNamespace,
		"namespace":    lookupConfig.Namespace,
	}).Info("starting cnames-lookup...")

	network := dht.NetworkFromString(lookupConfig.Network)
	kadProtocol := network.KadPrefix()

	h, err := libp2p.New(
		libp2p.UserAgent(dht.CustomUserAgent),
		libp2p.Identity(dht.LoadPrivKey()),
		// libp2p.NATPortMap(), // enable upnp
		libp2p.DisableRelay(),
	)
	if err != nil {
		return err
	}

	dhtOpts := []kad.Option{
		kad.Mode(kad.ModeClient),
		kad.BootstrapPeers(dht.BootstrapPeers(network)...),
		kad.ProtocolPrefix(kadProtocol),
	}
	dhtCli, err := kad.New(ctx, h, dhtOpts...)
	if err != nil {
		return err
	}

	bootnodes := 0
	for _, bootstrapper := range dht.BootstrapPeers(network) {
		if err := h.Connect(ctx, bootstrapper); err != nil {
			log.Warn("couldn't connect to", bootstrapper, ":", err)
		} else {
			bootnodes++
		}
	}

	log.Info("HOST info:")
	log.Info("- Peer ID:			", h.ID())
	log.Info("- Network:			", network)
	log.Info("- Protocols:			", h.Mux().Protocols())
	log.Info("- Agent Version:		", dht.CustomUserAgent)
	log.Info("- Bootnodes:			", bootnodes)

	err = dhtCli.Bootstrap(ctx)
	if err != nil {
		return nil
	}
	time.Sleep(5 * time.Second)

	log.Info("- Routing table size:	", dhtCli.RoutingTable().Size())

	disc := routingdisc.NewRoutingDiscovery(dhtCli)

	findCtx, cancelFunc := context.WithTimeout(ctx, 15*time.Second)
	defer cancelFunc()

	peers, err := disc.FindPeers(findCtx, lookupConfig.Namespace, discovery.Limit(0))
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	log.Info("Found peers:")
	c := 1
	for p := range peers {
		log.Infof("%d -> peer_id: %s", c, p.ID.String())
		c += 1
	}
	log.Info("Total peers found:", c-1)

	return nil
}
