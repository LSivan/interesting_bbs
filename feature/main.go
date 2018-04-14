package feature

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"github.com/huichen/sego"
	"math"
	"github.com/astaxie/beego"
)
var dictionary = "feature/dic.txt"
var KeywordMap map[string]int // 存放关键字和数组对应关系的m
var segment = &sego.Segmenter{}
func init() {
	KeywordMap = make(map[string]int, 300)
	segment.LoadDictionary(dictionary)
	// 记录关键字和数组对应下标映射关系
	for _, file := range strings.Split(dictionary, ",") {
		dictFile, _ := os.Open(file)
		defer dictFile.Close()

		reader := bufio.NewReader(dictFile)
		var text string
		i := 0
		// 逐行读入分词
		for {
			size, _ := fmt.Fscanln(reader, &text)
			if size == 0 {
				// 文件结束
				break
			}
			if KeywordMap[strings.ToLower(text)] != 0 {
				beego.BeeLogger.Info("duplicate text:%v",text)
			}
			KeywordMap[strings.ToLower(text)] = i
			i++
		}
	}
	beego.BeeLogger.Info("keywordMap:%v",len(KeywordMap))
}

/**
salt 影响特征值的一个数
 */
func SegmentsToFeatureSlice(segs []sego.Segment, salt float64) (feature []float64) {
	output := make([]int,len(KeywordMap))
	feature = make([]float64,len(KeywordMap))
	total := 0
	for _, seg := range segs {
		word := seg.Token().Text()
		if index,exist := KeywordMap[strings.ToLower(word)];exist {
			i := output[index]
			output[index] = i + 1
			total++
		}
	}
	var log10 = func(num float64) float64 {
		num = math.Log10(num+10+salt)
		// 防止恶意刷关键字,高于200的就会多做一次取对数的操作
		for num > 2.33 {
			num = math.Log10(num)
		}
		return num - 1
	}
	for index, i := range output {
		feature[index] = log10(float64(i))+float64(i)/float64(total)
	}
	return
}