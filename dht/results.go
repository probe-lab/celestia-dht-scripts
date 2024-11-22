package dht

import (
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

type CrawlResults struct {
	m                sync.RWMutex
	succPeers        map[peer.ID]peer.AddrInfo
	failedPeers      map[peer.ID]peer.AddrInfo
	provPeers        map[peer.ID]peer.AddrInfo
	agentVersionDist map[string]int
	initTime         time.Time
	finishTime       time.Time
}

func NewCrawlerResults() *CrawlResults {
	return &CrawlResults{
		succPeers:        make(map[peer.ID]peer.AddrInfo),
		failedPeers:      make(map[peer.ID]peer.AddrInfo),
		provPeers:        make(map[peer.ID]peer.AddrInfo),
		agentVersionDist: make(map[string]int),
	}
}

func (r *CrawlResults) addAgentVersion(av string) {
	r.m.Lock()
	defer r.m.Unlock()

	// if the peer wasn't already in the map, add it straight away
	_, ok := r.agentVersionDist[av]
	if !ok {
		// add it straight away
		r.agentVersionDist[av] = 1
	} else {
		r.agentVersionDist[av]++
	}
}

func (r *CrawlResults) addSuccessfullPeer(p peer.ID, ai peer.AddrInfo) {
	r.m.Lock()
	defer r.m.Unlock()

	// if the peer wasn't already in the map, add it straight away
	_, ok := r.succPeers[p]
	if !ok {
		// add it straight away
		r.succPeers[p] = ai
	}
}

func (r *CrawlResults) addProvider(p peer.ID, ai peer.AddrInfo) {
	r.m.Lock()
	defer r.m.Unlock()

	// if the peer wasn't already in the map, add it straight away
	_, ok := r.provPeers[p]
	if !ok {
		// add it straight away
		r.provPeers[p] = ai
	}
}

func (r *CrawlResults) addFailedPeer(p peer.ID, ai peer.AddrInfo) {
	r.m.Lock()
	defer r.m.Unlock()

	// if the peer wasn't already in the map, add it straight away
	_, ok := r.failedPeers[p]
	if !ok {
		// add it straight away
		r.failedPeers[p] = ai
	}
}

// retrievals
func (r *CrawlResults) GetSuccPeers() map[peer.ID]peer.AddrInfo {
	r.m.RLock()
	defer r.m.RUnlock()

	total := make(map[peer.ID]peer.AddrInfo)

	for k, v := range r.succPeers {
		total[k] = v
	}
	return total
}

func (r *CrawlResults) GetProvPeers() map[peer.ID]peer.AddrInfo {
	r.m.RLock()
	defer r.m.RUnlock()

	total := make(map[peer.ID]peer.AddrInfo)

	for k, v := range r.provPeers {
		total[k] = v
	}
	return total
}

func (r *CrawlResults) GetFailedPeers() map[peer.ID]peer.AddrInfo {
	r.m.RLock()
	defer r.m.RUnlock()

	total := make(map[peer.ID]peer.AddrInfo)

	for k, v := range r.failedPeers {
		total[k] = v
	}
	return total
}

func (r *CrawlResults) GetAgentDistributions() map[string]int {
	r.m.RLock()
	defer r.m.RUnlock()

	final := make(map[string]int)
	total := 0

	for k, v := range r.agentVersionDist {
		final[k] = v
		total = total + v
	}
	final["total"] = total
	return final
}

func (c *CrawlResults) GetCrawlerDuration() time.Duration {
	return c.finishTime.Sub(c.initTime)
}
