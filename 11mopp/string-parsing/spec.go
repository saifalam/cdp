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
	result        bool
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
	//fmt.Println(grammar.analyzeString)
	//fmt.Println(grammar.terminalSymbols)
	//fmt.Println(grammar.nonTerminalSymbols)
	//fmt.Println(grammar.startSymbol)
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

	/*for i := 0; i < len(grammar.productions); i++ {
		fmt.Println("From Outside loop: ", grammar.productions[i])
	}*/
	return grammar
}

func any_empty(s Stack) bool {
	if s.inputString == nil || !(len(s.inputString) > 0) {
		return true
	}
	if s.currentString == nil || !(len(s.currentString) > 0) {
		return true
	}
	return false
}
func reduce(s Stack) Stack {
	if any_empty(s) {
		return s
	}
	s.currentString, s.inputString = eliminate_common_prefix(s.currentString, s.inputString)
	fmt.Println("From reduce: ", s.currentString, " ", s.inputString)
	return s
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

func eliminate_common_prefix(currentString, inputString []string) ([]string, []string) {
	fmt.Println(currentString, "  ", inputString)
	commonPrefix := common_prefix(currentString, inputString)
	//fmt.Println(commonPrefix, "  ", len(commonPrefix))

	current := remove_prefix(commonPrefix, currentString)
	input := remove_prefix(commonPrefix, inputString)

	return current, input
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

func evaluate_string_parsing(grammar Grammar) Stack {
	eval := Stack{}
	eval.inputString = grammar.analyzeString
	for i := 0; i < len(grammar.productions); i++ {
		eval.currentString = grammar.productions[i].right

		//Rule no: 01 (Eliminate common prefix)
		eval = reduce(eval)
		fmt.Println("From Eval: ", eval.currentString, " ", eval.inputString)

		//Rule no: 02 (Is the initial symbol terminal)

		if begins_with_terminal(grammar, eval.currentString) {
			eval.result = false
			return eval
		}

		//Rule no: 03 (Both current and input string empty or not)
		if eval.inputString == nil {
			if eval.currentString == nil {
				eval.result = true
				return eval //return (Need to return 	<eval.productions, true>)
			}
			if !exists_non_terminal(grammar, eval.currentString) {
				eval.result = false
				return eval //return (Need to return 	<eval.productions, false>)
			}
		}
		//Rule no: 05 (Try all derivations recursively)
	}
	return eval
}
func parse_recursive_descent(g Grammar) Stack {
	eval := Stack{}
	eval.inputString = g.analyzeString
	// Need to get current String according to the starting g.startSymbol
	eval.currentString = g.productions[0].right
	return evaluate_string_parsing(g)
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	if err := reader.Err(); err != nil {
		log.Fatal(err)
	} else {
		grammar := read_input(reader)
		//fmt.Println(grammar)
		//eval := parse_recursive_descent(grammar)
		eval := evaluate_string_parsing(grammar)
		fmt.Println("From main: ", eval.productions, "  ", eval.result)
	}
}
