# Auth

Auth is a authentication service for the [Chat](https://github.com/anychatname/) proyect. 
This is a microservices that handles the **Sign Up** and **Login** actions from several platforms.

This service is currently in **PROCESS**.

### External Platforms Available:

 - Google
 - Facebook

## Personal Goal

This project was created to explore new technologies, architectures, methodologies and to reinforce previous knowledge. It has no other purpose than personal growth and added experience.

## Arquitecture

This service is not completly based in the Clean Code arquitecture, rather it is more of a personal architecture based on individual organization of priorities and responsibilities.

## API

### REST

The REST API is based on the [Open API](https://swagger.io/specification/) specification.
A common example of a route is the following:

    .../api/v[versionNumber]/[route]
 
 - `versionNumber`: API version number. Must be a natural number greater than zero ([1, 2, 3, ....]).
- `route`
