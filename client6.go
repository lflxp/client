package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"10.6.200.8:2379"},
		DialTimeout: 5*time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), "/", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	fmt.Println("ok")
}