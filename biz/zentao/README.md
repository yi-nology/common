# Zentao SDK

禅道(Zentao)项目管理系统 Go SDK，提供完整的 REST API v1 封装。

## 安装

```go
import "github.com/yi-nology/common/biz/zentao"
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/yi-nology/common/biz/zentao"
)

func main() {
    // 创建客户端
    client := zentao.NewClient("https://your-zentao-server.com")
    
    // 认证获取Token
    token, err := client.GetToken("account", "password")
    if err != nil {
        panic(err)
    }
    fmt.Println("Token:", token)
    
    // 获取产品列表
    products, err := client.GetProducts()
    if err != nil {
        panic(err)
    }
    for _, p := range products {
        fmt.Printf("产品: %s (ID: %d)\n", p.Name, p.ID)
    }
}
```

## API 列表

### 认证
| 方法 | 说明 |
|------|------|
| `GetToken(account, password)` | 获取认证Token |
| `SetToken(token)` | 设置Token |
| `SetTimeout(duration)` | 设置请求超时 |

### 项目集 (Programs)
| 方法 | 说明 |
|------|------|
| `GetPrograms()` | 获取项目集列表 |
| `GetProgram(id)` | 获取项目集详情 |
| `CreateProgram(req)` | 创建项目集 |
| `UpdateProgram(id, req)` | 更新项目集 |
| `DeleteProgram(id)` | 删除项目集 |

### 产品 (Products)
| 方法 | 说明 |
|------|------|
| `GetProducts()` | 获取产品列表 |
| `GetProduct(id)` | 获取产品详情 |
| `CreateProduct(req)` | 创建产品 |
| `UpdateProduct(id, req)` | 更新产品 |
| `DeleteProduct(id)` | 删除产品 |

### 项目 (Projects)
| 方法 | 说明 |
|------|------|
| `GetAllProjects(limit)` | 获取所有项目 |
| `GetProjectsByProduct(productID)` | 获取产品关联的项目 |
| `GetProject(id)` | 获取项目详情 |
| `CreateProject(req)` | 创建项目 |
| `UpdateProject(id, req)` | 更新项目 |
| `DeleteProject(id)` | 删除项目 |

### 执行/迭代 (Executions)
| 方法 | 说明 |
|------|------|
| `GetExecutions(projectID)` | 获取执行列表 |
| `GetExecution(id)` | 获取执行详情 |
| `CreateExecution(projectID, req)` | 创建执行 |
| `UpdateExecution(id, req)` | 更新执行 |
| `DeleteExecution(id)` | 删除执行 |

### 任务 (Tasks)
| 方法 | 说明 |
|------|------|
| `GetTasks(executionID, limit)` | 获取任务列表(支持分页) |
| `GetTask(id)` | 获取任务详情 |
| `CreateTask(executionID, req)` | 创建任务 |
| `UpdateTask(id, req)` | 更新任务 |
| `DeleteTask(id)` | 删除任务 |
| `StartTask(id, req)` | 开始任务 |
| `FinishTask(id, req)` | 完成任务 |
| `PauseTask(id, req)` | 暂停任务 |
| `ActivateTask(id, consumed, left)` | 激活任务 |
| `AssignTask(id, req)` | 指派任务 |

### Bug
| 方法 | 说明 |
|------|------|
| `GetBugs(productID, limit)` | 获取Bug列表 |
| `GetBug(id)` | 获取Bug详情 |
| `GetBugsByProject(productID, projectID, limit)` | 按项目过滤Bug |
| `GetBugsByStatus(productID, status, limit)` | 按状态过滤Bug |
| `GetActiveBugs(productID, limit)` | 获取激活状态的Bug |
| `SearchBugs(params)` | 搜索Bug(多条件) |
| `ConfirmBug(id, req)` | 确认Bug |
| `ResolveBug(id, req)` | 解决Bug |
| `CloseBug(id, req)` | 关闭Bug |
| `ActivateBug(id, req)` | 激活Bug |
| `AssignBug(id, req)` | 指派Bug |

### 需求 (Stories)
| 方法 | 说明 |
|------|------|
| `GetStoriesByProduct(productID, limit)` | 获取产品需求 |
| `GetStoriesByProject(projectID, limit)` | 获取项目需求 |
| `GetStoriesByExecution(executionID, limit)` | 获取执行需求 |
| `GetStory(id)` | 获取需求详情 |
| `CreateStory(productID, req)` | 创建需求 |
| `UpdateStory(id, req)` | 更新需求 |
| `DeleteStory(id)` | 删除需求 |
| `ChangeStory(id, spec, verify)` | 变更需求 |

### 计划 (Plans)
| 方法 | 说明 |
|------|------|
| `GetPlans(productID)` | 获取计划列表 |
| `GetPlan(id)` | 获取计划详情 |
| `CreatePlan(productID, req)` | 创建计划 |
| `UpdatePlan(id, req)` | 更新计划 |
| `DeletePlan(id)` | 删除计划 |
| `LinkStoriesToPlan(planID, storyIDs)` | 关联需求到计划 |
| `UnlinkStoriesFromPlan(planID, storyIDs)` | 取消关联需求 |
| `LinkBugsToPlan(planID, bugIDs)` | 关联Bug到计划 |
| `UnlinkBugsFromPlan(planID, bugIDs)` | 取消关联Bug |

### 发布 (Releases)
| 方法 | 说明 |
|------|------|
| `GetReleasesByProduct(productID)` | 获取产品发布列表 |
| `GetReleasesByProject(projectID)` | 获取项目发布列表 |
| `GetRelease(id)` | 获取发布详情 |

### 版本 (Builds)
| 方法 | 说明 |
|------|------|
| `GetBuildsByProject(projectID)` | 获取项目版本列表 |
| `GetBuildsByExecution(executionID)` | 获取执行版本列表 |
| `GetBuild(id)` | 获取版本详情 |
| `CreateBuild(projectID, req)` | 创建版本 |
| `UpdateBuild(id, req)` | 更新版本 |
| `DeleteBuild(id)` | 删除版本 |

### 用户 (Users)
| 方法 | 说明 |
|------|------|
| `GetUsers()` | 获取用户列表 |
| `GetCurrentUser()` | 获取当前登录用户 |
| `GetUserByID(id)` | 根据ID获取用户 |

### 工时 (Effort)
| 方法 | 说明 |
|------|------|
| `RecordEffort(taskID, date, consumed, left, work)` | 记录工时 |
| `GetTaskEfforts(taskID)` | 获取任务工时日志 |

## 测试

运行测试需要设置环境变量:

```bash
export ZENTAO_URL="https://your-zentao-server.com"
export ZENTAO_ACCOUNT="your_account"
export ZENTAO_PASSWORD="your_password"

go test -v ./...
```

## 文件结构

```
zentao/
├── client.go       # 基础客户端
├── types.go        # 类型定义
├── programs.go     # 项目集API
├── products.go     # 产品API
├── projects.go     # 项目API
├── executions.go   # 执行API
├── tasks.go        # 任务API
├── bugs.go         # Bug API
├── stories.go      # 需求API
├── plans.go        # 计划API
├── releases.go     # 发布API
├── builds.go       # 版本API
├── users.go        # 用户API
├── effort.go       # 工时API
└── client_test.go  # 测试文件
```

## 依赖

- `github.com/imroc/req/v3` - HTTP客户端
