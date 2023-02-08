server:
	kitex --thrift-plugin validator -module github.com/seed30/TikTok idl/user.thrift # execute in the project root directory
	kitex --thrift-plugin validator -module github.com/seed30/TikTok -service userserver -use github.com/seed30/TikTok/kitex_gen ../../idl/user.thrift # execute in cmd/user