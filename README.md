# 🏢 Leave Management System - Backend (Golang)

This repository contains the **backend service** for the **Leave Management System**, built using **Go (Golang)**.  
It provides RESTful APIs for employees and HR to manage leave applications efficiently.

---

## 🚀 Tech Stack

- **Go (Golang)** — Backend language  
- **Gorilla Mux** — HTTP routing  
- **MongoDB Atlas** — Cloud database  
- **JWT Authentication** — Secure HR login  
- **dotenv** — Environment configuration  
- **CORS** — Frontend-backend communication  

---

## ⚙️ Features

### 👩‍💼 Employee
- Apply for leave  
- View leave status  
- Retrieve profile details  

### 🧍‍♀️ HR
- Login securely with JWT  
- View pending leave requests  
- Approve or decline leaves  

---

## 📂 Folder Structure

```bash
backend/
│
├── main.go              # Entry point for the server
├── handlers.go          # Contains all route handlers
├── db.go                # MongoDB connection setup
├── models.go            # Structs and data models
├── middleware.go        # JWT middleware for HR authentication
├── go.mod / go.sum      # Go module files
└── .env                 # Environment variables (not committed)

```

## 🔑 Environment Variables (`.env`)

Create a `.env` file inside the backend directory with these values:

```env
MONGO_URI=your_mongo_connection_string
DB_NAME=leave_management
HR_PASSWORD=admin
JWT_SECRET=your_secret_key
```
## 🧩 API Endpoints

### 👨‍💼 Employee Routes

| Method | Endpoint | Description |
|--------|-----------|-------------|
| `GET` | `/employee/{id}` | Get employee details |
| `POST` | `/employee/{id}/apply-leave` | Apply for leave |
| `GET` | `/employee/{id}/leaves` | Get all leaves applied by employee |



### 🧍‍♀️ HR Routes

| Method | Endpoint | Description |
|--------|-----------|-------------|
| `POST` | `/hr/login` | HR login — returns JWT token |
| `GET` | `/hr/pending-leaves` | View all pending leave requests |
| `POST` | `/hr/leave/{index}/grant` | Approve leave |
| `POST` | `/hr/leave/{index}/decline` | Decline leave |

---

## 🔒 Authentication

All HR routes require a valid JWT token passed in the request header:

Authorization: Bearer <your_token>

yaml
Copy code

You can obtain the token by logging in through the `/hr/login` endpoint.

---

## 👥 Test Data

### 🧑‍💻 Employees

| Name   | ID  |
|--------|-----|
| Nithin | 22  |
| Kalyan | 23  |

### 👩‍💼 HR Login

| Role | Password |
|------|-----------|
| HR   | admin     |

---

## 🧾 Example API Requests

### 🔑 HR Login

**Request**

```bash
POST /hr/login
Content-Type: application/json

{
  "password": "admin"
}
Response

json
Copy code
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5..."
}
```
📝 Apply for Leave
Request

```bash
Copy code
POST /employee/22/apply-leave
Content-Type: application/json

{
  "reason": "Family event",
  "from_date": "2025-11-10",
  "to_date": "2025-11-12"
}
Response

json
Copy code
{
  "message": "Leave request submitted successfully"
}
```

