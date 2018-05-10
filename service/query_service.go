package service

import (
	"log"

	"golang.org/x/net/context"
	pb "github.com/magicgravity/chatserver/proto"
)


func (s *QueryServer)QueryUser(ctx context.Context,in *pb.QueryUserInfoRequest)(*pb.QueryUserInfoResponse,error){
	log.Printf("request content ===> %v \r\n",in)
	return &pb.QueryUserInfoResponse{},nil
}


func (s *QueryServer)GetSmsCode(ctx context.Context,in *pb.GetSmsCodeRequest)(*pb.GetSmsCodeResponse,error){
	log.Printf("request content ===> %v \r\n",in)
	return &pb.GetSmsCodeResponse{},nil
}
