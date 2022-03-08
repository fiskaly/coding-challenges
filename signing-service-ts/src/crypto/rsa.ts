import crypto, { KeyFormat, RSAKeyPairOptions } from 'crypto';

// DTO which holds RSA private and public keys
export interface RsaKeyPair {
  public: string;
  private: string;
}

// TODO: generate RSA key pair

