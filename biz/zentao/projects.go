package zentao

import "fmt"

// ========== 项目(Project)管理 ==========

// GetAllProjects 获取所有项目列表（支持分页）
func (c *Client) GetAllProjects(page, limit int) (*ProjectListResponse, error) {
	var result ProjectListResponse
	path := fmt.Sprintf("/api.php/v1/projects?page=%d&limit=%d", page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取项目列表失败: %v", err)
	}
	return &result, nil
}

// GetProjectsByProduct 获取产品关联的项目列表（支持分页）
func (c *Client) GetProjectsByProduct(productID int, page, limit int) (*ProjectListResponse, error) {
	var result ProjectListResponse
	path := fmt.Sprintf("/api.php/v1/products/%d/projects?page=%d&limit=%d", productID, page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取项目列表失败: %v", err)
	}
	return &result, nil
}

// GetProject 获取项目详情
func (c *Client) GetProject(projectID int) (*Project, error) {
	var project Project
	path := fmt.Sprintf("/api.php/v1/projects/%d", projectID)
	if err := c.doGet(path, &project); err != nil {
		return nil, fmt.Errorf("获取项目详情失败: %v", err)
	}
	return &project, nil
}

// CreateProject 创建项目
func (c *Client) CreateProject(req ProjectCreateRequest) (*Project, error) {
	var project Project
	if err := c.doPost("/api.php/v1/projects", req, &project); err != nil {
		return nil, fmt.Errorf("创建项目失败: %v", err)
	}
	return &project, nil
}

// UpdateProject 更新项目
func (c *Client) UpdateProject(projectID int, req ProjectCreateRequest) (*Project, error) {
	var project Project
	path := fmt.Sprintf("/api.php/v1/projects/%d", projectID)
	if err := c.doPut(path, req, &project); err != nil {
		return nil, fmt.Errorf("更新项目失败: %v", err)
	}
	return &project, nil
}

// DeleteProject 删除项目
func (c *Client) DeleteProject(projectID int) error {
	path := fmt.Sprintf("/api.php/v1/projects/%d", projectID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除项目失败: %v", err)
	}
	return nil
}
