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
			//fmt.Println("scanner Text: ", scanner.Text())
		}
		//fmt.Println(grammar.productions)
	}

	for i := 0; i < len(grammar.productions); i++ {
		fmt.Println("From Outside loop: ", grammar.productions[i])
	}
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

func remove_prefix(prefix, s []string) []string {
	var result []string
	if prefix != nil && len(prefix) > 0 {
		for i := 0; i < len(s); i++ {
			if i >= len(prefix) {
				result = append(result, s[i])
			}
		}
	}
	return result
}

func eliminate_common_prefix(currentString, inputString []string) []string {
	fmt.Println(currentString, "  ", inputString)
	commonPrefix := common_prefix(currentString, inputString)
	//fmt.Println(commonPrefix, "  ", len(commonPrefix))

	current := remove_prefix(commonPrefix, currentString)
	/*if commonPrefix != nil && len(commonPrefix) > 0 {
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
	}*/
	return current
}

func is_terminal(g Grammar, s []string) bool {
	for i := 0; i < len(g.terminalSymbols); i++ {
		if g.terminalSymbols[i] == s[0] {
			//fmt.Println("From is terminal: ", g.terminalSymbols[i])
			return true
		}
	}
	return false
}

func is_nonterminal(g Grammar, s []string) bool {
	return !is_terminal(g, s)
}

func exists_non_terminal(g Grammar, s []string) bool {
	if s != nil && len(s) > 0 {
		if is_nonterminal(g, s) {
			return true
		}
	}
	return false
}

func begins_with_terminal(g Grammar, s []string) bool {
	if s != nil && len(s) > 0 {
		if is_terminal(g, s) {
			return true
		}
	}
	return false
}

func evaluate_string_parsing(grammar Grammar) {
	eval := Stack{}
	eval.inputString = grammar.analyzeString
	for i := 0; i < len(grammar.productions); i++ {
		eval.currentString = grammar.productions[i].right

		//Rule no: 01 (Eliminate common prefix)
		eval.currentString = eliminate_common_prefix(eval.currentString, eval.inputString)
		fmt.Println(eval.currentString)

		//Rule no: 02 (Is the initial symbol terminal)
		res := begins_with_terminal(grammar, eval.currentString)
		fmt.Println("Is terminal result: ", res)

		//Rule no: 03 (Both current and input string empty or not)
		if eval.inputString == nil {
			if eval.currentString == nil {
				//return (Need to return 	<eval.productions, true>)
			}
			if !exists_non_terminal(grammar, eval.currentString) {
				//return (Need to return 	<eval.productions, false>)
			}
		}

		//Rule no: 05 (Try all derivations recursively)

	}
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	if err := reader.Err(); err != nil {
		log.Fatal(err)
	} else {
		grammar := read_input(reader)
		fmt.Println(grammar)
		evaluate_string_parsing(grammar)
	}
}
