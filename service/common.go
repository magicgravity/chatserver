package service

import (
	"net"
	"log"
	"reflect"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	pb "github.com/magicgravity/chatserver/proto"
	"github.com/magicgravity/chatserver/common"
	"errors"
)

func StartServer(bs common.Server)(*grpc.Server,error){
	lis, err := net.Listen("tcp", bs.GetPort())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	serverType := reflect.ValueOf(bs)
	log.Printf("serverType ==> %s \r\n",serverType.Elem().String())
	switch serverType.Elem().String(){
		case "<service.LoginServer Value>":
			pb.RegisterLoginServer(s,&LoginServer{})
		case "<service.RegistServer Value>":
			pb.RegisterRegistServer(s,&RegistServer{})
		case "<service.ChatServer Value>":
			pb.RegisterChatServer(s,&ChatServer{})
		case "<service.QueryServer Value>":
			pb.RegisterQueryServer(s,&QueryServer{})
		default:
			return nil,errors.New("start server fail ")
	}

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return s,err
	}else{
		bs.Start()
		log.Printf("Server [%s] Start Ok ...\r\n",serverType.String())
		return s,nil
	}
}


type LoginServer struct{
	common.BaseServer
}

type RegistServer struct {
	common.BaseServer
}

type ChatServer struct {
	common.BaseServer
}

type QueryServer struct {
	common.BaseServer
}
