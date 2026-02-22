package zentao

import (
	"fmt"
	"strings"
)

// ========== Bug管理 ==========

// GetBugs 获取产品的Bug列表
func (c *Client) GetBugs(productID int, limit int) ([]Bug, error) {
	if limit <= 0 {
		limit = 100
	}
	var result BugListResponse
	path := fmt.Sprintf("/api.php/v1/products/%d/bugs?limit=%d", productID, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取Bug列表失败: %v", err)
	}
	return result.Bugs, nil
}

// GetBug 获取Bug详情
func (c *Client) GetBug(bugID int) (*Bug, error) {
	var bug Bug
	path := fmt.Sprintf("/api.php/v1/bugs/%d", bugID)
	if err := c.doGet(path, &bug); err != nil {
		return nil, fmt.Errorf("获取Bug详情失败: %v", err)
	}
	return &bug, nil
}

// GetBugsByProject 根据项目ID过滤Bug列表
func (c *Client) GetBugsByProject(productID, projectID int, limit int) ([]Bug, error) {
	bugs, err := c.GetBugs(productID, limit)
	if err != nil {
		return nil, err
	}

	var filtered []Bug
	for _, bug := range bugs {
		if bug.Project == projectID {
			filtered = append(filtered, bug)
		}
	}
	return filtered, nil
}

// GetBugsByStatus 根据状态过滤Bug列表
func (c *Client) GetBugsByStatus(productID int, status string, limit int) ([]Bug, error) {
	bugs, err := c.GetBugs(productID, limit)
	if err != nil {
		return nil, err
	}

	var filtered []Bug
	for _, bug := range bugs {
		if bug.Status == status {
			filtered = append(filtered, bug)
		}
	}
	return filtered, nil
}

// GetActiveBugs 获取激活状态的Bug列表
func (c *Client) GetActiveBugs(productID int, limit int) ([]Bug, error) {
	return c.GetBugsByStatus(productID, "active", limit)
}

// SearchBugs 搜索Bug（支持多条件过滤）
func (c *Client) SearchBugs(params BugSearchParams) ([]Bug, error) {
	if params.Limit <= 0 {
		params.Limit = 500
	}

	bugs, err := c.GetBugs(params.ProductID, params.Limit)
	if err != nil {
		return nil, err
	}

	var filtered []Bug
	for _, bug := range bugs {
		if params.Status != "" && bug.Status != params.Status {
			continue
		}
		if params.AssignedTo != "" && bug.AssignedTo.Account != params.AssignedTo {
			continue
		}
		if params.Keyword != "" {
			if !containsIgnoreCase(bug.Title, params.Keyword) && fmt.Sprintf("%d", bug.ID) != params.Keyword {
				continue
			}
		}
		if params.Severity > 0 && bug.Severity != params.Severity {
			continue
		}
		if params.Pri > 0 && bug.Pri != params.Pri {
			continue
		}
		filtered = append(filtered, bug)
	}
	return filtered, nil
}

// ConfirmBug 确认Bug
func (c *Client) ConfirmBug(bugID int, req BugConfirmRequest) error {
	path := fmt.Sprintf("/api.php/v1/bugs/%d/confirm", bugID)
	if err := c.doPost(path, req, nil); err != nil {
		return fmt.Errorf("确认Bug失败: %v", err)
	}
	return nil
}

// ResolveBug 解决Bug
func (c *Client) ResolveBug(bugID int, req BugResolveRequest) error {
	path := fmt.Sprintf("/api.php/v1/bugs/%d/resolve", bugID)
	if err := c.doPost(path, req, nil); err != nil {
		return fmt.Errorf("解决Bug失败: %v", err)
	}
	return nil
}

// CloseBug 关闭Bug
func (c *Client) CloseBug(bugID int, req BugCloseRequest) error {
	path := fmt.Sprintf("/api.php/v1/bugs/%d/close", bugID)
	if err := c.doPost(path, req, nil); err != nil {
		return fmt.Errorf("关闭Bug失败: %v", err)
	}
	return nil
}

// ActivateBug 激活Bug
func (c *Client) ActivateBug(bugID int, req BugActivateRequest) error {
	path := fmt.Sprintf("/api.php/v1/bugs/%d/activate", bugID)
	if err := c.doPost(path, req, nil); err != nil {
		return fmt.Errorf("激活Bug失败: %v", err)
	}
	return nil
}

// AssignBug 指派Bug
func (c *Client) AssignBug(bugID int, req BugAssignRequest) error {
	path := fmt.Sprintf("/api.php/v1/bugs/%d/assign", bugID)
	if err := c.doPost(path, req, nil); err != nil {
		return fmt.Errorf("指派Bug失败: %v", err)
	}
	return nil
}

// containsIgnoreCase 忽略大小写的字符串包含检查
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
