# React Challenge

Welcome challenger!
This challenge is here to see your React skills.

## Setup

Jump into the `app` directory and install the dependencies. The scaffold was created with
`npx create-react-app my-app --template typescript`.

You can start the project with `npm start`.

## Instructions

Create a **Hex viewer** component (see `src/components/HexViewer.tsx`).
Basically a [Hex editor](https://en.wikipedia.org/wiki/Hex_editor) without the ability to change its content.

- The component can retrieve either a `string` or `Uint8Array`.
- The Hex viewer must be responsive and each part of the line must match (hex value must match with the text value line by line).
- The component must be able to display the bytes as hex and readable text side by side.
- Non-readable characters should be replaced with a special character.

Optional features:
- It would be nice to be able to select one part of the hex viewer and automatically the other part gets highlighted as well.
- The possibility to then copy the selected hex value or text (depending on what was selected) would be nice as well.
- Add the ability to display `Uint16Array` and `ArrayBuffer`.

---

Feel free to alter the structure of this challenge if you see fit.

We wish you all the best and are looking forward to your results!
Cheers
