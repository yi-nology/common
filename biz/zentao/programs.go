package zentao

import "fmt"

// ========== 项目集(Program)管理 ==========

// GetPrograms 获取项目集列表（支持分页）
func (c *Client) GetPrograms(page, limit int) (*ProgramListResponse, error) {
	var result ProgramListResponse
	path := fmt.Sprintf("/api.php/v1/programs?page=%d&limit=%d", page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取项目集列表失败: %v", err)
	}
	return &result, nil
}

// GetProgram 获取项目集详情
func (c *Client) GetProgram(programID int) (*Program, error) {
	var program Program
	path := fmt.Sprintf("/api.php/v1/programs/%d", programID)
	if err := c.doGet(path, &program); err != nil {
		return nil, fmt.Errorf("获取项目集详情失败: %v", err)
	}
	return &program, nil
}

// CreateProgram 创建项目集
func (c *Client) CreateProgram(req ProgramCreateRequest) (*Program, error) {
	var program Program
	if err := c.doPost("/api.php/v1/programs", req, &program); err != nil {
		return nil, fmt.Errorf("创建项目集失败: %v", err)
	}
	return &program, nil
}

// UpdateProgram 更新项目集
func (c *Client) UpdateProgram(programID int, req ProgramCreateRequest) (*Program, error) {
	var program Program
	path := fmt.Sprintf("/api.php/v1/programs/%d", programID)
	if err := c.doPut(path, req, &program); err != nil {
		return nil, fmt.Errorf("更新项目集失败: %v", err)
	}
	return &program, nil
}

// DeleteProgram 删除项目集
func (c *Client) DeleteProgram(programID int) error {
	path := fmt.Sprintf("/api.php/v1/programs/%d", programID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除项目集失败: %v", err)
	}
	return nil
}
