package exception

import "fmt"

func PanicIfNeeded(err interface{}) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
