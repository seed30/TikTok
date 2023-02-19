package constants

const (
	APIServiceName         = "douyin.api"
	BaseServiceName        = "douyin.base"
	InteractionServiceName = "douyin.interaction"
	SocialServiceName      = "douyin.social"
	EtcdAddress            = "127.0.0.1:2379"
	BaseTCPAddr            = "127.0.0.1:8889"
	InteractionTCPAddr     = "127.0.0.1:8890"
	SocialTCPAddr          = "127.0.0.1:8891"
	MySQLDefaultDSN        = "tiktok:tiktok-7306@tcp(localhost:3309)/tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	UploadAddr             = "http://192.168.254.84:8888/upload/" //客户端测试时ip地址改为自己的无线局域网ipv4地址
	VideoCountLimit        = 30
	UserNameMaxLen         = 32
	PassWordMaxLen         = 32
)
