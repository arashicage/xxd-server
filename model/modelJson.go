package model

type Omit *struct{}

// not working
type SysUserJson struct {
	*SysUser
	Password Omit `json:"-"`
	Nickname Omit `json:"-"`
	Birthday Omit `json:"-"`
	Skype    Omit `json:"-"`
	Qq       Omit `json:"-"`
	Yahoo    Omit `json:"-"`
	Gtalk    Omit `json:"-"`
	Wangwang Omit `json:"-"`
	Address  Omit `json:"-"`
	Zipcode  Omit `json:"-"`
	Visits   Omit `json:"-"`
	Ip       Omit `json:"-"`
	Last     Omit `json:"-"`
	Ping     Omit `json:"-"`
	Fails    Omit `json:"-"`
	Join     Omit `json:"-"`
	Locked   Omit `json:"-"`
	Deleted  Omit `json:"-"`
}
