package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	etcd "go.etcd.io/etcd/client/v3" // <-- v3.5
)

var (
	beatCmd = &cobra.Command{
		Use:   "beat key [value]",
		Short: "periodically put a <value> to <key>",
		RunE:  doBeat,
	}

	periodFlag = time.Second
)

func init() {
	rootCmd.AddCommand(beatCmd)
	beatCmd.Flags().DurationVarP(&periodFlag, "period", "p", time.Second,
		"period of time between 'beats'")
}

func doBeat(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return errors.New("argument required: <key>")
	}
	beatKey := args[0]
	beatVal := "" // default, "" ==> timestamp
	if len(args) > 1 {
		beatVal = args[1] // if specified, arg 2 is the value to write
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	etcdClient, err := newEtcd3Client()
	if err != nil {
		log.Error().Err(err).Msg("failed to create etcd client")
		return err
	}

	opts := []etcd.OpOption{}
	opts = append(opts, etcd.WithPrevKV())

	beatValIsTime := (len(beatVal) < 1)
	for ctx.Err() == nil { // ever (until context canceled)
		if beatValIsTime {
			beatVal = time.Now().Format(time.RFC3339)
		}

		resp, err := etcdClient.Put(ctx, beatKey, beatVal, opts...)
		if err != nil {
			log.Warn().Err(err).Str("key", beatKey).Msg("Put error")
		} else {
			fmt.Printf("PUT:(rev: %d): (%s: %s) <-- %s\n",
				resp.Header.Revision, beatKey, beatVal, FmtKv(resp.PrevKv))
		}
		time.Sleep(periodFlag)
	}

	return nil
}
