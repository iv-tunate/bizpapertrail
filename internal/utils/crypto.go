package utils

import (
	"crypto/rand"
	"log/slog"
	"math/big"
)

func GenerateRandomToken(length int) (string, error) {
	const pool = "AB01234PQRC5VWXYZ6KLMNO789FGHIJSTDEU"
	token := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))

		if err != nil{
			slog.Info("[ERROR] GenerateRandomToken: An error occured while generating a random token", "Error Details", err)
			return "", err
		}

		token[i] = pool[num.Int64()]
	}
	
	return  string(token), nil
}