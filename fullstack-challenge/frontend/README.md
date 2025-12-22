# Signature Service - Frontend

This is the Next.js frontend for the signature service challenge.

## Getting Started

### Prerequisites

- Node.js 20+
- Backend server running (Go or TypeScript)

### Installation

```bash
npm install
```

### Running Locally

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Running with Docker

```bash
docker-compose up
```

### Environment Variables

Copy the example file and adjust as needed:

```bash
cp .env.local.example .env.local
```

- `NEXT_PUBLIC_API_URL`: Backend API URL (default: http://localhost:8080)

## Project Structure

```
frontend/
├── src/
│   ├── app/              # Next.js App Router pages
│   │   ├── devices/      # Devices management page
│   │   ├── transactions/ # Transactions page
│   │   ├── layout.tsx    # Root layout
│   │   ├── page.tsx      # Home page
│   │   └── globals.css   # Global styles
│   ├── components/       # React components
│   │   ├── DeviceList.tsx
│   │   ├── CreateDeviceForm.tsx
│   │   ├── TransactionForm.tsx
│   │   └── TransactionTable.tsx
│   └── lib/
│       └── api.ts        # API client
├── package.json
├── next.config.js
├── tailwind.config.js
├── tsconfig.json
├── vitest.config.ts
├── vitest.setup.ts
├── Dockerfile
└── docker-compose.yml
```

## Features to Implement

### Device Management Page (`/devices`)
- [ ] Display all signature devices in a grid/list
- [ ] Show device information (label, algorithm, signature count, creation date)
- [ ] Form to create new devices (algorithm selector + label input)
- [ ] Handle loading and error states

### Transactions Page (`/transactions`)
- [ ] Form to create transactions (select device, input data to sign)
- [ ] Display signature result after signing
- [ ] Table of all transactions with device filtering
- [ ] Show transaction details (counter, timestamp, data, signature)

## Testing

Run tests with:

```bash
npm test
```

Run tests in watch mode:

```bash
npm run test:watch
```

## Building for Production

```bash
npm run build
npm start
```

## Tech Stack

- **Next.js 15+** with App Router
- **React 19**
- **TypeScript**
- **Tailwind CSS**
- **Vitest** for testing
