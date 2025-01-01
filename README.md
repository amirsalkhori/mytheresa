# MyTheresa Promotions Test

## Overview
This repository contains the implementation of the MyTheresa Promotions Test. The project demonstrates efficient handling of large datasets with cursor-based pagination, dynamic discount application, and performance optimization using Redis caching.

---

## Key Features
- **Cursor-Based Pagination:** Efficient navigation for datasets with over 20,000 products.
- **Dynamic Discount Application:** Real-time calculation of discounts with priority rules.
- **Redis Caching:** Optimization of discount retrieval for reduced database load.
- **Filtering:** Supports category and price filters for flexible product queries.
- **Scalability:** Designed to handle large datasets while maintaining low response times.

---

## Quickstart

### Prerequisites
- Docker
- Docker Compose
- Makefile (optional for easier commands)
- Golang 1.20+

### Running the Project
1. Clone the repository:
   ```bash
   git clone https://github.com/amirsalkhori/mytheresa.git
   cd mytheresa
   ```
2. Start the application:
   ```bash
   make build-up
   ```
3. Seed the database with sample data:
   ```bash
   make seed-products
   ```
4. Access the API at `http://localhost:8000/v1/products`.

### Stopping the Application
To stop and clean up the Docker containers, run:
```bash
make down
```

---

## API Endpoints

### `GET /v1/products`
Retrieve a list of products with optional filters and discounts applied.

#### Query Parameters:
- `category` (optional): Filter by category.
- `priceLessThan` (optional): Filter products by price (before discounts).
- `next` (optional): Fetch the next page using the cursor.
- `prev` (optional): Fetch the previous page using the cursor.

#### Example Response:
```json
{
  "products": [
    {
      "sku": "000001",
      "name": "BV Lean leather ankle boots",
      "category": "boots",
      "price": {
        "original": 89000,
        "final": 62300,
        "discount_percentage": "30%",
        "currency": "EUR"
      }
    }
  ],
  "pagination": {
    "pageSize": 5,
    "next": "abcdef123",
    "prev": ""
  }
}
```
The full swagger doc is available in [docs/swagger/swagger.yaml](docs/swagger/swagger.yaml)

---

## Documentation
Detailed documentation is available in the `docs/` folder:
- [Implementation Details](docs/implementations.md)
- [Performance Testing](docs/performance.md)
- [Challenges and Solutions](docs/challenges-and-solutions.md)
- [Future Improvements](docs/future.md)

---

## Project Structure
```
- cmd/                      # Application entry points
  - mytheresa/main.go       # Main app starter
  - seeder/main.go          # Seeder for populating data
- config/                   # Loading default app configurations
  - configs.go              # Configuration file for app settings
- deploy/                   # Deployment configurations
- docs/                     # Documentation
  - implementations.md      # Implementation details
  - performance.md          # Performance testing
  - future.md               # Future improvements
  - challenges-and-solutions.md  # Challenges faced and solutions
  - task.md                 # Original task definition
  - assets/                 # Images and diagrams related to docs
- internal/                 # Application internal services
  - app/                    # Application layer
    - application.go        # Core orchestrator
    - dto/                  # Data Transfer Objects (input/output formats)
  - domain/                 # Core business logic
    - derrors/              # Domain-specific errors
  - handler/                # Adapters for external interactions (HTTP, gRPC, etc.)
  - infra/                  # Infrastructure layer
    - db/
      - mysql/
        - config/mysql.go   # MySQL configuration
        - model/            # MySQL models
        - repository/       # MySQL repositories
      - redis/
        - redis.go          # Redis configuration and caching
  - ports/                  # Interfaces (ports) defining contracts
  - services/               # Business logic implementations
    - test/mocks/           # Mocks for unit testing
- migrations/               # Migrations for db
- docker-compose.yml        # Multi-container orchestration
- migrations/               # Database migrations
- Makefile                  # Helper commands for setup
- README.md                 # Main project entry point

```

---

## Testing
To run tests, use:
```bash
make test-ginkgo
```
This will execute all unit and integration tests and provide a summary of results.

---

## Future Improvements
- **Hybrid Pagination:** Combine traditional and cursor-based pagination for initial pages.
- **Elasticsearch Integration:** Add advanced search capabilities.
- **Filter-Based Caching:** Optimize filters by caching their results.

For more details, see [Future Improvements](docs/future.md).

---

## License
This project is licensed under the MIT License.
