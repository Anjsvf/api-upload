# Database Schema & Operations

Documentação detalhada sobre o esquema do MongoDB e operações disponíveis.

---

## 📊 Visão Geral do Banco

**Banco**: `vicio_api` (configurável via `MONGO_DB_NAME`)

**Coleções**:
- `posts` - Postagens de conteúdo
- `tags` - Categorias de conteúdo

---

## 🏷️ Coleção: `tags`

### Esquema

```javascript
{
  _id: ObjectId,
  name: String,          // Único, nome legível
  slug: String,          // Único, URL-friendly
  emoji: String,         // Emoji único de 1 caractere
  created_at: Date       // Timestamp de criação
}
```

### Exemplo de Documento

```json
{
  "_id": ObjectId("507f1f77bcf86cd799439011"),
  "name": "Programação",
  "slug": "programacao",
  "emoji": "💻",
  "created_at": ISODate("2024-01-15T10:30:00.000Z")
}
```

### Validações

| Campo | Tipo | Obrigatório | Validações |
|-------|------|-------------|-----------|
| `_id` | ObjectId | Sim | Auto-gerado |
| `name` | String | Sim | Único, 3-50 caracteres |
| `slug` | String | Sim | Único, lowercase, URL-friendly |
| `emoji` | String | Sim | 1 caractere emoji |
| `created_at` | Date | Sim | Auto-gerado no insert |

### Índices

```javascript
// Índices recomendados
db.tags.createIndex({ "slug": 1 }, { unique: true })
db.tags.createIndex({ "name": 1 }, { unique: true })

// Ver índices
db.tags.getIndexes()

// Dropear índice
db.tags.dropIndex("slug_1")
```

### Operações

#### Inserir Tag
```javascript
db.tags.insertOne({
  name: "Jogos",
  slug: "jogos",
  emoji: "🎮",
  created_at: new Date()
})
```

#### Buscar Tag por Slug
```javascript
db.tags.findOne({ slug: "jogos" })
```

#### Buscar Tag por ID
```javascript
db.tags.findOne({ 
  _id: ObjectId("507f1f77bcf86cd799439011") 
})
```

#### Listar Todas as Tags
```javascript
db.tags.find({}).sort({ created_at: -1 })
```

#### Atualizar Tag
```javascript
db.tags.updateOne(
  { slug: "jogos" },
  { $set: { emoji: "🕹️" } }
)
```

#### Deletar Tag
```javascript
db.tags.deleteOne({ slug: "jogos" })
```

---

## 📰 Coleção: `posts`

### Esquema

```javascript
{
  _id: ObjectId,
  title: String,         // Título principal
  caption: String,       // Descrição detalhada
  link: String,          // URL do conteúdo
  type: String,          // "youtube" | "whatsapp"
  tag_id: ObjectId,      // Referência para tags._id
  tag_name: String,      // Desnormalizado (denormalization)
  tag_slug: String,      // Desnormalizado
  tag_emoji: String,     // Desnormalizado
  ip: String,            // IP do criador (para rate limit)
  created_at: Date       // Timestamp de criação
}
```

### Exemplo de Documento

```json
{
  "_id": ObjectId("507f1f77bcf86cd799439020"),
  "title": "Tutorial Go Avançado",
  "caption": "Aprenda técnicas avançadas em Go para produção",
  "link": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
  "type": "youtube",
  "tag_id": ObjectId("507f1f77bcf86cd799439013"),
  "tag_name": "Programação",
  "tag_slug": "programacao",
  "tag_emoji": "💻",
  "ip": "192.168.1.100",
  "created_at": ISODate("2024-01-20T10:30:00.000Z")
}
```

### Validações

| Campo | Tipo | Obrigatório | Validações |
|-------|------|-------------|-----------|
| `_id` | ObjectId | Sim | Auto-gerado |
| `title` | String | Sim | 3-100 caracteres |
| `caption` | String | Sim | 3-500 caracteres |
| `link` | String | Sim | URL válida, protocolo https |
| `type` | String | Sim | "youtube" ou "whatsapp" |
| `tag_id` | ObjectId | Sim | Deve existir em tags |
| `tag_name` | String | Sim | Cópia de tags.name |
| `tag_slug` | String | Sim | Cópia de tags.slug |
| `tag_emoji` | String | Sim | Cópia de tags.emoji |
| `ip` | String | Sim | Formato IP válido |
| `created_at` | Date | Sim | Auto-gerado no insert |

### Índices Recomendados

```javascript
// Índices para queries comuns
db.posts.createIndex({ "tag_slug": 1 })
db.posts.createIndex({ "type": 1 })
db.posts.createIndex({ "created_at": -1 })

// Índice composto para rate limiting
db.posts.createIndex({ 
  "ip": 1, 
  "created_at": 1 
})

// Índice para paginação eficiente
db.posts.createIndex({ 
  "tag_slug": 1,
  "created_at": -1
})

// Ver todos os índices
db.posts.getIndexes()
```

### Performance dos Índices

```javascript
// Analisar query performance
db.posts.find({ 
  tag_slug: "programacao",
  created_at: { $gte: ISODate("2024-01-01") }
}).explain("executionStats")

// Resultados esperados:
// - "executionStages": { "stage": "COLLSCAN" } = SEM índice (ruim)
// - "executionStages": { "stage": "IXSCAN" } = COM índice (bom)
```

### Operações

#### Inserir Post
```javascript
db.posts.insertOne({
  title: "Tutorial Go",
  caption: "Um tutorial excelente",
  link: "https://www.youtube.com/watch?v=...",
  type: "youtube",
  tag_id: ObjectId("507f1f77bcf86cd799439013"),
  tag_name: "Programação",
  tag_slug: "programacao",
  tag_emoji: "💻",
  ip: "192.168.1.100",
  created_at: new Date()
})
```

#### Buscar Posts por Tag
```javascript
db.posts.find({ 
  tag_slug: "programacao" 
}).sort({ created_at: -1 })
```

#### Buscar Posts por Tipo
```javascript
db.posts.find({ 
  type: "youtube" 
}).sort({ created_at: -1 })
```

#### Buscar Posts de um IP
```javascript
db.posts.find({ 
  ip: "192.168.1.100" 
})
```

#### Rate Limit - Contar Posts do IP em 24h
```javascript
// Contar quantos posts este IP criou no último dia
db.posts.countDocuments({
  ip: "192.168.1.100",
  created_at: {
    $gte: new Date(new Date().getTime() - 24 * 60 * 60 * 1000)
  }
})

// Se retornar > 4, aplicar rate limit
```

#### Paginação
```javascript
// Page 1, 10 itens por página
db.posts.find({ tag_slug: "programacao" })
  .sort({ created_at: -1 })
  .skip(0)
  .limit(10)

// Page 2, 10 itens por página
db.posts.find({ tag_slug: "programacao" })
  .sort({ created_at: -1 })
  .skip(10)
  .limit(10)

// Page 3, 5 itens por página
db.posts.find({ tag_slug: "programacao" })
  .sort({ created_at: -1 })
  .skip(10)
  .limit(5)
```

#### Deletar Post
```javascript
db.posts.deleteOne({ 
  _id: ObjectId("507f1f77bcf86cd799439020") 
})
```

#### Deletar Posts Antigos (Limpeza)
```javascript
// Deletar posts com mais de 30 dias
db.posts.deleteMany({
  created_at: {
    $lt: new Date(new Date().getTime() - 30 * 24 * 60 * 60 * 1000)
  }
})
```

---

## 🔗 Relacionamentos

### Relacionamento Tags ↔ Posts

**Tipo**: One-to-Many (Uma tag pode ter muitos posts)

```
tags (1)
  └─── _id
  └─── name: "Programação"
  └─── slug: "programacao"

posts (Many)
  └─── tag_id: <referência para tag._id>
  └─── tag_name: "Programação" (denormalizado)
  └─── tag_slug: "programacao" (denormalizado)
```

**Vantagens da Denormalização**:
- ✅ Queries mais rápidas (não precisa JOIN)
- ✅ Menos operações de banco
- ✅ Tag info sempre disponível no post
- ❌ Inconsistência se tag for renomeada (problema menor)

**Alternativa com Lookup** (se necessário):
```javascript
db.posts.aggregate([
  {
    $lookup: {
      from: "tags",
      localField: "tag_id",
      foreignField: "_id",
      as: "tag_details"
    }
  }
])
```

---

## 📈 Estatísticas e Análises

### Total de Posts
```javascript
db.posts.countDocuments({})
```

### Posts por Tag
```javascript
db.posts.aggregate([
  {
    $group: {
      _id: "$tag_slug",
      count: { $sum: 1 }
    }
  },
  {
    $sort: { count: -1 }
  }
])
```

### Posts por Tipo
```javascript
db.posts.aggregate([
  {
    $group: {
      _id: "$type",
      count: { $sum: 1 }
    }
  }
])
```

### Posts por Dia (Últimos 7 dias)
```javascript
db.posts.aggregate([
  {
    $match: {
      created_at: {
        $gte: new Date(new Date().getTime() - 7 * 24 * 60 * 60 * 1000)
      }
    }
  },
  {
    $group: {
      _id: {
        $dateToString: {
          format: "%Y-%m-%d",
          date: "$created_at"
        }
      },
      count: { $sum: 1 }
    }
  },
  {
    $sort: { _id: 1 }
  }
])
```

### IPs com Maior Atividade
```javascript
db.posts.aggregate([
  {
    $group: {
      _id: "$ip",
      count: { $sum: 1 }
    }
  },
  {
    $sort: { count: -1 }
  },
  {
    $limit: 10
  }
])
```

### Posts Populares (Mais compartilhados hoje)
```javascript
db.posts.find({
  created_at: {
    $gte: new Date(new Date().setHours(0, 0, 0, 0))
  }
}).sort({ created_at: -1 }).limit(10)
```

---

## 🔧 Manutenção

### Backup
```bash
# Backup completo do banco
mongodump --uri="mongodb://localhost:27017/vicio_api" --out=./backup

# Backup de coleção específica
mongodump --uri="mongodb://localhost:27017/vicio_api" --collection posts

# Restore
mongorestore --uri="mongodb://localhost:27017/vicio_api" ./backup/vicio_api
```

### Validação de Índices
```javascript
// Verificar índices
db.posts.getIndexes()

// Remover índice não utilizado
db.posts.dropIndex("index_name")

// Recriar índices
db.posts.reIndex()
```

### Limpeza de Dados Redundantes
```javascript
// Remover posts que referem tags que não existem mais
var invalidPosts = [];
db.posts.find().forEach(function(post) {
  if (!db.tags.findOne({ _id: post.tag_id })) {
    invalidPosts.push(post._id);
  }
});

// Deletar posts inválidos
db.posts.deleteMany({ _id: { $in: invalidPosts } })
```

---

## 🔒 Segurança

### Proteção no MongoDB

```javascript
// Criar usuário com permissões limitadas (recomendado para produção)
db.createUser({
  user: "api_user",
  pwd: "senha_segura",
  roles: [
    {
      role: "readWrite",
      db: "vicio_api"
    }
  ]
})

// Usar em MONGO_URI
mongodb://api_user:senha_segura@localhost:27017/vicio_api?authSource=admin
```

### Validação de Entrada
- Links sempre validados antes de insert
- IP sempre capturado
- Tags referenciadas validadas

---

## 📊 Queries Úteis para Development

### Resetar Banco
```javascript
// Cuidado: Deleta tudo!
db.dropDatabase()
```

### Popular com Dados de Teste
```javascript
// Criar tags
db.tags.insertMany([
  { name: "Programação", slug: "programacao", emoji: "💻", created_at: new Date() },
  { name: "Jogos", slug: "jogos", emoji: "🎮", created_at: new Date() }
])

// Criar posts de teste
db.posts.insertOne({
  title: "Test Post",
  caption: "Test",
  link: "https://www.youtube.com/watch?v=jNQXAC9IVRw",
  type: "youtube",
  tag_id: ObjectId("..."),
  tag_name: "Programação",
  tag_slug: "programacao",
  tag_emoji: "💻",
  ip: "127.0.0.1",
  created_at: new Date()
})
```

### Verificar Tamanho da Coleção
```javascript
db.posts.stats()
db.tags.stats()

// Em bytes
db.posts.stats().size

// Número de documentos
db.posts.countDocuments({})
```

---

## 🛠️ Troubleshooting

### Problema: Query Lenta

**Solução**:
1. Adicionar índice na query (`explain("executionStats")`)
2. Verificar se índice existe
3. Recriar índice se necessário

### Problema: Espaço Em Disco

**Solução**:
```javascript
// Deletar posts antigos
db.posts.deleteMany({
  created_at: {
    $lt: new Date("2023-01-01")
  }
})

// Compactar banco
db.repairDatabase()
```

### Problema: Rate Limit Não Funciona

**Diagnóstico**:
```javascript
// Verificar se há índice em (ip, created_at)
db.posts.getIndexes()

// Criar se não existir
db.posts.createIndex({ "ip": 1, "created_at": 1 })
```

---

## 📚 Recursos

- [MongoDB Manual](https://docs.mongodb.com/manual/)
- [MongoDB Query Language](https://docs.mongodb.com/manual/reference/operator/query/)
- [MongoDB Aggregation](https://docs.mongodb.com/manual/reference/operator/aggregation/)
- [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver)

---

**Database Schema & Operations - v1.0**
