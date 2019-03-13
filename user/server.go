package main
import (
	"log"
	"context"
	"google.golang.org/grpc"
	"net"
	pb "github.com/libgo/rbac/proto/user"
	pbrole "github.com/libgo/rbac/proto/role"
	pbaccess "github.com/libgo/rbac/proto/access"
)

//implement UserServiceServer
type UserServiceServer struct{
	userPool UserPool
	//roleclient by roleservice
	roleClient pbrole.RoleServiceClient
	//accessclient by accessservice
	accessClient pbaccess.AccessServiceClient
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
	user,err := c.userPool.CreateUser(u)
	if err != nil {
		log.Printf("Create user error: %v",err)
	}

	return &pb.Response{Created:true,User:user},nil
}


func (c *UserServiceServer) GetRole(ctx context.Context,user *pb.User) (*pb.Role,error) {
	return user.Role,nil
}

//get permission by user
func (c *UserServiceServer) GetAccessPermission(ctx context.Context,user *pb.User) (*pb.Permission,error) {
	user_role := user.Role
	//get role from roleid from roleservice
	role,err := c.roleClient.GetRole(context.Background(), &pbrole.Id{Id:user_role.Id})
	if err != nil {
		log.Fatalf("GetRole error:%v",err)
	}
	//get permissionid from roleservice by role
	permission,err := c.roleClient.GetAccessPermission(context.Background(),role)
	if err != nil {
		log.Fatalf("Get permission error:%v",err)
	}
	accessid := pbaccess.Id{Id:permission.Id}
	//get permission from permissionservice by permissionid
	access,err := c.accessClient.GetPermissionById(context.Background(),&accessid)

	if err != nil {
		log.Fatalf("Get Permission error:%v",err)
	}
	return &pb.Permission{Id:access.Id,Name:access.Name},nil
}

//define port
const (
	port = ":9999"
	roleAddress = "localhost:9998"
	accessAddress = "localhost:9997"
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
	//init roleclient
	roleConn,err := grpc.Dial(roleAddress,grpc.WithInsecure())
	if err != nil {
		log.Fatalf("role service connect error:%v",err)
	}
	defer roleConn.Close()
	roleClient := pbrole.NewRoleServiceClient(roleConn)
	//init accessclient
	accessConn,err := grpc.Dial(accessAddress,grpc.WithInsecure())
	if err != nil {
		log.Fatalf("access service connect error:%v",err)
	}
	defer accessConn.Close()
	accessClient := pbaccess.NewAccessServiceClient(accessConn)

	service := UserServiceServer{userPool:user_pool,roleClient:roleClient,accessClient:accessClient}
	pb.RegisterUserServiceServer(server,&service)
	server.Serve(Listen)
}


