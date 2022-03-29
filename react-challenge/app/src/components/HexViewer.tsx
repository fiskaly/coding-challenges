import React from 'react';

interface HexViewerProps {
  data: string | Uint8Array;
}

export default function HexViewer(props: HexViewerProps) {
  /*
   * This component is the main challenge. You can be wild here and change
   * everything!
   */
  return (
    <pre style={{ overflowWrap: 'break-word', whiteSpace: 'pre-wrap', wordBreak: 'break-all' }}>
      Here comes the HexViewer<br />{ props.data }
    </pre>
  );
}
