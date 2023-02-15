package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// GetSHA256HashCode 基于 SHA256 算法生成哈希值
func GetSHA256HashCode(file *os.File) string {
	// 创建一个基于 SHA256 算法的 hash.Hash 接口的对象
	hash := sha256.New()
	_, _ = io.Copy(hash, file)
	// 计算哈希值
	bytes := hash.Sum(nil)
	// 将字符串编码为 16 进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	// 返回哈希值
	return hashCode

}

// EncodeMd5 基于 md5 加密
func EncodeMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
