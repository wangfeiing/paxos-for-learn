package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"paxos/storage"
	"strconv"
	"sync"
	"time"
)

var BaseProposerID = time.Now().Unix()
var AcceptedProposerID int64 =  0
var ProposerUrlList = []string{"http://127.0.0.1:8080/proposer"}
var AcceptedUrlList = []string{"http://127.0.0.1:8080/accepted"}
var acceptedValue = ""
var lock sync.Mutex
var dataFilePath = ""
var curServePort = ""

func init() {

}

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
func Accepted(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	proposerId :=  request.FormValue("proposerId")
	val := request.FormValue("proposerValue")
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
	AcceptedProposerID = BaseProposerID
	storage.Save( dataFilePath,proposerId , val)
	writer.Write([]byte("yes"))
}
func Put(writer http.ResponseWriter , request * http.Request) {
	request.ParseForm()
	data := request.FormValue("data")
	fmt.Println(data)
	tryCount := 5
	for {
		tryCount--
		if tryCount<= 0 {
			break
		}
		proposerId := getProposerID()
		if ok :=  requestProposer(proposerId , data) ; !ok {
			continue
		}
		if yes := requestAccepted(proposerId , data) ; !yes{
			continue
		}
		writer.Write([]byte("OK"))
		return
	}
	writer.Write([]byte("FAIL"))
	return
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

func requestAccepted(proposerId int64 , proposerValue string ) bool {
	yesCount := 0
	if len(AcceptedUrlList) <= 0 {
		return true
	}
	for _ , u := range AcceptedUrlList {
		formData :=  url.Values{
			"proposerId": []string{strconv.Itoa(int(proposerId))},
			"proposerValue": []string{proposerValue},
		}
		resp , err :=  http.PostForm(u,  formData)
		if err != nil {
			fmt.Println("postError,u=", u)
			continue
		}
		if  resp.StatusCode != http.StatusOK {
			continue
		}
		defer  resp.Body.Close()
		bodyByte , err :=  ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		ret := string(bodyByte)
		if ret == "yes" {
			yesCount++
		}
	}
	if yesCount > len(ProposerUrlList)/2 {
		return true
	}
	return false
}

func main() {

	if len(os.Args) < 0 {
		fmt.Println("You should input port")
		return
	}
	dataPath :=  flag.String("d","" ,"Data path for user")
	port := flag.String("p" , "" , "Compute port")
	flag.Parse()

	fmt.Printf("serverListenPort=%+v,dataPath=%+v",*port,*dataPath)
	if len(*port) <= 0 {
		fmt.Println("Input port")
		return
	}
	if len(*dataPath) <= 0 {
		fmt.Println("Input data path")
		return
	}
	dataFilePath = *dataPath
	curServePort = *port
	http.HandleFunc("/put" , Put )
	http.HandleFunc("/proposer" ,  Proposer)
	http.HandleFunc("/accepted" , Accepted)

	if err :=  http.ListenAndServe(":"+*port , nil); err != nil {
		fmt.Println(err)
	}
}



