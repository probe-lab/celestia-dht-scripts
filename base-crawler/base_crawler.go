package base_crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-kad-dht/crawler"
	pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	mh "github.com/multiformats/go-multihash"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// Crawler is simply a wrappper on top of the official go-lip2p-kad-dht/crawler
// It is mainly used to crawl the network searching for Hydra-Booster peers.
// with the final intention of blacklisting them in the Kad-Dht process if the Avoid-Hydras flags gets triggered
type BaseCrawler struct {
	crawler *crawler.DefaultCrawler
	h       host.Host
	pm      *pb.ProtocolMessenger
	results *CrawlResults
}

func New(h host.Host, ptcls []protocol.ID, pm *pb.ProtocolMessenger) (*BaseCrawler, error) {
	// create the official crawler
	c, err := crawler.NewDefaultCrawler(
		h,
		crawler.WithParallelism(300),
		crawler.WithMsgTimeout(10*time.Second),
		crawler.WithConnectTimeout(10*time.Second),
		crawler.WithProtocols(ptcls),
	)
	if err != nil {
		return nil, err
	}

	return &BaseCrawler{
		crawler: c,
		h:       h,
		pm:      pm,
		results: NewCrawlerResults(),
	}, nil
}

func (c *BaseCrawler) Run(ctx context.Context, startingNodes []*peer.AddrInfo, recordKey string) *CrawlResults {
	recordCid, err := keyToCid(recordKey)
	if err != nil {
		return nil
	}

	// set up the handle Success function for the crawler
	handleSucc := func(p peer.ID, rtPeers []*peer.AddrInfo) {
		c.results.addSuccessfullPeer(p, peer.AddrInfo{})

		// get the agent version from the peer
		av := "unknown"
		avIntf, err := c.h.Peerstore().Get(p, "AgentVersion")
		if err == nil {
			av = avIntf.(string)
		}
		c.results.addAgentVersion(av)

		// fmt.Printf("successfull connection to %s, requesting provs for %s key\n", p.String(), recordKey)
		// on each successfull connection, request the PRs from the key
		provs, _, err := c.pm.GetProviders(ctx, p, recordCid.Hash())
		if err != nil {
			return
		}
		if len(provs) > 0 {
			for _, provider := range provs {
				c.results.addProvider(provider.ID, *provider)
			}
			fmt.Printf("peer %s reported %d providers for %s nodes\n", p.String(), len(provs), recordKey)
		}
	}

	// set up the handle Fail function for the crawler
	handleFail := func(p peer.ID, err error) {
		c.results.addFailedPeer(p, peer.AddrInfo{})
	}

	c.results.initTime = time.Now()
	c.crawler.Run(ctx, startingNodes, handleSucc, handleFail)
	c.results.finishTime = time.Now()

	return c.results
}

func (c *BaseCrawler) Close() {
	// check is there is host
	if c.h == nil {
		return
	}
	c.h.Close()
}

func keyToCid(ns string) (cid.Cid, error) {
	h, err := mh.Sum([]byte(ns), mh.SHA2_256, -1)
	if err != nil {
		return cid.Undef, err
	}

	return cid.NewCidV1(cid.Raw, h), nil
}
