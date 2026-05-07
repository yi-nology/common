package zentao

import (
	"os"
	"testing"
)

// 测试配置 - 通过环境变量传入
// ZENTAO_BASE_URL - 禅道服务器地址
// ZENTAO_ACCOUNT - 登录账号
// ZENTAO_PASSWORD - 登录密码

func getTestConfig(t *testing.T) (baseURL, account, password string) {
	baseURL = os.Getenv("ZENTAO_BASE_URL")
	account = os.Getenv("ZENTAO_ACCOUNT")
	password = os.Getenv("ZENTAO_PASSWORD")

	if baseURL == "" {
		t.Skip("ZENTAO_BASE_URL 环境变量未设置，跳过集成测试")
	}
	if account == "" {
		t.Skip("ZENTAO_ACCOUNT 环境变量未设置，跳过集成测试")
	}
	if password == "" {
		t.Skip("ZENTAO_PASSWORD 环境变量未设置，跳过集成测试")
	}
	return
}

// TestGetToken 测试获取Token
func TestGetToken(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	token, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}
	if token == "" {
		t.Fatal("Token为空")
	}
	t.Logf("Token获取成功: %s...", token[:20])
}

// TestGetProducts 测试获取产品列表
func TestGetProducts(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	products, err := client.GetProducts(1, 20)
	if err != nil {
		t.Fatalf("获取产品列表失败: %v", err)
	}

	t.Logf("产品总数: %d", products.Total)
	for _, p := range products.Products {
		t.Logf("  - [%d] %s (%s)", p.ID, p.Name, p.Status)
	}
}

// TestGetBugs 测试获取Bug列表（分页）
func TestGetBugs(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	// 获取第一个产品
	products, err := client.GetProducts(1, 1)
	if err != nil {
		t.Fatalf("获取产品列表失败: %v", err)
	}
	if len(products.Products) == 0 {
		t.Skip("没有产品，跳过测试")
	}

	productID := products.Products[0].ID

	// 测试第一页
	page1, err := client.GetBugs(productID, 1, 5)
	if err != nil {
		t.Fatalf("获取Bug列表失败: %v", err)
	}

	t.Logf("产品 %d 的Bug列表 (第1页, 每页5条):", productID)
	t.Logf("  总数: %d, 当前页: %d, 每页: %d", page1.Total, page1.Page, page1.Limit)

	for i, bug := range page1.Bugs {
		t.Logf("  %d. [%d] %s (状态: %s, 严重程度: %d, 优先级: %d)",
			i+1, bug.ID, bug.Title, bug.Status, bug.Severity, bug.Pri)
	}

	// 如果Bug数量足够，测试第二页
	if page1.Total > 5 {
		page2, err := client.GetBugs(productID, 2, 5)
		if err != nil {
			t.Fatalf("获取Bug列表第二页失败: %v", err)
		}
		t.Logf("第二页: %d 条记录", len(page2.Bugs))
		if page2.Page != 2 {
			t.Errorf("分页参数错误: 期望 page=2, 实际 page=%d", page2.Page)
		}
	}
}

// TestGetBugDetail 测试获取Bug详情
func TestGetBugDetail(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	// 获取一个产品
	products, err := client.GetProducts(1, 1)
	if err != nil {
		t.Fatalf("获取产品列表失败: %v", err)
	}
	if len(products.Products) == 0 {
		t.Skip("没有产品，跳过测试")
	}

	productID := products.Products[0].ID

	// 获取Bug列表
	bugs, err := client.GetBugs(productID, 1, 1)
	if err != nil {
		t.Fatalf("获取Bug列表失败: %v", err)
	}
	if len(bugs.Bugs) == 0 {
		t.Skip("没有Bug，跳过测试")
	}

	bugID := bugs.Bugs[0].ID

	// 获取Bug详情
	bug, err := client.GetBug(bugID)
	if err != nil {
		t.Fatalf("获取Bug详情失败: %v", err)
	}

	t.Logf("Bug详情:")
	t.Logf("  ID: %d", bug.ID)
	t.Logf("  标题: %s", bug.Title)
	t.Logf("  状态: %s", bug.Status)
	t.Logf("  严重程度: %d", bug.Severity)
	t.Logf("  优先级: %d", bug.Pri)
	t.Logf("  类型: %s", bug.Type)
	t.Logf("  创建人: %+v", bug.OpenedBy)
	t.Logf("  指派给: %+v", bug.AssignedTo)
	t.Logf("  影响版本: %v", bug.OpenedBuild)
}

// TestGetProjects 测试获取项目列表
func TestGetProjects(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	projects, err := client.GetAllProjects(1, 20)
	if err != nil {
		t.Fatalf("获取项目列表失败: %v", err)
	}

	t.Logf("项目总数: %d", projects.Total)
	for _, p := range projects.Projects {
		t.Logf("  - [%d] %s (%s, 状态: %s)", p.ID, p.Name, p.Code, p.Status)
	}
}

// TestGetExecutions 测试获取执行列表
func TestGetExecutions(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	// 获取一个项目
	projects, err := client.GetAllProjects(1, 1)
	if err != nil {
		t.Fatalf("获取项目列表失败: %v", err)
	}
	if len(projects.Projects) == 0 {
		t.Skip("没有项目，跳过测试")
	}

	projectID := projects.Projects[0].ID

	// 获取执行的执行列表
	executions, err := client.GetExecutions(projectID, 1, 20)
	if err != nil {
		t.Fatalf("获取执行列表失败: %v", err)
	}

	t.Logf("项目 %d 的执行/迭代列表:", projectID)
	t.Logf("  总数: %d", executions.Total)
	for _, e := range executions.Executions {
		t.Logf("  - [%d] %s (%s ~ %s, 状态: %s)", e.ID, e.Name, e.Begin, e.End, e.Status)
	}
}

// TestGetTasks 测试获取任务列表（分页）
func TestGetTasks(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	// 获取一个项目
	projects, err := client.GetAllProjects(1, 1)
	if err != nil {
		t.Fatalf("获取项目列表失败: %v", err)
	}
	if len(projects.Projects) == 0 {
		t.Skip("没有项目，跳过测试")
	}

	projectID := projects.Projects[0].ID

	// 获取执行列表
	executions, err := client.GetExecutions(projectID, 1, 1)
	if err != nil {
		t.Fatalf("获取执行列表失败: %v", err)
	}
	if len(executions.Executions) == 0 {
		t.Skip("没有执行，跳过测试")
	}

	executionID := executions.Executions[0].ID

	// 获取任务列表（第一页）
	tasks, err := client.GetTasks(executionID, 1, 10)
	if err != nil {
		t.Fatalf("获取任务列表失败: %v", err)
	}

	t.Logf("执行 %d 的任务列表:", executionID)
	t.Logf("  总数: %d, 当前页: %d, 每页: %d", tasks.Total, tasks.Page, tasks.Limit)

	for i, task := range tasks.Tasks {
		t.Logf("  %d. [%d] %s (状态: %s, 指派给: %s)",
			i+1, task.ID, task.Name, task.Status, task.AssignedTo.Realname)
	}

	// 测试分页
	if tasks.Total > 10 {
		page2, err := client.GetTasks(executionID, 2, 10)
		if err != nil {
			t.Fatalf("获取任务列表第二页失败: %v", err)
		}
		t.Logf("第二页: %d 条记录, page=%d", len(page2.Tasks), page2.Page)
	}
}

// TestGetStoriesByProduct 测试获取产品需求列表
func TestGetStoriesByProduct(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	// 获取一个产品
	products, err := client.GetProducts(1, 1)
	if err != nil {
		t.Fatalf("获取产品列表失败: %v", err)
	}
	if len(products.Products) == 0 {
		t.Skip("没有产品，跳过测试")
	}

	productID := products.Products[0].ID

	// 获取需求列表
	stories, err := client.GetStoriesByProduct(productID, 1, 10)
	if err != nil {
		t.Fatalf("获取需求列表失败: %v", err)
	}

	t.Logf("产品 %d 的需求列表:", productID)
	t.Logf("  总数: %d, 当前页: %d, 每页: %d", stories.Total, stories.Page, stories.Limit)

	for i, s := range stories.Stories {
		t.Logf("  %d. [%d] %s (状态: %s, 优先级: %d)", i+1, s.ID, s.Title, s.Status, s.Pri)
	}
}

// TestGetCasesByProduct 测试获取用例列表
func TestGetCasesByProduct(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	// 获取一个产品
	products, err := client.GetProducts(1, 1)
	if err != nil {
		t.Fatalf("获取产品列表失败: %v", err)
	}
	if len(products.Products) == 0 {
		t.Skip("没有产品，跳过测试")
	}

	productID := products.Products[0].ID

	// 获取用例列表
	cases, err := client.GetCasesByProduct(productID, 1, 10)
	if err != nil {
		t.Fatalf("获取用例列表失败: %v", err)
	}

	t.Logf("产品 %d 的用例列表:", productID)
	t.Logf("  总数: %d", cases.Total)

	for i, c := range cases.Testcases {
		t.Logf("  %d. [%d] %s (类型: %s, 状态: %s)", i+1, c.ID, c.Title, c.Type, c.Status)
	}
}

// TestGetUsers 测试获取用户列表
func TestGetUsers(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	users, err := client.GetUsers(1, 20)
	if err != nil {
		t.Fatalf("获取用户列表失败: %v", err)
	}

	t.Logf("用户总数: %d", users.Total)
	for _, u := range users.Users {
		t.Logf("  - [%d] %s (%s, 角色: %s)", u.ID, u.Realname, u.Account, u.Role)
	}
}

// TestGetTestTasks 测试获取测试单列表
func TestGetTestTasks(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	testTasks, err := client.GetTestTasks(1, 10)
	if err != nil {
		t.Fatalf("获取测试单列表失败: %v", err)
	}

	t.Logf("测试单总数: %d", testTasks.Total)
	for _, tt := range testTasks.Testtasks {
		t.Logf("  - [%d] %s (状态: %s)", tt.ID, tt.Name, tt.Status)
	}
}

// TestPagination 测试分页功能
func TestPagination(t *testing.T) {
	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	// 测试产品列表分页
	t.Run("Products", func(t *testing.T) {
		page1, err := client.GetProducts(1, 3)
		if err != nil {
			t.Fatalf("获取产品列表失败: %v", err)
		}

		t.Logf("产品列表分页测试:")
		t.Logf("  第1页: %d 条 (limit=3)", len(page1.Products))
		t.Logf("  总数: %d", page1.Total)

		if page1.Total > 3 {
			page2, err := client.GetProducts(2, 3)
			if err != nil {
				t.Fatalf("获取产品列表第2页失败: %v", err)
			}
			t.Logf("  第2页: %d 条", len(page2.Products))

			// 验证分页正确性
			if page2.Page != 2 {
				t.Errorf("页码错误: 期望 2, 实际 %d", page2.Page)
			}
			if page1.Limit != 3 || page2.Limit != 3 {
				t.Errorf("每页数量错误: page1.Limit=%d, page2.Limit=%d", page1.Limit, page2.Limit)
			}
		}
	})

	// 测试项目列表分页
	t.Run("Projects", func(t *testing.T) {
		page1, err := client.GetAllProjects(1, 5)
		if err != nil {
			t.Fatalf("获取项目列表失败: %v", err)
		}

		t.Logf("项目列表分页测试:")
		t.Logf("  第1页: %d 条 (limit=5)", len(page1.Projects))
		t.Logf("  总数: %d", page1.Total)

		if page1.Page != 1 {
			t.Errorf("页码错误: 期望 1, 实际 %d", page1.Page)
		}
		if page1.Limit != 5 {
			t.Errorf("每页数量错误: 期望 5, 实际 %d", page1.Limit)
		}
	})
}

// TestCreateAndDeleteStory 测试创建和删除需求（可选测试）
// 注意：此测试会创建真实的需求，默认跳过
func TestCreateAndDeleteStory(t *testing.T) {
	// 设置环境变量 ZENTAO_RUN_CREATE_TESTS=1 来运行此测试
	if os.Getenv("ZENTAO_RUN_CREATE_TESTS") != "1" {
		t.Skip("跳过创建测试，设置 ZENTAO_RUN_CREATE_TESTS=1 来运行")
	}

	baseURL, account, password := getTestConfig(t)

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}

	// 获取一个产品
	products, err := client.GetProducts(1, 1)
	if err != nil {
		t.Fatalf("获取产品列表失败: %v", err)
	}
	if len(products.Products) == 0 {
		t.Skip("没有产品，跳过测试")
	}

	productID := products.Products[0].ID

	// 创建需求
	createReq := StoryCreateRequest{
		Title:    "【单元测试】测试需求标题",
		Product:  productID,
		Pri:      2,
		Category: "feature",
		Spec:     "【单元测试】需求描述",
	}

	story, err := client.CreateStory(createReq)
	if err != nil {
		t.Fatalf("创建需求失败: %v", err)
	}

	t.Logf("创建需求成功: ID=%d, 标题=%s", story.ID, story.Title)

	// 清理：删除创建的需求
	err = client.DeleteStory(story.ID)
	if err != nil {
		t.Logf("警告：删除需求失败: %v (需要手动删除 Story #%d)", err, story.ID)
	} else {
		t.Logf("删除需求成功: ID=%d", story.ID)
	}
}
