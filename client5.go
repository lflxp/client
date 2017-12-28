package main

import (
	"context"
	"flag"
	"log"

	"github.com/smallnest/rpcx/client"
)

var (
	addr135 = flag.String("addr1", "tcp@10.6.200.8:8972", "server1 address")
	// addr23 = flag.String("addr2", "tcp@localhost:9981", "server2 address")
)

func main() {
	flag.Parse()

	// d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr13}, {Key: *addr23}})
	d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr135},{Key:"tcp@10.6.200.8:8973"},{Key:"tcp@10.6.200.8:8974"},{Key:"tcp@10.6.200.8:8975"}})
	option := client.DefaultOption
	option.Retries = 10
	xclient := client.NewXClient("PxeTemplate", client.Failover, client.RandomSelect, d, option)
	defer xclient.Close()

	hehe := map[string]string{
		"$SERVER_IP":  "10.10.5.9",
		"$OS_VERSION": "vxlan_centos7.2",
		"$SN":         "G4JA923FAS",
		"$IPADDR":     "10.23.72.52",
		"$NETMASK":    "255.255.255.0",
		"$GATEWAY":    "10.23.72.254",
	}

	type PxeReply struct {
		Result 		string
		Status 		bool
	}

	for i := 0; i < 1000; i++ {
		reply := &PxeReply{}
		call,err := xclient.Go(context.Background(), "AutoTemplate", hehe, reply,nil)
		// if err != nil {
		// 	log.Printf("failed to call: %v", err)
		// } else {
		// 	log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		// }
		if err != nil {
			log.Printf("error : %s",err.Error())
		}

		replyCall := <- call.Done
		if replyCall.Error != nil {
			log.Fatalf("failed to call: %v",replyCall.Error)
		} else {
			log.Printf("sn : %s %d", reply.Result,i)
		}
		
	}
}