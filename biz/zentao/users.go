package zentao

import "fmt"

// ========== 用户(User)管理 ==========

// GetUsers 获取用户列表
func (c *Client) GetUsers() ([]User, error) {
	var result UserListResponse
	if err := c.doGet("/api.php/v1/users?limit=500", &result); err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %v", err)
	}
	return result.Users, nil
}

// GetCurrentUser 获取当前登录用户信息
func (c *Client) GetCurrentUser() (*User, error) {
	var user User
	if err := c.doGet("/api.php/v1/user", &user); err != nil {
		return nil, fmt.Errorf("获取当前用户信息失败: %v", err)
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户信息
func (c *Client) GetUserByID(userID int) (*User, error) {
	var user User
	path := fmt.Sprintf("/api.php/v1/users/%d", userID)
	if err := c.doGet(path, &user); err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	return &user, nil
}
