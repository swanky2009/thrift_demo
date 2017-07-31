package main

import (
	"bufio"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"thrift_demo/rpc"
	"time"
)

const (
	HOST = "127.0.0.1"
	PORT = "9090"
)

func main() {

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	trans, err := thrift.NewTSocket(net.JoinHostPort(HOST, PORT))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}
	userTrans, err := transportFactory.GetTransport(trans)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving trans port:", err)
		os.Exit(1)
	}
	userService := rpc.NewUserServiceClientFactory(userTrans, protocolFactory)
	if err := userTrans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to "+HOST+":"+PORT, " ", err)
		os.Exit(1)
	}
	defer userTrans.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("input command:")
		comm, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		comm = strings.Replace(comm, "\r", "", -1)
		comm = strings.Replace(comm, "\n", "", -1)
		if comm == "" {
			continue
		}
		if comm == "quit" {
			break
		}
		ExecCommand(comm, userService)
	}
}

func ExecCommand(comm string, userService *rpc.UserServiceClient) {
	startTime := currentTimeMillis()

	comms := strings.Split(comm, " ")

	switch comms[0] {
	case "list":
		if len(comms) < 2 {
			fmt.Println("param error")
			return
		}
		page, _ := strconv.Atoi(comms[1])
		res, _ := userService.GetAllUsers(10, int32(page))
		fmt.Println("Rpc -> GetList:")
		for i, u := range res {
			if i == 0 {
				fmt.Println("| uid |", " name        |", " age |")
			}
			fmt.Println("| ", u.UID, "| ", u.Name, repeat(10-len(u.Name), ' '), "| ", u.Pro.Age, " | ")
		}
	case "add":
		if len(comms) < 3 {
			fmt.Println("param error")
			return
		}
		name := comms[1]
		age, _ := strconv.Atoi(comms[2])
		res, _ := userService.AddUser(name, int16(age))
		fmt.Println("Rpc -> AddUser:", res)
	case "update":
		if len(comms) < 4 {
			fmt.Println("param error")
			return
		}
		uid, _ := strconv.Atoi(comms[1])
		name := comms[2]
		age, _ := strconv.Atoi(comms[3])
		res, _ := userService.UpdateUser(int32(uid), name, int16(age))
		fmt.Println("Rpc -> UpdateUser:", res)
	case "get":
		if len(comms) < 2 {
			fmt.Println("param error")
			return
		}
		uid, _ := strconv.Atoi(comms[1])
		res, err := userService.GetUser(int32(uid))
		if err != nil {
			fmt.Println("Rpc -> GetUser:", err)
			break
		}
		fmt.Println("Rpc -> GetUser:", res)
	case "del":
		if len(comms) < 2 {
			fmt.Println("param error")
			return
		}
		uid, _ := strconv.Atoi(comms[1])
		res, _ := userService.DeleteUser(int32(uid))
		fmt.Println("Rpc -> DeleteUser:", res)
	default:
		fmt.Println("command error")
		return
	}

	endTime := currentTimeMillis()
	fmt.Printf("used time:%d-%d=%dms\n", endTime, startTime, (endTime - startTime))
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func repeat(time int, char rune) string {
	var s = make([]rune, time)
	for i := range s {
		s[i] = char
	}
	return string(s)
}
