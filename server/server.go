package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Cotacao representa a estrutura da cotação do dólar
type Cotacao struct {
	ID        uint   `gorm:"primaryKey"`
	Rate      string `gorm:"not null"`
	Timestamp time.Time
}

// Configura o banco de dados SQLite utilizando GORM
func setupDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("cotacoes.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Montando o AutoMigrate para criar a tabela
	db.AutoMigrate(&Cotacao{})
	return db
}

// Função para buscar a cotação do dólar da API externa
func getDollarRate(ctx context.Context) (string, error) {
	// configurando context Timeout para 200ms
	client := &http.Client{Timeout: 200 * time.Millisecond}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// convertendo o resultado da request para json usando map
	var result map[string]map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	// capturando a propriedade do json
	cotacao, ok := result["USDBRL"]
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}
	// retornando o value da key bid
	return cotacao["bid"], nil
}

// Função para salvar a cotação no banco de dados
func saveRate(ctx context.Context, db *gorm.DB, rate string) error {
	cotacao := Cotacao{Rate: rate, Timestamp: time.Now()}
	// configurando o context para 10ms
	ctxDB, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	// salvando utilizando o context e retornando erro caso ocorra
	return db.WithContext(ctxDB).Create(&cotacao).Error
}

// Handler para a rota /cotacao
func cotacaoHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// criando um context com timeout de 200ms
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()
	// buscando a cotação do dólar
	rate, err := getDollarRate(ctx)
	if err != nil {
		log.Printf("Error fetching dollar rate: %v", err)
		http.Error(w, "failed to fetch dollar rate", http.StatusInternalServerError)
		return
	}
	// salvando a cotação no db
	if err := saveRate(r.Context(), db, rate); err != nil {
		log.Printf("Error saving rate to database: %v", err)
		http.Error(w, "failed to save rate to database", http.StatusInternalServerError)
		return
	}
	// retornando para o client um json com o bid ( cotação do dólar em reais )
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"bid": rate})
}

func main() {
	// criando ou startando db
	db := setupDatabase()
	// criando a rota
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		cotacaoHandler(w, r, db)
	})

	log.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
