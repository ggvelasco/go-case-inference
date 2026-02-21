package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Receber do cliente
type InferRequest struct {
	PolicyDOT string                 `json:"policy_dot"`
	Input     map[string]interface{} `json:"input"`
}

// Retornar ao cliente
type InferResponse struct {
	Output map[string]interface{} `json:"output"`
}

func inferHandler(w http.ResponseWriter, r *http.Request) {
	//Fazer com que aceite apenas o POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	//Ler o JSON do r.Body
	var req InferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	//Rodar a inferencia
	output, err := runInference(req.PolicyDOT, req.Input)
	if err != nil {
		http.Error(w, "Erro na inferência: "+err.Error(), http.StatusBadRequest)
		return
	}

	//Retornar resultado
	resp := InferResponse{Output: output}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	//Registrar endpoint
	http.HandleFunc("/infer", inferHandler)

	// Iniciar o servidor
	port := ":8080"
	fmt.Printf("Servidor rodando em http://localhost%s\n", port)
	fmt.Println("Endpoint: POST /infer")
	fmt.Println()

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Erro ao iniciar o servidor : %v\n", err)
	}
}
