# Starter SaaS Template with Go, HTMX, Pocketbase and Tailwind

Features:

- [ ] HTMX, Tailwind, daisyUI for hte frontend.
- [ ] Pocketbase backend
- [ ] Authentication with OAuth2: Google
- [ ] Authentication with username/password
- [ ] Header

## Structure

```
❯ tree . -L 1
.
├── air.toml            // Air config, for file watching
├── dev.sh              // Start dev server
├── index.go            // Main entry point of the application
├── lib                 // Utility go code.
├── middleware          // Middleware for common routes (e.g. authentication)
├── public              // Static public files that should be accessible by raw path.
├── routes              // Backend routes for components and pages. Organized by groups.
├── styles.css          // Custom styles to merge into the tailwind-generated css file
├── tailwind.config.js // Tailwind config and plugins
└── views
```


## How to use

### Installation

```
npm install
```

### OAuth2

- Create a clientID and clientSecret from the OAuth2 provider (e.g. [Google](https://developers.google.com/identity/protocols/oauth2)) .
- Configure the OAuth2 authentication in PocketBase ([link](https://pocketbase.io/docs/authentication#oauth2-integration)).
- Create an `AuthProvider` entry in `routes/auth.go`.
- Add a button in `views/pages/login.html` to htmx the route to the new auth provider.


### Deployment

- For now, just replace all the `localhost` you see with the actual domain, please.
