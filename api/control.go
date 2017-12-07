package api

// 这一层为中间层，间接调用 model.go，组织返回的消息

import (
	"fmt"
	"net/http"
	"xxd-server/model"
	"xxd-server/util"

	"strings"
	"strconv"
	"github.com/mitchellh/mapstructure"
	"time"
)

type RequestData struct {
	UserID int         `json:"userid"` // 用户id，xxd -> rzs 非登录时必须
	Module string      `json:"module"` // 模块名称,必须
	Method string      `json:"method"` // 方法名称,必须
	Test   bool        `json:"test"`   // 可选参数，bool,默认为false。
	Params interface{} `json:"params"` // 参数对象,可选
	Data   interface{} `json:"data"`   // 请求数据,可选,与params配合使用,通常data传输是对象
}

// 启动 chat server
func ChatServerStart(module string, method string) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["status"] = http.StatusInternalServerError
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	if err := ResetUserStatus("offline"); err != nil {
		PanicWith("reset user status to offline fail.", err)
	} else if err := CreateSystemChat(); err != nil {
		PanicWith("create system chat fail.", err)
	} else {
		data["status"] = http.StatusOK
	}

	return data
}

// Login 入口
func ChatLogin(module string, method string, account string, password string, status string) (data map[string]interface{}) {

	// 来自 客户端的请求，status = ""
	if status == "" {
		return VerifyLogin(module, method, account, password)
	} else {
		return Login(module, method, account, password, status)
	}
}

// 验证密码
func VerifyLogin(module string, method string, account string, password string) (data map[string]interface{}) {

	fmt.Println(`Request: >>> {"method":"verifylogin","module":"chat"}`)

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	md5 := util.GetMD5(password + account)

	user := new(model.SysUser)
	if has, err := engine.Where("account = ?", account).Get(user); err != nil {
		PanicWith("query user from db fail.", err)
		return data
	} else if !has {
		PanicWith("user not exists.")
		return data
	}

	pass := md5 == user.Password

	if pass {
		// 用户密码验证通过
		data["result"] = "success"

	} else {
		// 用户密码验证未通过
		data["result"] = "fail"
		data["data"] = "Login failed. Check you account and password."
	}

	return data

}

// Login
// 返回当前在线用户数据, 内容包含在 model.SysUser 中，没有 time 类型字段
/*
{
    module: 'chat',
    method: 'login',
    result,
    users[]，
    data:
    {             // 当前在线的用户数据
        id,       // ID
        account,  // 用户名
        realname, // 真实姓名
        avatar,   // 头像URL
        role,     // 角色
        dept,     // 部门ID
        status,   // 当前状态
        admin,    // 是否超级管理员，super 超级管理员 | no 普通用户
        gender,   // 性别，u 未知 | f 女 | m 男
        email,    // 邮箱
        mobile,   // 手机
        site,     // 网站
        phone     // 电话
    }
}
*/
func Login(module string, method string, account string, password string, status string) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	md5 := util.GetMD5(password + account)

	user := new(model.SysUser)
	if has, err := engine.Where("account = ?", account).Get(user); err != nil {
		PanicWith("query user from db fail.", err)
		return data
	} else if !has {
		PanicWith("user not exists.")
		return data
	}
	pass := md5 == user.Password

	/*
	if user.Status == "online" {
		//	kickoff
	}
	*/

	if pass {
		// 用户密码验证通过
		data["result"] = "success"

		if status == "online" {

			// user with status
			// 更新数据库库和内存中的 user 状态
			userOnline := new(model.SysUser)
			userOnline.Id = user.Id
			userOnline.Status = status
			user, err := EditUser(userOnline)
			if err != nil {
				PanicWith(err)
				return data
			}

			onLines, err := GetUserList("online", nil)
			if err != nil {
				PanicWith("user not exists or Get user fail.", err)
				return data
			}

			onLinesID := make([]int, 0)
			for _, user := range onLines {
				onLinesID = append(onLinesID, user.Id)
			}

			data["users"] = onLinesID
			data["data"] = user

		}

	} else {
		// 用户密码验证未通过
		data["result"] = "fail"
		data["message"] = "Login failed. Check you account and password."
	}

	return data
}

// ChatLogout
func ChatLogout(module string, method string, userID int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	user := new(model.SysUser)
	user.Id, user.Status = userID, "offline"

	user, err := EditUser(user)
	if err != nil {
		PanicWith(err)
	}

	onLines, err := GetUserList("online", nil)
	if err != nil {
		PanicWith(err)
	}

	data["result"] = "success"
	data["users"] = []int{userID}
	data["data"] = onLines

	return data

}

// 获取所有用户列表，不附加 status, id 条件
func ChatUserGetlist(module string, method string, userid int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	// 获取所有用户列表，不附加 status, id 条件
	users, err := GetUserList("", nil)
	if err != nil {
		PanicWith("Call GetUserList in ChatUserGetlist fail.", err)
	} else {
		data["result"] = "success"
		data["users"] = []int{userid}
		data["data"] = users
	}

	return data
}

func ChatGetList(module string, method string, mm *RequestData) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	userid := fmt.Sprintf("%d", mm.UserID)
	chatList, err := GetChatListByUserID(userid, false)
	if err != nil {
		PanicWith("GetListByUserID failed.", err)
		return data
	}

	type RespondData struct {
		*model.ImChat
		Star    string `json:"star"`
		Hide    string `json:"hide"`
		Mute    string `json:"-" field:"mute"`
		Members []int  `json:"members"`
	}

	dataSlice := make([]RespondData, 0)

	for _, chat := range chatList {

		item := RespondData{}
		item.ImChat = chat.Chat
		item.Star = chat.Star
		item.Hide = chat.Hide
		item.Mute = chat.Mute

		members, err := GetChatMemberByGID(chat.Chat.Gid)

		if err != nil {
			PanicWith("GetMemberListByGID failed.", err)
			return data
		}

		item.Members = members

		dataSlice = append(dataSlice, item)
	}

	data["result"] = "success"
	data["users"] = []int{mm.UserID}
	data["data"] = dataSlice

	return data
}

/*
		$messages = $this->chat->getOfflineMessages($userID);
        if(dao::isError())
        {
            $this->output->result  = 'fail';
            $this->output->message = 'Get offline messages fail.';
        }
        else
        {
            $this->output->result = 'success';
            $this->output->users  = array($userID);
            $this->output->data   = $messages;
        }
        $this->output->method = 'message';
        die($this->app->encrypt($this->output));
*/
func ChatGetOfflineMessages(module string, method string, userid int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = "message"

	offlineMessages, err := GetOfflineMessages(userid)
	if err != nil {
		PanicWith("Get offline messages fail.", err)
		return data
	} else {
		data["result"] = "success"
		data["users"] = []int{userid}
		data["data"] = offlineMessages
	}

	return data
}

/*
public function settings($account = '', $settings = '', $userID = 0)
{
	if($settings)
	{
		$this->loadModel('setting')->setItem("system.sys.chat.settings.$account", helper::jsonEncode($settings));
	}

	if(dao::isError())
	{
		$this->output->result  = 'fail';
		$this->output->message = 'Save settings failed.';
	}
	else
	{
		$this->output->result = 'success';
		$this->output->users  = array($userID);
		$this->output->data   = !empty($settings) ? $settings : json_decode($this->config->chat->settings->$account);
	}

	die($this->app->encrypt($this->output));
}
*/

func ChatSettings(module string, method string, userid int, account string, settings string) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	//s, _ := DoSettings(account, settings)

	data["users"] = []int{userid}
	data["result"] = "success"
	data["data"] = ""

	return data
}

/*
public function ping($userID = 0)
{
	$this->output->result = 'success';
	$this->output->users  = array($userID);

	die($this->app->encrypt($this->output));
}
*/
func ChatPing(module string, method string, userid int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	data["result"] = "success"
	data["users"] = []int{userid}

	return data

}

func ChatUserChange(module string, method string, userMap map[string]interface{}, userID int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	// {"module":"chat","method":"userChange","params":[{"status":"away"}],"userID":3}

	user, err := EditUserFromJSON(userMap, userID)
	if err != nil {
		fmt.Println(err)
		PanicWith(err)
		return data
	}

	onLines, err := GetUserList("online", nil)
	if err != nil {
		fmt.Println(err)
		PanicWith(err)
		return data
	}

	onLinesID := make([]int, 0)
	for _, user := range onLines {
		onLinesID = append(onLinesID, user.Id)
	}

	data["users"] = onLinesID
	data["data"] = user

	return data
}

func ChatMessage(module string, method string, messages map[string]interface{}, userid int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	// {"module":"chat","method":"message","params":{"messages":[{"gid":"7af873cf-8db7-4663-9779-0edce57afc6d","cgid":"2&3","type":"normal","contentType":"text","content":"online message 1","date":"","user":2}]},"userID":2}

	// params 的实际构造，decode 出来
	xMessages := struct {
		Messages []model.ImMessage `mapstructure:"messages"`
	}{}
	mapstructure.Decode(messages, &xMessages) // decode 出来的 model.ImMessage 里 id 是 int 0

	// 检查每条消息的 cgid 是否有重复，重复的只记一次。消息内容涉及多个cgid 时返回失败。
	chats := map[interface{}]bool{}
	for _, m := range xMessages.Messages {
		chats[m.Cgid] = true // 相同的 cgid 只记一次
	}
	if len(chats) > 1 {
		data["result"] = "fail"
		data["data"] = "Messages belong to different chats."
		return data
	}

	// 原 php 代码中只取了 current，即第一个
	xerr := make([]map[string]interface{}, 0)
	message := xMessages.Messages[0]
	_, chat, err := GetChatByGID(message.Cgid)
	if err != nil {
		errm := make(map[string]interface{})
		errm["gid"] = message.Cgid
		errm["message"] = err
		xerr = append(xerr, errm)
	} else if chat == nil {
		errm := make(map[string]interface{})
		errm["gid"] = message.Cgid
		errm["message"] = "Chat do not exist."
		xerr = append(xerr, errm)
	} else if chat.Admins != "" {
		admins := strings.Split(chat.Admins, ",")
		if !contains(admins, strconv.Itoa(userid)) {
			errm := make(map[string]interface{})
			errm["gid"] = message.Cgid
			errm["message"] = "No rights to chat."
			xerr = append(xerr, errm)
		}
	}
	// 检查上面 xerr 内容，如果发生了错误就返回
	if len(xerr) != 0 {
		data["result"] = "fail"
		data["data"] = err
		return data
	}

	members, err := GetChatMemberByGID(message.Cgid)
	if err != nil {
		PanicWith(err)
	}

	onLines, offLines := []int{userid}, make([]int, 0)
	users, err := GetUserList("", members)
	if err != nil {
		PanicWith("Send message failed.")
	} else {
		for _, user := range users {
			if user.Id == userid {
				continue
			} else if user.Status == "offline" {
				offLines = append(offLines, user.Id)
			} else {
				onLines = append(onLines, user.Id)
			}
		}
	}

	messageSend, err := CreateMessage(xMessages.Messages, userid)
	if err != nil {
		PanicWith(err)
	}
	SaveOfflineMessages(messageSend, offLines)

	data["result"] = "success"
	data["users"] = onLines
	data["data"] = messageSend

	return data
}

func ChatJoinChat(module string, method string, gid string, userid int, join bool) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
			fmt.Println(data)
		}
	}()

	data = make(map[string]interface{}, 0)
	data["module"] = module
	data["method"] = method

	_, chat, err := GetChatByGID(gid)
	if err != nil {
		// orm 报错 或 chat 不存在
		PanicWith(err)
	} else if chat.Type != "group" && false { // todo  system one2one group 必须是group？
		PanicWith("It is not a group chat.")
	} else if join && chat.Public == "0" {
		PanicWith("It is not a public chat.")
	}

	// join == true 加入会话; join == false 退出会话
	JoinChat(gid, userid, join)

	// todo 获取chat 的 members，和chat 组合在一起
	// 获取memeber 里 online 的部分
	members, err := GetChatMemberByGID(chat.Gid)
	if err != nil {
		PanicWith(err)
	}

	type RespondData struct {
		*model.ImChat
		Members []int `json:"members"`
	}

	crd := RespondData{
		ImChat:  chat,
		Members: members,
	}

	data["result"] = "success"
	data["data"] = crd //members

	return data
}

func ChatCreate(module string, method string, params []interface{}, userid int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	/*
	0	gid,     					// 会话的全局id,
	1	name,    					// 会话的名称
	2	type,    					// 会话的类型
	3	members: [{id}, {id}...] 	// 会话的成员列表
	4	subject, 					//可选,主题会话的关联主题ID,默认为0
	5	pulic    					//可选,是否公共会话,默认为false
	*/

	x, chat, err := GetChatByGID(params[0].(string))
	if x && err != nil {
		// orm 报错
		fmt.Println(err)
		PanicWith("Create chat fail.")
	} else if !x && err != nil {
		// orm 无错，err 非空，无数据。创建 chat，填充 chatuser
		chat, err = Create(params, userid)
		if err != nil {
			fmt.Println(err)
			PanicWith("Create chat fail.")
		}
	}
	// chat 已存在
	members, err := GetChatMemberByGID(chat.Gid)
	if err != nil {
		fmt.Println(err)
		PanicWith("Create chat fail.")
	}

	type RespondData struct {
		*model.ImChat
		Members []int `json:"members"`
	}

	crd := RespondData{
		ImChat:  chat,
		Members: members,
	}

	data["result"] = "success"
	data["users"] = []int{userid}
	data["data"] = crd

	return data
}

func ChatChangePublic(module string, method string, params []interface{}, userid int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	// {"module":"chat","method":"changePublic","params":["7adacd1c-01ac-444d-aef3-5daf22401fed",true],"userID":2}

	gid, public := params[0].(string), params[1].(bool)

	_, chat, err := GetChatByGID(gid)
	if err != nil {
		fmt.Println(err)
		PanicWith(err)
	}

	// exist
	if chat.Type != "group" {
		PanicWith("It is not a group chat.")
	}

	chat.Public = BoolAsString(public)

	chat, err = UpdateChat(chat, userid, []string{"public"})
	if err != nil {
		fmt.Println(err)
		PanicWith("update chat fail.")
	}

	onLines, err := GetChatOnlineMemberByGID(gid)
	if err != nil {
		fmt.Println(err)
		PanicWith("Create chat fail.")
	}

	type RespondData struct {
		*model.ImChat
		Members []int `json:"members"`
	}

	crd := RespondData{
		ImChat:  chat,
		Members: onLines,
	}

	data["result"] = "success"
	data["users"] = onLines
	data["data"] = crd

	return data
}

func ChatChangeName(module string, method string, gid string, name string, userid int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	// {"module":"chat","method":"changename","params":["021d859c-85b3-43b5-902e-7bd5b192b376","123456"],"userID":2}

	_, chat, err := GetChatByGID(gid)
	if err != nil {
		fmt.Println(err)
		PanicWith(err)
	}

	// exist
	if chat.Type != "group" {
		PanicWith("It is not a group chat.")
	}

	chat.Name = name

	chat, err = UpdateChat(chat, userid, []string{"name"})
	if err != nil {
		fmt.Println(err)
		PanicWith("update chat fail.")
	}

	onLines, err := GetChatOnlineMemberByGID(gid)
	if err != nil {
		fmt.Println(err)
		PanicWith("Create chat fail.")
	}

	type RespondData struct {
		*model.ImChat
		Members []int `json:"members"`
	}

	crd := RespondData{
		ImChat:  chat,
		Members: onLines,
	}

	data["result"] = "success"
	data["users"] = onLines
	data["data"] = crd

	return data
}

// ChatAddMember 添加或移除 chat 成员
func ChatAddMember(module string, method string, gid string, newMembers []int, join bool, userID int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("Get public chat list failed. %v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	_, chat, err := GetChatByGID(gid)
	if err != nil {
		fmt.Println(err)
		PanicWith(err)
	}

	// exist
	if chat.Type != "group" {
		PanicWith("It is not a group chat.")
	}

	for _, member := range newMembers {
		JoinChat(gid, member, join)
	}

	members, err := GetChatMemberByGID(chat.Gid)
	if err != nil {
		PanicWith(err)
	}

	type RespondData struct {
		*model.ImChat
		Members []int `json:"members"`
	}

	crd := RespondData{
		ImChat:  chat,
		Members: members,
	}

	data["result"] = "success"
	data["data"] = crd

	return data

}

// ChatMembers 获取会话的成员
func ChatMembers(module string, method string, gid string, userid int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("Get member list failed. %v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	members, err := GetChatMemberByGID(gid)
	if err != nil {
		PanicWith(err)
	}

	data["result"] = "success"
	data["users"] = []int{userid}
	data["data"] = members

	return data
}

// ChatGetPublicList 获取公共的 chatList
func ChatGetPublicList(module string, method string, userid int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("Get public chat list failed. %v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	chatList, err := GetChatList(BoolAsString(true))
	if err != nil {
		fmt.Println(err)
		PanicWith(err)
	}

	data["result"] = "success"
	data["users"] = []int{userid}
	data["data"] = chatList

	return data
}

func ChatStar(module string, method string, gid string, star bool, userID int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	if err := StarChat(gid, BoolAsString(star), userID); err != nil {
		if star {
			PanicWith("Star chat failed.")
		} else {
			PanicWith("Cancel star chat failed.")
		}
	}

	type RespondData struct {
		Gid  string `json:"gid"`
		Star bool   `json:"star"`
	}

	rd := RespondData{
		Gid:  gid,
		Star: star,
	}

	data["result"] = "success"
	data["users"] = []int{userID}
	data["data"] = rd

	return data
}

func ChatHide(module string, method string, gid string, star bool, userID int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	if err := HideChat(gid, BoolAsString(star), userID); err != nil {
		if star {
			PanicWith("Star chat failed.")
		} else {
			PanicWith("Cancel star chat failed.")
		}
	}

	type RespondData struct {
		Gid  string `json:"gid"`
		Hide bool   `json:"star"`
	}

	rd := RespondData{
		Gid:  gid,
		Hide: star,
	}

	data["result"] = "success"
	data["users"] = []int{userID}
	data["data"] = rd

	return data
}

func ChatSetCommiters(module string, method string, gid string, committers string, userID int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("Set committers failed. %v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	_, chat, err := GetChatByGID(gid)
	if err != nil {
		fmt.Println(err)
		PanicWith(err)
	}

	// exist
	if chat.Type != "group" {
		PanicWith("It is not a group chat.")
	}

	chat.Committers = committers

	chat, err = UpdateChat(chat, userID, []string{"committers"})
	if err != nil {
		fmt.Println(err)
		PanicWith(err)
	}

	onLines, err := GetChatOnlineMemberByGID(gid)
	if err != nil {
		PanicWith(err)
	}

	type RespondData struct {
		*model.ImChat
		Members []int `json:"members"`
	}

	crd := RespondData{
		ImChat:  chat,
		Members: onLines,
	}

	data["result"] = "success"
	data["users"] = onLines
	data["data"] = crd

	return data

}

//
func ChatHistory(module string, method string, gid string, recPerPage int, pageID int, recTotal int, continued bool, userID int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("Set committers failed. %v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	// {"module":"chat","method":"history","params":["2&4",50,1,0,true],"userID":2}

	messageList, err := GetMessageListByCGID(gid)
	if err != nil {
		fmt.Println(err)
		PanicWith("Get history failed.")
	}

	pager := make(map[string]interface{}, 0)
	if len(messageList) == 0 {
		pager["recPerPage"] = recPerPage
		pager["pageID"] = pageID
		pager["recTotal"] = len(messageList)
		pager["gid"] = gid
		pager["continued"] = false
	} else {
		pager["recPerPage"] = recPerPage
		pager["pageID"] = pageID
		pager["recTotal"] = len(messageList)
		pager["gid"] = gid
		pager["continued"] = recPerPage*pageID < len(messageList)
	}

	data["result"] = "success"
	data["users"] = []int{userID}
	data["data"] = messageList
	data["pager"] = pager

	return data

}

func ChatUploadFile(module string, method string, fileName string, path string, size float64, timestamp float64, gid string, userID int) (data map[string]interface{}) {

	defer func() {
		if rec := recover(); rec != nil {
			data["result"] = "fail"
			data["message"] = fmt.Sprintf("%v", rec)
		}
	}()

	data = make(map[string]interface{})
	data["module"] = module
	data["method"] = method

	_, chat, err := GetChatByGID(gid)
	if chat == nil {
		fmt.Println(err)
		PanicWith("Chat do not exist.")
	}

	user, err := GetUserByUserID(userID)
	if err != nil {
		fmt.Println(err)
		PanicWith("Upload file failed.")
	}

	onLines, err := GetChatOnlineMemberByGID(gid)
	if err != nil {
		fmt.Println(err)
		PanicWith("Upload file failed.")
	}

	file := new(model.SysFile)
	file.Pathname = path
	xFile := strings.Split(fileName, ".")
	file.Title = strings.Join(xFile[:len(xFile)-1], ".")
	file.Extension = xFile[len(xFile)-1]
	file.Size = int(size)
	file.Objecttype = "chat"
	file.Objectid = chat.Id
	file.Createdby = user.Account
	file.Createddate = time.Unix(int64(timestamp), 0)

	_, err = engine.Insert(file)
	if err != nil {
		fmt.Println(err)
		PanicWith("Upload file failed.")
	}

	fileID := file.Id
	pathx := path + util.GetMD5(fmt.Sprintf("%s%d%d", fileName, fileID, timestamp))

	sql := "update `sys_file` set pathname = ? where id = ?"
	if _, err := engine.Exec(sql, pathx, fileID); err != nil {
		fmt.Println(err)
		PanicWith("Upload file failed.")
	}

	data["result"] = "success"
	data["users"] = onLines
	data["data"] = fmt.Sprintf("%d", fileID)

	return data

}
