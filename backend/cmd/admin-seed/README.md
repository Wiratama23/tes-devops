# Admin Seeder

A standalone command-line tool to create admin accounts in the database.

## Usage

### Build

```bash
cd backend
go build -o admin-seed ./cmd/admin-seed
```

### Run - Interactive Mode (Recommended)

Simply run the command and you'll be prompted for credentials:

```bash
./admin-seed
```

You'll be prompted for:
1. Admin username
2. Admin email
3. Admin password
4. Password confirmation

Example output:
```
Enter admin username: admin2
Enter admin email: admin2@example.com
Enter admin password: 
Confirm admin password: 

✅ Admin user created successfully!
   UID:      550e8400-e29b-41d4-a716-446655440000
   Username: admin2
   Email:    admin2@example.com
   Is Admin: true
```

### Run - Command-Line Arguments (Non-interactive)

Pass credentials directly via flags:

```bash
./admin-seed -username=admin2 -email=admin2@example.com -password=SecurePassword123
```

Flags:
- `-username string` - Admin username
- `-email string` - Admin email address
- `-password string` - Admin password (will be prompted if not provided)

## Environment Variables

The tool requires a database connection. Configure via `.env` file or environment:

```properties
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

## Docker Usage

If running inside Docker:

```bash
docker exec -it <backend-container-id> go run ./cmd/admin-seed/main.go
```

Or with arguments:

```bash
docker exec -it <backend-container-id> go run ./cmd/admin-seed/main.go \
  -username=admin3 \
  -email=admin3@example.com \
  -password=MySecurePassword
```

## Duplicate Prevention

If a username already exists, the operation will fail with an error. You must use a unique username for each admin account.

## Password Requirements

- Must not be empty
- Minimum recommended length: 8 characters
- Stored as bcrypt hash (salted and secured)
- Passwords must match when confirming (interactive mode)

## Related Commands

- `seed` - Create test data (users, articles, products)
- `api` - Main API server
