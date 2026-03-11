# GPS Real-World Examples

This document provides practical, real-world examples of using GPS with AI coding assistants. Each example demonstrates a specific workflow where GPS provides token-efficient project analysis.

---

## Example 1: Project Exploration

### Scenario
A new developer joins a team and needs to quickly understand an unfamiliar codebase before starting work on a feature.

### GPS Commands
```bash
# Get quick project overview
gps --summary

# Full structure at depth 2
gps -L 2

# Focus on key files only
gps --entry-points
```

### AI Prompt
```
I just joined a new team and need to understand this codebase. 
Can you give me an overview of the project structure and tell me 
where I should start exploring?
```

### How AI Uses GPS
The AI runs `gps --summary` first to get a high-level overview, then `gps -L 2` for directory structure, and `gps --entry-points` to identify where the application starts.

### Expected Output
```
project[payment-service]{
  type: go
  files: 47
  size: 312KB
  lines: 8920
  entry: cmd/server/main.go
  tests: 12 files
  docs: README.md, docs/api.md
}
root[47]{
  cmd[3]{
    server[1]{
      main.go [1.2KB, 42L, go]
    }
    cli[1]{
      main.go [856B, 28L, go]
    }
  }
  internal[32]{
    api[8]{
      handlers[5]{
        payment.go [4.2KB, 156L, go]
        refund.go [3.1KB, 112L, go]
      }
      middleware[3]{
        auth.go [2.8KB, 98L, go]
      }
    }
    service[6]{
      processor.go [8.4KB, 312L, go]
      validator.go [5.1KB, 189L, go]
    }
    repository[4]{
      payment_repo.go [6.2KB, 228L, go]
    }
    models[4]{
      payment.go [3.4KB, 125L, go]
    }
  }
  pkg[8]{
    stripe[3]{
      client.go [4.8KB, 178L, go]
    }
  }
}
keyfiles{
  entry: cmd/server/main.go
  config: go.mod, config.yaml
  tests: 12 files
  docs: README.md
}
```

### AI Response
```
Welcome to the payment-service codebase! Here's what you need to know:

**Project Overview**
This is a Go-based payment processing service with 47 files and ~8,900 lines of code.

**Architecture**
- **Entry Point**: `cmd/server/main.go` - Start here to see how the service boots
- **API Layer**: `internal/api/handlers/` - HTTP request handlers
- **Business Logic**: `internal/service/` - Core payment processing
- **Data Layer**: `internal/repository/` - Database operations
- **External Integrations**: `pkg/stripe/` - Stripe payment gateway

**Recommended Exploration Path**
1. Start with `cmd/server/main.go` to understand initialization
2. Review `internal/models/payment.go` for data structures
3. Check `internal/service/processor.go` for core business logic
4. Look at `internal/api/handlers/payment.go` for API endpoints

The project has good test coverage (12 test files) - reviewing tests is a great way to understand expected behavior.
```

---

## Example 2: Entry Point Discovery

### Scenario
A developer needs to find all entry points in a monorepo to understand what services can be deployed independently.

### GPS Commands
```bash
# Find all entry points
gps --entry-points

# With more context
gps -L 3 --focus cmd
```

### AI Prompt
```
This is a monorepo with multiple services. Can you identify all 
the entry points and tell me what services we can deploy independently?
```

### How AI Uses GPS
The AI runs `gps --entry-points` to get all detected entry points, then uses `gps -L 3 --focus cmd` to understand the command structure.

### Expected Output
```
entry-points{
  api-server: cmd/api/main.go
  worker: cmd/worker/main.go
  scheduler: cmd/scheduler/main.go
  cli: cmd/cli/main.go
  migrations: cmd/migrate/main.go
}
```

### AI Response
```
I found 5 independent entry points in this monorepo:

| Service | Entry Point | Purpose |
|---------|-------------|---------|
| API Server | `cmd/api/main.go` | Main HTTP API service |
| Worker | `cmd/worker/main.go` | Background job processor |
| Scheduler | `cmd/scheduler/main.go` | Cron job scheduler |
| CLI | `cmd/cli/main.go` | Command-line tool |
| Migrations | `cmd/migrate/main.go` | Database migration tool |

**Deployment Strategy**
- **API Server** - Primary service, public-facing
- **Worker** - Can scale independently based on queue depth
- **Scheduler** - Single instance, triggers timed jobs
- **CLI** - Development/admin tool, not deployed
- **Migrations** - Run during deployments only

Each can be containerized and deployed independently. The worker is a good candidate for horizontal scaling.
```

---

## Example 3: Focused Module Analysis

### Scenario
A developer needs to understand a specific module without loading the entire project context.

### GPS Commands
```bash
# Focus on specific directory
gps --focus internal/auth -L 3

# Get module summary only
gps --focus internal/auth --summary
```

### AI Prompt
```
I need to understand the authentication module at internal/auth. 
Can you analyze its structure and explain how authentication works?
```

### How AI Uses GPS
The AI uses `--focus` to limit output to the auth module, reducing token usage while getting relevant structure.

### Expected Output
```
project[auth-module]{
  type: go
  files: 12
  size: 48KB
  lines: 1820
}
internal/auth[12]{
  jwt[3]{
    generator.go [4.2KB, 156L, go]
    validator.go [3.8KB, 142L, go]
    claims.go [1.2KB, 45L, go]
  }
  oauth[4]{
    google.go [5.1KB, 189L, go]
    github.go [4.8KB, 178L, go]
    provider.go [2.1KB, 78L, go]
    callback.go [3.4KB, 125L, go]
  }
  middleware[2]{
    require_auth.go [3.2KB, 118L, go]
    optional_auth.go [1.8KB, 67L, go]
  }
  handlers[3]{
    login.go [2.9KB, 108L, go]
    logout.go [1.4KB, 52L, go]
    refresh.go [2.6KB, 96L, go]
  }
}
```

### AI Response
```
The authentication module has a clean, layered architecture:

**Components**

1. **JWT Package** (`jwt/`)
   - `generator.go` - Creates JWT tokens
   - `validator.go` - Validates token signatures and expiration
   - `claims.go` - Defines custom claims structure

2. **OAuth Package** (`oauth/`)
   - `provider.go` - Common OAuth interface
   - `google.go` / `github.go` - Provider-specific implementations
   - `callback.go` - Handles OAuth callbacks

3. **Middleware** (`middleware/`)
   - `require_auth.go` - Blocks unauthenticated requests
   - `optional_auth.go` - Adds user context if authenticated

4. **Handlers** (`handlers/`)
   - `login.go` - Initiates authentication
   - `logout.go` - Invalidates sessions
   - `refresh.go` - Renews tokens

**Flow**: Login → OAuth Provider → Callback → JWT Generation → Token Return

The module supports both OAuth (Google, GitHub) and JWT-based authentication with token refresh capability.
```

---

## Example 4: Architecture Understanding

### Scenario
A team lead needs to verify that a project follows clean architecture principles before a code review.

### GPS Commands
```bash
# Full structure to verify layers
gps -L 2

# Check dependency direction
gps --focus internal -L 3
```

### AI Prompt
```
Review the project structure and tell me if this follows clean 
architecture principles. Identify any potential layer violations.
```

### How AI Uses GPS
The AI analyzes directory structure to verify layered architecture: domain → use cases → interface adapters → frameworks.

### Expected Output
```
project[order-service]{
  type: go
  files: 34
  size: 198KB
}
root[34]{
  cmd[2]{
    server[1]{
      main.go [1.1KB, 38L, go]
    }
  }
  internal[28]{
    domain[6]{
      order.go [4.2KB, 156L, go]
      customer.go [2.8KB, 104L, go]
      repository.go [1.9KB, 72L, go]
    }
    usecase[8]{
      create_order.go [5.1KB, 189L, go]
      cancel_order.go [3.4KB, 125L, go]
      get_order.go [2.8KB, 104L, go]
    }
    infrastructure[6]{
      postgres[3]{
        order_repo.go [6.2KB, 228L, go]
      }
      rabbitmq[2]{
        publisher.go [3.8KB, 142L, go]
      }
    }
    delivery[8]{
      http[5]{
        handler.go [4.8KB, 178L, go]
        routes.go [2.1KB, 78L, go]
      }
      grpc[3]{
        service.go [5.2KB, 192L, go]
      }
    }
  }
  pkg[4]{
    validator[2]{
      order_validator.go [3.2KB, 118L, go]
    }
  }
}
```

### AI Response
```
**Architecture Analysis: Clean Architecture Compliance ✓**

The project follows clean architecture with clear layer separation:

```
┌─────────────────────────────────────────┐
│  Delivery (HTTP/gRPC)                   │  ← Interface Adapters
├─────────────────────────────────────────┤
│  Use Cases (CreateOrder, CancelOrder)   │  ← Application Business Rules
├─────────────────────────────────────────┤
│  Domain (Order, Customer, Repository)   │  ← Enterprise Business Rules
├─────────────────────────────────────────┤
│  Infrastructure (Postgres, RabbitMQ)    │  ← Frameworks & Drivers
└─────────────────────────────────────────┘
```

**Compliance Check**
| Principle | Status | Notes |
|-----------|--------|-------|
| Domain isolation | ✓ | No external dependencies in domain/ |
| Use case independence | ✓ | Business logic in usecase/ |
| Interface adapters | ✓ | HTTP/gRPC in delivery/ |
| Framework isolation | ✓ | DB/messaging in infrastructure/ |
| Dependency direction | ✓ | Dependencies point inward |

**No layer violations detected.** The structure properly isolates:
- Business rules from frameworks
- Domain logic from delivery mechanisms
- Infrastructure concerns from core logic

The `pkg/validator` is appropriately placed as shared utility code.
```

---

## Example 5: Code Review Preparation

### Scenario
A developer is preparing for a code review and wants to understand the files that were changed in a PR.

### GPS Commands
```bash
# Get quick overview for context
gps --summary

# Focus on changed areas
gps --focus internal/api -L 2 --no-meta
```

### AI Prompt
```
I'm reviewing a PR that changes files in internal/api/handlers. 
Give me context about this area of the codebase so I can do a thorough review.
```

### How AI Uses GPS
The AI gets overall context with `--summary`, then focuses on the changed area with minimal output using `--no-meta`.

### Expected Output
```
project[user-service]{
  type: go
  files: 28
  size: 186KB
  entry: cmd/server/main.go
}
internal/api[14]{
  handlers[8]{
    user_handler.go
    auth_handler.go
    profile_handler.go
    search_handler.go
  }
  middleware[3]{
    auth.go
    logging.go
    ratelimit.go
  }
  routes[2]{
    router.go
    routes_test.go
  }
  dto[1]{
    requests.go
  }
}
```

### AI Response
```
**Code Review Context: internal/api**

**Area Overview**
The `internal/api` package contains 14 files handling HTTP concerns:

**Handlers** (8 files)
- User CRUD operations
- Authentication endpoints
- Profile management
- Search functionality

**Supporting Components**
- `middleware/` - Auth, logging, rate limiting
- `routes/` - URL routing configuration
- `dto/` - Request/response structures

**Review Checklist for Handler Changes**
1. ✓ Proper error handling and status codes
2. ✓ Input validation via DTOs
3. ✓ Auth middleware applied where needed
4. ✓ Consistent response format
5. ✓ No business logic in handlers (should delegate to service layer)

**Key Files to Review**
- Check `routes/router.go` if new endpoints added
- Verify `dto/requests.go` for new input structures
- Review middleware chain in `middleware/`

**Watch for**: Handlers calling repositories directly instead of services - this violates layering.
```

---

## Example 6: Refactoring Planning

### Scenario
A team is planning a major refactoring effort and needs to understand dependencies between modules.

### GPS Commands
```bash
# Full structure analysis
gps -L 3

# Token-efficient for large project
gps -L 2 --no-meta
```

### AI Prompt
```
We need to refactor the user module. Analyze the project structure 
and identify which other modules might be affected by changes to 
internal/user.
```

### How AI Uses GPS
The AI maps the full structure to understand module boundaries and identify coupling points.

### Expected Output
```
root[89]{
  cmd[3]{
    server[1]{ main.go }
  }
  internal[72]{
    user[12]{
      service.go
      repository.go
      model.go
      cache.go
    }
    auth[8]{
      service.go
      jwt.go
      password.go
    }
    order[10]{
      service.go
      repository.go
      model.go
    }
    notification[6]{
      service.go
      email.go
      push.go
    }
    api[14]{
      handlers[8]{
        user_handler.go
        auth_handler.go
        order_handler.go
      }
      middleware[3]{
        auth.go
      }
    }
    search[4]{
      indexer.go
      user_search.go
    }
  }
  pkg[8]{
    events[3]{
      publisher.go
      user_events.go
    }
  }
}
```

### AI Response
```
**Refactoring Impact Analysis: internal/user**

**Direct Dependencies** (modules that import user)
| Module | Files | Risk Level |
|--------|-------|------------|
| `internal/auth` | service.go | High - Auth depends on user validation |
| `internal/api/handlers` | user_handler.go | High - Direct API coupling |
| `internal/search` | user_search.go | Medium - Indexes user data |
| `internal/notification` | service.go | Low - Only uses user ID |
| `internal/order` | model.go | Low - References user FK |
| `pkg/events` | user_events.go | Medium - Event publishing |

**Refactoring Strategy**

1. **Phase 1: Interface Extraction**
   - Define `UserRepository` interface
   - Define `UserService` interface
   - Keep existing implementations

2. **Phase 2: Consumer Updates**
   - Update `auth/service.go` to use interface
   - Update `api/handlers/user_handler.go`
   - Update `search/user_search.go`

3. **Phase 3: Core Refactor**
   - Modify user module internals
   - Add new repository implementations
   - Update event schemas

4. **Phase 4: Cleanup**
   - Remove deprecated code
   - Update tests

**Breaking Changes to Avoid**
- Don't change `user.Model` public fields without migration
- Keep `pkg/events/user_events.go` event structure compatible
- Maintain `UserService` method signatures

**Test Coverage Required**
- Auth integration tests
- API handler tests
- Search indexing tests
```

---

## Example 7: Documentation Generation

### Scenario
A developer needs to generate architecture documentation for an internal wiki.

### GPS Commands
```bash
# Human-readable format
gps -f tree -L 3

# Summary for overview
gps --summary
```

### AI Prompt
```
Generate architecture documentation for this project that I can 
add to our internal wiki. Include a directory structure overview 
and describe the purpose of each major component.
```

### How AI Uses GPS
The AI uses tree format for human-readable structure and summary for project metadata.

### Expected Output
```
project[inventory-service]{
  type: go
  files: 42
  size: 287KB
  lines: 9520
  entry: cmd/server/main.go
}
inventory-service/
├── cmd/
│   └── server/
│       ├── main.go [1.2KB, 42L]
│       └── config.go [2.1KB, 78L]
├── internal/
│   ├── domain/
│   │   ├── product.go [3.8KB, 142L]
│   │   ├── inventory.go [4.2KB, 156L]
│   │   └── warehouse.go [2.9KB, 108L]
│   ├── service/
│   │   ├── inventory_service.go [8.4KB, 312L]
│   │   ├── product_service.go [6.1KB, 225L]
│   │   └── warehouse_service.go [5.2KB, 192L]
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── inventory_repo.go [7.8KB, 288L]
│   │   │   └── product_repo.go [6.4KB, 236L]
│   │   └── redis/
│   │       └── cache.go [3.2KB, 118L]
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── inventory_handler.go [5.6KB, 208L]
│   │   │   └── product_handler.go [4.8KB, 178L]
│   │   └── middleware/
│   │       ├── auth.go [2.8KB, 104L]
│   │       └── logging.go [1.9KB, 72L]
│   └── events/
│       ├── publisher.go [4.1KB, 152L]
│       └── handlers.go [5.8KB, 214L]
├── pkg/
│   ├── client/
│   │   └── erp_client.go [6.2KB, 228L]
│   └── utils/
│       └── pagination.go [2.4KB, 89L]
├── go.mod [312B, 12L]
└── go.sum [1.4KB, 32L]
```

### AI Response
```markdown
# Inventory Service Architecture

## Overview

The Inventory Service is a Go-based microservice responsible for managing product inventory across multiple warehouses. It provides real-time inventory tracking, synchronization with ERP systems, and event-driven updates.

| Metric | Value |
|--------|-------|
| Files | 42 |
| Size | 287KB |
| Lines of Code | 9,520 |
| Entry Point | `cmd/server/main.go` |

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         API Layer                           │
│  ┌─────────────────┐  ┌──────────────────────────────────┐ │
│  │    Handlers     │  │         Middleware               │ │
│  │  • Inventory    │  │  • Authentication                │ │
│  │  • Product      │  │  • Request Logging               │ │
│  └────────┬────────┘  └──────────────────────────────────┘ │
└───────────┼─────────────────────────────────────────────────┘
            │
┌───────────┼─────────────────────────────────────────────────┐
│           │           Service Layer                          │
│  ┌────────┴────────┐  ┌────────────────┐  ┌──────────────┐ │
│  │ Inventory Svc   │  │ Product Svc    │  │ Warehouse Svc│ │
│  │ • Stock mgmt    │  │ • Catalog      │  │ • Locations  │ │
│  │ • Reservations  │  │ • Pricing      │  │ • Capacity   │ │
│  └────────┬────────┘  └───────┬────────┘  └──────┬───────┘ │
└───────────┼───────────────────┼─────────────────┼──────────┘
            │                   │                 │
┌───────────┼───────────────────┼─────────────────┼──────────┐
│           │           Repository Layer                        │
│  ┌────────┴────────┐  ┌──────┴───────┐  ┌──────┴───────┐   │
│  │  PostgreSQL     │  │   Redis      │  │   Events     │   │
│  │  • Inventory    │  │  • Cache     │  │  • Publish   │   │
│  │  • Products     │  │  • Sessions  │  │  • Consume   │   │
│  └─────────────────┘  └──────────────┘  └──────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

| Directory | Purpose |
|-----------|---------|
| `cmd/server/` | Application entry point and configuration |
| `internal/domain/` | Core business entities and value objects |
| `internal/service/` | Business logic and orchestration |
| `internal/repository/` | Data access layer (PostgreSQL, Redis) |
| `internal/api/` | HTTP handlers and middleware |
| `internal/events/` | Event publishing and consumption |
| `pkg/client/` | External service clients (ERP integration) |
| `pkg/utils/` | Shared utility functions |

## Key Components

### Domain Layer (`internal/domain/`)
- **product.go** - Product entity with SKU, pricing, attributes
- **inventory.go** - Stock levels, reservations, allocations
- **warehouse.go** - Warehouse locations and capacity

### Service Layer (`internal/service/`)
- **InventoryService** - Core inventory management, stock adjustments
- **ProductService** - Product catalog operations
- **WarehouseService** - Warehouse management

### API Layer (`internal/api/`)
- REST endpoints for inventory and product operations
- Authentication middleware validates JWT tokens
- Request logging for audit trail

### Event System (`internal/events/`)
- Publishes inventory change events
- Consumes order events for reservations
- Kafka-based messaging

## External Dependencies

| Service | Purpose | Client |
|---------|---------|--------|
| ERP System | Master data sync | `pkg/client/erp_client.go` |
| PostgreSQL | Persistent storage | `internal/repository/postgres/` |
| Redis | Caching & sessions | `internal/repository/redis/` |
| Kafka | Event streaming | `internal/events/` |
```

---

## Example 8: CI/CD Integration

### Scenario
A DevOps engineer wants to add automated project analysis to a CI pipeline to validate project structure.

### GPS Commands
```bash
# JSON for programmatic processing
gps -f json

# Parse with jq
gps -f json | jq '.project.key_files.entry_points | length'
```

### AI Prompt
```
Help me create a CI/CD pipeline step that validates our Go project 
structure. I need to check that entry points exist and that we have 
test coverage.
```

### How AI Uses GPS
The AI uses JSON output format for programmatic access and demonstrates jq integration.

### Expected Output
```json
{
  "project": {
    "name": "myapp",
    "type": "go",
    "stats": {
      "file_count": 34,
      "total_lines": 8920,
      "by_type": {
        "go": 28,
        "yaml": 4,
        "md": 2
      }
    },
    "key_files": {
      "entry_points": ["cmd/server/main.go"],
      "configs": ["go.mod", "go.sum"],
      "tests": [
        "internal/service/order_service_test.go",
        "internal/repository/postgres/order_repo_test.go",
        "internal/api/handlers/order_handler_test.go"
      ],
      "docs": ["README.md"]
    }
  }
}
```

### AI Response
```yaml
# .github/workflows/validate-structure.yml
name: Validate Project Structure

on:
  push:
    branches: [main, develop]
  pull_request:

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install GPS
        run: |
          curl -sL https://github.com/wesbragagt/gps/releases/latest/download/gps-Linux-x86_64 -o gps
          chmod +x gps
          sudo mv gps /usr/local/bin/
      
      - name: Check entry points exist
        run: |
          ENTRY_COUNT=$(gps -f json | jq '.project.key_files.entry_points | length')
          if [ "$ENTRY_COUNT" -eq 0 ]; then
            echo "Error: No entry points found"
            exit 1
          fi
          echo "Found $ENTRY_COUNT entry point(s)"
      
      - name: Validate test coverage
        run: |
          TEST_COUNT=$(gps -f json | jq '.project.key_files.tests | length')
          if [ "$TEST_COUNT" -lt 3 ]; then
            echo "Warning: Low test coverage ($TEST_COUNT test files)"
          fi
          echo "Found $TEST_COUNT test file(s)"
      
      - name: Check project type
        run: |
          PROJECT_TYPE=$(gps -f json -r '.project.type')
          if [ "$PROJECT_TYPE" != "go" ]; then
            echo "Error: Expected Go project, found $PROJECT_TYPE"
            exit 1
          fi
          echo "Project type: $PROJECT_TYPE"
      
      - name: Generate structure report
        run: |
          gps --summary > structure-report.txt
          cat structure-report.txt
          
      - name: Upload structure report
        uses: actions/upload-artifact@v4
        with:
          name: structure-report
          path: structure-report.txt
```

**Additional Validation Scripts**

```bash
#!/bin/bash
# scripts/validate-structure.sh

set -e

echo "Validating project structure..."

# Check for required entry points
ENTRY_POINTS=$(gps -f json | jq -r '.project.key_files.entry_points[]')
for entry in $ENTRY_POINTS; do
    if [ ! -f "$entry" ]; then
        echo "ERROR: Entry point $entry not found"
        exit 1
    fi
    echo "✓ Entry point: $entry"
done

# Verify test coverage
GO_FILES=$(gps -f json | jq '.project.stats.by_type.go')
TEST_FILES=$(gps -f json | jq '.project.key_files.tests | length')
RATIO=$(echo "scale=2; $TEST_FILES / $GO_FILES * 100" | bc)

echo "Test ratio: ${RATIO}% ($TEST_FILES tests / $GO_FILES source files)"

if (( $(echo "$RATIO < 20" | bc -l) )); then
    echo "WARNING: Test coverage below 20%"
fi

# Check for documentation
DOCS=$(gps -f json | jq '.project.key_files.docs | length')
if [ "$DOCS" -eq 0 ]; then
    echo "WARNING: No documentation files found"
fi

echo "Structure validation complete!"
```

---

## Example 9: Multi-Language Projects

### Scenario
A developer is working with a polyglot codebase containing Go services, Node.js frontend, and Python scripts.

### GPS Commands
```bash
# Analyze full monorepo
gps -L 2

# Focus on specific language area
gps --focus services -L 2
gps --focus frontend -L 2
```

### AI Prompt
```
This is a monorepo with multiple languages. Help me understand the 
structure and how the different components interact.
```

### How AI Uses GPS
The AI uses GPS to map the polyglot structure and identify language-specific areas.

### Expected Output
```
project[platform-monorepo]{
  type: mixed
  files: 156
  size: 2.8MB
  lines: 48200
}
root[156]{
  services[68]{
    api[24]{
      cmd[2]{
        main.go [1.4KB, 52L, go]
      }
      internal[18]{
        handlers[8]{
          user_handler.go [4.2KB, 156L, go]
        }
        service[6]{
          user_service.go [6.8KB, 252L, go]
        }
      }
      go.mod [412B, 15L, go]
    }
    worker[12]{
      cmd[1]{
        main.go [1.1KB, 42L, go]
      }
      internal[8]{
        processor.go [8.2KB, 304L, go]
      }
      go.mod [386B, 14L, go]
    }
  }
  frontend[52]{
    src[38]{
      components[14]{
        Dashboard.tsx [3.8KB, 142L, tsx]
        UserList.tsx [2.9KB, 108L, tsx]
      }
      pages[8]{
        index.tsx [2.1KB, 78L, tsx]
        users.tsx [3.4KB, 125L, tsx]
      }
      api[6]{
        client.ts [2.8KB, 104L, ts]
        auth.ts [1.9KB, 72L, ts]
      }
      utils[4]{
        format.ts [1.2KB, 45L, ts]
      }
    }
    package.json [1.8KB, 68L, json]
    tsconfig.json [856B, 32L, json]
  }
  scripts[18]{
    deployment[6]{
      deploy.py [4.2KB, 156L, py]
      rollback.py [3.1KB, 112L, py]
    }
    data[8]{
      migrate.py [5.8KB, 214L, py]
      seed.py [4.4KB, 162L, py]
    }
    requirements.txt [512B, 18L, txt]
  }
  infra[12]{
    terraform[8]{
      main.tf [3.2KB, 118L, tf]
      variables.tf [1.8KB, 67L, tf]
    }
    docker[4]{
      Dockerfile.api [1.2KB, 45L, dockerfile]
      docker-compose.yml [2.4KB, 89L, yaml]
    }
  }
}
keyfiles{
  entry: services/api/cmd/main.go, services/worker/cmd/main.go
  config: package.json, go.mod, requirements.txt
  tests: 24 files
  docs: README.md, docs/architecture.md
}
```

### AI Response
```
**Monorepo Architecture Analysis**

This is a mixed-language monorepo with 156 files across 4 main areas:

## Language Distribution

| Area | Language | Files | Purpose |
|------|----------|-------|---------|
| `services/` | Go | 36 | Backend microservices |
| `frontend/` | TypeScript | 52 | React web application |
| `scripts/` | Python | 18 | DevOps and data scripts |
| `infra/` | HCL/YAML | 12 | Infrastructure as code |

## Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (TypeScript)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │  Components  │  │    Pages     │  │   API Client     │  │
│  │  Dashboard   │  │  /users      │  │  auth, requests  │  │
│  └──────────────┘  └──────────────┘  └────────┬─────────┘  │
└─────────────────────────────────────────────────┼───────────┘
                                                  │ HTTP
┌─────────────────────────────────────────────────┼───────────┐
│                  Services (Go)                  │           │
│  ┌──────────────┐  ┌──────────────┐            │           │
│  │   API Svc    │◄─┤   Worker     │◄─── Events ┘           │
│  │  :8080       │  │  Background  │                        │
│  └──────────────┘  └──────────────┘                        │
└─────────────────────────────────────────────────────────────┘
         │                      │
┌────────┴────────┐    ┌────────┴────────┐
│ Python Scripts  │    │   Terraform     │
│ Deploy, Migrate │    │   AWS Infra     │
└─────────────────┘    └─────────────────┘
```

## Key Interactions

1. **Frontend → API Service**
   - TypeScript API client (`src/api/client.ts`)
   - Communicates with `services/api` via HTTP

2. **API Service → Worker**
   - Event-driven communication
   - Async job processing

3. **Python Scripts**
   - Deployment automation
   - Database migrations
   - Interact with both Go services

## Development Workflows

| Task | Entry Point |
|------|-------------|
| Start API | `cd services/api && go run cmd/main.go` |
| Start Frontend | `cd frontend && npm run dev` |
| Run Migrations | `python scripts/data/migrate.py` |
| Deploy | `python scripts/deployment/deploy.py` |
```

---

## Example 10: Performance Optimization - Token Efficiency

### Scenario
A developer is working with a large codebase and needs to maximize token efficiency when providing context to an AI assistant.

### GPS Commands
```bash
# Compare token counts
gps --compare

# Maximum token efficiency
gps -f toon -L 2 --no-meta

# Ultra-minimal structure
gps -L 1 -d --no-meta
```

### AI Prompt
```
I'm working with a large codebase and need to understand its 
structure without using too many tokens. Show me the most 
token-efficient way to get project context.
```

### How AI Uses GPS
The AI demonstrates GPS's token efficiency features and compares formats.

### Expected Output
```
# gps --compare
Token Comparison (150 files, 450KB project):
  JSON:  4,820 tokens
  TOON:  1,928 tokens  (60% reduction)
  Tree:  2,640 tokens  (45% reduction)
  Flat:  2,180 tokens  (55% reduction)

# gps -f toon -L 2 --no-meta
project[large-service]{
  type: go
  files: 150
}
root[150]{
  cmd[4]{
    server[2]{}
    cli[2]{}
  }
  internal[98]{
    api[24]{}
    service[32]{}
    repository[18]{}
    domain[12]{}
    events[12]{}
  }
  pkg[24]{
    utils[12]{}
    client[12]{}
  }
  configs[8]{}
  scripts[12]{}
}

# gps -L 1 -d --no-meta
root[12]{
  cmd[2]
  internal[5]
  pkg[2]
  configs[1]
  scripts[1]
}
```

### AI Response
```
**Token-Efficient Context Strategies**

## Format Comparison for This Project (150 files)

| Format | Tokens | vs JSON | Best For |
|--------|--------|---------|----------|
| JSON | 4,820 | baseline | Scripting, CI/CD |
| TOON | 1,928 | **60% savings** | AI context |
| Tree | 2,640 | 45% savings | Human reading |
| Flat | 2,180 | 55% savings | Data analysis |

## Token Optimization Techniques

### 1. Use TOON Format (Default)
```bash
gps  # Already uses TOON by default
```
**Saves 60%** vs JSON while keeping all structure info.

### 2. Limit Depth
```bash
gps -L 2   # 2 levels deep
gps -L 1   # Top-level only
```
**Saves 40-70%** depending on project depth.

### 3. Skip Metadata
```bash
gps --no-meta
```
**Saves 15-25%** when you only need structure.

### 4. Directories Only
```bash
gps -d
```
**Saves 50-80%** when file details aren't needed.

### 5. Focus on Area
```bash
gps --focus internal/service -L 2
```
**Saves 70-90%** by limiting scope.

### 6. Combined Maximum Efficiency
```bash
gps -L 2 -d --no-meta
```
**Saves 85-95%** - ultra-minimal structure.

## Recommended Workflow

For large projects, use progressive disclosure:

```bash
# Step 1: Ultra-quick overview (50-100 tokens)
gps -L 1 -d --no-meta

# Step 2: Focused exploration (200-500 tokens)
gps --focus internal/api -L 2

# Step 3: Detailed when needed (500-1000 tokens)
gps --focus internal/api/handlers -L 3
```

## Real Savings Example

| Task | Naive Approach | GPS Optimized | Savings |
|------|----------------|---------------|---------|
| Full project context | `tree -L 5` (5,200 tokens) | `gps -L 2` (1,900 tokens) | 63% |
| Find entry points | Read all main.go files | `gps --entry-points` (50 tokens) | 95% |
| Module overview | `find . -type f` (3,800 tokens) | `gps --focus module -L 2` (400 tokens) | 89% |

**Tip**: Start minimal, expand as needed. GPS lets you iteratively explore without blowing your context budget.
```

---

## Example 11: Node.js Project Analysis

### Scenario
A developer needs to understand a Node.js/Express API project structure.

### GPS Commands
```bash
# Node.js project overview
gps -L 3 -I "node_modules,dist,.git"

# Focus on source
gps --focus src -L 3
```

### AI Prompt
```
Analyze this Node.js Express API project and explain its structure. 
Help me understand how the routes and middleware are organized.
```

### How AI Uses GPS
The AI maps the Node.js project structure, excluding node_modules for efficiency.

### Expected Output
```
project[express-api]{
  type: node
  files: 42
  size: 186KB
  lines: 5420
  entry: src/index.js
}
root[42]{
  src[32]{
    routes[10]{
      index.js [1.2KB, 42L, js]
      users.js [2.8KB, 104L, js]
      products.js [3.4KB, 125L, js]
      orders.js [4.1KB, 152L, js]
      auth.js [2.2KB, 82L, js]
    }
    controllers[8]{
      userController.js [4.2KB, 156L, js]
      productController.js [5.1KB, 189L, js]
      orderController.js [4.8KB, 178L, js]
      authController.js [3.6KB, 134L, js]
    }
    middleware[6]{
      auth.js [3.2KB, 118L, js]
      validate.js [2.4KB, 89L, js]
      errorHandler.js [1.8KB, 67L, js]
      rateLimit.js [1.4KB, 52L, js]
    }
    models[4]{
      User.js [2.8KB, 104L, js]
      Product.js [2.4KB, 89L, js]
      Order.js [3.1KB, 112L, js]
    }
    services[4]{
      emailService.js [3.8KB, 142L, js]
      paymentService.js [4.2KB, 156L, js]
    }
    config[2]{
      database.js [1.6KB, 58L, js]
      app.js [2.1KB, 78L, js]
    }
    index.js [1.8KB, 68L, js]
  }
  tests[8]{
    routes[4]{
      users.test.js [2.4KB, 89L, js]
      products.test.js [2.8KB, 104L, js]
    }
    middleware[2]{
      auth.test.js [1.9KB, 72L, js]
    }
  }
  package.json [1.4KB, 52L, json]
}
keyfiles{
  entry: src/index.js
  config: package.json
  tests: 6 files
  docs: README.md
}
```

### AI Response
```
**Node.js Express API Structure**

## Architecture Overview

This is a well-organized Express.js API following the MVC pattern with 42 files and ~5,400 lines of code.

## Directory Structure

```
src/
├── index.js           # Application entry point
├── config/            # Configuration (database, app settings)
├── routes/            # Route definitions (URL → Controller)
├── controllers/       # Request handlers
├── models/            # Mongoose/Sequelize models
├── middleware/        # Express middleware
└── services/          # Business logic & external integrations
```

## Request Flow

```
HTTP Request
    │
    ▼
┌─────────────┐
│ Middleware  │  auth.js → validate.js → rateLimit.js
└─────┬───────┘
      │
      ▼
┌─────────────┐
│   Routes    │  /api/users → users.js
└─────┬───────┘
      │
      ▼
┌─────────────┐
│ Controllers │  userController.js
└─────┬───────┘
      │
      ▼
┌─────────────┐
│   Models    │  User.js (MongoDB/PostgreSQL)
└─────────────┘
```

## Key Components

### Routes (`src/routes/`)
- `index.js` - Mounts all route modules
- Resource-specific files (users, products, orders)
- `auth.js` - Authentication endpoints

### Controllers (`src/controllers/`)
- Handle HTTP request/response
- Validate input
- Call services for business logic

### Middleware (`src/middleware/`)
- `auth.js` - JWT verification
- `validate.js` - Request validation (Joi/express-validator)
- `errorHandler.js` - Global error handling
- `rateLimit.js` - Rate limiting

### Services (`src/services/`)
- `emailService.js` - Email sending (SendGrid/AWS SES)
- `paymentService.js` - Payment processing (Stripe)

## Entry Point Analysis

`src/index.js` likely:
1. Loads environment variables
2. Connects to database (`config/database.js`)
3. Configures Express app (`config/app.js`)
4. Mounts routes
5. Starts server

## Quick Start Commands

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Run tests
npm test
```
```

---

## Example 12: Python Project Analysis

### Scenario
A data scientist needs to understand a Python ML project structure for collaboration.

### GPS Commands
```bash
# Python project structure
gps -L 3 -I "__pycache__,*.pyc,.venv,venv"

# Focus on source code
gps --focus src -L 3
```

### AI Prompt
```
Help me understand this Python machine learning project. 
I need to know how the data pipelines and models are organized.
```

### How AI Uses GPS
The AI analyzes the Python project structure, identifying ML-specific patterns.

### Expected Output
```
project[ml-pipeline]{
  type: python
  files: 38
  size: 245KB
  lines: 6820
  entry: src/train.py
}
root[38]{
  src[26]{
    data[8]{
      loader.py [4.2KB, 156L, py]
      preprocessor.py [5.8KB, 214L, py]
      feature_engineer.py [6.4KB, 236L, py]
      validator.py [2.8KB, 104L, py]
    }
    models[6]{
      base.py [2.4KB, 89L, py]
      random_forest.py [4.8KB, 178L, py]
      xgboost_model.py [5.2KB, 192L, py]
      neural_net.py [8.4KB, 312L, py]
    }
    training[4]{
      trainer.py [6.2KB, 228L, py]
      evaluator.py [4.6KB, 172L, py]
    }
    inference[4]{
      predictor.py [3.8KB, 142L, py]
      api.py [4.2KB, 156L, py]
    }
    utils[4]{
      config.py [2.1KB, 78L, py]
      logging.py [1.8KB, 67L, py]
    }
    train.py [2.4KB, 89L, py]
    predict.py [1.6KB, 58L, py]
  }
  notebooks[4]{
    exploration.ipynb [45KB, 892L, ipynb]
    experiments.ipynb [62KB, 1240L, ipynb]
  }
  tests[6]{
    test_data.py [3.2KB, 118L, py]
    test_models.py [4.1KB, 152L, py]
  }
  configs[4]{
    model_config.yaml [1.2KB, 45L, yaml]
    features.yaml [2.4KB, 89L, yaml]
  }
  requirements.txt [856B, 32L, txt]
  setup.py [1.4KB, 52L, py]
}
keyfiles{
  entry: src/train.py, src/predict.py
  config: requirements.txt, setup.py, configs/
  tests: 4 files
  docs: README.md
}
```

### AI Response
```
**Python ML Pipeline Structure**

## Project Overview

A machine learning pipeline project with 38 files organized for data processing, model training, and inference.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Entry Points                           │
│  ┌────────────────┐           ┌────────────────────────┐   │
│  │   train.py     │           │     predict.py         │   │
│  │   Training     │           │     Inference          │   │
│  └───────┬────────┘           └───────────┬────────────┘   │
└──────────┼────────────────────────────────┼────────────────┘
           │                                │
┌──────────┼────────────────────────────────┼────────────────┐
│          │        Data Pipeline           │                │
│  ┌───────┴───────┐  ┌─────────────┐  ┌────┴───────────┐   │
│  │   loader.py   │→ │preprocessor │→ │feature_eng.py │   │
│  └───────────────┘  └─────────────┘  └────────────────┘   │
└─────────────────────────────────────────────────────────────┘
           │
┌──────────┼────────────────────────────────────────────────┐
│          │         Model Layer                            │
│  ┌───────┴───────────────────────────────────────────┐    │
│  │                    base.py                        │    │
│  └───────┬───────────┬────────────┬──────────────────┘    │
│  │       │           │            │                       │
│  │  random_forest  xgboost    neural_net                 │
│  └───────────────────────────────────────────────────────┘
└─────────────────────────────────────────────────────────────┘
           │
┌──────────┼────────────────────────────────────────────────┐
│          │      Training & Evaluation                     │
│  ┌───────┴───────┐           ┌─────────────────┐          │
│  │   trainer.py  │──────────►│  evaluator.py   │          │
│  └───────────────┘           └─────────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

## Key Directories

### `src/data/` - Data Pipeline
| File | Purpose |
|------|---------|
| `loader.py` | Load data from various sources (CSV, DB, S3) |
| `preprocessor.py` | Clean and normalize data |
| `feature_engineer.py` | Create and transform features |
| `validator.py` | Data quality checks |

### `src/models/` - Model Definitions
| File | Purpose |
|------|---------|
| `base.py` | Abstract base class for models |
| `random_forest.py` | Sklearn RandomForest implementation |
| `xgboost_model.py` | XGBoost model wrapper |
| `neural_net.py` | PyTorch/TensorFlow neural network |

### `src/training/` - Training Pipeline
| File | Purpose |
|------|---------|
| `trainer.py` | Model training orchestration |
| `evaluator.py` | Metrics and model evaluation |

### `src/inference/` - Prediction Service
| File | Purpose |
|------|---------|
| `predictor.py` | Load model and make predictions |
| `api.py` | REST API for model serving |

## Configuration

- `configs/model_config.yaml` - Model hyperparameters
- `configs/features.yaml` - Feature definitions
- `requirements.txt` - Python dependencies

## Typical Workflow

```bash
# 1. Install dependencies
pip install -r requirements.txt

# 2. Explore data (optional)
jupyter notebook notebooks/exploration.ipynb

# 3. Train model
python src/train.py --config configs/model_config.yaml

# 4. Run predictions
python src/predict.py --model models/trained_model.pkl

# 5. Start inference API
python src/inference/api.py
```

## Extension Points

1. **Add new model**: Create file in `src/models/` extending `base.py`
2. **Add features**: Modify `feature_engineer.py` or update `features.yaml`
3. **Add data source**: Extend `loader.py` with new connector
```

---

## Summary

These examples demonstrate GPS's versatility in real-world AI assistant workflows:

| Use Case | Key GPS Feature | Token Savings |
|----------|-----------------|---------------|
| Project Exploration | `--summary`, `-L 2` | 60-70% |
| Entry Point Discovery | `--entry-points` | 90-95% |
| Focused Analysis | `--focus` | 70-90% |
| Architecture Review | `-L 2` with analysis | 50-60% |
| Code Review Prep | `--no-meta`, `--focus` | 60-80% |
| Refactoring Planning | Full `-L 3` scan | 40-50% |
| Documentation | `-f tree`, `--summary` | N/A (human use) |
| CI/CD Integration | `-f json` with jq | N/A (automation) |
| Multi-Language | Auto-detection | 50-60% |
| Token Optimization | `--compare`, `--no-meta` | 85-95% |
| Node.js Projects | `-I "node_modules"` | 80-90% |
| Python Projects | `-I "__pycache__,*.pyc"` | 70-80% |

**Key Takeaway**: GPS provides flexible, token-efficient project analysis that adapts to any workflow. Start with minimal context (`--summary`), then progressively expand (`--focus`, `-L`) as needed.
