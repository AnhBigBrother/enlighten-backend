package interceptor

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func StreamInterceptor(
	srv any,
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Println("-----")
	log.Println("stream interceptor:", info.FullMethod)
	md, _ := metadata.FromIncomingContext(stream.Context())
	log.Println("authorization", md["authorization"])
	log.Println("description", md["description"])

	return handler(srv, stream)
}
