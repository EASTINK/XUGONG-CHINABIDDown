package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	cookiejar "github.com/juju/persistent-cookiejar"
	"github.com/xuri/excelize/v2"
)

func LogPrintln(v ...any) {
	log.Println(v)
}

// save logfile
func Log_save() *os.File {
	LogPrintln(GetCurrentAbPath() + `/log/`)
	if err := os.MkdirAll(filepath.Dir(GetCurrentAbPath()+`/log/`), os.ModePerm); err == nil {
		//格式改动 win下不支持名字带":" | 注意在Linux下不要使用go run .,你应该通过build或者IDE的debug,不然拿到的目录是错的
		if logFile, err := os.OpenFile(GetCurrentAbPath()+`/log/`+time.Now().Format("2006-01-02 15-04-05")+`.log`, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
			multiWriter := io.MultiWriter(os.Stdout, logFile)
			log.SetOutput(multiWriter)
			return logFile
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}

}

// 最终方案-全兼容
func GetCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// 返回目录下已获取成功的页数表
func okpage_A() []int {
	var okpage []int
	//------------
	file, err := os.ReadFile("okpage.txt")
	if err == nil {
		okpagestrs := strings.Split(string(file), "\n")
		for _, v := range okpagestrs {
			num, err := strconv.Atoi(v)
			if err == nil || (v == "" && len(okpagestrs) == 1) {
				okpage = append(okpage, num)
			}
		}
	} else {
		if os.IsNotExist(err) {
			okpage = make([]int, 1)
			okpage[0] = 0
		} else {
			LogPrintln(err)
			return nil
		}
	}
	return okpage
}

// 返回目录下已处理成功的qyid表
func okqyids_A() []string {
	var okqyids []string
	//------------
	file, err := os.ReadFile("okqyids.txt")
	if err == nil {
		okpagestrs := strings.Split(string(file), "\n")
		for _, v := range okpagestrs {
			if err == nil { //|| (v != "" && len(okpagestrs) >= 1) {
				okqyids = append(okqyids, v)
			}
		}
	} else {
		if os.IsNotExist(err) {
			okqyids = make([]string, 1)
			okqyids[0] = ""
		} else {
			LogPrintln(err)
			return nil
		}
	}
	return okqyids
}

// 返回目录下输入qyid数据表
func idpage() []string {
	var idpage []string
	//-------------
	file, err := os.ReadFile("id.txt")
	if err == nil {
		idpagestrs := strings.Split(string(file), "\n")
		for _, v := range idpagestrs {
			if v != "" {
				idpage = append(idpage, v)
			}
		}
	} else {
		LogPrintln(err)
		return nil
	}
	return idpage
}

func savebytes(data []byte, savepath string, msg string) {
	if err := os.MkdirAll(filepath.Dir(savepath), os.ModePerm); err == nil {
		if err = os.WriteFile(savepath, data, os.ModePerm); err == nil {
			// LogPrintln(msg)
			LogPrintln(msg)
		}
	}
}

// ------------
type Lxr_A_Result struct {
	Infocount int          `json:"Infocount"`
	Msg       string       `json:"Msg"`
	Ret       bool         `json:"Ret"`
	Status    int          `json:"Status"`
	Other     string       `json:"Other"`
	other     []lxr_A_data `json:"Other"`
}

// 鬼知道为什么一层会出错？
func (res *Lxr_A_Result) Unmarshal() error {
	var data []lxr_A_data
	if err := json.Unmarshal([]byte(res.Other), &data); err != nil {
		return err
	}
	res.other = data
	return nil
}

type lxr_A_data struct {
	//内存省点是一点吧 .done
	// LxrCount string `json:"lxrcount"`
	QyID string `json:"qyid"`
	// QyName   string `json:"qyname"`
	// QyType   string `json:"qytype"`
}

type Lxr_B_Result struct {
	Msg    string       `json:"Msg"`
	Ret    bool         `json:"Ret"`
	Status int          `json:"Status"`
	Other  string       `json:"Other"`
	other  []lxr_B_data `json:"Other"`
}

// 鬼知道为什么一层会出错？
func (res *Lxr_B_Result) Unmarshal() error {
	var data []lxr_B_data
	if err := json.Unmarshal([]byte(res.Other), &data); err != nil {
		return err
	}
	res.other = data
	return nil
}

type lxr_B_data struct {
	Bumen   string `json:"bumen"`
	Code    string `json:"code"`
	Company string `json:"company"`
	Dizhi   string `json:"dizhi"`
	Faren   string `json:"faren"`
	Fax     string `json:"fax"`
	Jianjie string `json:"jianjie"`
	LxrList []struct {
		Lianxiren  string `json:"lianxiren"`
		Lxremail   string `json:"lxremail"`
		Lxrlyid    int    `json:"lxrlyid"`
		Lxrshouji  string `json:"lxrshouji"`
		Lxrtel     string `json:"lxrtel"`
		Lxrzhineng string `json:"lxrzhineng"`
		Lxrzhiwu   string `json:"lxrzhiwu"`
	} `json:"lxrlist"`
	Lytp   int    `json:"lytp"`
	Mail   string `json:"mail"`
	Mobile string `json:"mobile"`
	Other  string `json:"other"`
	Tel    string `json:"tel"`
	Yewu   string `json:"yewu"`
}

func isChanClose(ch chan struct {
	chdata []lxr_A_data
	page   int
	msg    string
}) bool {
	select {
	case _, received := <-ch:
		return !received
	default:
	}
	return false
}

// 拿公司联系人的json信息 结构如上
func getData(c *req.Client, qyid string) (error, Lxr_B_Result) {
	urlB := "https://www.bidcenter.com.cn/JsonHandler/LianXiRenDataSearchHandler.aspx" //{qyid}
	t := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	var qy Lxr_B_Result
	err := c.Get(urlB + "?v=" + t + "&yzID=" + qyid).
		//?v=1695476776007&yzID=
		// SetBody("v=" + t + "&yzID=" + qyid).
		Do().Into(&qy)
	if err != nil {
		LogPrintln(err)
		return err, Lxr_B_Result{}
	}
	err = qy.Unmarshal()
	if err != nil {
		LogPrintln(err)
		return err, Lxr_B_Result{}
	}
	// LogPrintln(qy)
	return nil, qy
}

func getId(c *req.Client, p int, ch chan struct {
	chdata []lxr_A_data
	page   int
	msg    string
}) {
	urlA := "https://www.bidcenter.com.cn/JsonHandler/BuserCenter/LxrDataSearchHandler.aspx"
	// urlA := "https://www.baidu.com"
	var qy Lxr_A_Result
	t := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	err := c.Post(urlA). // + "?v=" + t + "&keyword=&checkstyle=0&type=1&pageindex=" + strconv.Itoa(p)).
				SetBody("v=" + t + "&keyword=&checkstyle=0&type=1&pageindex=" + strconv.Itoa(p)).
				Do().
				Into(&qy)
	if err != nil {
		ch <- struct {
			chdata []lxr_A_data
			page   int
			msg    string
		}{chdata: nil, page: p, msg: "请求异常"}
		return
	}
	//不知道会不会返回查询成功消息格式不同的数据 预防一手
	//["Infocount":0,"Msg":"您的设计院数据库搜索过多....."Ret":false,"Status"]
	if !qy.Ret || qy.Msg != "" || qy.Status != 0 || qy.Infocount == 0 {
		//查询限制
		// return nil
		ch <- struct {
			chdata []lxr_A_data
			page   int
			msg    string
		}{chdata: nil, page: p, msg: qy.Msg}
		return
	}
	if err := qy.Unmarshal(); err != nil {
		//解析格式有问题
		ch <- struct {
			chdata []lxr_A_data
			page   int
			msg    string
		}{chdata: nil, page: p, msg: "解析出错"}
		return
	}
	//pass ->
	ch <- struct {
		chdata []lxr_A_data
		page   int
		msg    string
	}{chdata: qy.other, page: p, msg: ""}
	return
}

/*
返回带Cookes的HTTP客户端

@params: 请求客户端，请求最大数，失败最大数
*/
func client() *req.Client {
	jar, err := cookiejar.New(&cookiejar.Options{
		Filename: "cookies.json",
	})
	if err != nil {
		log.Fatalf("failed to create persistent cookiejar: %s\n", err.Error())
	}
	//默认跳转代理
	// var proxies = []string{
	// 	"socks5://localhost:7890",
	// }
	client := req.C()
	client.
		SetCookieJar(jar).
		SetTimeout(10 * time.Second)
		// SetProxyURL(proxies[0]) //
		// SetCommonRetryCount(len(proxies)).
		// SetCommonRetryCondition(func(resp *req.Response, err error) bool {
		// 	return err != nil || resp.StatusCode == http.StatusTooManyRequests
		// }).
		// SetCommonRetryHook(func(resp *req.Response, err error) {
		// 	c := client.Clone().SetProxyURL(proxies[resp.Request.RetryAttempt-1])
		// 	resp.Request.SetClient(c)
		// })
	return client
}

/*
自动跳过已完成请求页数

@params: 请求客户端，请求最大数，失败最大数
*/
func action_ID(c *req.Client, post_max int, failed_max int) {
	ok_page := okpage_A() //拿目录下已完成的页码 若没有okpage.txt 会自动创建
	if ok_page == nil {
		LogPrintln("获取已完成页码数据失败，正在退出")
		return
	}
	var dataAll []lxr_A_data
	post_index := 1
	nok := 0
	ok := 0
	ch := make(chan struct {
		chdata []lxr_A_data
		page   int
		msg    string
	})
	for i := post_index; i <= post_max; i++ {
		isok := false
		for _, n := range ok_page {
			//如果i页已经成功过了
			if i == n {
				isok = true
			}
		}
		if !isok {
			go getId(c, i, ch)
		} else {
			continue
		}

		if ok >= post_max {
			break
		}

		data := <-ch
		if data.chdata != nil {
			dataAll = append(dataAll, data.chdata...)
			ok_page = append(ok_page, data.page)
			ok += 1
		} else {
			if data.msg != "" {
				LogPrintln("查询页数"+strconv.Itoa(data.page)+"：", data.msg)
			}
			nok += 1
			if nok >= failed_max {
				LogPrintln("查询异常，停止查询")
				break
			}
		}
	}

	if file, err := os.OpenFile("okpage.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		str := ""
		for _, k := range ok_page {
			str += strconv.Itoa(k) + "\n"
		}
		if xx, err := os.ReadFile("okpage.txt"); err == nil && string(xx) != (str) {
			file.WriteString(str)
		}
	} else {
		LogPrintln(err)
	}

	if file, err := os.OpenFile("id.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		str := ""
		for _, k := range dataAll {
			str += (k.QyID) + "\n"
		}
		file.WriteString(str)
	} else {
		LogPrintln(err)
	}

	LogPrintln("获取ID:"+"成功"+strconv.Itoa(ok), "失败:"+strconv.Itoa(nok), "目录下已保存数据:"+strconv.Itoa(len(okqyids_A())-1))

	// for _, chgs := range chs {
	// 	for _, ch := range chgs {
	// 		if !isChanClose(ch) {
	// 			close(ch)
	// 		}
	// 	}
	// }
	if !isChanClose(ch) {
		close(ch)
	}
}

/*
启用goruntine并发拿取json ;

@params: 客户端	要处理的qyid	已处理的qyid	最大失败重试数
*/
func action_Data(c *req.Client, renextmax int) {
	var okqyid []string
	//remax := 3 //遇见异常后往下重新尝试3次
	faildX := 0
	faild := 0
	okqyids := okqyids_A() //拿目录下已处理的qyid 如上
	if okqyids == nil {
		LogPrintln("获取已处理QYID数据失败，正在退出")
		return
	}
	idpage := idpage() //拿目录下的所有qyid
	if idpage == nil {
		LogPrintln("获取QYID数据输入失败，正在退出")
		return
	}
	for i := 0; i < len(idpage); i++ {
		qyid := idpage[i]
		isok := false
		//qyid去重
		//....if qyid != xx().....{...}
		for _, v := range okqyids {
			if qyid == v {
				isok = true
			}
		}
		if isok {
			break
		}
		if err, data := getData(c, qyid); err == nil && qyid != "" {
			//如果一直遇见异常
			if faild >= renextmax {
				break
			}
			//如果异常不连续 恢复计数
			faild = 0
			//拿到数据 做点事 ....done
			//.....记录okqyid
			okqyid = append(okqyid, qyid)
			//.........做一个导出的数据处理格式
			type Result struct {
				Msg    string `json:"Msg"`
				Ret    bool   `json:"Ret"`
				Status int    `json:"Status"`
				// Other  string       `json:"Other"`
				Other []lxr_B_data `json:"Other"`
			}
			dataB := &Result{
				Msg:    data.Msg,
				Ret:    data.Ret,
				Status: data.Status,
				Other:  data.other,
			}
			//lxrlist := data.other[0].LxrList
			//不知道干点什么好 先扔到目录下吃灰吧
			jsondata, err := json.Marshal(dataB)
			if err == nil {
				savebytes(jsondata,
					"./data/"+qyid+"_"+data.other[0].Company+".json",
					qyid+"		保存到 : ./data/"+data.other[0].Company+".json")
			}
		} else {
			//返回异常：比如后端限制了请求数
			faildX += 1
			faild += 1
		}
	}
	//-------------
	if len(okqyid) > 0 {
		//保存处理完的qyid
		if file, err := os.OpenFile("okqyids.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
			str := ""
			for _, k := range okqyid {
				str += k + "\n"
			}
			//保证本地文件数据和预写入数据不一致
			if xx, err := os.ReadFile("okqyids.txt"); err == nil && string(xx) != (str) {
				file.WriteString(str)
			}
		} else {
			//抛出异常
			LogPrintln("保存okqyid失败")
			LogPrintln(err)
		}
	}

	LogPrintln("处理数据 成功:"+strconv.Itoa(len(okqyid)), "请求异常:"+strconv.Itoa(faildX), "目录下已保存数据："+strconv.Itoa(len(okqyids_A())-1))
}

// ----------------
// 还没写完
// ...保存到excel
func savexlsx() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Set value of a cell.
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
