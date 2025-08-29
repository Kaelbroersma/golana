# Golana

[![CI](https://github.com/Kaelbroersma/golana/actions/workflows/ci.yml/badge.svg)](https://github.com/Kaelbroersma/golana/actions/workflows/ci.yml)

A Go-based trading application for Solana tokens.

## Here's the vision...

As a capstone project for completing the boot.dev python & golang course, I am building 2 applications in one..

I will be utilizing docker to spin up both the server in the background(currently in the works) and the client at the same time (a cli tool).
The client (cli) will query the server to complete actions and retrieve information.

Goal:
To make a web-friendly paper trading server.

Spin:
Make the paper trading server accessible from a CLI (Because i'm not that versatile in frontend yet.. I need AI for that. Ironically, it is my job currently.)

<details>
    <summary> The Inspiration</summary>

I had this idea months ago, probably february of 2025- but was beat to market by a company called paper.fun. I actually found them before they launched because I was looking to buy the exact domain and found it was registered.
    
## Oh well...
I took a few months to learn about programming (again), and more specifically, RESTful services, so that I could build something like this when I next got the idea or opportunity.
    
## You may have thought...
when reading that, "If it's your job to build applications with AI currently, why would you spend months re-learning programming and best practices?" Well, that's exactly why.. *best practices*.
    
I want to make sure that everything I build is secure, and omnipotent. When I use AI, it takes a lot of tweaking to get to a point that's even close to omnipotence. This way, I control my own fate.
</details>
---

<details>
  <summary>Example requests and responses for server</summary>

## /api/trades

  <details>
    <summary>PATCH</summary>

    ```http
    PATCH http://localhost:8081/api/trades HTTP/1.1
    Content-Type: application/json
    Authorization: Bearer <your-token-here>

    {
      "trade": "9e458414-1466-4ffa-b30a-4209cace2089",
      "percent": 100
    }
    ```

    ```json
    {
      "ID": "9e458414-1466-4ffa-b30a-4209cace2089",
      "UserID": "ae92cedf-9302-46b1-b47d-89faaa62138e",
      "Contract": "9BB6NFEcjBCtnNLFko2FqVQBq8HHM13kCyYcdQbgpump",
      "Quantity": 1000,
      "OpenPrice": {
        "Float64": 0.817049,
        "Valid": true
      },
      "ClosePrice": {
        "Float64": 0.8064550298592587,
        "Valid": true
      },
      "UnrealizedProfit": {
        "Float64": 0,
        "Valid": false
      },
      "RealizedProfit": {
        "Float64": -10.593970140741373,
        "Valid": true
      },
      "CreatedAt": "2025-08-26T22:40:26Z",
      "UpdatedAt": "2025-08-26T22:40:26Z"
    }
    ```

  </details>

  <details>
    <summary>POST</summary>

    ```http
    POST http://localhost:8081/api/trades HTTP/1.1
    Content-Type: application/json
    Authorization: Bearer <your-token-here>

    {
      "Contract": "<token-contract>",
      "Quantity": <quantity of token to purchase>,
    }
    ```

    ```json
    {
      "ID": "new-trade-id",
      "UserID": "<user-id>",
      "Contract": "<token-contract>",
      "Quantity": 1000,
      "OpenPrice": {
        "Float64": 0.817049,
        "Valid": true
      },
      "ClosePrice": {
        "Float64": 0,
        "Valid": false
      },
      "UnrealizedProfit": {
        "Float64": 0,
        "Valid": false
      },
      "RealizedProfit": {
        "Float64": 0,
        "Valid": false
      },
      "CreatedAt": "2025-08-27T00:00:00Z",
      "UpdatedAt": "2025-08-27T00:00:00Z"
    }
    ```

  </details>
</details>

