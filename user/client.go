package main
import (
	"log"
	"context"
	"google.golang.org/grpc"
	pb "github.com/libgo/rbac/proto/user"
)

func main() {
	conn,err := grpc.Dial("localhost:9999",grpc.WithInsecure())
	if err!= nil {
		log.Printf("connect error: %v",err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	role_a := &pb.Role{Id:1,Name:"a"}
	resp,err := client.CreateUser(context.Background(),&pb.User{Id:1,Name:"izhaohu",Role:role_a})
	if err != nil {
		log.Printf("Create user failed: %v",err)
	}
	log.Printf("%v",resp.User)
}


