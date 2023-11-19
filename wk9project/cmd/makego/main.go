package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jeremycruzz/msds301-wk8/pkg/chatgpt"
)

func main() {
	var promptBuilder strings.Builder
	defaultPrompt := "I need a program that analyzes all four sets of the AnscombeQuartet dataset using linear regression and prints 'Set I: m= b=' for each of the four sets."

	// get from flags
	apiKey := flag.String("apikey", "", "API key for chatgpt")
	name := flag.String("name", "newProject", "Name for go module")
	prompt := flag.String("prompt", defaultPrompt, "prompt for chat gpt program")

	flag.Parse()

	if *apiKey == "" {
		log.Fatal("API key is required. Start with -apikey flag.")
	}

	// string to only get code
	promptBuilder.WriteString("I am going to ask you to write a go program for me all contained within a main.go file. Only respond with the contents of this main.go file in raw text and nothing else. I'm going to paste your response directly into a go file. Do not include anything before or after the code including comments explaining the code. ")
	promptBuilder.WriteString(*prompt)

	newDir := fmt.Sprintf("../%v", *name)
	chatgpt := chatgpt.New(*apiKey)

	// create new directory
	fmt.Printf("Creating directory: %v...\n", newDir)
	err := os.Mkdir(newDir, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// cd new directory
	err = os.Chdir(newDir)
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}

	// init go mod
	fmt.Printf("Creating go module: %v...\n", *name)
	cmd := exec.Command("go", "mod", "init", *name)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error initializing go module:", err)
		return
	}

	//ask chat gpt for code
	fmt.Printf("Asking chatgpt: \n%v", promptBuilder.String())
	code, err := chatgpt.AskCustom(promptBuilder.String())
	if err != nil {
		fmt.Println("Error with chatgpt", err)
		return
	}

	//remove ```go and ```
	code = code[6 : len(code)-3]

	// write code
	fmt.Println("Writing to main.go...")
	err = os.WriteFile("main.go", []byte(code), 0644)
	if err != nil {
		fmt.Println("Error writing Go code to file:", err)
		return
	}

	// tidy module
	fmt.Println("Tidying dependencies...")
	cmd = exec.Command("go", "mod", "tidy")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running go mod tidy:", err)
		return
	}

	// build
	cmd = exec.Command("go", "build")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error building the project:", err)
		return
	}

	fmt.Println("Project setup and build complete.")
}
