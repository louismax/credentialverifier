package credentialverifier

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// generateCode 生成数据
func generateCode(id string) (map[string]string, error) {
	length := len(id)
	if length == 15 {
		return generateShortCode(id)
	}
	if length == 18 {
		return generateLongCode(id)
	}
	return map[string]string{}, errors.New("invalid ID card number length")
}

// generateShortCode 生成短数据
func generateShortCode(id string) (map[string]string, error) {
	if len(id) != 15 {
		return map[string]string{}, errors.New("invalid ID card number length")
	}

	mustCompile := regexp.MustCompile("(.{6})(.{6})(.{3})")
	subMatch := mustCompile.FindStringSubmatch(strings.ToLower(id))

	return map[string]string{
		"body":         subMatch[0],
		"addressCode":  subMatch[1],
		"birthdayCode": "19" + subMatch[2],
		"order":        subMatch[3],
		"checkBit":     "",
		"type":         "15",
	}, nil
}

// generateLongCode 生成长数据
func generateLongCode(id string) (map[string]string, error) {
	if len(id) != 18 {
		return map[string]string{}, errors.New("invalid ID card number length")
	}
	mustCompile := regexp.MustCompile("((.{6})(.{8})(.{3}))(.)")
	subMatch := mustCompile.FindStringSubmatch(strings.ToLower(id))

	return map[string]string{
		"body":         subMatch[1],
		"addressCode":  subMatch[2],
		"birthdayCode": subMatch[3],
		"order":        subMatch[4],
		"checkBit":     subMatch[5],
		"type":         "18",
	}, nil
}

//checkOrderCode 检查顺序码
func checkOrderCode(orderCode string) bool {
	return len(orderCode) == 3
}

// checkBirthdayCode 检查出生日期码
func checkBirthdayCode(birthdayCode string) bool {
	year, _ := strconv.Atoi(substr(birthdayCode, 0, 4))
	if year < 1800 {
		return false
	}

	nowYear := time.Now().Year()
	if year > nowYear {
		return false
	}

	_, err := time.Parse("20060102", birthdayCode)

	return err == nil
}
// checkAddressCode 检查地址码
func checkAddressCode(addressCode string, birthdayCode string, strict bool) bool {
	return getAddressInfo(addressCode, birthdayCode, strict)["province"] != ""
}