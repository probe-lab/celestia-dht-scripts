# celestia-dht-scripts
Set of scripts to interact with the Celestia DHT to monitor the PRs for `archival` and `full` nodes

## requiremets 
1. go installed (v0.22.2 works fine)
2. install the dependencies:
```
go get ./...
```
3. configure the network and the lookup key at each script
```go
// at ./dht-crawler/crawler.go or ./dht-lookup/lookup.go  
celestiaNet := commons.Mocha      // "commons.Mainnet" or "commons.Arabica" or "commons.Mocha"
recordKey := commons.ArchivalNode // "commons.ArchivalNode" or "commons.ArchivalNode"
```

## running the crawler
```
go run dht-crawler/crawler.go
```

## running the lookup
```
go run dht-lookup/lookup.go
```