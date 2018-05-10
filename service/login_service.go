package service

import (
	"log"

	"golang.org/x/net/context"
	pb "github.com/magicgravity/chatserver/proto"
	"github.com/magicgravity/chatserver/common"
)


func (s *LoginServer) VerifyLogin(ctx context.Context, in *pb.VerifyUserLoginByUidPwdRequest) (*pb.VerifyUserLoginByUidPwdResponse, error) {
	log.Printf("request content ===> %v \r\n",in)
	reqHeader := in.CommonHeader


	resp := pb.VerifyUserLoginByUidPwdResponse{}
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

	resp.VerifyResultFlag = 0
	resp.User= nil

	return &resp,nil
}

func (s *LoginServer) VerifyLoginBySms (ctx context.Context,in *pb.VerifyUserLoginByMobileSmsRequest)(*pb.VerifyUserLoginByMobileSmsResponse,error){
	log.Printf("request content ===> %v \r\n",in)
	reqHeader := in.CommonHeader


	resp := pb.VerifyUserLoginByMobileSmsResponse{}
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

	resp.VerifyResultFlag = 0
	resp.User= nil

	return &resp,nil
}