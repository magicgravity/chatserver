package dao

import (
	"log"
	"github.com/magicgravity/chatserver/pojo"
	"github.com/magicgravity/chatserver/db"
)

func AddNewUser(u *pojo.User)error{
	log.Printf("prepare to add user >> %v \r\n",u)
	session := db.GetConn().NewSession()
	defer session.Close()

	/*
	insert into user (id,username,userid,email,mobileno,address,sex,introduce,avatar,bgimgurl,job,city,country,createtime,updatetime,password)
        values (?1,?2,?3,?4,?5,?6,?7,?8,?9,?10,?11,?12,?13,?14,?15,?16)
	 */

	paramMap := map[string]interface{}{"1": u.Id, "2": u.UserName,
										"3":u.UserId,"4":u.Email,"5":u.MobileNo,
										"6":u.Address,"7":u.Sex,"8":u.Introduce,
										"9":u.Avatar,"10":u.BgImgUrl,"11":u.Job,
										"12":u.City,"13":u.Country,"14":u.CreateTime,
										"15":u.UpdateTime,"16":u.Password}
	_,err := session.SqlMapClient("addNewUser", &paramMap).Execute()
	if err!= nil{
		log.Fatalf("add new user fail >> %v",err)
		return err
	}else{
		return nil
	}
}


func LoginVerifyByUidPwd(uid ,pwd string) (bool,error){
	log.Printf("prepare to execute login verify by uid[%s] and pwd[%s] \r\n",uid,pwd)
	session := db.GetConn().NewSession()
	defer session.Close()

	paramMap :=  map[string]interface{}{
		"1":uid,"2":pwd,
	}
	var users []pojo.User
	result:=session.SqlMapClient("queryUserByUidPwd",&paramMap).Search(&users)
	if result.Error!= nil{
		log.Fatalf("login verify user fail >> %v",result.Error)
		return false,result.Error
	}else{
		log.Printf("result ==> %v",result.Result)
		for _,u :=range users {
			log.Printf("user==> %v", u)
		}
		return true,nil
	}

}
