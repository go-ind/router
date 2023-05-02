package example

import (
	"fmt"
	"net/http"

	goind "github.com/go-ind/go-ind-router"
)

func main() {
	r := goind.SetupDefaultRouter()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "Hello, world!")
	})

	v1 := r.Group("/parent")
	v1.Get("/ase", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "Hello, world!")
	})

	v2 := r.Group("/parent2")
	v2.Get("/asep", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "Hello, worldssss!")
		goind.ResponseJSON(w, ctx, 200, true, "Succes", nil, nil)
	})

	r.Post("/", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		// id := req.URL.Query().Get("id")
		// fmt.Println("id =>", id)
		// fmt.Fprint(w, "POST succes cok")
		goind.ResponseJSON(w, ctx, 200, true, "Succes", nil, nil)
	})

	r.Put("/", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "PUT succes cok")
	})

	r.Delete("/", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "DElete succes cok")
	})
	r.Patch("/testt", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "Patch succes cok")
	})

	// v2 := r.Group("/parent2")
	// v2.Get("/asep", func(w http.ResponseWriter, req *http.Request) {
	// 	id := req.URL.Query().Get("id")
	// 	fmt.Println("id =>", id)
	// 	fmt.Fprint(w, "Hello, worldssss!")
	// })
	http.ListenAndServe(":8080", r)
}
