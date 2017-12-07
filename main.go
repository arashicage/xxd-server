package main

import (
	"log"
	"flag"
	"net/http"

	"xxd-server/api"
	"xxd-server/daemon"
	"fmt"
	//"encoding/json"

	"github.com/json-iterator/go"
)

/*

http://172.30.11.230/ranzhi/sys/user-create.html

*/

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {

	requestCN := map[string]string{
		"chat.serverStart":        "启动服务器(来自xxd),创建 system chat,重置所有用户状态为 offline",
		"chat.login":              "验证客户端用户密码; 通过验证后, 获取在线用户列表(来自xxd)",
		"chat.userGetlist":        "获取所有用户列表(来自xxd)",
		"chat.getlist":            "获取讨论组信息(来自xxd)",
		"chat.getOfflineMessages": "获取本用户的离线消息",
		"chat.userChange":         "修改用户信息在线状态",
		"chat.message":            "发送消息",
		"chat.create":             "创建 讨论组",
		"chat.changePublic":       "修改 讨论组 为 公开/私密",
		"chat.getpubliclist":      "获取 公共 讨论组",
		"chat.joinchat":           "加入/退出 讨论组",
		"chat.changename":         "重命名 讨论组",
		"chat.star":               "收藏/取消收藏 讨论组",
		"chat.hide":               "隐藏/显示 讨论组",
		"chat.members":            "获取 chat 的成员",
		"chat.logout":             "登出，注销",
		"chat.kickoff":            "重复登录时, 踢出前一次登陆, xxd 功能",
		"chat.history":            "获取聊天历史记录",
		"chat.uploadFile":         "上传文件",
	}

	//daemon.(daemon.DaemonTask(api.PingMysql))
	if false {
		daemon.DaemonTask(api.PingMysql).Tick()
	}

	conf := flag.String("conf", "xxd-server.yaml", "configuration file")

	http.HandleFunc("/xuanxuan", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "POST":

			data := make(map[string]interface{}, 0)
			req := new(api.RequestData)

			body := make([]byte, r.ContentLength)
			r.Body.Read(body)

			if decryptedReq, err := api.AesDecrypt(body, []byte(api.ServerToken)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - StatusInternalServerError!, AesDecrypt fail."))
				return
			} else {
				if err := json.Unmarshal(decryptedReq, &req); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 - StatusInternalServerError!, Unmarshal fail."))
					return
				} else {

					module := req.Module
					method := req.Method
					entrypoint := module + "." + method

					debug := entrypoint != "chat.ping" && entrypoint != "chat.settings"
					if true || debug {
						fmt.Println()
						fmt.Printf("Request: >>> %s\n", string(decryptedReq))
						if true {
							fmt.Printf("Request: >>> %s\n", requestCN[entrypoint])
						}
					}

					switch entrypoint {
					case "chat.serverStart":
						{
							// 来自 xxd 端的请求 chat.serverStart
							data = api.ChatServerStart(module, method)
							w.WriteHeader(data["status"].(int))

						}
					case "chat.login":
						{

							/* 登录时，先后调用 chat.login 两次，两次请求来源和参数不一样
							1. 先从 桌面客户端或浏览器 发送 chat.login 请求。
							 {"method":"login","module":"chat","params":["","admin","e10adc3949ba59abbe56e057f20f883e",""]}
							此时的 userID = "" status = "" ，后台的处理逻辑为验证用户的 用户和密码是否匹配，
							并获取用户的 userID，并用于后面的实际登陆逻辑

							2. 获取 userID 后，从 xxd 端发出 chat.login 请求。
							{"method":"login","module":"chat","params":["","admin","e10adc3949ba59abbe56e057f20f883e","online"],"userID":1}
							此时的 userID 有值 status = "online" ，后台的处理逻辑为登录系统，获取在线用户的列表。
							*/

							params := req.Params.([]interface{})
							account, password, status := params[1].(string), params[2].(string), params[3].(string)

							data = api.ChatLogin(module, method, account, password, status)

						}
					case "chat.logout":
						{
							// {"method":"logout","module":"chat","userID":2}

							data = api.ChatLogout(module, method, req.UserID)

						}
					case "chat.userGetlist":
						{
							// 获取所有用户列表，不附加 status, id 条件
							// 这里获取的信息，显示为客户端的左侧的 "联系人" 列表
							data = api.ChatUserGetlist(module, method, req.UserID)

						}
					case "chat.getlist":
						{
							// 获取该用户参与的 chat 列表。即客户端界面中的 "讨论组" 信息，其中的 members
							// 在客户端用来显示讨论组中有多少参与者
							api.SetShowSQL(true)
							data = api.ChatGetList(module, method, req)
							api.SetShowSQL(false)

						}
					case "chat.getOfflineMessages":
						{
							// 从 usermessage 获取到 user 的离线信息
							data = api.ChatGetOfflineMessages(module, method, req.UserID)

						}
					case "chat.userChange":
						{

							// {"module":"chat","method":"userChange","params":[{"status":"away"}],"userID":2}

							params := req.Params.([]interface{})
							userMap := params[0].(map[string]interface{})

							data = api.ChatUserChange(module, method, userMap, req.UserID)

						}
					case "chat.message":
						{

							// {"module":"chat","method":"message","params":{"messages":[{"gid":"7af873cf-8db7-4663-9779-0edce57afc6d","cgid":"2&3","type":"normal","contentType":"text","content":"online message 1","date":"","user":2}]},"userID":2}

							messages := req.Params.(interface{}).(map[string]interface{})

							data = api.ChatMessage(module, method, messages, req.UserID)

						}
					case "chat.create":
						{
							/*
							情况1. 100001 直接找 100002 开聊
								{"module":"chat","method":"create","params":["2&4","[Chat-NaN]","one2one",[2,4],0,false],"userID":2}

							情况2. 100001 选定 100002， 然后选择添加其他用户开讨论组，会前后两次 create
								{"module":"chat","method":"create","params":["2&4","[Chat-NaN]","one2one",[2,4],0,false],"userID":2}
								{"module":"chat","method":"create","params":["462e1d89-f10d-4c53-a613-09517ceccace","cccc","group",[3,2,4],0,false],"userID":2}
							*/

							data = api.ChatCreate(module, method, req.Params.([]interface{}), req.UserID)

						}
					case "chat.changePublic":
						{
							// {"module":"chat","method":"changePublic","params":["7adacd1c-01ac-444d-aef3-5daf22401fed",true],"userID":2}

							data = api.ChatChangePublic(module, method, req.Params.([]interface{}), req.UserID)

						}
					case "chat.getpubliclist":
						{
							// {"module":"chat","method":"getpubliclist","userID":2}

							data = api.ChatGetPublicList(module, method, req.UserID)

						}
					case "chat.joinchat":
						{
							// {"module":"chat","method":"joinchat","params":["1d4785d2-c69d-4957-a92f-df8a83653453",true],"userID":5}

							data = api.ChatJoinChat(module, method, req.Params.([]interface{})[0].(string), req.UserID, req.Params.([]interface{})[1].(bool))

						}
					case "chat.changename":
						{
							// {"module":"chat","method":"changename","params":["021d859c-85b3-43b5-902e-7bd5b192b376","123456"],"userID":2}

							gid, name, userid := req.Params.([]interface{})[0].(string), req.Params.([]interface{})[1].(string), req.UserID

							data = api.ChatChangeName(module, method, gid, name, userid)

						}
					case "chat.addmember":
						{
							// {"module":"chat","method":"addmember","params":["4a44ef96-bd78-40d9-b41e-43c6b4559e49",[1],true],"userID":2}

							gid, join, userID := req.Params.([]interface{})[0].(string),
								req.Params.([]interface{})[2].(bool),
								req.UserID

							members := make([]int, 0)
							for _, member := range req.Params.([]interface{})[1].([]interface{}) {
								members = append(members, int(member.(float64)))
							}
							api.SetShowSQL(true)
							data = api.ChatAddMember(module, method, gid, members, join, userID)
							api.SetShowSQL(false)

						}
					case "chat.members":
						{
							// {"module":"chat","method":"addmember","params":["4a44ef96-bd78-40d9-b41e-43c6b4559e49",[1],true],"userID":2}

							gid, userID := req.Params.([]interface{})[0].(string), req.UserID

							data = api.ChatMembers(module, method, gid, userID)

						}
					case "chat.star":
						{
							// {"module": "chat", "method": "star", "params": ["021d859c-85b3-43b5-902e-7bd5b192b376", true], "userID":2}

							gid, star, userID := req.Params.([]interface{})[0].(string), req.Params.([]interface{})[1].(bool), req.UserID

							data = api.ChatStar(module, method, gid, star, userID)

						}
					case "chat.hide":
						{
							// {"module": "chat", "method": "hide", "params": ["021d859c-85b3-43b5-902e-7bd5b192b376", true], "userID":2}

							gid, star, userID := req.Params.([]interface{})[0].(string), req.Params.([]interface{})[1].(bool), req.UserID

							data = api.ChatHide(module, method, gid, star, userID)

						}
					case "chat.history":
						{
							// {"module":"chat","method":"history","params":["2&4",50,1,0,true],"userID":2}

							gid, recPerPage, pageID, recTotal, continued, userID := req.Params.([]interface{})[0].(string),
								int(req.Params.([]interface{})[1].(float64)),
								int(req.Params.([]interface{})[2].(float64)),
								int(req.Params.([]interface{})[3].(float64)),
								req.Params.([]interface{})[4].(bool), req.UserID

							data = api.ChatHistory(module, method, gid, int(recPerPage), int(pageID), int(recTotal), continued, userID)

						}
					case "chat.setCommitters":
						{
							// {"module":"chat","method":"setCommitters","params":["4a44ef96-bd78-40d9-b41e-43c6b4559e49","2,1,5"],"userID":2}

							gid, commiters, userID := req.Params.([]interface{})[0].(string), req.Params.([]interface{})[1].(string), req.UserID

							data = api.ChatSetCommiters(module, method, gid, commiters, userID)

						}
					case "chat.settings":
						{

							/*
							获取客户端设置，上传，下载
							*/
							params := req.Params.([]interface{})
							account, settings := params[0].(string), params[1].(interface{})
							s, _ := json.Marshal(settings)
							api.SetShowSQL(true)
							data = api.ChatSettings(module, method, req.UserID, account, string(s))
							api.SetShowSQL(false)
						}
					case "chat.uploadFile":
						{
							// {"userID":2,"module":"chat","method":"uploadFile","params":["LayIM-v3.0.2 Pro版.zip","tmpfile/xuanxuan/2017/12/07/",651341,1512622303,"d2dfe785-edd6-4e13-a749-dffe9e8623ca"]}
							fileName, path, size, time, gid := req.Params.([]interface{})[0].(string),
								req.Params.([]interface{})[1].(string),
								req.Params.([]interface{})[2].(float64),
								req.Params.([]interface{})[3].(float64),
								req.Params.([]interface{})[4].(string)

							data = api.ChatUploadFile(module, method, fileName, path, size, time, gid, req.UserID)

						}
					case "chat.ping":
						{
							// 来自页面客户端或桌面客户端的 chat.ping 请求

							data = api.ChatPing(module, method, req.UserID)

						}
					default:
						{
							data["module"] = module
							data["method"] = method
							data["status"] = "fail"
							data["message"] = "unexpected entrypoint."
						}
					}

					if resp, err := api.EncrypFrom(data, true); err != nil {
						log.Fatal(err)
					} else {
						w.Write(resp)
					}
				}
			}

		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad Request!"))
		}

	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("pong"))
	})

	log.Printf("Starting Server with config file %s on port 4000", *conf)
	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal(err)
	}
}
