package main
import (
	"log"
	"net"
	"fmt"
	"context"
	"google.golang.org/grpc"
	pb "github.com/libgo/rbac/proto/role"
)

type RolePool struct {
	role []*pb.Role
}

func (c *RolePool) CreateRole(role *pb.Role) (*pb.Role,error) {
	c.role = append(c.role,role)
	return role,nil
}

type RoleServiceServer struct {
	rolePool RolePool
}

func (c *RoleServiceServer) CreateRole(ctx context.Context,role *pb.Role) (*pb.Response,error){
	role,error := c.rolePool.CreateRole(role)
	if error != nil {
		log.Fatalf("create role error:%v",error)
	}
	resp := pb.Response{Created:true}
	return &resp,nil
}

func (c *RoleServiceServer) GetAccessPermission(ctx context.Context,role *pb.Role) (*pb.Permission,error) {
	return role.GetPermission(),nil
}

func (c *RoleServiceServer) GetRole(ctx context.Context,id *pb.Id)(*pb.Role,error){
	for i :=0; i<len(c.rolePool.role); i++ {
		if id.Id == c.rolePool.role[i].Id {
			return c.rolePool.role[i],nil
		}
	}
	return &pb.Role{},nil
}

const (
	port = ":9998"
)
func main(){
	Listen,err := net.Listen("tcp",port)
	if err != nil {
		log.Fatalf("net listen error:%v",err)
	}
	server := grpc.NewServer()

	rolePool := RolePool{}
	service := RoleServiceServer{rolePool}
	fmt.Println("Start listen on :",port)

	pb.RegisterRoleServiceServer(server,&service)
	server.Serve(Listen)
}
