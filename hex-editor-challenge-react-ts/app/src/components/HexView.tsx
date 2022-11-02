import React from 'react';

interface HexViewProps {
  data: string[];
  hoverIndex: number;
  activeIndex: number;
  hoverIndexChange: React.Dispatch<React.SetStateAction<number>>;
  activeIndexChange: React.Dispatch<React.SetStateAction<number>>;
}

export default function HexView(props: HexViewProps) {
  function HandleKeyDown(ev: React.KeyboardEvent<HTMLDivElement>) {
    if (props.activeIndex > -1 && (ev.ctrlKey || ev.metaKey) && ev.key === 'c') {

      navigator.clipboard.writeText(props.data[props.activeIndex]).then(function () {
        // Copied the value to clipboard
      }, function (err) {
        console.error('Could not copy text: ', err);
      });
    }
  }

  return (
    <div className="hex-view" onKeyDown={HandleKeyDown} tabIndex={0}>{
      props.data.map((x, i) =>
        <span
          onMouseOver={() => { props.hoverIndexChange(i) }}
          onMouseLeave={() => { props.hoverIndexChange(-1) }}
          onMouseUp={() => { props.activeIndexChange(i) }}
          className={`${props.hoverIndex === i ? 'hover' : ''}${props.activeIndex === i ? ' active' : ''}`}
          key={i}>{x}</span>)
    }</div>
  );
}
