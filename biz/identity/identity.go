package identity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidateIDCard 校验身份证号码
func ValidateIDCard(idCard string) bool {
	// 检查长度和基本格式
	if len(idCard) != 18 || !regexp.MustCompile(`^\d{17}[\dX]$`).MatchString(idCard) {
		return false
	}

	// 检查出生日期
	birthDate, err := time.Parse("20060102", idCard[6:14])
	if err != nil || birthDate.After(time.Now()) || birthDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local)) {
		return false
	}

	// 计算校验位
	sum := 0
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCode := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
	for i := 0; i < 17; i++ {
		num, _ := strconv.Atoi(string(idCard[i]))
		sum += num * weights[i]
	}
	if checkCode[sum%11] != strings.ToUpper(string(idCard[17])) {
		return false
	}

	return true
}

// ExtractInfoFromIDCard 从身份证号码中提取信息 行政区号 生日 性别1男2女 错误
func ExtractInfoFromIDCard(idCard string) (string, time.Time, int64, error) {
	// 检查长度
	if len(idCard) != 18 || !regexp.MustCompile(`^\d{17}[\dX]$`).MatchString(idCard) {
		return "", time.Time{}, 0, fmt.Errorf("invalid id card length")
	}

	// 提取行政区划代码
	// 注意：这里需要一个有效的行政区划代码映射表来将代码转换为实际的行政区划名称
	administrativeDivisionCode := idCard[:6]

	// 提取并解析出生日期
	birthDateString := idCard[6:14]
	birthDate, err := time.Parse("20060102", birthDateString)
	if err != nil {
		return "", time.Time{}, 0, fmt.Errorf("invalid birth date")
	}

	// 提取性别
	// 注意：身份证号码的第17位代表性别，奇数代表男性，偶数代表女性
	genderCode, err := strconv.Atoi(string(idCard[16]))
	if err != nil {
		return "", time.Time{}, 0, fmt.Errorf("invalid gender code")
	}
	gender := int64(2)
	if genderCode%2 == 1 {
		gender = int64(1)
	}

	return administrativeDivisionCode, birthDate, gender, nil
}
