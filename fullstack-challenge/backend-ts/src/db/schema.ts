import { pgTable, uuid, text, integer, timestamp } from 'drizzle-orm/pg-core';
import { relations } from 'drizzle-orm';

// Devices table (signature devices)
export const devices = pgTable('devices', {
  id: uuid('id').primaryKey().defaultRandom(),
  label: text('label').notNull(),
  algorithm: text('algorithm', { enum: ['RSA', 'ECC'] }).notNull(),
  publicKey: text('public_key').notNull(),
  privateKey: text('private_key').notNull(),
  signatureCounter: integer('signature_counter').default(0).notNull(),
  status: text('status', { enum: ['active', 'deactivated'] }).default('active').notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
});

// Transactions table
export const transactions = pgTable('transactions', {
  id: uuid('id').primaryKey().defaultRandom(),
  deviceId: uuid('device_id').references(() => devices.id).notNull(),
  counter: integer('counter').notNull(),
  timestamp: timestamp('timestamp').defaultNow().notNull(),
  data: text('data').notNull(),
  signature: text('signature').notNull(),
  signedData: text('signed_data').notNull(),
  previousSignature: text('previous_signature').notNull(),
});

// Relations
export const devicesRelations = relations(devices, ({ many }) => ({
  transactions: many(transactions),
}));

export const transactionsRelations = relations(transactions, ({ one }) => ({
  device: one(devices, {
    fields: [transactions.deviceId],
    references: [devices.id],
  }),
}));

// Type exports for use in application
export type Device = typeof devices.$inferSelect;
export type NewDevice = typeof devices.$inferInsert;
export type Transaction = typeof transactions.$inferSelect;
export type NewTransaction = typeof transactions.$inferInsert;
