package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"math"
	"regexp"
	"time"
	//"errors"
	"github.com/philsong/ohlala/golink"
	"io"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
)

// hash a string
func PasswordHash(pwd string) string {
	hasher := sha1.New()
	hasher.Write([]byte(pwd))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

var defaultPagesize = golink.PAGE_SIZE

// 检查分页参数。
// page第一页为1；
// pagesize默认值为20，范围为 5~200.
// return page, pagesize
// 返回的page是从0开始的值
func PageCheck(page, pagesize int) (int, int) {
	if page < 1 {
		page = 1
	}
	page = page - 1
	if pagesize == 0 {
		pagesize = defaultPagesize
	} else if pagesize < 5 {
		pagesize = 5
	} else if pagesize > 200 {
		pagesize = 200
	}
	return page, pagesize
}

// 从请求参数中获取分页相关的参数
func PagerParams(r *http.Request) (page, pagesize int) {
	pagesize, _ = strconv.Atoi(r.FormValue("pagesize"))
	page, _ = strconv.Atoi(r.FormValue("page"))
	if page == 0 {
		page = 1
	}
	if pagesize == 0 {
		pagesize = defaultPagesize
	}
	return
}

/** 微博时间格式化显示
 * @param timestamp，标准时间戳
 */
func SmcTimeSince(timeAt time.Time) string {
	now := time.Now()
	since := math.Abs(float64(now.Unix() - timeAt.Unix()))

	output := ""
	switch {
	case since < 60:
		output = "刚刚"
	case since < 60*60:
		output = fmt.Sprintf("%v分钟前", math.Floor(since/60))
	case since < 60*60*24:
		output = fmt.Sprintf("%v小时前", math.Floor(since/3600))
	case since < 60*60*24*2:
		output = fmt.Sprintf("昨天%v", timeAt.Format("15:04"))
	case since < 60*60*24*3:
		output = fmt.Sprintf("前天%v", timeAt.Format("15:04"))
	case timeAt.Format("2006") == now.Format("2006"):
		output = timeAt.Format("1月2日 15:04")
	default:
		output = timeAt.Format("2006年1月2日 15:04")
	}
	// if math.Floor(since/3600) > 0 {
	//     if timeAt.Format("2006-01-02") == now.Format("2006-01-02") {
	//         output = "今天 "
	//         output += timeAt.Format("15:04")
	//     } else {
	//         if timeAt.Format("2006") == now.Format("2006") {
	//             output = timeAt.Format("1月2日 15:04")
	//         } else {
	//             output = timeAt.Format("2006年1月2日 15:04")
	//         }
	//     }
	// } else {
	//     m := math.Floor(since / 60)
	//     if m > 0 {
	//         output = fmt.Sprintf("%v分钟前", m)
	//     } else {
	//         output = "刚刚"
	//     }
	// }
	return output
}

//获取这个小时的开始点
func ThisHour() time.Time {
	t := time.Now()
	year, month, day := t.Date()
	hour, _, _ := t.Clock()

	return time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
}

//获取今天的开始点
func ThisDate() time.Time {
	t := time.Now()
	year, month, day := t.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

//获取这周的开始点
func ThisWeek() time.Time {
	t := time.Now()
	year, month, day := t.AddDate(0, 0, -1*int(t.Weekday())).Date()

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

//获取这月的开始点
func ThisMonth() time.Time {
	t := time.Now()
	year, month, _ := t.Date()

	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

//获取今年的开始点
func ThisYear() time.Time {
	t := time.Now()
	year, _, _ := t.Date()

	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}

func GetEmailRegexp() (*regexp.Regexp, error) {
	return regexp.Compile(`^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((\.[a-zA-Z0-9_-]{2,3}){1,2})$`)
}

// 对字符串进行md5哈希,
// 返回32位小写md5结果
func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 对字符串进行md5哈希,
// 返回16位小写md5结果
func MD5_16(s string) string {
	return MD5(s)[8:24]
}

/**
* user : example@example.com login smtp server user
* password: xxxxx login smtp server password
* host: smtp.example.com:port   smtp.163.com:25
* to: example@example.com;example1@163.com;example2@sina.com.cn;...
* subject:The subject of mail
* body: The content of mail
* mailtyoe: mail type html or text
 */
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func GetSensitiveInfoRemovedEmail(email string) string {
	const (
		mail_separator_sign = "@"
		min_mail_id_length  = 2
	)

	emailSepPos := strings.Index(email, mail_separator_sign)

	if emailSepPos < 0 {
		return email
	}

	mailId, mailDomain := email[:emailSepPos], email[emailSepPos+1:]

	if mailIdLength := len(mailId); mailIdLength > min_mail_id_length {
		firstChar, lastChar := string(mailId[0]), string(mailId[mailIdLength-1])
		stars := "***"
		switch mailIdLength - min_mail_id_length {
		case 1:
			stars = "*"
		case 2:
			stars = "**"
		}
		mailId = firstChar + stars + lastChar
	}

	result := mailId + mail_separator_sign + mailDomain
	return result
}

//获取url的host
func GetUrlHost(cUrl string) string {

	u, err := url.Parse(cUrl)
	if err != nil {
		return ""
	}
	if strings.Index(u.Host, "www.") == 0 {
		return u.Host[4:]
	}

	return u.Host
}
