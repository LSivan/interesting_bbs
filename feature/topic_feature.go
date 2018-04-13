package main

import (
	"git.oschina.net/gdou-geek-bbs/models"
	"github.com/huichen/sego"
	"time"
	"encoding/json"
)

// TODO 话题的特征值提取与存储
// TODO 用户的特征值提取与存储
// TODO 更改推荐的算法
// 将所有的文章分析好特征并存到redis中
func init() {

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
		sego.SegmentsToFeatureSlice(segmenter.Segment([]byte(topic.Title+topic.Content+topic.Section.Name)),0),
		topic.InTime,
		topic.Id,
	}
}