# Chirpy

A RESTful web API that implements a backend server for a twitter-like website.

Features:
- User and 'chirp' creation
- Authorization with JWT refresh and access tokens
- JSON data storage safely exposed in an internal database package

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

### /v1/users Endpoint

**POST** `http://localhost:<Port>/v1/users`

Creates a new user database entry and returns it.

- Headers: None
- Request Body:
```json
{
  "name": "<User Name>"
}
```
- Response Body:
```json
{
  "id": "<User ID>",
  "created_at": "<Timestamp>",
  "updated_at": "<Timestamp>",
  "name": "<User Name>",
  "apikey": "<API Key>"
}
```


**GET** `http://localhost:<Port>/v1/users`

Returns a user's database entry.

- Headers: Requires authentication header:
```bash
Authentication: APIKey <API Key>
```
- Request Body: None
- Response Body:
```json
{
  "id": "<User ID>",
  "created_at": "<Timestamp>",
  "updated_at": "<Timestamp>",
  "name": "<User Name>",
  "apikey": "<API Key>"
}
```

---
### /v1/feeds Endpoint

