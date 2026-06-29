import Link from "next/link";

export default function HomePage() {
  return (
    <div style={{ padding: "2rem", maxWidth: "1200px", margin: "0 auto" }}>
      <header
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          marginBottom: "3rem",
        }}
      >
        <h1 style={{ fontSize: "1.5rem", fontWeight: 700 }}>Spark</h1>
        <nav style={{ display: "flex", gap: "1rem" }}>
          <Link href="/login">Sign In</Link>
          <Link href="/register">Get Started</Link>
        </nav>
      </header>

      <section style={{ textAlign: "center", marginBottom: "4rem" }}>
        <h2 style={{ fontSize: "3rem", fontWeight: 800, marginBottom: "1rem" }}>
          The future of live streaming
        </h2>
        <p
          style={{
            fontSize: "1.25rem",
            color: "#666",
            marginBottom: "2rem",
          }}
        >
          Create, share, and discover live content with your community.
        </p>
        <Link
          href="/register"
          style={{
            background: "#0070c4",
            color: "#fff",
            padding: "0.75rem 2rem",
            borderRadius: "0.5rem",
            fontWeight: 600,
            display: "inline-block",
          }}
        >
          Start Streaming
        </Link>
      </section>

      <section
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(300px, 1fr))",
          gap: "2rem",
        }}
      >
        {[
          { title: "Live Streaming", desc: "Go live in seconds with ultra-low latency streaming." },
          { title: "Monetize", desc: "Earn through subscriptions, tips, gifts, and ads." },
          { title: "Community", desc: "Build real connections with chat, communities, and events." },
        ].map((feature) => (
          <div
            key={feature.title}
            style={{
              padding: "2rem",
              border: "1px solid #e5e5e5",
              borderRadius: "0.75rem",
            }}
          >
            <h3 style={{ fontSize: "1.25rem", fontWeight: 600, marginBottom: "0.5rem" }}>
              {feature.title}
            </h3>
            <p style={{ color: "#666" }}>{feature.desc}</p>
          </div>
        ))}
      </section>
    </div>
  );
}
