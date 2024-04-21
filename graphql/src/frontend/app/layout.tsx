import type { Metadata } from "next";
import { Inter as FontSans } from "next/font/google";
import "./globals.css";
import { cn } from "@/lib/utils";
import Providers from "./providers";
import Link from "next/link";
import { Toaster } from "@/components/ui/sonner";

const fontSans = FontSans({
  subsets: ["latin"],
  variable: "--font-sans",
});

export const metadata: Metadata = {
  title: "My GraphQL Site",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={cn(
          "min-h-screen bg-background font-sans antialiased p-4",
          fontSans.variable
        )}
      >
        <nav className="flex gap-4 py-4">
          <Link href="/" className="hover:underline text-muted-foreground">
            Home
          </Link>

          <Link
            href="/products"
            className="hover:underline text-muted-foreground"
          >
            Products
          </Link>

          <Link
            href="/basket"
            className="hover:underline text-muted-foreground"
          >
            Basket
          </Link>

          <Link href="/admin" className="hover:underline text-muted-foreground">
            Admin
          </Link>
        </nav>

        <Providers>{children}</Providers>
        <Toaster />
      </body>
    </html>
  );
}
