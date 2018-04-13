package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"github.com/huichen/sego"
)
var dictionary = "feature/dic.txt"
var keywordMap map[string]int // 存放关键字和数组对应关系的m
// 载入词典
var segmenter sego.Segmenter

func init() {
	keywordMap = make(map[string]int, 300)
	// 载入词典
	segmenter.LoadDictionary(dictionary)
	// 记录关键字和数组对应下标映射关系
	for _, file := range strings.Split(dictionary, ",") {
		dictFile, _ := os.Open(file)
		defer dictFile.Close()

		reader := bufio.NewReader(dictFile)
		var text string
		var freqText string
		var pos string
		i := 0
		// 逐行读入分词
		for {
			size, _ := fmt.Fscanln(reader, &text, &freqText,&pos)
			if size == 0 {
				// 文件结束
				break
			} else if size < 2 {
				// 无效行
				continue
			}
			keywordMap[text] = i
			i++
		}
	}
}

func main() {
	fmt.Println("main")
	fmt.Println(keywordMap["中国人"])
}