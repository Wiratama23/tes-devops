tech stack used on golang is PGX, CorazaWAF, Goose

#raw postgresql schema
users:
uid : google.uuid
username
email

articles:
articlesid
uid //from users database
title
articleText
dateCreated

products:
product_id : VARCHAR(20) PRIMARY KEY (SKU format)
product_name
product_quantity
product_prices
product_type
created_at
created_by : UUID reference to users.uid
image_path



### Users API

#### 1. Create User
```
POST /users
Content-Type: application/json

Request Body:
{
  "username": "john_doe",
  "email": "john@example.com"
}

Response (201 Created):
{
  "uid": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "email": "john@example.com",
  "created_at": "2026-04-19T10:30:00Z",
  "updated_at": "2026-04-19T10:30:00Z"
}
```

#### 2. Get All Users
```
GET /users

Response (200 OK):
[
  {
    "uid": "550e8400-e29b-41d4-a716-446655440000",
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2026-04-19T10:30:00Z",
    "updated_at": "2026-04-19T10:30:00Z"
  },
  {
    "uid": "660e8400-e29b-41d4-a716-446655440001",
    "username": "jane_smith",
    "email": "jane@example.com",
    "created_at": "2026-04-19T11:00:00Z",
    "updated_at": "2026-04-19T11:00:00Z"
  }
]
```

#### 3. Get User by ID
```
GET /users/{uid}

Example: GET /users/550e8400-e29b-41d4-a716-446655440000

Response (200 OK):
{
  "uid": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "email": "john@example.com",
  "created_at": "2026-04-19T10:30:00Z",
  "updated_at": "2026-04-19T10:30:00Z"
}
```

#### 4. Update User
```
PUT /users/{uid}
Content-Type: application/json

Request Body:
{
  "username": "john_updated",
  "email": "john.new@example.com"
}

Response (200 OK):
{
  "uid": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_updated",
  "email": "john.new@example.com",
  "created_at": "2026-04-19T10:30:00Z",
  "updated_at": "2026-04-19T10:35:00Z"
}
```

#### 5. Delete User
```
DELETE /users/{uid}

Response (204 No Content)
```

---

### Articles API

#### 1. Create Article
```
POST /articles
Content-Type: application/json

Request Body:
{
  "uid": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Getting Started with Go",
  "article_text": "Go is a powerful modern programming language..."
}

Response (201 Created):
{
  "articles_id": 1,
  "uid": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Getting Started with Go",
  "article_text": "Go is a powerful modern programming language...",
  "date_created": "2026-04-19T10:45:00Z",
  "updated_at": "2026-04-19T10:45:00Z"
}
```

#### 2. Get All Articles
```
GET /articles
GET /articles?limit=10
GET /articles?limit=10&offset=0

Query Parameters:
- limit: Number of articles to fetch (default: 50, max recommended: 100)
- offset: Number of articles to skip (default: 0)

Examples:
- GET /articles - Fetch 10 random articles
- GET /articles?page=2 - Fetch 10 articles (default number) for page 1
- GET /articles?page=2&limit=20 - Fetch 20 articles skipping the first 10

Response (200 OK):
{
  "data": [
    {
      "articles_id": 1,
      "uid": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Getting Started with Go",
      "article_text": "Go is a powerful modern programming language...",
      "date_created": "2026-04-19T10:45:00Z",
      "updated_at": "2026-04-19T10:45:00Z"
    },
    {
      "articles_id": 2,
      "uid": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Advanced Go Patterns",
      "article_text": "In this post we'll explore some advanced patterns...",
      "date_created": "2026-04-19T11:15:00Z",
      "updated_at": "2026-04-19T11:15:00Z"
    }
  ],
  "total_count": 42,
  "limit": 10,
  "offset": 0
}
```

#### 3. Get Article by ID
```
GET /articles/{id}

Example: GET /articles/1

Response (200 OK):
{
  "articles_id": 1,
  "uid": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Getting Started with Go",
  "article_text": "Go is a powerful modern programming language...",
  "date_created": "2026-04-19T10:45:00Z",
  "updated_at": "2026-04-19T10:45:00Z"
}
```

#### 4. Get User's Articles
```
GET /users/{uid}/articles

Example: GET /users/550e8400-e29b-41d4-a716-446655440000/articles

Response (200 OK):
[
  {
    "articles_id": 1,
    "uid": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Getting Started with Go",
    "article_text": "Go is a powerful modern programming language...",
    "date_created": "2026-04-19T10:45:00Z",
    "updated_at": "2026-04-19T10:45:00Z"
  }
]
```

#### 5. Update Article
```
PUT /articles/{id}
Content-Type: application/json

Request Body:
{
  "title": "Getting Started with Go (Updated)",
  "article_text": "Updated content..."
}

Response (200 OK):
{
  "articles_id": 1,
  "uid": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Getting Started with Go (Updated)",
  "article_text": "Updated content...",
  "date_created": "2026-04-19T10:45:00Z",
  "updated_at": "2026-04-19T10:50:00Z"
}
```

#### 6. Delete Article
```
DELETE /articles/{id}

Response (204 No Content)
```

---

### Products API

#### 1. Create Product
```
POST /products
Content-Type: application/json

Request Body:
{
  "product_id": "SKU10001",
  "product_name": "Premium Coffee Beans",
  "product_quantity": 100,
  "product_prices": "29.99",
  "product_type": "10",
  "created_by": "550e8400-e29b-41d4-a716-446655440000",
  "image_path": "assets/coffee.jpg"
}

Response (201 Created):
{
  "product_id": "SKU10001",
  "product_name": "Premium Coffee Beans",
  "product_quantity": 100,
  "product_prices": "29.99",
  "product_type": "10",
  "created_at": "2026-04-19T12:00:00Z",
  "created_by": "550e8400-e29b-41d4-a716-446655440000",
  "image_path": "assets/coffee.jpg"
}
```

#### 2. Get All Products
```
GET /products
GET /products?limit=25
GET /products?limit=25&offset=0

Query Parameters:
- limit: Number of products to fetch (default: 100, max recommended: 200)
- offset: Number of products to skip (default: 0)

Examples:
- GET /products - Fetch 10 random products
- GET /products?limit=25 - Fetch 25 random products
- GET /products?limit=20&offset=20 - Fetch 20 products skipping the first 20

Response (200 OK):
{
  "data": [
    {
      "product_id": "SKU10001",
      "product_name": "Premium Coffee Beans",
      "product_quantity": 100,
      "product_prices": "29.99",
      "product_type": "10",
      "created_at": "2026-04-19T12:00:00Z",
      "created_by": "550e8400-e29b-41d4-a716-446655440000",
      "image_path": "assets/coffee.jpg"
    },
    {
      "product_id": "SKU05001",
      "product_name": "Programming Book",
      "product_quantity": 50,
      "product_prices": "49.99",
      "product_type": "05",
      "created_at": "2026-04-19T12:15:00Z",
      "created_by": "550e8400-e29b-41d4-a716-446655440000",
      "image_path": "assets/book.jpg"
    }
  ],
  "total_count": 156,
  "limit": 25,
  "offset": 0
}
```

#### 3. Get Product by ID
```
GET /products/{id}

Example: GET /products/SKU10001

Response (200 OK):
{
  "product_id": "SKU10001",
  "product_name": "Premium Coffee Beans",
  "product_quantity": 100,
  "product_prices": "29.99",
  "product_type": "10",
  "created_at": "2026-04-19T12:00:00Z",
  "created_by": "550e8400-e29b-41d4-a716-446655440000",
  "image_path": "assets/coffee.jpg"
}
```

#### 4. Update Product
```
PUT /products/{id}
Content-Type: application/json

Request Body:
{
  "product_name": "Premium Coffee Beans - Updated",
  "product_quantity": 85,
  "product_prices": "34.99",
  "product_type": "10",
  "image_path": "assets/coffee_updated.jpg"
}

Response (200 OK):
{
  "product_id": "SKU10001",
  "product_name": "Premium Coffee Beans - Updated",
  "product_quantity": 85,
  "product_prices": "34.99",
  "product_type": "10",
  "created_at": "2026-04-19T12:00:00Z",
  "created_by": "550e8400-e29b-41d4-a716-446655440000",
  "image_path": "assets/coffee_updated.jpg"
}
```

#### 5. Delete Product
```
DELETE /products/{id}

Example: DELETE /products/SKU10001

Response (204 No Content)
```

#### Product Type Codes
- `10` - Drinks (e.g., coffee, tea, beverages)
- `05` - Books (e.g., programming books, novels)
- `20` - Electronics
- Other codes can be defined based on your product categories

#### SKU Format
Product IDs follow the format: `SKU-CC-NUM` where:
- `SKU` - Static prefix
- `CC` - Category Code (e.g., 10 for drinks, 05 for books)
- `NUM` - Random 3-digit number

Example: `SKU10001`, `SKU05042`, `SKU20123`

---

## 🔒 Error Responses

### 400 Bad Request
```json
Invalid request body or parameters
```

### 404 Not Found
```json
{
  "error": "user not found"
}

{
  "error": "article not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "database connection failed"
}
```


## 💾 Database Seeding

To populate the database with sample data:

```bash
# Run the seed script
go run cmd/seed/main.go

# This creates:
# - 2 sample users (alice, bob)
# - 3 sample articles (2 by alice, 1 by bob)
# - 3 sample products (2 by alice, 1 by bob)
```

### Manual Testing with cURL

```bash
# Create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com"}'

# Get all users
curl http://localhost:8080/users

# Get specific user
curl http://localhost:8080/users/{uid}

# Update user
curl -X PUT http://localhost:8080/users/{uid} \
  -H "Content-Type: application/json" \
  -d '{"username":"updated","email":"new@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/users/{uid}

# Create article
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -d '{
    "uid":"{user-uid}",
    "title":"My Article",
    "article_text":"Article content here"
  }'

# Get all articles
curl http://localhost:8080/articles

# Get user's articles
curl http://localhost:8080/users/{uid}/articles

# Create product
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "product_id":"SKU10001",
    "product_name":"Premium Coffee",
    "product_quantity":100,
    "product_prices":"29.99",
    "product_type":"10",
    "created_by":"{user-uid}",
    "image_path":"assets/coffee.jpg"
  }'

# Get all products
curl http://localhost:8080/products

# Get specific product
curl http://localhost:8080/products/SKU10001

# Update product
curl -X PUT http://localhost:8080/products/SKU10001 \
  -H "Content-Type: application/json" \
  -d '{
    "product_name":"Premium Coffee - Updated",
    "product_quantity":85,
    "product_prices":"34.99",
    "product_type":"10",
    "image_path":"assets/coffee_updated.jpg"
  }'

# Delete product
curl -X DELETE http://localhost:8080/products/SKU10001
```


---

## 🔧 Tech Stack

- **Language**: Go 1.26.2
- **Database Driver**: PGX v5
- **Database**: PostgreSQL 15
- **Migrations**: Goose v3
- **WAF**: Coraza v3 with OWASP CRS
- **UUID**: google/uuid
- **Environment**: godotenv

---

## 📁 Project Structure

```
/tesdb
├── cmd/
│   ├── api/main.go              # API entry point
│   └── seed/main.go             # Database seeder
├── internal/
│   ├── database/db.go           # Database connection
│   ├── handlers/                # HTTP handlers
│   ├── middleware/              # WAF middleware
│   ├── models/                  # Data models
│   └── repository/              # Database operations
├── migrations/                  # Database migrations (Goose)
├── .env                         # Environment variables
├── docker-compose.yaml          # Docker setup
├── Dockerfile                   # Container build
├── go.mod                       # Go modules
├── go.sum                       # Dependency checksums
└── readme.md                    # This file
```

cmd:
run the seeder
sudo docker compose exec api ./seeder