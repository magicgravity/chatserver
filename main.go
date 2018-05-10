package main

import (
	"reflect"
	"log"
	"github.com/magicgravity/chatserver/service"
	"github.com/magicgravity/chatserver/common"
	//"time"
	//"google.golang.org/grpc"
	//pb "github.com/magicgravity/chatserver/proto"
	//"golang.org/x/net/context"
)



func main(){



	var rs common.Server = new(service.RegistServer)
	serverType := reflect.ValueOf(rs)
	log.Printf("serverType ==> %s ",serverType.Elem().String())
	//<service.RegistServer Value>

	serverType2 := reflect.ValueOf(&rs)
	log.Printf("serverType2 ==> %s ",serverType2.Elem().String())
	//<service.RegistServer Value>


	var rs2 common.Server = new(service.RegistServer)
	rs2.SetPort(":50050")

	service.StartServer(rs2)

	//go func(){
	//
	//	server,err := service.StartServer(rs2)
	//	if err!= nil{
	//		panic(err)
	//	}else{
	//		log.Printf("server start ok!")
	//	outpoint:
	//		for{
	//			select {
	//			case <-time.After(time.Minute):
	//				server.Stop()
	//				log.Printf("now server is stoping ...")
	//				break outpoint
	//			}
	//		}
	//		log.Printf("end test")
	//	}
	//}()




	//
	//log.Printf("client wait to connect")
	//
	//log.Printf("client begin to connect")
	//// Set up a connection to the server.
	//conn, err := grpc.Dial("127.0.0.1:50050", grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//log.Printf("client  connect ok")
	//
	//defer conn.Close()
	//c := pb.NewRegistClient(conn)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//
	//testReq := pb.CheckUserExistRequest{}
	//testReq.MobileNo = "13800138000"
	//testReq.UserId = "xc01"
	//commonHeader := pb.CommonHeaderRequest{}
	//commonHeader.IpAddr = "192.156.33.28"
	//commonHeader.ReqTime = common.FormatCurrentDateYYYYMMdd()
	//commonHeader.DeviceType = common.DeviceType_PC
	//commonHeader.VersionNo = common.CurVersionNo
	//commonHeader.TransCode = common.TransCode_CheckUserExist
	//commonHeader.ReqSeqNo = 1
	//testReq.CommonHeader = &commonHeader
	//r, err := c.IsExist(ctx, &testReq)
	//if err != nil {
	//	log.Fatalf("could not connect: %v", err)
	//}
	//log.Printf("response header === > : %v  , response flag ====> %v", r.CommonHeader,r.ExistFlag)

}



