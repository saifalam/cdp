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

type Stack struct {
	currentString []string
	inputString   []string
	productions   []ProductionRule
}

type Grammar struct {
	analyzeString, terminalSymbols, nonTerminalSymbols []string
	productions                                        []ProductionRule
	startSymbol                                        string
}

func read_input(scanner *bufio.Scanner) Grammar {
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
	fmt.Println(grammar.analyzeString)
	fmt.Println(grammar.terminalSymbols)
	fmt.Println(grammar.nonTerminalSymbols)
	fmt.Println(grammar.startSymbol)
	for scanner.Scan() {
		if scanner.Text() != "-" {
			content := strings.Split(scanner.Text(), ":")
			productionRule = ProductionRule{left: content[0], right: strings.Fields(content[1])}
			//fmt.Println("left rule: ", productionRule.left, "right rule: ", productionRule.right)
			grammar.productions = append(grammar.productions, productionRule)
		}

	}
	/*fmt.Println(grammar.analyzeString)
	fmt.Println(grammar.terminalSymbols)
	fmt.Println(grammar.nonTerminalSymbols)
	fmt.Println(grammar.startSymbol)
	fmt.Println(len(grammar.productions))

	for i := 0; i < len(grammar.productions); i++ {
		fmt.Println(grammar.productions[i])
	}*/
	return grammar

}

func common_prefix(s1, s2 []string) []string {
	var commonPrefix []string
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return commonPrefix
		} else {
			commonPrefix = append(commonPrefix, s1[i])
		}
	}
	return commonPrefix

}

func eliminate_common_prefix(currentString, inputString []string) ([]string, []string) {
	var current []string
	var input []string

	fmt.Println(currentString, "  ", inputString)

	commonPrefix := common_prefix(currentString, inputString)
	//fmt.Println(commonPrefix, "  ", len(commonPrefix))

	if commonPrefix != nil && len(commonPrefix) > 0 {
		for i := 0; i < len(currentString); i++ {
			if i >= len(commonPrefix) {
				current = append(current, currentString[i])
			}
		}
		for i := 0; i < len(inputString); i++ {
			if i >= len(commonPrefix) {
				input = append(input, inputString[i])
			}
		}
		//fmt.Println(current, " ", input)
	}
	return current, input
}

func evaluate_string_parsing(grammar Grammar) {
	eval := Stack{}
	eval.inputString = grammar.analyzeString
	for i := 0; i < len(grammar.productions); i++ {
		eval.currentString = grammar.productions[i].right
		// Rule no: 01 (Eliminate common prefix)
		current, input := eliminate_common_prefix(eval.currentString, eval.inputString)
		fmt.Println(current, " ", input)
	}

}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	if err := reader.Err(); err != nil {
		log.Fatal(err)
	} else {
		grammar := read_input(reader)
		evaluate_string_parsing(grammar)
	}
}
