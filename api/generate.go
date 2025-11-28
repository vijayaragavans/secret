package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vijayaragavans/secret/config"
	"github.com/vijayaragavans/secret/internal"
)

func Generate(w http.ResponseWriter, r *http.Request) {

	type Input struct {
		Data string `json:"data"`
	}
	var (
		input                Input
		payloadBytes, secret []byte
		err                  error
		req                  *http.Request
		resp                 *http.Response
	)

	// Prepare data for Vault
	type VaultData struct {
		Data map[string]string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if secret, err = internal.Encrypt([]byte(config.EncryptKey), input.Data); err != nil {
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return
	}

	vaultPayload := VaultData{
		Data: map[string]string{
			"secret": string(secret),
		},
	}

	if payloadBytes, err = json.Marshal(vaultPayload); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Use provided key or default
	key := fmt.Sprintf("%d", time.Now().UnixNano())

	// Create request to Vault
	if req, err = http.NewRequest("POST", config.VAULT_URL+key, bytes.NewBuffer(payloadBytes)); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-Vault-Token", config.VAULT_TOKEN)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if resp, err = client.Do(req); err != nil || resp.StatusCode != 200 {
		http.Error(w, "Vault connection failed", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{
		"message": config.SUCCESS_MSG,
		"key":     key,
	})
}
