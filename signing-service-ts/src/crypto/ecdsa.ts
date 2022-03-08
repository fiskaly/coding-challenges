import crypto, { ECKeyPairOptions, KeyFormat } from 'crypto';

// DTO which holds EC private and public keys
export interface EcKeyPair {
  public: string;
  private: string;
}

// TODO: generate EC key pair
