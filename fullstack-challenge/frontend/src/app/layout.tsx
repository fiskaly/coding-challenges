import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Signature Service",
  description: "Fullstack coding challenge - Signature service dashboard",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
