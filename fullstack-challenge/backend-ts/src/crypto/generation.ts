import * as crypto from 'crypto';

export interface KeyPair {
  publicKey: string;  // PEM format
  privateKey: string; // PEM format
}

export class RSAKeyGenerator {
  /**
   * Generate a new RSA key pair
   * Returns keys in PEM format
   */
  generate(): KeyPair {
    const { publicKey, privateKey } = crypto.generateKeyPairSync('rsa', {
      modulusLength: 2048,
      publicKeyEncoding: {
        type: 'spki',
        format: 'pem'
      },
      privateKeyEncoding: {
        type: 'pkcs8',
        format: 'pem'
      }
    });

    return { publicKey, privateKey };
  }
}

export class ECCKeyGenerator {
  /**
   * Generate a new ECC key pair using P-256 curve
   * Returns keys in PEM format
   */
  generate(): KeyPair {
    const { publicKey, privateKey } = crypto.generateKeyPairSync('ec', {
      namedCurve: 'P-256',
      publicKeyEncoding: {
        type: 'spki',
        format: 'pem'
      },
      privateKeyEncoding: {
        type: 'pkcs8',
        format: 'pem'
      }
    });

    return { publicKey, privateKey };
  }
}
