/*
 * @Package captcha
 * @Author Quan Chen
 * @Date 2020/3/29
 * @Description
 *
 */
package captcha

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const defaultChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

//生成随机字符.
func RandText(num int, chars ...string) string {
	var str string
	if len(chars) != 0 {
		// 自定义字符集
		var stringBuilder strings.Builder
		for _, v := range chars {
			stringBuilder.WriteString(v)
		}
		str = stringBuilder.String()
	} else {
		// 默认字符集
		str = defaultChars
	}

	textNum := len(str)
	text := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < num; i++ {
		text = text + string(str[r.Intn(textNum)])
	}
	return text
}

//生成指定范围的随机数.
func Random(min int64, max int64) float64 {

	if max <= min {
		panic(fmt.Sprintf("invalid range %d >= %d", max, min))
	}
	decimal := rand.Float64()

	if max <= 0 {
		return (float64(rand.Int63n((min*-1)-(max*-1))+(max*-1)) + decimal) * -1
	}
	if min < 0 && max > 0 {
		if rand.Int()%2 == 0 {
			return float64(rand.Int63n(max)) + decimal
		} else {
			return (float64(rand.Int63n(min*-1)) + decimal) * -1
		}
	}
	return float64(rand.Int63n(max-min)+min) + decimal
}
