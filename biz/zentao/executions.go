package zentao

import "fmt"

// ========== 执行/迭代(Execution)管理 ==========

// GetExecutions 获取项目的执行/迭代列表
func (c *Client) GetExecutions(projectID int) ([]Execution, error) {
	var result ExecutionListResponse
	path := fmt.Sprintf("/api.php/v1/projects/%d/executions?limit=100", projectID)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取执行列表失败: %v", err)
	}
	return result.Executions, nil
}

// GetExecution 获取执行详情
func (c *Client) GetExecution(executionID int) (*Execution, error) {
	var execution Execution
	path := fmt.Sprintf("/api.php/v1/executions/%d", executionID)
	if err := c.doGet(path, &execution); err != nil {
		return nil, fmt.Errorf("获取执行详情失败: %v", err)
	}
	return &execution, nil
}

// CreateExecution 创建执行
func (c *Client) CreateExecution(projectID int, req ExecutionCreateRequest) (*Execution, error) {
	var execution Execution
	path := fmt.Sprintf("/api.php/v1/projects/%d/executions", projectID)
	if err := c.doPost(path, req, &execution); err != nil {
		return nil, fmt.Errorf("创建执行失败: %v", err)
	}
	return &execution, nil
}

// UpdateExecution 更新执行
func (c *Client) UpdateExecution(executionID int, req ExecutionCreateRequest) (*Execution, error) {
	var execution Execution
	path := fmt.Sprintf("/api.php/v1/executions/%d", executionID)
	if err := c.doPut(path, req, &execution); err != nil {
		return nil, fmt.Errorf("更新执行失败: %v", err)
	}
	return &execution, nil
}

// DeleteExecution 删除执行
func (c *Client) DeleteExecution(executionID int) error {
	path := fmt.Sprintf("/api.php/v1/executions/%d", executionID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除执行失败: %v", err)
	}
	return nil
}
