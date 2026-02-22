package zentao

import "fmt"

// ========== 发布(Release)管理 ==========

// GetReleasesByProduct 获取产品的发布列表
func (c *Client) GetReleasesByProduct(productID int) ([]Release, error) {
	var result ReleaseListResponse
	path := fmt.Sprintf("/api.php/v1/products/%d/releases?limit=100", productID)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取发布列表失败: %v", err)
	}
	return result.Releases, nil
}

// GetReleasesByProject 获取项目的发布列表
func (c *Client) GetReleasesByProject(projectID int) ([]Release, error) {
	var result ReleaseListResponse
	path := fmt.Sprintf("/api.php/v1/projects/%d/releases?limit=100", projectID)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取发布列表失败: %v", err)
	}
	return result.Releases, nil
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
