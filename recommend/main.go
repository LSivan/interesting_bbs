package recommend

import (
	"bufio"
	"fmt"
	"git.oschina.net/gdou-geek-bbs/common"
	"github.com/astaxie/beego"
	"github.com/huichen/sego"
	"math"
	"os"
	"strings"
	"encoding/json"
)

var dictionary = "recommend/dic.txt"
var KeywordMap map[string]int // 存放关键字和数组对应关系的m
var segment = &sego.Segmenter{}
var defaultFeature []float64

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
				beego.BeeLogger.Info("duplicate text:%v", text)
			}
			KeywordMap[strings.ToLower(text)] = i
			i++
		}
	}
	beego.BeeLogger.Info("keywordMap:%v", len(KeywordMap))
	defaultFeature = make([]float64, len(KeywordMap))
}

/**
salt 影响特征值的一个数
*/
func SegmentsToFeatureSlice(segs []sego.Segment, salt float64) (feature []float64) {
	output := make([]int, len(KeywordMap))
	feature = make([]float64, len(KeywordMap))
	total := 0
	for _, seg := range segs {
		word := seg.Token().Text()
		if index, exist := KeywordMap[strings.ToLower(word)]; exist {
			i := output[index]
			output[index] = i + 1
			total++
		}
	}
	var log10 = func(num float64) float64 {
		num = math.Log10(num + 10 + salt)
		// 防止恶意刷关键字,高于200的就会多做一次取对数的操作
		for num > 2.33 {
			num = math.Log10(num)
		}
		return num - 1
	}
	for index, i := range output {
		feature[index] = log10(float64(i)) + float64(i)/float64(total)
	}
	return
}

// 定时触发
// 假设N个用户,M篇文章,M为长度的ID列表,共N个,每天更新,user-favorite 1 [2,4,12,64,24]
// [0,1,23,43,4,54,5,4,] [3,4,3,5,67,8,7,8,98,]
func GetUsersFavoriteList() {
	userResult, err := common.Redis.HGetAll("user-feature").Result()
	if err != nil {
		beego.BeeLogger.Error("fail to get user-feature :%v",err)
		return
	}
	topicResult, err := common.Redis.HGetAll("topic-feature").Result()
	if err != nil {
		beego.BeeLogger.Error("fail to get topic-feature :%v",err)
		return
	}
	for userId, userFeature := range userResult {
		uf := UserFeature([]float64{})
		uf.UnmarshalBinary([]byte(userFeature))
		ft := &FavoriteTopic{
			TopicIDS:      make([]int,0,500),
			TopicFeatures: make([]float64,0, 500),
		}
		index := 0
		for _, topicFeature := range topicResult {
			tf := TopicFeature{}
			tf.UnmarshalBinary([]byte(topicFeature))

			total := 0.00
			for i, v := range uf {
				total += tf.Tokens[i] * v
			}
			// 算上时间特征 TODO 是否需要,Top-n的列表的维护难度要变大了
			//total += math.Log10(time.Now().Truncate(24 * time.Hour).Sub(tf.T.Truncate(24 * time.Hour)).Hours()+11) - 1
			// 构建一个Top-N的堆
			if index < 500 {
				ft.TopicIDS = append(ft.TopicIDS, tf.ID)
				ft.TopicFeatures = append(ft.TopicFeatures, total)
				index++
				if index == 500 {
					ft.buildHeap()
				}
			} else {
				// 数组已填满,进行堆的调整
				if total > ft.TopicFeatures[0] {
					ft.TopicFeatures[0] = total
					ft.adjustHeap( 0)
				}
			}
			//  quickSort
			ft.quickSort(0,len(ft.TopicFeatures)-1)
		}
		//ft.TopicFeatures = ft.TopicFeatures[500-index:]
		//ft.TopicIDS = ft.TopicIDS[500-index:]
		common.Redis.HSet("user-favorite",userId,ft)
	}
}

type FavoriteTopic struct {
	TopicIDS      []int
	TopicFeatures []float64
}
func (ft FavoriteTopic) MarshalBinary() (data []byte, err error){
	return json.Marshal(ft)
}
func (ft *FavoriteTopic) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, ft)
}
func (ft *FavoriteTopic)quickSort(begin,end int) {
	i,j := begin, end
	index := ft.TopicFeatures[begin]
	if i < j {
		for i < j {
			if ft.TopicFeatures[i] < index {
				ft.TopicFeatures[i], ft.TopicFeatures[j] = ft.TopicFeatures[j], ft.TopicFeatures[i]
				ft.TopicIDS[i], ft.TopicIDS[j] = ft.TopicIDS[j], ft.TopicIDS[i]
				j--
			} else {
				i++
			}
		}
		if ft.TopicFeatures[i] <= index {
			i--
		}
		ft.TopicFeatures[i], ft.TopicFeatures[begin] = ft.TopicFeatures[begin], ft.TopicFeatures[i]
		ft.TopicIDS[i], ft.TopicIDS[begin] = ft.TopicIDS[begin], ft.TopicIDS[i]
		ft.quickSort(begin,i)
		ft.quickSort(j,end)
	}
}

func (ft *FavoriteTopic)buildHeap() {
	for i := len(ft.TopicFeatures) - 1; i >= 0; i-- {
		ft.adjustHeap(i)
	}
}
func (ft *FavoriteTopic)adjustHeap( pos int) {
	node := pos
	length := len(ft.TopicFeatures)
	left, right := 2*node+1, 2*node+2
	child := 0
	if left > length { //无子节点
		return
	} else if right < length { // 两个子节点
		if ft.TopicFeatures[left] > ft.TopicFeatures[right] {
			child = right
		} else {
			child = left
		}
	} else { // 一个子节点
		child = left
	}
	if child > 0 && ft.TopicFeatures[child] < ft.TopicFeatures[node] {
		ft.TopicFeatures[node], ft.TopicFeatures[child] = ft.TopicFeatures[child], ft.TopicFeatures[node]
		ft.TopicIDS[node], ft.TopicIDS[child] = ft.TopicIDS[child], ft.TopicIDS[node]
		ft.adjustHeap(child)
	}
}
