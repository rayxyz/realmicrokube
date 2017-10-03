package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	pb "realmicrokube/service/proto"

	micro "realmicrokube/micro"
)

func call(w http.ResponseWriter, r *http.Request) {
	client, err := micro.NewServiceClient("com-shendu-service-usercenter-user", pb.NewSdUserClient)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Successfully created service client => ", client.Config.Name)
	log.Println("New service client succeed. service host => ", client.Config.Host, " port => ", client.Config.Port)
	resp, err := client.Call("GetUserInfo", context.TODO(), &pb.UserReq{Id: 123456})
	if err != nil {
		log.Println(err)
		return
	}
	uinfo, err := json.Marshal(resp.(*pb.UserResp).User)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(uinfo))
	w.Write(uinfo)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am the hi API."))
	})
	http.HandleFunc("/hiapi/call", func(w http.ResponseWriter, r *http.Request) {
		log.Println("In call...")
		call(w, r)
	})
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		call(w, r)
	})
	log.Println("Hi api server running on port: 6767")
	http.ListenAndServe(":8989", nil)
}
