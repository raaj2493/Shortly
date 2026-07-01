````md
# 🚀 Shortly

> **A production-inspired full-stack URL Shortener built with Go, PostgreSQL, Redis, and Flutter.**

Build, manage, and analyze short links while learning modern backend engineering, system design, and production-ready development practices.

---

## ✨ Features

- 🔐 JWT Authentication & Authorization
- 🔗 URL Shortening with Custom Aliases
- 📊 Click Analytics Dashboard
- ⚡ Redis Caching
- 🔍 Search, Filtering & Pagination
- ⏰ URL Expiration
- 🧪 Unit & Integration Testing
- 🐳 Dockerized Development
- 📱 Flutter Mobile Application
- 🚀 CI/CD & Production Deployment

---

## 🛠 Tech Stack

### Backend

- Go
- net/http
- PostgreSQL
- pgxpool
- Redis
- JWT
- bcrypt
- Docker
- Docker Compose
- golang-migrate
- Swagger/OpenAPI

### Frontend

- Flutter
- Riverpod
- Dio
- GoRouter
- flutter_secure_storage
- Material 3

### DevOps

- Docker
- GitHub Actions
- Nginx
- Ubuntu VPS
- Let's Encrypt

---

## 📂 Project Structure

```text
shortly/

├── backend/
│   ├── cmd/
│   ├── internal/
│   ├── migrations/
│   ├── tests/
│   └── docker/
│
└── flutter_app/
    ├── lib/
    ├── assets/
    └── test/
````

---

## 🏗 Architecture

```text
Flutter App
      │
      ▼
 REST API
      │
      ▼
 Middleware
      │
      ▼
 Handlers
      │
      ▼
 Services
      │
      ▼
 Repository
      │
      ├──────────────► Redis
      │
      ▼
 PostgreSQL
```

---

## 🚀 Getting Started

### Clone the repository

```bash
git clone https://github.com/your-username/shortly.git
cd shortly
```

### Start the development environment

```bash
docker compose up --build
```

---

## 📚 What You'll Learn

* Clean Architecture
* REST API Design
* Authentication & Authorization
* PostgreSQL
* Redis Caching
* Docker
* Background Workers
* Testing
* CI/CD
* Flutter API Integration
* Practical System Design

---

## 🎯 Project Goals

This project focuses on building software the way professional engineering teams do.

* Design before coding
* Build one feature at a time
* Understand every line of code
* Write clean, maintainable software
* Apply production-inspired engineering practices

---

## 🧠 AI Development Philosophy

AI is used as a **mentor**, **reviewer**, and **pair programmer**—not as an autopilot.

Every feature follows the same workflow:

```text
Understand
    ↓
Design
    ↓
Implement
    ↓
Review
    ↓
Test
    ↓
Improve
```

The goal is to become an **AI-augmented engineer**, not a vibe coder.

---

## 🗺 Roadmap

* ✅ Project Foundation
* ✅ Authentication
* ✅ URL Management
* ✅ Analytics
* ✅ Redis Caching
* ✅ Background Workers
* ✅ Testing
* ✅ Flutter Application
* ✅ Deployment & CI/CD

---

## 🌟 Why Shortly?

Unlike a typical CRUD project, Shortly emphasizes:

* Production-inspired architecture
* Real-world backend engineering
* Practical system design
* Modern Flutter development
* Deployment & DevOps fundamentals
* Clean, maintainable code

---

## 🤝 Contributing

Contributions, ideas, and suggestions are always welcome.

Feel free to fork the project, open an issue, or submit a pull request.

---

## 📄 License

This project is licensed under the MIT License.

---

## 👨‍💻 Author

**Ashutosh Giri**

Building software from first principles while learning modern backend engineering.

If you found this project helpful, consider giving it a ⭐ on GitHub!

```
```
