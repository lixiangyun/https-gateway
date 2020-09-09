module github.com/lixiangyun/https-gateway

go 1.14

require (
	github.com/astaxie/beego v1.12.2
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/google/uuid v1.1.2
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
