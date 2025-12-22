const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface Device {
  id: string;
  label: string;
  algorithm: 'RSA' | 'ECC';
  publicKey: string;
  signatureCounter: number;
  createdAt: string;
}

export interface Transaction {
  id: string;
  deviceId: string;
  counter: number;
  timestamp: string;
  data: string;
  signature: string;
  signedData: string;
}

export interface ApiResponse<T> {
  data: T;
}

export interface ApiError {
  errors: string[];
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    if (!response.ok) {
      const error: ApiError = await response.json();
      throw new Error(error.errors.join(', ') || 'API request failed');
    }

    const result: ApiResponse<T> = await response.json();
    return result.data;
  }

  // Device endpoints
  async getDevices(): Promise<Device[]> {
    return this.request<Device[]>('/api/v1/devices');
  }

  async createDevice(data: {
    label: string;
    algorithm: 'RSA' | 'ECC';
  }): Promise<Device> {
    return this.request<Device>('/api/v1/devices', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Transaction endpoints
  async getTransactions(deviceId?: string): Promise<Transaction[]> {
    const url = deviceId
      ? `/api/v1/transactions?deviceId=${deviceId}`
      : '/api/v1/transactions';
    return this.request<Transaction[]>(url);
  }

  async createTransaction(data: {
    deviceId: string;
    data: string;
  }): Promise<Transaction> {
    return this.request<Transaction>('/api/v1/transactions', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }
}

export const apiClient = new ApiClient(API_URL);
