// this tool performs create,get,update and lock issues in github.
package main

/**
You should use the commandline to start the tool.

After you start it,you will see some reminders you need to fill in.

	please enter repo and owner: pro-client kingmore96(this is my exaple)
	please enter the mode:(c u g l) g (this is my example)
if you enter an read-mode g(get), there's no need to fill in the token for Authorization
otherwise you need to provide your token to the tool.
So you will see another reminder like below :
	please enter your token: XXX (your token)

Now the tool will get the issues for you and show in the commandline.
When it show some results to you. You will see another remainder:
	do you want to do anythingelse in the last repo and owner?(y n)
if you print y, the tool will reuse the infos you last provided
and ask the mode you need this round.
if you print n, the tool will begin a new round so you need to provide the new repo and owner info.

## other important things
1. the mode you can use in the tool.
	c: create
	u: update
	g: get
	l: lock the issue
	ul : unlock the issue
	q : quit the program
2. write-mode and read-mode
	g is read-mode and c,u,l,ul are the write-mod.
3. if you choose c or u mode, there are somethingelse you need to provide.
	when you type c, you will see another remider
	please enter your issue title:
	if you entered "how to use it?", the tools will do the folowing things:
	1. create a md file named "how to use it?" in the filesystem(use the ongoing path).
	2. open the file with your default editor.
	3. After you finish editting the issue's body and save the file. you need to go to the
	command line and type yes to let the tool know you have finish your edit. (after a reminder finished body?)

	when you type u,you need to provide the same thing as you type c
	Besides,you need to provide the issue number to the tool,so it will know the exact issue you need to update

	when you type l,you don't need to provide the body,you only need to provide the issume number
	so as to ul.
4. want to quit?
	you can type q to quit the tool.

Have fun with the tool!
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/repos/"

// 映射go中的成员名和json返回名时，如果是一个单词，可以不用写Tag，自动进行匹配
// 如果有下划线链接的，需要额外写一下Tag
type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
	Locked    bool
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type SearchInfo struct {
	repo   string
	owner  string
	token  string
	mode   string
	issueN string
}

var client http.Client

func main() {
	si := &SearchInfo{}
	firstIn := true
	var needRefresh bool

	for {
		//读取输入
		if firstIn || needRefresh {
			firstIn = false
			fmt.Print("repo and owner: ")
			_, err := fmt.Scan(&si.repo, &si.owner)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		fmt.Print("mode: ")
		_, err := fmt.Scan(&si.mode)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//处理输入
		switch si.mode {
		case "c", "u", "l", "ul":
			if (*si).token == "" {
				fmt.Print("token: ")
				_, err := fmt.Scan(&si.token)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			if si.mode != "c" {
				fmt.Print("issue numer: ")
				_, err = fmt.Scan(&si.issueN)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			writeIssue(si)
		case "g":
			getIssues(si)
		case "q":
			fmt.Println("bye bye")
			return
		default:
			fmt.Println("wrong mode")
		}

		//判断是否需要刷新配置信息
	s:
		for {
			fmt.Print("reuse repo and owner?(y n) ")

			var r string
			_, err = fmt.Scan(&r)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			switch r {
			case "n", "y":
				if r == "n" {
					needRefresh = true
				}
				break s
			default:
				fmt.Println("wrong param", r)
				continue
			}
		}
	}
}

func getIssues(si *SearchInfo) {
	//fmt.Printf("start deal with %s mode\n", si.mode)
	//发送请求
	//1. 拼接URL
	url := generateURL(*si)
	// fmt.Println(s)
	//2. 发送请求
	rp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rp.Body.Close()

	//处理结果,不为200
	if rp.StatusCode != http.StatusOK {
		fmt.Println(rp.Status)
		return
	}

	//解析body
	var bodyResult []*Issue
	if err := json.NewDecoder(rp.Body).Decode(&bodyResult); err != nil {
		fmt.Println(err)
		return
	}

	//打印结果
	fmt.Println("Get Result:")
	for _, v := range bodyResult {
		fmt.Printf("number=%d user=%s state=%s locked=%t title=%q\n", v.Number, v.User.Login, v.State, v.Locked, v.Title)
	}
}

func writeIssue(si *SearchInfo) {
	if si.mode == "l" {
		lockOrUnlockAnIssue(si, true)
	} else if si.mode == "ul" {
		lockOrUnlockAnIssue(si, false)
	} else if si.mode == "c" {
		createOrUpdateAnIssue(si, true)
	} else {
		createOrUpdateAnIssue(si, false)
	}
}

func lockOrUnlockAnIssue(si *SearchInfo, lock bool) {
	url := generateURL(*si)
	var rq *http.Request
	var err error
	var method string
	if lock {
		method = "PUT"
	} else {
		method = "DELETE"
	}
	if rq, err = http.NewRequest(method, url, nil); err != nil {
		fmt.Println(err)
		return
	}

	rq.Header.Set("Authorization", "token "+si.token)

	rp, err := client.Do(rq)
	if err != nil {
		fmt.Println(err)
		return
	}
	if rp.StatusCode == http.StatusNoContent {
		fmt.Printf("issue %s:锁定(解锁)成功\n", si.issueN)
		return
	}

	fmt.Println(rp.Status)
}

const EXP = ".md"

func createOrUpdateAnIssue(si *SearchInfo, create bool) {
	//读取RequestBody，支持title和body两个字段
	fmt.Printf("issue title: ")
	var titlename string
	var err error
	_, err = fmt.Scan(&titlename)
	if err != nil {
		fmt.Println(err)
		return
	}

	filename := titlename + EXP

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	//首先关闭，否则用户无法打开文件
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//利用命令行打开默认文本编辑器
	if err := openFileWithDefaultEditor(filename); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("finish writing?(y n) ")
	var finish string
	_, err = fmt.Scan(&finish)
	if err != nil {
		fmt.Println(err)
		return
	}

	//建立文件的reader
	// file, err = os.Open(filename)
	filedata, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	//构建请求body，然后利用json包转成json格式
	rqmap := map[string]string{
		"title": titlename,
		"body":  string(filedata),
	}

	jbody, err := json.Marshal(rqmap)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", jbody)

	//建立rq
	var method string
	if si.mode == "c" {
		method = "POST"
	} else {
		method = "PATCH"
	}

	url := generateURL(*si)
	rq, err := http.NewRequest(method, url, bytes.NewReader(jbody))
	if err != nil {
		fmt.Println(err)
		return
	}
	rq.Header.Set("Authorization", "token "+si.token)
	rp, err := client.Do(rq)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rp.Body.Close()

	//处理结果,不为201

	if si.mode == "c" && rp.StatusCode != http.StatusCreated {
		fmt.Println(rp.Status)
		return
	}

	if si.mode == "u" && rp.StatusCode != http.StatusOK {
		fmt.Println(rp.Status)
		return
	}

	//解析body
	var bodyResult *Issue
	if err := json.NewDecoder(rp.Body).Decode(&bodyResult); err != nil {
		fmt.Println(err)
		return
	}

	//打印结果
	fmt.Println("Create Or Update success!")
	//for _, v := range bodyResult {
	fmt.Printf("number=%d user=%s state=%s title=%q\n", bodyResult.Number, bodyResult.User.Login, bodyResult.State, bodyResult.Title)
}

func openFileWithDefaultEditor(filename string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir + "\\" + filename)
	cmd := exec.Command("CMD", "/c", "start ", dir+"\\"+filename)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func generateURL(si SearchInfo) string {
	//拼接repo和owner
	ro := []string{si.owner, si.repo, "issues"}
	r := IssuesURL + strings.Join(ro, "/")
	if si.mode == "c" || si.mode == "g" {
		return r
	}

	r = r + "/" + si.issueN
	if si.mode == "u" {
		return r
	}

	r = r + "/" + "lock"
	return r
}
