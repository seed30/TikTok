package main

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/seed30/TikTok/dal"
	feed "github.com/seed30/TikTok/kitex_gen/feed/feedservice"
	"github.com/seed30/TikTok/pkg/consts"
	"github.com/seed30/TikTok/pkg/jwt"
	"github.com/seed30/TikTok/pkg/mw"
	"net"
)

var Jwt *jwt.JWT

func Init() {
	dal.Init()
	klog.SetLogger(kitexlogrus.NewLogger())
	klog.SetLevel(klog.LevelInfo)
	Jwt = jwt.NewJWT([]byte("JWT.signingKey"))
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{consts.ETCDAddress})
	if err != nil {
		klog.Fatal(err)
	}
	addr, err := net.ResolveTCPAddr(consts.TCP, consts.FeedServiceAddr)
	if err != nil {
		klog.Fatal(err)
	}

	Init()

	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(consts.FeedServiceName),
		provider.WithExportEndpoint(consts.ExportEndpoint),
		provider.WithInsecure(),
	)

	svr := feed.NewServer(new(FeedServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		server.WithMiddleware(mw.CommonMiddleware),
		server.WithMiddleware(mw.ServerMiddleware),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.UserServiceName}),
	)

	err = svr.Run()

	if err != nil {
		klog.Fatal(err)
	}
}
