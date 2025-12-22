import * as crypto from 'crypto';

export interface Signer {
  sign(data: string): string; // Returns base64 encoded signature
}

export class RSASigner implements Signer {
  constructor(private privateKey: string) {}

  sign(data: string): string {
    const sign = crypto.createSign('SHA256');
    sign.update(data);
    sign.end();
    const signature = sign.sign(this.privateKey);
    return signature.toString('base64');
  }
}

export class ECCSigner implements Signer {
  constructor(private privateKey: string) {}

  sign(data: string): string {
    const sign = crypto.createSign('SHA256');
    sign.update(data);
    sign.end();
    const signature = sign.sign(this.privateKey);
    return signature.toString('base64');
  }
}

/**
 * Factory function to create a signer based on algorithm type
 */
export function createSigner(algorithm: 'RSA' | 'ECC', privateKey: string): Signer {
  switch (algorithm) {
    case 'RSA':
      return new RSASigner(privateKey);
    case 'ECC':
      return new ECCSigner(privateKey);
    default:
      throw new Error(`Unsupported algorithm: ${algorithm}`);
  }
}
