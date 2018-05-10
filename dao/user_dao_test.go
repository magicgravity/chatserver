package dao

import (
	"testing"
	"github.com/magicgravity/chatserver/pojo"
	"time"
	"github.com/magicgravity/chatserver/db"
)

func TestAddNewUser(t *testing.T) {
	defer db.CloseDb()
	u := new(pojo.User)
	u.Id = "4"
	u.UserId = "skyhigh"
	u.UserName = "西园寺世界"
	u.City = "beijing"
	u.Country = "China"
	u.Job = "都是世界的错"
	u.BgImgUrl ="www.baidu.com"
	u.Avatar = "www.taobao.com"
	u.Sex = 0
	u.Address = "太阳系地球中国"
	u.Introduce = "别砍我啊"
	u.Email = "332133@163.com"
	u.MobileNo ="13800138021"
	u.Password = "111111"
	u.CreateTime = time.Now()
	u.UpdateTime = time.Now()

	err := AddNewUser(u)
	if err!=nil{
		t.Fatal(err)
	}


}


func TestLoginVerifyByUidPwd(t *testing.T) {
	LoginVerifyByUidPwd("4","111111")
}