package test

import (
	"github.com/huichen/sego"
	"fmt"
	"testing"
	"bufio"
	"os"
	"strings"
)

func TestSegment(t *testing.T) {
	// 载入词典
	var segmenter sego.Segmenter
	segmenter.LoadDictionary("/home/sivan/go/src/git.oschina.net/gdou-geek-bbs/tests/dic.txt")

	// 分词
	text := []byte("数据库测试报告:测试通过,表示很优良,Mysql数据库广告的性能真的很不错,通过使用XML配置,B2C SNS-EC C++ PR值 爱上当借口了Java F#B2B是一种完美的商业模式,我真的觉得超级棒")
	segments := segmenter.Segment(text)
	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	fmt.Println(sego.SegmentsToSlice(segments,false))

	fmt.Println(sego.SegmentsToFeatureSlice(segments,0))
}

func TestFinDic(t *testing.T) {
	dic, err := os.Open("/home/sivan/go/src/git.oschina.net/gdou-geek-bbs/tests/dic.txt")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(dic)
	line,_, _ := r.ReadLine()
	for string(line) != "" {
		ll := strings.Split(string(line)," ")
		if len(ll) == 3 {
			fmt.Println(string(line))
		}
		line,_, _ = r.ReadLine()
	}
}