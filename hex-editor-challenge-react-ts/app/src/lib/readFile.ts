import React from "react";

export default function readFile(
  e: React.FormEvent<HTMLInputElement>
): Promise<null | string | Uint8Array> {
  return new Promise((resolve, reject) => {
    const files = (e.target as HTMLInputElement).files;
    if (files == null || !files.length) {
      return resolve(null);
    }
    if (FileReader == null) {
      return reject('Your browser is outdated!');
    }

    const file = files[0];
    const reader = new FileReader();
    reader.addEventListener('load', (e) => {
      const raw = (e.target as FileReader).result;
      if (raw != null && typeof raw !== 'string') {
        return resolve(new Uint8Array(raw));
      }
      return resolve(raw);
    });

    if (/\.txt$/.test(file.name)) {
      reader.readAsText(file);
    } else {
      reader.readAsArrayBuffer(file);
    }
  })
}
