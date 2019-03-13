package main
import (
	"net"
	"log"
	"context"
	"google.golang.org/grpc"
	pb "github.com/libgo/rbac/proto/access"
)

type AccessPool struct {
	permission []*pb.Permission
}

func (c *AccessPool) CreatePermission(p *pb.Permission) (*pb.Permission,error){
	c.permission = append(c.permission,p)
	return p,nil
}

type AccessServiceServer struct {
	accessPool AccessPool
}

func (c *AccessServiceServer) CreateAccessPermission(ctx context.Context,p *pb.Permission) (*pb.Response,error) {
	resp,err := c.accessPool.CreatePermission(p)
	if err != nil {
		log.Fatalf("Create permission failed:%v",err)
	}
	log.Fatalf("%v created",resp)
	return &pb.Response{Created:true},nil
}

func (c *AccessServiceServer) GetPermissionById(ctx context.Context,id *pb.Id) (*pb.Permission,error) {
	for i := 0; i < len(c.accessPool.permission); i++ {
		if id.Id == c.accessPool.permission[i].Id {
			return c.accessPool.permission[i],nil
		}
	}
	return &pb.Permission{},nil
}

const (
	port = ":9997"
)

func main(){
	Listen,err := net.Listen("tcp",port)

	if err != nil {
		log.Fatalf("net listen failed:%v",err)
	}

	server := grpc.NewServer()

	accessPool := AccessPool{}

	service := AccessServiceServer{accessPool}
	pb.RegisterAccessServiceServer(server,&service)
	server.Serve(Listen)
}


