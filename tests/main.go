package main

import (
	"github.com/huichen/sego"
	"fmt"
)

func main() {
	// 载入词典
	var segmenter sego.Segmenter
	segmenter.LoadDictionary("/home/sivan/go/src/git.oschina.net/gdou-geek-bbs/tests/dictionary.txt")

	// 分词
	text := []byte("中华人民共和国")
	segments := segmenter.Segment(text)

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	fmt.Println(sego.SegmentsToSlice(segments, true))
}
