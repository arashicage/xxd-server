D:\Workspace\workspace.env\nginx-1.13.1
D:\Dropbox\workspace\workspace.go\gocode\src\xxd-server
D:\Dropbox\workspace\workspace.go\gocode\src\xxd-server

setTimeout(()=>{this.markClose(),this.send("logout") 要改成 setTimeout(()=>{this.send("logout")


//xorm cmd tool
//xorm reverse mysql "root:123456@tcp(172.30.11.230:3306)/xxd?charset=utf8" d:\Dropbox\workspace\workspace.go\gopath\src\github.com\go-xorm\cmd\xorm\templates/goxorm



系统组里的任意一条消息，都会发给组内的每一个人。 usermsesage.message 的内容是一条信息，不是多条。

1对1 的离线消息，只发给对方。

im_message 也是存单条信息。



在不创建group 的时候，直接和 xx 直接开聊，发消息， 会 create  gid:xx&yy  name:chaNaN type:one2one 的 chat。

选定xx1， 添加 xx2， xx3 时，创建一个 gid:uuid  name:ccc type:group.
                     随后对xx1组员，创建一个 常规组


create
1. 根据gid查 chat 含 member
2. 如果chat 不存在，调用create
3. 取出 users = memeber 中 online 的 users
4. data= chat， users = users id


xxd 启动	                                    OK
登陆	                                        OK
注销	                                        OK- 	browser 端注销时未见调用*1
重复登陆	                                    OK- 	xxd 内部实现
获取所有用户列表	                            OK
获取当前登录用户参与的所有讨论组	                OK
获取当前登录用户所有离线消息	                    OK
更改当前登录用户的信息	                        OK
创建聊天会话	                                OK
加入或退出聊天会话	                            OK	退出后其他人右侧的人员列表没变化，browser 前端有问题
更改会话名称	                                OK
收藏或取消收藏会话	                            OK
邀请新的用户到回话	                            OK	其他人右侧的人员列表没，browser 前端有问题
踢出用户	                                    OK-	前端无入口
向会话发送消息	                                OK
获取会话的所有消息记录
获取会话的所有成员信息	                        OK
隐藏或显示会话	                                OK-	前端无入口
设置公共会话	                                OK	从 + 处可以看到
设置会话管理员	                                OK-	前端无入口
设置会话发言人	                                OK
上传下载setting
上传文件

*1 browser-client dist\budle.js setTimeout(()=>{this.markClose(),this.send("logout") 要改成 setTimeout(()=>{this.send("logout")
