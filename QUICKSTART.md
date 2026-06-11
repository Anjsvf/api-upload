# Quick Start Guide

Guia rápido para começar a desenvolver com a API Upload.

## 🚀 Começar em 5 minutos

### Passo 1: Clonar e Instalar
```bash
git clone https://github.com/Anjsvf/api-upload.git
cd api-upload
go mod download
```

### Passo 2: Configurar Ambiente
Criar `.env` na raiz:
```env
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=vicio_api
PORT=8080
YOUTUBE_API_KEY=  # opcional
```

### Passo 3: Executar MongoDB Localmente

**Com Docker:**
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

**Ou instalado localmente:**
```bash
mongod
```

### Passo 4: Rodar a API
```bash
go run main.go
```

Saída esperada:
```
Servidor rodando na alturas na porta 8080
```

### Passo 5: Testar
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

✅ **Pronto! API rodando localmente**

---

## 📝 Primeiros Testes

### 1. Verificar Tags
```bash
curl http://localhost:8080/api/v1/tags
```

Resposta esperada:
```json
[
  {
    "id": "...",
    "name": "Jogos",
    "slug": "jogos",
    "emoji": "🎮",
    "created_at": "2024-..."
  },
  ...
]
```

### 2. Obter Feed Vazio
```bash
curl http://localhost:8080/api/v1/feed
```

Resposta: `[]` (array vazio no início)

### 3. Criar Primeiro Post
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Primeiro Post",
    "caption": "Testando a API",
    "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
    "type": "youtube",
    "tag_id": "jogos"
  }'
```

✅ Se receber status 201 com o post, funcionou!

### 4. Ver o Feed Atualizado
```bash
curl http://localhost:8080/api/v1/feed
```

Agora deve retornar o post criado.

---

## 🔧 Desenvolvendo Localmente

### Estrutura de um Novo Endpoint

1. **Criar handler** em `handlers/`
```go
func MyNewEndpoint(c *gin.Context) {
    // Lógica
    c.JSON(200, gin.H{"message": "ok"})
}
```

2. **Registrar rota** em `routes/routes.go`
```go
v1.GET("/my-endpoint", handlers.MyNewEndpoint)
```

3. **Testar localmente**
```bash
go run main.go
curl http://localhost:8080/api/v1/my-endpoint
```

### Adionar Validação

Exemplo: Adicionar campo novo no Post

1. **Modificar modelo** em `models/post.go`
```go
type Post struct {
    // ... campos existentes
    NewField string `json:"new_field" bson:"new_field"`
}
```

2. **Adicionar validação** em `validators/`
```go
func ValidateNewField(field string) error {
    if len(field) < 3 {
        return fmt.Errorf("new_field deve ter no mínimo 3 caracteres")
    }
    return nil
}
```

3. **Usar no handler**
```go
if err := validators.ValidateNewField(req.NewField); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

---

## 🛠️ Ferramentas Úteis

### Postman
1. Importar coleção (exemplo em `docs/postman_collection.json`)
2. Configurar variável `base_url`
3. Executar requisições

### VS Code
Extensões recomendadas:
- Go (golang.go)
- REST Client
- Thunder Client
- MongoDB for VS Code

### CLI Tools
```bash
# Instalar Go tools
go install github.com/cosmtrek/air@latest   # Hot reload
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Hot reload durante desenvolvimento
air
```

---

## 🐛 Debugging

### Ver logs do servidor
```bash
go run main.go  # Logs aparecem no terminal
```

### Verificar conexão MongoDB
```bash
# No terminal do MongoDB
mongo
> show dbs
> use vicio_api
> db.posts.find()
```

### Testar Rate Limit
```bash
# Executar 5 vezes (limite é 4):
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/v1/posts \
    -H "Content-Type: application/json" \
    -d '{
      "title": "Post '$i'",
      "caption": "Teste",
      "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
      "type": "youtube",
      "tag_id": "jogos"
    }'
  echo ""
done
```

A 5ª deve retornar erro 429.

---

## 📚 Referências Rápidas

| Comando | Descrição |
|---------|-----------|
| `go run main.go` | Rodar API |
| `go build` | Compilar para executável |
| `go mod tidy` | Limpar dependências não usadas |
| `go fmt ./...` | Formatar código |
| `go vet ./...` | Verificar erros comuns |

---

## 🚢 Deploy Rápido

### Heroku
```bash
# Criar app
heroku create api-upload

# Configurar variáveis
heroku config:set MONGO_URI=...
heroku config:set YOUTUBE_API_KEY=...

# Deploy
git push heroku main
```

### Railway
```bash
# Login
railway login

# Deploy
railway up
```

---

## ❓ FAQ

### P: MongoDB local não conecta
**R**: Verifique se MongoDB está rodando:
```bash
# macOS/Linux
mongod

# Windows
# Iniciar MongoDB Service via Services
```

### P: "Module not found" error
**R**: Rode `go mod download` e `go mod tidy`

### P: Porta 8080 já está em uso
**R**: Mude no `.env`: `PORT=8081`

### P: YOUTUBE_API_KEY é obrigatório?
**R**: Não. Se vazio, links YouTube não serão validados (apenas links WhatsApp funcionarão)

### P: Como resetar o banco?
**R**: 
```bash
# No MongoDB shell
use vicio_api
db.dropDatabase()
# Reinicie a API para recriar tags
```

---

## 🎯 Próximos Passos

Depois de rodar localmente:

1. ✅ Entenda a [Arquitetura](ARCHITECTURE.md)
2. ✅ Leia a [Documentação Completa](README.md)
3. ✅ Explore o código em `handlers/`, `models/`, `database/`
4. ✅ Implemente sua primeira feature
5. ✅ Escreva testes
6. ✅ Faça commit e push
7. ✅ Abra um Pull Request

---

## 📞 Precisa de Ajuda?

- 📖 Veja [README.md](README.md)
- 🏗️ Veja [ARCHITECTURE.md](ARCHITECTURE.md)
- 🐛 Abra uma issue no GitHub
- 💬 Pergunte no Discord (se houver comunidade)

---

**Quick Start - v1.0**
