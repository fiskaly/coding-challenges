# React Challenge

Welcome challenger!
This challenge is here to see your React skills.

## Instructions

Create a **Hex viewer** component and display it in storybook.
Basically a [Hex editor](https://en.wikipedia.org/wiki/Hex_editor) without the ability to change its content.

- The component can retrieve either a `string` or `Uint8Array`.
- The Hex viewer must be responsive and each part of the line must match (hex value must match with the text value on the other side).
- The component must be able to display any sort of blob file and displays the bytes as hex on the left and if the
  characters are printable (like numbers and letters) on the right.
- Non-printable characters should be replaced with a special character and displayed on the right.
- The component should be viewable with multiple variations on storybook.

Optional features:
- It would be nice to be able to select one part of the hex viewer and automatically the other part gets highlighted as well.
- The possibility to then copy the selected hex value or text (depending on what was selected) would be nice as well.
- Add the ability to display a `Uint16Array`.

We wish you all the best and are looking forward to your results!
Cheers
