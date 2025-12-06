# External API Documentation

## Overview
External API yang dibuat untuk memungkinkan aplikasi lain mengakses data user melalui API Key authentication.

## Authentication
External API menggunakan **API Key** authentication yang dikonfigurasi melalui environment variable.

### Setup API Key
Tambahkan ke file `.env`:
```bash
API_KEY=your-secret-api-key-here
```

### Cara Menggunakan
Setiap request ke external API harus menyertakan header:
```http
X-API-Key: your-secret-api-key-here
```

## Base URL
```
https://your-domain.com/api/v1/external
```

## Available Endpoints

### 1. Get List Users
Mendapatkan daftar user dengan pagination dan filter.

**Endpoint:** `GET /api/v1/external/users`

**Headers:**
```
X-API-Key: your-secret-api-key-here
Content-Type: application/json
```

**Query Parameters:**
- `page` (int, optional): Halaman saat ini (default: 1)
- `limit` (int, optional): Jumlah data per halaman (default: 20, max: 100)
- `search` (string, optional): Pencarian general
- `employee_id` (string, optional): Filter berdasarkan employee ID (NIP)
- `username` (string, optional): Filter berdasarkan username
- `is_active` (boolean, optional): Filter berdasarkan status aktif
- `start_time` (datetime, optional): Filter data setelah tanggal ini
- `end_time` (datetime, optional): Filter data sebelum tanggal ini

**Example Request:**
```bash
curl -X GET "https://your-domain.com/api/v1/external/users?page=1&limit=10&employee_id=123456" \
  -H "X-API-Key: your-secret-api-key-here"
```

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "uuid-user-1",
      "employee_id": "123456789",
      "username": "john.doe",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@company.com",
      "phone_number": "08123456789",
      "is_active": true,
      "last_login_at": "2023-12-01T10:30:00Z",
      "organization": {
        "id": "uuid-org-1",
        "name": "PT. Example",
        "type": "Headquarters"
      },
      "roles": [
        {
          "id": "uuid-role-1",
          "name": "Admin",
          "description": "Administrator role"
        }
      ],
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-12-01T10:30:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total": 50,
    "total_pages": 5
  },
  "meta": {
    "api_version": "v1",
    "timestamp": "2023-12-01T10:30:00Z"
  }
}
```

### 2. Get User Detail
Mendapatkan detail user berdasarkan UUID.

**Endpoint:** `GET /api/v1/external/users/{id}`

**Headers:**
```
X-API-Key: your-secret-api-key-here
Content-Type: application/json
```

**Path Parameters:**
- `id` (string, required): UUID user

**Example Request:**
```bash
curl -X GET "https://your-domain.com/api/v1/external/users/uuid-user-1" \
  -H "X-API-Key: your-secret-api-key-here"
```

**Success Response (200):**
```json
{
  "data": {
    "id": "uuid-user-1",
    "employee_id": "123456789",
    "username": "john.doe",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@company.com",
    "phone_number": "08123456789",
    "is_active": true,
    "last_login_at": "2023-12-01T10:30:00Z",
    "organization": {
      "id": "uuid-org-1",
      "name": "PT. Example",
      "type": "Headquarters"
    },
    "roles": [
      {
        "id": "uuid-role-1",
        "name": "Admin",
        "description": "Administrator role"
      }
    ],
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-12-01T10:30:00Z"
  },
  "meta": {
    "api_version": "v1",
    "timestamp": "2023-12-01T10:30:00Z"
  }
}
```

### 3. Search Users
Pencarian user dengan filter yang lebih spesifik.

**Endpoint:** `GET /api/v1/external/users/search`

**Headers:**
```
X-API-Key: your-secret-api-key-here
Content-Type: application/json
```

**Query Parameters:**
Sama seperti endpoint list users, dengan tambahan filter spesifik:
- `employee_id` (string): Pencarian exact match employee ID
- `username` (string): Pencarian partial username (case insensitive)
- `is_active` (boolean): Filter status aktif

**Example Request:**
```bash
curl -X GET "https://your-domain.com/api/v1/external/users/search?employee_id=123456789&is_active=true" \
  -H "X-API-Key: your-secret-api-key-here"
```

## Error Responses

### 401 Unauthorized
**Missing or Invalid API Key:**
```json
{
  "message": "API Key is required"
}
```
atau
```json
{
  "message": "Invalid API Key"
}
```

### 400 Bad Request
**Validation Error:**
```json
{
  "error": "Validation error: validation failed"
}
```

**Invalid Query Parameters:**
```json
{
  "error": "Invalid query parameters: invalid parameter format"
}
```

### 404 Not Found
**User Not Found:**
```json
{
  "error": "User not found"
}
```

### 500 Internal Server Error
**Server Error:**
```json
{
  "error": "Failed to fetch users: database connection error"
}
```

## Usage Examples

### 1. Mengambil semua user aktif
```bash
curl -X GET "https://your-domain.com/api/v1/external/users?is_active=true&page=1&limit=50" \
  -H "X-API-Key: your-secret-api-key-here"
```

### 2. Mencari user berdasarkan employee ID (NIP)
```bash
curl -X GET "https://your-domain.com/api/v1/external/users?employee_id=123456789" \
  -H "X-API-Key: your-secret-api-key-here"
```

### 3. Mencari user berdasarkan nama
```bash
curl -X GET "https://your-domain.com/api/v1/external/users/search?username=john" \
  -H "X-API-Key: your-secret-api-key-here"
```

### 4. Filter berdasarkan rentang waktu
```bash
curl -X GET "https://your-domain.com/api/v1/external/users?start_time=2023-01-01T00:00:00Z&end_time=2023-12-31T23:59:59Z" \
  -H "X-API-Key: your-secret-api-key-here"
```

## Security Considerations

1. **API Key Confidentiality**: Jangan pernah share API key di client-side code
2. **HTTPS**: Selalu gunakan HTTPS untuk komunikasi API
3. **Rate Limiting**: Pertimbangkan untuk mengimplementasikan rate limiting
4. **IP Whitelisting**: Batasi akses API hanya dari IP address yang trusted
5. **Key Rotation**: Lakukan rotasi API key secara berkala

## Rate Limiting & Best Practices

1. **Pagination**: Gunakan pagination untuk data yang besar
2. **Filtering**: Gunakan filter yang spesifik untuk mengurangi data transfer
3. **Caching**: Implementasikan caching untuk data yang tidak sering berubah
4. **Error Handling**: Implementasikan proper error handling di client
5. **Retry Logic**: Implementasikan retry logic untuk temporary failures

## Implementation Notes

- **Employee ID**: Field ini bisa digunakan sebagai NIP (Nomor Induk Pegawai)
- **Soft Delete**: Users yang dihapus tidak akan muncul di API response
- **Timestamps**: Semua timestamps menggunakan format ISO 8601 UTC
- **Null Values**: Fields yang kosong akan direpresentasikan sebagai `null` atau dihilangkan
- **Data Types**: Pastikan client menangani tipe data yang sesuai (boolean, datetime, etc.)