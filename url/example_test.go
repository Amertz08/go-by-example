package url_test

import (
	"fmt"
	"log"

	"github.com/Amertz08/go-by-example/url"
)

func ExampleParse() {
	uri, err := url.Parse("https://github.com/inancgumus")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(uri)
	// Output:
	// https://github.com/inancgumus
}
