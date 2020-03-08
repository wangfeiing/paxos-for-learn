package paxos

import "net/http"

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

