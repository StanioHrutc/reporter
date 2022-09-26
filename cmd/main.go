package main

import (
	"BoardReporter/logic"
)

func main() {
	config := logic.GetConfig()

	rc := logic.RabbitConsumer{
		Conf: *config,
	}
	rc.Consume(logic.UpdateStatusBoard)
}
