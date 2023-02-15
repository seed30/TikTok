package rpc

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/seed30/TikTok/kitex_gen/feed"
	"github.com/seed30/TikTok/kitex_gen/feed/feedservice"
	"github.com/seed30/TikTok/pkg/consts"
	"github.com/seed30/TikTok/pkg/errno"
	"github.com/seed30/TikTok/pkg/mw"
	"time"
)

var feedClient feedservice.Client

// Feed RPC 客户端初始化
func initFeed() {
	r, err := etcd.NewEtcdResolver([]string{consts.ETCDAddress})
	if err != nil {
		panic(err)
	}

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(consts.ApiServiceName),
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	c, err := feedservice.NewClient(
		consts.ApiServiceName,
		client.WithMiddleware(mw.CommonMiddleware),
		client.WithInstanceMW(mw.ClientMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(30*time.Second),             // rpc timeout
		client.WithConnectTimeout(30000*time.Millisecond), // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(tracing.NewClientSuite()),        // tracer
		client.WithResolver(r),                            // resolver
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.ApiServiceName}),
	)
	if err != nil {
		panic(err)
	}
	feedClient = c
}

// GetUserFeed 传递获取视频流操作的上下文, 并获取 RPC Server 端的响应
func GetUserFeed(ctx context.Context, req *feed.FeedRequest) (resp *feed.FeedResponse, err error) {
	resp, err = feedClient.GetUserFeed(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(int64(resp.StatusCode), *resp.StatusMsg)
	}
	return resp, nil
}
