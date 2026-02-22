package zentao

import "fmt"

// ========== 计划(Plan)管理 ==========

// GetPlans 获取产品的计划列表
func (c *Client) GetPlans(productID int) ([]Plan, error) {
	var result PlanListResponse
	path := fmt.Sprintf("/api.php/v1/products/%d/plans?limit=100", productID)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取计划列表失败: %v", err)
	}
	return result.Plans, nil
}

// GetPlan 获取计划详情
func (c *Client) GetPlan(planID int) (*Plan, error) {
	var plan Plan
	path := fmt.Sprintf("/api.php/v1/productplans/%d", planID)
	if err := c.doGet(path, &plan); err != nil {
		return nil, fmt.Errorf("获取计划详情失败: %v", err)
	}
	return &plan, nil
}

// CreatePlan 创建计划
func (c *Client) CreatePlan(productID int, req PlanCreateRequest) (*Plan, error) {
	var plan Plan
	path := fmt.Sprintf("/api.php/v1/products/%d/plans", productID)
	if err := c.doPost(path, req, &plan); err != nil {
		return nil, fmt.Errorf("创建计划失败: %v", err)
	}
	return &plan, nil
}

// UpdatePlan 更新计划
func (c *Client) UpdatePlan(planID int, req PlanCreateRequest) (*Plan, error) {
	var plan Plan
	path := fmt.Sprintf("/api.php/v1/productplans/%d", planID)
	if err := c.doPut(path, req, &plan); err != nil {
		return nil, fmt.Errorf("更新计划失败: %v", err)
	}
	return &plan, nil
}

// DeletePlan 删除计划
func (c *Client) DeletePlan(planID int) error {
	path := fmt.Sprintf("/api.php/v1/productplans/%d", planID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除计划失败: %v", err)
	}
	return nil
}

// LinkStoriesToPlan 关联需求到计划
func (c *Client) LinkStoriesToPlan(planID int, storyIDs []int) error {
	body := map[string]interface{}{
		"stories": storyIDs,
	}
	path := fmt.Sprintf("/api.php/v1/productplans/%d/linkstories", planID)
	if err := c.doPost(path, body, nil); err != nil {
		return fmt.Errorf("关联需求失败: %v", err)
	}
	return nil
}

// UnlinkStoriesFromPlan 取消关联需求
func (c *Client) UnlinkStoriesFromPlan(planID int, storyIDs []int) error {
	body := map[string]interface{}{
		"stories": storyIDs,
	}
	path := fmt.Sprintf("/api.php/v1/productplans/%d/unlinkstories", planID)
	if err := c.doPost(path, body, nil); err != nil {
		return fmt.Errorf("取消关联需求失败: %v", err)
	}
	return nil
}

// LinkBugsToPlan 关联Bug到计划
func (c *Client) LinkBugsToPlan(planID int, bugIDs []int) error {
	body := map[string]interface{}{
		"bugs": bugIDs,
	}
	path := fmt.Sprintf("/api.php/v1/productplans/%d/linkbugs", planID)
	if err := c.doPost(path, body, nil); err != nil {
		return fmt.Errorf("关联Bug失败: %v", err)
	}
	return nil
}

// UnlinkBugsFromPlan 取消关联Bug
func (c *Client) UnlinkBugsFromPlan(planID int, bugIDs []int) error {
	body := map[string]interface{}{
		"bugs": bugIDs,
	}
	path := fmt.Sprintf("/api.php/v1/productplans/%d/unlinkbugs", planID)
	if err := c.doPost(path, body, nil); err != nil {
		return fmt.Errorf("取消关联Bug失败: %v", err)
	}
	return nil
}
