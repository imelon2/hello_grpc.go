## About `protoc` (Protocol Compiler):
- Google에서 개발한 공식 프로토콜 버퍼 컴파일러입니다.
- .proto 파일을 컴파일하여 해당 언어(예: C++, Java, Python, Go 등)의 코드를 생성합니다.
- 주요 기능은 프로토콜 버퍼 정의를 컴파일하여 다양한 언어로 **Compile**하는 것입니다.

## Before Begin
#### (1) Proto Budders(.proto)를 Go로 Compile하기 위해 `protoc`를 설치합니다.

```` shell
# https://grpc.io/docs/protoc-installation/
brew install protobuf
protoc --version
````

#### (2) Backend와 grcp 통신을 위한 CLI Tool `grpcurl`와 `buf`를 설치합니다.
```` shell
# https://github.com/fullstorydev/grpcurl
brew install grpcurl
grpcurl --version


# https://buf.build/docs/installation
brew install bufbuild/buf/buf
buf --version
````

<br/>
<br/>

## Run the example 
### (1) Generate go proto (.pd.go)
`proto/greet.proto`에 구현된 인터페이스에 따라 go proto 파일(.pd.go)로 Compile 합니다.

```` shell
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/greet.proto
````

### (2) Backend 실행
port : *50051*
```` shell
go run cmd/main.go 
````

### (3) grpc list 조회
서버에서 호출 가능한 grpc Class를 조회합니다
````shell
grpcurl --plaintext localhost:50051 list 
````
> 서버(main.go)에서 `reflection.Register(s)`을 설정해줬을 경우, 호출이 가능합니다.

#### *reflection 란?
많은 gRPC 관련 도구는 호출자가 런타임에 서비스의 Protobuf 스키마에 액세스할 수 있도록 해주는 서버 reflection에 의존합니다. Connect는 connectrpc.com/grpcreflect 패키지로 서버 리플렉션을 지원합니다. 

### (4) grpc 메시지 전송
`grpcurl` 또는 `buf curl` CLI를 사용하여 gRPC를 호출합니다.
```` shell
# grpcurl version
grpcurl --plaintext -d '{"name":"HI grpcurl"}' localhost:50051 greet.Greeter/SayHello

# buf version
buf curl \
  --schema proto \
  --protocol grpc \
  --http2-prior-knowledge \
  --data '{"name": "HI buf curl"}' \
  http://localhost:50051/greet.Greeter/SayHello
````