package discovery

import (
	"context"
	"testing"
	"time"
)

func TestServiceDiscovery(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ctx := context.Background()
	ser := NewServiceDiscovery(&ctx, endpoints)
	defer ser.Close()
	err := ser.WatchService("/web/", func(key, value string) {}, func(key, value string) {})
	if err != nil {
		return
	}
	err = ser.WatchService("/gRPC/", func(key, value string) {}, func(key, value string) {})
	if err != nil {
		return
	}
	for {
		select {
		case <-time.Tick(10 * time.Second):
		}
	}
}
