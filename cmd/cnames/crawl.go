package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v3"

	"github.com/probe-lab/celestia-dht-scripts/dht"
)

var crawlConfig = dht.LookupCmdConfig{
	Network:           dht.DefaultNetwork.String(),
	IsCustomNamespace: dht.DefaultIsNamespace,
	Namespace:         dht.DefaultNamespace.String(),
}

var cmdCrawl = &cli.Command{
	Name:   "scan",
	Usage:  "estimates the uplink BW from the active list of nodes in the network",
	Flags:  cmdCrawlFlags,
	Action: cmdCrawlAction,
}

var cmdCrawlFlags = []cli.Flag{
	&cli.StringFlag{
		Name: "network",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_NETWORK")},
		},
		Usage:       "TODO",
		Value:       crawlConfig.Network,
		Destination: &crawlConfig.Network,
	},
	&cli.BoolFlag{
		Name: "is-custom",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_IS_CUSTOM")},
		},
		Usage:       "TODO",
		Value:       crawlConfig.IsCustomNamespace,
		Destination: &crawlConfig.IsCustomNamespace,
	},
	&cli.StringFlag{
		Name: "namespace",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_NAMESPACE")},
		},
		Usage:       "TODO",
		Value:       crawlConfig.Namespace,
		Destination: &crawlConfig.Namespace,
	},
}

func cmdCrawlAction(ctx context.Context, cmd *cli.Command) error {
	log.WithFields(log.Fields{
		"network":      crawlConfig.Network,
		"is-custom-ns": crawlConfig.IsCustomNamespace,
		"namespace":    crawlConfig.Namespace,
	}).Info("starting cnames-crawl...")
	defer log.Infof("stopped cnames-crawl")

	network := dht.NetworkFromString(crawlConfig.Network)
	kadProtocol := network.KadProtocol()

	// get bootstrappers
	bootstrapers := dht.BootstrapPeers(network)
	startingPeers := make([]*peer.AddrInfo, len(bootstrapers))
	for idx, bootstrapper := range bootstrapers {
		startingPeers[idx] = &bootstrapper
	}

	// libp2p host
	h, err := libp2p.New(
		libp2p.UserAgent(dht.CustomUserAgent),
		libp2p.Identity(dht.LoadPrivKey()),
		// libp2p.NATPortMap(), // enable upnp
		libp2p.DisableRelay(),
	)
	if err != nil {
		return err
	}
	log.Info("HOST info:")
	log.Info("- Peer ID:      ", h.ID())
	log.Info("- Network:      ", network)
	log.Info("- Protocol:     ", kadProtocol)
	log.Info("- Protocols:    ", h.Mux().Protocols())
	log.Info("- Agent Version:", dht.CustomUserAgent)

	// protocol messenger for the DHT queries
	prots := []protocol.ID{kadProtocol}
	pm, err := pb.NewProtocolMessenger(&dht.MessageSender{H: h, Protocols: []protocol.ID{kadProtocol}, Timeout: 30 * time.Second})
	if err != nil {
		panic(err)
	}

	// for the crawler
	dhtCrawler, err := dht.New(h, prots, pm)
	if err != nil {
		panic(err)
	}

	results := dhtCrawler.Run(ctx, startingPeers, crawlConfig.Namespace)

	succPeers := results.GetSuccPeers()
	failedPeers := results.GetFailedPeers()
	providers := results.GetProvPeers()
	agentVersions := results.GetAgentDistributions()

	log.Info("Found %s nodes:\n", crawlConfig.Namespace)
	for p := range providers {
		fmt.Println(p)
	}

	log.Infof("Summary of the crawl on %s:", network)
	log.Infof(" - Duration: %s", results.GetCrawlerDuration())
	log.Infof(" - Total discovered nodes: %d", len(succPeers)+len(failedPeers))
	log.Infof(" - Successful connected nodes: %d", len(succPeers))
	log.Infof(" - Failed to connect nodes: %d", len(failedPeers))
	log.Infof(" - Advertised %s nodes: %d", crawlConfig.Namespace)
	log.Infof(" - AgentVersion distribution:", len(providers))

	printTable(agentVersions)

	return nil
}

func printTable(data map[string]int) {
	// Determine the maximum key length for formatting
	var maxKeyLength int
	for key := range data {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
	}

	// Print headerd
	log.Infof("\n%-*s | nodes\n", maxKeyLength, "agent_version")
	log.Info(strings.Repeat("-", maxKeyLength+8))

	// Print key-value pairs
	for key, value := range data {
		if key == "total" {
			continue
		}
		log.Infof("%-*s | %v\n", maxKeyLength, key, value)
	}
	log.Info(strings.Repeat("-", maxKeyLength+8))
	log.Infof("%-*s | %v\n", maxKeyLength, "total", data["total"])
}
