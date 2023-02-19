package utils

import (
	"github.com/importcjj/sensitive"
	"log"
)

var Filter *sensitive.Filter

func FilterInit() {
	Filter = sensitive.New()
	err := Filter.LoadWordDict("utils/SensitiveDict.txt")
	if err != nil {
		log.Println("敏感词库加载失败")
	}

	comment := "这篇文章真的好垃圾"
	comment = Filter.Replace(comment, '*')
	log.Println(comment)
	// output => 这篇文章真的好**
}
