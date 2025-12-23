# Signature Service - Fullstack Coding Challenge

## Instructions

This challenge is part of the fullstack software engineering interview process at fiskaly.

If you see this challenge, you've passed the first round of interviews and are now at the second and last stage.

We would like you to attempt the challenge below. You will then be able to discuss your solution in the skill-fit interview with two of our colleagues from the development department.

**The quality of your code is more important to us than the quantity.**

### Project Setup

For the challenge, we provide you with two implementation options. **Choose ONE backend option** and implement the frontend:

**Option A: Go Backend**
- Go project structure with basic HTTP server setup
- Encoding/decoding utilities for different key types (RSA, ECC)
- Key generation algorithms
- PostgreSQL database schema
- Basic domain models and repository pattern examples

**Option B: TypeScript Backend**
- TypeScript project with Express setup
- Crypto utilities for RSA and ECC signing
- Drizzle ORM configuration
- Database schema and migrations
- Basic API structure

**Frontend (Required for both options):**
- Next.js 15+ App Router scaffold
- React 19
- Tailwind CSS configuration
- Vitest + Testing Library setup
- API client utilities
- Component examples with TypeScript

You can use these foundations as-is or modify them to fit your design. The goal is to evaluate your architectural decisions and implementation quality, not to test your ability to set up boilerplate.

### Prerequisites & Tooling

**Backend (choose one):**
- Go 1.21+ (Option A), OR
- Node.js 20+ (Option B - TypeScript)

**Frontend (required):**
- Next.js 15+

**Infrastructure:**
- Docker & Docker Compose
- PostgreSQL 15+

### The Challenge

The goal is to implement an API service with a web frontend that allows customers to create **signature devices** and use them to sign transaction data.

#### Domain Description

##### Signature Device

The service can manage multiple **signature devices**. Each device is identified by a unique identifier (UUID). For this challenge, assume a single organization is using the system - you do not need to implement user management or multi-tenancy.

When creating a signature device, the client must specify:
- **Algorithm**: The signature algorithm (`RSA` or `ECC`)
- **Label**: A human-readable name (e.g., "Device #1", "Store POS")

During creation, the system must:
- Generate a new cryptographic key pair (public key & private key)
- Initialize a `signature_counter` at 0
- Store the device with status `active`

The signature device maintains:
- `id` (UUID)
- `label` (string)
- `algorithm` ('RSA' | 'ECC')
- `publicKey` (PEM format for client verification)
- `privateKey` (securely stored, used for signing)
- `signatureCounter` (monotonically increasing integer)
- `status` ('active' | 'deactivated')
- `createdAt` (timestamp)

##### Transaction Signing

A **transaction** represents data that must be cryptographically signed.

When creating a transaction, the client provides:
- `deviceId`: Which signature device should sign this transaction
- `data`: The raw data to be signed (string)

The signing process follows this algorithm:

1. **Build the secured data string:**
   ```
   For counter = 0:
   secured_data = "0_<data>_<base64(device.id)>"
   
   For counter > 0:
   secured_data = "<counter>_<data>_<last_signature>"
   ```

2. **Sign the secured data:**
   ```
   signature = sign(secured_data, device.privateKey)
   signature_base64 = base64(signature)
   ```

3. **Store the transaction:**
   ```
   {
     id: UUID,
     deviceId: string,
     counter: current_counter,
     timestamp: ISO8601,
     data: string,
     signature: signature_base64,
     signedData: secured_data,
     previousSignature: previous_signature_or_base64_device_id
   }
   ```

4. **Increment the counter:**
   ```
   device.signatureCounter += 1
   ```

This creates a cryptographic chain where each signature depends on the previous one, making it tamper-evident.

#### Backend API Requirements

Implement a RESTful HTTP API with the following endpoints:

**Signature Device Management:**
```
POST   /api/v1/devices
  Body: { algorithm: 'RSA' | 'ECC', label: string }
  Response: { id, label, algorithm, publicKey, signatureCounter, createdAt }

GET    /api/v1/devices
  Response: [ { id, label, algorithm, signatureCounter, createdAt }, ... ]
```

**Transaction Management:**
```
POST   /api/v1/transactions
  Body: { deviceId, data }
  Response: { id, deviceId, counter, timestamp, signature, signedData }

GET    /api/v1/transactions
  Response: [ { id, deviceId, counter, timestamp, data, signature }, ... ]
```

#### Frontend Requirements

Build a web dashboard using **Next.js 15+ with App Router** that provides:

**Pages:**

1. **Device Management (`/` or `/devices`)**
   - Display all signature devices in a grid/list
   - Show: label, algorithm, signature count, creation date
   - Button/form to create new device (inline or modal)
   - Form fields: algorithm selector (RSA/ECC) and label input

2. **Transactions (`/transactions`)**
   - Form to create new transaction: Select device, input data to sign
   - Display result after signing: signature and signed data
   - Table of all signed transactions with filtering by device
   - Display: counter, timestamp, data preview, signature

**UI Requirements:**
- Responsive design with Tailwind CSS
- Loading states for async operations
- Error handling with user-friendly messages
- Form validation

**Components to implement (suggested):**
- `DeviceList` - Display all devices
- `CreateDeviceForm` - Form to create new device
- `TransactionForm` - Form to sign data
- `TransactionTable` - List signed transactions

#### Performance & Caching

Your solution must be designed with **optimized performance** in mind, focusing on efficient request processing between client and server.

**Requirements:**

You are expected to implement **caching approaches** to achieve **low latency** across the application. Consider where caching can improve performance and choose appropriate strategies for your architecture.

**Documentation:**

Add a section to this README titled `## Caching Strategy` that includes:
- **Caching approaches chosen**: Describe what caching mechanisms you implemented and where
- **Implementation locations**: Specify which layers (backend, frontend, database, etc.) use caching
- **TTL and strategy justification**: Explain why you chose this/these cache strategie/s, specific TTL values and invalidation strategies
- **Data consistency considerations**: Discuss how you handle cache consistency, especially for critical data like signature counters

**Evaluation Criteria:**

Your caching implementation will be evaluated on:
- **Correctness of cache invalidation**: Ensuring stale data doesn't cause incorrect behavior
- **Justification of strategies**: Clear reasoning for chosen TTLs and caching mechanisms
- **Data consistency handling**: Proper management of cache coherence with the database

#### QA / Testing

**Correctness is critical**. Your implementation must be verifiable.

**Required:**
- Tests demonstrating correct signature format
- Tests verifying monotonic counter behavior

**Recommended:**
- Unit tests for signing logic (backend)
- Component tests for forms (frontend)
- API endpoint tests
- E2E tests for critical flows (optional)

Provide a clear way to run your tests (e.g., `npm test`, `go test ./...`).

#### Technical Constraints & Considerations

- **Concurrency**: The system will be used by multiple concurrent clients. The `signatureCounter` must be strictly monotonically increasing with no gaps, even under concurrent load.
- **Thread Safety**: Use appropriate mechanisms (mutexes in Go, database transactions in TypeScript) to ensure counter consistency.
- **Algorithm Extensibility**: Design the signing mechanism to allow easy addition of new algorithms (e.g., Ed25519) without modifying core domain logic.
- **Persistence**: Use PostgreSQL for all data. The provided schema is a starting point - modify as needed.
- **Error Handling**: Properly handle and communicate errors (invalid device, invalid data, etc.).
- **API Design**: Follow RESTful conventions and HTTP standards.

#### Docker Setup

We provide Docker Compose configurations for easy development setup.

**Quick Start (Run everything from project root):**
```bash
# From fullstack-challenge/ directory
docker-compose up

# This will start:
# - PostgreSQL on localhost:5432
# - Backend (Go by default) on localhost:8080
# - Frontend on localhost:3000
```

**To use TypeScript backend instead:**
1. Edit `docker-compose.yml` in project root
2. Comment out the `backend-go` service
3. Uncomment the `backend-ts` service
4. Update frontend's `depends_on` to use `backend-ts`
5. Run `docker-compose up`

**Individual services:**

Backend Go:
```bash
cd backend-go
docker-compose up
# Backend on http://localhost:8080
# PostgreSQL on localhost:5432
```

Backend TypeScript:
```bash
cd backend-ts
docker-compose up
# Backend on http://localhost:8080
# PostgreSQL on localhost:5432
```

Frontend:
```bash
cd frontend
docker-compose up
# Frontend on http://localhost:3000
# Configure NEXT_PUBLIC_API_URL in .env.local
```

Or run everything with:
```bash
# From project root
docker-compose up
```

### AI Tools Usage

The use of AI tools (GitHub Copilot, ChatGPT, Claude, etc.) to aid in completing the challenge is **permitted and encouraged**. At fiskaly, we embrace AI as a productivity tool.

**However, you must document your AI usage clearly.** Add a section to this README with:

```markdown
## AI Tools Used

**GitHub Copilot:**
- Used for: [describe what you used it for]
- Why: [explain your reasoning]
- Modifications made: [describe significant changes to AI-generated code]

**ChatGPT/Claude:**
- Used for: [describe what you used it for]
- Why: [explain your reasoning]
- What you learned: [key insights or decisions informed by AI]

**Other Tools:**
- [tool name]: [usage and reasoning]
```

### Credits

This challenge is heavily influenced by the regulations for `KassenSichV` (Germany) as well as the `RKSV` (Austria) and our solutions for them.

---

**Good luck! We're excited to see your solution. 🚀**