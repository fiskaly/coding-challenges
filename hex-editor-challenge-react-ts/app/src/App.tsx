import React from 'react';
import readFile from "./lib/readFile";
import HexViewer from "./components/HexViewer";

function App() {
  const [file, setFile] = React.useState<null | string | Uint8Array>(null);
  
  const updateFileState = async (e: React.FormEvent<HTMLInputElement>) => {
    const result = await readFile(e);
    setFile(result);
  }
  
  const fileInputWrapper = {
    display: "flex",
    gap: ".4em",
    padding: "1em 2em",
    width: "100%",
    border: "3px dotted #bbc",
    background: "#dedeef",
    margin: "0 0 .5em"
  };

  function toArray(data: string | Uint8Array): number[] {
    return typeof data === 'string' ? Array.from(data).map(el => el.charCodeAt(0)) : Array.from(data);
  }

  function arrayToHex(data: number[]): string[] {
    return data.map(num => num.toString(16).toUpperCase().padStart(2, '0'))
  }

  const renderComponents = () => {
    if (!file) {
      return (
        <input
          name="file"
          type="file"
          role="button"
          data-testid="file"
          onInput={updateFileState}
          style={fileInputWrapper}
        />
      )
    }

    const isBinary = typeof file !== 'string';
    return (
      <>
        <div style={fileInputWrapper}>
          <span>Loaded {isBinary ? 'binary' : 'text'} file</span>
          {' '}
          <button onClick={() => setFile(null)}>Reset</button>
        </div>
        <HexViewer data={arrayToHex(toArray(file))} />
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
