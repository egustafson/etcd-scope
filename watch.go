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

	// mimic the caller of this func's context
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	etcdClient, err := newEtcd3Client()
	if err != nil {
		log.Error().Err(err).Msg("failed to create etcd client")
		return err
	}

	opts := []etcd.OpOption{}
	opts = append(opts, etcd.WithPrevKV())
	if prefixFlag {
		opts = append(opts, etcd.WithPrefix())
	}

	//watchCh := etcdClient.Watch(ctx, watchKey, opts...)

	var lastRev int64 = 0
	for ctx.Err() == nil { // ever (until context canceled)
		lastRev, err = watchUntilDisconnect(ctx, etcdClient, lastRev, watchKey, opts)
		log.Info().Int64("last-rev", lastRev).Msg("broke out of watchUntilDisconnect()")
	}

	return ctx.Err()
}

func watchUntilDisconnect(ctx context.Context, etcdClient *etcd.Client, rev int64, key string, opts []etcd.OpOption) (lastRev int64, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	lastRev = rev
	if lastRev > 0 { // pick up where we left off
		log.Info().Int64("last-rev", lastRev+1).Msg("setting OpOption.WithRev()")
		opts = append(opts, etcd.WithRev(lastRev+1))
	}

	watchCh := etcdClient.Watch(ctx, key, opts...)

	for ctx.Err() == nil { // ever (until context canceled)
		var wr etcd.WatchResponse

		// read next event from etcd
		select {
		case wr = <-watchCh:
		case <-ctx.Done():
			return lastRev, ctx.Err()
		}

		if wr.Created {
			log.Info().Msg("WatchResponse.Created = TRUE")
		}

		if wr.Canceled {
			log.Warn().Int64("rev", wr.Header.GetRevision()).Msg("received WatchResponse.Canceled = TRUE")
			return lastRev, nil

		}

		if len(wr.Events) <= 0 {
			log.Warn().Msg("received WatchResponse from etcd with ZERO events")
			continue
		}

		for ii, ev := range wr.Events {
			if ev.Type == etcd.EventTypePut {
				fmt.Printf("watch[%d](rev: %d): PUT: %s <-- %s\n",
					ii, wr.Header.Revision, FmtKv(ev.Kv), FmtKv(ev.PrevKv))
			} else if ev.Type == etcd.EventTypeDelete {
				fmt.Printf("watch[%d](rev: %d): DEL <-- %s\n",
					ii, wr.Header.Revision, FmtKv(ev.PrevKv))
			} else {
				log.Warn().Msgf("UNKNOWN Event Type: %v", ev)
			}
		}
		lastRev = wr.Header.Revision
	}
	return lastRev, ctx.Err()
}
