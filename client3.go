package main

import (
	"context"
	"flag"
	"log"

	example "github.com/rpcx-ecosystem/rpcx-examples3"
	"github.com/smallnest/rpcx/client"
)

var (
	addr13 = flag.String("addr1", "tcp@localhost:8972", "server1 address")
	addr23 = flag.String("addr2", "tcp@localhost:9981", "server2 address")
)

func main() {
	flag.Parse()

	d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr13}, {Key: *addr23}})
	option := client.DefaultOption
	option.Retries = 10
	xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, option)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 10; i++ {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		// if err != nil {
		// 	log.Printf("failed to call: %v", err)
		// } else {
		// 	log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		// }
		if err == nil {
			log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		} 
	}
}