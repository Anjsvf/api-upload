# API Upload - Backend Documentation

## 📋 Visão Geral

API Backend para compartilhamento de conteúdo (vídeos do YouTube e grupos do WhatsApp) com sistema de tags, filtros e rate limiting. Desenvolvida em **Go** com **Gin** e **MongoDB**.

O projeto permite que usuários criem postagens categorizadas e as compartilhem através de uma API RESTful com proteção contra abuso.

---

## 🏗️ Estrutura do Projeto

```
api-upload/
├── main.go                 # Ponto de entrada da aplicação
├── go.mod                  # Módulo Go e dependências
├── Read.md                 # Documentação adicional
│
├── config/                 # Configurações
│   └── config.go          # Gerenciamento de variáveis de ambiente
│
├── database/              # Camada de dados
│   └── mongodb.go         # Conexão e operações MongoDB
│
├── models/                # Modelos de dados
│   ├── post.go            # Modelo Post e requisições
│   └── tag.go             # Modelo Tag
│
├── handlers/              # Controladores HTTP
│   ├── post_handler.go    # Endpoints de posts
│   └── tag_handler.go     # Endpoints de tags
│
├── routes/               # Definição de rotas
│   └── routes.go         # Setup das rotas da API
│
├── middleware/           # Middlewares
│   └── ratelimit.go      # Rate limiting por IP
│
├── validators/           # Validadores
│   └── link_validator.go # Validação de URLs
│
├── filters/              # Filtros de conteúdo
│   └── word_filter.go    # Filtro de palavras-chave
│
└── seed/                 # Dados iniciais
    └── seed.go           # Seed de tags padrão
```

---

## 🛠️ Tecnologias Utilizadas

- **Linguagem**: Go 1.25.0
- **Framework Web**: Gin Gonic v1.9.1
- **Banco de Dados**: MongoDB v1.13.1
- **CORS**: Gin CORS v1.5.0
- **Environment**: godotenv v1.5.1

### Dependências Principais
```go
github.com/gin-gonic/gin v1.9.1           // Framework web
go.mongodb.org/mongo-driver v1.13.1       // Driver MongoDB
github.com/joho/godotenv v1.5.1           // Gerenciamento de .env
github.com/gin-contrib/cors v1.5.0        // CORS
```

---

## 📦 Instalação e Setup

### Pré-requisitos
- Go 1.25.0 ou superior
- MongoDB rodando localmente ou em cloud (MongoDB Atlas)
- Uma chave de API do YouTube (opcional, para validação de vídeos)

### 1. Clonar o Repositório
```bash
git clone https://github.com/Anjsvf/api-upload.git
cd api-upload
```

### 2. Instalar Dependências
```bash
go mod download
go mod tidy
```

### 3. Configurar Variáveis de Ambiente
Criar arquivo `.env` na raiz do projeto:
```env
# Configurações MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=vicio_api

# Porta do servidor
PORT=8080

# Chave da API do YouTube (opcional)
YOUTUBE_API_KEY=sua_chave_aqui
```

### 4. Executar a Aplicação
```bash
go run main.go
```

Saída esperada:
```
Servidor rodando na alturas na porta 8080
```

---

## ⚙️ Configuração

### Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `MONGO_URI` | String de conexão MongoDB | `mongodb://localhost:27017` |
| `MONGO_DB_NAME` | Nome do banco de dados | `vicio_api` |
| `PORT` | Porta do servidor | `8080` |
| `YOUTUBE_API_KEY` | Chave API do YouTube | `` (vazio) |

### CORS

O servidor aceita requisições de qualquer origem (CORS habilitado para `*`):
- **Origens Permitidas**: `*` (todas)
- **Métodos Permitidos**: `GET`, `POST`
- **Headers Permitidos**: `Content-Type`
- **Credentials**: Desabilitado

---

## 📚 Modelos de Dados

### Post
```go
type Post struct {
    ID        ObjectID  // ID único MongoDB
    Title     string    // Título (3-100 caracteres)
    Caption   string    // Descrição (3-500 caracteres)
    Link      string    // URL validada
    Type      PostType  // "youtube" ou "whatsapp"
    TagID     ObjectID  // ID da tag
    TagName   string    // Nome da tag
    TagSlug   string    // Slug da tag (ex: "jogos")
    TagEmoji  string    // Emoji associado
    IP        string    // IP de origem (não retornado)
    CreatedAt time.Time // Data de criação
}
```

**PostType**:
- `youtube` - Vídeo do YouTube
- `whatsapp` - Grupo do WhatsApp

### Tag
```go
type Tag struct {
    ID        ObjectID  // ID único MongoDB
    Name      string    // Nome da tag
    Slug      string    // Slug (URL friendly)
    Emoji     string    // Emoji representativo
    CreatedAt time.Time // Data de criação
}
```

---

## 🔗 Endpoints da API

### Health Check
Verifica se o servidor está ativo.

**Endpoint**: `GET /health`

**Resposta**:
```json
{
  "status": "ok"
}
```

---

### Obter Tags
Lista todas as tags disponíveis.

**Endpoint**: `GET /api/v1/tags`

**Resposta**:
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "Jogos",
    "slug": "jogos",
    "emoji": "🎮",
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": "507f1f77bcf86cd799439012",
    "name": "Cursos",
    "slug": "cursos",
    "emoji": "📚",
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

---

### Obter Feed
Retorna postagens com filtros opcionais.

**Endpoint**: `GET /api/v1/feed`

**Query Parameters**:
| Parâmetro | Tipo | Descrição | Padrão |
|-----------|------|-----------|--------|
| `tag` | string | Slug da tag (ex: "jogos") | - |
| `type` | string | Tipo de conteúdo ("youtube" ou "whatsapp") | - |
| `page` | number | Número da página | 1 |
| `limit` | number | Itens por página | 10 |

**Exemplo de Requisição**:
```
GET /api/v1/feed?tag=jogos&type=youtube&page=1&limit=5
```

**Resposta**:
```json
[
  {
    "id": "507f1f77bcf86cd799439013",
    "title": "Tutorial Go Avançado",
    "caption": "Aprenda técnicas avançadas em Go",
    "link": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
    "type": "youtube",
    "tag_id": "507f1f77bcf86cd799439011",
    "tag_name": "Programação",
    "tag_slug": "programacao",
    "tag_emoji": "💻",
    "created_at": "2024-01-20T14:30:00Z"
  }
]
```

---

### Criar Post
Cria uma nova postagem (com rate limiting).

**Endpoint**: `POST /api/v1/posts`

**Rate Limiting**: Máximo 4 posts por IP a cada 24 horas

**Body**:
```json
{
  "title": "Novo Vídeo Interessante",
  "caption": "Um vídeo muito legal sobre Go",
  "link": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "type": "youtube",
  "tag_id": "jogos"
}
```

**Validações**:
- `title`: Obrigatório, 3-100 caracteres
- `caption`: Obrigatório, 3-500 caracteres
- `link`: Obrigatório, URL válida (YouTube ou WhatsApp)
- `type`: Obrigatório, "youtube" ou "whatsapp"
- `tag_id`: Obrigatório, slug ou ID válido da tag

**Filtros de Conteúdo**:
- Palavras-chave proibidas são verificadas em `title` e `caption`

**Validação de Links**:
- **YouTube**: Valida usando a API do YouTube (requer chave)
- **WhatsApp**: Valida formato de link do grupo

**Resposta de Sucesso** (201):
```json
{
  "id": "507f1f77bcf86cd799439014",
  "title": "Novo Vídeo Interessante",
  "caption": "Um vídeo muito legal sobre Go",
  "link": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "type": "youtube",
  "tag_id": "507f1f77bcf86cd799439011",
  "tag_name": "Jogos",
  "tag_slug": "jogos",
  "tag_emoji": "🎮",
  "created_at": "2024-01-20T15:45:00Z"
}
```

**Erros Possíveis**:
- `400 Bad Request` - Dados inválidos
- `422 Unprocessable Entity` - Conteúdo filtrado, link inválido ou tag não encontrada
- `429 Too Many Requests` - Limite de taxa excedido

---

## 🔒 Rate Limiting

O rate limiting é aplicado ao endpoint `POST /api/v1/posts`:

- **Limite**: 4 posts por IP a cada 24 horas
- **Identificação**: Por IP do cliente (com suporte para proxies usando headers `X-Forwarded-For` e `X-Real-IP`)

Quando o limite é excedido, retorna:
```json
{
  "error": "Limite de postagens atingido. Máximo de 4 por dia."
}
```

---

## 🗄️ Banco de Dados

### Coleções MongoDB

#### `posts`
```json
{
  "_id": ObjectId,
  "title": "string",
  "caption": "string",
  "link": "string",
  "type": "youtube|whatsapp",
  "tag_id": ObjectId,
  "tag_name": "string",
  "tag_slug": "string",
  "tag_emoji": "string",
  "ip": "string",
  "created_at": ISODate
}
```

#### `tags`
```json
{
  "_id": ObjectId,
  "name": "string",
  "slug": "string",
  "emoji": "string",
  "created_at": ISODate
}
```

### Seed de Tags

Tags padrão criadas automaticamente na inicialização:
- Jogos 🎮
- Cursos 📚
- Programação 💻
- Arte 🎨
- Música 🎵
- Cinema 🎬

---

## 🚀 Build e Deployment

### Build Executável

```bash
# Linux/macOS
GOOS=linux GOARCH=amd64 go build -o api-upload

# Windows
go build -o api-upload.exe
```

### Docker (Exemplo)

Se quiser containerizar:

```dockerfile
FROM golang:1.25.0 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o api-upload

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-upload .
EXPOSE 8080
CMD ["./api-upload"]
```

---

## 🧪 Testando a API

### Usando curl

```bash
# Health check
curl http://localhost:8080/health

# Obter tags
curl http://localhost:8080/api/v1/tags

# Obter feed
curl http://localhost:8080/api/v1/feed?tag=jogos&limit=5

# Criar post
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Novo Vídeo",
    "caption": "Descrição",
    "link": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
    "type": "youtube",
    "tag_id": "jogos"
  }'
```

### Usando Postman

1. Importar requisições
2. Configurar variável `base_url` = `http://localhost:8080`
3. Executar requisições da coleção

---

## 🔍 Estrutura de Código

### Package `main`
- Inicializa a aplicação
- Conecta ao MongoDB
- Configura handlers e seeds
- Inicia o servidor

### Package `config`
- Carrega variáveis de ambiente via `.env`
- Fornece valores padrão

### Package `database`
- Gerencia conexão com MongoDB
- Fornece acesso às coleções

### Package `handlers`
- `CreatePost()` - Handler para criar posts
- `GetFeed()` - Handler para obter feed com filtros
- `GetTags()` - Handler para listar tags

### Package `middleware`
- `RateLimitByIP()` - Middleware de rate limiting
- `GetClientIP()` - Detecta IP do cliente (com suporte a proxies)

### Package `validators`
- `ValidateYouTubeLink()` - Valida URLs do YouTube
- `ValidateWhatsAppLink()` - Valida URLs de grupos WhatsApp

### Package `filters`
- `CheckFields()` - Verifica palavras-chave proibidas

### Package `models`
- `Post` - Modelo de postagem
- `Tag` - Modelo de categoria
- `CreatePostRequest` - DTO para criação de posts
- `FeedFilter` - Filtros para o feed

### Package `routes`
- `Setup()` - Configura todas as rotas e middlewares

### Package `seed`
- `RunTags()` - Popula tags padrão no banco

---

## 📝 Fluxo de Criação de Post

```
POST /api/v1/posts
    ↓
[Body Binding]
    ↓
[Validação de Campos]
    ↓
[Rate Limiting Middleware]
    ↓
[Filtro de Palavras-chave]
    ↓
[Validação de Tipo]
    ↓
[Validação de Link (YouTube/WhatsApp)]
    ↓
[Busca da Tag]
    ↓
[Inserção no MongoDB]
    ↓
[Resposta 201 Created]
```

---

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

---

## 📄 Licença

Este projeto está sob licença MIT. Veja o arquivo LICENSE para mais detalhes.

---

## 👤 Autor

**Anjsvf**  
GitHub: [@Anjsvf](https://github.com/Anjsvf)

---

## ❓ Suporte

Para dúvidas ou problemas:
1. Abra uma issue no GitHub
2. Verifique a documentação acima
3. Cheque os logs do servidor

---

## 📊 Endpoints Resumo

| Método | Endpoint | Descrição | Rate Limit |
|--------|----------|-----------|------------|
| `GET` | `/health` | Health check | ❌ |
| `GET` | `/api/v1/tags` | Listar tags | ❌ |
| `GET` | `/api/v1/feed` | Listar posts com filtros | ❌ |
| `POST` | `/api/v1/posts` | Criar novo post | ✅ (4/24h) |

---

**Última Atualização**: 2024
