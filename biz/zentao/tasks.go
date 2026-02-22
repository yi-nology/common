package zentao

import "fmt"

// ========== 任务(Task)管理 ==========

// GetTasks 获取执行的任务列表(支持分页)
func (c *Client) GetTasks(executionID int, limit int) ([]Task, error) {
	if limit <= 0 {
		limit = 500
	}

	var allTasks []Task
	page := 1
	for {
		var result TaskListResponse
		path := fmt.Sprintf("/api.php/v1/executions/%d/tasks?limit=%d&page=%d", executionID, limit, page)
		if err := c.doGet(path, &result); err != nil {
			return nil, fmt.Errorf("获取任务列表失败: %v", err)
		}

		allTasks = append(allTasks, result.Tasks...)

		if len(allTasks) >= result.Total || len(result.Tasks) == 0 {
			break
		}
		page++
	}

	return allTasks, nil
}

// GetTask 获取任务详情
func (c *Client) GetTask(taskID int) (*Task, error) {
	var task Task
	path := fmt.Sprintf("/api.php/v1/tasks/%d", taskID)
	if err := c.doGet(path, &task); err != nil {
		return nil, fmt.Errorf("获取任务详情失败: %v", err)
	}
	return &task, nil
}

// CreateTask 创建任务
func (c *Client) CreateTask(executionID int, req TaskCreateRequest) (*Task, error) {
	var task Task
	path := fmt.Sprintf("/api.php/v1/executions/%d/tasks", executionID)
	if err := c.doPost(path, req, &task); err != nil {
		return nil, fmt.Errorf("创建任务失败: %v", err)
	}
	return &task, nil
}

// UpdateTask 更新任务
func (c *Client) UpdateTask(taskID int, req TaskUpdateRequest) (*Task, error) {
	var task Task
	path := fmt.Sprintf("/api.php/v1/tasks/%d", taskID)
	if err := c.doPut(path, req, &task); err != nil {
		return nil, fmt.Errorf("更新任务失败: %v", err)
	}
	return &task, nil
}

// DeleteTask 删除任务
func (c *Client) DeleteTask(taskID int) error {
	path := fmt.Sprintf("/api.php/v1/tasks/%d", taskID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除任务失败: %v", err)
	}
	return nil
}

// StartTask 开始任务
func (c *Client) StartTask(taskID int, req TaskStartRequest) error {
	path := fmt.Sprintf("/api.php/v1/tasks/%d/start", taskID)
	if err := c.doPost(path, req, nil); err != nil {
		return fmt.Errorf("开始任务失败: %v", err)
	}
	return nil
}

// FinishTask 完成任务
func (c *Client) FinishTask(taskID int, req TaskFinishRequest) error {
	path := fmt.Sprintf("/api.php/v1/tasks/%d/finish", taskID)
	if err := c.doPost(path, req, nil); err != nil {
		return fmt.Errorf("完成任务失败: %v", err)
	}
	return nil
}

// PauseTask 暂停任务
func (c *Client) PauseTask(taskID int, req TaskPauseRequest) error {
	path := fmt.Sprintf("/api.php/v1/tasks/%d/pause", taskID)
	if err := c.doPost(path, req, nil); err != nil {
		return fmt.Errorf("暂停任务失败: %v", err)
	}
	return nil
}

// ActivateTask 激活任务
func (c *Client) ActivateTask(taskID int, consumed float64, left float64) error {
	body := map[string]interface{}{
		"consumed": consumed,
		"left":     left,
	}
	path := fmt.Sprintf("/api.php/v1/tasks/%d/restart", taskID)
	if err := c.doPost(path, body, nil); err != nil {
		return fmt.Errorf("激活任务失败: %v", err)
	}
	return nil
}

// AssignTask 指派任务
func (c *Client) AssignTask(taskID int, req TaskAssignRequest) error {
	body := map[string]interface{}{
		"assignedTo": req.AssignedTo,
	}
	if req.Left > 0 {
		body["left"] = req.Left
	}
	path := fmt.Sprintf("/api.php/v1/tasks/%d", taskID)
	if err := c.doPut(path, body, nil); err != nil {
		return fmt.Errorf("指派任务失败: %v", err)
	}
	return nil
}
