package utils

import "log"

func LogError(action string, err error) {
	if err != nil {
		log.Printf("%s失败，err：%v\n", action,err)
	}
	log.Printf("%s成功\n", action)
}