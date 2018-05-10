package service


import (
	"log"

	"golang.org/x/net/context"
	pb "github.com/magicgravity/chatserver/proto"
)

func (s *ChatServer)SearchFriend(ctx context.Context,in *pb.SearchOtherPersonsRequest)(*pb.SearchOtherPersonsResponse,error){
	log.Printf("request content ===> %v \r\n",in)
	return &pb.SearchOtherPersonsResponse{},nil
}


func (s *ChatServer)AddFriend(ctx context.Context,in *pb.SendMakeFriendRequest)(*pb.SendMakeFriendResponse,error){
	log.Printf("request content ===> %v \r\n",in)
	return &pb.SendMakeFriendResponse{},nil
}

func (s *ChatServer)SendMsg(ss pb.Chat_SendMsgServer)error{
	log.Printf("sendMsg  ===>  \r\n")
	return nil
}
