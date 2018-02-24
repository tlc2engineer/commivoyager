package main

import (
	"github.com/my/commi/matrix"
	"testing"

	"fmt"
)

func TestCalc(t *testing.T) {
	points := matrix.GetPFromFile("./test_data/points2.txt")
	controlMoves, controlResult := matrix.GetMovesFromFile("./test_data/moves2.txt")
	res, it := matrix.GetResult(points)
	bestMoves, bestResult := res.Moves, res.Res
	fmt.Println("Число итераций", it)
	if matrix.IsResEquals(bestMoves, bestResult, controlMoves, controlResult) {
		fmt.Println("Результаты совпадают!")
	} else {
		t.Errorf("results not match %d!=%d", bestResult, controlResult)
	}
}

func BenchmarkCalc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		points := matrix.GetPFromFile("./test_data/points2.txt")
		matrix.GetResult(points)
	}

}
