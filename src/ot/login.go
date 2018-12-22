package ot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

const (
	loginUrl  = ""
	username  = ""
	password  = ""
	k8sUrl    = ""
	queryMenu = ""
	logoutUrl = ""
)

var client *http.Client

func init() {
	jar := NewJar()
	client = &http.Client{Jar: jar}
}

func login() {

	values := make(url.Values)
	values.Set("username", username)
	values.Set("password", password)

	resp, err1 := client.PostForm(loginUrl, values)

	if err1 != nil {
		fmt.Println(err1)
	}

	defer resp.Body.Close()
	bs, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		panic(err2)
	}
	//str := string(bs)
	//fmt.Println(str)

	login := &LoginResp{}
	json.Unmarshal(bs, login)

	if login.Success {
		fmt.Println("\n---------->登录成功")
	} else {
		fmt.Println("\n---------->登录失败！")
		return
	}

}

func logout()  {
	client.PostForm(logoutUrl,make(url.Values))
}

func QryApp() (app *AppInfo) {

	login()

	app = &AppInfo{}

	values := make(url.Values)
	values.Set("env", "msit")
	values.Set("system", "order")
	values.Set("type", "bss-order-shopcart-dubbo")
	values.Set("page", "1")
	values.Set("rows", "20")

	resp, err1 := client.PostForm(k8sUrl, values)

	if err1 != nil {
		fmt.Println(err1)
	}

	defer resp.Body.Close()
	bs, e := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(bs))
	if e != nil {
		panic(e)
	}
	var k8sBean K8sResp
	json.Unmarshal(bs, &k8sBean)

	defer logout()
	if k8sBean.TotalCount == 1 {
		app.name = k8sBean.Item[0].Name
		app.hostIp = k8sBean.Item[0].HostIp
		return app
	} else {
		fmt.Println("\n程序正在部署中，请稍后再试")
		return app
	}

}

type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

type LoginResp struct {
	Success     bool   `json:"success,omitempty"`
	Data        string `json:"data,omitempty"`
	Message     string `json:"message,omitempty"`
	MessageType string `json:"message_type,omitempty"`
}

type K8sResp struct {
	TotalCount int8       `json:"totalCount"`
	Item       []K8sServe `json:"items"`
}

type K8sServe struct {
	Name         string `json:"name"`
	HostIp       string `json:"hostIp"`
	PodIp        string `json:"podIP"`
	StartTime    string `json:"startTime"`
	Image        string `json:"image"`
	Phase        string `json:"phase"`
	ContainerNum int8 `json:"containerNum"`
	ReadyNum     int8 `json:"readyNum"`
	Restarts     int8 `json:"restarts"`
}

//应用信息
type AppInfo struct {
	name   string
	hostIp string
}
