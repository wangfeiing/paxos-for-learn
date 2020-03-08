package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)



func init() {

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



