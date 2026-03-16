package zentao

import "fmt"

// ========== 测试单(TestTask)相关 ==========

// TestTask 测试单结构
type TestTask struct {
	ID               int       `json:"id"`
	Project          int       `json:"project"`
	Product          int       `json:"product"`
	Name             string    `json:"name"`
	Execution        int       `json:"execution"`
	Build            int       `json:"build"`
	Type             string    `json:"type"`
	Owner            UserRef   `json:"owner"`
	Pri              int       `json:"pri"`
	Begin            string    `json:"begin"`
	End              string    `json:"end"`
	RealFinishedDate *string   `json:"realFinishedDate"`
	Mailto           string    `json:"mailto"`
	Desc             string    `json:"desc"`
	Report           string    `json:"report"`
	Status           string    `json:"status"`
	TestReport       int       `json:"testreport"`
	Auto             string    `json:"auto"`
	SubStatus        string    `json:"subStatus"`
	Deleted          string    `json:"deleted"`
	ProductName       string   `json:"productName"`
	ExecutionName     string   `json:"executionName"`
	BuildName         string   `json:"buildName"`
	Branch           int       `json:"branch"`
}

// TestTaskListResponse 获取测试单列表响应
type TestTaskListResponse struct {
	Page      int         `json:"page"`
	Total     int         `json:"total"`
	Limit     int         `json:"limit"`
	Testtasks []TestTask `json:"testtasks"`
}

// ProjectTestTaskListResponse 获取项目测试单列表响应
type ProjectTestTaskListResponse struct {
	Page      int         `json:"page"`
	Total     int         `json:"total"`
	Limit     int         `json:"limit"`
	Testtasks []TestTask `json:"testtasks"`
}

// ========== 测试单(TestTask) API 方法 ==========

// GetTestTasks 获取测试单列表（支持分页）
func (c *Client) GetTestTasks(page, limit int) (*TestTaskListResponse, error) {
	var result TestTaskListResponse
	path := fmt.Sprintf("/api.php/v1/testtasks?page=%d&limit=%d", page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取测试单列表失败: %v", err)
	}
	return &result, nil
}

// GetTestTasksByProject 获取项目的测试单列表
func (c *Client) GetTestTasksByProject(projectID int, page, limit int) (*ProjectTestTaskListResponse, error) {
	var result ProjectTestTaskListResponse
	path := fmt.Sprintf("/api.php/v1/projects/%d/testtasks?page=%d&limit=%d", projectID, page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取项目测试单列表失败: %v", err)
	}
	return &result, nil
}

// GetTestTask 获取测试单详情
func (c *Client) GetTestTask(testTaskID int) (*TestTask, error) {
	var testTask TestTask
	path := fmt.Sprintf("/api.php/v1/testtasks/%d", testTaskID)
	if err := c.doGet(path, &testTask); err != nil {
		return nil, fmt.Errorf("获取测试单详情失败: %v", err)
	}
	return &testTask, nil
}
