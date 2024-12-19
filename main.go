package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	errorhandler "github.com/vishrudh-raj-rs-14/lox/errorHandler"
	"github.com/vishrudh-raj-rs-14/lox/operations"
	"github.com/vishrudh-raj-rs-14/lox/parser"
	"github.com/vishrudh-raj-rs-14/lox/resolver"
	"github.com/vishrudh-raj-rs-14/lox/scanner"
)

func main() {
	start := time.Now() // Start timing
	cmdArgs := os.Args[1:]
	if len(cmdArgs) > 1 {
		fmt.Println("Usage: Glox [script]")
	} else if len(cmdArgs) == 1 {
		runFile(cmdArgs[0])
	} else {
		runPrompt()
	}
	duration := time.Since(start) // Calculate duration
	fmt.Printf("Execution time: %v\n", duration)
}

func runFile(path string){
	data, err := os.ReadFile(path);
	interpreter := operations.CreateInterpreter();
	if(err != nil){
		fmt.Println("Error: Unable to readfile");
	}
	run(string(data), interpreter);
}

func runPrompt(){
	scanner := bufio.NewScanner(os.Stdin);
	interpreter := operations.CreateInterpreter();
	for{
		fmt.Printf("> ")
		if (!scanner.Scan()){
			return;
		}

		code := scanner.Text();
		if(strings.TrimSpace(strings.ToLower(code))=="exit"){
			fmt.Println("Exiting...");
			return
		}
		run(code, interpreter);
	}
}



func run(code string, interpreter operations.Interpreter){

	scanner := scanner.CreateScanner(code);
	tokens := scanner.ScanTokens();
	gloxParser := parser.CreateParser(tokens);
	ast := gloxParser.Parse();
	if(errorhandler.HadError){
		errorhandler.HadError = false;
		return
	}
	resolver := resolver.CreateResolver(interpreter);
	resolver.ResolveStatements(ast);
	if(errorhandler.HadError){
		errorhandler.HadError = false;
		return
	}
	interpreter.Interpret(ast);
	errorhandler.HadError = false;
}