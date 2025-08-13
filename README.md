# 🐊 Gator Project

Gator is a backend project built in **Go (Golang)**.
It provides a simple RSS feed aggregation system with user authentication and feed management features.  

---

## ✨ Features / Commands

### 👤 User Management
- `register` – 📝 Register a new user.
    ```
        Gator register <username>
    ```
- `login` – 🔑 Login as an existing user.
    ```
        Gator login <username>
    ```
- `logout` – 🚪 Logout the current user.
    ```
        Gator logout
    ```
- `users` – 👥 Display all user names in the database.
    ```
        Gator users
    ```

### 📰 Feed Management
- `addfeed` – ➕ Add a feed with a name and URL (requires login).
    ```
        Gator addfeed <name> <url>
    ```
- `feeds` – 📃 List all feeds in the system (requires login).
    ```
        Gator feeds
    ```
- `follow` – ⭐ Follow a feed created by another user (requires login).
   ```
        Gator follow <url>
    ```
- `unfollow` – ❌ Unfollow a feed or remove a feed created by the current user (requires login).
   ```
        Gator unfollow <url>
    ```
- `following` – 👀 Show all feeds the current user is following (requires login).
   ```
        Gator following
    ```

### ⚡ Aggregation & Browsing
- `agg` – 🔄 Main background command. Continuously fetches posts from all feeds, using `timebetween` to control the interval between requests (requires login).
   ```
        Gator agg <time_between_requests example(1h, 1m, 1s)>
    ```
- `browse` – 📖 Browse posts from feeds followed by the current user (requires login).
   ```
        Gator browse
    ```

### 🛠 Maintenance
- `clear` – 🧹 Completely clears all posts for a fresh start (requires login).
   ```
        Gator clear
    ```
- `reset` – 🔁 Resets the entire database, removing all users, feeds, and posts, and logs out the current user. Resets the program state.
   ```
        Gator reset
    ```

---

## 🧩 Dependencies

- Go 1.24.5 🟢
- PostgreSQL 🐘
- Goose – database migration tool 🛠
- SQLC – type-safe SQL query generator 🔧

---

## 🚀 Setup & Usage

1. **Create config file**  

Create a file at `~/.gatorconfig.json` with the following structure:  

```
{
  "db_url": "postgres://<postgresUsername>:<postgresPassword>@localhost:5432/gator?sslmode=disable"
}
```

2. **Create PostgreSQL database**  

```
CREATE DATABASE gator;
```

3. **Run migrations**  

Navigate to `Gator/sql/schema` and run:  

```
goose <postgres_username> <postgres_connection_string> up
```

4. **Start the project**  

From the root of the project directory:  

```
go run .
```

5. **Build the project**  

```
go build
```

---

Gator is now ready to use! 🐊  
Start by registering a user, adding feeds, and following feeds to begin aggregating content. 📈

