package zentao

import (
	"encoding/json"
	"fmt"
	"time"
)

// ========== 工时(Effort)管理 ==========

// RecordEffort 记录任务工时
// date 参数格式为 YYYY-MM-DD，若为空则使用当天日期
func (c *Client) RecordEffort(taskID int, date string, consumed float64, left float64, work string) error {
	if c.Token == "" {
		return fmt.Errorf("token为空，请先调用GetToken获取token")
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	effortBody := EffortRequest{
		ID:         []int{0},
		ObjectID:   []int{taskID},
		Dates:      []string{date},
		Consumed:   []float64{consumed},
		Left:       []float64{left},
		Work:       []string{work},
		ObjectType: []string{"task"},
	}

	path := fmt.Sprintf("/api.php/v1/tasks/%d/estimate", taskID)
	if err := c.doPost(path, effortBody, nil); err != nil {
		return fmt.Errorf("记录工时失败: %v", err)
	}

	return nil
}

// GetTaskEfforts 获取任务的工时日志列表
func (c *Client) GetTaskEfforts(taskID int) ([]EffortEntry, error) {
	path := fmt.Sprintf("/api.php/v1/tasks/%d/estimate", taskID)
	data, err := c.getRawBytes(path)
	if err != nil {
		return nil, fmt.Errorf("获取工时日志失败: %v", err)
	}

	// 禅道返回格式: {"effort": {"1": {...}, "2": {...}}} 或 {"effort": []}
	var rawResp struct {
		Effort json.RawMessage `json:"effort"`
	}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("解析工时日志响应失败: %v", err)
	}

	var entries []EffortEntry
	if len(rawResp.Effort) == 0 || string(rawResp.Effort) == "[]" || string(rawResp.Effort) == "null" {
		return entries, nil
	}

	// 尝试解析为 map[string]EffortEntry
	var effortMap map[string]EffortEntry
	if err := json.Unmarshal(rawResp.Effort, &effortMap); err != nil {
		// 可能是数组格式
		if err2 := json.Unmarshal(rawResp.Effort, &entries); err2 != nil {
			return nil, fmt.Errorf("解析工时日志数据失败: %v", err)
		}
		return entries, nil
	}

	for _, e := range effortMap {
		if e.Deleted != "1" {
			entries = append(entries, e)
		}
	}
	return entries, nil
}
