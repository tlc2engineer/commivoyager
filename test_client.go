package main

//import (
//	"net"
//	"fmt"
//	"encoding/gob"
//	"time"
//)

//func main() {
//	conn,err:=net.Dial("tcp", "127.0.0.1:2233")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	err = gob.NewEncoder(conn).Encode("In vino veritas!")
//	if err != nil { fmt.Println(err) }
//	time.Sleep(100)
//	fmt.Println("End")
//}
