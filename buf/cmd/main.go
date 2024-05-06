package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	greetv1 "buf/greet/gen/greet/v1"
	greetv1connect "buf/greet/gen/greet/v1/greetv1connect"

	connect "connectrpc.com/connect"
	grpcreflect "connectrpc.com/grpcreflect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/proto"
)

const address = "localhost:50051"

func main() {
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreeterServiceHandler(&greeterServiceServer{})
	mux.Handle(path, handler)

	// reflection
	reflector := grpcreflect.NewStaticReflector(
		greetv1connect.GreeterServiceName,
	)

	mux.Handle(grpcreflect.NewHandlerV1(reflector))

	fmt.Println("... Listening on", address)
	http.ListenAndServe(
		address,
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

type greeterServiceServer struct {
	greetv1connect.UnimplementedGreeterServiceHandler
}

func (s *greeterServiceServer) SayHello(
	ctx context.Context,
	req *connect.Request[greetv1.SayHelloRequest],
) (*connect.Response[greetv1.SayHelloResponse], error) {

	msg := req.Msg
	fmt.Println("SayHello original data:", msg)

	data, err := proto.Marshal(req.Msg)

	if err != nil {
		log.Fatalf("protobuf marshal error: %v", err)
	}

	// 직렬화된 바이너리 데이터를 출력합니다.
	encodedData := base64.StdEncoding.EncodeToString(data)
	fmt.Println("Serialized binary data:", encodedData)

	name := req.Msg.GetName()
	petType := req.Msg.ProtoReflect()
	log.Printf("Got a request to create a %v named %s", petType, name)
	return connect.NewResponse(&greetv1.SayHelloResponse{}), nil
}
