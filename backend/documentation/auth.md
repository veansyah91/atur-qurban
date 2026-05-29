# Auth API Documentation

Base URL: http://localhost:8080/api/v1

Notes:
- Phone normalization: numbers starting with `08` are converted to `62...` (e.g. `081234...` → `628123...`).
- All request/response bodies use JSON.
- Access token must be sent in header: `Authorization: Bearer <ACCESS_TOKEN>`.

---

## Endpoints

### 1) Register
- Method: POST
- Path: /auth/register
- Description: Create account (phone + password)
- Request JSON:
```json
{
  "name": "John Doe",
  "phone": "081234567890",
  "password": "password123"
}
```
- Success (200) sample:
```json
{
  "success": true,
  "message": "register berhasil",
  "data": {
    "user": { "id":"...","name":"John Doe","phone":"6281234567890" },
    "access_token": "<ACCESS_TOKEN>",
    "refresh_token": "<REFRESH_TOKEN>"
  }
}
```

### 2) Login
- Method: POST
- Path: /auth/login
- Request JSON:
```json
{ "phone": "081234567890", "password": "password123" }
```
- Response: same format as Register (contains access + refresh tokens)

### 3) Refresh Token
- Method: POST
- Path: /auth/refresh
- Request JSON:
```json
{ "refresh_token": "<REFRESH_TOKEN>" }
```
- Response: new access & refresh tokens

### 4) Logout
- Method: POST
- Path: /auth/logout
- Auth: Bearer access token
- Body: none
- Effect: blacklists access token (and optionally refresh token)

### 5) Me
- Method: GET
- Path: /auth/me
- Auth: Bearer access token
- Response: current user data

### 6) Forgot Password (dev)
- Method: POST
- Path: /auth/forgot-password
- Request:
```json
{ "phone": "081234567890" }
```
- Response (dev): returns OTP in response. In production, OTP must be sent via SMS/WhatsApp.

### 7) Reset Password
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
2. Register or use seeded user (admin phone: `6285271766661`, password `9968siskom`).
3. In Tests tab (Login/Register) use script to save tokens:
```javascript
let d = pm.response.json().data;
if (d?.access_token) pm.environment.set('access_token', d.access_token);
if (d?.refresh_token) pm.environment.set('refresh_token', d.refresh_token);
```
4. Use `{{base_url}}/auth/me` with header `Authorization: Bearer {{access_token}}` to validate.
5. To refresh, POST to `{{base_url}}/auth/refresh` with JSON `{ "refresh_token": "{{refresh_token}}" }` and update tokens from response.

---

## Frontend / Mobile Integration Notes
- Store `access_token` securely (in-memory or secure storage on mobile). Use refresh token only when access token expired.
- Always send `Authorization: Bearer <ACCESS_TOKEN>` for protected endpoints.
- Normalize phone input on client or server; server normalizes automatically but sending in international form (`62...`) is recommended.
- For "forgot password" flow: call `/forgot-password` to request OTP, prompt user to enter OTP, then call `/reset-password` with OTP and new password.
