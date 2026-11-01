# 🏷️ PAASET — Backend API

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Framework-008ECF?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)

Backend service for **PAASET**, an online catalog platform to showcase and manage consigned assets such as cars, trucks, heavy equipment, and properties.

---

## 📖 About PAASET

**PAASET** helps connect asset owners with buyers through a clean online catalog managed by an Admin team:

- 🛒 **For Visitors (Public Catalog)**: Browse available units, search by keyword or location, filter by category and price range, sort listings (newest, lowest/highest price, most popular), and view detailed specifications and photos before contacting Admin via WhatsApp.
- 🔐 **For Admins (Dashboard)**: Manage categories and custom specifications (e.g., Year, Transmission, Land Size), publish or update asset listings, upload unit photos, and store **private owner contact details & internal notes** safely (hidden from public visitors).

---

## ✨ Main Features

- 🗂️ **Flexible Categories & Specifications** — Create any asset category and define custom specification fields for each category.
- 📸 **Asset & Photo Management** — Manage asset details, pricing, unit status (`ACTIVE`, `SOLD`, etc.), and multiple photos.
- 🛡️ **Owner Privacy Protection** — Owner names, phone numbers, and internal workshop notes are visible only to logged-in Admins.
- ⚙️ **Website Settings** — Manage site name and Admin WhatsApp contact number dynamically.

---

## 🧰 Prerequisites

- 🐳 **Docker & Docker Compose** *(Recommended)*
- 💻 Or **Go 1.26+**, **PostgreSQL**, and **Redis** *(if running manually without Docker)*

---

## 🚀 How to Run (Quick Start with Docker)

### 1️⃣ Prepare the `.env` file
Copy `.env.example` into `.env` and fill in your configuration values:
```bash
cp .env.example .env
```

### 2️⃣ Start the application
```bash
docker compose up -d --build
```
This starts PostgreSQL, Redis, and the Backend API, and automatically creates all database tables.

### 3️⃣ Open API Documentation
- 🌐 **API Base URL**: `http://localhost:8080/api/v1`
- 📚 **Swagger UI (Interactive Docs)**: `http://localhost:8080/swagger`

---

## 🛠️ Running Without Docker (Local Development)

1. Make sure PostgreSQL and Redis are running on your machine, and update `DB_URL` & `REDIS_ADDR` in `.env`.
2. Run the server:
   ```bash
   go mod download
   go run ./cmd/tjus/main.go
   ```

---

## 🧪 Running Tests

```bash
go test -v ./...
```
