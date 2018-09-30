package main

import (
	"fmt"
	"os"
	"strconv"
)

var arg string
var position int

func getChar() string {
	l := len(arg)
	if l <= position {
		return "error"
	}
	s := string(arg[position])
	position++
	return s
}

const (
	tINT    = "tINT"
	tADD    = "tADD"
	tSUB    = "tSUB"
	tMUL    = "tMUL"
	tLPAREN = "tLPAREN"
	tRPAREN = "rTPAREN"
	tEND    = "tEND"
)

// Token ...
type Token struct {
	Type     string
	IntValue int
}

func isDigit(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	return n, err == nil
}

func lex() Token {
	var token Token

	c := getChar()

	if n, ok := isDigit(c); ok {
		for {
			d := getChar()
			m, ok := isDigit(d)
			if !ok {
				break
			}
			n = n*10 + m
		}
		token.Type = tINT
		token.IntValue = n
	} else if c == "+" {
		token.Type = tADD
	} else if c == "-" {
		token.Type = tSUB
	} else if c == "*" {
		token.Type = tMUL
	} else if c == "(" {
		token.Type = tLPAREN
	} else if c == ")" {
		token.Type = tRPAREN
	} else {
		os.Exit(1)
	}

	return token
}

var has_next_token = false
var next_token Token

func peekToken() Token {
	if has_next_token {
		return next_token
	}
	has_next_token = true
	next_token = lex()
	return next_token
}

func getToken() Token {
	if has_next_token {
		has_next_token = false
		return next_token
	}
	return lex()
}

func primaryExpression() {
	token := getToken()

	if token.Type == tINT {
		fmt.Printf("  sub $4, %%rsp\n  movl $%d, 0(%%rsp)\n", token.IntValue)
	} else if token.Type == tLPAREN {
		additiveExpression()
		if getToken().Type != tRPAREN {
			os.Exit(1)
		}
	} else {
		os.Exit(1)
	}
}

// *
func multiplicativeExpression() {
	primaryExpression()

	for {
		op := peekToken()
		if op.Type != tMUL {
			break
		}
		getToken()

		primaryExpression()

		fmt.Printf("  movl 0(%%rsp), %%edx\n  add $4, %%rsp\n")
		fmt.Printf("  movl 0(%%rsp), %%eax\n  add $4, %%rsp\n")

		if op.Type == tMUL {
			fmt.Printf("  imull %%edx\n")
		}

		fmt.Printf("  sub $4, %%rsp\n  movl %%eax, 0(%%rsp)\n")
	}
}

// +
func additiveExpression() {
	multiplicativeExpression()

	for {
		op := peekToken()
		if op.Type != tADD && op.Type != tSUB {
			break
		}
		getToken()

		multiplicativeExpression()

		fmt.Printf("  movl 0(%%rsp), %%edx\n  add $4, %%rsp\n")
		fmt.Printf("  movl 0(%%rsp), %%eax\n  add $4, %%rsp\n")

		if op.Type == tADD {
			fmt.Printf("  addl %%edx, %%eax\n")
		}

		if op.Type == tSUB {
			fmt.Printf("  subl %%edx, %%eax\n")
		}

		fmt.Printf("  sub $4, %%rsp\n  movl %%eax, 0(%%rsp)\n")
	}
}

func main() {
	arg = os.Args[1]
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  push %%rbp\n")
	fmt.Printf("  mov %%rsp, %%rbp\n")

	additiveExpression()

	fmt.Printf("  movl 0(%%rsp), %%eax\n  add $4, %%rsp\n")
	fmt.Printf("  pop %%rbp\n")
	fmt.Printf("  ret\n")
}
