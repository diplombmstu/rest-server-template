-----

## Features
- Simple, flexible and testable architecture
  - Light-weight server component
  - Easy to replace components with your own library of choice (router, database driver, etc)
  - Guaranteed thread-safety for each request: uses `gorilla/context` for per-request context
  - Uses dependency-inversion principle (DI) to reduce complexity and manage package dependencies.
- Not a framework
  - More like a project to quickly kick-start your own REST API server, customized to your own needs.
  - Easily extend the project with your own REST resources, in addition to the built-in `users` and `sessions` resources.
- Each REST resource is a separate package
  - Modular approach
  - Separation between `model`, `controller` and `data` layers
  - Clear abstraction from `server` package 
  - Take a look at built-in resources for examples: `users` and `sessions`
  - More example projects coming soon!
- Batteries come included
  - API versioning using using Accept header, for e.g: `Accept=application/json;version=1.0,*/*`
  - Default resources for `users` and `sessions`
  - Access control using activity-based access control (ABAC)
  - Authentication and session management using JWT token
  - Context middleware using `gorilla/context` for per-request context
  - JSON response rendering using `unrolled/render`; extensible to XML or other formats for response
  - MongoDB middleware for database; extensible for other database drivers
- Highly-testable code base
  - Unit-tested `server`; 100% code coverage
  - Easily test REST resources routes
  - Parallelizable test suite
  - Uses `ginkgo` for test framework; optional.


## Architecture
<a href="http://i.imgur.com/HwIhPz7.png"><img src="http://i.imgur.com/HwIhPz7.png"/ height="750"/></a>
