package zentao

import "fmt"

// ========== 版本(Build)管理 ==========

// GetBuildsByProject 获取项目的版本列表（支持分页）
func (c *Client) GetBuildsByProject(projectID int, page, limit int) (*BuildListResponse, error) {
	var result BuildListResponse
	path := fmt.Sprintf("/api.php/v1/projects/%d/builds?page=%d&limit=%d", projectID, page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取版本列表失败: %v", err)
	}
	return &result, nil
}

// GetBuildsByExecution 获取执行的版本列表（支持分页）
func (c *Client) GetBuildsByExecution(executionID int, page, limit int) (*BuildListResponse, error) {
	var result BuildListResponse
	path := fmt.Sprintf("/api.php/v1/executions/%d/builds?page=%d&limit=%d", executionID, page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取版本列表失败: %v", err)
	}
	return &result, nil
}

// GetBuild 获取版本详情
func (c *Client) GetBuild(buildID int) (*Build, error) {
	var build Build
	path := fmt.Sprintf("/api.php/v1/builds/%d", buildID)
	if err := c.doGet(path, &build); err != nil {
		return nil, fmt.Errorf("获取版本详情失败: %v", err)
	}
	return &build, nil
}

// CreateBuild 创建版本
func (c *Client) CreateBuild(projectID int, req BuildCreateRequest) (*Build, error) {
	var build Build
	path := fmt.Sprintf("/api.php/v1/projects/%d/builds", projectID)
	if err := c.doPost(path, req, &build); err != nil {
		return nil, fmt.Errorf("创建版本失败: %v", err)
	}
	return &build, nil
}

// UpdateBuild 更新版本
func (c *Client) UpdateBuild(buildID int, req BuildCreateRequest) (*Build, error) {
	var build Build
	path := fmt.Sprintf("/api.php/v1/builds/%d", buildID)
	if err := c.doPut(path, req, &build); err != nil {
		return nil, fmt.Errorf("更新版本失败: %v", err)
	}
	return &build, nil
}

// DeleteBuild 删除版本
func (c *Client) DeleteBuild(buildID int) error {
	path := fmt.Sprintf("/api.php/v1/builds/%d", buildID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除版本失败: %v", err)
	}
	return nil
}
