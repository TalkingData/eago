package sso

import (
	"eago/auth/conf"
	"github.com/jda/go-crowd"
)

var Crowd *crowd.Crowd

// InitCrowd 初始化Crowd
func InitCrowd() error {
	client, err := crowd.New(conf.Conf.CrowdAppName, conf.Conf.CrowdAppPassword, conf.Conf.CrowdAddress)

	if err != nil {
		return err
	}

	Crowd = &client

	return nil
}
