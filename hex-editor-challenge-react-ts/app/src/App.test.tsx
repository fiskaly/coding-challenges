import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

test('renders file button', () => {
  render(<App />);
  const fileElement = screen.getByRole('button');
  expect(fileElement).toBeInTheDocument();
});
