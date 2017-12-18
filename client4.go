package main

import (
	"context"
	"flag"
	"log"

	"github.com/smallnest/rpcx/share"
	example "github.com/rpcx-ecosystem/rpcx-examples3"
	"github.com/smallnest/rpcx/client"
)

var (
	addr134 = flag.String("addr1", "tcp@10.6.200.8:8972", "server1 address")
	addr234 = flag.String("addr2", "tcp@localhost:9981", "server2 address")
)

func main() {
	flag.Parse()

	d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr134}, {Key: *addr234}})
	option := client.DefaultOption
	option.Retries = 10
	xclient := client.NewXClient("Arith", client.Failover, client.RandomSelect, d, option)
	defer xclient.Close()

	xclient.Auth("bearer tGzv3JOkF0XG5Qx2TlKWIA")

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 1000; i++ {
		reply := &example.Reply{}
		ctx := context.WithValue(context.Background(),share.ReqMetaDataKey,make(map[string]string))
		err := xclient.Call(ctx, "Mul", args, reply)
		// if err != nil {
		// 	log.Printf("failed to call: %v", err)
		// } else {
		// 	log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		// }
		if err == nil {
			log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		} else {
			log.Fatalf("failed to call: %v",err)
		}
	}
}