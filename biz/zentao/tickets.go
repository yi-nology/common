package zentao

import "fmt"

// ========== 工单(Ticket)管理 ==========

// Ticket 工单结构
type Ticket struct {
	ID              int       `json:"id"`
	Product         int       `json:"product"`
	Module          int       `json:"module"`
	Title           string    `json:"title"`
	Type            string    `json:"type"`
	Desc            string    `json:"desc"`
	OpenedBuild     string    `json:"openedBuild"`
	Feedback        int       `json:"feedback"`
	AssignedTo      *string   `json:"assignedTo"`
	AssignedDate    string    `json:"assignedDate"`
	RealStarted     string    `json:"realStarted"`
	StartedBy       string    `json:"startedBy"`
	StartedDate     string    `json:"startedDate"`
	Deadline        *string   `json:"deadline"`
	Pri             int       `json:"pri"`
	Estimate        float64   `json:"estimate"`
	Left            float64   `json:"left"`
	Status          string    `json:"status"`
	OpenedBy        UserRef   `json:"openedBy"`
	OpenedDate      string    `json:"openedDate"`
	ActivatedCount  int       `json:"activatedCount"`
	ActivatedBy     *string   `json:"activatedBy"`
	ActivatedDate   *string   `json:"activatedDate"`
	ClosedBy        *string   `json:"closedBy"`
	ClosedDate      *string   `json:"closedDate"`
	ClosedReason    string    `json:"closedReason"`
	FinishedBy      *string   `json:"finishedBy"`
	FinishedDate    *string   `json:"finishedDate"`
	ResolvedBy      string    `json:"resolvedBy"`
	ResolvedDate    string    `json:"resolvedDate"`
	Resolution      string    `json:"resolution"`
	EditedBy        *string   `json:"editedBy"`
	EditedDate      *string   `json:"editedDate"`
	Keywords        string    `json:"keywords"`
	RepeatTicket    int       `json:"repeatTicket"`
	Mailto          []string  `json:"mailto"`
	Deleted         bool      `json:"deleted"`
	Consumed        float64   `json:"consumed"`
}

// TicketListResponse 工单列表响应
type TicketListResponse struct {
	Page    int      `json:"page"`
	Total   int      `json:"total"`
	Limit   int      `json:"limit"`
	Tickets []Ticket `json:"tickets"`
}

// TicketCreateRequest 创建工单请求
type TicketCreateRequest struct {
	Product int    `json:"product"`
	Module  int    `json:"module"`
	Title   string `json:"title"`
	Type    string `json:"type,omitempty"`
}

// TicketUpdateRequest 修改工单请求
type TicketUpdateRequest struct {
	Module  int    `json:"module,omitempty"`
	Title   string `json:"title,omitempty"`
	Type    string `json:"type,omitempty"`
	Desc    string `json:"desc,omitempty"`
	Pri     int    `json:"pri,omitempty"`
	AssignedTo string `json:"assignedTo,omitempty"`
	Deadline   string `json:"deadline,omitempty"`
	Estimate   float64 `json:"estimate,omitempty"`
}

// TicketDeleteResponse 删除工单响应
type TicketDeleteResponse struct {
	Message string `json:"message"`
}

// GetTickets 获取工单列表（支持分页和筛选）
func (c *Client) GetTickets(browseType string, param int, page, limit int) (*TicketListResponse, error) {
	var result TicketListResponse
	path := fmt.Sprintf("/api.php/v1/tickets?browseType=%s&param=%d&page=%d&limit=%d", browseType, param, page, limit)
	if err := c.doGet(path, &result); err != nil {
		return nil, fmt.Errorf("获取工单列表失败: %v", err)
	}
	return &result, nil
}

// GetTicket 获取工单详情
func (c *Client) GetTicket(ticketID int) (*Ticket, error) {
	var ticket Ticket
	path := fmt.Sprintf("/api.php/v1/tickets/%d", ticketID)
	if err := c.doGet(path, &ticket); err != nil {
		return nil, fmt.Errorf("获取工单详情失败: %v", err)
	}
	return &ticket, nil
}

// CreateTicket 创建工单
func (c *Client) CreateTicket(req TicketCreateRequest) (*Ticket, error) {
	var ticket Ticket
	if err := c.doPost("/api.php/v1/tickets", req, &ticket); err != nil {
		return nil, fmt.Errorf("创建工单失败: %v", err)
	}
	return &ticket, nil
}

// UpdateTicket 修改工单
func (c *Client) UpdateTicket(ticketID int, req TicketUpdateRequest) (*Ticket, error) {
	var ticket Ticket
	path := fmt.Sprintf("/api.php/v1/tickets/%d", ticketID)
	if err := c.doPut(path, req, &ticket); err != nil {
		return nil, fmt.Errorf("修改工单失败: %v", err)
	}
	return &ticket, nil
}

// DeleteTicket 删除工单
func (c *Client) DeleteTicket(ticketID int) error {
	path := fmt.Sprintf("/api.php/v1/tickets/%d", ticketID)
	if err := c.doDelete(path); err != nil {
		return fmt.Errorf("删除工单失败: %v", err)
	}
	return nil
}
