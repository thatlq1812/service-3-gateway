# Contributing to API Gateway Service

Thank you for considering contributing to the API Gateway Service!

## Development Setup

### Prerequisites

- Go 1.22 or higher
- Docker & Docker Compose
- Git
- Make (optional)

### Getting Started

1. Fork the repository
2. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/service-3-gateway.git
cd service-3-gateway
```

3. Create a branch:
```bash
git checkout -b feature/your-feature-name
```

4. Start development environment:
```bash
docker-compose up -d
```

## Code Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for code formatting
- Run `go vet` to catch common mistakes
- Write clear, self-documenting code with comments for complex logic

### Project Structure

```
service-3-gateway/
├── cmd/server/         # Application entry point
├── internal/           # Private application code
│   ├── client/        # gRPC clients (User, Article)
│   ├── handler/       # HTTP handlers
│   ├── middleware/    # HTTP middleware (auth, logging)
│   └── response/      # Response formatting
├── config/            # Configuration files
└── vendor/            # Dependencies (Go modules)
```

### Commit Messages

Use conventional commit format:
- `feat(handler): add user endpoint`
- `fix(middleware): correct JWT validation`
- `docs(readme): update API examples`
- `refactor(client): improve error handling`
- `test(handler): add article tests`

## Development Workflow

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/handler/...
```

### Building Locally

```bash
# Build binary
go build -o bin/gateway ./cmd/server

# Run locally (requires backend services)
./bin/gateway
```

### Testing Changes

1. Start backend services:
```bash
docker-compose up -d postgres-user postgres-article redis user-service article-service
```

2. Run gateway locally:
```bash
go run ./cmd/server/main.go
```

3. Test endpoints:
```bash
# Health check
curl http://localhost:8080/health

# Test user endpoints
curl http://localhost:8080/api/v1/users

# Test article endpoints
curl http://localhost:8080/api/v1/articles
```

## Pull Request Process

1. Update README.md with details of changes if needed
2. Update CHANGELOG.md with your changes
3. Ensure all tests pass
4. Update documentation for new features
5. Request review from maintainers

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Tests added/updated for changes
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Commit messages follow conventional format
- [ ] No merge conflicts
- [ ] All CI checks passing

## Code Review

- Address all review comments
- Keep discussions focused and professional
- Be open to suggestions and alternatives
- Update PR based on feedback

## API Changes

When modifying HTTP endpoints:

1. Maintain backward compatibility when possible
2. Document breaking changes clearly
3. Update API Reference section in README.md
4. Add migration guide if needed
5. Version endpoints appropriately

## Middleware Changes

When adding/modifying middleware:

1. Keep middleware focused and single-purpose
2. Document middleware behavior
3. Add tests for middleware logic
4. Consider performance impact
5. Update middleware chain documentation

## Questions?

- Open an issue for bugs or feature requests
- Discuss major changes before implementing
- Ask questions in pull requests

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.
