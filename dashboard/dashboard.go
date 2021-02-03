package dashboard

import (
	"fmt"
	"strconv"
)

func Run(port int) {
	addr := "127.0.0.1:" + strconv.Itoa(port)
	fmt.Println("Web UI opening...")
	err := Open("http://" + addr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Starting server at %s...\n", addr)
	Server(addr)

}
