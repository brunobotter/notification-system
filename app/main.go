package main

import (
	"github.com/brunobotter/notification-system/main/app"
	"github.com/brunobotter/notification-system/main/providers"
)

func main() {
	app.NewApplication(providers.List()).Bootstrap()
}

/*Seu método bootstrapProvider() já chama todos os métodos Boot dos providers, então está correto.
Explique que Providers podem ter métodos de inicialização (Boot) para rodar processos background.
Mostre que o Hub precisa rodar em background para processar conexões.
Mostre como a injeção de dependências facilita a testabilidade e desacoplamento.*/
