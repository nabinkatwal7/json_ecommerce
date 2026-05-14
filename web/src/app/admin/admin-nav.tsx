"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { SITE_NAME } from "@/lib/site";

const links: { href: string; label: string }[] = [
  { href: "/admin", label: "Dashboard" },
  { href: "/admin/products", label: "Products" },
  { href: "/admin/categories", label: "Categories" },
  { href: "/admin/discounts", label: "Discounts" },
  { href: "/admin/banners", label: "Banners" },
  { href: "/admin/tags", label: "Tags" },
  { href: "/admin/orders", label: "Orders" },
  { href: "/admin/inventory/low-stock", label: "Low stock" },
  { href: "/admin/rmas", label: "RMAs" },
  { href: "/admin/search", label: "Search" },
];

export function AdminNav() {
  const pathname = usePathname();
  return (
    <aside className="w-56 shrink-0 border-r border-border/80 bg-muted/15 p-5">
      <Link href="/" className="mb-6 block">
        <span className="font-display text-lg font-semibold text-foreground">{SITE_NAME}</span>
        <span className="mt-0.5 block text-xs font-medium text-muted-foreground">← Back to shop</span>
      </Link>
      <nav className="flex flex-col gap-0.5 text-sm">
        {links.map((l) => (
          <Link
            key={l.href}
            href={l.href}
            className={cn(
              "rounded-lg px-3 py-2 font-medium text-muted-foreground transition-colors hover:bg-background hover:text-foreground",
              pathname === l.href && "bg-background text-foreground shadow-sm",
            )}
          >
            {l.label}
          </Link>
        ))}
      </nav>
    </aside>
  );
}
