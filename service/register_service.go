package service


import (
	"log"

	"golang.org/x/net/context"
	pb "github.com/magicgravity/chatserver/proto"
	"github.com/magicgravity/chatserver/common"
)


func (s *RegistServer)Register(ctx context.Context,in *pb.UserRegisterRequest)(*pb.UserRegisterResponse,error){
	log.Printf("request content ===> %v \r\n",in)
	reqHeader := in.CommonHeader


	resp := pb.UserRegisterResponse{}
	respHeader := new(pb.CommonHeaderResponse)
	respHeader.ReqSeqNo = reqHeader.ReqSeqNo
	respHeader.TransCode = reqHeader.TransCode
	respHeader.VersionNo = common.CurVersionNo
	respHeader.MsgCode = common.MsgCode_SUCCESS
	respHeader.MsgInfo = "success"
	respHeader.ResTime = common.FormatCurrentDateYYYYMMdd()
	//TODO now for test
	respHeader.SessionId = "1111"
	resp.CommonHeader = respHeader

	resp.RegisterFlag = true

	return &resp,nil
}


func (s *RegistServer)IsExist(ctx context.Context,in *pb.CheckUserExistRequest)(*pb.CheckUserExistResponse,error){
	log.Printf("request content ===> %v \r\n",in)
	reqHeader := in.CommonHeader


	resp := pb.CheckUserExistResponse{}
	respHeader := new(pb.CommonHeaderResponse)
	respHeader.ReqSeqNo = reqHeader.ReqSeqNo
	respHeader.TransCode = reqHeader.TransCode
	respHeader.VersionNo = common.CurVersionNo
	respHeader.MsgCode = common.MsgCode_SUCCESS
	respHeader.MsgInfo = "success"
	respHeader.ResTime = common.FormatCurrentDateYYYYMMdd()
	//TODO now for test
	respHeader.SessionId = "1111"
	resp.CommonHeader = respHeader

	resp.ExistFlag = false

	return &resp,nil
}


