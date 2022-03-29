import { KeyPair } from './generation';
import crypto from 'crypto';

export default function sign(dataToBeSigned: string, keyPair: KeyPair): string {
  try {
    const cipher = crypto.createCipheriv('aes-256-gcm', crypto.scryptSync(keyPair.private, 'fiskaly is AWESOME!', 32), crypto.randomBytes(16));
    let encrypted = cipher.update(dataToBeSigned);
    encrypted = Buffer.concat([encrypted, cipher.final()]);
    return encrypted.toString('base64');
  } catch (e) {
    console.log(e);
    throw e;
  }
}
