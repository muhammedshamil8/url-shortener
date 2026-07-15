# Contributing to Snippy URL Shortener

Thank you for your interest in contributing to the Snippy URL Shortener! Contributions make the open-source community an amazing place to learn, inspire, and create.

This guide outlines our standards, workflows, and best practices to help get your changes merged smoothly.

---

## 🛠️ Getting Started

### 1. Fork & Clone
1. Fork this repository to your own GitHub account.
2. Clone the fork to your local machine:
   ```bash
   git clone https://github.com/YOUR_USERNAME/url-shortener.git
   cd url-shortener
   ```

### 2. Set Up Local Environment
* **Backend**: Follow the [Getting Started guide in README.md](./README.md#getting-started) to run the Go backend and database containers.
* **Frontend**: Follow the [Frontend Setup guide in frontend/README.md](./frontend/README.md) to set up Vite and run the React client.

---

## 🌿 Branching Strategy

We use branch prefixes to keep our repository organized:
- `feat/` for new features (e.g., `feat/password-reset`)
- `fix/` for bug fixes (e.g., `fix/jwt-expiration`)
- `docs/` for documentation updates (e.g., `docs/contributing-guide`)
- `refactor/` for code refactoring with no behavior changes (e.g., `refactor/repository-interface`)

Create a branch from `dev` before starting work:
```bash
git checkout dev
git checkout -b feat/your-feature-name
```

---

## 📝 Coding Standards & Guidelines

### Go Backend Guidelines
- **Formatting**: Always format your code with `go fmt` (or `make fmt`).
- **Linting**: Ensure your code passes basic safety checks with `go vet` (or `make vet`).
- **Architecture**: Keep backend logic structured within Clean Architecture layers:
  - `internal/handlers` for HTTP routing, binding, and responses.
  - `internal/repository` for SQL queries and DB interaction.
  - `internal/models` for reusable structs.
- **Tests**: Write unit and integration tests under `internal/handlers` or `internal/repository`. Your code should maintain or improve the existing code coverage. Run tests using:
  ```bash
  make test
  ```

### React Frontend Guidelines
- **TypeScript**: Strictly type your props, API responses, and states. Avoid `any`.
- **Styling**: Use Tailwind CSS utilities. Keep custom CSS to a minimum and declare it in `index.css`.
- **Components**: Group reusable components under `frontend/src/components`. Keep views/pages separate under `frontend/src/views`.
- **API Requests**: Use the custom Axios instance from `api.ts` to make calls. This ensures token-refresh flows are handled automatically.

---

## 💬 Commit Message Guidelines

We use conventional commit messages to keep project logs clean:
- **`feat:`** A new feature (e.g., `feat: add custom shortcode alias support`)
- **`fix:`** A bug fix (e.g., `fix: resolve token expiry edge case`)
- **`docs:`** Documentation changes (e.g., `docs: add contributing guidelines`)
- **`style:`** Changes that do not affect the meaning of the code (formatting, whitespace)
- **`refactor:`** A code change that neither fixes a bug nor adds a feature
- **`test:`** Adding or correcting tests

Keep commit messages concise, starting with a lowercase verb in the imperative mood.

---

## 🚀 Submitting a Pull Request

1. Push your branch to your forked repository:
   ```bash
   git push origin feat/your-feature-name
   ```
2. Open a Pull Request on GitHub against the `dev` branch of the main repository.
3. Provide a clear description of the problem solved, testing performed, and any breaking changes.
4. Ensure your branch is updated with the latest changes from the upstream `dev` branch.
