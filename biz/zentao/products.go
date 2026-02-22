package zentao

import "fmt"

// ========== 产品(Product)管理 ==========

// GetProducts 获取产品列表
func (c *Client) GetProducts() ([]Product, error) {
	var result ProductListResponse
	if err := c.doGet("/api.php/v1/products?limit=500", &result); err != nil {
		return nil, fmt.Errorf("获取产品列表失败: %v", err)
	}
	return result.Products, nil
}

// GetProduct 获取产品详情
func (c *Client) GetProduct(productID int) (*Product, error) {
	var product Product
	path := fmt.Sprintf("/api.php/v1/products/%d", productID)
	if err := c.doGet(path, &product); err != nil {
		return nil, fmt.Errorf("获取产品详情失败: %v", err)
	}
	return &product, nil
}

// CreateProduct 创建产品
func (c *Client) CreateProduct(req ProductCreateRequest) (*Product, error) {
	var product Product
	if err := c.doPost("/api.php/v1/products", req, &product); err != nil {
		return nil, fmt.Errorf("创建产品失败: %v", err)
	}
	return &product, nil
}

// UpdateProduct 更新产品
func (c *Client) UpdateProduct(productID int, req ProductCreateRequest) (*Product, error) {
	var product Product
	path := fmt.Sprintf("/api.php/v1/products/%d", productID)
	if err := c.doPut(path, req, &product); err != nil {
		return nil, fmt.Errorf("更新产品失败: %v", err)
	}
	return &product, nil
}

// DeleteProduct 删除产品
func (c *Client) DeleteProduct(productID int) error {
	path := fmt.Sprintf("/api.php/v1/products/%d", productID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除产品失败: %v", err)
	}
	return nil
}
