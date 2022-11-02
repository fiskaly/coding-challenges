# Hex Editor Progress

Author: Sławomir Amielucha

## The objective

The primary objective of the challenge is to show the capability of working with a React application.

Knowing that, I have placed the main focus on the overall app architecture rather than on producing a fully-fledged hex viewer.

## Meeting the requirements

I consider my solution to meet the primary requirements of the challenge.

From the list of optional features I have selected two: 

1. selecting an item from one side of the viewer also selects the matching character on the opposite side. This, of course, could be extended to allow for a range selection using the `Window.getSelection()` method.
1. Copying the selected value. Analogically, this could work with a range selection, composing a string from selected nodes.

## The approach

Due to time constraints I have built a MVP version of a viewer. Styling has been kept to minimum, making the viewer as non-opinionated as possible and suitable to be used within any design system.

Styles file supports CSS variables in order to easily adjust key aspects of the application.

Performance has not been treated as a priority and obviously can be greatly improved, avoiding excessive updates or avoiding expensive operations. Hover function, for example, has been kept only as a proof of concept and in a real-world application should be implemented using a more performant method.

I added tests validating key functionality of the application.
