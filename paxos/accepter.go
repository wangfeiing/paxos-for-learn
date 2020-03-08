package paxos

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"paxos/storage"
	"strconv"
)


var AcceptedProposerID int64 =  0
var AcceptedUrlList = []string{"http://127.0.0.1:8080/accepted"}
var acceptedValue = ""
var dataFilePath = ""

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