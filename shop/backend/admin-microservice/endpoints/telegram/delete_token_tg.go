package telegram

import (
	"encoding/json"
	"net/http"

	auth "github.com/AndreyLebedev1998/auth-grpc"
)

func RemoveToken(w http.ResponseWriter, r *http.Request, clientAuth auth.AuthServiceClient) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	resp, err := clientAuth.RemoveTelegramToken(ctx, &auth.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp != nil {
		w.Header().Set("Content-Type", "application-json")
		json.NewEncoder(w).Encode(resp)
	}
}
