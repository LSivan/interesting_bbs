package recommend

import (
	"encoding/json"
	"git.oschina.net/gdou-geek-bbs/common"
	"git.oschina.net/gdou-geek-bbs/models"
	"github.com/astaxie/beego"
	"strconv"
	"sync"
	"fmt"
)
var DefaultUserFeature UserFeature

func init() {
	DefaultUserFeature = UserFeature(defaultFeatures)
	go func() {
		beego.BeeLogger.Info("watch the user's operations and change his recommend")
		for {
			select {
			case topicOperation := <-change:
				beego.BeeLogger.Info(" change his recommend : %v",topicOperation)

				// 加盐算出用户操作的那个话题的特征
				features := SegmentsToFeatureSlice(segment.Segment([]byte(topicOperation.topic.Title+topicOperation.topic.Content+topicOperation.topic.Section.Name)), topicOperation.operation)

				// 取出用户的特征数组
				// 假设redis中没有用户的特征值
				if !common.Redis.HExists("user-feature", strconv.Itoa(topicOperation.userId)).Val() {
					// 存一个默认的特征值数组到redis中
					common.Redis.HSet("user-feature", strconv.Itoa(topicOperation.userId), UserFeature(defaultFeatures))
				}
				b, err := common.Redis.HGet("user-feature", strconv.Itoa(topicOperation.userId)).Bytes()
				if err != nil {
					beego.BeeLogger.Warn("fail to get user feature:%v", err)
				}
				uf := UserFeature([]float64{})
				uf.UnmarshalBinary(b)
				// 重新计算
				for i, v := range features {
					uf[i] = (uf[i] + v) / 2
				}
				// 放回到redis中
				common.Redis.HSet("user-feature", strconv.Itoa(topicOperation.userId), uf)
			}
		}
	}()
}

type UserFeature []float64

// 实现了MarshalBinary的struct才能放进redis中
func (uf UserFeature) MarshalBinary() (data []byte, err error) {
	return json.Marshal(uf)
}
func (uf *UserFeature) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, uf)
}

var change = make(chan *TopicOperation, 100)

type TopicOperation struct {
	userId    int
	topic     *models.Topic
	operation float64 // 用户对话题的操作,作为盐值,拉黑=-2,查看=0.5,点赞=1,评论=2,收藏=3,发表=4
}

func (t *TopicOperation)String() string {
	return fmt.Sprintf("userId:%v,topic:%v-%v,operation:%v",t.userId,t.topic.Id,t.topic.Title,t.operation)
}

func ChangeUserFeature(userId int, operation float64, topic *models.Topic) {
	change <- &TopicOperation{
		userId:    userId,
		operation: operation,
		topic:     topic,
	}
}

// 假设没有的话,从数据库获取UserTopicList的记录,可以取出发表,收藏和拉黑的所有记录,然后通过ChangeUserFeature方法进行用户特征的抽取
// 但是解决不了没有任何操作的用户的特征列表
func InitUserFeature() {
	// 取出所有的文章
	if common.Redis.Exists("user-feature").Val() == 0 {
		userTopicList := models.FindAllUserTopicList()
		wg := &sync.WaitGroup{}
		for _, ut := range *userTopicList {
			wg.Add(1)
			go func(t models.UserTopicList) {
				ChangeUserFeature(t.Topic.User.Id, 4, t.Topic)
				if t.ActionType == 1 {
					ChangeUserFeature(t.User.Id, 3, t.Topic)
				} else if t.ActionType == 0 {
					ChangeUserFeature(t.User.Id, -2, t.Topic)
				}
				wg.Done()
			}(ut)
		}
		wg.Wait()
	}
}
