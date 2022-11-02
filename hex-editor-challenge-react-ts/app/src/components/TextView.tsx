import React from 'react';

interface TextViewProps {
  data: string[];
  hoverIndex: number;
  activeIndex: number;
  hoverIndexChange: React.Dispatch<React.SetStateAction<number>>;
  activeIndexChange: React.Dispatch<React.SetStateAction<number>>;
}

export default function TextView(props: TextViewProps) {
  function isReadable(hex: string): boolean {
    let codepoint = parseInt(hex, 16);

    // is not a low ASCII character and is not DEL
    return codepoint >= 32 && codepoint !== 127;
  }

  function HandleKeyDown(ev: React.KeyboardEvent<HTMLDivElement>) {
    if (props.activeIndex > -1 && (ev.ctrlKey || ev.metaKey) && ev.key === 'c' && isReadable(props.data[props.activeIndex])) {
      navigator.clipboard.writeText(String.fromCharCode(parseInt(props.data[props.activeIndex], 16))).then(function () {
        // Copied the value to clipboard
      }, function (err) {
        console.error('Could not copy text: ', err);
      });
    }
  }

  return (
    <div className="text-view" onKeyDown={HandleKeyDown} tabIndex={0}>{
      props.data.map((hex, i) =>
        <span
          onMouseOver={() => { props.hoverIndexChange(i) }}
          onMouseLeave={() => { props.hoverIndexChange(-1) }}
          onMouseUp={() => { props.activeIndexChange(i) }}
          className={`${props.hoverIndex === i ? 'hover' : ''}${props.activeIndex === i ? ' active' : ''}`}
          key={i}
          >{isReadable(hex) ? String.fromCharCode(parseInt(hex, 16)) : <span className="faded">&middot;</span>}</span>)
    }</div>
  );
}