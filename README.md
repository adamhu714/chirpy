# Chirpy

A RESTful web API that implements a backend server for a twitter-like website.

Features:
- User and 'chirp' creation
- Authorization with JWT refresh and access tokens
- JSON data storage safely exposed in an internal database package
- Subscribed users and non-subscribed users

## Getting Started

### Prerequisites
Ensure you have Go v1.22+ installed on your system.

### Environment Variables
Create a `.env` file in your project root directory with the following environment variables:

```bash
JWT_SECRET=<Your JWT Secret Key>
POLKA_APIKEY=<Your Polka API Key>
```

### Building the Application
From the root directory, use the Go command-line tool to build the executable:

```bash
go build -o chirpy
```

This command generates an executable named `chirpy`, which starts the web API server on the specified port.

### Running the Application

Execute the built binary:

```bash
./chirpy
```

## API Endpoints

### /api/users

**POST** `http://localhost:<Port>/api/users`

Creates a new user database entry and returns it.

- Headers: None
- Request Body:
```json
{
  "email": "<User Name>",
  "password": "<User Password>"
}
```
- Response Body:
```json
{
  "id": "<User ID>",
  "email": "<User Name>",
  "is_chirpy_red": "<Boolean Value>"
}
```


**PUT** `http://localhost:<Port>/api/users`

Updates a user's database details. Updates 
