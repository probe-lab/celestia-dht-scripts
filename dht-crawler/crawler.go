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

	bc "github.com/probe-lab/celestia-dht-scripts/base-crawler"
	"github.com/probe-lab/celestia-dht-scripts/commons"
)

func main() {
	// "commons.Mainnet" or "commons.Arabica" or "commons.Mocha"
	celestiaNet := commons.Mocha
	// "commons.ArchivalNode" or "commons.ArchivalNode"
	recordKey := commons.ArchivalNode

	celestiaProtocolBase := protocol.ID(fmt.Sprintf("/celestia/%s", celestiaNet))
	celestiaKadProtocol := protocol.ID(fmt.Sprintf("%s/kad/1.0.0", string(celestiaProtocolBase)))

	fmt.Println(celestiaKadProtocol)

	ctx := context.Background()

	// get bootstrappers
	bootstrapers := commons.BootstrapPeers(celestiaNet)
	startingPeers := make([]*peer.AddrInfo, len(bootstrapers))
	for idx, bootstrapper := range bootstrapers {
		startingPeers[idx] = &bootstrapper
	}

	// libp2p host
	h, err := libp2p.New(
		libp2p.UserAgent("celestia-celestia"),
		libp2p.Identity(commons.LoadPrivKey()),
		// libp2p.NATPortMap(), // enable upnp
		libp2p.DisableRelay(),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Peer ID of host:", h.ID())

	// protocol messenger for the DHT queries
	prots := []protocol.ID{celestiaKadProtocol}
	pm, err := pb.NewProtocolMessenger(&bc.MessageSender{H: h, Protocols: []protocol.ID{celestiaKadProtocol}, Timeout: 30 * time.Second})
	if err != nil {
		panic(err)
	}

	// for the crawler
	dhtCrawler, err := bc.New(h, prots, pm)
	if err != nil {
		panic(err)
	}

	results := dhtCrawler.Run(ctx, startingPeers, recordKey.String())

	succPeers := results.GetSuccPeers()
	failedPeers := results.GetFailedPeers()
	providers := results.GetProvPeers()
	agentVersions := results.GetAgentDistributions()

	fmt.Printf("Found %s nodes:\n", recordKey.String())
	for p := range providers {
		fmt.Println(p)
	}

	fmt.Printf(`Summary of the crawl on %s:
- Duration: %s
- Total discovered nodes: %d
- Successful connected nodes: %d
- Failed to connect nodes: %d
- Advertised %s nodes: %d
- AgentVersion distribution:`,
		celestiaNet,
		results.GetCrawlerDuration(),
		len(succPeers)+len(failedPeers),
		len(succPeers),
		len(failedPeers),
		recordKey.String(),
		len(providers),
	)
	printTable(agentVersions)
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
	fmt.Printf("\n%-*s | nodes\n", maxKeyLength, "agent_version")
	fmt.Println(strings.Repeat("-", maxKeyLength+8))

	// Print key-value pairs
	for key, value := range data {
		if key == "total" {
			continue
		}
		fmt.Printf("%-*s | %v\n", maxKeyLength, key, value)
	}
	fmt.Println(strings.Repeat("-", maxKeyLength+8))
	fmt.Printf("%-*s | %v\n", maxKeyLength, "total", data["total"])
}
