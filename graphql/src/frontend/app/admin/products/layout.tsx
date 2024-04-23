import Link from "next/link";

export default function AdminLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div>
      <nav className="flex gap-4">
        <Link
          href="/admin/products/create"
          className="hover:underline text-muted-foreground"
        >
          Create
        </Link>
      </nav>

      {children}
    </div>
  );
}
