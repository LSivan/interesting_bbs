package main

import (
	"github.com/huichen/sego"
	"fmt"
)

func main() {
	// 载入词典
	var segmenter sego.Segmenter
	segmenter.LoadDictionary("/home/sivan/go/src/git.oschina.net/gdou-geek-bbs/tests/dic.txt")

	// 分词
	text := []byte("数据库测试报告:测试通过,表示很优良,Mysql数据库的性能真的很不错,通过使用XML配置,B2C SNS-EC C++ PR值 爱上当借口了Java")
	segments := segmenter.Segment(text)

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	fmt.Println(sego.SegmentsToFreqMap(segments))
}
