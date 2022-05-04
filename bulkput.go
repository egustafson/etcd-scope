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
	bulkputCmd = &cobra.Command{
		Use:   "bulkput key-prefix",
		Short: "put a whole bunch of keys",
		RunE:  doBulkput,
	}

	countFlag = 100
)

func init() {
	rootCmd.AddCommand(bulkputCmd)
	bulkputCmd.Flags().IntVarP(&countFlag, "count", "n", 100,
		"number of key/values to put")
}

func doBulkput(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return errors.New("argument required: <key-prefix>")
	}
	bulkKeyPrefix := args[0]

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	etcdClient, err := newEtcd3Client()
	if err != nil {
		log.Error().Err(err).Msg("failed to create etcd client")
		return err
	}

	opts := []etcd.OpOption{}
	// no opts at this time

	startTime := time.Now()
	for ii := 0; ii < countFlag; ii++ {
		key := fmt.Sprintf("%s/k-%d", bulkKeyPrefix, ii)
		val := fmt.Sprintf("v-%09d", ii)
		_, err := etcdClient.Put(ctx, key, val, opts...)
		if err != nil {
			log.Warn().Err(err).Str("key", key).Msg("Put error")
		}
	}
	duration := time.Since(startTime)
	fmt.Printf("Took %s to PUT %v k/v's\n", duration, countFlag)
	return nil
}
