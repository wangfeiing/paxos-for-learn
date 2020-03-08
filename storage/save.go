package storage

import (
	"fmt"
	"os"
)

func Save(filePath ,key string , data string ) error  {
	logData := fmt.Sprintf("%s\t%s\n",key , data)
	writeToFlie(filePath , logData)
	return nil
}

func writeToFlie(filePath string , data string ) {

	var f *os.File
	f , _  = os.OpenFile(filePath, os.O_CREATE | os.O_APPEND | os.O_RDWR , 0660)
	fmt.Println(*f)
	defer f.Close()
	_ , err :=  f.WriteString(data)
	if err != nil {
		fmt.Println("fwriteStringError=",err)
		return
	}
	return
}