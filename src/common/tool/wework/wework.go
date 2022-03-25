package wework

import (
	"bytes"
	"eago/common/log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	HTTP_TIMEOUT_SECONDS                = 5
	TOKEN_EXPIRATION_ADVANCE_SECONDS    = 60
	WEWORK_GET_DEPARTMENT_BASEURL       = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s&id=1"
	WEWORK_GET_DEPARTMENT_USERS_BASEURL = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?access_token=%s&department_id=1&fetch_child=1"
	WEWORK_GET_USER_INFO_BASEURL        = "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s"
	WEWORK_SEND_MESSAGE_BASEURL         = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
	WEWORK_GEN_TOKEN_BASEURL            = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)

type Wework interface {
	ListDepartments() ([]*WeworkDepartment, error)
	ListUsersWithDepartment() ([]*WeworkUsersDepartment, error)
	SendWework(contentType, subject string, content interface{}, to []string) error
}

// wework struct
type wework struct {
	token           string
	tokenExpireTime time.Time

	agentId    string
	corpId     string
	corpSecret string
}

// NewWework
func NewWework(agentId, corpId, corpSecret string) Wework {
	return &wework{
		agentId:    agentId,
		corpId:     corpId,
		corpSecret: corpSecret,
	}
}

// SendWework 发送微信消息
func (w *wework) SendWework(contentType, subject string, content interface{}, to []string) error {
	log.InfoWithFields(log.Fields{
		"content_type": contentType,
		"subject":      subject,
		"to":           to,
	}, "Wework.SendWework called.")

	message := make(map[string]interface{})
	message["touser"] = strings.Join(to, "|")
	message["msgtype"] = contentType
	message["agentid"] = w.agentId

	// 区分不同的内容类型进行处理
	switch contentType {
	case "textcard":
		log.Debug("Wework sent content type: textcard.")
		// 获得textcard
		textCard := content.(map[string]interface{})
		textCard["title"] = subject
		message["textcard"] = textCard

	default:
		log.Debug("Wework sent default content type: text.")
		message["safe"] = 0
		text := make(map[string]interface{})
		text["content"] = content.(string)
		message["text"] = text
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	//准备发送数据
	w.validateToken()
	client := &http.Client{Timeout: HTTP_TIMEOUT_SECONDS * time.Second}
	url := fmt.Sprintf(WEWORK_SEND_MESSAGE_BASEURL, w.token)
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "Error in http.Client.Get.")
		return err
	}
	// 结束后关闭IO
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("Http reponse status code is %v not 200.", resp.StatusCode)
		log.ErrorWithFields(log.Fields{
			"status_code": resp.StatusCode,
			"error":       err,
		}, "Error in http.Client.Post.")
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"status_code": resp.StatusCode,
			"error":       err,
		}, "Error in ioutil.ReadAll(resp.Body).")
	}

	jsonBody := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "Error in json.Unmarshal.")
		return err
	}
	if jsonBody["errcode"].(float64) != 0 {
		err := fmt.Errorf("Http reponse errcode is %v not 0.", jsonBody["errcode"].(float64))
		log.ErrorWithFields(log.Fields{
			"errcode": jsonBody["errcode"].(float64),
			"errmsg":  jsonBody["errmsg"].(string),
			"error":   err,
		}, "Error in http.Client.Get.")
		return err
	}

	return nil
}

// ListUsersWithDepartment 列出用户和其部门信息
func (w *wework) ListUsersWithDepartment() ([]*WeworkUsersDepartment, error) {
	wud := make([]*WeworkUsersDepartment, 0)

	// 超时时间：5秒
	client := &http.Client{Timeout: HTTP_TIMEOUT_SECONDS * time.Second}

	w.validateToken()
	url := fmt.Sprintf(WEWORK_GET_DEPARTMENT_USERS_BASEURL, w.token)
	resp, err := client.Get(url)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "Error in http.Client.Get.")
		return wud, err
	}
	// 结束后关闭IO
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"status_code": resp.StatusCode,
			"error":       err,
		}, "Error in ioutil.ReadAll(resp.Body).")
	}

	jsonBody := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "Error in json.Unmarshal.")
		return wud, err
	}
	if jsonBody["errcode"].(float64) != 0 {
		err := fmt.Errorf("Http reponse errcode is %v not 0.", jsonBody["errcode"].(float64))
		log.ErrorWithFields(log.Fields{
			"errcode": jsonBody["errcode"].(float64),
			"errmsg":  jsonBody["errmsg"].(string),
			"error":   err,
		}, "Error in http.Client.Get.")
		return wud, err
	}

	for _, v := range jsonBody["userlist"].([]interface{}) {
		mp := v.(map[string]interface{})
		tmp := &WeworkUsersDepartment{
			Username:    mp["userid"].(string),
			DisplayName: mp["name"].(string),
		}
		depts := mp["department"].([]interface{})
		if len(depts) >= 1 {
			tmp.DepartmentId = int(depts[0].(float64))
		}

		wud = append(wud, tmp)
	}

	return wud, nil
}

// ListDepartments 列出部门信息
func (w *wework) ListDepartments() ([]*WeworkDepartment, error) {
	wd := make([]*WeworkDepartment, 0)

	// 超时时间：5秒
	client := &http.Client{Timeout: HTTP_TIMEOUT_SECONDS * time.Second}

	w.validateToken()
	url := fmt.Sprintf(WEWORK_GET_DEPARTMENT_BASEURL, w.token)
	resp, err := client.Get(url)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "Error in http.Client.Get.")
		return wd, err
	}
	// 结束后关闭IO
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"status_code": resp.StatusCode,
			"error":       err,
		}, "Error in ioutil.ReadAll(resp.Body).")
	}

	jsonBody := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "Error in json.Unmarshal.")
		return wd, err
	}
	if jsonBody["errcode"].(float64) != 0 {
		err := fmt.Errorf("Http reponse errcode is %v not 0.", jsonBody["errcode"].(float64))
		log.ErrorWithFields(log.Fields{
			"errcode": jsonBody["errcode"].(float64),
			"errmsg":  jsonBody["errmsg"].(string),
			"error":   err,
		}, "Error in http.Client.Get.")
		return wd, err
	}

	for _, v := range jsonBody["department"].([]interface{}) {
		mp := v.(map[string]interface{})
		wd = append(wd, &WeworkDepartment{
			Id:       int(mp["id"].(float64)),
			Name:     mp["name"].(string),
			ParentId: int(mp["parentid"].(float64)),
		})
	}

	return wd, nil
}

// validateToken 验证Token是否有效
func (w *wework) validateToken() {
	// 如果token不为空，并且当前时间小于tokenExpireTime，则不用重新获取token
	if w.token != "" && time.Now().Before(w.tokenExpireTime) {
		return
	}

	log.DebugWithFields(log.Fields{
		"token":       w.token,
		"expire_time": w.tokenExpireTime.Format("2006-01-02 15:04:05"),
	}, "Token expired, try to get new.")
	_ = w.getToken()
}

// getToken 登录获得token
func (w *wework) getToken() error {
	// 超时时间：5秒
	client := &http.Client{Timeout: HTTP_TIMEOUT_SECONDS * time.Second}

	url := fmt.Sprintf(WEWORK_GEN_TOKEN_BASEURL, w.corpId, w.corpSecret)

	resp, err := client.Get(url)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "Error in http.Client.Get.")
		return err
	}
	// 结束后关闭IO
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("Http reponse status code is %d not 200.", resp.StatusCode)
		log.ErrorWithFields(log.Fields{
			"status_code": resp.StatusCode,
			"error":       err,
		}, "Error in http.Client.Get.")
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"status_code": resp.StatusCode,
			"error":       err,
		}, "Error in ioutil.ReadAll(resp.Body).")
	}

	jsonBody := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "Error in json.Unmarshal.")
		return err
	}
	if jsonBody["errcode"].(float64) != 0 {
		err := fmt.Errorf("Http reponse errcode is %v not 0.", jsonBody["errcode"].(float64))
		log.ErrorWithFields(log.Fields{
			"errcode": jsonBody["errcode"].(float64),
			"errmsg":  jsonBody["errmsg"].(string),
			"error":   err,
		}, "Error in http.Client.Get.")
		return err
	}

	w.token = jsonBody["access_token"].(string)
	ttl := time.Duration(jsonBody["expires_in"].(float64) - TOKEN_EXPIRATION_ADVANCE_SECONDS)
	w.tokenExpireTime = time.Now().Add(time.Second * ttl)
	return nil
}
