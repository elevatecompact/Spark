"use client";

import Link from "next/link";
import type { ReactNode } from "react";

const navItems = [
  { href: "/stream", label: "Streams" },
  { href: "/creator", label: "Creator" },
  { href: "/wallet", label: "Wallet" },
  { href: "/messages", label: "Messages" },
  { href: "/community", label: "Communities" },
  { href: "/events", label: "Events" },
  { href: "/search", label: "Search" },
  { href: "/notifications", label: "Notifications" },
  { href: "/settings", label: "Settings" },
];

export default function MainLayout({ children }: { children: ReactNode }) {
  return (
    <div style={{ display: "flex", minHeight: "100vh" }}>
      <aside
        style={{
          width: "240px",
          borderRight: "1px solid #e5e5e5",
          padding: "1rem",
          display: "flex",
          flexDirection: "column",
          gap: "0.25rem",
        }}
      >
        <Link
          href="/"
          style={{
            fontSize: "1.25rem",
            fontWeight: 700,
            marginBottom: "1.5rem",
            display: "block",
          }}
        >
          Spark
        </Link>
        {navItems.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            style={{
              padding: "0.5rem 0.75rem",
              borderRadius: "0.375rem",
              color: "#555",
              fontSize: "0.9rem",
            }}
          >
            {item.label}
          </Link>
        ))}
      </aside>
      <main style={{ flex: 1, padding: "1.5rem" }}>{children}</main>
    </div>
  );
}
