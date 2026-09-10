# Weblog

A full-stack weblog application built with **Go**, **Echo**, and **PostgreSQL**.

🌐 **Live Demo:** https://weblog-application.onrender.com/

---

## 📖 About

Weblog is a simple and modern blogging platform where users can create accounts, publish blog posts, interact with posts, and manage their content.

The project was built as a practical backend-focused web application using Go and Echo, with PostgreSQL as the database.

---

## ✨ Features

* User registration and login
* Session-based authentication
* User accounts
* Create blog posts
* Delete blog posts
* Public and private posts
* Image upload for blog posts
* Comment on blog posts
* Share blog posts
* View blog posts
* PostgreSQL database
* Responsive UI with Tailwind CSS

---

## 🛠️ Tech Stack

### Backend

* **Go**
* **Echo**
* **PostgreSQL**
* **HTML Templates**
* **Session-based Authentication**

### Frontend

* **HTML**
* **Tailwind CSS**
* **JavaScript**

### Development & Deployment

* **Docker**
* **Docker Compose**
* **Render**

---

## 📂 Project Structure

```text
.
├── docker-compose.yml
├── go.mod
├── go.sum
│
├── internal
│   ├── config
│   │   └── config.go
│   │
│   ├── database
│   │   └── database.go
│   │
│   ├── handler
│   │   ├── auth_handler.go
│   │   ├── comment_handler.go
│   │   ├── home_handler.go
│   │   ├── template_renderer.go
│   │   ├── weblog_handler.go
│   │   └── weblog_share_handler.go
│   │
│   ├── middleware
│   │   ├── auth.go
│   │   └── guest.go
│   │
│   ├── model
│   │   ├── comment.go
│   │   ├── user.go
│   │   └── weblog.go
│   │
│   ├── repo
│   │   ├── comment_repo.go
│   │   ├── user_repo.go
│   │   ├── weblog_repo.go
│   │   └── weblog_share_repo.go
│   │
│   ├── service
│   │   ├── auth_service.go
│   │   ├── comment_service.go
│   │   ├── weblog_service.go
│   │   └── weblog_share_service.go
│   │
│   ├── session
│   │   └── session.go
│   │
│   └── upload
│       └── image.go
│
├── LICENSE
├── main.go
│
├── migrations
│   └── schema.sql
│
├── README.md
│
└── templates
    ├── auth
    │   ├── login.html
    │   └── register.html
    │
    └── weblog
        ├── create.html
        ├── detail.html
        └── home.html
```

### 📁 Directory Overview

| Directory / File      | Description                                                      |
| --------------------- | ---------------------------------------------------------------- |
| `main.go`             | Application entry point and server initialization                |
| `internal/config`     | Application configuration and environment variables              |
| `internal/database`   | PostgreSQL database connection                                   |
| `internal/handler`    | HTTP handlers for authentication, weblogs, comments, and sharing |
| `internal/middleware` | Authentication and guest middleware                              |
| `internal/model`      | Application data models                                          |
| `internal/repo`       | Database access and repository layer                             |
| `internal/service`    | Application business logic                                       |
| `internal/session`    | User session management                                          |
| `internal/upload`     | Image upload and file handling                                   |
| `migrations`          | PostgreSQL database schema                                       |
| `templates`           | Server-side rendered HTML templates                              |
| `docker-compose.yml`  | Local PostgreSQL development environment                         |

---

## ⚙️ Environment Variables

The application requires the following environment variables:

```env
DATABASE_URL=your_postgresql_connection_string
SESSION_SECRET=your_session_secret
```

For local development, the PostgreSQL database runs through Docker Compose.

The default local database configuration is:

```text
Host: localhost
Port: 5433
Database: web
User: web
Password: web
```

Therefore, the local `DATABASE_URL` can be configured according to the application's database connection format.

---

## 🚀 Running Locally

### Prerequisites

Make sure you have the following installed:

* Go
* Docker
* Docker Compose

### 1. Clone the repository

```bash
git clone https://github.com/NooraPanahi/Weblog-Application
cd Weblog-Application
```

### 2. Start PostgreSQL with Docker Compose

Start the PostgreSQL container:

```bash
docker compose up -d
```

The PostgreSQL server will be available on:

```text
localhost:5433
```

The database configuration is defined in `docker-compose.yml`:

```yaml
POSTGRES_USER: web
POSTGRES_PASSWORD: web
POSTGRES_DB: web
```

### 3. Create the database schema

Execute the SQL schema located at:

```text
migrations/schema.sql
```

The schema creates the required tables.

### 4. Configure environment variables

Set the required environment variables:

```env
DATABASE_URL=your_database_url
SESSION_SECRET=your_session_secret
```

### 5. Install Go dependencies

```bash
go mod download
```

### 6. Run the application

```bash
go run .
```

The application should then be available at:

```text
http://localhost:8080
```

### 7. Stop PostgreSQL

When you are finished, stop the PostgreSQL container:

```bash
docker compose down
```

The PostgreSQL data is persisted using the Docker volume defined in `docker-compose.yml`.

---

## 🔐 Authentication

The application provides user authentication through registration and login.

After successful authentication, a session is created for the user. Protected routes use the session to determine whether a user is authenticated and to associate actions such as creating posts and comments with the correct account.

---

## 📝 Blog Posts

Users can create blog posts containing:

* Title
* Content
* Optional image
* Privacy setting
* Author
* Creation date

Posts can be either:

* **Public** — visible to other users
* **Private** — restricted according to the application's access rules

---

## 💬 Comments & Sharing

Users can interact with blog posts through comments and sharing.

Comments are associated with both the post and the user who created them.

Post sharing is tracked in the `weblog_shares` table, with a unique constraint preventing duplicate shares by the same user.

---

## 🎨 UI

The frontend uses **Tailwind CSS** to provide a clean and responsive interface across the main pages of the application, including authentication, post creation, and blog viewing.

---

### Live Application

🌐 https://weblog-application.onrender.com/

---

## 📚 Tutorial

This project was developed with the help of the following tutorial:

👉 **Tutorial:** [Goraz](https://github.com/DKeshavarz/goraz)

---

## 📄 License

This project is available for educational and personal use.
