import crypto, { ECKeyPairOptions, RSAKeyPairOptions, KeyFormat } from 'crypto';

export interface KeyPair {
  public: string;
  private: string;
}

export const generateEcKeyPair = (): Promise<KeyPair> => {
  const options: ECKeyPairOptions<KeyFormat, KeyFormat> = {
    namedCurve: 'secp256k1',
    publicKeyEncoding: {
      type: 'spki',
      format: 'der'
    },
    privateKeyEncoding: {
      type: 'pkcs8',
      format: 'der'
    }
  };

  return new Promise<KeyPair>((resolve, reject) => {
    try {
      crypto.generateKeyPair('ec', options, (err, publicKey, privateKey) => {
        if (err) {
          console.log(`Encountered an error during EC key pair generation: ${err.message}`);
        }

        const keyPair: KeyPair = {
          public: publicKey.toString(),
          private: privateKey.toString(),
        };
        // const keyPair: CryptoKeyPair = {
        //   public: publicKey,
        //   private: privateKey,
        // };
        console.log(`Created EC key pair: ${keyPair}`);
        resolve(keyPair);
      });
    } catch (e) {
      console.log(`Failed to generated EC key pair`);
      reject(e);
    }
  });
}

export const generateRsaKeyPair = (): Promise<KeyPair> => {
  const options: RSAKeyPairOptions<KeyFormat, KeyFormat> = {
    modulusLength: 2048,
    publicExponent: 0x10101,
    publicKeyEncoding: {
      type: 'pkcs1',
      format: 'der'
    },
    privateKeyEncoding: {
      type: 'pkcs8',
      format: 'der',
      cipher: 'aes-192-cbc',
      passphrase: 'fiskaly is AWESOME'
    }
  };

  return new Promise<KeyPair>((resolve, reject) => {
    try {
      crypto.generateKeyPair('rsa', options, (err, publicKey, privateKey) => {
        if (err) {
          console.log(`Encountered an error during RSA key pair generation: ${err.message}`);
        }

        const keyPair: KeyPair = {
          public: publicKey.toString(),
          private: privateKey.toString(),
        };
        console.log(`Created RSA key pair: ${keyPair}`);
        resolve(keyPair);
      });
    } catch (e) {
      console.log(`Failed to generated RSA key pair`);
      reject(e);
    }
  });
}

export default function generateKeyPair() {
  // TODO: implement key pair generation
}
