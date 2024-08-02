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
// "commons.Mainnet" or "commons.Arabica" or "commons.Mocha"
celestiaNet := commons.Mocha      
// "commons.ArchivalNode" or "commons.ArchivalNode"
recordKey := commons.ArchivalNode
```

## running the crawler
```
go run dht-crawler/crawler.go
```

## running the lookup
```
go run dht-lookup/lookup.go
```