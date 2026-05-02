package telegram

import (
	"admin-microservice/models"
	"encoding/json"
	"net/http"

	auth "github.com/AndreyLebedev1998/auth-grpc"
)

func AddTokenTg(w http.ResponseWriter, r *http.Request, clientAuth auth.AuthServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var token models.NewTelegramToken
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := clientAuth.AddTelegramToken(ctx, &auth.NewTelegramToken{
		Token: token.Token,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res != nil {
		w.Header().Set("Content-Type", "application-json")
		json.NewEncoder(w).Encode(res)
	} else {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}
