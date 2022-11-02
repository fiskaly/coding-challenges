import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import App from './App';

test('renders file button', () => {
  render(<App />);
  const fileElement = screen.getByRole('button');
  expect(fileElement).toBeInTheDocument();
});

test(`doesn't render the Hex Viewer without a file`, () => {
  render(<App />);
  const hexViewerWrapper = document.querySelector('.hex-viewer');
  expect(hexViewerWrapper).not.toBeInTheDocument();
});

test(`uploads the file`, async () => {
  const file = new File(['example text'], 'hello.text', { type: 'text/plain' })

  render(<App />);
  const input: HTMLInputElement = screen.getByTestId('file');

  userEvent.upload(input, file);

  await waitFor(() => expect(input?.files![0]).toBe(file))
  expect(input?.files![0]).toBe(file)
  expect(input.files?.item(0)).toBe(file)
  expect(input.files).toHaveLength(1)
});

test(`renders correct hex and text content`, async () => {
  const fileText = 'example text';
  const expectedFileHex = '6578616D706C652074657874';
  const file = new File([fileText], 'hello.text', { type: 'text/plain' })

  render(<App />);
  const input: HTMLInputElement = screen.getByTestId('file');

  userEvent.upload(input, file);
  const hexViewer = await screen.findByTestId('viewer');

  expect([...hexViewer.querySelectorAll('.hex-view > span')].map(el => el.innerHTML).join('')).toBe(expectedFileHex)
  expect([...hexViewer.querySelectorAll('.text-view > span')].map(el => el.innerHTML).join('')).toBe(fileText)
});
