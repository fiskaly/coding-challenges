import express from 'express';
import cors from 'cors';
import { createServer } from './api/server';

const PORT = process.env.PORT || 8080;

const app = express();

// Middleware
app.use(cors());
app.use(express.json());

// Create and setup server
const server = createServer(app);

// Start server
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
