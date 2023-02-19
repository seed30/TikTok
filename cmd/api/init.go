package api

import (
	"github.com/seed30/TikTok/cmd/api/biz/handler"
	"github.com/seed30/TikTok/cmd/api/biz/router"
	"github.com/seed30/TikTok/pkg/tracer"
	"github.com/seed30/TikTok/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	hertztracer "github.com/hertz-contrib/tracer/hertz"
)

func NewDouyinApiHertz() *server.Hertz {
	hTracer, _ := tracer.InitTracer("douyin.api")
	svc := server.Default(
		server.WithTracer(hertztracer.NewTracer(hTracer, func(c *app.RequestContext) string {
			return "hertz.server" + "::" + c.FullPath()
		})))
	svc.Use(hertztracer.ServerCtx())
	h := handler.NewHandler(rpc.NewRPC())
	router.Register(svc, h)
	return svc
}
