package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	values, err := json.Marshal(map[string]string{
		"ApiKey":    "haha",
		"SecretKey": "haha",
		"Brandname": "haha",
		"Content":   "haha",
		"Phone":     "haha",
		"SmsType":   "2",
		// "RequestId": fmt.Sprintf("phone:sendSMS:user:%d", user.ID),
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(values))

}
