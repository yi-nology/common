package phone

import (
	"github.com/emirpasic/gods/sets/hashset"
	"regexp"
)

// AirtelCheck 获取手机号运营商
func AirtelCheck(phone string) string {
	//手机运营商判断
	//中国电信：133，153，177，180，181，189；
	//中国联通：130，131，132，155，156，185，186，145，176；
	//中国电信：134，135，136，137，138，139，147，150，151，152，157，158，159，178，182，183，184，187，188。
	dx := hashset.New()
	dx.Add("133", "153", "177", "180", "181", "189")
	lt := hashset.New()
	lt.Add("130", "131", "132", "155", "156", "185", "186", "145", "176")
	yd := hashset.New()
	yd.Add("134", "135", "137", "138", "139", "147", "150", "151", "152", "157", "158", "159", "178", "182", "183", "184", "187", "188")
	if len(phone) != 13 && len(phone) != 11 {
		return ""
	}

	if len(phone) == 13 {
		phone = phone[2:]
	}
	pre := phone[0:3]
	if dx.Contains(pre) {
		return "dx"
	} else if lt.Contains(pre) {
		return "lt"
	} else {
		return "yd"
	}

}

// ValidatePhoneNumber 校验中国大陆手机号
func ValidatePhoneNumber(phoneNumber string) bool {

	// 使用正则表达式匹配手机号码
	// ^ 表示字符串的开始
	// [1] 表示第一位必须是1
	// [3-9] 表示第二位可以是3、4、5、6、7、8或9
	// \d{9} 表示后面必须跟着9个数字
	// $ 表示字符串的结束
	reg := `^1[3-9]\d{9}$`
	match, _ := regexp.MatchString(reg, phoneNumber)

	return match
}
