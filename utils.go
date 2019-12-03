package authcenter

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var GolangBirth = "2006-01-02 15:04:05"

type EmptyStruct struct{}

type ResponseContentTpl struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func MakeRespContent(code int, desc string, data interface{}) (rct ResponseContentTpl) {
	rct.Code = code
	rct.Msg = desc
	rct.Data = data
	if data == nil {
		rct.Data = EmptyStruct{}
	}
	return
}

func Response403() ResponseContentTpl {
	return MakeRespContent(403, "403 Forbidden", nil)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetTime() string {
	const shortForm = "2006-01-02 15:04:05"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	return str
}

func GetMD5Hash(text string) string {
	haser := md5.New()
	haser.Write([]byte(text))
	return hex.EncodeToString(haser.Sum(nil))
}

func NowUnixTime() int64 {
	return time.Now().Unix()
}

func GetNowtimeMD5() string {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	return GetMD5Hash(timestamp)
}


// 调用这个方法，先判断 err 是否为空，
// err 不为空，证明是其他错误，需要处理错误（比如普通用户访问 root 下的目录会显示权限不足）。
// err 为空，在判断 uint 值: 0表示不存在, 1表示存在是文件, 2表示存在是文件夹。
func CheckFileIfExists(fileName string) (error, uint) {
    info, err := os.Stat(file)
    if err == nil {
        if info.IsDir(){
            reutrn nil, 2
        }
        return nil, 1
    }
    if os.IsNotExist(err) {
        return nil, 0
    }
    return err, 0
}

func GetFileMd5(filePath string) (MD5Str string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Open", err)
		return
	}
	defer file.Close()
	md5hash := md5.New()
	if _, err = io.Copy(md5hash, file); err != nil {
		log.Println("Copy", err)
		return
	}
	MD5Str = hex.EncodeToString(md5hash.Sum(nil))
	return
}

// 创建文件夹
func CreateFolderIfNotExists(dir string) error {
	info, err := os.Stat(dir)
	if err == nil {
		if info.IsDir() {
			return nil
		}
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// 日期时间转时间戳
// DatetimeString2Unix("2017-07-07 09:00:00") => 1499389200
func DatetimeString2Unix(dateTime string) int64 {
	local, _ := time.LoadLocation("Local")
	unixTimestamp, _ := time.ParseInLocation(GolangBirth, dateTime, local)
	return unixTimestamp.Unix()
}

// 获取指定日期的0点时间戳
func GetTheDayZeroTime(date string) int64 {
	datetime := fmt.Sprintf("%s 00:00:00", date)
	unixTs, _ := time.Parse(GolangBirth, datetime)
	return unixTs.Unix()
}

// 返回当前月份的第一天0点时间戳
func GetFirstDateOfCurrentMonth() int64 {
	d := time.Now()
	d = d.AddDate(0, 0, -d.Day()+1)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()).Unix()
}

// 时间字符串节点-->时间戳
type DatetimeNodes struct {
	year  int
	month time.Month
	day   int
	hour  int
	min   int
	sec   int
}

func TimeNodeString2Time(dtn DatetimeNodes) int64 {
	theTime := time.Date(dtn.year, dtn.month, dtn.day, dtn.hour, dtn.min, dtn.sec, 0, time.Local)
	return theTime.Unix()
}

// 时间戳--> 字符串
func Unixtime2String(unixTime int64) string {
	return time.Unix(unixTime, 0).Format(GolangBirth)
}

// 转两位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func UnknownJson(data string) interface{} {
	if data == `` {
		return nil
	}
	r := strings.NewReader(data)
	dec := json.NewDecoder(r)
	switch data[0] {
	case 91:
		// "[" 开头的Json
		param := []interface{}{}
		dec.Decode(&param)
		return param
	case 123:
		// "{" 开头的Json
		param := make(map[string]interface{})
		dec.Decode(&param)
		return param
	default:
		return nil
	}
}

// 异常捕获
func CatchError(skip int, err error) {
	if err == nil {
		return
	}
	pc, file, line, ok := runtime.Caller(skip)
	errorMessage := err.Error()
	if ok {
		//获取函数名
		pcName := runtime.FuncForPC(pc).Name()
		pcName = strings.Join(strings.Split(pcName, "/")[1:], "/")
		file = strings.Join(strings.Split(file, "/")[4:], "/")
		errorMessage = fmt.Sprintf("%s   %d   %s   %s", file, line, pcName, errorMessage)
	}
	log.Println(errorMessage)
}
