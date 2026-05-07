package zentao

import "fmt"

// ========== 反馈(Feedback)相关 ==========

// Feedback 反馈结构
type Feedback struct {
	ID            int        `json:"id"`
	Product       int        `json:"product"`
	Module        int        `json:"module"`
	Title         string     `json:"title"`
	Type          string     `json:"type"`
	Solution      string     `json:"solution"`
	Desc          string     `json:"desc"`
	Status        string     `json:"status"`
	SubStatus     string     `json:"subStatus"`
	Public        int        `json:"public"`
	Notify        int        `json:"notify"`
	NotifyEmail   string     `json:"notifyEmail"`
	Likes         string     `json:"likes"`
	Result        int        `json:"result"`
	Faq           int        `json:"faq"`
	OpenedBy      UserRef    `json:"openedBy"`
	OpenedDate    string     `json:"openedDate"`
	ReviewedBy    *string    `json:"reviewedBy"`
	ReviewedDate  *string    `json:"reviewedDate"`
	ProcessedBy    *string    `json:"processedBy"`
	ProcessedDate  *string    `json:"processedDate"`
	ClosedBy      *string    `json:"closedBy"`
	ClosedDate    *string    `json:"closedDate"`
	ClosedReason  string     `json:"closedReason"`
	EditedBy      *UserRef   `json:"editedBy"`
	EditedDate    string     `json:"editedDate"`
	AssignedTo    *UserRef   `json:"assignedTo"`
	AssignedDate  string     `json:"assignedDate"`
	FeedbackBy     string    `json:"feedbackBy"`
	Mailto         []string  `json:"mailto"`
	Deleted       int        `json:"deleted"`
	LikesCount     int        `json:"likesCount"`
	ProductName    string     `json:"productName"`
	ModuleName     string     `json:"moduleName"`
}

// FeedbackListResponse 获取反馈列表响应
type FeedbackListResponse struct {
	Page       int         `json:"page"`
	Total      int         `json:"total"`
	Limit      int         `json:"limit"`
	Feedbacks []Feedback `json:"feedbacks"`
}

// FeedbackCreateRequest 创建反馈请求
type FeedbackCreateRequest struct {
	Product       int       `json:"product"`
	Module        int       `json:"module,omitempty"`
	Title         string    `json:"title"`
	Type          string    `json:"type,omitempty"`
	Desc          string    `json:"desc,omitempty"`
	Public        int       `json:"public,omitempty"`
	Notify        int       `json:"notify,omitempty"`
	NotifyEmail   string    `json:"notifyEmail,omitempty"`
	FeedbackBy     string   `json:"feedbackBy,omitempty"`
}

// FeedbackUpdateRequest 修改反馈请求
type FeedbackUpdateRequest struct {
	Product       int       `json:"product,omitempty"`
	Module        int       `json:"module,omitempty"`
	Title         string    `json:"title,omitempty"`
	Type          string    `json:"type,omitempty"`
	Desc          string    `json:"desc,omitempty"`
	Public        int       `json:"public,omitempty"`
	Notify        int       `json:"notify,omitempty"`
	NotifyEmail   string    `json:"notifyEmail,omitempty"`
	FeedbackBy     string   `json:"feedbackBy,omitempty"`
}

// FeedbackAssignRequest 指派反馈请求
type FeedbackAssignRequest struct {
	AssignedTo string `json:"assignedTo"`
	Comment    string `json:"comment,omitempty"`
	Mailto     string `json:"mailto,omitempty"`
}

// FeedbackCloseRequest 关闭反馈请求
type FeedbackCloseRequest struct {
	ClosedReason string `json:"closedReason"`
	Comment      string `json:"comment,omitempty"`
}

// FeedbackDeleteResponse 删除反馈响应
type FeedbackDeleteResponse struct {
	Message string `json:"message"`
}

// ========== 反馈(Feedback) API 方法 ==========

// GetFeedbacks 获取反馈列表（支持分页）
func (c *Client) GetFeedbacks(page, limit int) (*FeedbackListResponse, error) {
	var result FeedbackListResponse
	path := fmt.Sprintf("/api.php/v1/feedbacks?page=%d&limit=%d", page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取反馈列表失败: %v", err)
	}
	return &result, nil
}

// GetFeedback 获取反馈详情
func (c *Client) GetFeedback(feedbackID int) (*Feedback, error) {
	var feedback Feedback
	path := fmt.Sprintf("/api.php/v1/feedbacks/%d", feedbackID)
	if err := c.doGet(path, &feedback); err != nil {
		return nil, fmt.Errorf("获取反馈详情失败: %v", err)
	}
	return &feedback, nil
}

// CreateFeedback 创建反馈
func (c *Client) CreateFeedback(req FeedbackCreateRequest) (*Feedback, error) {
	var feedback Feedback
	if err := c.doPost("/api.php/v1/feedbacks", req, &feedback); err != nil {
		return nil, fmt.Errorf("创建反馈失败: %v", err)
	}
	return &feedback, nil
}

// UpdateFeedback 修改反馈
func (c *Client) UpdateFeedback(feedbackID int, req FeedbackUpdateRequest) (*Feedback, error) {
	var feedback Feedback
	path := fmt.Sprintf("/api.php/v1/feedbacks/%d", feedbackID)
	if err := c.doPut(path, req, &feedback); err != nil {
		return nil, fmt.Errorf("修改反馈失败: %v", err)
	}
	return &feedback, nil
}

// AssignFeedback 指派反馈
func (c *Client) AssignFeedback(feedbackID int, req FeedbackAssignRequest) (*Feedback, error) {
	var feedback Feedback
	path := fmt.Sprintf("/api.php/v1/feedbacks/%d/assign", feedbackID)
	if err := c.doPost(path, req, &feedback); err != nil {
		return nil, fmt.Errorf("指派反馈失败: %v", err)
	}
	return &feedback, nil
}

// CloseFeedback 关闭反馈
func (c *Client) CloseFeedback(feedbackID int, req FeedbackCloseRequest) (*Feedback, error) {
	var feedback Feedback
	path := fmt.Sprintf("/api.php/v1/feedbacks/%d/close", feedbackID)
	if err := c.doPost(path, req, &feedback); err != nil {
		return nil, fmt.Errorf("关闭反馈失败: %v", err)
	}
	return &feedback, nil
}

// DeleteFeedback 删除反馈
func (c *Client) DeleteFeedback(feedbackID int) error {
	path := fmt.Sprintf("/api.php/v1/feedbacks/%d", feedbackID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除反馈失败: %v", err)
	}
	return nil
}
