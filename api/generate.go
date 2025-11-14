package api

import (
	"log"
	"net/http"

	"github.com/vijayaragavans/secret/config"
	"github.com/vijayaragavans/secret/internal"
)

func Generate(w http.ResponseWriter, r *http.Request) {

	var (
		password []byte
		err      error
	)

	if password, err = internal.Encrypt([]byte(config.EncryptKey), "akjdhaksdjhaksd"); err != nil {
		log.Println("Error:", err)
	}

	w.WriteHeader(200)
	w.Write(password)

}
