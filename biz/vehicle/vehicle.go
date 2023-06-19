package vehicle

import (
	"regexp"

	"github.com/yi-nology/common/utils/xlogger"
)

// Regular expression patterns
const (
	newEnergyPattern    = "^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼][A-Z][DF]\\d{4}[DF]$"
	regularPlatePattern = "^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼][A-Z][A-Z0-9]{5}$"
)

// validatePlate 校验车牌号
func validatePlate(log xlogger.Logger, plate string, pattern string) bool {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		log.Errorf("车牌校验 正则错误 %+v", err)
		return false
	}
	return reg.MatchString(plate)
}

// ValidateLicenseNewPlate 校验中国大陆新能源汽车车牌号
func ValidateLicenseNewPlate(log xlogger.Logger, licensePlate string) bool {
	if len([]rune(licensePlate)) != 8 {
		return false
	}
	return validatePlate(log, licensePlate, newEnergyPattern)
}

// ValidateLicensePlate 验证中国大陆普通车牌号
func ValidateLicensePlate(log xlogger.Logger, licensePlate string) bool {
	if len([]rune(licensePlate)) != 7 {
		return false
	}
	return validatePlate(log, licensePlate, regularPlatePattern)
}

// ValidatePlate 验证车牌号，包括新能源和普通车牌号
func ValidatePlate(log xlogger.Logger, plate string) bool {
	switch len([]rune(plate)) {
	case 8:
		return ValidateLicenseNewPlate(log, plate)
	case 7:
		return ValidateLicensePlate(log, plate)
	default:
		return false
	}
}
