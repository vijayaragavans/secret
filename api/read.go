package api

import (
	"encoding/json"
	"net/http"

	"github.com/vijayaragavans/secret/config"
	"github.com/vijayaragavans/secret/internal"
)

func Read(w http.ResponseWriter, r *http.Request) {

	type VaultResponse struct {
		Data struct {
			Data struct {
				Secret string `json:"secret"`
			} `json:"data"`
		} `json:"data"`
	}
	var (
		req           *http.Request
		resp          *http.Response
		err           error
		client        = &http.Client{}
		responseBytes []byte
		vaultResp     VaultResponse
		output        string
		response      = map[string]interface{}{
			"data": output,
		}
	)

	if req, err = http.NewRequest("GET", config.VAULT_URL+"generated-secret", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-Vault-Token", config.VAULT_TOKEN)

	if resp, err = client.Do(req); err != nil || resp.StatusCode != 200 {
		http.Error(w, "Vault connection failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&vaultResp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if response["data"], err = internal.Decrypt([]byte(config.EncryptKey), vaultResp.Data.Data.Secret); err != nil {
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return
	}

	if responseBytes, err = json.Marshal(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
