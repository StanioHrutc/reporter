package logic

import (
	"bytes"
	"fmt"
	"net/http"
)

func UpdateStatusBoard(config Config, msg []byte) {
	host, port := config.BoardConfig.StatusBoardHost, config.BoardConfig.StatusBoardPort

	requestURL := fmt.Sprintf("%s:%d", host, port)
	// var jsonStr = []byte(`{"slot_number": 1, "status": "free", "timestamp": 1673232362}`)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(msg))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error while updating slots information: %s\n", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Got response code: %v", resp.StatusCode)
}
