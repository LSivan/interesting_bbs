package feature

import (
	"git.oschina.net/gdou-geek-bbs/models"
	"time"
	"encoding/json"
	"git.oschina.net/gdou-geek-bbs/common"
	"strconv"
)

// TODO 话题的特征值提取与存储 - fin
// TODO 用户特征值如何变化 每做一次对文章的操作,就把特征值变化一次,每天凌晨根据用户的特征值和文章的特征值来计算各个用户喜好文章列表
// 假设N个用户,M篇文章,M为长度的ID列表,共N个,每天更新,user-favorite 1 [2,4,12,64,24]
// [0,1,23,43,4,54,5,4,] [3,4,3,5,67,8,7,8,98,]
// TODO 用户的特征值提取与存储
// TODO 更改推荐的算法
// 将所有的文章分析好特征并存到redis中
func InitTopicFeature() {
	// 取出所有的文章
	topics := models.FindTopicFrom(0,models.CountTopicFromID(0))
	for _, topic := range topics {
		go func(t *models.Topic) {
			feature := GetTopicFeature(t)
			common.Redis.HSet("topic-feature",strconv.Itoa(feature.ID),feature)
		}(topic)

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