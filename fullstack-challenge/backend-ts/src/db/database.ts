import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import * as schema from './schema';
import { eq, and } from 'drizzle-orm';

// Database connection
const connectionString = process.env.DATABASE_URL || 
  'postgresql://postgres:postgres@localhost:5432/signature_service';

const client = postgres(connectionString);
export const db = drizzle(client, { schema });

// Type-safe database instance
export type Database = typeof db;

// Example query helpers for candidates to implement
export class DeviceRepository {
  // TODO: Implement device persistence methods
  // Example methods to implement:
  //
  // async createDevice(device: schema.NewDevice): Promise<schema.Device> {
  //   const [created] = await db.insert(schema.devices).values(device).returning();
  //   return created;
  // }
  //
  // async getDevices(): Promise<schema.Device[]> {
  //   return db.select().from(schema.devices);
  // }
  //
  // async getDeviceById(id: string): Promise<schema.Device | undefined> {
  //   const [device] = await db.select().from(schema.devices).where(eq(schema.devices.id, id));
  //   return device;
  // }
  //
  // async updateDeviceCounter(id: string, counter: number): Promise<void> {
  //   await db.update(schema.devices)
  //     .set({ signatureCounter: counter })
  //     .where(eq(schema.devices.id, id));
  // }
  //
  // async deactivateDevice(id: string): Promise<void> {
  //   await db.update(schema.devices)
  //     .set({ status: 'deactivated' })
  //     .where(eq(schema.devices.id, id));
  // }
}

export class TransactionRepository {
  // TODO: Implement transaction persistence methods
  // Example methods to implement:
  //
  // async createTransaction(transaction: schema.NewTransaction): Promise<schema.Transaction> {
  //   const [created] = await db.insert(schema.transactions).values(transaction).returning();
  //   return created;
  // }
  //
  // async getTransactions(deviceId?: string): Promise<schema.Transaction[]> {
  //   if (deviceId) {
  //     return db.select().from(schema.transactions)
  //       .where(eq(schema.transactions.deviceId, deviceId));
  //   }
  //   return db.select().from(schema.transactions);
  // }
  //
  // async getLastTransactionForDevice(deviceId: string): Promise<schema.Transaction | undefined> {
  //   const [transaction] = await db.select().from(schema.transactions)
  //     .where(eq(schema.transactions.deviceId, deviceId))
  //     .orderBy(desc(schema.transactions.counter))
  //     .limit(1);
  //   return transaction;
  // }
  //
  // async getTransactionById(id: string): Promise<schema.Transaction | undefined> {
  //   const [transaction] = await db.select().from(schema.transactions)
  //     .where(eq(schema.transactions.id, id));
  //   return transaction;
  // }
}
