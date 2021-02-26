package sso

import (
	"eago-auth/config"
	"github.com/jda/go-crowd"
)

var Crowd *crowd.Crowd

func InitCrowd() error {
	client, err := crowd.New(config.Config.CrowdAppName, config.Config.CrowdAppPassword, config.Config.CrowdAddress)

	if err != nil {
		return err
	}

	Crowd = &client

	return nil
}
