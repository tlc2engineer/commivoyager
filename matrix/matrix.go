package matrix

import (

	"math"
	"io/ioutil"
	"log"
	"strings"
	"strconv"
	"fmt"
)
var bestResult float64=-1
var bestMoves []Move;
var countIt =0

/**
Установка лучшего результата
 */
func setBestResult(result float64,moves []Move){
	if bestResult==-1 || result<bestResult{
		bestResult=result
		bestMoves=moves

	}
}



type  Matrix struct{
	data [][]float64
	rowNums,columnNums []int
	size int
	sum float64
}
/*
* Создание пустой матрицы
  */
func Create(size int) Matrix {
	o:=make([][]float64,size);
	for i:=0;i<size;i++{
		o[i]=make([]float64,size)
	}
	return Matrix{o,make([]int,size),make([]int,size),size,0}
}

/*
Создание матрицы из массива данных
 */
func CreateFromPoint(points []Point) Matrix {
	size:=len(points)
	o:=Create(size)
	for i:=0;i<size;i++{
		o.rowNums[i]=points[i].Num
		o.columnNums[i]=points[i].Num
	}
	for row:=0;row<size;row++{
		for column:=0;column<size;column++{
			if row==column {
				o.set(row, column, -1)
			}else{
				o.set(row, column, distance(points[row],points[column]))
			}

		}
	}
	return o

}
/*
Установка элемента по номерам колонок и сторок.
 */
func (m*Matrix) set(row int,column int,data float64){
	m.data[row][column]=data;
}
/*
Получение элемента по номерам колонок и сторок.
 */
func (m Matrix) get(row int,column int) float64{
	return m.data[row][column]
}
/*
Точка структура
 */
type Point struct{
	Num int
	X, Y float64
}
/*
Расстояние между точками
 */
func distance(p1,p2 Point) float64{
	return math.Sqrt(math.Pow(p1.X-p2.X,2)+math.Pow(p1.Y-p2.Y,2))
}
/**
нахождение минимума массива и уменьшение всех значений на минимум. Если все значения  -1 возврат -1
 */
func reduceArr(arr []float64) float64{
	norm:=false
	min:=-1.0
	for _,v:=range arr{
		if v>=0.0{
			min=v
			norm=true
			break
		}
	}
	if !norm{
		return -1.0
	}
	for _,v:=range arr{
		if v>-1 && v<min{
			min=v
		}
	}
	for i:=0;i<len(arr);i++{
		if arr[i]>-1 {
			arr[i] = arr[i] - min
		}
	}
	return min
}
/*
Функция приведения без создания новой матрицы
 */
func  (m*Matrix) reduce() bool {
	sum := 0.0
	for i := 0; i < m.size; i++ {
		add := reduceArr(m.data[i])
		if add < 0 {

			return false
		}
		sum += add
	}
	for i:=0;i<len(m.columnNums);i++{
		col:=m.getColumn(m.columnNums[i])
		add:=reduceArr(col)
		if add<0 {

			return false
			}
		sum+=add
		m.setColumn(col,m.columnNums[i])
	}

	m.sum+=sum
	return true
}
/*
Функция приведения с созданием новой матрицы
 */
func (m* Matrix) reduceM() Matrix{
	matr:=m.copy()
	sum:=0.0
	for i := 0; i < matr.size; i++ {
		add := reduceArr(matr.data[i])
		if add < 0 {
			return Matrix{}
		}
		sum += add
	}
	//
	for i:=0;i<len(matr.columnNums);i++{
		col:=matr.getColumn(matr.columnNums[i])
		add:=reduceArr(col)
		if add<0 {return Matrix{}}
		sum+=add
		matr.setColumn(col,matr.columnNums[i])
	}
	matr.sum+=sum
	return matr
}
func (m* Matrix) getRow(num int) []float64{
	nr:=m.getAbsNRow(num)
	return m.data[nr];
}
/*
Возвращает массив колонки.
num- номер точки
 */
func (m*Matrix) getColumn(num int) []float64{
	i:=m.getAbsNColumn(num)
	// в i номер точки
	o:=make([]float64,m.size) // новый срез
	for j:=0;j<len(m.rowNums);j++{
		o[j]=m.data[j][i]

	}
	return o
}
/*
Установка значения колонки
 */
func (m*Matrix) setColumn(col []float64,num int){
	i:=m.getAbsNColumn(num)
	for j:=0;j<len(m.rowNums);j++{
		m.data[j][i]=col[j]

	}
}
/**
Строковое представление матрицы
 */
func (m* Matrix) String() string{
	s:="  "
	for i:=0;i<len(m.columnNums);i++{
		s+=fmt.Sprintf("       %d",m.columnNums[i])
	}
	s+="\n"
	for i:=0;i<len(m.rowNums);i++{
		s+=fmt.Sprintf("%d ",m.rowNums[i])
		for j:=0;j<len(m.columnNums);j++{
			s+=fmt.Sprintf("   %3.3f",m.data[i][j])
		}
		s+="\n"
	}
	return s
}
/*
Копия матрицы
 */
func (m* Matrix) copy() Matrix{
	data:=make([][]float64,len(m.data))
	for i:=0;i<len(data);i++{
		data[i]=make([]float64,m.size)
		copy(data[i],m.data[i])
	}
	rowNums:=make([]int,len(m.rowNums))
	columnNums:=make([]int,len(m.columnNums))
	copy(rowNums,m.rowNums)
	copy(columnNums,m.columnNums)
	return Matrix{data,rowNums,columnNums,m.size,m.sum}
}

/*
Получение оптимальной точки для ветвления
 */
func (m* Matrix) getOptRC() (int,int) {
	 optRow:=0;
	 optColumn:=0;
	 optDist:=-1.0;
	for row:=0;row<len(m.data);row++ {
		for column := 0; column <len(m.data); column++ {
			if(m.data[row][column]==0.0){
				 arow:=m.getRow(m.rowNums[row])
				 acolumn:=m.getColumn(m.columnNums[column])
				// Устанавливаем в  бесконечность точку
				arow[column]=-1;
				acolumn[row]=-1;
				minr:=getArrMin(arow)
				minc:=getArrMin(acolumn)
				if (minr+minc)>optDist {
					optDist = minr+minc;
					optRow=row;
					optColumn=column;
				}
				arow[column]=0;
				acolumn[row]=0;
			}
		}
	}
	return m.rowNums[optRow],m.columnNums[optColumn]
}
/*
Получение элементов по порядковым номерам точек
 */
func (m* Matrix) getEl(row int,column int) float64{
	nr:=m.getAbsNRow(row);
	nc:=m.getAbsNColumn(column);
	if nr>=len(m.rowNums) || nc>=len(m.columnNums) {panic("Длина не соответствует!")};
	return m.data[nr][nc];
}

/*
Получение минимального значения массива
 */
 func getArrMin(slice []float64) float64{
	 min:=-1.0
	 for i:=0;i<len(slice);i++{
		 if slice[i]==-1 {continue}
		 if slice[i]<min || min==-1.0 {
			 min=slice[i]
		 }
	 }
	 if min==-1 {min=0}
	 return min
 }

 /*
 Удаление строки
  */
func (m* Matrix) removeRow(row int)   {
	 nr:=m.getAbsNRow(row);
	 if nr<len(m.rowNums)-1{
		 copy(m.data[nr:],m.data[nr+1:])
		 copy(m.rowNums[nr:],m.rowNums[nr+1:])
	 }
	 m.data=m.data[:len(m.data)-1]
	 m.rowNums=m.rowNums[:len(m.rowNums)-1]

}
/*
Удаление колонки
 */
func (m* Matrix) removeColumn(col int)   {
	nc:=m.getAbsNColumn(col)
	if nc<len(m.columnNums)-1{
		copy(m.columnNums[nc:],m.columnNums[nc+1:])
		for i:=0;i<len(m.rowNums);i++ {
			copy(m.data[i][nc:], m.data[i][nc+1:])
		}
	}
	m.columnNums=m.columnNums[:len(m.columnNums)-1]
	for i:=0;i<len(m.rowNums);i++ {
		m.data[i]= m.data[i][:len(m.data)]
	}

}

/*
Удаление элемента - колонки и строки
 */

 func (m* Matrix) removeNode(row int,col int){
 	m.removeRow(row)
 	m.removeColumn(col)
 	m.size=m.size-1
 	if m.exists(col,row){
 		m.setEl(col,row,-1)
	}
 }
/*
Получение абсолютного номера строки по номеру точки
*/
func (m* Matrix) getAbsNRow(nrow int) int{
	nr:=0;
	for nr=0;nr<len(m.rowNums);nr++{
		if m.rowNums[nr]==nrow {break}
	}
	return nr
}

func (m* Matrix) exists(row int,col int) bool{
	nr:=m.getAbsNRow(row)
	nc:=m.getAbsNColumn(col)
	return nr<len(m.rowNums)&& nc<len(m.columnNums)
}
/*
Получение абсолютного номера колонки по номеру точки
 */
func (m* Matrix) getAbsNColumn(ncol int) int{
	nc:=0
	for nc=0;nc<len(m.columnNums);nc++{
		if m.columnNums[nc]==ncol {break}
	}
	return nc
}
/**
Установка элементов по номерам точек
 */
func (m* Matrix) setEl(row int,col int,val float64){
	nr,nc:=m.getAbsNRow(row),m.getAbsNColumn(col)

	if nr>=len(m.rowNums) || nc>=len(m.columnNums){
		panic("Длина больше заданной")
	}
	m.data[nr][nc]=val
}


type Move struct{
	From, To int
}
/**
Проверка замыкания цепи ходов
 */
func isClosePath(moves []Move,bg int,end int) bool{
	for _,move:=range moves{
		if move.From ==end{
			if move.To ==bg {return true}
			return isClosePath(moves,bg,move.To);
		}
	}
	return false
}
/**
Установка переходов по которым будут замыкания в -1.
 */
func removeClosePath(matrix* Matrix,moves []Move) bool{
	for row:=0;row<len(matrix.data);row++ {
		for column:=0;column<len(matrix.data);column++{
			if matrix.getEl(matrix.rowNums[row],matrix.columnNums[column])==-1 {continue};
			if(isClosePath(moves,matrix.rowNums[row],matrix.columnNums[column])){
				matrix.setEl(matrix.rowNums[row],matrix.columnNums[column],-1);
			}
		}

	}
	return true

}

/*
Обработка матрицы с размером 2
 */

 func h2matrix(m* Matrix,moves []Move) []Move{
 	if m.size!=2 {panic("Размер матрицы не равен 2!")}
 	if !m.reduce() {return moves}
	  c1,c2,r1,r2:=m.columnNums[0],m.columnNums[1],m.rowNums[0],m.rowNums[1]
	  v1,v2,v3,v4:=m.getEl(r1,c1),m.getEl(r2,c2),m.getEl(r2,c1),m.getEl(r1,c2)
	 if v1!=-1 && v2!=-1 && v3!=-1 && v4!=-1 {
		 if (v1+v2)<(v3+v4){
			moves=append(moves, Move{r1,c1})
			moves=append(moves, Move{r2,c2})
			setBestResult(m.sum,moves)
			return moves
		 }else{
			 moves=append(moves, Move{r2,c1})
			 moves=append(moves, Move{r1,c2})
			 setBestResult(m.sum,moves)
			 return moves
		 }
	 }
	 if v1!=-1 && v2!=-1 {
		 moves=append(moves, Move{r1,c1})
		 moves=append(moves, Move{r2,c2})
		 setBestResult(m.sum,moves)
		 return moves
	 }
	 if v3!=-1 && v4!=-1{
		 moves=append(moves, Move{r2,c1})
		 moves=append(moves, Move{r1,c2})
		 setBestResult(m.sum,moves)
		 return moves
	 }
	 return moves;
 }
/*
Основная рекурсивная функция
 */
 func mainRecFunc(m* Matrix,moves []Move){
 	countIt++
 	if m.size==2{
		h2matrix(m,moves)
		return
	}
	 // Удаление замкнутых циклов
	 if !removeClosePath(m, moves) {return}
	 if !m.reduce() {return}
	 // Если хуже лучшего результата возврат.
	 if m.sum>bestResult && bestResult!=-1 {return};
	 // Оптимальный элемент ветвления
	  r,c := m.getOptRC()
	 // Копия перемещений для левой ветки
	 cmoves:=make([]Move,len(moves))
	 copy(cmoves,moves)
	 // Создаем копию матрицы для левой ветки
	 cm:=m.copy()
	 //Удаляем элемент для правой ветки
	 m.removeNode(r,c)
	 // Добавляем ход
	 moves=append(moves, Move{r,c})
	 // -1 для левой ветки
	 cm.setEl(r,c,-1)
	 mainRecFunc(m,moves)
	 mainRecFunc(&cm,cmoves)
 }

 /*
 Рекурсивная многопоточная функция
  */
func mainRecFuncT(m* Matrix,moves []Move) Result {
	countIt++
	if m.size==2{
		moves=h2matrix(m,moves)
		return Result{moves,m.sum}
	}
	// Удаление замкнутых циклов
	if !removeClosePath(m, moves) {return Result{moves,-1}}
	if !m.reduce() {return Result{moves,-1}}
	// Если хуже лучшего результата возврат.
	if m.sum>bestResult && bestResult!=-1 {return Result{moves,-1}};
	// Оптимальный элемент ветвления
	r,c := m.getOptRC()
	// Копия перемещений для левой ветки
	cmoves:=make([]Move,len(moves))
	copy(cmoves,moves)
	// Создаем копию матрицы для левой ветки
	cm:=m.copy()
	//Удаляем элемент для правой ветки
	m.removeNode(r,c)
	// Добавляем ход
	moves=append(moves, Move{r,c})
	// -1 для левой ветки
	cm.setEl(r,c,-1)
	ch1:=make(chan Result)
	//ch2:=make(chan Result)
	eval:= func(c chan Result,mt* Matrix,mvs []Move){
		//time.Sleep(2)
		c<- mainRecFuncT(mt,mvs)
	}
	go eval(ch1,&cm,cmoves)
	res2:= mainRecFuncT(m,moves)
	res1:=<-ch1
	if res1.Res==-1 && res2.Res==-1 {return res1}
	if res1.Res==-1 {return res2}
	if res2.Res==-1 {return res1}
	if res1.Res<res2.Res { return res1}
	return res2
}

/*
Вывод последовательности ходов в текстовом виде
 */
 func  String(moves []Move) string{
 	s:=""
 	for _,move:=range moves{
 		s+=fmt.Sprintf("%d-%d ",move.From,move.To)
	 }
	 return s
 }
 /*
 Ввод точек из текстового файла
  */

func GetPFromFile(name string)  []Point{
	data,err:=ioutil.ReadFile(name);
	if err!=nil{
		log.Fatal("Неправильно считаны точки.")
	}
	s:=string(data)
	lines:=strings.Split(s,"\n")
	size:=len(lines)
	points:=make([]Point,0,size)
	for _,line:=range lines{
		d:=strings.Split(line," ")
		if d[0]==""{continue}
		num,err:=strconv.Atoi(d[0])
		if err!=nil{
			log.Fatal(err)
			}
		s:=strings.Replace(d[1],",",".",1)
		x,err:=strconv.ParseFloat(s,64)
		if err!=nil{log.Fatal(err)}
		s=strings.TrimSpace(strings.Replace(d[2],",",".",1))
		y,err:=strconv.ParseFloat(s,64)
		if err!=nil{log.Fatal(err)}
		points=append(points,Point{num,x,y})

	}
	return points
}

func GetMovesFromFile(name string) ([]Move,float64){
	data,err:=ioutil.ReadFile(name);
	result:=0.0
	if err!=nil{
		log.Fatal("Неправильно считаны переходы.")
	}
	s:=string(data)
	lines:=strings.Split(s,"\n")
	size:=len(lines)
	moves:=make([]Move,0,size)
	for _,line:=range lines{
		dat:=strings.Split(line," ")
		if len(dat)<2 {continue}
		froms:=string(dat[0])
		if froms!="Result:"{
			tos:=strings.TrimSpace(string(dat[1]))
			from,err:=strconv.Atoi(froms)
			if err!=nil{
				log.Fatal(err)
			}
			to,err:=strconv.Atoi(tos)
			if err!=nil{
				log.Fatal(err)
			}
			moves=append(moves, Move{from,to})
		}else{
			s=strings.TrimSpace(strings.Replace(dat[1],",",".",1))
			result,err=strconv.ParseFloat(s,64)
			if err!=nil{log.Fatal(err)}
		}
	}
	return moves,result
}

func IsResEquals(moves []Move,result float64,targetMoves []Move,targetResult float64) bool{
	flag:=false
	for _,move:=range moves{
		flag=false
		for _,movet:=range targetMoves{
			if (move.From ==movet.From && move.To ==movet.To) || (move.From ==movet.To && move.To ==movet.From){
				flag=true
				break
			}
		}
		if flag==false {return false}
	}
	return math.Abs(result-targetResult)<0.01

}

type Result struct {
	Moves []Move
	Res   float64

}

func GetResult(points []Point) (Result,int){
	countIt=0
	m:=CreateFromPoint(points)
	moves:=make([]Move,0,len(points))
	bestResult=-1
	return mainRecFuncT(&m,moves),countIt

}
/**
Получение списка точек из строки
 */
func GetPFromString(s string)  []Point{
	lines:=strings.Split(s,"\n")
	size:=len(lines)
	points:=make([]Point,0,size)
	for _,line:=range lines{
		d:=strings.Split(line," ")
		if d[0]==""{continue}
		num,err:=strconv.Atoi(d[0])
		if err!=nil{
			continue
		}
		s:=strings.Replace(d[1],",",".",1)
		x,err:=strconv.ParseFloat(s,64)
		if err!=nil{log.Fatal(err)}
		s=strings.TrimSpace(strings.Replace(d[2],",",".",1))
		y,err:=strconv.ParseFloat(s,64)
		if err!=nil{log.Fatal(err)}
		points=append(points,Point{num,x,y})

	}
	return points
}