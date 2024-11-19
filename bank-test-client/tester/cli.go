package tester

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ShowMenu() {
	fmt.Println("Welcome to the Bank API Tester!")
	fmt.Println("Choose an option:")
	fmt.Println("1. Create a new account")
	fmt.Println("2. List all accounts")
	fmt.Println("3. Get account details")
	fmt.Println("4. Make a deposit/withdrawal")
	fmt.Println("5. List account transactions")
	fmt.Println("6. Transfer funds")
	fmt.Println("7. Exit")
}

func InputString() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func InputFloat() float64 {
	for {
		input := InputString()
		value, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Print("Invalid input. Please enter a numeric value: ")
			continue
		}
		return value
	}
}

func GetMenuChoice() int {
	fmt.Println("Please, select the number associated to the method you'd like to test (1-7)")
	input := InputString()
	choice, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input, please enter a number.")
		return -1
	}
	return choice
}
