# Go-Postgres Stock API

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green.svg)

A robust RESTful API built with **Golang** and **PostgreSQL** to manage stock market data. This project demonstrates full CRUD (Create, Read, Update, Delete) operations using Go's standard `database/sql` library and the `gorilla/mux` router.

## 📂 Project Structure

The project follows a modular architecture, separating data models, routing, and database logic.

```bash
go-postgres/
├── main.go                 # Application entry point; starts server on Port 8000
├── go.mod                  # Module definition and dependencies
├── .env                    # Environment variables (Database credentials)
├── models/
│   └── models.go           # Defines the 'Stock' struct and JSON tagging
├── router/
│   └── router.go           # Defines API routes and links them to middleware
└── middleware/
    └── handler.go          # Contains core logic, database connections, and SQL queries
```

    ## 🛠️ Libraries & Implementation Details

This project utilizes specific libraries to handle routing, configuration, and database connectivity. Here is a breakdown of how they function and how they are implemented in this project:

### 1. `github.com/gorilla/mux` (Routing)
* **How it works:** This library implements a request router and dispatcher for matching incoming requests to their respective handler. It offers more flexibility than Go's standard `http.ServeMux`.
* **My Implementation:**
    * In `router/router.go`, I used `mux.NewRouter()` to initialize the router.
    * I utilized `.Methods("GET", "POST", ...)` to enforce strict HTTP verbs for each endpoint.
    * In `middleware/handler.go`, I used `mux.Vars(r)` to extract dynamic ID parameters from the URL (e.g., extracting `1` from `/api/stock/1`).

### 2. `github.com/lib/pq` (Postgres Driver)
* **How it works:** This is a pure Go driver for PostgreSQL that interfaces with the standard `database/sql` package.
* **My Implementation:**
    * It is imported in `middleware/handler.go` using a blank identifier (`_ "github.com/lib/pq"`).
    * This registers the driver silently as a side-effect, allowing the `sql.Open("postgres", ...)` function to recognize and communicate with the Postgres database without needing to call `pq` functions directly.

### 3. `github.com/joho/godotenv` (Configuration)
* **How it works:** This library loads environment variables from a `.env` file into the system's environment variables, making them accessible via `os.Getenv`.
* **My Implementation:**
    * In the `CreateConnection()` function, `godotenv.Load(".env")` is called immediately.
    * This ensures that sensitive credentials (like `POSTGRES_URL`) are read securely from the file rather than being hardcoded in the Go scripts.

### 4. `database/sql` (Standard Library)
* **How it works:** Provides a generic interface around SQL (or SQL-like) databases. It manages a pool of connections.
* **My Implementation:**
    * **`QueryRow`**: Used for `INSERT` and `SELECT` (single row) operations. I chained `.Scan()` to map the returned data to the `Stock` struct.
    * **`Exec`**: Used for `UPDATE` and `DELETE` operations. I utilized `.RowsAffected()` on the result to confirm how many records were modified.
    * **`Scan`**: Crucial for mapping raw SQL columns (`stock_id`, `name`, etc.) into the Go struct fields (`StockID`, `Name`, etc.).


## 🚀 Getting Started

### Prerequisites
* **Go** (v1.25 or higher)
* **PostgreSQL** installed and running locally or in the cloud.

### 1. Database Setup
Run the following SQL command in your PostgreSQL database to create the necessary table:

```sql
CREATE TABLE stocks (
    stock_id SERIAL PRIMARY KEY,
    name TEXT,
    price DECIMAL,
    company TEXT
);
```

### 2. Environment Configuration
Create a .env file in the root directory of the project:
```env
POSTGRES_URL="postgres://postgres:yourpassword@localhost:5432/yourdbname?sslmode=disable"
```
Replace yourpassword and yourdbname with your actual credentials.

### 3.Installation
Clone the repository and install the required dependencies:
```bash
git clone https://github.com/JangidRkt08/go-postgres.git
cd go-postgres
go mod tidy
```

### 4. Running the Application
Start the server

```bash
go run main.go
```
You should see the output:
    Hello, World! Starting Server on PORT 8000...

## 🔌 API Endpoints

The API runs on `http://localhost:8000`.

| Method   | Endpoint           | Description            | Request Body (JSON Example)                               |
| :------- | :----------------- | :--------------------- | :-------------------------------------------------------- |
| `GET`    | `/api/stocks`      | Get all stocks         | N/A                                                       |
| `GET`    | `/api/stock/{id}`  | Get stock by ID        | N/A                                                       |
| `POST`   | `/api/stock`       | Create a new stock     | `{"name": "Tesla", "price": 250, "company": "Tesla Inc"}` |
| `PUT`    | `/api/stock/{id}`  | Update existing stock  | `{"name": "Tesla", "price": 300, "company": "Tesla Inc"}` |
| `DELETE` | `/api/stock/{id}`  | Delete a stock         | N/A                                                       |


---

### **Analysis of the Code & Libraries (Study Notes)**

Since I Wrote the README based on standard practices, here is a deeper explanation of the **"why"** and **"how"** behind the code you likely wrote:

1.  **The `main` function**:
    * It acts as the orchestrator. It first **connects to the database** (using `sql.Open`), checks the connection (using `db.Ping`), and then **initializes the router** to start listening for requests.

2.  **Handling JSON**:
    * You likely used `encoding/json`.
    * **Decoding**: When a user sends data (POST request), you use `json.NewDecoder(r.Body).Decode(&user)` to turn their raw JSON text into a Go Struct you can work with.
    * **Encoding**: When sending data back, you use `json.NewEncoder(w).Encode(response)` to turn your Go Structs back into JSON text for the browser/client.

3.  **Database Connection String**:
    * The formatting `postgres://user:password@host:port/dbname?sslmode=disable` is crucial. It tells the `lib/pq` driver exactly where to find your Postgres server. The `sslmode=disable` is often used in local development to avoid SSL certificate errors.

    