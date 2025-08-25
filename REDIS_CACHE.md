# Redis Cache Implementation

Implementasi cache Redis yang sederhana untuk aplikasi Go Event.

## 📋 Overview

Cache Redis ini dirancang dengan prinsip **KISS (Keep It Simple, Stupid)** - hanya 3 method utama yang mudah dipahami dan digunakan.

### 🎯 Fitur Utama
- **Simple**: Hanya 3 method (Set, Get, Delete)
- **Automatic**: Cache otomatis terintegrasi dengan database models
- **Graceful**: Aplikasi tetap berjalan jika Redis tidak tersedia
- **Fast**: Response time lebih cepat untuk data yang sering diakses

## 🛠️ Setup

### 1. Install Redis
```bash
# Windows
choco install redis-64

# macOS
brew install redis

# Linux (Ubuntu/Debian)
sudo apt update
sudo apt install redis-server
```

### 2. Start Redis Server
```bash
# Windows/macOS/Linux
redis-server

# Linux (service)
sudo systemctl start redis-server
```

### 3. Test Redis Connection
```bash
redis-cli ping
# Should return: PONG
```

### 4. Set Environment Variables
Buat file `.env` di root project:
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=event_app

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key
```

### 5. Install Dependencies & Run
```bash
go mod tidy
go run cmd/api/main.go
```

## 📁 File Structure

### Core Implementation
- **`internal/cache/cache.go`** - Implementasi cache Redis sederhana (hanya 3 method)

### Environment Configuration  
- **`internal/env/redis.go`** - Konfigurasi Redis dari environment variables

### Database Integration
- **`internal/database/models.go`** - Inisialisasi cache di models
- **`internal/database/users.go`** - Cache integration untuk user operations
- **`internal/database/events.go`** - Cache integration untuk event operations

### Dependencies
- **`go.mod`** - Menambahkan dependency `github.com/redis/go-redis/v9`

## 🔧 Cara Kerja

### Database Caching (Otomatis)
Cache sudah terintegrasi dengan model database. Tidak perlu modifikasi kode tambahan:

```go
// Get user dari cache atau database
user, err := models.Users.Get(userID)

// Get event dari cache atau database  
event, err := models.Events.Get(eventID)

// Get all events dari cache atau database
events, err := models.Events.GetAll()
```

### Cache Methods (Hanya 3)
```go
// Set - Menyimpan data ke cache
cache.Set(ctx, "key", data, 30*time.Minute)

// Get - Mengambil data dari cache  
cache.Get(ctx, "key", &data)

// Delete - Menghapus data dari cache
cache.Delete(ctx, "key1", "key2")
```

### Cache Keys
- `user:{id}` - User berdasarkan ID
- `user:email:{email}` - User berdasarkan email
- `event:{id}` - Event berdasarkan ID
- `events:list` - Daftar semua events

### TTL (Time To Live)
- **User Cache**: 30 menit
- **Event Cache**: 30 menit
- **Events List Cache**: 15 menit

### Cache Invalidation
Cache otomatis di-invalidate saat:
- **User Insert**: Invalidate events list
- **Event Insert**: Invalidate events list dan user cache
- **Event Update**: Invalidate event cache dan events list
- **Event Delete**: Invalidate event cache dan events list

## 💻 Manual Cache Operations

Jika ingin menggunakan cache secara manual:

```go
// Set cache
err := models.Cache.Set(ctx, "custom:key", data, 10*time.Minute)

// Get cache
var data MyStruct
err := models.Cache.Get(ctx, "custom:key", &data)

// Delete cache
err := models.Cache.Delete(ctx, "custom:key")
```

## 🧪 Testing & Monitoring

### Test Redis Connection
```bash
redis-cli ping
# Should return: PONG
```

### Monitor Cache dengan Redis CLI
```bash
# Connect ke Redis
redis-cli

# Monitor semua Redis commands
MONITOR

# Lihat semua keys
KEYS *

# Get specific keys
GET "user:1"
GET "event:1"
GET "events:list"

# Check memory usage
INFO memory

# Check database size
DBSIZE
```

### Test API Endpoints
```bash
# Get all events (akan di-cache)
curl http://localhost:8080/api/v1/events

# Get specific event (akan di-cache)
curl http://localhost:8080/api/v1/events/1

# Create new event (akan invalidate cache)
curl -X POST http://localhost:8080/api/v1/events \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Event","description":"Test Description","date":"2024-01-01","location":"Test Location"}'
```

## ✅ Benefits

1. **Performance**: Response time lebih cepat untuk data yang sering diakses
2. **Scalability**: Mengurangi beban database
3. **Reliability**: Graceful degradation jika Redis tidak tersedia
4. **Simplicity**: Implementasi yang sederhana dan mudah dipahami
5. **Maintainability**: Kode yang bersih dan mudah maintain

## 🚨 Troubleshooting

### Redis Connection Failed
```
Warning: Failed to connect to Redis: connection refused. Continuing without cache.
```

**Solutions:**
- Pastikan Redis server berjalan
- Cek port Redis (default: 6379)
- Cek firewall settings
- Cek Redis configuration

### Cache Not Working
**Check:**
- Redis connection status di log aplikasi
- Cache keys di Redis CLI (`KEYS *`)
- TTL settings
- Data serialization

### High Memory Usage
**Solutions:**
- Monitor Redis memory usage (`INFO memory`)
- Adjust TTL settings
- Implement cache eviction policy
- Consider Redis cluster untuk scaling

## 🏭 Production Setup

### Redis Configuration
Edit `/etc/redis/redis.conf`:
```conf
# Memory management
maxmemory 256mb
maxmemory-policy allkeys-lru

# Persistence
save 900 1
save 300 10
save 60 10000

# Security
requirepass your-strong-password
```

### Environment Variables
```env
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-strong-password
REDIS_DB=0
```

### Monitoring
- Setup Redis monitoring dengan Redis INFO
- Monitor cache hit/miss ratio
- Setup alerts untuk Redis memory usage
- Consider Redis Sentinel untuk high availability

## 📈 Performance Tips

1. **TTL Optimization**: Sesuaikan TTL berdasarkan frekuensi update data
2. **Cache Keys**: Gunakan prefix yang konsisten untuk easy management
3. **Memory Management**: Monitor dan limit Redis memory usage
4. **Connection Pooling**: Redis client sudah menggunakan connection pooling
5. **Serialization**: Gunakan efficient serialization format

## 🔒 Security Considerations

1. **Redis Authentication**: Enable Redis password di production
2. **Network Security**: Restrict Redis access dengan firewall
3. **Data Encryption**: Consider Redis encryption untuk sensitive data
4. **Access Control**: Limit Redis access hanya untuk aplikasi yang membutuhkan

---

## 🎉 Summary

Implementasi cache Redis ini sangat sederhana namun powerful:

- **Hanya 3 method utama** (Set, Get, Delete)
- **Cache otomatis terintegrasi** dengan database models
- **Graceful degradation** jika Redis tidak tersedia
- **Response time lebih cepat** untuk data yang sering diakses
- **Mudah dipahami dan maintain**

Cache akan otomatis bekerja tanpa perlu konfigurasi tambahan! 🚀
