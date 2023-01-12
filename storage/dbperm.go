package storage

// perm 表相关操作

// GetSelfAlwaysToken 创建自己前发给别人的一条永久 token
// 数据库里面只会有一条这样的数据
// TokenType = TokenTypeFromSelf, 且 PermType = PermTypeAlways
func GetSelfAlwaysToken() {
	// TODO 实现
}
