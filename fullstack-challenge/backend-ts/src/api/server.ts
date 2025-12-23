import { Express, Request, Response } from 'express';

export interface ApiResponse<T> {
  data: T;
}

export interface ApiError {
  errors: string[];
}

export const createServer = (app: Express): Express => {
  // Health check endpoint
  app.get('/api/v1/health', (req: Request, res: Response) => {
    const response: ApiResponse<{ status: string }> = {
      data: { status: 'ok' }
    };
    res.json(response);
  });

  // TODO: Add device endpoints
  // POST /api/v1/devices - Create a new signature device
  // GET /api/v1/devices - List all devices

  // TODO: Add transaction endpoints
  // POST /api/v1/transactions - Sign data with a device
  // GET /api/v1/transactions - List all transactions

  // Error handling middleware
  app.use((err: Error, req: Request, res: Response, next: any) => {
    console.error(err.stack);
    const errorResponse: ApiError = {
      errors: [err.message || 'Internal server error']
    };
    res.status(500).json(errorResponse);
  });

  return app;
};
