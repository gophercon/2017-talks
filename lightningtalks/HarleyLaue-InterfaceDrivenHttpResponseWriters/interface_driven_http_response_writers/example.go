var (
	jsonAccept    = response.Accept{"application/json", json.Marshal}
	xmlAccept     = response.Accept{"application/xml", xml.Marshal}
)

type GopherCon struct {
	Id   int64
	Name string
}

// Mixed handler that can return JSON & XML
func MixedHandler(w http.ResponseWriter, r *http.Request) {
	if err := response.Write(w,
		GopherCon{456, "GopherCon 2017"},
		response.Acceptable(r, jsonAccept, xmlAccept),
	); err != nil {
		log.Printf("Error: %+v", err)
	}
}

func main() {
	http.HandleFunc("/mixed", MixedHandler)
	http.ListenAndServe("127.0.0.1:8888", nil)
}
