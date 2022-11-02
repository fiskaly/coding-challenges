import React from 'react';
import { useState } from 'react';
import "./HexViewer.css";
import HexView from "./HexView";
import TextView from "./TextView";

interface HexViewerProps {
  data: string[];
}

export default function HexViewer(props: HexViewerProps) {
  const [hovered, setHovered] = useState<number>(-1);
  const [active, setActive] = useState<number>(-1);

  return (
    <div className="hex-viewer" data-testid="viewer">
      <HexView data={props.data} hoverIndex={hovered} hoverIndexChange={setHovered} activeIndex={active} activeIndexChange={setActive} />
      <TextView data={props.data} hoverIndex={hovered} hoverIndexChange={setHovered} activeIndex={active} activeIndexChange={setActive} />
    </div>
  );
}
