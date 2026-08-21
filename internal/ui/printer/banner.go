package printer

import "fmt"

// PrintBanner prints the CLI banner when running `mdev` without arguments.
func PrintBanner() {
	fmt.Print(`
███╗   ███╗██████╗ ███████╗██╗   ██╗
████╗ ████║██╔══██╗██╔════╝██║   ██║
██╔████╔██║██║  ██║█████╗  ██║   ██║
██║╚██╔╝██║██║  ██║██╔══╝  ╚██╗ ██╔╝
██║ ╚═╝ ██║██████╔╝███████╗ ╚████╔╝ 
╚═╝     ╚═╝╚═════╝ ╚══════╝  ╚═══╝  

mdev — macOS Development Environment Manager
`)
}
