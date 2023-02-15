package sso

import "github.com/jda/go-crowd"

// NewCrowdSso 初始化CrowdSso
func NewCrowdSso(appName, appPass, crowdAddr string) (*crowd.Crowd, error) {
	cli, err := crowd.New(appName, appPass, crowdAddr)
	return &cli, err
}
