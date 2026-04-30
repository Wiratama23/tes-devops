# Troubleshooting Guide

## Issue: ERR_CONNECTION_REFUSED on Image Loading

### Root Cause
The `CONNECTION_REFUSED` error occurs when the frontend tries to fetch images/data from the backend API, but the backend service (`api` container) is not running or not accessible.

### Symptoms
- Browser console shows: `Failed to load resource: net::ERR_CONNECTION_REFUSED`
- Error URLs: `localhost/api/auth`, `localhost/api/uploads/images`, etc.
- Images don't load, showing the fallback default image instead
- Network tab shows failed requests to API endpoints

### Why It Happens

1. **Backend not running**: The Docker backend services aren't started
2. **Network mismatch**: Frontend is running on host but backend is in Docker
3. **DNS resolution**: Docker DNS resolver not resolving service names correctly
4. **Port not exposed**: API port 8080 not accessible from frontend

### Solutions

#### Solution 1: Start All Services with Docker Compose (Recommended for Production)

```powershell
# From the root directory
docker-compose up -d

# Verify all services are running
docker-compose ps

# Check backend logs
docker-compose logs api

# Check if backend is healthy
curl http://localhost/api/health
```

**Why this works**: All services run in the same Docker network, so `http://api:8080` DNS resolution works correctly through the Docker internal resolver at `127.0.0.11`.

#### Solution 2: Development Setup (Frontend on host, Backend in Docker)

```powershell
# Start only backend services
docker-compose -f docker-compose.dev.yaml up -d

# In a separate terminal, start frontend
cd frontend
bun run dev
```

**Configuration for this setup** (already in `.env.local`):
```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost/api
INTERNAL_API_URL=http://localhost/api
```

**Why it works**: nginx (running in Docker) proxies requests from `localhost/api` to the backend container.

#### Solution 3: Ensure Backend Connectivity

```powershell
# Test if backend is responding
curl http://localhost/api/health

# If using Docker, check container network
docker inspect <container-id> | grep -A 10 Networks

# Verify DNS resolution in Docker
docker run --rm -it nicolaka/netshoot nslookup api
```

### CORS Configuration

Your backend is correctly configured with CORS support. The configuration file is in `backend/cmd/api/main.go`:

```go
router.Use(cors.Handler(cors.Options{
    AllowedOrigins: utils.GetAllowedOrigins(), // From ALLOWED_ORIGINS env var
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    AllowCredentials: true,
}))
```

**To add custom origins**, set the `ALLOWED_ORIGINS` environment variable:

```bash
# In docker-compose.yaml or .env
ALLOWED_ORIGINS=http://localhost:4000,http://localhost:3000,https://yourdomain.com
```

---

## Issue: Images Not Using Default Fallback

### Root Cause
The `SmartImage` component has built-in error handling, but there are several reasons it might not work:

1. **Default image file missing**: `/public/assets/default_image.jpg` not found
2. **Environment variable not set**: `NEXT_PUBLIC_DEFAULT_IMAGE` not configured
3. **Image error not triggering**: Backend returns 500 instead of proper 404
4. **Caching issues**: Stale image in browser cache

### Solutions

#### Solution 1: Verify Default Image Exists

```powershell
# Check if the default image is in the frontend
Test-Path "d:\python coding\mitesdevops\frontend\public\assets\default_image.jpg"

# Check if it's in the backend
Test-Path "d:\python coding\mitesdevops\assets\default_image.jpg"
```

#### Solution 2: Verify Environment Variable

Check your `.env.local`:
```bash
NEXT_PUBLIC_DEFAULT_IMAGE=/assets/default_image.jpg
```

If missing, add it.

#### Solution 3: Clear Browser Cache

```powershell
# Hard refresh in browser: Ctrl+Shift+R (Windows)
# Or clear application cache in DevTools
```

#### Solution 4: Verify Backend Image Serving

The backend should serve images at: `/api/assets/{filename}`

Test it:
```powershell
# Test with a known product image
curl http://localhost/api/assets/default_image.jpg -v

# If using Docker
docker-compose exec api wget -qO- http://localhost:8080/health
```

### Image Flow

1. **Product has image_path** (e.g., `"assets/coffee.jpg"`)
   - Frontend resolves to: `http://localhost/api/assets/coffee.jpg`
   - Backend returns the image or 404
   - On error, falls back to `/assets/default_image.jpg`

2. **Product has no image_path** (null/empty)
   - Frontend skips API call entirely
   - Uses fallback: `/assets/default_image.jpg` immediately

3. **API unreachable** (CONNECTION_REFUSED)
   - Image fails to load
   - Falls back to `/assets/default_image.jpg`

---

## Checklist for Troubleshooting

- [ ] Backend service is running: `docker-compose ps` shows all services UP
- [ ] Backend is healthy: `curl http://localhost/api/health` returns 200
- [ ] Nginx is routing correctly: `curl http://localhost/api/products?page=1` works
- [ ] CORS headers present: Check browser DevTools Network tab for `Access-Control-Allow-*` headers
- [ ] Default image exists: `/public/assets/default_image.jpg` is present
- [ ] Environment variables set: Check `.env.local` for `NEXT_PUBLIC_DEFAULT_IMAGE`
- [ ] No browser cache issues: Hard refresh with Ctrl+Shift+R
- [ ] Firewall/Proxy not blocking: Corporate proxy might block localhost connections

---

## Quick Commands

```powershell
# Full reset and restart
docker-compose down
docker-compose up --build

# Check logs in real-time
docker-compose logs -f api

# Test backend directly
curl -v http://localhost/api/health

# Test frontend
bun run dev

# Build and run in Docker
bun run build
docker-compose -f docker-compose.test.yaml up --build
```

