package service


import (
	"testing"
	"log"
	//"time"
	//"google.golang.org/grpc"
	//pb "github.com/magicgravity/chatserver/proto"
	//"golang.org/x/net/context"
	"github.com/magicgravity/chatserver/common"
	"time"
)

func TestRegistServer_IsExist2(t *testing.T) {
	var rs common.Server = new(RegistServer)
	rs.SetPort(":50050")




	server,err := StartServer(rs)
	if err!= nil{
		panic(err)
	}else{
		log.Printf("server start ok!")
		outpoint:
		for{
			select {
				case <-time.After(time.Minute):
					server.Stop()
					log.Printf("now server is stoping ...")
					break outpoint
			}
		}
		log.Printf("end test")
	}



}
