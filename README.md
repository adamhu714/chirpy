# Chirpy

A RESTful web API that implements a backend server for a twitter-like website.

Features:
1. User and 'chirp' creation
2. Authorization with JWT refresh and access tokens
3. JSON data storage safely exposed in an internal database package

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
