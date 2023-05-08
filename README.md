# router
Go ind router


# Example

func main() {
	r := goind.SetupDefaultRouter()

	r.Post("/", func(w http.ResponseWriter, req *http.Request) {
        goind.ResponseJSON(w, ctx, 200, true, "Succes", nil, nil)
	})

	v1 := r.Group("/parent")
	v1.Get("/child", func(w http.ResponseWriter, req *http.Request) {
        goind.ResponseJSON(w, ctx, 200, true, "Succes", nil, nil)
	})

	http.ListenAndServe(":8080", r)
