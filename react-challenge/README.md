# React Challenge

Welcome challenger!
This challenge is here to see your React skills.

## Instructions

Create a **Hex viewer** component and display it in storybook.
Basically a [Hex editor](https://en.wikipedia.org/wiki/Hex_editor) without the ability to change its content.

- The component can retrieve either a `string` or `Uint8Array`.
- The Hex viewer must be responsive and each part of the line must match (hex value must match with the text value line by line).
- The component must be able to display the bytes as hex and readable text side by side.
- Non-readable characters should be replaced with a special character.
- The component should be viewable with multiple variations/examples on storybook.

Optional features:
- Create a wrapper component that can load binary files (e.g. the example binary files in this folder)
- It would be nice to be able to select one part of the hex viewer and automatically the other part gets highlighted as well.
- The possibility to then copy the selected hex value or text (depending on what was selected) would be nice as well.
- Add the ability to display a `Uint16Array`.

We wish you all the best and are looking forward to your results!
Cheers
