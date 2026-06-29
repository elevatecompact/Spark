export function Button({
  children,
  variant = "primary",
  ...props
}: React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost";
}) {
  const styles: Record<string, React.CSSProperties> = {
    primary: {
      background: "#0070c4",
      color: "#fff",
      border: "none",
    },
    secondary: {
      background: "transparent",
      color: "#0070c4",
      border: "1px solid #0070c4",
    },
    ghost: {
      background: "transparent",
      color: "#555",
      border: "none",
    },
  };

  return (
    <button
      style={{
        padding: "0.5rem 1rem",
        borderRadius: "0.375rem",
        fontWeight: 600,
        cursor: "pointer",
        ...styles[variant],
        ...(props.disabled ? { opacity: 0.5, cursor: "not-allowed" } : {}),
      }}
      {...props}
    >
      {children}
    </button>
  );
}
