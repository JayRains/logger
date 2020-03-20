package main

import (
	"fmt"
	"logger/log"
)

func main()  {
	logger := log.DefaultLogger(false)
	logger.ReceiveLog(func(log string) {
		fmt.Println("receive: ",log)
	})

	logger.DEBUG("hello debug")
	logger.INFO("hello info")
	logger.ERROR("hello error")
	logger.WARN("hello warn")
	logger.SERIOUS("hello serious")
}
