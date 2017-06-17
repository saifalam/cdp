package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
)

var splitTask = make(chan data)
var combineTask = make(chan data)

type GOL struct {
	board [][]int
	size  int
}

type data struct {
	splitBoard         [][]int
	chunkSize, chunkid int
}

func make_chunk(chunkid, chunkSize, col int) data {
	newchunk := data{}
	newchunk.chunkSize = chunkSize //middle rows (need to work with how may rows)
	newchunk.chunkid = chunkid
	newchunk.splitBoard = make([][]int, (chunkSize + 2))
	for i, _ := range newchunk.splitBoard {
		newchunk.splitBoard[i] = make([]int, col)
	}
	return newchunk
}

func assign_data(board GOL, chunkId, i, totalRow, chunkSize int) data {
	dataChunk := make_chunk(chunkId, chunkSize, board.size)

	if i-1 < 0 {
		dataChunk.splitBoard[0] = make([]int, board.size)
	} else {
		copy(dataChunk.splitBoard[0], board.board[i-1])
	}

	for r, c := i, 1; r < i+totalRow && r < board.size; r = r + 1 {
		copy(dataChunk.splitBoard[c], board.board[r])
		c = c + 1
	}

	if i+chunkSize+1 > board.size {
		dataChunk.splitBoard[chunkSize+1] = make([]int, board.size)
	} else {
		copy(dataChunk.splitBoard[chunkSize+1], board.board[i+chunkSize])
	}
	return dataChunk
}

func split_board(board GOL, splitTask chan data) {
	chunkId := 0
	chunkSize := runtime.NumCPU()
	for i := 0; i < board.size; i = i + chunkSize {
		splitTask <- assign_data(board, chunkId, i, chunkSize, chunkSize)
		chunkId = chunkId + 1
	}
}

func combine_board(board GOL, combineTask <-chan data) GOL {
	newBoard := make_board(board.size)
	total_chunk := int(math.Ceil(float64(float32(board.size) / float32(runtime.NumCPU()))))
	for i := 0; i < total_chunk; i++ {
		result := <-combineTask
		k := result.chunkid * result.chunkSize
		for r, s := k, 1; r < k+result.chunkSize && r < board.size; r++ {
			copy(newBoard.board[r], result.splitBoard[s])
			s = s + 1
		}
	}
	return newBoard
}

func adjacent_cell(board [][]int, col, row, i, j int) int {
	var count, sk, ek, sl, el int

	if i > 0 {
		sk = i - 1
	} else {
		sk = i
	}

	if i+1 < row {
		ek = i + 1
	} else {
		ek = i
	}

	if j > 0 {
		sl = j - 1
	} else {
		sl = j
	}

	if j+1 < col {
		el = j + 1
	} else {
		el = j
	}

	for k := sk; k <= ek; k++ { //rows
		for l := sl; l <= el; l++ { //cols
			count += board[k][l]
		}
	}
	count -= board[i][j]
	return count
}

func game_rules(input data) data {
	col := len(input.splitBoard[0])
	row := len(input.splitBoard)
	newChunk := make_chunk(input.chunkid, input.chunkSize, col)
	for i := 1; i <= input.chunkSize; i++ {
		for j := 0; j < col; j++ {
			count := adjacent_cell(input.splitBoard, col, row, i, j)

			if count == 2 {
				newChunk.splitBoard[i][j] = input.splitBoard[i][j]
			} else if count == 3 {
				newChunk.splitBoard[i][j] = 1
			} else if count < 2 {
				newChunk.splitBoard[i][j] = 0
			} else { //if count > 3
				newChunk.splitBoard[i][j] = 0
			}
		}
	}
	return newChunk
}

func worker(splitTask <-chan data, combineTask chan<- data) {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				input := <-splitTask
				result := game_rules(input)
				combineTask <- result
			}
		}()
	}
}

func play(board GOL) GOL {
	newBoard := make_board(board.size)
	go split_board(board, splitTask)
	newBoard = combine_board(board, combineTask)
	return newBoard
}

//print the life board
func printBoard(board GOL) {
	// for each row
	for j := 0; j < board.size; j++ {
		// print each column position
		for i := 0; i < board.size; i++ {
			if board.board[j][i] == 1 {
				fmt.Print("x")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}
func make_board(size int) GOL {
	newGOLBoard := GOL{}
	newGOLBoard.size = size
	newGOLBoard.board = make([][]int, size)
	for i, _ := range newGOLBoard.board {
		newGOLBoard.board[i] = make([]int, size)
	}
	return newGOLBoard
}

func read_file(scanner *bufio.Scanner) (int, GOL) {
	var size, steps int
	if scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d%d", &size, &steps)
	}
	initialBoard := make_board(size)
	initialBoard.size = size
	for i := 0; i < size; i++ {
		if scanner.Scan() {
			input := scanner.Text()
			for j, x := range input {
				if x == 'x' {
					initialBoard.board[i][j] = 1
				}
			}
		}
	}
	return steps, initialBoard
}

func main() {
	reader := bufio.NewScanner(os.Stdin)

	go worker(splitTask, combineTask)

	if err := reader.Err(); err != nil {
		log.Fatal(err)
	} else {
		steps, board := read_file(reader)
		for i := 0; i < steps; i++ {
			board = play(board)
			//printBoard(board)
		}
		printBoard(board)
	}
}
