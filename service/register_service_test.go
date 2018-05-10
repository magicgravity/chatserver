package service


import (
	"testing"
	"log"
	"time"
	"google.golang.org/grpc"
	pb "github.com/magicgravity/chatserver/proto"
	"golang.org/x/net/context"
	"github.com/magicgravity/chatserver/common"
)

func TestRegistServer_IsExist(t *testing.T) {
	isOk := make(chan int)
	isEnd := make(chan int)
	var rs common.Server = new(RegistServer)
	rs.SetPort(":50050")


	go func() {
		log.Printf("client wait to connect")
		<-isOk
		log.Printf("client begin to connect")
		// Set up a connection to the server.
		conn, err := grpc.Dial("127.0.0.1:50050", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		log.Printf("client  connect ok")

		defer conn.Close()
		c := pb.NewRegistClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		testReq := pb.CheckUserExistRequest{}
		testReq.MobileNo = "13800138000"
		testReq.UserId = "xc01"
		commonHeader := pb.CommonHeaderRequest{}
		commonHeader.IpAddr = "192.156.33.28"
		commonHeader.ReqTime = common.FormatCurrentDateYYYYMMdd()
		commonHeader.DeviceType = common.DeviceType_PC
		commonHeader.VersionNo = common.CurVersionNo
		commonHeader.TransCode = common.TransCode_CheckUserExist
		commonHeader.ReqSeqNo = 1
		testReq.CommonHeader = &commonHeader
		defer func() {
			isEnd<-1
		}()
		r, err := c.IsExist(ctx, &testReq)
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		log.Printf("response content === > : %v", r)

	}()


	server,err := StartServer(rs)
	if err!= nil{
		panic(err)
	}else{
		log.Printf("server start ok!")
		isOk<-1
		for{
			select{
				case <-isEnd:
					server.Stop()
					log.Printf("now server is stoping ...")
				default:
					log.Printf("server is runing...")

			}
		}
	}

	//for{
	//	select {
	//	case <-isStop:
	//		return
	//	default:
	//		time.Sleep(time.Second)
	//	}
	//}
	//


}


func TestRegistServer_Register(t *testing.T) {

}

