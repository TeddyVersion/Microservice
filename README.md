# Gosmart Golang Monorepo

This is a Golang monorepo boilerplate following standard Go project layout.

## Structure
- `/cmd` - Main applications for this project
    - `auth-service` - User registration, login, PIN, biometrics
    - `account-service` - User accounts, balances, statements
    - `transfer-service` - Money transfers
    - `chat-service` - Chat banking, P2P
    - `billpay-service` - Bill payments
    - `topup-service` - Airtime/data top-ups
    - `merchant-service` - Merchant payments
    - `profile-service` - Profile management
    - `loan-service` - Loan management
    - `finance-service` - Personal finance
    - `miniapp-service` - Mini-apps
    - `notification-service` - Notifications
- `/pkg` - Shared libraries/utilities (JWT, Kafka, DB, etc.)
- `/internal` - Private code

## Running a Service

```bash
go run ./cmd/<service-name>
# Example:
go run ./cmd/auth-service
```

Each service exposes a `/healthz` endpoint for health checks.

## Extending
- Add new microservices under `/cmd`.
- Place shared code in `/pkg`.
- Place internal-only code in `/internal`.

---

For more details, see the technical specification document.
