package main

import "fmt"

func main()  {
	logger := DefaultLogger(false)
	logger.ReceiveLog(func(log string) {
		fmt.Println("receive: ",log)
	})

	logger.DEBUG("hello debug")
	logger.INFO("hello info")
	logger.ERROR("hello error")
	logger.WARN("hello warn")
	logger.SERIOUS("hello serious")
}
