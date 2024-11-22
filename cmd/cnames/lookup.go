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
	Name:   "scan",
	Usage:  "estimates the uplink BW from the active list of nodes in the network",
	Flags:  cmdLookupFlags,
	Action: cmdLookupAction,
}

var cmdLookupFlags = []cli.Flag{
	&cli.StringFlag{
		Name: "network",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_NETWORK")},
		},
		Usage:       "TODO",
		Value:       lookupConfig.Network,
		Destination: &lookupConfig.Network,
	},
	&cli.BoolFlag{
		Name: "is-custom",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_IS_CUSTOM")},
		},
		Usage:       "TODO",
		Value:       lookupConfig.IsCustomNamespace,
		Destination: &lookupConfig.IsCustomNamespace,
	},
	&cli.StringFlag{
		Name: "namespace",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_NAMESPACE")},
		},
		Usage:       "TODO",
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
	defer log.Infof("stopped cnames-lookup")

	network := dht.NetworkFromString(lookupConfig.Network)
	kadProtocol := network.KadProtocol()

	h, err := libp2p.New(
		libp2p.UserAgent(dht.CustomUserAgent),
		libp2p.Identity(dht.LoadPrivKey()),
		// libp2p.NATPortMap(), // enable upnp
		libp2p.DisableRelay(),
	)
	if err != nil {
		panic(err)
	}

	dhtOpts := []kad.Option{
		kad.Mode(kad.ModeClient),
		kad.BootstrapPeers(dht.BootstrapPeers(network)...),
		kad.ProtocolPrefix(kadProtocol),
	}
	dhtCli, err := kad.New(ctx, h, dhtOpts...)
	if err != nil {
		panic(err)
	}

	for _, bootstrapper := range dht.BootstrapPeers(network) {
		if err := h.Connect(ctx, bootstrapper); err != nil {
			log.Warn("couldn't connect to", bootstrapper, ":", err)
		}
	}

	log.Info("HOST info:")
	log.Info("- Peer ID:      ", h.ID())
	log.Info("- Network:      ", network)
	log.Info("- Protocols:    ", h.Mux().Protocols())
	log.Info("- Agent Version:", dht.CustomUserAgent)

	dhtCli.Bootstrap(ctx)

	time.Sleep(5 * time.Second)

	log.Info("Routing table size:", dhtCli.RoutingTable().Size())

	disc := routingdisc.NewRoutingDiscovery(dhtCli)

	findCtx, cancelFunc := context.WithTimeout(ctx, 15*time.Second)
	defer cancelFunc()

	peers, err := disc.FindPeers(findCtx, lookupConfig.Namespace, discovery.Limit(0))
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	log.Info("Found peers:")
	c := 0
	for p := range peers {
		log.Info(p.ID)
		c += 1
	}
	log.Info("Total peers found:", c)

	return nil
}
