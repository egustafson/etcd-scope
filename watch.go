package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	etcd "go.etcd.io/etcd/client/v3" // <-- v3.5
)

var (
	watchCmd = &cobra.Command{
		Use:   "watch key",
		Short: "watch <key> (or keys with prefix)",
		RunE:  doWatch,
	}

	prefixFlag = false
)

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().BoolVarP(&prefixFlag, "prefix", "p", false, "watch all keys prefixed by key")
}

func doWatch(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return errors.New("argument required: <key>")
	}
	watchKey := args[0]
	if prefixFlag {
		log.Debug().Str("key-prefix", watchKey).Msg("watching prefix")
	} else {
		log.Debug().Str("key", watchKey).Msg("watching single key")
	}

	opts := []etcd.OpOption{}
	opts = append(opts, etcd.WithPrevKV())
	if prefixFlag {
		opts = append(opts, etcd.WithPrefix())
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	etcdClient, err := newEtcd3Client()
	if err != nil {
		log.Error().Err(err).Msg("failed to create etcd client")
		return err
	}

	watchCh := etcdClient.Watch(ctx, watchKey, opts...)

	for ctx.Err() == nil { // ever (until context canceled)
		var wr etcd.WatchResponse

		// read next event from etcd
		select {
		case wr = <-watchCh:
		case <-ctx.Done():
			return ctx.Err()
		}

		if len(wr.Events) <= 0 {
			log.Warn().Msg("received WatchResponse from etcd with ZERO events")
			continue
		}

		for ii, ev := range wr.Events {
			key := string(ev.Kv.Key)
			val := string(ev.Kv.Value)
			if ev.Type == etcd.EventTypePut {
				fmt.Printf("watch[%d]: PUT{k: %s, v: %s}\n", ii, key, val)
			} else if ev.Type == etcd.EventTypeDelete {
				fmt.Printf("watch[%d]: DEL{k: %s}\n", ii, key)
			} else {
				log.Warn().Msgf("UNKNOWN Event Type: %v", ev)
			}
		}
	}

	fmt.Println("watch stub...")
	return nil
}
