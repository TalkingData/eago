package permission

import auth "eago/auth/srv/proto"

var authCli auth.AuthService

// 设置auth cli
func SetAuthClient(c auth.AuthService) {
	authCli = c
}