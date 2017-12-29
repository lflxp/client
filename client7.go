package main

import (
	"fmt"
	"log"
	"sync"
	"context"
	"strings"
	"github.com/lflxp/dbui/etcd"
	// "github.com/smallnest/rpcx/share"
	"github.com/smallnest/rpcx/client"
	"github.com/coreos/etcd/clientv3"
)

var etcdDD string = "10.23.70.80:2379,10.123.4.46:2379,10.123.4.38:2379"
var Clientxx []*client.KVPair
var Test map[string]string

func init() {
	Test = map[string]string{}
	Clientxx = []*client.KVPair{}
	st := &etcd.EtcdUi{Endpoints: strings.Split(etcdDD,",")}
	st.InitClientConn()

	//获取客户端
	resp := st.More("/ams/main/services/install_ping")
	for _,v := range resp.Kvs {
		if strings.Contains(string(v.Key),"tcp@") {
			log.Println(string(v.Key))
			tmp := strings.Split(string(v.Key),"/")
			log.Println("Old",tmp[len(tmp)-1])
			Clientxx = append(Clientxx,&client.KVPair{Key:tmp[len(tmp)-1]})
			Test[tmp[len(tmp)-1]] = string(v.Value)
		}
	}

	rch := st.ClientConn.Watch(context.Background(),"/ams/main/services/install_ping",clientv3.WithPrefix())
	go func(ccc clientv3.WatchChan) {
		for wresp := range rch {
			for _,ev := range wresp.Events {
				log.Println("机型配置更改")
				if strings.Contains(string(ev.Kv.Key),"tcp@") {
					log.Println("before",Test)
					log.Println("###########",string(ev.Kv.Key),string(ev.Kv.Value))
					tmp := strings.Split(string(ev.Kv.Key),"/")
					switch ev.Type.String() {
					case "PUT":
						log.Println("PUT",tmp[len(tmp)-1],string(ev.Kv.Value))
						Test[tmp[len(tmp)-1]] = string(ev.Kv.Value)
					case "DELETE":
						delete(Test,tmp[len(tmp)-1])
					default:
						log.Println(fmt.Sprintf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value))
					}
					log.Println("after",Test)
				}
			}
			Tclient := []*client.KVPair{} 
			for key,_ := range Test {
				log.Println("New",key)
				Tclient = append(Tclient,&client.KVPair{Key:key})
			}
			Clientxx = Tclient	
			log.Println(Test)
		}
	}(rch)
}

type Config struct {
	User 		string
	Password 	string
	Type 		string
	Sn 			string
	Url 		string
	ServerApi   string
}

type Reply struct {
	PXE 		string
	REBOOT  	string
	Ping		bool
	Err 		error
}

var waitgroup sync.WaitGroup

func main() { 
	chans := make(chan int)
	for {
		for i:=0;i<100;i++ {	
			waitgroup.Add(1)
			go func(){
				d := client.NewMultipleServersDiscovery(Clientxx)
				option := client.DefaultOption
				// option.Retries = 10
				// option.ConnectTimeout = 5
				// option.ReadTimeout = 5
				// option.WriteTimeout = 5
				xclient := client.NewXClient("Config",client.Failfast,client.RandomSelect,d,option)
				// defer xclient.Close()
				args := Config{
					User:"users",
					Password:"passwords",
					Type:"R720XD",
					Sn:"G4MTGY1",
					ServerApi:"http://portal.qiyi.domain",
				}
		
			
				reply := &Reply{}
				
				call,err := xclient.Go(context.Background(),"Run",args,reply,nil)
				if err != nil {
					log.Println("error:",err.Error())
				}
			
				replyCall := <- call.Done
				if replyCall.Error != nil {
					log.Fatalf("failed to call: %v",replyCall.Error)
				} else {
					if reply.Err != nil {
						log.Printf(reply.Err.Error())
					} else {
						log.Printf("%t %s %s",reply.Ping,reply.PXE,reply.REBOOT)
					}
				} 
		
				xclient.Close()
				waitgroup.Done()
			}()
	
			// reply := &Reply{}
			// ctx := context.WithValue(context.Background(),share.ReqMetaDataKey,make(map[string]string))	
			// err := xclient.Call(ctx,"Run",args,reply)
			// if err != nil {
			// 	log.Fatal("error:",err.Error())
			// }
		
			// log.Printf("%t %s %s",reply.Ping,reply.PXE,reply.REBOOT)
		}
		waitgroup.Wait()
	}
	

	<-chans
	// for i:=0;i<1000;i++ {
	// 	reply := &Reply{}
	// 	ctx := context.WithValue(context.Background(),share.ReqMetaDataKey,make(map[string]string))	
	// 	err := xclient.Call(ctx,"Run",args,reply)
	// 	if err != nil {
	// 		log.Fatal("error:",err.Error())
	// 	}
	
	// 	log.Printf("%t %s %s",reply.Ping,reply.PXE,reply.REBOOT)
	// }
}