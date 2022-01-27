# Web Challenge

Welcome challenger!
The goal of this challenge is to see you ability in creating UI components with React.

## Instructions

Create a **Hex Viewer** component and display it in a storybook installation.
Basically a [Hex editor](https://en.wikipedia.org/wiki/Hex_editor) without the ability to change its content.

- The Hex viewer must be responsive and each part of the line must match (hex value must match with the text value on the other side).
- The component must be able to display any sort of blob file and displays the bytes as hex on the left and if the
  characters are printable (like numbers and letters) on the right.
- Non-printable characters should be replaced with a special character and displayed on the right.
- It would be nice to be able to select one part of the hex viewer and automatically the other part gets highlighted as well.
- Also would be nice to then be able to copy the selected hex value or text (depending what was selected).

The component should then be viewable with multiple variations on storybook.