# Trello OAuth API

A Go backend API that allows users to authenticate with their Trello accounts using OAuth and access their Trello data.

## Features

- OAuth 2.0 authentication with Trello
- JWT-based authentication for API endpoints
- Secure state parameter validation to prevent CSRF attacks
- Protected API endpoints for accessing Trello data
- Environment-based configuration

## Prerequisites

- Go 1.24 or higher
- Trello API key and secret
  - Create a Trello account and get API credentials from [Trello Developer Portal](https://developer.atlassian.com/cloud/trello/)

## Installation

1. Clone the repository

```bash
git clone https://github.com/mihiryadav20/goauth.git
cd goauth
```

2. Install dependencies

```bash
go mod download
```

3. Configure environment variables

Create a `.env` file in the project root with the following variables:

```
TRELLO_API_KEY=your_trello_api_key
TRELLO_API_SECRET=your_trello_api_secret
CALLBACK_URL=http://localhost:3000/auth/trello/callback
JWT_SECRET=your_secure_jwt_secret_key
PORT=3000
```

## Running the Application

```bash
go run main.go
```

The server will start on port 3000 (or the port specified in your `.env` file).

## API Endpoints

### Public Endpoints

- `GET /` - Welcome message
- `GET /auth/trello` - Initiates the Trello OAuth flow
- `GET /auth/trello/callback` - Handles the OAuth callback from Trello

### Protected Endpoints

All protected endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer your_jwt_token
```

- `GET /api/profile` - Get the authenticated user's Trello profile
- `GET /api/boards` - Get the authenticated user's Trello boards

## OAuth Flow

1. Client requests `/auth/trello`
2. Server returns an authorization URL and state parameter
3. Client redirects user to the authorization URL
4. User authorizes the application on Trello
5. Trello redirects back to `/auth/trello/callback` with token and state
6. Server validates state, gets user info from Trello, and issues a JWT
7. Client uses JWT for subsequent API requests

## Security Considerations

- The JWT secret should be a strong, random string
- In production, use a distributed cache (Redis, etc.) for state storage
- Consider adding rate limiting for public endpoints
- Never commit your `.env` file to version control

## License

MIT

