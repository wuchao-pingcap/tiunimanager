package client

import (
	"crypto/tls"
	"tcp/addon/logger"
	"tcp/addon/tracer"
	"tcp/config"
	tcpPb "tcp/proto/tcp"

	_ "github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/transport"
)

// Make request
/*
	rsp, err := TcpClient.Hello(context.Background(), &pb.HelloRequest{
		Name: "Foo",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Greeting)
*/
var TcpClient tcpPb.TcpService

func init() {
	// tls
	log := logger.WithContext(nil).WithField("init", "client")
	cert, err := tls.LoadX509KeyPair(config.GetCertificateCrtFilePath(), config.GetCertificateKeyFilePath())
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err:%s", err)
		return
	}
	tlsConfigPtr := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	// create a new service
	service := micro.NewService(
		micro.WrapClient(opentracing.NewClientWrapper(tracer.GlobalTracer)),
		micro.Transport(transport.NewHTTPTransport(transport.Secure(true), transport.TLSConfig(tlsConfigPtr))),
	)

	// parse command line flags
	service.Init()

	// Use the generated client stub
	TcpClient = tcpPb.NewTcpService("go.micro.tcp", service.Client())
}