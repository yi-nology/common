package zentao

import (
	"fmt"
	"os"
	"testing"
)

func getTestClient(t *testing.T) *Client {
	baseURL := os.Getenv("ZENTAO_URL")
	account := os.Getenv("ZENTAO_ACCOUNT")
	password := os.Getenv("ZENTAO_PASSWORD")

	if baseURL == "" || account == "" || password == "" {
		t.Skip("跳过测试: 请设置环境变量 ZENTAO_URL, ZENTAO_ACCOUNT, ZENTAO_PASSWORD")
	}

	client := NewClient(baseURL)
	_, err := client.GetToken(account, password)
	if err != nil {
		t.Fatalf("获取Token失败: %v", err)
	}
	return client
}

func TestGetPrograms(t *testing.T) {
	client := getTestClient(t)
	programs, err := client.GetPrograms()
	if err != nil {
		t.Fatalf("获取项目集列表失败: %v", err)
	}
	fmt.Printf("项目集数量: %d\n", len(programs))
	for i, p := range programs {
		if i >= 3 {
			fmt.Println("...")
			break
		}
		fmt.Printf("  [%d] %s (ID: %d)\n", i+1, p.Name, p.ID)
	}
}

func TestGetProducts(t *testing.T) {
	client := getTestClient(t)
	products, err := client.GetProducts()
	if err != nil {
		t.Fatalf("获取产品列表失败: %v", err)
	}
	fmt.Printf("产品数量: %d\n", len(products))
	for i, p := range products {
		if i >= 3 {
			fmt.Println("...")
			break
		}
		fmt.Printf("  [%d] %s (ID: %d)\n", i+1, p.Name, p.ID)
	}
}

func TestGetProjects(t *testing.T) {
	client := getTestClient(t)
	projects, err := client.GetAllProjects(10)
	if err != nil {
		t.Fatalf("获取项目列表失败: %v", err)
	}
	fmt.Printf("项目数量: %d\n", len(projects))
	for i, p := range projects {
		if i >= 3 {
			fmt.Println("...")
			break
		}
		fmt.Printf("  [%d] %s (ID: %d)\n", i+1, p.Name, p.ID)
	}
}

func TestGetExecution(t *testing.T) {
	client := getTestClient(t)
	projects, err := client.GetAllProjects(10)
	if err != nil {
		t.Fatalf("获取项目列表失败: %v", err)
	}
	if len(projects) == 0 {
		t.Skip("没有项目可测试")
	}

	for _, proj := range projects {
		executions, err := client.GetExecutions(proj.ID)
		if err != nil {
			continue
		}
		if len(executions) > 0 {
			exec, err := client.GetExecution(executions[0].ID)
			if err != nil {
				t.Fatalf("获取执行详情失败: %v", err)
			}
			fmt.Printf("执行详情: ID=%d, Name=%s, Status=%s\n", exec.ID, exec.Name, exec.Status)
			return
		}
	}
	t.Skip("没有找到可用的执行")
}

func TestGetUsers(t *testing.T) {
	client := getTestClient(t)
	users, err := client.GetUsers()
	if err != nil {
		t.Fatalf("获取用户列表失败: %v", err)
	}
	fmt.Printf("用户数量: %d\n", len(users))
	for i, u := range users {
		if i >= 3 {
			fmt.Println("...")
			break
		}
		fmt.Printf("  [%d] %s (%s)\n", i+1, u.Realname, u.Account)
	}
}
