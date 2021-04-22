package sso

import (
	"eago-auth/conf"
	"github.com/jda/go-crowd"
)

var Crowd *crowd.Crowd

// InitCrowd 初始化Crowd
func InitCrowd() error {
	client, err := crowd.New(conf.Config.CrowdAppName, conf.Config.CrowdAppPassword, conf.Config.CrowdAddress)

	if err != nil {
		return err
	}

	Crowd = &client

	return nil
}
