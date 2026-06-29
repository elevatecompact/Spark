import Link from "next/link";

export default function ForgotPasswordPage() {
  return (
    <div
      style={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <div
        style={{
          width: "100%",
          maxWidth: "400px",
          padding: "2rem",
          display: "flex",
          flexDirection: "column",
          gap: "1rem",
        }}
      >
        <h1 style={{ fontSize: "1.5rem", fontWeight: 700, marginBottom: "0.5rem" }}>
          Reset Password
        </h1>
        <p style={{ color: "#666" }}>
          Enter your email and we&apos;ll send you a reset link.
        </p>
        <input
          type="email"
          placeholder="Email"
          required
          style={{
            padding: "0.75rem",
            border: "1px solid #e5e5e5",
            borderRadius: "0.5rem",
          }}
        />
        <button
          style={{
            padding: "0.75rem",
            background: "#0070c4",
            color: "#fff",
            border: "none",
            borderRadius: "0.5rem",
            fontWeight: 600,
            cursor: "pointer",
          }}
        >
          Send Reset Link
        </button>
        <p style={{ fontSize: "0.875rem", textAlign: "center" }}>
          <Link href="/login">Back to sign in</Link>
        </p>
      </div>
    </div>
  );
}
