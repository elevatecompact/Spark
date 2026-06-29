import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Spark",
  description: "Live streaming platform",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
