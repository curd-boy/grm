package logo

import "fmt"

var logo = `
 _____
|  __ \ 
| |  \/_ __ _ __ ___
| | __| '__| '_ ' _ \
| |_\ \ |  | | | | | |
 \____/_|  |_| |_| |_|
                       %v
`

func PrintLogo(ver interface{}) {
	fmt.Printf(logo, ver)
}
