package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Response representa a estrutura da resposta do servidor
type Response struct {
	Bid string `json:"bid"`
}

// Função para obter a cotação do servidor
func getCotacao(ctx context.Context) (string, error) {
	// Criando um cliente HTTP
	client := &http.Client{Timeout: 300 * time.Millisecond}
	// Criando uma requisição HTTP com o context
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// estruturando o json da request na struct
	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Bid, nil
}

func main() {
	// Criando o context de 300ms
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	// Obtendo a cotação do dólar
	bid, err := getCotacao(ctx)
	if err != nil {
		log.Printf("Error fetching cotacao: %v", err)
		return
	}
	// Criando o txt de cotação
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return
	}
	defer file.Close()
	// Escrevendo a cotação no txt
	content := fmt.Sprintf("Dólar: %s", bid)
	if _, err := file.WriteString(content); err != nil {
		log.Printf("Error writing to file: %v", err)
	}
}
