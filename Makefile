build:
	protoc -I. --go_out=plugins=grpc:. proto/user/user.proto
	protoc -I. --go_out=plugins=grpc:. proto/role/role.proto
	protoc -I. --go_out=plugins=grpc:. proto/access/access.proto
