package zentao

import "fmt"

// ========== 需求(Story)管理 ==========

// GetStoriesByProduct 获取产品的需求列表
func (c *Client) GetStoriesByProduct(productID int, limit int) ([]Story, error) {
	if limit <= 0 {
		limit = 100
	}
	var result StoryListResponse
	path := fmt.Sprintf("/api.php/v1/products/%d/stories?limit=%d", productID, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取需求列表失败: %v", err)
	}
	return result.Stories, nil
}

// GetStoriesByProject 获取项目的需求列表
func (c *Client) GetStoriesByProject(projectID int, limit int) ([]Story, error) {
	if limit <= 0 {
		limit = 100
	}
	var result StoryListResponse
	path := fmt.Sprintf("/api.php/v1/projects/%d/stories?limit=%d", projectID, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取需求列表失败: %v", err)
	}
	return result.Stories, nil
}

// GetStoriesByExecution 获取执行的需求列表
func (c *Client) GetStoriesByExecution(executionID int, limit int) ([]Story, error) {
	if limit <= 0 {
		limit = 100
	}
	var result StoryListResponse
	path := fmt.Sprintf("/api.php/v1/executions/%d/stories?limit=%d", executionID, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取需求列表失败: %v", err)
	}
	return result.Stories, nil
}

// GetStory 获取需求详情
func (c *Client) GetStory(storyID int) (*Story, error) {
	var story Story
	path := fmt.Sprintf("/api.php/v1/stories/%d", storyID)
	if err := c.doGet(path, &story); err != nil {
		return nil, fmt.Errorf("获取需求详情失败: %v", err)
	}
	return &story, nil
}

// CreateStory 创建需求
func (c *Client) CreateStory(productID int, req StoryCreateRequest) (*Story, error) {
	var story Story
	path := fmt.Sprintf("/api.php/v1/products/%d/stories", productID)
	if err := c.doPost(path, req, &story); err != nil {
		return nil, fmt.Errorf("创建需求失败: %v", err)
	}
	return &story, nil
}

// UpdateStory 更新需求
func (c *Client) UpdateStory(storyID int, req StoryCreateRequest) (*Story, error) {
	var story Story
	path := fmt.Sprintf("/api.php/v1/stories/%d", storyID)
	if err := c.doPut(path, req, &story); err != nil {
		return nil, fmt.Errorf("更新需求失败: %v", err)
	}
	return &story, nil
}

// DeleteStory 删除需求
func (c *Client) DeleteStory(storyID int) error {
	path := fmt.Sprintf("/api.php/v1/stories/%d", storyID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除需求失败: %v", err)
	}
	return nil
}

// ChangeStory 变更需求
func (c *Client) ChangeStory(storyID int, spec string, verify string) error {
	body := map[string]interface{}{
		"spec":   spec,
		"verify": verify,
	}
	path := fmt.Sprintf("/api.php/v1/stories/%d/change", storyID)
	if err := c.doPost(path, body, nil); err != nil {
		return fmt.Errorf("变更需求失败: %v", err)
	}
	return nil
}
