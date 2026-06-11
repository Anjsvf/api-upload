# Development Guide

Guia para desenvolvedores que querem contribuir com o projeto.

---

## 🛠️ Setup de Desenvolvimento

### Pré-requisitos
- Go 1.25.0+
- MongoDB 5.0+
- Git
- Visual Studio Code (recomendado)

### Extensões VS Code Recomendadas
```json
{
  "recommendations": [
    "golang.go",
    "ms-vscode.makefile-tools",
    "mongodb.mongodb-vscode",
    "humao.rest-client"
  ]
}
```

### Clonar e Configurar
```bash
git clone https://github.com/Anjsvf/api-upload.git
cd api-upload
go mod download
cp .env.example .env  # ou criar .env manualmente
```

---

## 📊 Estrutura de Código

### Convenções

#### Nomeação
- **Pacotes**: `lowercase` (ex: `handlers`, `database`)
- **Funções Exportadas**: `PascalCase` (ex: `CreatePost`)
- **Funções Privadas**: `camelCase` (ex: `getClientIP`)
- **Constantes**: `UPPER_CASE` (ex: `MaxPostsPerDay`)
- **Variáveis**: `camelCase` (ex: `clientIP`)

#### Estruturas
```go
type Post struct {
    // Public fields first
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Title     string             `json:"title" bson:"title"`
    
    // Private fields last
    ip        string             `json:"-" bson:"ip"`
}
```

#### Tags de Struct
- `json:"..."` - Serialização JSON
- `bson:"..."` - Serialização MongoDB
- `binding:"..."` - Validação Gin
- `json:"-"` - Omitir em JSON
- `bson:"-"` - Omitir em BSON

### Organização de Pacotes

```
handlers/       - Controladores HTTP
├── post_handler.go
├── tag_handler.go
└── common.go    (funções compartilhadas)

models/         - Modelos de dados
├── post.go
├── tag.go
└── dto.go       (Data Transfer Objects)

validators/     - Validadores
├── link_validator.go
└── content_validator.go

database/       - Camada de dados
├── mongodb.go
└── operations.go (CRUD operations)

middleware/     - Middlewares
├── ratelimit.go
└── cors.go

routes/         - Roteamento
└── routes.go

config/         - Configuração
└── config.go

filters/        - Filtros
└── word_filter.go

seed/           - Dados iniciais
└── seed.go
```

---

## ✏️ Adicionar Novo Endpoint

### Passo 1: Criar Model (se necessário)
**Arquivo**: `models/new_model.go`
```go
package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type NewModel struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name      string             `json:"name" bson:"name"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type CreateNewModelRequest struct {
    Name string `json:"name" binding:"required,min=3,max=100"`
}
```

### Passo 2: Criar Validador (se necessário)
**Arquivo**: `validators/new_validator.go`
```go
package validators

import "fmt"

func ValidateNewModel(name string) error {
    if len(name) < 3 {
        return fmt.Errorf("name deve ter no mínimo 3 caracteres")
    }
    return nil
}
```

### Passo 3: Criar Handler
**Arquivo**: `handlers/new_handler.go`
```go
package handlers

import (
    "context"
    "net/http"
    "time"

    "github.com/Anjsvf/api-upload/database"
    "github.com/Anjsvf/api-upload/models"
    "github.com/Anjsvf/api-upload/validators"
    
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateNewModel(c *gin.Context) {
    var req models.CreateNewModelRequest
    
    // 1. Parse e valida JSON
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Dados inválidos",
            "details": err.Error(),
        })
        return
    }
    
    // 2. Validações de negócio
    if err := validators.ValidateNewModel(req.Name); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
        return
    }
    
    // 3. Operação no banco
    collection := database.GetCollection("new_models")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    model := models.NewModel{
        ID:        primitive.NewObjectID(),
        Name:      req.Name,
        CreatedAt: time.Now(),
    }
    
    _, err := collection.InsertOne(ctx, model)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar"})
        return
    }
    
    // 4. Retorna resultado
    c.JSON(http.StatusCreated, model)
}

func GetNewModels(c *gin.Context) {
    collection := database.GetCollection("new_models")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar"})
        return
    }
    defer cursor.Close(ctx)
    
    var models []models.NewModel
    if err := cursor.All(ctx, &models); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar"})
        return
    }
    
    if models == nil {
        models = []models.NewModel{}  // Retornar array vazio em vez de null
    }
    
    c.JSON(http.StatusOK, models)
}
```

### Passo 4: Registrar Rota
**Arquivo**: `routes/routes.go`
```go
func Setup() *gin.Engine {
    r := gin.Default()
    
    // ... middlewares existentes ...
    
    v1 := r.Group("/api/v1")
    {
        // ... rotas existentes ...
        
        // Novas rotas
        v1.GET("/new-models", handlers.GetNewModels)
        v1.POST("/new-models", handlers.CreateNewModel)
    }
    
    return r
}
```

### Passo 5: Testar
```bash
# 1. Rodar aplicação
go run main.go

# 2. Testar novo endpoint
curl http://localhost:8080/api/v1/new-models

# 3. Criar novo modelo
curl -X POST http://localhost:8080/api/v1/new-models \
  -H "Content-Type: application/json" \
  -d '{"name": "Exemplo"}'
```

---

## 🧪 Testes

### Estrutura de Testes

Criar arquivo `handlers/post_handler_test.go`:

```go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
    // Arrange
    req := models.CreatePostRequest{
        Title:   "Test Post",
        Caption: "Test Caption",
        Link:    "https://www.youtube.com/watch?v=test",
        Type:    "youtube",
        TagID:   "test-tag",
    }
    
    // Act
    // ... fazer operação ...
    
    // Assert
    assert.NotNil(t, result)
    assert.Equal(t, "Test Post", result.Title)
}
```

### Executar Testes
```bash
# Todos os testes
go test ./...

# Com coverage
go test -cover ./...

# Com detalhes
go test -v ./...

# Teste específico
go test -run TestCreatePost ./handlers
```

---

## 🚀 Build & Deployment

### Build Local
```bash
# Build para o SO atual
go build -o api-upload

# Build para Linux (de Windows/Mac)
GOOS=linux GOARCH=amd64 go build -o api-upload

# Build para Windows
GOOS=windows GOARCH=amd64 go build -o api-upload.exe

# Build com versão
VERSION=$(git describe --tags)
go build -ldflags "-X main.Version=$VERSION" -o api-upload
```

### Docker
```dockerfile
# Dockerfile
FROM golang:1.25.0 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-upload .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-upload .
COPY .env .
EXPOSE 8080
CMD ["./api-upload"]
```

### Makefile (Exemplo)
```makefile
.PHONY: help build test run clean

help:
	@echo "Comandos disponíveis:"
	@echo "  make build    - Compilar aplicação"
	@echo "  make test     - Rodar testes"
	@echo "  make run      - Executar aplicação"
	@echo "  make clean    - Limpar artifacts"
	@echo "  make lint     - Verificar código"

build:
	go build -o api-upload

test:
	go test -v ./...

run:
	go run main.go

clean:
	rm -f api-upload
	go clean

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...
	goimports -w .
```

---

## 🔍 Debugging

### Debug com Delve
```bash
# Instalar
go install github.com/go-delve/delve/cmd/dlv@latest

# Rodar com debugger
dlv debug

# Breakpoint interativo
(dlv) break main.main
(dlv) continue
(dlv) print cfg
(dlv) next
```

### Logs Estruturados (Recomendado)
```bash
go get github.com/sirupsen/logrus
```

**Exemplo de uso**:
```go
import "github.com/sirupsen/logrus"

func CreatePost(c *gin.Context) {
    log.WithFields(logrus.Fields{
        "endpoint": "CreatePost",
        "ip": c.ClientIP(),
    }).Info("Tentando criar post")
    
    // ... código ...
    
    log.WithError(err).Error("Erro ao criar post")
}
```

---

## 📝 Code Style

### Go Style Guide
```bash
# Formatar
go fmt ./...

# Lint
go vet ./...

# Mais rigoroso com golangci-lint
golangci-lint run
```

### Exemplo de Código Bem Estruturado
```go
package handlers

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/Anjsvf/api-upload/database"
    "github.com/Anjsvf/api-upload/models"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
)

// CreatePost cria um novo post no banco
func CreatePost(c *gin.Context) {
    // 1. Bind and validate
    var req models.CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondError(c, http.StatusBadRequest, "Dados inválidos", err)
        return
    }
    
    // 2. Business logic
    post := models.Post{
        Title:     req.Title,
        CreatedAt: time.Now(),
    }
    
    // 3. Database operation
    collection := database.GetCollection("posts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    result, err := collection.InsertOne(ctx, post)
    if err != nil {
        respondError(c, http.StatusInternalServerError, "Erro ao salvar", err)
        return
    }
    
    // 4. Response
    c.JSON(http.StatusCreated, post)
}

// Helper function
func respondError(c *gin.Context, status int, message string, err error) {
    c.JSON(status, gin.H{
        "error": message,
        "details": err.Error(),
    })
}
```

---

## 📋 Checklist para PR

Antes de fazer um Pull Request:

- [ ] Código formatado (`go fmt ./...`)
- [ ] Sem erros de lint (`go vet ./...`)
- [ ] Testes passando (`go test ./...`)
- [ ] Testes novos para novo código
- [ ] Documentação atualizada
- [ ] Commits com mensagens descritivas
- [ ] Sem quebra de compatibilidade backward
- [ ] Variáveis de ambiente documentadas
- [ ] Tratamento de erros apropriado

---

## 🔄 CI/CD (Recomendado)

### GitHub Actions Workflow
```yaml
name: Go

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.25.0
    - name: Test
      run: go test -v ./...
    - name: Lint
      run: go vet ./...
```

---

## 📚 Recursos

- [Go Official Docs](https://golang.org/doc/)
- [Gin Documentation](https://github.com/gin-gonic/gin)
- [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

---

## ❓ FAQ

### P: Como adicionar novo middleware?
**R**: Criar função em `middleware/` e usar em `routes/routes.go`:
```go
r.Use(middleware.NovoMiddleware())
```

### P: Como testar com dados reais?
**R**: Use MongoDB Compass ou Mongo CLI:
```bash
mongo
> use vicio_api
> db.posts.find()
```

### P: Como resetar o banco durante testes?
**R**: 
```go
// No teste
collection.DeleteMany(ctx, bson.M{})
```

### P: Qual é o padrão de error handling?
**R**: Sempre retornar status HTTP apropriado + JSON estruturado:
```json
{
  "error": "Mensagem",
  "details": "Detalhes opcionais"
}
```

---

## 🎯 Boas Práticas

1. **Sempre fazer defer cancel()** após `context.WithTimeout`
2. **Usar timeouts** em operações de banco (5-10 segundos)
3. **Validar entrada** antes de usar
4. **Erros com contexto**: Include what went wrong + why
5. **Não expor segredos**: Use variáveis de ambiente
6. **Testes**: Cobertura mínima 80%
7. **Documentação**: Comentar funções exportadas
8. **Logging**: Usar estruturado, não print()

---

**Development Guide - v1.0**
