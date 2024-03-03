# Dinner Done Better Frontend

This is a turborepo.

## What's inside?

This turborepo uses yarn as a package manager. It includes the following packages/apps:

### Apps and Packages

- `web`: the main user-facing [Next.js](https://nextjs.org) app
- `admin`: a [Next.js](https://nextjs.org) app for managing content available to the web app
- `api-client`: a Typescript library used by the applications to talk to the API.
- `models`: a Typescript library used by the applications and the API containing established models.
- `ui`: a React component library used by the applications
- `eslint-config-custom`: `eslint` configurations (includes `eslint-config-next` and `eslint-config-prettier`)
- `tsconfig`: `tsconfig.json`s used throughout the monorepo

Each package/app is 100% [TypeScript](https://www.typescriptlang.org/).

### Utilities

- [TypeScript](https://www.typescriptlang.org/) for static type checking
- [ESLint](https://eslint.org/) for code linting
- [Prettier](https://prettier.io) for code formatting

### Build

To build all apps and packages, run the following command:

```
cd my-turborepo
make build
```

### Develop

To develop all apps and packages, run the following command:

```
cd my-turborepo
make dev
```

## Useful Turbo Links

Learn more about Turborepo:

- [Pipelines](https://turborepo.org/docs/core-concepts/pipelines)
- [Caching](https://turborepo.org/docs/core-concepts/caching)
- [Scoped Tasks](https://turborepo.org/docs/core-concepts/scopes)
- [Configuration Options](https://turborepo.org/docs/reference/configuration)
- [CLI Usage](https://turborepo.org/docs/reference/command-line-reference)
