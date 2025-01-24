package main

import (
	"context"

	"github.com/probe-lab/celestia-dht-scripts/dht"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var dhtKeysConf struct{
	Key string
}

var cmdDHTKeys = &cli.Command{
	Name:   "key-info",
	Usage:  "show all info for the given DHT key",
	Flags:  []cli.Flag{
		&cli.StringFlag{
			Name: "key",
			Required: true,
			Usage:       "target key",
			Value:       dhtKeysConf.Key,
			Destination: &dhtKeysConf.Key,
		},
	},
	Action: cmdKeyAction,
}


func cmdKeyAction(ctx context.Context, c *cli.Command) error {
	recordCid, err := dht.KeyToCid(dhtKeysConf.Key)
	if err != nil {
		return nil
	}
	hash := recordCid.Hash()
	logrus.Info("key:            ", dhtKeysConf.Key)
	logrus.Info("b58-cid:        ", recordCid.String())
	logrus.Info("key-string-cid: ", recordCid.KeyString())
	logrus.Info("b58-hash-cid:   ", hash.B58String())
	logrus.Info("hex-hash-cid:   ", hash.HexString())
	return nil
}

