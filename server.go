package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"log"
	"os"
	"thrift_demo/models"
	"thrift_demo/rpc" //导入Thrift生成的接口包
	//"time"
)

const (
	NetworkAddr = "127.0.0.1:9090" //监听地址&端口
)

type userService struct {
}

func (this *userService) AddUser(name string, age int16) (r int64, err error) {
	Log("client Call --> AddUser:", name, "\t", age)
	r = models.AddUser(name, age)
	return
}

func (this *userService) UpdateUser(uid int32, name string, age int16) (r *rpc.User, err error) {
	Log("client Call --> UpdateUser:", uid, "\t", name, "\t", age)
	u, er := models.UpdateUser(int(uid), name, age)
	r = &rpc.User{UID: int32(u.Uid), Name: u.Name, Pro: &rpc.Profile{UID: int32(u.Profile.Uid), Age: int32(u.Profile.Age)}}
	err = er
	return
}

func (this *userService) DeleteUser(uid int32) (r int64, err error) {
	Log("client Call --> DeleteUser:", uid)
	r = models.DeleteUser(int(uid))
	return
}

func (this *userService) GetUser(uid int32) (r *rpc.User, err error) {
	Log("client Call --> GetUser:", uid)
	u, er := models.GetUser(int(uid))
	if er != nil {
		err = er
		return
	}
	r = &rpc.User{UID: int32(u.Uid), Name: u.Name, Pro: &rpc.Profile{UID: int32(u.Profile.Uid), Age: int32(u.Profile.Age)}}
	err = er
	return
}

func (this *userService) GetAllUsers(rows int32, page int32) (r []*rpc.User, err error) {
	Log("client Call --> GetAllUsers:rows:", rows, "\tpage:", page)
	users := models.GetAllUsers(int(rows), int(page))
	for _, u := range users {
		r = append(r, &rpc.User{UID: int32(u.Uid), Name: u.Name, Pro: &rpc.Profile{UID: int32(u.Profile.Uid), Age: int32(u.Profile.Age)}})
	}
	return
}

func main() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//protocolFactory := thrift.NewTCompactProtocolFactory()

	serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	handler := &userService{}
	processor := rpc.NewUserServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("thrift server in", NetworkAddr)
	server.Serve()
}

func Log(v ...interface{}) {
	log.Println(v...)
}
