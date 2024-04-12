# Lightweight Identity Server in GoLang

This is a very lightweight stand-alone server for managing user identities using JWTs. 
The server currently supports both HTTP and gRPC, and comes with a pre-built admin
web page for creating, deleting and managing users

The main end-points are:

1. User-details via the `/me` end-point
2. Creates, validates and refreshes user JWTs via `/jwt/*`
3. Allows users to change passwords via `/change-password`

## Tech Stack

1. PostgreSQL for persistent storage
2. HTMX and Tailwind for admin web-page
3. Swagger for auto-documentation

## Building and Running

1. Create your `.env` using the attached `.env.template`
2. Check commands in `Makefile` for available commands.