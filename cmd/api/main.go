package main

import (
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"os"

	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/jamalyusuf/avoxi-challenge/pkg/gateway"
	"github.com/jamalyusuf/avoxi-challenge/pkg/insecure"
	"github.com/jamalyusuf/avoxi-challenge/pkg/server"
	pbService "github.com/jamalyusuf/avoxi-challenge/proto"

	// Static files
	_ "github.com/jamalyusuf/avoxi-challenge/statik"
)

// getOpenAPIHandler serves an OpenAPI UI.
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")

	statikFS, err := fs.New()
	if err != nil {
		panic("creating OpenAPI filesystem: " + err.Error())
	}

	return http.FileServer(statikFS)
}

func main() {
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	logger := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, os.Stderr)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	grpclog.SetLoggerV2(logger)

	addr := "0.0.0.0:10000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer(
		// TODO: Replace with real certificate
		grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
	)
	//TODO: should be only on for dev
	reflection.Register(s)

	pbService.RegisterIPFilterServiceServer(s, server.New(infoLog, errorLog))

	// Serve gRPC Server
	infoLog.Printf("Serving gRPC on https://%s", addr)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	err = gateway.Run("dns:///" + addr)
	log.Fatalln(err)
}
