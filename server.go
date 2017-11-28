package main

import (
	//"time"
	"fmt"
	//"github.com/my/commi/matrix"
	"net"
	"github.com/my/commi/matrix"
	"encoding/gob"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"text/template"

	"github.com/ajstarks/svgo"


	//"strings"
	"math/rand"
	"strconv"
)

var timeCalc float64
var itime=0
var npoint=0
func worker(conn net.Conn){

	defer conn.Close()
	buffer:=make([]byte,1024)
	size,err:=conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Считано %d байт\n",size)
	s:=string(buffer)
	fmt.Println(s)
	points:=matrix.GetPFromString(s)
	bg:=time.Now()
	res,countIt:=matrix.GetResult(points)
	tn:=time.Since(bg).Seconds()
	bestMoves,bestResult:=res.Moves,res.Res
	out:=fmt.Sprintf("Результат: %5.2f Путь: %s  Число итераций: %d  Время: %f\n",bestResult,matrix.String(bestMoves),countIt,tn)
	fmt.Println(out)
	gob.NewEncoder(conn).Encode(out)

}

func main() {
	go func() {
		serv, err := net.Listen("tcp", ":2233")
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			conn, err := serv.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Accept and Go!")
			go worker(conn)
		}
	}()

	http.HandleFunc("/html/", genStaticHandler("html"))
	http.HandleFunc("/css/", genStaticHandler("css"))
	http.HandleFunc("/js/", genStaticHandler("js"))
	http.HandleFunc("/time", hTime)
	http.HandleFunc("/jtime",getTime)
	http.HandleFunc("/getSVG",getSVG)
	http.HandleFunc("/commi",commi)
	http.HandleFunc("/getCalcData",getCalcData)
	http.ListenAndServe(":8080", nil)
	//bg:=time.Now()
	//fmt.Println("Go")
	//points:=matrix.GetPFromFile("points.txt")
	//controlMoves,controlResult:=matrix.GetMovesFromFile("moves.txt")
	//res,countIt:=matrix.GetResult(points)
	//bestMoves,bestResult:=res.Moves,res.Res
	//fmt.Println("Время итоговое мс: ",time.Since(bg))
	//fmt.Printf("Результат %5.2f Путь: %s  Число итераций: %d\n",bestResult,matrix.String(bestMoves),countIt)
	//if matrix.IsResEquals(bestMoves,bestResult,controlMoves,controlResult){
	//	fmt.Println("Результаты совпадают!")
	//} else{
	//	fmt.Println("Результаты не совпадают!")
	//}

}


func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func genStaticHandler(t string) func(w http.ResponseWriter, r *http.Request){

	return func(w http.ResponseWriter, r *http.Request){
		name:=r.URL.Path[1:]
		fileName:="./"+name+".html"
		switch(t){
		case "css":{fileName="./"+name; w.Header().Set("Content-Type","text/css")}
		case "js": fileName="./"+name;w.Header().Set("Content-Type","text/javascript")
		}
		res,err:=ioutil.ReadFile(fileName)
		if err!=nil{
			fmt.Println(err)
		}
		fmt.Fprintf(w, "%s",string(res))

	}
}
func hTime(w http.ResponseWriter, r *http.Request){
	Time:=time.Now()
	UTime:=<-time.After(time.Duration(3600))
	WTime:=<-time.After(time.Duration(7200))
	times:=make([]LocalTime,10)
	times[0]=LocalTime{Time.Format("15:04:05"),"Московское время"}
	times[1]=LocalTime{UTime.Format("15:04:05"),"Киевское время"}
	times[2]=LocalTime{WTime.Format("15:04:05"),"Венское время"}

	t, _ := template.ParseFiles("./templ/tbase.html")
	data:=OutData{times[:3],true}
	t.Execute(w,data )
}

func getTime(w http.ResponseWriter, r *http.Request){
	Time:=time.Now()
	dat,_:=json.Marshal(Time)
	w.Header().Set("Content-Type","text/json")
	w.Write(dat)
}

type OutData struct{
	Times []LocalTime
	Descr bool
}

type LocalTime struct{
	Time string
	Descr string
}


func formatTime(t int64) string{
	tm:=time.Unix(t/1000,0)
	return fmt.Sprintf("%d-%d-%d",  tm.Year(),tm.Month(),tm.Day())
}

func f2(s string) string{
	if len(s)>40 {
		return s[:40]
	}
	return s
}

func getSVG(w http.ResponseWriter, req *http.Request) {
	 np,err:=strconv.Atoi(req.FormValue("num"))
	 if err!=nil{
	 	np=20}
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)

	points:=genPoints(np)
	//p0:=points[0]

	s.Start(500, 500)
	for _,p:=range points{
		s.Circle(int(p.X),int(p.Y),2, "fill:none;stroke:black")
		s.Text(int(p.X),int(p.Y),fmt.Sprintf("%d",p.Num),"stroke:red; fill: red")
	}
	bg:=time.Now()
	res,countIt:=matrix.GetResult(points)
	itime=countIt
	npoint=len(points)
	timeCalc=time.Since(bg).Seconds()
	moves:=res.Moves
	for _,move:=range moves{
		from:=points[move.From-1]
		to:=points[move.To-1]

		s.Line(int(from.X),int(from.Y),int(to.X),int(to.Y),"fill:none;stroke:blue")
	}
	//s.Circle(250, 250, 50, "fill:none;stroke:black")
	s.End()
}

func commi(w http.ResponseWriter, req *http.Request){
	t, _ := template.ParseFiles("./templ/commi.html")
	data:=commiData{"",0,0.0}
	t.Execute(w,data )
}

type commiData struct {
	Pic string
	Count int
	Time float64
}

/*
Генерирует заданное количество случайных точек
 */
func genPoints(n int) []matrix.Point{
	points:=make([]matrix.Point,0,n)
	for i:=0;i<n;i++{
		points=append(points, matrix.Point{i+1,rand.Float64()*500.0,rand.Float64()*500.0})
	}
	return points

}

func getCalcData(w http.ResponseWriter, req *http.Request){
	data:=CalcData{npoint,itime,timeCalc, }
	res,err:=json.Marshal(data)
	if err!=nil{fmt.Println(err)}
	w.Header().Set("Content-Type","text/json")
	w.Write(res)

}

type CalcData struct{
	N int
	Count int
	Time float64
}
