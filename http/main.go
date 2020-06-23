package main

import (
	"fmt"
	"net/http"
)

type H struct {
	name string
}

var hy H

func main() {
	//go sub()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover 1/0 error")
		}
	}()
	http.Handle("/", &hy)
	http.HandleFunc("/info", IndexHandle)

	http.ListenAndServe(":8080", nil)
}

func (h *H) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.name = "120"
	w.Write([]byte(fmt.Sprint(h.name)))
}

func IndexHandle(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(hy.name + "\n hello"))
	i := 100
	for j := 0; j < 3; j++ {
		go sub(i, j)
	}

}

func sub(i, j int) {
	defer func() {
		if recover() != nil {
			fmt.Println("recover err")
		}
	}()
	fmt.Println(i / j)
}
