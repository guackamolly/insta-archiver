package main

import (
	"fmt"

	"github.com/guackamolly/insta-archiver/internal/data/client/http"
	"github.com/guackamolly/insta-archiver/internal/data/user"
)

func main() {
	r := user.ViewIGStoryUserRepository{
		Client: http.Native(),
	}
	res, err := r.Stories("cristiano")

	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", res)
	}
}
