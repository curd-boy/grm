package logo

import "gopkg.in/ffmt.v1"

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
	ffmt.Printf(logo, ver)
}
