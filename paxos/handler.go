package paxos

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)
var curServePort = ""
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

func Start() {
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
	dataFilePath =  "./data/" + *dataPath
	curServePort = *port
	http.HandleFunc("/put" , Put )
	http.HandleFunc("/proposer" ,  Proposer)
	http.HandleFunc("/accepted" , Accepted)

	if err :=  http.ListenAndServe(":"+*port , nil); err != nil {
		fmt.Println(err)
	}
}