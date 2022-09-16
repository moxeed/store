package common

import "log"

func Log(err error) {
	if err != nil {
		log.Println(err)
	}
}
