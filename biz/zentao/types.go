package main
package zentao

// ========== 认证相关 ==========

// TokenRequest 获取token的请求结构
type TokenRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// TokenResponse 获取token的响应结构
type TokenResponse struct {
	Token string `json:"token"`
}

// ========== 用户相关 ==========

// User 用户结构
type User struct {
	ID       int    `json:"id"`
	Account  string `json:"account"`
	Realname string `json:"realname"`
	Role     string `json:"role"`
	Dept     int    `json:"dept"`
	Email    string `json:"email"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Page  int    `json:"page"`
	Total int    `json:"total"`
	Limit int    `json:"limit"`
	Users []User `json:"users"`
}

// ========== 项目集(Program)相关 ==========

// Program 项目集结构
type Program struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Code     string      `json:"code"`
	Parent   int         `json:"parent"`
	Type     string      `json:"type"`
	Status   string      `json:"status"`
	Begin    string      `json:"begin"`
	End      string      `json:"end"`
	Desc     string      `json:"desc"`
	PM       interface{} `json:"PM"`
	OpenedBy interface{} `json:"openedBy"`
}

// ProgramListResponse 项目集列表响应
type ProgramListResponse struct {
	Page     int       `json:"page"`
	Total    int       `json:"total"`
	Limit    int       `json:"limit"`
	Programs []Program `json:"programs"`
}

// ProgramCreateRequest 创建项目集请求
type ProgramCreateRequest struct {
	Name   string `json:"name"`
	Code   string `json:"code,omitempty"`
	Parent int    `json:"parent,omitempty"`
	Begin  string `json:"begin,omitempty"`
	End    string `json:"end,omitempty"`
	Desc   string `json:"desc,omitempty"`
	PM     string `json:"PM,omitempty"`
}

// ========== 产品(Product)相关 ==========

// Product 产品结构
type Product struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Desc   string `json:"desc"`
}

// ProductListResponse 产品列表响应
type ProductListResponse struct {
	Page     int       `json:"page"`
	Total    int       `json:"total"`
	Limit    int       `json:"limit"`
	Products []Product `json:"products"`
}

// ProductCreateRequest 创建/更新产品请求
type ProductCreateRequest struct {
	Name   string `json:"name"`
	Code   string `json:"code,omitempty"`
	Type   string `json:"type,omitempty"`
	Desc   string `json:"desc,omitempty"`
	PO     string `json:"PO,omitempty"`
	QD     string `json:"QD,omitempty"`
	RD     string `json:"RD,omitempty"`
	Status string `json:"status,omitempty"`
}

// ========== 项目(Project)相关 ==========

// Project 项目结构
type Project struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Code     string      `json:"code"`
	Model    string      `json:"model"`
	Type     string      `json:"type"`
	Status   string      `json:"status"`
	PM       interface{} `json:"PM"`
	Begin    string      `json:"begin"`
	End      string      `json:"end"`
	Progress interface{} `json:"progress"`
}

// ProjectListResponse 项目列表响应
type ProjectListResponse struct {
	Page     int       `json:"page"`
	Total    int       `json:"total"`
	Limit    int       `json:"limit"`
	Projects []Project `json:"projects"`
}

// ProjectCreateRequest 创建项目请求
type ProjectCreateRequest struct {
	Name   string `json:"name"`
	Code   string `json:"code,omitempty"`
	Model  string `json:"model,omitempty"`
	Type   string `json:"type,omitempty"`
	Begin  string `json:"begin,omitempty"`
	End    string `json:"end,omitempty"`
	Desc   string `json:"desc,omitempty"`
	PM     string `json:"PM,omitempty"`
	Parent int    `json:"parent,omitempty"`
	Status string `json:"status,omitempty"`
}

// ========== 执行/迭代(Execution)相关 ==========

// Execution 执行/迭代结构
type Execution struct {
	ID        int    `json:"id"`
	Project   int    `json:"project"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Begin     string `json:"begin"`
	End       string `json:"end"`
	RealBegan string `json:"realBegan"`
	RealEnd   string `json:"realEnd"`
	Desc      string `json:"desc"`
}

// ExecutionListResponse 执行列表响应
type ExecutionListResponse struct {
	Page       int         `json:"page"`
	Total      int         `json:"total"`
	Limit      int         `json:"limit"`
	Executions []Execution `json:"executions"`
}

// ExecutionCreateRequest 创建执行请求
type ExecutionCreateRequest struct {
	Name   string `json:"name"`
	Code   string `json:"code,omitempty"`
	Type   string `json:"type,omitempty"`
	Begin  string `json:"begin,omitempty"`
	End    string `json:"end,omitempty"`
	Desc   string `json:"desc,omitempty"`
	PM     string `json:"PM,omitempty"`
	Status string `json:"status,omitempty"`
}

// ========== 任务(Task)相关 ==========

// UserRef 用户引用结构
type UserRef struct {
	ID       int    `json:"id"`
	Account  string `json:"account"`
	Avatar   string `json:"avatar"`
	Realname string `json:"realname"`
}

// Task 任务结构
type Task struct {
	ID           int         `json:"id"`
	Project      int         `json:"project"`
	Execution    int         `json:"execution"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Pri          int         `json:"pri"`
	Status       string      `json:"status"`
	AssignedTo   UserRef     `json:"assignedTo"`
	EstStarted   string      `json:"estStarted"`
	Deadline     string      `json:"deadline"`
	Estimate     float64     `json:"estimate"`
	Consumed     float64     `json:"consumed"`
	Left         float64     `json:"left"`
	Desc         string      `json:"desc"`
	OpenedBy     interface{} `json:"openedBy"`
	OpenedDate   string      `json:"openedDate"`
	FinishedBy   interface{} `json:"finishedBy"`
	FinishedDate interface{} `json:"finishedDate"`
	ClosedBy     interface{} `json:"closedBy"`
	ClosedDate   interface{} `json:"closedDate"`
	StatusName   string      `json:"statusName"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Page  int    `json:"page"`
	Total int    `json:"total"`
	Limit int    `json:"limit"`
	Tasks []Task `json:"tasks"`
}

// TaskCreateRequest 创建任务请求
type TaskCreateRequest struct {
	Name       string  `json:"name"`
	Type       string  `json:"type,omitempty"`
	AssignedTo string  `json:"assignedTo,omitempty"`
	EstStarted string  `json:"estStarted,omitempty"`
	Deadline   string  `json:"deadline,omitempty"`
	Estimate   float64 `json:"estimate,omitempty"`
	Desc       string  `json:"desc,omitempty"`
	Pri        int     `json:"pri,omitempty"`
}

// TaskUpdateRequest 更新任务请求
type TaskUpdateRequest struct {
	Name       string  `json:"name,omitempty"`
	Type       string  `json:"type,omitempty"`
	AssignedTo string  `json:"assignedTo,omitempty"`
	EstStarted string  `json:"estStarted,omitempty"`
	Deadline   string  `json:"deadline,omitempty"`
	Estimate   float64 `json:"estimate,omitempty"`
	Consumed   float64 `json:"consumed,omitempty"`
	Left       float64 `json:"left,omitempty"`
	Desc       string  `json:"desc,omitempty"`
	Pri        int     `json:"pri,omitempty"`
}

// TaskStartRequest 开始任务请求
type TaskStartRequest struct {
	RealStarted string  `json:"realStarted,omitempty"`
	Consumed    float64 `json:"consumed,omitempty"`
	Left        float64 `json:"left,omitempty"`
	Comment     string  `json:"comment,omitempty"`
}

// TaskFinishRequest 完成任务请求
type TaskFinishRequest struct {
	Consumed     float64 `json:"consumed,omitempty"`
	FinishedDate string  `json:"finishedDate,omitempty"`
	Comment      string  `json:"comment,omitempty"`
}

// TaskPauseRequest 暂停任务请求
type TaskPauseRequest struct {
	Comment string `json:"comment,omitempty"`
}

// TaskAssignRequest 指派任务请求
type TaskAssignRequest struct {
	AssignedTo string  `json:"assignedTo"`
	Left       float64 `json:"left,omitempty"`
	Comment    string  `json:"comment,omitempty"`
}

// ========== Bug相关 ==========

// Bug Bug结构
type Bug struct {
	ID            int         `json:"id"`
	Project       int         `json:"project"`
	Product       int         `json:"product"`
	Title         string      `json:"title"`
	Keywords      string      `json:"keywords"`
	Severity      int         `json:"severity"`
	Pri           int         `json:"pri"`
	Type          string      `json:"type"`
	OS            string      `json:"os"`
	Browser       string      `json:"browser"`
	Hardware      string      `json:"hardware"`
	Steps         string      `json:"steps"`
	Status        string      `json:"status"`
	SubStatus     string      `json:"subStatus"`
	Color         string      `json:"color"`
	Confirmed     int         `json:"confirmed"`
	PlanTime      string      `json:"planTime"`
	OpenedBy      UserRef     `json:"openedBy"`
	OpenedDate    string      `json:"openedDate"`
	OpenedBuild   string      `json:"openedBuild"`
	AssignedTo    UserRef     `json:"assignedTo"`
	AssignedDate  string      `json:"assignedDate"`
	Deadline      interface{} `json:"deadline"`
	ResolvedBy    interface{} `json:"resolvedBy"`
	Resolution    string      `json:"resolution"`
	ResolvedBuild string      `json:"resolvedBuild"`
	ResolvedDate  interface{} `json:"resolvedDate"`
	ClosedBy      interface{} `json:"closedBy"`
	ClosedDate    interface{} `json:"closedDate"`
	StatusName    string      `json:"statusName"`
	LifeCycle     string      `json:"lifeCycle"`
}

// BugListResponse Bug列表响应
type BugListResponse struct {
	Page  int   `json:"page"`
	Total int   `json:"total"`
	Limit int   `json:"limit"`
	Bugs  []Bug `json:"bugs"`
}

// BugSearchParams Bug搜索参数
type BugSearchParams struct {
	ProductID  int
	Status     string
	AssignedTo string
	Keyword    string
	Severity   int
	Pri        int
	Limit      int
	Page       int
}

// BugConfirmRequest 确认Bug请求
type BugConfirmRequest struct {
	Comment string `json:"comment,omitempty"`
}

// BugResolveRequest 解决Bug请求
type BugResolveRequest struct {
	Resolution    string `json:"resolution"`
	ResolvedBuild string `json:"resolvedBuild,omitempty"`
	Comment       string `json:"comment,omitempty"`
}

// BugCloseRequest 关闭Bug请求
type BugCloseRequest struct {
	Comment string `json:"comment,omitempty"`
}

// BugActivateRequest 激活Bug请求
type BugActivateRequest struct {
	AssignedTo string `json:"assignedTo,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// BugAssignRequest 指派Bug请求
type BugAssignRequest struct {
	AssignedTo string `json:"assignedTo"`
	Comment    string `json:"comment,omitempty"`
}

// ========== 需求(Story)相关 ==========

// Story 需求结构
type Story struct {
	ID           int     `json:"id"`
	Product      int     `json:"product"`
	Module       int     `json:"module"`
	Plan         int     `json:"plan"`
	Source       string  `json:"source"`
	Title        string  `json:"title"`
	Spec         string  `json:"spec"`
	Verify       string  `json:"verify"`
	Type         string  `json:"type"`
	Status       string  `json:"status"`
	Stage        string  `json:"stage"`
	Pri          int     `json:"pri"`
	Estimate     float64 `json:"estimate"`
	Version      int     `json:"version"`
	OpenedBy     string  `json:"openedBy"`
	OpenedDate   string  `json:"openedDate"`
	AssignedTo   string  `json:"assignedTo"`
	AssignedDate string  `json:"assignedDate"`
	ClosedBy     string  `json:"closedBy"`
	ClosedDate   string  `json:"closedDate"`
	ClosedReason string  `json:"closedReason"`
}

// StoryListResponse 需求列表响应
type StoryListResponse struct {
	Page    int     `json:"page"`
	Total   int     `json:"total"`
	Limit   int     `json:"limit"`
	Stories []Story `json:"stories"`
}

// StoryCreateRequest 创建需求请求
type StoryCreateRequest struct {
	Title      string  `json:"title"`
	Spec       string  `json:"spec,omitempty"`
	Verify     string  `json:"verify,omitempty"`
	Type       string  `json:"type,omitempty"`
	Pri        int     `json:"pri,omitempty"`
	Estimate   float64 `json:"estimate,omitempty"`
	AssignedTo string  `json:"assignedTo,omitempty"`
	Module     int     `json:"module,omitempty"`
	Plan       int     `json:"plan,omitempty"`
	Source     string  `json:"source,omitempty"`
}

// ========== 计划(Plan)相关 ==========

// Plan 计划结构
type Plan struct {
	ID       int    `json:"id"`
	Product  int    `json:"product"`
	Parent   int    `json:"parent"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Begin    string `json:"begin"`
	End      string `json:"end"`
	Status   string `json:"status"`
	ClosedBy string `json:"closedBy"`
}

// PlanListResponse 计划列表响应
type PlanListResponse struct {
	Page  int    `json:"page"`
	Total int    `json:"total"`
	Limit int    `json:"limit"`
	Plans []Plan `json:"plans"`
}

// PlanCreateRequest 创建计划请求
type PlanCreateRequest struct {
	Title  string `json:"title"`
	Begin  string `json:"begin,omitempty"`
	End    string `json:"end,omitempty"`
	Desc   string `json:"desc,omitempty"`
	Parent int    `json:"parent,omitempty"`
}

// ========== 发布(Release)相关 ==========

// Release 发布结构
type Release struct {
	ID        int    `json:"id"`
	Product   int    `json:"product"`
	Build     int    `json:"build"`
	Name      string `json:"name"`
	Marker    string `json:"marker"`
	Date      string `json:"date"`
	Stories   string `json:"stories"`
	Bugs      string `json:"bugs"`
	Desc      string `json:"desc"`
	Status    string `json:"status"`
	SubStatus string `json:"subStatus"`
}

// ReleaseListResponse 发布列表响应
type ReleaseListResponse struct {
	Page     int       `json:"page"`
	Total    int       `json:"total"`
	Limit    int       `json:"limit"`
	Releases []Release `json:"releases"`
}

// ========== 版本(Build)相关 ==========

// Build 版本结构
type Build struct {
	ID        int    `json:"id"`
	Product   int    `json:"product"`
	Project   int    `json:"project"`
	Execution int    `json:"execution"`
	Name      string `json:"name"`
	ScmPath   string `json:"scmPath"`
	FilePath  string `json:"filePath"`
	Date      string `json:"date"`
	Stories   string `json:"stories"`
	Bugs      string `json:"bugs"`
	Builder   string `json:"builder"`
	Desc      string `json:"desc"`
	Deleted   string `json:"deleted"`
}

// BuildListResponse 版本列表响应
type BuildListResponse struct {
	Page   int     `json:"page"`
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Builds []Build `json:"builds"`
}

// BuildCreateRequest 创建版本请求
type BuildCreateRequest struct {
	Name      string `json:"name"`
	Execution int    `json:"execution,omitempty"`
	ScmPath   string `json:"scmPath,omitempty"`
	FilePath  string `json:"filePath,omitempty"`
	Date      string `json:"date,omitempty"`
	Desc      string `json:"desc,omitempty"`
}

// ========== 工时(Effort)相关 ==========

// EffortEntry 工时日志条目
type EffortEntry struct {
	ID         int     `json:"id"`
	ObjectType string  `json:"objectType"`
	ObjectID   int     `json:"objectID"`
	Product    string  `json:"product"`
	Project    int     `json:"project"`
	Execution  int     `json:"execution"`
	Account    string  `json:"account"`
	Work       string  `json:"work"`
	Date       string  `json:"date"`
	Left       float64 `json:"left"`
	Consumed   float64 `json:"consumed"`
	Begin      string  `json:"begin"`
	End        string  `json:"end"`
	Deleted    string  `json:"deleted"`
}

// EffortRequest 记录工时请求
type EffortRequest struct {
	ID         []int     `json:"id"`
	ObjectID   []int     `json:"objectID"`
	Dates      []string  `json:"dates"`
	Work       []string  `json:"work"`
	Consumed   []float64 `json:"consumed"`
	Left       []float64 `json:"left"`
	ObjectType []string  `json:"objectType"`
}
