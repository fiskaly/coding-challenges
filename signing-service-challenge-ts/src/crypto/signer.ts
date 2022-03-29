export interface Signer {
  sign: (dataToBeSigned: string[]) => string[] | Error;
}
