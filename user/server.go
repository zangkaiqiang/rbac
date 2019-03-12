package main
import (
	"log"
	"context"
	"google.golang.org/grpc"
	"net"
	pb "github.com/libgo/rbac/proto/user"
)

//implement UserServiceServer
type UserServiceServer struct{
	user_pool UserPool
}

//Create UserPool for store user
type UserPool struct {
	users []*pb.User
}

//add createuser to userpool
func (c *UserPool) CreateUser(user *pb.User)(*pb.User,error){
	c.users = append(c.users,user)
	return user,nil
}

//implement createuser for userserviceserver
func (c *UserServiceServer) CreateUser(ctx context.Context,u *pb.User) (*pb.Response,error){
	user,err := c.user_pool.CreateUser(u)
	if err != nil {
		log.Printf("Create user error: %v",err)
	}

	return &pb.Response{Created:true,User:user},nil
}

//define port
const (
	port = ":9999"
)

//implement main func
func main(){
	Listen,err := net.Listen("tcp",port)
	if err!= nil {
		log.Printf("net listen error: %v",err)
	}
	log.Printf("Start listening port:%v",port)
	server := grpc.NewServer()
	user_pool := UserPool{}
	service := UserServiceServer{user_pool:user_pool}
	pb.RegisterUserServiceServer(server,&service)
	server.Serve(Listen)

}


