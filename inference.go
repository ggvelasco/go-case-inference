package main

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
)

func parseValue(s string) interface{} {
	// Tenta bool
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}

	// Tenta número inteiro
	var num int
	if _, err := fmt.Sscanf(s, "%d", &num); err == nil {
		return num
	}

	// Se não for nada disso, é string
	return s
}

func evaluateCondition(cond string, variables map[string]interface{}) bool {
	if cond == "" {
		return true
	}

	expression, err := govaluate.NewEvaluableExpression(cond)
	if err != nil {
		// fmt.Println("Erro ao criar expressão:", err)
		return false
	}

	//Avaliar
	result, err := expression.Evaluate(variables)
	if err != nil {
		// fmt.Println("Erro ao avaliar:", err)
		return false
	}

	//Converter para bool e retornar o result
	boolResult, ok := result.(bool)
	if !ok {
		// fmt.Println("Resultado não é booleano")
		return false
	}
	return boolResult
}

// Aplica um resultado no output
// Ex: result = "approved=true,segment=prime"
// Adiciona: output["approved"] = true, output["segment"] = "prime"
func applyResult(result string, output map[string]interface{}) {
	if result == "" {
		return // nada pra fazer
	}

	// Quebra por vírgula: "approved=true,segment=prime" → ["approved=true", "segment=prime"]
	pairs := strings.Split(result, ",")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair) // remove espaços extras

		if pair == "" {
			continue // pula vazios
		}

		// Quebra por =: "approved=true" → ["approved", "true"]
		parts := strings.SplitN(pair, "=", 2)

		if len(parts) != 2 {
			continue // se não tem =, pula
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Converte o valor e adiciona no output
		output[key] = parseValue(value)

		// fmt.Println("  →", key, "=", output[key])
	}
}

func runInference(policyDOT string, input map[string]interface{}) (map[string]interface{}, error) {
	// 1. Parsear o DOT (NOVO jeito)
	graph, err := parseDOT(policyDOT)
	if err != nil {
		return nil, err
	}

	// 2. Criar output inicial (cópia do input)
	output := make(map[string]interface{})
	for chave, valor := range input {
		output[chave] = valor
	}

	// 3. Começar do nó "start"
	currentNode := "start"

	// 4. Loop de navegação
	for {
		// fmt.Println("Estou no nó:", currentNode)

		if node, exists := graph.Nodes[currentNode]; exists {
			if node.Result != "" {
				// fmt.Println("Aplicando resultado:", node.Result)
				applyResult(node.Result, output)
			}
		}

		// Buscar setas que saem do nó atual
		var minhasSetas []*Edge

		for _, seta := range graph.Edges {
			if seta.From == currentNode {
				minhasSetas = append(minhasSetas, seta)
			}
		}

		// fmt.Println("Encontrei", len(minhasSetas), "setas")

		if len(minhasSetas) == 0 {
			// fmt.Println("Não existem mais setas")
			break
		}

		proximoNo := ""

		for _, seta := range minhasSetas {
			condicao := seta.Cond
			// fmt.Println("Testando condição:", condicao)

			if evaluateCondition(condicao, output) {
				// fmt.Println("Condição verdadeira")
				proximoNo = seta.To
				break
			} else {
				// fmt.Println("Condição falsa, próxima...")
			}
		}

		if proximoNo == "" {
			// fmt.Println("Nenhuma condição foi verdadeira, encerrado.")
			break
		}
		// fmt.Println("Indo para o nó:", proximoNo)
		currentNode = proximoNo

	}

	return output, nil
}
