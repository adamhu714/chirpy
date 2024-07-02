# Chirpy

A RESTful web API that implements a backend server in Golang for a twitter-like website.

Features:
- JSON data storage safely exposed in an internal database package
- User and 'chirp' creation and storage
- Authorization with JWT refresh and access tokens
- Subscribed users and non-subscribed users

## Contents
* [Getting Started](#getting-started)<br>
  * [Prerequisites](#prerequisites)<br>
  * [Environment Variables](#environment-variables)<br>
  * [Building and Running The Application](#building-and-running-the-application)<br>
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
* [/api/login](#apilogin)<br>
* [/api/chirps](#apichirps)<br>

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

Updates a user's database entry details - their email and password. 

- Headers: Requires authentication header:
```bash
Authentication: Bearer <Access Token>
```
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

*[Back To Top](#chirpy)* &nbsp; *[Back To Endpoints](#api-endpoints)*<br>

### /api/login
**POST** `http://localhost:<Port>/api/login`

Authenticates a user and returns to them new access and refresh tokens.

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
  "is_chirpy_red": "<Boolean Value>",
  "token": "<JWT Access Token>",
  "refresh_token": "<JWT Refresh Token>"
}
```

*[Back To Top](#chirpy)* &nbsp; *[Back To Endpoints](#api-endpoints)*<br>

### /api/chirps
**POST** `http://localhost:<Port>/api/chirps`

Creates a new chirp database entry associated with a specific user and returns it.

- Headers: Requires authentication header:
```bash
Authentication: Bearer <Access Token>
```
- Request Body:
```json
{
  "body": "<Chrip Message>"
}
```
- Response Body:
```json
{
  "id": "<Chirp ID>",
  "author_id": "<User ID>",
  "body": "<Chirp Message>"
}
```

**GET** `http://localhost:<Port>/api/chirps`

Returns a list of all chirp database entries. <br>
Defaults sorting by ascending id. Use optional query parameter `sort` to sort chirps in descending order.
For example: `http://localhost:8080/api/chirps?sort=desc`<br>
Use optional query parameter `author_id` to limit chirps to a specific author.
For example: `http://localhost:8080/api/chirps?author_id=1`

- Headers: None
- Request Body: None
- Response Body:
```json
[
  {
    "id": "<Chirp ID>",
    "author_id": "<User ID>",
    "body": "<Chirp Message>"
  }
]
```

**GET** `http://localhost:<Port>/api/chirps/{id}`

Returns a specific chirp database entry with the ID that is provided in the URL.

- Headers: None
- Request Body: None
- Response Body:
```json
{
  "id": "<Chirp ID>",
  "author_id": "<User ID>",
  "body": "<Chirp Message>"
}
```

**DELETE** `http://localhost:<Port>/api/chirps/{id}`

Deletes a specific chirp database entry with the ID that is provided in the URL, for an authenticated user.

- Headers: Requires authentication header:
```bash
Authentication: Bearer <Access Token>
```
- Request Body: None
- Response Body: None

*[Back To Top](#chirpy)* &nbsp; *[Back To Endpoints](#api-endpoints)*<br>
