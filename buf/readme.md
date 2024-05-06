## About `buf` (Protocol Compiler):
- 버퍼프로젝트(Buffer)에서 개발한 오픈 소스 도구입니다.
- 프로토콜 버퍼 관리와 유지 보수를 위한 도구로서, `protoc`보다 더 강력한 기능을 제공합니다.
- .proto 파일의 구문 분석, 검증, 정적 분석을 수행하여 코드의 일관성을 유지하고 코드 품질을 향상시킵니다.
- 코드 생성 뿐만 아니라 linting, breaking change detection 등의 추가 기능을 제공합니다.
- 버전 관리 시스템과의 통합을 위한 기능도 있습니다.

## Before Begin
실습을 위해 buf를 설치합니다.
```
# https://github.com/fullstorydev/grpcurl
brew install grpcurl
grpcurl --version


# https://buf.build/docs/installation
brew install bufbuild/buf/buf
buf --version
```

## Run the example
### (1) init buf.yaml
`buf.yaml` 파일은 Buf 도구에서 사용하는 설정 파일입니다. Buf는 프로토콜 버퍼(Protocol Buffers) 프로젝트를 관리하고 유지 보수하기 위한 도구로, 프로토콜 버퍼 파일의 구조를 검증하고 코드 생성에 필요한 설정을 지정하는 등의 작업을 수행합니다. </br># [공식문서 바로가기](https://buf.build/docs/configuration/v1/buf-yaml)
1. 버퍼 룰(Rules): **프로토콜 버퍼 파일의 구조를 검증**하는 데 사용되는 규칙을 지정합니다. `lint` 명령어를 통해 정적인 규칙 검사를 지원합니다.
2. 설정(Configuration): 코드 생성 및 빌드와 관련된 설정을 지정합니다. 예를 들어, 언어별 코드 생성 경로, **프로토콜 버퍼 컴파일러 설정** 등을 포함할 수 있습니다.
3. 추가적인 설정(Additional Configurations): Buf의 플러그인이나 다른 도구와의 통합을 위한 추가적인 설정을 지정할 수 있습니다.

</br>

````
cd ./proto
buf mod init
````
> `.proto` 파일 경로가 확인되는 방식이므로 항상 `.proto` 파일 계층 `root`의 루트에 `buf.yaml` 파일을 두는 것이 좋습니다.


## (2) init buf.gen.yaml
`buf.gen.yaml`은 buf generate 명령이 선택한 언어에 대한 통합 코드를 생성하는 데 사용하는 구성 파일입니다.  </br># [공식문서 바로가기](https://buf.build/docs/configuration/v1/buf-gen-yaml)

````
cd ../

echo "version: v1
managed:
  enabled: true
  go_package_prefix:
    default: buf/greet/gen
plugins:
  - plugin: buf.build/protocolbuffers/go
    out: gen
    opt: paths=source_relative
  - plugin: buf.build/connectrpc/go
    out: gen
    opt: paths=source_relative" \
> buf.gen.yaml
````

## (3) Generate code
`buf.gen.yaml`에 작성된 내용에 따라 Protocol Buffers를 Compile합니다.
````
buf generate proto
````
> 해당 실습은 GO 언어로 진행되며, `buf.build/protocolbuffers/go` 컴파일러를 사용합니다.

> `out: gen`으로 설정했기 때문에, `./gen`디렉토리에 Compile된 코드(.pd.go)가 생성됩니다.

#### *template 사용법
buf generate는 기본적으로 현재 위치의 `buf.gen.yaml`를 추적합니다. 다른 `.yaml` 적용을 원하는 경우 `--template [PATH]` flag를 사용하면 됩니다.
```
buf generate --template ./proto/buf.gen.go.yaml 
```

## (3) Lint API
`buf lint`를 실행하면 buf는 `buf.yaml`의  `lint:use: -DEFAULT` 에 설정된 옵션에 따라 모든 Protobuf 파일에 대해 일련의 **lint rules**을 검사합니다.

linit rules란? : 프로젝트에서 정의한 디렉토리 구조(ex. pet/v1/pet.proto) 또는 Prifix/Suffix 관례 등의 규칙을 의미합니다.
<br/>[**lnit rules** docs 바로가기](https://buf.build/docs/lint/rules)

````
buf lint proto
````

### (5) Backend 실행
port : *50051*
```` shell
go run cmd/main.go 
````

### (6) grpc list 조회
서버에서 호출 가능한 grpc Class를 조회합니다
````shell
grpcurl --plaintext localhost:50051 list 
````
> 서버(main.go)에서 `reflection.`을 설정해줬을 경우, 호출이 가능합니다.

#### *reflection 란?
많은 gRPC 관련 도구는 호출자가 런타임에 서비스의 Protobuf 스키마에 액세스할 수 있도록 해주는 서버 reflection에 의존합니다. Connect는 connectrpc.com/grpcreflect 패키지로 서버 리플렉션을 지원합니다. 

### (7) grpc 메시지 전송
`grpcurl` 또는 `buf curl` CLI를 사용하여 gRPC를 호출합니다.
```` shell
# grpcurl version
grpcurl --plaintext -d '{"name":"HI buf"}' localhost:50051 greet.v1.GreeterService/SayHello

# buf version
buf curl \
  --schema proto \
  --protocol grpc \
  --http2-prior-knowledge \
  --data '{"name": "HI buf curl"}' \
  http://localhost:50051/greet.v1.GreeterService/SayHello
````