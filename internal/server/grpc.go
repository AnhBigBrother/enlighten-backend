package server

import (
	"crypto/tls"
	"log"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/interceptor"
	"github.com/AnhBigBrother/enlighten-backend/internal/pb"
	"github.com/AnhBigBrother/enlighten-backend/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewGrpcServer() *grpc.Server {
	tlsCredentials, err := loadTSLCredentials()
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	server := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(
			interceptor.LoggingInterceptor,
			interceptor.AuthInterceptor,
			interceptor.UnaryAuthGuard(
				[]string{
					"pb.User",
					"pb.Post",
					"pb.Comment",
				},
				[]string{}),
		),
	)

	publicServer := services.NewPublicServer()
	userServer := services.NewUserServer()
	postServer := services.NewPostServer()
	commentServer := services.NewCommentServer()
	gameServer := services.NewGameServer()
	oauthServer := services.NewOauthServer()

	pb.RegisterPublicServer(server, publicServer)
	pb.RegisterUserServer(server, userServer)
	pb.RegisterPostServer(server, postServer)
	pb.RegisterCommentServer(server, commentServer)
	pb.RegisterGameServer(server, gameServer)
	pb.RegisterOauthServer(server, oauthServer)

	return server
}

func loadTSLCredentials() (credentials.TransportCredentials, error) {
	// load server certificate and private key
	serverCert, err := tls.LoadX509KeyPair(cfg.ServerCertificateFile, cfg.ServerKeyFile)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}
