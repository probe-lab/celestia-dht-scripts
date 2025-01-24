# cnames (celestia-dht-scripts)
CLI tool to interact with the Celestia DHT and monitor the PRs for `archival` and `full` nodes.

## Installation
There are two options to build or install the cnames tool:

1. Simple way of installing it:
```
make install
cnames lookup
```


2. Simple way of building it:
```
make build
./build/cnames lookup
```

## How to use it

Root options:
```
NAME:
   cnames - A Celestia's DHT namespace scrapper

USAGE:
   cnames [global options] [command [command options]]

COMMANDS:
   lookup   TODO
   crawl    estimates the uplink BW from the active list of nodes in the network
   key-info  show all info for the given DHT key
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

   Logging Configuration:

   --log.format value  Sets the format to output the log statements in: text, json (default: "text") [$CNAMES_LOG_FORMAT]
   --log.level value   Sets an explicity logging level: debug, info, warn, error. Takes precedence over the verbose flag. (default: "info") [$CNAMES_LOG_LEVEL]
```

Subcommands:
1. `lookup`: makes a DHT lookup for the given namespace
2. `crawl`: asks each node in the network for the PRs they have asociated with the namespace

```
NAME:
   cnames lookup/crawl - makes a DHT lookup for the given namespace

USAGE:
   cnames lookup/crawl [command [command options]]

OPTIONS:
   --network value    celestia network where the cname will run (default: "celestia") [$CNAMES_NETWORK]
   --is-custom        is the namespace custom? (default: false) [$CNAMES_IS_CUSTOM]
   --namespace value  DHT key or namespace the will be searched (default: "/full/v0.1.0") [$CNAMES_NAMESPACE]
   --help, -h         show help

GLOBAL OPTIONS:
   --log.level value   Sets an explicity logging level: debug, info, warn, error. Takes precedence over the verbose flag. (default: "info") [$CNAMES_LOG_LEVEL]
   --log.format value  Sets the format to output the log statements in: text, json (default: "text") [$CNAMES_LOG_FORMAT]
```

3. `key.info`: returns all the different formatting types for a given DHT Key (CID and Hash of the CID)  

```
NAME:
   cnames key-info - show all info for the given DHT key

USAGE:
   cnames key-info [command [command options]]

OPTIONS:
   --key value  target key
   --help, -h   show help

GLOBAL OPTIONS:
   --log.level value   sets an explicity logging level: debug, info, warn, error. Takes precedence over the verbose flag. (default: "info") [$CNAMES_LOG_LEVEL]
   --log.format value  sets the format to output the log statements in: text, json (default: "text") [$CNAMES_LOG_FORMAT]
```


