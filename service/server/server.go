package main

import (
	"log"
	micro "realmicrokube/micro"
	pb "realmicrokube/service/proto"
	// db "realmicrokube/service/sddb"
	db "realmicrokube/service/db"

	"golang.org/x/net/context"
)

type userServer struct {
	db *db.DB
}

func newServer() pb.SdUserServer {
	return &userServer{
		db: db.NewDB(),
	}
}

func newService() {
	config := &micro.ServiceConfig{
		Name: "com.shendu.service.usercenter.user",
		Port: 9999,
	}
	micro.NewService(config, newServer(), pb.RegisterSdUserServer)
}

func (s *userServer) GetUserInfo(ctx context.Context, in *pb.UserReq) (*pb.UserResp, error) {
	log.Println("Client request is arraving...")
	userCount, err := s.db.QueryUserCount()
	if err != nil {
		log.Println("Query user count error => ", err.Error())
	}
	log.Println("The user count is => ", userCount)
	return &pb.UserResp{User: &pb.User{Id: 123456, Name: "Xiaoming"}}, nil
}

func main() {
	newService()
}
