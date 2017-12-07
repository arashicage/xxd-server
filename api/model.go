package api

// 这个文件里的 method 直接操作数据库的数据

import (
	"fmt"
	"time"
	"xxd-server/model"

	"github.com/satori/go.uuid"
	"encoding/json"
	"strconv"
	"bytes"
	"github.com/kataras/go-errors"
)

// ResetUserStatus 重置全部用户的状态
// model.php resetUserStatus
func ResetUserStatus(status string) error {

	sql := "update `sys_user` set status = ?"
	if _, err := engine.Exec(sql, status); err != nil {
		return fmt.Errorf("[ResetUserStatus]: reset user status fail. %v", err)
	}
	return nil
}

// CreateSystemChat 创建一个 system chat
// model.php createSystemChat
func CreateSystemChat() error {

	// 首先检查 system chat 是否存在，不存在就重新创建一个
	chat := new(model.ImChat)
	if has, err := engine.Where("type = ?", "system").Get(chat); err != nil {
		// orm 报错
		return fmt.Errorf("[CreateSystemChat]: check system chat existence fail with orm fault. %v", err)
	} else if !has {
		// orm 没报错，但是没找到 system chat，则创建一个
		chat.Gid = fmt.Sprintf("%s", uuid.NewV4())
		chat.Name = "system group"
		chat.Type = "system"
		chat.Public = BoolAsString(true)
		chat.Createdby = "system"
		chat.Createddate = model.Time(time.Now())

		if _, err := engine.Insert(chat); err != nil {
			return fmt.Errorf("[CreateSystemChat]: system chat does not exists, insert(create) new system chat fail. %v", err)
		}
	} // else {} // 如果有，什么也不用做
	return nil
}

// GetUserByUserID 根据 id 获取 user
// model.php getUserByUserID
func GetUserByUserID(id int) (*model.SysUser, error) {

	user := new(model.SysUser)
	cols := []string{"id", "account", "realname", "avatar", "role", "dept", "status", "admin", "gender", "email", "mobile", "phone", "site"}

	has, err := engine.Cols(cols...).ID(id).Get(user)
	if err != nil {
		return nil, fmt.Errorf("[GetUserByUserID]: get user by id fail. %v", err)
	} else if !has {
		return nil, errors.New("[GetUserByUserID]: get user by id fail, no user found.")
	} else {
		return user, nil
	}

}

// GetUserList 根据 status, id 来获取 user list
// model.php getUserList idAsKey=false
func GetUserList(status string, id []int) ([]*model.SysUser, error) {

	cols := []string{"id", "account", "realname", "avatar", "role", "dept", "status", "admin", "gender", "email", "mobile", "phone", "site"}

	// 只在未(软)删除的用户里查找
	session := engine.Cols(cols...).Where("deleted = ?", "0")
	switch status {
	case "":
		{
			// 不附加条件
		}
	case "online":
		// 获取所有在线的 status != offline
		{
			session = session.And("status != ?", "offline")
		}
	default:
		// 获取指定 status 条件
		{
			session = session.And("status = ?", status)
		}
	}

	switch id {
	case nil:
		{
			// 不附加条件
		}
	default:
		{
			session = session.In("id", id)
		}
	}

	users := make([]*model.SysUser, 0)

	if err := session.Find(&users); err != nil {
		return nil, fmt.Errorf("[GetUserList]: get user list as slice fail. %v", err)
	}

	return users, nil

}

// GetUserList 根据 status, id 来获取 user list
// model.php getUserList idAsKey=true
func GetUserMap(status string, id []int) (map[int64]*model.SysUser, error) {

	cols := []string{"id", "account", "realname", "avatar", "role", "dept", "status", "admin", "gender", "email", "mobile", "phone", "site"}

	// 只在未(软)删除的用户里查找
	session := engine.Cols(cols...).Where("deleted = ?", "0")
	switch status {
	case "":
		{
			// 不附加条件
		}
	case "online":
		// 获取所有在线的 status != offline
		{
			session = session.And("status != ?", "offline")
		}
	default:
		// 获取指定 status 条件
		{
			session = session.And("status = ?", status)
		}
	}

	switch id {
	case nil:
		{
			// 不附加条件
		}
	default:
		{
			session = session.In("id", id)
		}
	}

	users := make(map[int64]*model.SysUser, 0)

	if err := session.Find(&users); err != nil {
		return nil, fmt.Errorf("[GetUserMap]: get user list as map fail. %v", err)
	}

	return users, nil

}

// EditUser 更新 user
// model.php editUser
func EditUser(user *model.SysUser) (*model.SysUser, error) {

	if user == nil {
		return nil, errors.New("[EditUser]: passed in a <nil> user.")
	} else {
		_, err := engine.Id(user.Id).Update(user)
		if err != nil {
			return nil, fmt.Errorf("[EditUser]: update user fail. %v", err)
		} else {
			return GetUserByUserID(user.Id)
		}
	}
}

// EditUser 更新 user, 从一个json 对象, UserChange 用户修改状态时, 传进的是有一个 json 对象
// model.php editUser
func EditUserFromJSON(userMap map[string]interface{}, userID int) (*model.SysUser, error) {

	if userMap == nil {
		return nil, errors.New("[EditUserFromJSON]: passed in a <nil> user.")
	} else {
		userId := userID
		_, err := engine.Table(new(model.SysUser)).Omit("id").ID(userId).Update(userMap)
		if err != nil {
			return nil, fmt.Errorf("[EditUserFromJSON]: update user fail. %v", err)
		} else {
			return GetUserByUserID(userId)
		}
	}
}

// GetChatByGID 根据 gid 获取 chat
// true 	nil 	err
// false 	nil 	err
// false 	xxx 	nil
// model.php getByGID 去掉了附加 members 的逻辑
func GetChatByGID(gid string) (bool, *model.ImChat, error) {

	chat := make([]*model.ImChat, 0)
	if err := engine.AllCols().Where("gid = ?", gid).Limit(1).Find(&chat); err != nil {
		return true, nil, fmt.Errorf("[GetChatByGID]: orm err: %v", err)
	} else if len(chat) == 0 {
		return false, nil, errors.New("[GetChatByGID]: chat not exist.")
	}

	return false, chat[0], nil

	/*
	chat := new(model.ImChat)

	has, err := engine.AllCols().Where("gid = ?", gid).Limit(1).Get(chat)
	return has, chat, nil
	*/

}

// GetAllSystemUsers 获取全部用户, 未删除的用户 deleted = "0"
func GetAllSystemUsers() ([]*model.SysUser, error) {

	xuser := make([]*model.SysUser, 0)
	if err := engine.Where("deleted = ?", "0").Find(&xuser); err != nil {
		return nil, fmt.Errorf("[GetAllSystemUser]: get system user fail %v", err)
	} else {
		return xuser, nil
	}
}

// GetChatMemberByGID 获取 chat 参与者 id
// model.php getMemberListByGID
func GetChatMemberByGID(gid string) ([]int, error) {

	members := make([]int, 0)

	chat := new(model.ImChat)
	chat.Gid = gid
	if has, err := engine.Get(chat); !has || err != nil {
		return nil, fmt.Errorf("[GetChatMemberByGID]: chat not exist or %v", err)
	}

	if chat.Type == "system" {
		// 获取全部用户
		user, err := GetAllSystemUsers()
		if err != nil {
			return nil, fmt.Errorf("[GetChatMemberByGID]: %v", err)
		}

		for _, user := range user {
			members = append(members, user.Id)
		}
	} else {
		chatUser := make([]*model.ImChatuser, 0)
		session := engine.Cols("user").Where("quit = ?", "0000-00-00 00:00:00")
		if gid != "" {
			session.And("cgid = ?", gid)
		}

		if err := session.Find(&chatUser); err != nil {
			return nil, fmt.Errorf("[GetChatMemberByGID]: get chat user fail %v", err)
		}

		for _, user := range chatUser {
			members = append(members, user.User)
		}
	}

	return members, nil

}

// 获取 chat 的在线参与者
func GetChatOnlineMemberByGID(gid string) ([]int, error) {

	members, err := GetChatMemberByGID(gid)
	if err != nil {
		return nil, fmt.Errorf("[GetChatOnlineMemberByGID] %v", err)
	}

	users, err := GetUserList("online", members)
	if err != nil {
		return nil, fmt.Errorf("[GetChatOnlineMemberByGID] %v", err)
	}

	onLines := make([]int, 0)
	for _, user := range users {
		if user.Status == "online" {
			onLines = append(onLines, user.Id)
		}
	}

	return onLines, nil
}

// CreateMessage 创建 message
// model.php createMessage
func CreateMessage(messages []model.ImMessage, userID int) ([]*model.ImMessage, error) {

	// message.id list, chat.gid list
	messageIdList, chatGidList := make([]int, 0), make([]string, 0)

	for _, m := range messages {
		message := new(model.ImMessage)
		has, err := engine.AllCols().Where("gid = ?", m.Gid).Get(message)
		if err != nil {
			return nil, nil
		} else if has {
			// 有则更新
			if message.Contenttype == "image" || message.Contenttype == "file" {
				sql := "update `im_message` set content = ? where id = ?"
				_, err := engine.Exec(sql, message.Content, message.Id)
				if err != nil {
					break
				}
				messageIdList = append(messageIdList, message.Id)
			}
		} else {
			// 无则插入
			if m.User == "" {
				m.User = strconv.Itoa(userID)
			}
			if time.Time(m.Date).IsZero() {
				m.Date = model.Time(time.Now())
			}

			// 传入的内容 id 为空，此处忽略让 db 自己填
			_, err := engine.Omit("id").Insert(&m)
			if err != nil {
				break
			}
			messageIdList = append(messageIdList, m.Id)
		}
		chatGidList = append(chatGidList, m.Cgid)
	}

	if 0 == len(messageIdList) {
		return nil, nil
	}

	sql := "update `im_chat` set lastActiveTime = now() where gid = ?"
	var err error
	for _, gid := range chatGidList {
		if _, err := engine.Exec(sql, gid); err != nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("[CreateMessage] %v", err)
	}

	return GetMessageListByID(messageIdList)

}

// SaveOfflineMessages 存储离线消息
// model.php saveOfflineMessages
func SaveOfflineMessages(messages []*model.ImMessage, users []int) (err error) {

	// 解决 Date 字段转 时间戳(秒数), 使用 model.Time
	// 解决 cgid 中 html escape 问题，使用下面的方法
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(messages); err != nil {
		return fmt.Errorf("[SaveOfflineMessages] ecodeing messages fail. %v", err)
	}

	for _, user := range users {
		userMessage := new(model.ImUsermessage)

		userMessage.User = user
		userMessage.Message = buffer.String()

		_, err = engine.Insert(userMessage)
		if err != nil {
			break
		}
	}
	return err

}

// GetMessageListByID 根据 messenge 的 id 列表获取 messages
// model.php getMessageList
func GetMessageListByID(id []int) ([]*model.ImMessage, error) {

	messages := make([]*model.ImMessage, 0)

	session := engine.AllCols().Where("1 = 1")
	if id != nil {
		session.In("id", id)
	}
	if err := session.Desc("id").Find(&messages); err != nil {
		return nil, fmt.Errorf("[GetMessageListByID]: get message list fail. %v", err)
	}
	return messages, nil
}

// GetMessageListByCGID 根据 messenge 的 cgid 获取 messages
func GetMessageListByCGID(cgid string) ([]*model.ImMessage, error) {

	messages := make([]*model.ImMessage, 0)
	if err := engine.AllCols().Where("cgid = ?", cgid).Desc("id").Find(&messages); err != nil {
		return nil, fmt.Errorf("[GetMessageListByCGID]: get message list by cgid fail. %v", err)
	}

	return messages, nil
}

// GetOfflineMessages 获取某个用户的离线信息 JSON
// model.php getOfflineMessages
func GetOfflineMessages(userID int) ([]map[string]interface{}, error) {

	offlineMessages := make([]map[string]interface{}, 0)

	messages := make([]*model.ImUsermessage, 0)

	if err := engine.AllCols().Where("user = ?", userID).OrderBy("level, id").Find(&messages); err != nil {
		return nil, fmt.Errorf("[GetOfflineMessages] get offline messages fail. %v", messages)
	}

	for _, message := range messages {

		jsonMessages := make([]map[string]interface{}, 0)
		if err := json.Unmarshal([]byte(message.Message), &jsonMessages); err != nil {
			return nil, fmt.Errorf("[GetOfflineMessages] assemble offline messages fail. %v", messages)
		}

		offlineMessages = append(offlineMessages, jsonMessages...)

	}

	// 取出后从数据库删除
	sql := "delete from `im_usermessage` where user = ?"
	if _, err := engine.Exec(sql, userID); err != nil {
		return nil, fmt.Errorf("[GetOfflineMessages]: delete user offline messages fail. %v", err)
	}

	/*
	userMessage := new(model.ImUsermessage)
	userMessage.User = userID
	if _, err := engine.Delete(userMessage); err != nil {
		return nil, fmt.Errorf("[GetOfflineMessages]: delete user offline messages fail. %v", err)
	}
	*/

	return offlineMessages, nil

}

// GetChatList 即 GetList 获取指定 public 状态的 chat
// model.php getList
func GetChatList(public string) ([]*model.ImChat, error) {

	xChat := make([]*model.ImChat, 0)

	if err := engine.AllCols().Where("public = ?", public).Find(&xChat); err != nil {
		return nil, fmt.Errorf("[GetList]: get chat list fail. %v", err)
	}

	return xChat, nil

}

type ChatList struct {
	Chat *model.ImChat `xorm:"extends"`
	Star string
	Hide string
	Mute string
}

func (ChatList) TableName() string {
	return "im_chat"
}

// GetSystemChat 获取 system chat
func GetSystemChat() (*model.ImChat, error) {

	chat := make([]*model.ImChat, 0)
	if err := engine.AllCols().Where("type = ?", "system").Limit(1).Find(&chat); err != nil {
		return nil, fmt.Errorf("[GetSystemChat]: get systemchat fail. %v", err)
	}

	return chat[0], nil
}

// GetChatListByUserID 获取 某个用户参与的 chat 列表，含 system chat，附加 hide star mute
// model.php getListByUserID
func GetChatListByUserID(userID string, star bool) ([]*ChatList, error) {

	systemChat, err := GetSystemChat()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	chatList := make([]*ChatList, 0)

	session := engine.Cols("im_chat.*, im_chatuser.star, im_chatuser.hide, im_chatuser.mute").Join("LEFT", "im_chatuser", "im_chat.gid = im_chatuser.cgid").Where("im_chatuser.quit = ?", "0000-00-00 00:00:00").And("im_chatuser.user = ?", userID)
	if star {
		session.And("im_chatuser.star = ?", "1")
	} else {
		// 不附加条件
		//session.And("im_chatuser.star = ?", "0")
	}

	err = session.Find(&chatList)
	if err != nil {
		return nil, fmt.Errorf("get list by user id fail %v", err)
	}

	chatList = append(chatList, &ChatList{systemChat, "0", "0", "0"})

	return chatList, nil

}

// 创建 chat，比如 chat 被检查到不存在时
// model.php create
func Create(params []interface{}, userID int) (*model.ImChat, error) {

	user, err := GetUserByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("[Create] %v", err)
	}

	chat := new(model.ImChat)
	chat.Gid = params[0].(string)
	chat.Name = params[1].(string)
	chat.Type = params[2].(string)
	chat.Subject = int(params[4].(float64))
	if params[5].(bool) == true {
		chat.Public = "1"
	} else {
		chat.Public = "0"
	}
	if user.Account != "" {
		chat.Createdby = user.Account
	} else {
		chat.Createdby = ""
	}
	chat.Createddate = model.Time(time.Now())

	_, err = engine.Insert(chat)
	if err != nil {
		return nil, fmt.Errorf("[Create]: insert chat fail. %v", err)
	}

	// 1.2 版的桌面客户端，组织 参数时是错误的,类型混杂 "1","2","3",2;也观察到全是 int的情况。 browser 版是正确的
	for _, id := range params[3].([]interface{}) {
		var userID int
		switch id.(type) {
		case float64:
			userID = int(id.(float64))
		case string:
			userID, _ = strconv.Atoi(id.(string))
		}
		JoinChat(chat.Gid, userID, true)
	}

	_, chat, err = GetChatByGID(chat.Gid)
	return chat, err

}

// UpdateChat 更新一个 chat 的某些字段
// model.php update
func UpdateChat(chat *model.ImChat, userID int, cols []string) (*model.ImChat, error) {

	if chat != nil {
		user, err := GetUserByUserID(userID)
		if err != nil {
			return nil, fmt.Errorf("[UpdateChat] %v", err)
		} else {
			chat.Editedby = user.Account
			chat.Editeddate = model.Time(time.Now())

			_, err := engine.Cols(cols...).Where("gid = ?", chat.Gid).Update(chat)
			if err != nil {
				return nil, fmt.Errorf("[UpdateChat] update chat fail %v", err)
			} else {
				_, chat, err := GetChatByGID(chat.Gid)
				return chat, err
			}
		}
	} else {
		return nil, fmt.Errorf("[UpdateChat] update fail, passed in with a nil. ")
	}
}

// StarChat 收藏或取消 chat, 更新 star 字段
// model.php starChat
func StarChat(gid string, star string, userid int) (error) {

	sql := "update `im_chatuser` set star = ? where cgid = ? and user = ?"
	_, err := engine.Exec(sql, star, gid, userid)
	if err != nil {
		return fmt.Errorf("[StarChat] star chat fail. %v", err)
	} else {
		return nil
	}

}

// HideChat 隐藏或显示 chat, 更新 hide 字段
// model.php hideChat
func HideChat(gid string, hide string, userID int) error {

	sql := "update `im_chatuser` set hide = ? where cgid = ? and user = ?"
	_, err := engine.Exec(sql, hide, gid, userID)
	if err != nil {
		return fmt.Errorf("[HideChat] hide chat fail. %v", err)
	} else {
		return nil
	}
}

// JoinChat 加入或退出 chat
// model.php joinChat
func JoinChat(gid string, userID int, join bool) error {

	chatUser := new(model.ImChatuser)
	xChatUser := make([]*model.ImChatuser, 0)
	if join {
		/* Join chat. */
		// 检查用户是否参与了该 chat
		if err := engine.AllCols().Where("cgid = ?", gid).And("user = ?", userID).Limit(1).Find(&xChatUser); err != nil {
			return err
		} else if len(xChatUser) != 0 {
			// user found
			if chatUser.Quit.Format("2006-01-02 15:04:05") == "0000-00-00 00:00:00" {
				// user exists and not quit
				return nil
			} else {
				chatUser = xChatUser[0]
				// 如果用户已经退出了，则更新该用户
				chatUser.Join = time.Now()
				q, _ := time.Parse("2006-01-02 15:04:05", "0000-00-00 00:00:00")
				chatUser.Quit = q
				_, err := engine.Cols("join", "quit").Where("cgid = ?", gid).And("user = ?", userID).Update(chatUser)
				if err != nil {
					return err
				} else {
					return nil
				}
			}
		} else {
			// user not found
			// 该讨论组里没有该有该用户
			chatUser.Cgid = gid
			chatUser.User = userID
			chatUser.Join = time.Now()
			chatUser.Star = "0"
			chatUser.Hide = "0"
			chatUser.Mute = "0"

			_, err := engine.Insert(chatUser)
			if err != nil {
				return err
			}

			// 设置 order = id
			sql := "update `im_chatuser` set `order` = ? where id = ?"
			if _, err := engine.Exec(sql, chatUser.Id, chatUser.Id); err != nil {
				fmt.Println(err)
				return err
			} else {
				return nil
			}

		}
	} else {
		/* Quit chat. */
		sql := "update `im_chatuser` set quit = ? where cgid = ? and user = ?"
		if _, err := engine.Exec(sql, time.Now(), gid, userID); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func DoSettings(userID string, settings string) (string, error) {

	users := make([]*model.SysUser, 0)
	if err := engine.Where("account = ?", userID).Limit(1).Find(&users); err != nil {
		return "", fmt.Errorf("[DoSettings] %v", err)
	}

	user := users[0]
	if settings == "" {
		// 取
		return user.Settings, nil
	} else {
		// 存
		if _, err := engine.Cols("settings").Where("account = ?", userID).Update(user); err != nil {
			return "", fmt.Errorf("[DoSettings] %v", err)
		} else {
			return settings, nil
		}
	}

}
