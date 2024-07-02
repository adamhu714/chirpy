# Chirpy

A RESTful web API that implements a backend server in Golang for a twitter-like website.

Features:
- JSON data storage safely exposed in an internal database package
- User and 'chirp' creation, storage and management
- Authorization with JWT refresh and access tokens
- Subscribed users and non-subscribed users
- Middleware for web app pages 

## Contents
* [Getting Started](#getting-started)<br>
  * [Prerequisites](#prerequisites)<br>
  * [Environment Variables](#environment-variables)<br>
  * [Building and Running The Application](#building-and-running-the-application)<br>
* [API Endpoints](#api-endpoints)<br>
  * [/api/users](#apiusers)<br>
  * [/api/login](#apilogin)<br>
  * [/api/chirps](#apichirps)<br>
  * [/api/refresh](#apirefresh)<br>
  * [/api/revoke](#apirevoke)<br>
  * [/api/healthz](#apihealthz)<br>
  * [/api/reset](#apireset)
* [Webpages](#webpages)<br>
  * [/app/index.html](#appindexhtml)
  * [/app/assets/logo.png](#appassetslogopng)
  * [/admin/metrics](#adminmetrics)

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
* [/api/refresh](#apirefresh)<br>
* [/api/revoke](#apirevoke)<br>
* [/api/healthz](#apihealthz)<br>
* [/api/reset](#apireset)

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
  "id": "<User ID>", // integer
  "email": "<User Name>",
  "is_chirpy_red": "<Boolean Value>" // boolean
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
  "id": "<User ID>", // integer
  "email": "<User Name>",
  "is_chirpy_red": "<Boolean Value>" // boolean
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
  "id": "<User ID>", // integer
  "email": "<User Name>",
  "is_chirpy_red": "<Boolean Value>", // boolean
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
  "id": "<Chirp ID>", // integer
  "author_id": "<User ID>", // integer
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
    "id": "<Chirp ID>", // integer
    "author_id": "<User ID>", // integer
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
  "id": "<Chirp ID>", // integer
  "author_id": "<User ID>", // integer
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

### /api/refresh
**POST** `http://localhost:<Port>/api/refresh`

Creates a new JWT access token and returns it in the response.

- Headers: Requires authentication header:
```bash
Authentication: Bearer <Refresh Token>
```
- Request Body: None
- Response Body:
```json
{
  "token": "<Access Token>"
}
```

*[Back To Top](#chirpy)* &nbsp; *[Back To Endpoints](#api-endpoints)*<br>

### /api/revoke
**POST** `http://localhost:<Port>/api/revoke`

Revokes a JWT token.

- Headers: Requires authentication header:
```bash
Authentication: Bearer <JWT Token>
```
- Request Body: None
- Response Body: None

*[Back To Top](#chirpy)* &nbsp; *[Back To Endpoints](#api-endpoints)*<br>

### /api/polka/webhooks
**POST** `http://localhost:<Port>/api/polka/webhooks`

Webhook endpoint for polka service that updates user's subscription status.

- Headers: Requires authentication header:
```bash
Authentication: Bearer <Polka API Key>
```
- Request Body:
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "<UserID>" // integer
  }
}
```
- Response Body: None

*[Back To Top](#chirpy)* &nbsp; *[Back To Endpoints](#api-endpoints)*<br>

### /api/healthz
**GET** `http://localhost:<Port>/api/healthz`

Returns status of the web server.

- Headers: None
- Request Body: None
- Response Body: 
```json
OK
```

*[Back To Top](#chirpy)* &nbsp; *[Back To Endpoints](#api-endpoints)*<br>

### /api/reset
**GET** `http://localhost:<Port>/admin/metrics`

Resets the middleware counter for the number of the times the app pages have been viewed.

- Headers: None
- Request Body: None
- Response Body: None

*[Back To Top](#chirpy)* &nbsp; *[Back To Endpoints](#api-endpoints)*<br>

## Webpages

* [/app/index.html](#appindexhtml)
* [/app/assets/logo.png](#appassetslogopng)
* [/admin/metrics](#adminmetrics)

### /app/index.html

Serves the web app welcome page.

*[Back To Top](#chirpy)* &nbsp; *[Back To Webpages](#webpages)*<br>

### /app/assets/logo.png

Serves the web app logo.

*[Back To Top](#chirpy)* &nbsp; *[Back To Webpages](#webpages)*<br>

### /admin/metrics

Displays how many times an app page is viewed.

*[Back To Top](#chirpy)* &nbsp; *[Back To Webpages](#webpages)*<br>
