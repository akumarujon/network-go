# network-go

A scalable social networking platform built with Go and GORM, enabling users to
share posts and interact, featuring efficient data modeling, RESTful APIs, and
optimized database queries for enhanced performance and user experience.

## Features

- User registration with email confirmation and password hashing
- User authentication
- Post creation and sharing
- User interactions
- Efficient data modeling using GORM
- RESTful API endpoints
- Optimized database queries for improved performance

## Technologies Used

- Go
- GORM (Go Object Relational Mapper)
- RESTful API architecture
- PostgreSQL

## Getting Started

### Prerequisites

- Go
- PostgreSQL

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/akumarujon/network-go.git
   ```

2. Navigate to the project directory:
   ```bash
   cd network-go
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Set up your PostgreSQL database and update the connection string in the
   configuration file.

5. Run the application:
   ```bash
   go run main.go
   ```

## Usage

The network-go application provides the following API endpoints:

### Authentication

- `/signin` - Sign in to an existing account
- `/signup` - Create a new account
- `/confirm/:token` - Confirm a user account, where `:token` is a UUID

### Posts

- `GET /` - Retrieve all posts
- `GET /posts/:id` - Get a specific post by ID
- `PATCH /posts/:id` - Update a specific post
- `DELETE /posts/:id` - Delete a specific post
- `POST /new` - Create a new post

### Users

- `GET /users/:id` - Get a specific user's profile
- `PATCH /users/:id` - Update a user's profile
- `DELETE /users/:id` - Delete a user account

Replace `:id` with the actual ID of the post or user you want to interact with.

To use these endpoints, send HTTP requests to the appropriate URL with the
necessary data. Make sure to include any required authentication tokens in your
requests for protected endpoints.

Example using curl:

# Get all posts

> Keep in mind, that every single route except the ones stated in
> `Authentication` section requires a UUID token in Headers.

```bash
curl http://localhost:8080/ -H "Token: UUID"
```

# Create a new post

```bash
curl -X POST -H "Content-Type: application/json" -H "Token: UUID" -d '{"title":"My New Post","content":"This is the content of my post"}' http://localhost:8080/new
```

# Get a specific user's profile

```bash
curl http://localhost:8080/users/123 -H "Token: UUID"
```

Replace `localhost:8080` with the actual host and port where your application is
running.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).
