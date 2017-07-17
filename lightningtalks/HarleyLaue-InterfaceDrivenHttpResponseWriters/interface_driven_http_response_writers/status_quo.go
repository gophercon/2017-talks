import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}
