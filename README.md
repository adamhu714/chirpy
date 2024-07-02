# Chirpy

A RESTful web API that implements a backend server in Golang for a twitter-like website.

Features:
- JSON data storage safely exposed in an internal database package
- User and 'chirp' creation and storage
- Authorization with JWT refresh and access tokens
- Subscribed users and non-subscribed users

## Contents
* [Getting Started](#getting-started)<br>
* [API Endpoints](#api-endpoints)<br>
* [Demonstration](#demonstration)

## Getting Started
### Prerequisites
Ensure you have Go v1.22+ installed on your system.

### Environment Variables
Create a `.env` file in your project root directory with the following environment variables:

```bash
JWT_SECRET=<Your JWT Secret Key>
POLKA_APIKEY=<Your Polka API Key>
```

### Building and Running the Application
From the root directory, use the Go command-line tool to build the executable:

```bash
go build -o chirpy
```

This command generates an executable named `chirpy`, which starts the web API server on the specified port.

Execute the built binary:

```bash
./chirpy
```
*[Back To Top](#chirpy)* <br>
## API Endpoints

* [/api/users](#apiusers)<br>

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


*[Back To Top](#chirpy)* <br>
