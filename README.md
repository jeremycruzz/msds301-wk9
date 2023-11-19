# msds301-wk8

### Setup (windows, gitbash)
- Clone repo with `git clone git@github.com:jeremycruzz/msds301-wk9.git`
- Get in project directory `cd wk9project`
    - you can build other projects as well
- Run `go mod tidy`

### Building executable
- Run `go build ./cmd/makego`

### Testing
- Only manual testing is done. Run each program manually

### Running Go executable
- Run `./makego.exe -apikey {API_KEY} -name {name} -prompt {program description}` 
    - `-name` and `prompt` flags are optional
    - `-name` is the name of project
    - `prompt` is the program description
    - I'll put the apikey as a comment in the assignment

### Background / Conclusion

For this assignment I made a program that would have chatgpt create a program for me. I started out with a simple `hello world` program just to see if I could get go to create and build a go program. Once that was complete I added my chatgpt package from wk8 and started building some prompts. These were the first few prompts I wrote to get the anscombe quartet program running.
<details>
<summary> "I am going to ask you to write a go program for me all contained within a main.go file. Only respond with the contents of this main.go file that can easily be pasted into a main.go file. I need a program that analyzes the AnscombeQuartet dataset using linear regression and prints 'Set I: m= b=' for the coefficients." </summary>

```
Below is the Go code that you can paste into a `main.go` file. This program imports the `gonum.org/v1/gonum/stat` package for linear regression, so you will need to install the Gonum package if you haven't done so already. You can install it by running `go get -u gonum.org/v1/gonum/...` in your terminal.

```go
package main

import (
	"fmt"
	"gonum.org/v1/gonum/stat"
)

func main() {
	type dataset struct {
		x []float64
		y []float64
	}

	anscombeQuartet := []dataset{
		{
			x: []float64{10.0, 8.0, 13.0, 9.0, 11.0, 14.0, 6.0, 4.0, 12.0, 7.0, 5.0},
			y: []float64{8.04, 6.95, 7.58, 8.81, 8.33, 9.96, 7.24, 4.26, 10.84, 4.82, 5.68},
		},
		// The rest of the quartet data would go here if you wanted to analyze sets II, III, and IV as well.
	}

	for i, set := range anscombeQuartet {
		alpha, beta := stat.LinearRegression(set.x, set.y, nil, false)
		fmt.Printf("Set I: m=%.2f b=%.2f\n", beta, alpha)
		if i == 0 { // We only calculate Set I as per your request
			break
		}
	}
}
```

Remember to install the Gonum package before running the program and replace the dataset values if needed. The program calculates the coefficients for the first dataset of the Anscombe quartet and prints them out in the required format. The rest of the datasets are omitted for brevity, but the structure is in place for including those as well.
```
</details>

<details>

<summary> "I am going to ask you to write a go program for me all contained within a main.go file. Only respond with the contents of this main.go file and nothing else. I need a program that analyzes all four sets of the AnscombeQuartet dataset using linear regression and prints 'Set I: m= b=' for each of the four sets." </summary>

\```go
package main

import (
	"fmt"
	"gonum.org/v1/gonum/stat"
)

func main() {
	// Anscombe's quartet data
	type Dataset struct {
		x []float64
		y []float64
	}

	quartet := []Dataset{
		{
			x: []float64{10.0, 8.0, 13.0, 9.0, 11.0, 14.0, 6.0, 4.0, 12.0, 7.0, 5.0},
			y: []float64{8.04, 6.95, 7.58, 8.81, 8.33, 9.96, 7.24, 4.26, 10.84, 4.82, 5.68},
		},
		{
			x: []float64{10.0, 8.0, 13.0, 9.0, 11.0, 14.0, 6.0, 4.0, 12.0, 7.0, 5.0},
			y: []float64{9.14, 8.14, 8.74, 8.77, 9.26, 8.10, 6.13, 3.10, 9.13, 7.26, 4.74},
		},
		{
			x: []float64{10.0, 8.0, 13.0, 9.0, 11.0, 14.0, 6.0, 4.0, 12.0, 7.0, 5.0},
			y: []float64{7.46, 6.77, 12.74, 7.11, 7.81, 8.84, 6.08, 5.39, 8.15, 6.42, 5.73},
		},
		{
			x: []float64{8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 19.0, 8.0, 8.0, 8.0},
			y: []float64{6.58, 5.76, 7.71, 8.84, 8.47, 7.04, 5.25, 12.50, 5.56, 7.91, 6.89},
		},
	}

	// Analyze each dataset with linear regression and print the results
	for i, set := range quartet {
		alpha, beta := stat.LinearRegression(set.x, set.y, nil, false)
		fmt.Printf("Set %s: m=%.2f b=%.2f\n", string('I'+rune(i)), beta, alpha)
	}
}
\```

</details>


The second prompt was easy to work with as I only had to get rid of the first and last line that chatgpt used to format the code. Once I ran the program to make sure it matched the output of the program that I wrote during week 2 which it did. After that worked I had it create the guesser program, blackjack, poker, poker2, then the crawler with the following prompts:

**Guesser** : "I want my program to pick a number between 1 and 100 and has the user guess. If the number is wrong print whether or not the number was higher or lower. Keep the user guessing until they make the correct guess."

**Blackjack** : "Program a full blackjack game where I can bet against the dealer split and double."

**Poker** : "Program a full five card draw poker game where I can bet."

**Poker2** : "Program a full five card draw poker game where I can bet. Use the suit characters and have proper payouts."

**Crawler** : "Write a web crawler that starts on a random wikipedia page and crawls at a depth of 2 using concurrency. I want the program to write to a json file for the results."

The guesser and poker2 program worked immediately without me needing to make any adjustments. Blackjack and poker added extra comments at the end of the program that didn't allow it to build. After removing the comments the programs built and ran correctly. Crawler did NOT work but created a good starting point. 

While the programs worked, there were a few problems with poker and poker 2. The first iteration did not payout any money even if the user won. The second iteration payed out money on every hand including a high card hand. With additional details in the prompt I'm sure that the AI could generate code to properly pay the user out but a developer could go in and easily make that change. 

As of right now I don't think LLMs like chatgpt could replace developers. However, I think chatgpt could be used as a tool that our developers can use to further enhance their productivity. LLM's could work as an alternative to google and figure out whats wrong with their exact code. While non developers are able to generate code, a developer could generate better code with their understanding of code. 

