package main

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"fmt"
	"xxd-server/api"
)

func main() {
	test3()
}


func test3(){
	gid := "25ac51b3-830c-6e21-3635-35b41dd697a8"
	//gid := "2&3"
	if members , err := api.GetChatMemberByGID(gid); err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(members)
	}
}

func test2() {
	api.SetShowSQL(true)
	if users, err := api.GetAllSystemUsers(); err != nil {
		fmt.Println("111", err)
	} else {
		fmt.Println("xxx",users)
		for _, user := range users {
			fmt.Println("xxxx", *user, "xxxx")
		}
	}

}

func test1() {

	startXXD := []byte(`{"module":"chat","method":"serverStart"}`)

	message, _ := api.AesEncrypt(startXXD, []byte("88888888888888888888888888888888"))

	fmt.Println(">>>", string(message))

	req, err := http.NewRequest("POST", "http://localhost:4000/xuanxuan", bytes.NewReader(message))
	if err != nil {
		fmt.Println(err)
	}

	var client *http.Client
	client = &http.Client{}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "easysoft-xxdClient/1.0.0")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("<<<", resp.Status, string(body))
}
