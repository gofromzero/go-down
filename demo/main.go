package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//var (
//	url_android string = "https://qd.myapp.com/myapp/qqteam/AndroidQQ/mobileqq_android.apk"
//	//url_pc = "https://dldir1.qq.com/qqfile/qq/PCQQ9.1.3/25326/QQ9.1.3.25326.exe"
//
//)

var progressList []float32
var fileList []string

func getLastIndex(s string, ch string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ch[0] {
			return i
		}
	}
	return 0
}

func getIndex(s string, ch string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ch[0] {
			return i
		}
	}
	return 0
}

func calcLength(L int) string {
	if L < 1024 {
		return fmt.Sprintf("%d Byte", L)
	}
	kb := float32(L) / 1024
	if kb < 1024 {
		return fmt.Sprintf("%f Kb", kb)
	}
	mb := kb / 1024
	if mb < 1024 {
		return fmt.Sprintf("%f Mb", mb)
	}
	gb := mb / 1024
	if mb < 1024 {
		return fmt.Sprintf("%f GB", gb)
	}
	return fmt.Sprintf("%f PB", gb/1024)
}

func initEnv(i ...int) {
	if len(i) == 0 {
		num := os.Getenv("number_of_processors")
		i[0], _ = strconv.Atoi(num)
	}
	runtime.GOMAXPROCS(i[0])
}

func strip(s string, chars string) string {
	length := len(s)
	max := len(s) - 1
	l, r := true, true //标记当左端或者右端找到正常字符后就停止继续寻找
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[uint8]bool) //创建字符集，也就是唯一的字符，方便后面判断是否存在
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r { // 如果左端和右端都没找到正常字符，那么表示该字符串没有正常字符
		return ""
	}
	return s[start : end+1]
}

func showProgress() {
	for {
		for i := 0; i < len(progressList); i++ {
			var size float32
			_ = filepath.Walk(fileList[i], func(path string, info os.FileInfo, err error) error {
				size = float32(info.Size())
				return nil
			})
			progress := size * 100 / progressList[i]
			//fmt.Printf("当前为第%d段,进度为",i)
			//fmt.Printf(	"\t%c[1;40;32m%.3f %% \r",0x1B,progress)
			//_, _ = fmt.Fprintf(os.Stdout, "当前为第%d段,进度  %.2f %% \r", i, progress)
			_, _ = fmt.Fprintf(os.Stdout, "当前进度: %.2f %% \r", progress)
		}
		time.Sleep(time.Millisecond * 500)
	}

}

func download(url string, filename string, dir string, msg chan int) {
	res, err := http.Get(string(url))
	if err != nil {
		panic(err)
	}
	//获取文件名
	if filename == "" {
		value, val := res.Header["Content-Disposition"] //从response的Header中获取文件名
		if val {
			tmpSplit := strings.Split(value[0], "filename=")
			if len(tmpSplit) > 1 {
				filename = tmpSplit[1]
			} else {
				filename = tmpSplit[0]
			}
		} else {
			lastIndex := getLastIndex(url, "/")
			filename = url[lastIndex+1:]
		}
	}
	filename = strip(filename, "< > / \\ | : \" * ?")
	if len(filename) < 1 {
		filename = "unknown"
	} else {
		reg := regexp.MustCompile(`[<\\>/|:"*?]`)
		filename = reg.ReplaceAllString(filename, "_")
	}
	//从response的Header中获取文件大小
	contentLenStr, exist := res.Header["Content-Length"]
	if !exist || len(contentLenStr) == 0 {
		contentLenStr = []string{"0"}
	}
	//转换字符串的文件大小为int的大小，方便后面计算
	contentLen, convertErr := strconv.Atoi(contentLenStr[0])
	if convertErr != nil {
		fmt.Println(convertErr)
		contentLen = 0
		fmt.Println("文件大小未知")
	}
	fmt.Println("文件大小:", calcLength(contentLen)) //计算并显示文件大小
	fmt.Println("文 件 名:", filename)
	//开始创建保存目录
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		e := os.MkdirAll(dir, os.ModePerm)
		if e != nil {
			fmt.Printf("不能创建目录")
			panic(e)
		}
	}
	filePath := path.Join(dir, filename) //拼接文件路径
	fileList[0] = filePath
	fmt.Println("保存位置:", dir)
	f, err := os.Create(filePath) //创建文件
	if err != nil {
		fmt.Println("文件创建失败")
		panic(err)
	}
	progressList[0] = float32(contentLen)
	//fmt.Println(f)
	//fmt.Println("------------------------------")
	//for k,v:=range res.Header {
	//	fmt.Println(k,v)
	//}
	//fmt.Println("------------------------------")
	go showProgress()
	io.Copy(f, res.Body)
	msg <- 0
}

/*
分配并创建任务
*/
func assignTask(url string, threadNum int) {
	taskList := make([][2]int, threadNum)
	/*
		获取并分配任务
	*/
	mergeChan := make(chan Part, threadNum)
	for i, taskArgs := range taskList {
		fmt.Println(i, taskArgs)
		go task(i, taskArgs, mergeChan)
	}
	go merge(threadNum, mergeChan)
}

func task(orderNum int, taskArgs [2]int, mergeChan chan Part) {
	//开始任务
	filePath := ""
	var part Part = Part{orderNum: orderNum, path: filePath}
	/*
		根据任务参数下载任务part
	*/
	mergeChan <- part
}

//文件的结构体里面存的是一个指针{ *file}，所以不需要指针
type Part struct {
	orderNum int
	path     string
}

func merge(threadNum int, mergeChan chan Part) {
	n := 0
	length := len(mergeChan)
	filePathList := make([]string, length)
	for part := range mergeChan {
		filePathList[part.orderNum] = part.path
		n++
		if n == threadNum {
			break
		}
	}
	close(mergeChan)
	/*
		开始合并文件
	*/
	//合并完毕
}

func main() {
	//保存路径（文件夹） ,url
	threadNum := 1
	progressList = make([]float32, threadNum)
	fileList = make([]string, threadNum)
	initEnv(1) //设置运行的最大核心数
	msg := make(chan int, 0)
	args := os.Args
	filename := ""
	if len(args) < 3 {
		fmt.Println("参数格式 [url] [dir] [filename] ")
		return
	}
	if len(args) > 3 {
		filename = args[3]
	}
	fmt.Println()
	url, dir := args[1], args[2]
	//go download(url,callback,"d:/downloads",msg)
	go download(url, filename, dir, msg)
	select {
	case _ = <-msg:
		_, _ = fmt.Fprintf(os.Stdout, "当前进度: %.2f %% \r", 100.00)
		fmt.Println("\n下载完成")
	}
}
