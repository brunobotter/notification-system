
📍 Fase 1 — Conexão WebSocket funcionando
🧱 1️⃣ Implementar o handler /ws

👉 Route que:

faz upgrade para WebSocket

instancia o Client com seu Hub e Logger

chama Hub.Register(client)

dispara os loops de leitura/escrita

(Conforme corrigido no código anterior.)

➡️ Por que isso importa?
Sem isso o Hub nunca recebe clientes — então todo o restante do sistema não poderá entregar notificações reais.

📌 Marcos de sucesso:

Você conecta com o navegador/postman web socket

Pode ver logs de conexão no Hub

ReadPump e WritePump funcionam sem erro

🧱 2️⃣ Verificação visual simples

Abra dois navegadores apontando pra:

ws://localhost:PORT/ws


E envie mensagens entre eles via Hub.
Se funcionar → tem base real de realtime.

📍 Fase 2 — Entrega de notificações via REST → WS
🧱 3️⃣ Criar endpoint REST /notifications

Esse endpoint será onde o mundo “produz” notificações.

Exemplo:

POST /notifications

{
  "user_id": "123",
  "message": "Notificação de teste"
}


Nesse endpoint:

Valida dados

Cria uma notificação

Publica no Hub (ou, mais tarde, no Redis)

Responde o resultado

Esse é o passo onde frontend/backend vão conversar.

🧱 4️⃣ Enviar para um único usuário conectado

Hoje seu broadcast entrega para todos.

Você precisa:

Mapear userID → client(s)

Enviar somente para quem está conectado

Isso é crucial para um sistema de notificações.

📍 Fase 3 — Redis Pub/Sub

Até agora o Hub (em memória) entrega mensagens no processo.

Mas para escalar:

🧱 5️⃣ Configurar Redis e conectar seu serviço

Redis Pub/Sub para receber eventos de notificações

Seu serviço assina (SUBSCRIBE)

Sempre que recebe no canal, envia ao Hub

📌 Isso permite:

múltiplas instâncias da aplicação

desacoplar produtores e consumidores

escalabilidade horizontal

📍 Fase 4 — Persistência com Database
🧱 6️⃣ Criar tabela notifications no PostgreSQL

Campos sugeridos:

id (uuid)

user_id

message

created_at

delivered (boolean)

🧱 7️⃣ Criar repositório e métodos

SaveNotification

FindNotificationsByUser

MarkAsRead

📍 Fase 5 — Filtro por usuário e histórico
🧱 8️⃣ Implementar rota para buscar histórico
GET /notifications?user_id=123


Retorna:

notificações não lidas

todas as notificações (com paginação opcional)

📍 Fase 7 — Testes e qualidade
🧱 10️⃣ Escrever testes

Com as interfaces que você já tem, pode testar:

Hub isolado (sem WS)

Client com uma conexão fake

Endpoint REST

Redis integration

📍 Fase 8 — Produção & escalabilidade
🧱 11️⃣ Deploy com Docker + Docker Compose

Serviços que você precisará:

app Go

Redis

PostgreSQL

🧱 12️⃣ Monitoramento

logs

métricas (Prometheus/Grafana)

🧠 Roadmap resumido (ordem)
Ordem	Tarefa
1	Rota /ws funcionando
2	Teste básico de conexão WS
3	Endpoint REST /notifications
4	Envio para usuário específico
5	Redis Pub/Sub
6	Persistência PostgreSQL
7	Rota para histórico
9	Testes automatizados
10	Deploy em container
11	Monitoramento/Métricas



Isso forma uma série coesa.
