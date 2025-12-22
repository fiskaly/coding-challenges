# Signature Service - TypeScript Backend

This is the TypeScript backend implementation for the signature service challenge.

## Getting Started

### Prerequisites

- Node.js 20+
- Docker & Docker Compose (for PostgreSQL)

### Installation

```bash
npm install
```

### Running with Docker

```bash
docker-compose up
```

The backend will be available at `http://localhost:8080`.

### Running Locally

1. Copy environment variables:
```bash
cp .env.example .env
```

2. Start PostgreSQL:
```bash
docker-compose up postgres
```

3. Run the application:
```bash
npm run dev
```

### Database

This project uses **Drizzle ORM** for database management.

#### Database Setup

The schema is defined in `src/db/schema.ts` using Drizzle ORM.

**Generate migrations from schema:**
```bash
npm run db:generate
```

**Apply migrations to database:**
```bash
npm run db:migrate
```

**Push schema directly to database (development):**
```bash
npm run db:push
```

**Open Drizzle Studio (visual database browser):**
```bash
npm run db:studio
```

#### Schema Overview

**Devices table:**
- `id` (UUID) - Primary key
- `label` (text) - Device label
- `algorithm` (enum: 'RSA' | 'ECC')
- `publicKey` (text) - PEM format public key
- `privateKey` (text) - PEM format private key
- `signatureCounter` (integer) - Monotonic counter
- `status` (enum: 'active' | 'deactivated')
- `createdAt` (timestamp)

**Transactions table:**
- `id` (UUID) - Primary key
- `deviceId` (UUID) - Foreign key to devices
- `counter` (integer) - Transaction counter
- `timestamp` (timestamp)
- `data` (text) - Raw data signed
- `signature` (text) - Base64 signature
- `signedData` (text) - Full signed data string
- `previousSignature` (text) - Previous signature (chain)

#### Using Drizzle in Your Code

```typescript
import { db } from './db/database';
import { devices, transactions } from './db/schema';
import { eq } from 'drizzle-orm';

// Query examples
const allDevices = await db.select().from(devices);
const device = await db.select().from(devices).where(eq(devices.id, deviceId));
const [created] = await db.insert(devices).values(newDevice).returning();
```

## Project Structure

```
backend-ts/
├── src/
│   ├── api/          # HTTP routes and handlers
│   ├── crypto/       # Cryptographic operations (RSA, ECC)
│   ├── domain/       # Domain models and interfaces
│   ├── db/           # Database connection and schemas
│   └── index.ts      # Application entry point
├── drizzle/          # Generated migration files
├── drizzle.config.ts # Drizzle Kit configuration
├── package.json
├── tsconfig.json
└── docker-compose.yml
```

## API Endpoints

### Health Check
```
GET /api/v1/health
```

### TODO: Implement these endpoints

```
POST   /api/v1/devices
GET    /api/v1/devices
POST   /api/v1/transactions
GET    /api/v1/transactions
```

## Testing

```bash
npm test
```

## Building

```bash
npm run build
```

## Running in Production

```bash
npm start
```
