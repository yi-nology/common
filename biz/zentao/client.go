// Package zentao provides a Go SDK for Zentao (禅道) project management system API.
package zentao

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
)

// Client 禅道API客户端
type Client struct {
	BaseURL string
	Token   string
	client  *req.Client
}

// NewClient 创建新的禅道客户端
func NewClient(baseURL string) *Client {
	client := req.C().
		SetTimeout(30*time.Second).
		SetCommonHeader("Content-Type", "application/json").
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return &Client{
		BaseURL: baseURL,
		client:  client,
	}
}

// SetTimeout 设置请求超时时间
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.client.SetTimeout(timeout)
	return c
}

// SetToken 设置认证Token
func (c *Client) SetToken(token string) *Client {
	c.Token = token
	return c
}

// GetToken 获取禅道token
func (c *Client) GetToken(account, password string) (string, error) {
	tokenReq := TokenRequest{
		Account:  account,
		Password: password,
	}

	var tokenResp TokenResponse
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBodyJsonMarshal(tokenReq).
		SetSuccessResult(&tokenResp).
		Post(c.BaseURL + "/api.php/v1/tokens")

	if err != nil {
		return "", fmt.Errorf("请求禅道token失败: %v", err)
	}

	if !resp.IsSuccessState() {
		return "", fmt.Errorf("获取禅道token失败, 状态码: %d, 响应: %s", resp.StatusCode, resp.String())
	}

	c.Token = tokenResp.Token
	return tokenResp.Token, nil
}

// doGet 执行GET请求
func (c *Client) doGet(path string, result interface{}) error {
	if c.Token == "" {
		return fmt.Errorf("token为空，请先调用GetToken获取token")
	}

	url := c.BaseURL + path
	resp, err := c.client.R().
		SetHeader("Token", c.Token).
		Get(url)

	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}

	if !resp.IsSuccessState() {
		return fmt.Errorf("请求失败, 状态码: %d, 响应: %s", resp.StatusCode, resp.String())
	}

	if result != nil {
		if err := json.Unmarshal(resp.Bytes(), result); err != nil {
			return fmt.Errorf("解析响应失败: %v", err)
		}
	}
	return nil
}

// doPost 执行POST请求
func (c *Client) doPost(path string, body interface{}, result interface{}) error {
	if c.Token == "" {
		return fmt.Errorf("token为空，请先调用GetToken获取token")
	}

	url := c.BaseURL + path
	resp, err := c.client.R().
		SetHeader("Token", c.Token).
		SetBodyJsonMarshal(body).
		Post(url)

	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}

	if !resp.IsSuccessState() {
		return fmt.Errorf("请求失败, 状态码: %d, 响应: %s", resp.StatusCode, resp.String())
	}

	if result != nil {
		if err := json.Unmarshal(resp.Bytes(), result); err != nil {
			return fmt.Errorf("解析响应失败: %v", err)
		}
	}
	return nil
}

// doPut 执行PUT请求
func (c *Client) doPut(path string, body interface{}, result interface{}) error {
	if c.Token == "" {
		return fmt.Errorf("token为空，请先调用GetToken获取token")
	}

	url := c.BaseURL + path
	resp, err := c.client.R().
		SetHeader("Token", c.Token).
		SetBodyJsonMarshal(body).
		Put(url)

	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}

	if !resp.IsSuccessState() {
		return fmt.Errorf("请求失败, 状态码: %d, 响应: %s", resp.StatusCode, resp.String())
	}

	if result != nil {
		if err := json.Unmarshal(resp.Bytes(), result); err != nil {
			return fmt.Errorf("解析响应失败: %v", err)
		}
	}
	return nil
}

// doDelete 执行DELETE请求
func (c *Client) doDelete(path string) error {
	if c.Token == "" {
		return fmt.Errorf("token为空，请先调用GetToken获取token")
	}

	url := c.BaseURL + path
	resp, err := c.client.R().
		SetHeader("Token", c.Token).
		Delete(url)

	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}

	if !resp.IsSuccessState() {
		return fmt.Errorf("请求失败, 状态码: %d, 响应: %s", resp.StatusCode, resp.String())
	}

	return nil
}

// getRawBytes 获取原始响应字节
func (c *Client) getRawBytes(path string) ([]byte, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("token为空，请先调用GetToken获取token")
	}

	url := c.BaseURL + path
	resp, err := c.client.R().
		SetHeader("Token", c.Token).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("请求失败, 状态码: %d, 响应: %s", resp.StatusCode, resp.String())
	}

	return resp.Bytes(), nil
}
