package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type ProductionRule struct {
	left  []string
	right []string
}

type Stack struct {
	currentString []string
	inputString   []string
	productions   []ProductionRule
	result        bool
}

type Grammar struct {
	startSymbol                                        []string
	analyzeString, terminalSymbols, nonTerminalSymbols []string
	productions                                        []ProductionRule
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
		grammar.startSymbol = strings.Fields(scanner.Text())
	}
	for scanner.Scan() {
		if strings.Trim(scanner.Text(), "\n\t") != "-" {
			content := strings.Split(strings.Trim(scanner.Text(), "\n\t"), ":")
			productionRule = ProductionRule{left: strings.Fields(content[0]), right: strings.Fields(content[1])}
			grammar.productions = append(grammar.productions, productionRule)
		} else {
			break
		}
	}
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
	return s
}

func common_prefix(s1, s2 []string) []string { // s1 = currentString && s2 = inputString
	var commonPrefix []string
	var lenght int
	if len(s1) >= len(s2) {
		lenght = len(s2)
	} else {
		lenght = len(s1)
	}
	for i := 0; i < lenght; i++ {
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
	commonPrefix := common_prefix(currentString, inputString)
	if (commonPrefix != nil) && (len(commonPrefix) > 0) {
		current := remove_prefix(commonPrefix, currentString)
		input := remove_prefix(commonPrefix, inputString)
		return current, input
	}
	return currentString, inputString
}

func is_terminal(g Grammar, s []string) bool {
	for i := 0; i < len(g.terminalSymbols); i++ {
		if g.terminalSymbols[i] == s[0] {
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
		for i := 0; i < len(g.nonTerminalSymbols); i++ {
			if g.nonTerminalSymbols[i] == s[0] {
				return true
			}
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

func to_string(strArray []string) string {
	s := strings.Join(strArray, " ")
	return s
}

func evaluate_leftmost_derivation(g Grammar, s []string) []string {
	var LMDerivations []string
	for i := 0; i < len(g.productions); i++ {
		if to_string(s) == to_string(g.productions[i].left) {
			LMDerivations = append(LMDerivations, to_string(g.productions[i].right))
		}
	}
	return LMDerivations
}

func find_left_most_nonterminal(g Grammar, currentString []string) []string {
	var LMNonTerminal []string
	for i := 0; i < len(currentString); i++ {
		if is_nonterminal(g, strings.Fields(strings.Trim(currentString[i], "\n\t"))) {
			LMNonTerminal = append(LMNonTerminal, strings.Trim(currentString[i], "\n\t"))
			return LMNonTerminal
		}
	}
	return LMNonTerminal
}

func get_all_strings_after_nonterminal(g Grammar, s []string) []string {
	var result []string
	for i := 1; i < len(s); i++ {
		result = append(result, s[i])
	}
	return result
}

func remove_all_rules(s Stack) Stack {
	if s.productions != nil {
		emptySlice := make([]string, 0)
		for i := 0; i < len(s.productions); i++ {
			s.productions[i].left = emptySlice
			s.productions[i].right = emptySlice
		}
	}
	return s
}

func evaluate_string_parsing(grammar Grammar, eval Stack) Stack {

	// Rule no: 01 (Eliminate common prefix)
	reducedEval := reduce(eval)

	//Rule no: 02 (Is the initial symbol terminal)
	if begins_with_terminal(grammar, reducedEval.currentString) {
		reducedEval.result = false
		return reducedEval
	}

	// Rule no: 03 (Both current and input string empty or not)
	if reducedEval.inputString == nil || to_string(reducedEval.inputString) == "" {
		if reducedEval.currentString == nil || to_string(reducedEval.currentString) == "" {
			reducedEval.result = true
			return reducedEval //return (Need to return 	<eval.productions, true>)
		}
		// Rule no: 04
		if !exists_non_terminal(grammar, reducedEval.currentString) {
			reducedEval.result = false
			return reducedEval //return (Need to return 	<eval.productions, false>)
		}
	}
	// Rule no: 05 (Try all derivations recursively)
	LMNonTerminal := find_left_most_nonterminal(grammar, reducedEval.currentString)

	if LMNonTerminal != nil {
		LMDerivations := evaluate_leftmost_derivation(grammar, LMNonTerminal)

		if (LMDerivations != nil) && (len(LMDerivations) > 0) {
			restPart := get_all_strings_after_nonterminal(grammar, reducedEval.currentString)

			for i := 0; i < len(LMDerivations); i++ {
				reducedEval.currentString = strings.Fields(LMDerivations[i])

				if restPart != nil && (len(restPart) > 0) {
					for j := 0; j < len(restPart); j++ {
						reducedEval.currentString = append(reducedEval.currentString, restPart[j])
					}

				}

				// Rule no: 6 (recursive calling)
				evaluation := evaluate_string_parsing(grammar, reducedEval)
				if evaluation.result {
					productionRule := ProductionRule{left: LMNonTerminal, right: strings.Fields(LMDerivations[i])}
					evaluation.productions = append(evaluation.productions, productionRule)
					return evaluation
				} else {
					evaluate_string_parsing(grammar, evaluation)
				}
			}
		}
	}
	return reducedEval
}

func parse_recursive_descent(g Grammar) Stack {
	eval := Stack{}
	eval.inputString = g.analyzeString
	eval.currentString = g.startSymbol
	return evaluate_string_parsing(g, eval)
}

func print_output(s Stack) {
	if s.result {
		for i := len(s.productions) - 1; i >= 0; i-- {
			fmt.Println(to_string(s.productions[i].left), ":", to_string(s.productions[i].right))
		}
		fmt.Println("SUCCESS")
	} else {
		fmt.Println("FAILED")
	}
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	if err := reader.Err(); err != nil {
		log.Fatal(err)
	} else {
		grammar := read_input(reader)
		eval := parse_recursive_descent(grammar)
		print_output(eval)
	}
}
