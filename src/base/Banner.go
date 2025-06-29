package base

import "fmt"

const Version = "4.12.19-stable"

func ShowBanner() {
	banner := `
 __      __  ______   __         ______    ______   __       __  ______  
/\ \    /\ \/\  ___\ /\ \       /\  ___\  /\  __ \ /\ \     /\ \/\  ___\ 
\ \ \___\ \ \ \  __\ \ \ \____  \ \ \____ \ \ \/\ \\ \ \____\ \ \ \  __\ 
 \ \_______\ \ \_____\\ \_____\  \ \_____\ \ \_____\\ \_____\\ \_\ \_____\
  \/_______/  \/_____/ \/_____/   \/_____/  \/_____/ \/_____/ \/_/\/_____/
																		                                 
`
	fmt.Println(banner)
	// Add welcome banner
	fmt.Println()
	fmt.Println("+------------------------------------------+")
	fmt.Println("|          Welcome to ContainDB CLI        |")
	fmt.Printf("|               Version: %-18s|\n", Version)
	fmt.Println("+------------------------------------------+")
	fmt.Println("|  A simple CLI to manage DB containers    |")
	fmt.Println("|                                          |")
	fmt.Println("|           Made by Ankan Saha             |")
	fmt.Println("+------------------------------------------+")
	fmt.Println()
}
