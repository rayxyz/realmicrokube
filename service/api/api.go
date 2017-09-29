package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	pb "realmicrokube/service/proto"

	micro "realmicrokube/micro"
)

func call(w http.ResponseWriter, r *http.Request) {
	client, err := micro.NewServiceClient("com.shendu.service.sduser.user", pb.NewSdUserClient)
	if err != nil {
		log.Println(err.Error())
		return
	}
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
	// call()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am the hi API."))
	})
	http.HandleFunc("/hiapi/call", func(w http.ResponseWriter, r *http.Request) {
		log.Println("In call...")
		call(w, r)
	})
	log.Println("Hi api server running on port: 6767")
	http.ListenAndServe(":8989", nil)
}
