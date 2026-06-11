# API Examples & Testing Guide

Exemplos práticos de como usar cada endpoint da API com diferentes cenários.

## 📌 Base URL
```
http://localhost:8080
```

---

## 🏥 Health Check

### Verificar Status do Servidor

**Requisição:**
```bash
GET /health
```

**cURL:**
```bash
curl -X GET http://localhost:8080/health
```

**Resposta (200 OK):**
```json
{
  "status": "ok"
}
```

---

## 🏷️ Tags Endpoints

### Listar Todas as Tags

**Requisição:**
```bash
GET /api/v1/tags
```

**cURL:**
```bash
curl -X GET http://localhost:8080/api/v1/tags
```

**Resposta (200 OK):**
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
  },
  {
    "id": "507f1f77bcf86cd799439013",
    "name": "Programação",
    "slug": "programacao",
    "emoji": "💻",
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": "507f1f77bcf86cd799439014",
    "name": "Arte",
    "slug": "arte",
    "emoji": "🎨",
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": "507f1f77bcf86cd799439015",
    "name": "Música",
    "slug": "musica",
    "emoji": "🎵",
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": "507f1f77bcf86cd799439016",
    "name": "Cinema",
    "slug": "cinema",
    "emoji": "🎬",
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

---

## 📰 Feed Endpoints

### Obter Todos os Posts (Feed Completo)

**Requisição:**
```bash
GET /api/v1/feed
```

**cURL:**
```bash
curl -X GET http://localhost:8080/api/v1/feed
```

**Resposta (200 OK):**
```json
[
  {
    "id": "507f1f77bcf86cd799439020",
    "title": "Tutorial Go Avançado",
    "caption": "Aprenda técnicas avançadas em Go para produção",
    "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
    "type": "youtube",
    "tag_id": "507f1f77bcf86cd799439013",
    "tag_name": "Programação",
    "tag_slug": "programacao",
    "tag_emoji": "💻",
    "created_at": "2024-01-20T10:30:00Z"
  }
]
```

### Filtrar Feed por Tag

**Requisição:**
```bash
GET /api/v1/feed?tag=jogos
```

**cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/feed?tag=jogos"
```

**Resposta (200 OK):**
```json
[
  {
    "id": "507f1f77bcf86cd799439021",
    "title": "Novo Jogo RPG",
    "caption": "Um jogo incrível de RPG para 2024",
    "link": "https://chat.whatsapp.com/grupo-de-jogos",
    "type": "whatsapp",
    "tag_id": "507f1f77bcf86cd799439011",
    "tag_name": "Jogos",
    "tag_slug": "jogos",
    "tag_emoji": "🎮",
    "created_at": "2024-01-20T14:30:00Z"
  }
]
```

### Filtrar Feed por Tipo

**Requisição:**
```bash
GET /api/v1/feed?type=youtube
```

**cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/feed?type=youtube"
```

### Filtrar por Tag E Tipo

**Requisição:**
```bash
GET /api/v1/feed?tag=programacao&type=youtube
```

**cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/feed?tag=programacao&type=youtube"
```

### Paginação

**Requisição:**
```bash
GET /api/v1/feed?page=2&limit=5
```

**Parâmetros:**
- `page`: Número da página (padrão: 1)
- `limit`: Itens por página (padrão: 10)

**cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/feed?page=2&limit=5"
```

### Combinar Todos os Filtros

**Requisição:**
```bash
GET /api/v1/feed?tag=programacao&type=youtube&page=1&limit=10
```

**cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/feed?tag=programacao&type=youtube&page=1&limit=10"
```

---

## ✏️ Post Endpoints

### Criar Post - YouTube

**Requisição:**
```bash
POST /api/v1/posts
Content-Type: application/json

{
  "title": "Tutorial Go Completo",
  "caption": "Um tutorial muito bom sobre Go do zero ao avançado",
  "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
  "type": "youtube",
  "tag_id": "programacao"
}
```

**cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Tutorial Go Completo",
    "caption": "Um tutorial muito bom sobre Go do zero ao avançado",
    "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
    "type": "youtube",
    "tag_id": "programacao"
  }'
```

**Resposta (201 Created):**
```json
{
  "id": "507f1f77bcf86cd799439030",
  "title": "Tutorial Go Completo",
  "caption": "Um tutorial muito bom sobre Go do zero ao avançado",
  "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
  "type": "youtube",
  "tag_id": "507f1f77bcf86cd799439013",
  "tag_name": "Programação",
  "tag_slug": "programacao",
  "tag_emoji": "💻",
  "created_at": "2024-01-20T15:45:00Z"
}
```

### Criar Post - WhatsApp

**Requisição:**
```bash
POST /api/v1/posts
Content-Type: application/json

{
  "title": "Grupo de Jogadores",
  "caption": "Comunidade de jogadores online para jogar juntos",
  "link": "https://chat.whatsapp.com/HUx1234567890",
  "type": "whatsapp",
  "tag_id": "jogos"
}
```

**cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Grupo de Jogadores",
    "caption": "Comunidade de jogadores online para jogar juntos",
    "link": "https://chat.whatsapp.com/HUx1234567890",
    "type": "whatsapp",
    "tag_id": "jogos"
  }'
```

**Resposta (201 Created):**
```json
{
  "id": "507f1f77bcf86cd799439031",
  "title": "Grupo de Jogadores",
  "caption": "Comunidade de jogadores online para jogar juntos",
  "link": "https://chat.whatsapp.com/HUx1234567890",
  "type": "whatsapp",
  "tag_id": "507f1f77bcf86cd799439011",
  "tag_name": "Jogos",
  "tag_slug": "jogos",
  "tag_emoji": "🎮",
  "created_at": "2024-01-20T15:50:00Z"
}
```

### Criar Post - Por ID da Tag

Você também pode usar o ID ObjectID da tag em vez do slug:

**cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Novo Video",
    "caption": "Descrição do video",
    "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
    "type": "youtube",
    "tag_id": "507f1f77bcf86cd799439013"
  }'
```

---

## ❌ Exemplos de Erros

### Erro: Campos Inválidos

**Requisição:**
```bash
POST /api/v1/posts
Content-Type: application/json

{
  "title": "Go",  # Menos de 3 caracteres
  "caption": "A",  # Menos de 3 caracteres
  "link": "https://www.youtube.com/watch?v=invalid",
  "type": "youtube"
  # Tag ID faltando
}
```

**Resposta (400 Bad Request):**
```json
{
  "error": "Dados inválidos",
  "details": "Key: 'CreatePostRequest.Title' Error:Field validation for 'Title' failed on the 'min' tag"
}
```

### Erro: Tag Não Encontrada

**Requisição:**
```bash
POST /api/v1/posts
Content-Type: application/json

{
  "title": "Video Legal",
  "caption": "Um video bem legal",
  "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
  "type": "youtube",
  "tag_id": "tag_inexistente"
}
```

**Resposta (400 Bad Request):**
```json
{
  "error": "Tag de vício não encontrada. Use o slug (ex: 'jogos') ou o ID da tag"
}
```

### Erro: Tipo Inválido

**Requisição:**
```bash
POST /api/v1/posts
Content-Type: application/json

{
  "title": "Video Interessante",
  "caption": "Um video interessante",
  "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
  "type": "tiktok",  # Tipo inválido
  "tag_id": "programacao"
}
```

**Resposta (400 Bad Request):**
```json
{
  "error": "Tipo inválido. Use 'youtube' ou 'whatsapp'"
}
```

### Erro: Link Inválido (YouTube)

**Requisição:**
```bash
POST /api/v1/posts
Content-Type: application/json

{
  "title": "Video Interessante",
  "caption": "Um video interessante",
  "link": "https://www.youtube.com/watch?v=INVALID123",
  "type": "youtube",
  "tag_id": "programacao"
}
```

**Resposta (422 Unprocessable Entity):**
```json
{
  "error": "Link do YouTube inválido ou privado"
}
```

### Erro: Link Inválido (WhatsApp)

**Requisição:**
```bash
POST /api/v1/posts
Content-Type: application/json

{
  "title": "Grupo",
  "caption": "Um grupo",
  "link": "https://invalid-link.com/group",
  "type": "whatsapp",
  "tag_id": "jogos"
}
```

**Resposta (422 Unprocessable Entity):**
```json
{
  "error": "Link do WhatsApp inválido"
}
```

### Erro: Conteúdo Filtrado

**Requisição:**
```bash
POST /api/v1/posts
Content-Type: application/json

{
  "title": "Conteudo censurado",  # Contém palavra proibida
  "caption": "Descrição",
  "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
  "type": "youtube",
  "tag_id": "programacao"
}
```

**Resposta (422 Unprocessable Entity):**
```json
{
  "error": "Conteúdo não permitido"
}
```

### Erro: Rate Limit Excedido

**Requisição:**
```
POST /api/v1/posts (5ª vez do mesmo IP em 24h)
```

**Resposta (429 Too Many Requests):**
```json
{
  "error": "Limite de postagens atingido. Máximo de 4 por dia."
}
```

---

## 🧪 Script de Teste Completo

Testar todos os endpoints em sequência:

**Salvar como `test.sh`:**
```bash
#!/bin/bash

BASE_URL="http://localhost:8080"
echo "=== TESTANDO API UPLOAD ==="
echo ""

echo "1. Health Check"
curl -s -X GET "$BASE_URL/health" | jq .
echo ""

echo "2. Listar Tags"
curl -s -X GET "$BASE_URL/api/v1/tags" | jq '.[] | {name, slug, emoji}'
echo ""

echo "3. Feed Vazio"
curl -s -X GET "$BASE_URL/api/v1/feed" | jq 'length'
echo ""

echo "4. Criar Post YouTube"
curl -s -X POST "$BASE_URL/api/v1/posts" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Tutorial Go Avançado",
    "caption": "Um tutorial excelente sobre Go",
    "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
    "type": "youtube",
    "tag_id": "programacao"
  }' | jq .
echo ""

echo "5. Feed com Posts"
curl -s -X GET "$BASE_URL/api/v1/feed" | jq 'length'
echo ""

echo "6. Feed Filtrado por Tag"
curl -s -X GET "$BASE_URL/api/v1/feed?tag=programacao" | jq '.[] | {title, type}'
echo ""

echo "=== TESTES COMPLETOS ==="
```

**Executar:**
```bash
chmod +x test.sh
./test.sh
```

---

## 📱 Exemplos com Postman

### Importar Collection JSON

```json
{
  "info": {
    "name": "API Upload",
    "description": "Coleção de testes para API Upload",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "url": "{{base_url}}/health"
      }
    },
    {
      "name": "Get Tags",
      "request": {
        "method": "GET",
        "url": "{{base_url}}/api/v1/tags"
      }
    },
    {
      "name": "Get Feed",
      "request": {
        "method": "GET",
        "url": "{{base_url}}/api/v1/feed"
      }
    },
    {
      "name": "Create Post YouTube",
      "request": {
        "method": "POST",
        "url": "{{base_url}}/api/v1/posts",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\"title\": \"Tutorial Go\", \"caption\": \"Um tutorial\", \"link\": \"https://www.youtube.com/watch?v=jNQXAC9IVRw\", \"type\": \"youtube\", \"tag_id\": \"programacao\"}"
        }
      }
    }
  ],
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080"
    }
  ]
}
```

---

## 💡 Dicas

1. **Use slugs de tag** ao invés de IDs (mais fácil de memorizar): `tag_id: "jogos"`
2. **Teste paginação** incrementando `page` e `limit`
3. **Combine filtros** para testes mais específicos
4. **Verifique rate limit** tentando criar 5 posts do mesmo IP
5. **Use ferramentas visuais** como Postman ou Thunder Client para facilitar testes

---

**API Examples - v1.0**
