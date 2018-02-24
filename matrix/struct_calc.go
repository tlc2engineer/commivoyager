package matrix

type Node struct {
	left   *Node
	right  *Node
	matrix *Matrix
	moves  *[]Move
}

func (node *Node) calc() Result {
	m := node.matrix
	moves := node.moves
	countIt++
	if m.size == 2 {
		*moves = h2matrix(m, moves)
		return Result{*moves, m.sum}
	}
	// Удаление замкнутых циклов
	if !removeClosePath(m, moves) {
		return Result{*moves, -1}
	}
	if !m.reduce() {
		return Result{*moves, -1}
	}
	// Если хуже лучшего результата возврат.
	if m.sum > bestResult && bestResult != -1 {
		return Result{*moves, -1}
	}
	// Оптимальный элемент ветвления
	r, c := m.getOptRC()
	// Копия перемещений для левой ветки
	cmoves := make([]Move, len(*moves))
	copy(cmoves, *moves)
	// Создаем копию матрицы для левой ветки
	cm := m.copy()
	//Удаляем элемент для правой ветки
	m.removeNode(r, c)
	// Добавляем ход
	*moves = append(*moves, Move{r, c})
	// -1 для левой ветки
	cm.setEl(r, c, -1)
	var res1, res2 Result
	eval := func(mt *Matrix, mvs []Move, ch chan Result) {
		//time.Sleep(2)
		ch <- mainRecFuncT(mt, &mvs)
	}
	if gorNumber < maxGorNumber && false {
		ch1 := make(chan Result)
		go eval(&cm, cmoves, ch1)
		gorNumber++
		res2 = mainRecFuncT(m, moves)
		res1 = <-ch1
		gorNumber--
		close(ch1)

	} else {
		node.right = &Node{nil, nil, m, moves}
		res2 = node.right.calc()
		node.right = nil
		node.left = &Node{nil, nil, &cm, &cmoves}
		res1 = node.left.calc()
		node.left = nil
		m = nil
		moves = nil
		cm.data = nil
		cmoves = nil
		//node.right=nil
		//node.left=nil
	}
	if res1.Res == -1 && res2.Res == -1 {
		return res1
	}
	if res1.Res == -1 {
		return res2
	}
	if res2.Res == -1 {
		return res1
	}
	if res1.Res < res2.Res {
		return res1
	}
	return res2
}
