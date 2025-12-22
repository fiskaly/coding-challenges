export interface Device {
  id: string;
  label: string;
  algorithm: 'RSA' | 'ECC';
  publicKey: string;
  privateKey: string;
  signatureCounter: number;
  createdAt: Date;
}

export interface CreateDeviceInput {
  label: string;
  algorithm: 'RSA' | 'ECC';
}

// TODO: Add validation and business logic for Device
