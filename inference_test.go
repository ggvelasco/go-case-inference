package main

import "testing"

func TestInferenceSimples(t *testing.T) {
	// 1. Define a policy (regra)
	policy := `digraph {
		start [result=""]
		aprovado [result="approved=true"]
		
		start -> aprovado [cond="age>=18"]
	}`

	// 2. Define o input (dados)
	input := map[string]interface{}{
		"age": 25,
	}

	// 3. Roda a inferência
	output, err := runInference(policy, input)

	// 4. Verifica se não deu erro
	if err != nil {
		t.Fatalf("Não esperava erro, mas recebi: %v", err)
	}

	// 5. Verifica se o output tem o age original
	if output["age"] != 25 {
		t.Errorf("Esperava age=25, mas recebi %v", output["age"])
	}

	// 6. Verifica se foi aprovado
	if output["approved"] != true {
		t.Errorf("Esperava approved=true, mas recebi %v", output["approved"])
	}
}

func TestInferenceMenorDeIdade(t *testing.T) {
	policy := `digraph {
		start [result=""]
		aprovado [result="approved=true"]
		negado [result="approved=false"]
		
		start -> aprovado [cond="age>=18"]
		start -> negado [cond="age<18"]
	}`

	// Pessoa MENOR de idade
	input := map[string]interface{}{
		"age": 15,
	}

	output, err := runInference(policy, input)

	if err != nil {
		t.Fatalf("Não esperava erro, mas recebi: %v", err)
	}

	// Deve ser NEGADO
	if output["approved"] != false {
		t.Errorf("Esperava approved=false, mas recebi %v", output["approved"])
	}
}

func TestInferenceComMultiplasCondicoes(t *testing.T) {
	policy := `digraph {
		start [result=""]
		prime [result="approved=true,segment=prime"]
		standard [result="approved=true,segment=standard"]
		negado [result="approved=false"]
		
		start -> prime [cond="age>=18 && score>700"]
		start -> standard [cond="age>=18 && score>=600 && score<=700"]
		start -> negado [cond="age<18"]
	}`

	input := map[string]interface{}{
		"age":   25,
		"score": 750,
	}

	output, err := runInference(policy, input)

	if err != nil {
		t.Fatalf("Não esperava erro, mas recebi: %v", err)
	}

	// Deve ser aprovado como PRIME
	if output["approved"] != true {
		t.Errorf("Esperava approved=true, mas recebi %v", output["approved"])
	}

	if output["segment"] != "prime" {
		t.Errorf("Esperava segment=prime, mas recebi %v", output["segment"])
	}
}
