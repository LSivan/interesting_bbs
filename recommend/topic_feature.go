package recommend

import (
	"git.oschina.net/gdou-geek-bbs/models"
	"time"
	"encoding/json"
	"git.oschina.net/gdou-geek-bbs/common"
	"strconv"
	"sync"
	"github.com/astaxie/beego"
)

// TODO 话题的特征值提取与存储 - fin
// TODO 用户特征值如何变化 每做一次对文章的操作,就把特征值变化一次,每天凌晨根据用户的特征值和文章的特征值来计算各个用户喜好文章列表 - fin
// TODO 用户的特征值提取与存储 - fin
// TODO 新用户怎么办,怎么推荐
// TODO 更改推荐的算法
// 将所有的文章分析好特征并存到redis中
func InitTopicFeature() {
	// 取出所有的文章
	if common.Redis.Exists("topic-feature").Val() == 0 {
		topics := models.FindTopicFrom(0,models.CountTopicFromID(0))
		wg := &sync.WaitGroup{}
		for _, topic := range topics {
			wg.Add(1)
			go func(t *models.Topic) {
				feature := GetTopicFeature(t)
				err := common.Redis.HSet("topic-feature",strconv.Itoa(t.Id),feature).Err()
				if err != nil {
					// 一般分词没有结果,导致feature.tokens为[NaN,NaN]而不能存进redis中,从而过滤无关的文章,TODO 可以用来更新关键字列表
					beego.BeeLogger.Info("recommend.ID:%v, err :%v",feature.ID, err)
				}
				defer wg.Done()
			}(topic)
		}
		wg.Wait()
	}
}

type TopicFeature struct {
	Tokens []float64
	T time.Time
	ID int
}

func (tf TopicFeature) MarshalBinary() (data []byte, err error){
	return json.Marshal(tf)
}
func (tf *TopicFeature) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, tf)
}
// 文章的特征存不存
// 存了
func GetTopicFeature(topic *models.Topic) *TopicFeature{
	return &TopicFeature{
		SegmentsToFeatureSlice(segment.Segment([]byte(topic.Title+topic.Content+topic.Section.Name)),0),
		topic.InTime,
		topic.Id,
	}
}