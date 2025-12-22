export interface Transaction {
  id: string;
  deviceId: string;
  counter: number;
  timestamp: Date;
  data: string;
  signature: string;        // Base64 encoded
  signedData: string;       // The full string that was signed
  previousSignature: string; // Base64 encoded or base64(device.id) for first transaction
}

export interface CreateTransactionInput {
  deviceId: string;
  data: string;
}

// TODO: Add validation and business logic for Transaction
