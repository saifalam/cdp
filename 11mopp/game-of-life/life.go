package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type GOL struct {
	board [][]int
	size  int
}

func adjacent_to(board GOL, i int, j int) int {
	var count, sk, ek, sl, el int

	if i > 0 {
		sk = i - 1
	} else {
		sk = i
	}

	if i+1 < board.size {
		ek = i + 1
	} else {
		ek = i
	}

	if j > 0 {
		sl = j - 1
	} else {
		sl = j
	}

	if j+1 < board.size {
		el = j + 1
	} else {
		el = j
	}

	for k := sk; k <= ek; k++ {
		for l := sl; l <= el; l++ {
			count += board.board[k][l]
		}
	}
	count -= board.board[i][j]
	return count
}

func play(board GOL) GOL {
	newboard := makeBoard(board.size)
	// for each cell, apply the rules of Life
	for i := 0; i < board.size; i++ {
		for j := 0; j < board.size; j++ {
			count := adjacent_to(board, i, j)
			if count == 2 {
				newboard.board[i][j] = board.board[i][j]
			}
			if count == 3 {
				newboard.board[i][j] = 1
			}
			if count < 2 {
				newboard.board[i][j] = 0
			}
			if count > 3 {
				newboard.board[i][j] = 0
			}
		}
	}
	return newboard
}

//print the life board
func printBoard(board GOL) {
	// for each row
	for j := 0; j < board.size; j++ {
		// print each column position
		for i := 0; i < board.size; i++ {
			if board.board[i][j] == 1 {
				fmt.Print("x")
			} else {
				fmt.Print("_")
			}
		}
		fmt.Print("\n")
	}
}
func makeBoard(size int) GOL {
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
	initialBoard := makeBoard(size)
	initialBoard.size = size
	for i := 0; i < size; i++ {
		if scanner.Scan() {
			ipText := scanner.Text()
			for j, x := range ipText {
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

	if err := reader.Err(); err != nil {
		log.Fatal(err)
	} else {
		steps, board := read_file(reader)
		printBoard(board)
		for i := 0; i < steps; i++ {
			board = play(board)
			printBoard(board)
		}
	}
}
