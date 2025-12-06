# Employee ID Integration Documentation

## Overview
Berikut adalah ringkasan perubahan yang telah dilakukan untuk menambahkan field `employee_id` sebagai identifier utama user (NIP) pada aplikasi identity system.

## Perubahan yang Dilakukan

### 1. Database Schema
- **Migration**: `000004_add_employee_id_to_users_table.up.sql`
  - Menambahkan kolom `employee_id VARCHAR(50)` pada tabel `users`
  - Menambahkan unique constraint untuk memastikan tidak ada duplicate employee_id
  - Menambahkan index untuk performa query lebih baik

### 2. Entity Layer
- **File**: `internal/entities/user.go`
  - Menambahkan field `EmployeeID nullable.NullString` pada struct `User`

### 3. Repository Layer
- **Files**:
  - `internal/user/repository/user_queries.go` - Update semua SQL queries
  - `internal/user/repository/user.go` - Update fungsi repository
  - `internal/user/repository.go` - Update interface repository

**Fungsi baru yang ditambahkan:**
- `FindByEmployeeID(ctx context.Context, employeeID string) (*entities.User, error)`

**Queries yang diupdate:**
- `insertUser` - menambahkan employee_id parameter
- `findByUsername` - menambahkan employee_id di SELECT
- `findByPhoneNumber` - menambahkan employee_id di SELECT
- `updateUser` - menambahkan update logic untuk employee_id
- `findUserById` - menambahkan employee_id di SELECT
- `listUsers` - menambahkan employee_id di SELECT
- `findByEmployeeId` - query baru untuk mencari berdasarkan employee_id

### 4. Use Case Layer
- **File**: `internal/user/usecase/user.go`
  - Menambahkan validasi uniqueness employee_id pada function `Create`
  - Menambahkan validasi uniqueness employee_id pada function `Update`

### 5. DTO Layer (Request/Response)
- **Files**:
  - `internal/user/dtos/store_user.go` - Update CreateNewUserReq
  - `internal/user/dtos/update_user.go` - Update UpdateUserReq
  - `internal/user/dtos/show_user.go` - Update ShowUserRes
  - `internal/user/dtos/list_user.go` - Update ListUserReq & ListUserRespData
  - `internal/user/dtos/whoami.go` - Update WhoamiRes

**Perubahan validasi:**
- Employee ID required dengan length 1-50 karakter
- Employee ID harus unique (tidak boleh duplikat)

### 6. API Endpoints
Semua existing endpoint otomatis mendukung employee_id tanpa perubahan kode karena menggunakan DTO yang sudah diupdate:
- `POST /api/v1/users` - Create user (dengan employee_id)
- `PUT /api/v1/users/{userUUID}` - Update user (dengan employee_id)
- `GET /api/v1/users/{userId}` - Show user detail (menampilkan employee_id)
- `GET /api/v1/users` - List users (menampilkan & filter employee_id)
- `GET /api/v1/auth/whoami` - Get current user info (menampilkan employee_id)

## Cara Penggunaan

### Create User
```json
{
  "employee_id": "123456789",
  "username": "john.doe",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "08123456789",
  "role_ids": ["role-uuid-1"],
  "organization_id": "org-uuid-1"
}
```

### Update User
```json
{
  "employee_id": "987654321",
  "first_name": "John Updated"
}
```

### Filter Users by Employee ID
```
GET /api/v1/users?employee_id=123456789
```

## Important Notes

1. **Employee ID adalah Unique**: Setiap user harus memiliki employee_id yang unik
2. **Required Field**: Employee ID bersifat required saat create user
3. **Search Support**: Employee ID bisa digunakan untuk filter di list users
4. **NIP Integration**: Employee ID dirancang untuk mengakomodasi NIP (Nomor Induk Pegawai)
5. **Backward Compatible**: Semua fitur existing tetap berjalan normal

## Database Migration
Jalankan migration untuk menambahkan field employee_id:
```bash
# Untuk PostgreSQL
psql -d your_database -f migrations/000004_add_employee_id_to_users_table.up.sql
```

## Testing
Build test sudah dilakukan dan berhasil tanpa error:
```bash
go build -v ./...
go fmt ./...
```

## Next Steps (Frontend)
Perlu update frontend forms untuk:
1. Menambahkan input field Employee ID di create user form
2. Menambahkan input field Employee ID di update user form
3. Menampilkan Employee ID di user list/detail view
4. Menambahkan filter by Employee ID di user list
5. Update validation messages untuk Employee ID