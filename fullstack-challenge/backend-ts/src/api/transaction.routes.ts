import { Router } from 'express';

const router = Router();

// TODO: Implement transaction endpoints
//
// POST /api/v1/transactions
//   - Sign data with a signature device
//   - Build secured data: <counter>_<data>_<last_signature>
//   - Increment device counter
//   - Store transaction in database
//   - Return transaction information
//
// GET /api/v1/transactions
//   - List all transactions
//   - Support filtering by deviceId query parameter

export default router;
