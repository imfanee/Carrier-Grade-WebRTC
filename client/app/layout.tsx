/**
 * Root layout for the Carrier-Grade WebRTC PWA client.
 * By:- Faisal Hanif | imfanee@gmail.com
 */

import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Carrier-Grade WebRTC',
  description: 'PWA client for real-time WebRTC communication',
  manifest: '/manifest.json',
  themeColor: '#0ea5e9',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <link rel="manifest" href="/manifest.json" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </head>
      <body>{children}</body>
    </html>
  );
}
