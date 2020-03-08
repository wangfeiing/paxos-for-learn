package paxos

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)
var BaseProposerID = time.Now().Unix()
var ProposerUrlList = []string{"http://127.0.0.1:8080/proposer"}
var lock sync.Mutex

func getProposerID() int64 {
	lock.Lock()
	defer lock.Unlock()
	BaseProposerID+=1
	return BaseProposerID
}

func Proposer(writer http.ResponseWriter , request *http.Request) {
	request.ParseForm()
	fmt.Println(request.Form)
	proposerId :=  request.FormValue("proposerId")
	val := request.FormValue("proposerValue")
	fmt.Println(proposerId)
	fmt.Println(val)
	if len(proposerId) <= 0 {
		writer.Write([]byte("no"))
		return
	}
	intProposerId , _  := strconv.Atoi(proposerId)
	if int64(intProposerId) < BaseProposerID {
		writer.Write([]byte("no"))
	}
	acceptedValue = val
	BaseProposerID = int64(intProposerId)
	writer.Write([]byte("yes"))
}

func requestProposer(proposerId int64 , proposerValue string ) bool {

	yesCount := 0
	if len(ProposerUrlList) <= 0 {
		return true
	}
	for _ , u := range ProposerUrlList {
		fmt.Println("url=", u)
		formData :=  url.Values{
			"proposerId": []string{strconv.Itoa(int(proposerId))},
			"proposerValue": []string{ proposerValue },
		}
		fmt.Println(formData)
		resp , err :=  http.PostForm(u,  formData)

		if err != nil {
			fmt.Println("postError,u=", u,",error=",err)
			continue
		}
		if  resp.StatusCode != http.StatusOK {
			fmt.Println("postError,u=", u,*resp)
			continue
		}
		defer  resp.Body.Close()
		bodyByte , err :=  ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("postError,ioutil.ReadAll,u=", u,err)
			continue
		}
		ret := string(bodyByte)
		fmt.Println(ret)
		if ret == "yes" {
			yesCount++
		}
	}
	if yesCount > len(ProposerUrlList)/2 {
		return true
	}
	return false
}
