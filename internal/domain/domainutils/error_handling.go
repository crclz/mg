package domainutils

import (
	"fmt"
)

func ShowErrorToUser(err error) {
	if err != nil {
		fmt.Printf("Error detail: %+v\n", err)
		fmt.Printf("\nError: %v\n", err)
	}
}
