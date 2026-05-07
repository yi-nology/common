package zentao

import "fmt"

// ========== 发布(Release)管理 ==========

// GetReleasesByProduct 获取产品的发布列表（支持分页）
func (c *Client) GetReleasesByProduct(productID int, page, limit int) (*ReleaseListResponse, error) {
	var result ReleaseListResponse
	path := fmt.Sprintf("/api.php/v1/products/%d/releases?page=%d&limit=%d", productID, page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取发布列表失败: %v", err)
	}
	return &result, nil
}

// GetReleasesByProject 获取项目的发布列表（支持分页）
func (c *Client) GetReleasesByProject(projectID int, page, limit int) (*ReleaseListResponse, error) {
	var result ReleaseListResponse
	path := fmt.Sprintf("/api.php/v1/projects/%d/releases?page=%d&limit=%d", projectID, page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取发布列表失败: %v", err)
	}
	return &result, nil
}

// GetRelease 获取发布详情
func (c *Client) GetRelease(releaseID int) (*Release, error) {
	var release Release
	path := fmt.Sprintf("/api.php/v1/releases/%d", releaseID)
	if err := c.doGet(path, &release); err != nil {
		return nil, fmt.Errorf("获取发布详情失败: %v", err)
	}
	return &release, nil
}
