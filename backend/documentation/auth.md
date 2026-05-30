# Auth API Documentation

Base URL: http://localhost:8080/api/v1

Notes:
- Phone normalization: numbers starting with `08` are converted to `62...` (e.g. `081234...` → `628123...`).
- All request/response bodies use JSON.
- Access token must be sent in header: `Authorization: Bearer <ACCESS_TOKEN>`.

---

## Endpoints

### 1) Register — Step 1: Request OTP
- Method: POST
- Path: /auth/register/request-otp
- Description: Kirim OTP ke nomor telepon. Data user (name, phone, password) disimpan sementara hingga OTP diverifikasi.
- Request JSON:
```json
{
  "name": "John Doe",
  "phone": "081234567890",
  "password": "password123",
  "confirm_password": "password123"
}
```
- Success (200) sample:
```json
{
  "success": true,
  "message": "OTP berhasil dikirim",
  "data": {
    "otp": "123456"
  }
}
```
> `otp` hanya muncul di response jika `APP_ENV != production`. Di production, OTP dikirim via WhatsApp.

### 2) Register — Step 2: Resend OTP (opsional)
- Method: POST
- Path: /auth/register/resend-otp
- Description: Kirim ulang OTP jika sudah expired. Membutuhkan data registrasi pending yang masih aktif (dalam 10 menit sejak request-otp). Tidak perlu kirim ulang nama/password.
- Request JSON:
```json
{ "phone": "081234567890" }
```
- Success (200) sample:
```json
{
  "success": true,
  "message": "OTP berhasil dikirim ulang",
  "data": {
    "otp": "654321"
  }
}
```
> `otp` hanya muncul jika `APP_ENV != production`.
> Jika data pending sudah expired (>10 menit sejak request-otp), endpoint ini mengembalikan 404 dan user harus daftar ulang dari awal via `/register/request-otp`.

### 3) Register — Step 3: Verify OTP
- Method: POST
- Path: /auth/register/verify
- Description: Verifikasi OTP untuk menyelesaikan registrasi. Mengembalikan user + token.
- Request JSON:
```json
{
  "phone": "081234567890",
  "otp": "123456"
}
```
- Success (200) sample:
```json
{
  "success": true,
  "message": "registrasi berhasil",
  "data": {
    "user": { "id":"...","name":"John Doe","phone":"6281234567890" },
    "access_token": "<ACCESS_TOKEN>",
    "refresh_token": "<REFRESH_TOKEN>"
  }
}
```

### 4) Login
- Method: POST
- Path: /auth/login
- Request JSON:
```json
{ "phone": "081234567890", "password": "password123" }
```
- Response: same format as Register step 2 (contains user + access + refresh tokens)

### 5) Refresh Token
- Method: POST
- Path: /auth/refresh
- Request JSON:
```json
{ "refresh_token": "<REFRESH_TOKEN>" }
```
- Response: new access & refresh tokens

### 6) Logout
- Method: POST
- Path: /auth/logout
- Auth: Bearer access token
- Body: none
- Effect: blacklists access token

### 7) Me
- Method: GET
- Path: /auth/me
- Auth: Bearer access token
- Response: current user data

### 8) Forgot Password (dev)
- Method: POST
- Path: /auth/forgot-password
- Request:
```json
{ "phone": "081234567890" }
```
- Response (dev): returns OTP in response. In production, OTP sent via WhatsApp.

### 9) Reset Password
- Method: POST
- Path: /auth/reset-password
- Request:
```json
{ "phone":"081234567890", "otp":"123456", "new_password":"newPass123" }
```
- Response: 200 OK when password updated

---

## Postman Tips
1. Create environment variable `base_url` = `http://localhost:8080/api/v1`.
2. To register: POST `{{base_url}}/auth/register/request-otp` with `name`, `phone`, `password`, and `confirm_password` → copy OTP from response → POST `{{base_url}}/auth/register/verify`. If OTP expires before you verify, POST `{{base_url}}/auth/register/resend-otp` with just the phone to get a new OTP.
3. Or use seeded user (admin phone: `6285271766661`, password `9968siskom`).
4. In Tests tab (Login/Register verify) use script to save tokens:
```javascript
let d = pm.response.json().data;
if (d?.access_token) pm.environment.set('access_token', d.access_token);
if (d?.refresh_token) pm.environment.set('refresh_token', d.refresh_token);
```
5. Use `{{base_url}}/auth/me` with header `Authorization: Bearer {{access_token}}` to validate.
6. To refresh, POST to `{{base_url}}/auth/refresh` with JSON `{ "refresh_token": "{{refresh_token}}" }` and update tokens from response.

---

## Frontend / Mobile Integration Notes
- Store `access_token` securely (in-memory or secure storage on mobile). Use refresh token only when access token expired.
- Always send `Authorization: Bearer <ACCESS_TOKEN>` for protected endpoints.
- Normalize phone input on client or server; server normalizes automatically but sending in international form (`62...`) is recommended.
- For "forgot password" flow: call `/forgot-password` to request OTP, prompt user to enter OTP, then call `/reset-password` with OTP and new password.
