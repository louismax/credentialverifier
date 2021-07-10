package credentialverifier

import (
	"github.com/louismax/credentialverifier/data"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// substr 截取字符串
func substr(source string, start int, end int) string {
	r := []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

// getAddressInfo 获取地址信息
func getAddressInfo(addressCode string, birthdayCode string, strict bool) map[string]string {
	addressInfo := map[string]string{
		"province": "",
		"city":     "",
		"district": "",
	}

	// 省级信息
	addressInfo["province"] = getAddress(substr(addressCode, 0, 2)+"0000", birthdayCode, strict)

	// 用于判断是否是港澳台居民居住证（8字开头）
	firstCharacter := substr(addressCode, 0, 1)
	// 港澳台居民居住证无市级、县级信息
	if firstCharacter == "8" {
		return addressInfo
	}

	// 市级信息
	addressInfo["city"] = getAddress(substr(addressCode, 0, 4)+"00", birthdayCode, strict)

	// 县级信息
	addressInfo["district"] = getAddress(addressCode, birthdayCode, strict)

	return addressInfo
}

// 获取省市区地址码
func getAddress(addressCode string, birthdayCode string, strict bool) string {
	address := ""
	addressCodeInt, _ := strconv.Atoi(addressCode)
	if _, ok := data.AddressCodeTimeline[addressCodeInt]; !ok {
		// 修复 \d\d\d\d01、\d\d\d\d02、\d\d\d\d11 和 \d\d\d\d20 的历史遗留问题
		// 以上四种地址码，现实身份证真实存在，但民政部历年公布的官方地址码中可能没有查询到
		// 如：440401 450111 等
		// 所以这里需要特殊处理
		// 1980年、1982年版本中，未有制定省辖市市辖区的代码，所有带县的省辖市给予“××××20”的“市区”代码。
		// 1984年版本开始对地级市（前称省辖市）市辖区制定代码，其中“××××01”表示市辖区的汇总码，同时撤销“××××20”的“市区”代码（追溯至1983年）。
		// 1984年版本的市辖区代码分为城区和郊区两类，城区由“××××02”开始排起，郊区由“××××11”开始排起，后来版本已不再采用此方式，已制定的代码继续沿用。
		suffixes := substr("123456", 4, 6)
		switch suffixes {
		case "20":
			address = "市区"
		case "01":
			address = "市辖区"
		case "02":
			address = "城区"
		case "11":
			address = "郊区"
		}

		return address
	}

	timeline := data.AddressCodeTimeline[addressCodeInt]
	year, _ := strconv.Atoi(substr(birthdayCode, 0, 4))
	startYear := "0001"
	endYear := "9999"
	for _, val := range timeline {
		if val["start_year"] != "" {
			startYear = val["start_year"]
		}
		if val["end_year"] != "" {
			endYear = val["end_year"]
		}
		startYearInt, _ := strconv.Atoi(startYear)
		endYearInt, _ := strconv.Atoi(endYear)
		if year >= startYearInt && year <= endYearInt {
			address = val["address"]
		}
	}

	if address == "" && !strict {
		for _, val := range timeline {
			// 由于较晚申请户口或身份证等原因，导致会出现地址码正式启用于2000年，但实际1999年出生的新生儿，由于晚了一年报户口，导致身份证上的出生年份早于地址码正式启用的年份
			// 由于某些地区的地址码已经废弃，但是实际上在之后的几年依然在使用
			// 这里就不做时间判断了
			address = val["address"]
			break
		}
	}

	return address
}

// getConstellation 获取星座信息
func getConstellation(birthdayCode string) string {
	monthStr := substr(birthdayCode, 4, 6)
	dayStr := substr(birthdayCode, 6, 8)
	month, _ := strconv.Atoi(monthStr)
	day, _ := strconv.Atoi(dayStr)
	startDate := data.Constellation[month]["start_date"]
	startDay, _ := strconv.Atoi(strings.Split(startDate, "-")[1])
	if day >= startDay {
		return data.Constellation[month]["name"]
	}

	tmpMonth := month - 1
	if month == 1 {
		tmpMonth = 12
	}

	return data.Constellation[tmpMonth]["name"]
}

// getChineseZodiac 获取生肖信息
func getChineseZodiac(birthdayCode string) string {
	// 子鼠
	start := 1900
	end, _ := strconv.Atoi(substr(birthdayCode, 0, 4))
	key := (end - start) % 12

	return data.ChineseZodiac[key]
}

// generatorAddressCode 生成地址码
func generatorAddressCode(address string) string {
	addressCode := ""
	for code, addressStr := range data.AddressCode {
		if address == addressStr {
			addressCode = strconv.Itoa(code)
			break
		}
	}

	classification := addressCodeClassification(addressCode)
	switch classification {
	case "country":
		// addressCode = getRandAddressCode("\\d{4}(?!00)[0-9]{2}$")
		addressCode = getRandAddressCode("\\d{4}(?)[0-9]{2}$")
	case "province":
		provinceCode := substr(addressCode, 0, 2)
		// pattern := "^" + provinceCode + "\\d{2}(?!00)[0-9]{2}$"
		pattern := "^" + provinceCode + "\\d{2}(?)[0-9]{2}$"
		addressCode = getRandAddressCode(pattern)
	case "city":
		cityCode := substr(addressCode, 0, 4)
		// pattern := "^" + cityCode + "(?!00)[0-9]{2}$"
		pattern := "^" + cityCode + "(?)[0-9]{2}$"
		addressCode = getRandAddressCode(pattern)
	}

	return addressCode
}

// getRandAddressCode 获取随机地址码
func getRandAddressCode(pattern string) string {
	mustCompile := regexp.MustCompile(pattern)
	var keys []string
	for key := range data.AddressCode {
		keyStr := strconv.Itoa(key)
		if mustCompile.MatchString(keyStr) && substr(keyStr, 4, 6) != "00" {
			keys = append(keys, keyStr)
		}
	}

	rand.Seed(time.Now().Unix())
	if keys != nil {
		return keys[rand.Intn(len(keys))]
	}
	return ""
}

// generatorBirthdayCode 生成出生日期码
func generatorBirthdayCode(addressCode string, address string, birthday string) string {
	sYear := "0001"
	endYear := "9999"
	year := datePipelineHandle(datePad(substr(birthday, 0, 4), "year"), "year")
	month := datePipelineHandle(datePad(substr(birthday, 4, 6), "month"), "month")
	day := datePipelineHandle(datePad(substr(birthday, 6, 8), "day"), "day")

	addressCodeInt, _ := strconv.Atoi(addressCode)
	if _, ok := data.AddressCodeTimeline[addressCodeInt]; ok {
		timeLine := data.AddressCodeTimeline[addressCodeInt]
		for _, val := range timeLine {
			if val["address"] == address {
				if val["start_year"] != "" {
					sYear = val["start_year"]
				}
				if val["end_year"] != "" {
					endYear = val["end_year"]
				}
			}
		}
	}

	yearInt, _ := strconv.Atoi(year)
	startYerInt, _ := strconv.Atoi(sYear)
	endYearInt, _ := strconv.Atoi(endYear)
	if yearInt < startYerInt {
		year = sYear
	}
	if yearInt > endYearInt {
		year = endYear
	}

	birthday = year + month + day
	_, err := time.Parse("20060102", birthday)
	// example: 195578
	if err != nil {
		year = datePad(year, "year")
		month = datePad(month, "month")
		day = datePad(day, "day")
	}

	return year + month + day
}

// datePad 日期补全
func datePad(date string, category string) string {
	padLength := 2
	if category == "year" {
		padLength = 4
	}

	for i := 0; i < padLength; i++ {
		length := len([]rune(date))
		if length < padLength {
			// date = fmt.Sprintf("%s%s", "0", date)
			date = "0" + date
		}
	}

	return date
}

// datePipelineHandle 日期处理
func datePipelineHandle(date string, category string) string {
	dateInt, _ := strconv.Atoi(date)

	switch category {
	case "year":
		nowYear := time.Now().Year()
		rand.Seed(time.Now().Unix())
		if dateInt < 1800 || dateInt > nowYear {
			randDate := rand.Intn(nowYear-1950) + 1950
			date = strconv.Itoa(randDate)
		}
	case "month":
		if dateInt < 1 || dateInt > 12 {
			randDate := rand.Intn(12-1) + 1
			date = strconv.Itoa(randDate)
		}

	case "day":
		if dateInt < 1 || dateInt > 31 {
			randDate := rand.Intn(28-1) + 1
			date = strconv.Itoa(randDate)
		}
	}

	return date
}

// addressCodeClassification 地址码分类
func addressCodeClassification(addressCode string) string {
	// 全国
	if addressCode == "" {
		return "country"
	}

	// 港澳台
	if substr(addressCode, 0, 1) == "8" {
		return "special"
	}

	// 省级
	if substr(addressCode, 2, 6) == "0000" {
		return "province"
	}

	// 市级
	if substr(addressCode, 4, 6) == "00" {
		return "city"
	}

	// 县级
	return "district"
}

// generatorOrderCode 生成顺序码
func generatorOrderCode(sex int) string {
	rand.Seed(time.Now().Unix())
	orderCode := rand.Intn(999-111) + 111
	if sex != orderCode%2 {
		orderCode--
	}

	return strconv.Itoa(orderCode)
}

// 生成Bit码
func generatorCheckBit(body string) string {
	// 位置加权
	var posWeight [19]float64
	for i := 2; i < 19; i++ {
		weight := int(math.Pow(2, float64(i-1))) % 11
		posWeight[i] = float64(weight)
	}

	// 累身份证号body部分与位置加权的积
	var bodySum int
	bodyArray := strings.Split(body, "")
	count := len(bodyArray)
	for i := 0; i < count; i++ {
		bodySub, _ := strconv.Atoi(bodyArray[i])
		bodySum += bodySub * int(posWeight[18-i])
	}

	// 生成校验码
	checkBit := (12 - (bodySum % 11)) % 11
	if checkBit == 10 {
		return "x"
	}
	return strconv.Itoa(checkBit)
}
