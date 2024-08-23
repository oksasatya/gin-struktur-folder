# Gin Struktur Folder

This repository provides a recommended folder structure for a Go (Golang) project using the Gin framework. It aims to help developers organize their code in a clean and maintainable way.

## Table of Contents

- [Introduction](#introduction)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Introduction

Gin is a web framework written in Go (Golang). It is known for its speed and minimalist approach. This repository demonstrates a structured way to organize a Gin project, ensuring clarity, scalability, and maintainability.

## Project Structure

The folder structure in this repository is organized as follows:

gin-struktur-folder/
│
├── cmd/ # Command line related files (if applicable)
├── config/ # Configuration files
├── controllers/ # Handlers for routing and request logic
├── middlewares/ # Custom middleware for Gin
├── models/ # Struct definitions and ORM models
├── routes/ # Application routes
├── services/ # Business logic and service layer
├── utils/ # Utility functions and helpers
├── public/ # Static files (images, JS, CSS, etc.)
├── templates/ # HTML templates for server-side rendering
├── .env.example # Example environment variables file
├── go.mod # Go module file
├── go.sum # Go dependencies
└── README.md # Project documentation

## Getting Started

To use this folder structure, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/oksasatya/gin-struktur-folder.git
    ```
2. Navigate to the project directory:

   ```bash
   cd gin-struktur-folder
   ```
3. Install the dependencies:
 
   ```bash
   go mod tidy
   ```
4. Run the application:

   ```bash
   go run main.go
   ```

Usage
This repository serves as a starter template. You can modify it to fit your specific needs. Below are some key points on using this structure:

Controllers: Add your route handlers in the controllers directory.
Middlewares: Place any custom middleware in the middlewares directory.
Models: Define your data models and ORM structures in the models directory.
Routes: Configure all your routes in the routes directory.
Contributing
Contributions are welcome! Please fork this repository and submit a pull request for any improvements.