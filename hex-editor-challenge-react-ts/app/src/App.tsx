import React from 'react';
import readFile from "./lib/readFile";
import HexViewer from "./components/HexViewer";

function App() {
  const [file, setFile] = React.useState<null | string | Uint8Array>(null);
  const updateFileState = async (e: React.FormEvent<HTMLInputElement>) => {
    const result = await readFile(e);
    setFile(result);
  }

  const renderComponents = () => {
    if (!file) {
      return (
        <input
          name="file"
          type="file"
          role="button"
          onInput={updateFileState}
        />
      )
    }

    const isBinary = typeof file !== 'string';
    return (
      <>
        <div>
          <span>Loaded {isBinary ? 'binary' : 'text'} file</span>
          {' '}
          <button onClick={() => setFile(null)}>Reset</button>
        </div>
        <HexViewer data={file} />
      </>
    )
  };

  return (
    <div className="App">
      {renderComponents()}
    </div>
  );
}

export default App;
