package zentao

import "fmt"

// ========== 用例(Case)相关 ==========

// CaseStep 用例步骤
type CaseStep struct {
	ID      int    `json:"id"`
	Parent  int    `json:"parent"`
	Case    int    `json:"case"`
	Version int    `json:"version"`
	Type    string `json:"type"`
	Desc    string `json:"desc"`
	Expect  string `json:"expect"`
}

// Case 用例结构
type Case struct {
	ID             int       `json:"id"`
	Project        int       `json:"project"`
	Product        int       `json:"product"`
	Execution      int       `json:"execution"`
	Branch         int       `json:"branch"`
	Lib            int       `json:"lib"`
	Module         int       `json:"module"`
	Path           int       `json:"path"`
	Story          int       `json:"story"`
	StoryVersion   int       `json:"storyVersion"`
	Title          string    `json:"title"`
	Precondition   string    `json:"precondition"`
	Keywords       string    `json:"keywords"`
	Pri            int       `json:"pri"`
	Type           string    `json:"type"`
	Auto           string    `json:"auto"`
	Frame          string    `json:"frame"`
	Stage          string    `json:"stage"`
	HowRun         string    `json:"howRun"`
	ScriptedBy     string    `json:"scriptedBy"`
	ScriptedDate   *string   `json:"scriptedDate"`
	ScriptStatus   string    `json:"scriptStatus"`
	ScriptLocation string    `json:"scriptLocation"`
	Status         string    `json:"status"`
	SubStatus      string    `json:"subStatus"`
	Color          string    `json:"color"`
	Frequency      string    `json:"frequency"`
	Order          int       `json:"order"`
	OpenedBy       UserRef   `json:"openedBy"`
	OpenedDate     string    `json:"openedDate"`
	ReviewedBy     *UserRef  `json:"reviewedBy"`
	ReviewedDate   *string   `json:"reviewedDate"`
	LastEditedBy   *UserRef  `json:"lastEditedBy"`
	LastEditedDate *string   `json:"lastEditedDate"`
	Version        int       `json:"version"`
	LinkCase       string    `json:"linkCase"`
	FromBug        int       `json:"fromBug"`
	FromCaseID     int       `json:"fromCaseID"`
	FromCaseVersion int     `json:"fromCaseVersion"`
	Deleted        bool      `json:"deleted"`
	LastRunner     string    `json:"lastRunner"`
	LastRunDate    *string   `json:"lastRunDate"`
	LastRunResult  string    `json:"lastRunResult"`
	StoryTitle      *string  `json:"storyTitle"`
	Needconfirm     bool     `json:"needconfirm"`
	Bugs            int      `json:"bugs"`
	Results         int      `json:"results"`
	CaseFails       int      `json:"caseFails"`
	StepNumber      int      `json:"stepNumber"`
	StatusName      string   `json:"statusName"`
	ToBugs          []int    `json:"toBugs"`
	Steps           []CaseStep `json:"steps"`
	Files           []string `json:"files"`
	CurrentVersion  int      `json:"currentVersion"`
}

// CaseListResponse 获取产品用例列表响应
type CaseListResponse struct {
	Page      int    `json:"page"`
	Total     int    `json:"total"`
	Limit     int    `json:"limit"`
	Testcases []Case `json:"testcases"`
}

// CaseCreateRequest 创建用例请求
type CaseCreateRequest struct {
	Branch       int       `json:"branch,omitempty"`
	Module       int       `json:"module,omitempty"`
	Story        int       `json:"story,omitempty"`
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	Stage        string    `json:"stage,omitempty"`
	Precondition string    `json:"precondition,omitempty"`
	Pri          int       `json:"pri,omitempty"`
	Steps        []CaseStep `json:"steps,omitempty"`
	Keywords     string    `json:"keywords,omitempty"`
}

// CaseUpdateRequest 修改用例请求
type CaseUpdateRequest struct {
	Branch       int       `json:"branch,omitempty"`
	Module       int       `json:"module,omitempty"`
	Story        int       `json:"story,omitempty"`
	Title        string    `json:"title,omitempty"`
	Type         string    `json:"type,omitempty"`
	Stage        string    `json:"stage,omitempty"`
	Precondition string    `json:"precondition,omitempty"`
	Pri          int       `json:"pri,omitempty"`
	Steps        []CaseStep `json:"steps,omitempty"`
	Keywords     string    `json:"keywords,omitempty"`
}

// CaseDeleteResponse 删除用例响应
type CaseDeleteResponse struct {
	Message string `json:"message"`
}

// ========== 用例(Case) API 方法 ==========

// GetCasesByProduct 获取产品用例列表（支持分页）
func (c *Client) GetCasesByProduct(productID int, page, limit int) (*CaseListResponse, error) {
	var result CaseListResponse
	path := fmt.Sprintf("/api.php/v1/products/%d/cases?page=%d&limit=%d", productID, page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取用例列表失败: %v", err)
	}
	return &result, nil
}

// GetCase 获取用例详情
func (c *Client) GetCase(caseID int) (*Case, error) {
	var cas Case
	path := fmt.Sprintf("/api.php/v1/cases/%d", caseID)
	if err := c.doGet(path, &cas); err != nil {
		return nil, fmt.Errorf("获取用例详情失败: %v", err)
	}
	return &cas, nil
}

// CreateCase 创建用例
func (c *Client) CreateCase(productID int, req CaseCreateRequest) (*Case, error) {
	var cas Case
	path := fmt.Sprintf("/api.php/v1/products/%d/cases", productID)
	if err := c.doPost(path, req, &cas); err != nil {
		return nil, fmt.Errorf("创建用例失败: %v", err)
	}
	return &cas, nil
}

// UpdateCase 修改用例
func (c *Client) UpdateCase(caseID int, req CaseUpdateRequest) (*Case, error) {
	var cas Case
	path := fmt.Sprintf("/api.php/v1/cases/%d", caseID)
	if err := c.doPut(path, req, &cas); err != nil {
		return nil, fmt.Errorf("修改用例失败: %v", err)
	}
	return &cas, nil
}

// DeleteCase 删除用例
func (c *Client) DeleteCase(caseID int) error {
	path := fmt.Sprintf("/api.php/v1/cases/%d", caseID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除用例失败: %v", err)
	}
	return nil
}
