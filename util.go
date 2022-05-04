package main

import (
	"fmt"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

func FmtKv(kv *mvccpb.KeyValue) string {
	if kv == nil {
		return "NONE"
	}
	return fmt.Sprintf("([%d] %s: %s)",
		kv.Version,
		string(kv.Key),
		string(kv.Value),
	)
}
