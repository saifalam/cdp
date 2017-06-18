package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type ProductionRule struct {
	left  string
	right []string
}

type Grammar struct {
	analyzeString, terminalSymbols, nonTerminalSymbols []string
	productions                                        []ProductionRule
	startSymbol                                        string
}

func read_input(scanner *bufio.Scanner) {
	grammar := Grammar{}
	productionRule := ProductionRule{}

	if scanner.Scan() {
		grammar.analyzeString = strings.Fields(scanner.Text())
	}
	if scanner.Scan() {
		grammar.terminalSymbols = strings.Fields(scanner.Text())

	}
	if scanner.Scan() {
		grammar.nonTerminalSymbols = strings.Fields(scanner.Text())
	}
	if scanner.Scan() {
		grammar.startSymbol = scanner.Text()
	}

	for scanner.Scan() {
		if scanner.Text() != "-" {
			content := strings.Split(scanner.Text(), ":")
			productionRule = ProductionRule{left: content[0], right: strings.Fields(content[1])}
			//fmt.Println("left rule: ", productionRule.left, "right rule: ", productionRule.right)
			grammar.productions = append(grammar.productions, productionRule)
		}

	}
	fmt.Println(grammar.analyzeString)
	fmt.Println(grammar.terminalSymbols)
	fmt.Println(grammar.nonTerminalSymbols)
	fmt.Println(grammar.startSymbol)
	fmt.Println(len(grammar.productions))

	for i := 0; i < len(grammar.productions); i++ {
		fmt.Println(grammar.productions[i])
	}

}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	if err := reader.Err(); err != nil {
		log.Fatal(err)
	} else {
		read_input(reader)
	}
}
