# Arquitetura da API Upload

## 📐 Visão Geral da Arquitetura

A aplicação segue a arquitetura em camadas com separação clara de responsabilidades:

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Requests                        │
└─────────────────────────────────┬───────────────────────┘
                                  │
┌─────────────────────────────────▼───────────────────────┐
│                  routes/routes.go                       │
│         (Define todas as rotas e middlewares)           │
└─────────────────────────────────┬───────────────────────┘
                                  │
┌─────────────────────────────────▼───────────────────────┐
│               Middleware Layer                          │
│                                                         │
│  - CORS Configuration                                   │
│  - Rate Limiting (middleware/ratelimit.go)              │
│  - IP Detection                                         │
└─────────────────────────────────┬───────────────────────┘
                                  │
┌─────────────────────────────────▼───────────────────────┐
│            Handler Layer (handlers/)                    │
│                                                         │
│  - CreatePost()                                         │
│  - GetFeed()                                            │
│  - GetTags()                                            │
└───┬──────────────┬──────────────┬───────────────────────┘
    │              │              │
    ▼              ▼              ▼
┌────────┐   ┌──────────┐   ┌──────────────┐
│Validator│  │Filters   │  │Models        │
│(links)  │  │(content) │  │(validation)  │
└────┬───┘   └────┬─────┘   └──────────────┘
     │            │
     └────────┬───┘
              │
┌─────────────▼───────────────────────────────────────────┐
│         Database Layer (database/mongodb.go)            │
│                                                         │
│     - Connection Management                             │
│     - Collection Access                                 │
└─────────────────────────────────┬───────────────────────┘
                                  │
┌─────────────────────────────────▼───────────────────────┐
│                   MongoDB                               │
│                                                         │
│     Database: vicio_api                                 │
│     Collections: posts, tags                            │
└─────────────────────────────────────────────────────────┘
```

---

## 🏗️ Componentes Principais

### 1. **Routes Layer** (`routes/routes.go`)
- Define todas as rotas HTTP
- Aplica middlewares globalmente ou por grupo
- Estrutura: `/api/v1` com versionamento

### 2. **Middleware Layer** (`middleware/`)
- **CORS**: Permite requisições de qualquer origem
- **Rate Limiting**: Controla o número de posts por IP
- Executa antes dos handlers

### 3. **Handler Layer** (`handlers/`)
Processa requisições HTTP e coordena a lógica:

```go
CreatePost(c *gin.Context)
├─ Valida estrutura JSON
├─ Aplica filtro de conteúdo
├─ Valida tipo de post
├─ Valida link (YouTube/WhatsApp)
├─ Busca tag no banco
├─ Salva post
└─ Retorna resposta

GetFeed(c *gin.Context)
├─ Parse de parâmetros de filtro
├─ Query ao MongoDB
└─ Retorna posts paginados

GetTags(c *gin.Context)
├─ Query todas as tags
└─ Retorna lista de tags
```

### 4. **Validation Layer** (`validators/`)
Valida URLs antes de salvar:
- YouTube: Verifica validade usando API do YouTube
- WhatsApp: Valida formato de URL

### 5. **Filter Layer** (`filters/`)
- Verifica palavras-chave proibidas
- Retorna erro se conteúdo for censurado

### 6. **Model Layer** (`models/`)
Define estruturas de dados com:
- Validações de binding (`binding:""`)
- Tags JSON para serialização
- Tags BSON para MongoDB

### 7. **Database Layer** (`database/`)
- Gerencia conexão com MongoDB
- Fornece acesso a coleções
- Lazy initialization de conexão

### 8. **Config Layer** (`config/`)
- Carrega variáveis de ambiente
- Fornece valores padrão
- Inicializado no main

---

## 🔄 Fluxos de Dados

### Criar Post
```
Client
   │
   ├─ POST /api/v1/posts
   │  └─ JSON Body
   │
   ▼ Routes
   │
   ├─ CORS Middleware ✓
   │
   ├─ Rate Limit Middleware
   │  ├─ GetClientIP()
   │  ├─ Query posts últimas 24h
   │  └─ Comparar com limite (4)
   │
   ▼ Handler (CreatePost)
   │
   ├─ c.ShouldBindJSON() - Parse e valida campos
   │  └─ Title, Caption, Link, Type, TagID
   │
   ├─ filters.CheckFields() - Verifica palavras proibidas
   │
   ├─ Valida Type (youtube | whatsapp)
   │
   ├─ validators.ValidateLink()
   │  ├─ Se YouTube: Valida com API do YouTube
   │  └─ Se WhatsApp: Valida padrão de URL
   │
   ├─ Busca Tag no MongoDB
   │  └─ Por slug ou ID
   │
   ├─ Cria objeto Post
   │  └─ Adiciona IP, data/hora
   │
   ├─ database.InsertOne(post) - Salva no MongoDB
   │
   └─ Retorna 201 Created + Post salvo
```

### Obter Feed
```
Client
   │
   ├─ GET /api/v1/feed?tag=jogos&type=youtube
   │
   ▼ Routes
   │
   ├─ CORS Middleware ✓
   │
   ▼ Handler (GetFeed)
   │
   ├─ Parse Query Parameters
   │  ├─ tag: string (tag_slug)
   │  ├─ type: string (youtube|whatsapp)
   │  ├─ page: int (padrão 1)
   │  └─ limit: int (padrão 10)
   │
   ├─ Monta filtro BSON
   │  └─ Se tag: {"tag_slug": tag}
   │  └─ Se type: {"type": type}
   │
   ├─ Query MongoDB com paginação
   │  ├─ Skip = (page - 1) * limit
   │  ├─ Limit = limit
   │  └─ Sort por created_at descending
   │
   ├─ Mapeia resultados para modelo
   │
   └─ Retorna 200 OK + Array de Posts
```

---

## 💾 Banco de Dados

### Schema MongoDB

#### Coleção: `posts`
```javascript
{
  _id: ObjectId,
  title: String,           // Indexado? Não (atualmente)
  caption: String,
  link: String,            // Indexado? Não (atualmente)
  type: String,            // Enum: youtube | whatsapp
  tag_id: ObjectId,        // Referência para tags
  tag_name: String,        // Desnormalizado
  tag_slug: String,        // Desnormalizado, indexado para filtros
  tag_emoji: String,
  ip: String,              // Indexado para rate limiting
  created_at: Date         // Indexado para paginação e rate limit
}
```

#### Coleção: `tags`
```javascript
{
  _id: ObjectId,
  name: String,            // Único
  slug: String,            // Único, indexado
  emoji: String,
  created_at: Date
}
```

### Índices Recomendados

Para otimizar queries:

```javascript
// Collection: posts
db.posts.createIndex({ "tag_slug": 1 })
db.posts.createIndex({ "type": 1 })
db.posts.createIndex({ "ip": 1, "created_at": 1 })  // Para rate limiting
db.posts.createIndex({ "created_at": -1 })          // Para paginação

// Collection: tags
db.tags.createIndex({ "slug": 1 }, { unique: true })
db.tags.createIndex({ "name": 1 }, { unique: true })
```

---

## 🔐 Segurança

### Rate Limiting
- Limite: 4 posts por IP a cada 24 horas
- Stored: No MongoDB (histórico de posts)
- IP Detection: Headers `X-Forwarded-For`, `X-Real-IP`, fallback para `ClientIP()`

### Validação de Entrada
- **Body Binding**: Gin valida tipos e estrutura
- **Tag Validation**: Confirma tag existe antes de salvar
- **Link Validation**: URLs validadas antes de armazenar
- **Content Filter**: Palavras-chave proibidas bloqueadas

### CORS
- Allowed Origins: `*` (qualquer origem)
- Allowed Methods: `GET`, `POST`
- Credentials: Desabilitado

### IP Privacidade
- IP armazenado no banco (não retornado no GET)
- Usado apenas para rate limiting

---

## 📊 Performance e Escalabilidade

### Otimizações Atuais
1. **Índices MongoDB**: Query em tag_slug, type, created_at
2. **Paginação**: Limite de resultados por página
3. **Timeouts**: 5 segundos para operações de banco
4. **Connection Pooling**: MongoDB driver gerencia pool

### Possíveis Melhorias
1. **Cache**: Redis para tags e feed popular
2. **Denormalização**: Dados de tag já armazenados (feito)
3. **Query Optimization**: Adicionar mais índices
4. **Async Processing**: Background jobs para validações
5. **CDN**: Cache de conteúdo estático

### Capacidade Estimada
- Com índices adequados: ~10k requests/segundo
- MongoDB local: Sem limitações até ~GB dados
- MongoDB Atlas: Escalável com tier apropriado

---

## 🔄 Padrões de Projeto

### 1. **Handler Pattern**
Cada endpoint é uma função que recebe `*gin.Context`:
```go
func CreatePost(c *gin.Context) {
    // Lógica
}
```

### 2. **Middleware Pattern**
Middlewares executam antes do handler:
```go
posts := v1.Group("/posts")
posts.Use(middleware.RateLimitByIP())
posts.POST("", handlers.CreatePost)
```

### 3. **Dependency Injection**
Gerenciador global de dependências:
```go
db := database.GetCollection("posts")
```

### 4. **Error Handling**
Erros retornam JSON estruturado:
```json
{
  "error": "Mensagem de erro",
  "details": "Detalhes opcionais"
}
```

### 5. **Context Timeout**
Operações de banco com timeout:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

---

## 🚀 Deployment

### Ambientes

#### Development
- MongoDB local em `localhost:27017`
- .env com valores de teste
- Logging verbose

#### Production
- MongoDB Atlas ou hosted
- .env com valores reais
- YOUTUBE_API_KEY necessária
- CORS pode ser restringido a domínios específicos

### Variáveis Críticas para Produção

```env
# OBRIGATÓRIO
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/
MONGO_DB_NAME=vicio_api_prod
YOUTUBE_API_KEY=AIza...

# RECOMENDADO
PORT=8080
```

---

## 📈 Monitoramento

### Métricas Recomendadas
1. **Requests por segundo**
2. **Taxa de erro (4xx, 5xx)**
3. **Tempo médio de resposta**
4. **Posts por dia**
5. **IPs únicos por dia**
6. **Taxa de acerto de rate limit**
7. **Conexões MongoDB ativas**
8. **Tamanho da coleção posts**

### Health Check
```
GET /health
Response: {"status": "ok"}
```

---

## 🔄 Ciclo de Vida da Aplicação

1. **Inicialização** (main.go)
   - Load config
   - Connect MongoDB
   - Setup YouTube API
   - Seed tags
   - Setup routes

2. **Runtime**
   - Aceita requests HTTP
   - Aplica middlewares
   - Processa handlers
   - Retorna respostas

3. **Shutdown**
   - Desconecta MongoDB
   - Encerra server

---

## 🛠️ Troubleshooting

### MongoDB não conecta
```
Erro: connection refused
Solução: Verificar MONGO_URI, MongoDB rodando
```

### Rate limit não funciona
```
Erro: Posts duplicados do mesmo IP
Solução: Verificar índice em posts(ip, created_at)
```

### YouTube API retorna erro
```
Erro: Invalid YouTube link
Solução: Verificar YOUTUBE_API_KEY válida e quota
```

---

**Documento de Arquitetura - v1.0**  
Última atualização: 2024
