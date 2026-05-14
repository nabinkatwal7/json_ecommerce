"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { SITE_NAME, SITE_TAGLINE } from "@/lib/site";
import { cn } from "@/lib/utils";
import { CartNavLink } from "@/components/cart-nav-link";

const navLink =
  "text-sm font-medium text-muted-foreground transition-colors hover:text-foreground";

export function SiteHeader() {
  const pathname = usePathname();
  const { user, logout, loading } = useAuth();

  if (pathname?.startsWith("/admin")) return null;

  return (
    <header className="sticky top-0 z-50 border-b border-border/70 bg-background/90 backdrop-blur-md">
      <div className="mx-auto flex h-[4.25rem] max-w-6xl items-center justify-between gap-6 px-4 sm:px-6">
        <Link href="/" className="group flex shrink-0 flex-col gap-0 leading-none">
          <span className="font-display text-2xl font-semibold tracking-tight text-foreground transition-colors group-hover:text-primary sm:text-[1.65rem]">
            {SITE_NAME}
          </span>
          <span className="hidden text-[11px] font-medium tracking-wide text-muted-foreground sm:block">
            {SITE_TAGLINE}
          </span>
        </Link>

        <nav className="hidden flex-1 items-center justify-center gap-8 md:flex">
          <Link href="/products" className={navLink}>
            Kits
          </Link>
          <Link href="/categories" className={navLink}>
            Categories
          </Link>
          <Link href="/search" className={navLink}>
            Search
          </Link>
          {user && (
            <>
              <CartNavLink className={navLink} />
              <Link href="/orders" className={navLink}>
                Orders
              </Link>
            </>
          )}
        </nav>

        <div className="flex items-center gap-2">
          {user?.role === "admin" && (
            <Button variant="outline" size="sm" className="hidden border-primary/25 sm:inline-flex" asChild>
              <Link href="/admin">Admin</Link>
            </Button>
          )}
          {loading ? (
            <span className="h-8 w-8 animate-pulse rounded-full bg-muted" aria-hidden />
          ) : user ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="secondary" size="sm" className="font-medium">
                  {user.name}
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-52">
                <DropdownMenuItem asChild>
                  <Link href="/account">Account</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href="/account/addresses">Addresses</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href="/account/insights">Insights</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href="/wishlist">Wishlist</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href="/save-later">Save for later</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href="/rmas">Returns</Link>
                </DropdownMenuItem>
                {user.role === "admin" && (
                  <DropdownMenuItem asChild className="sm:hidden">
                    <Link href="/admin">Admin</Link>
                  </DropdownMenuItem>
                )}
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => logout()}>Log out</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <div className="flex items-center gap-1">
              <Button variant="ghost" size="sm" className="font-medium" asChild>
                <Link href="/login">Log in</Link>
              </Button>
              <Button size="sm" className="font-semibold shadow-sm" asChild>
                <Link href="/register">Join</Link>
              </Button>
            </div>
          )}
        </div>
      </div>

      <div className="border-t border-border/50 bg-muted/30 px-4 py-2.5 md:hidden">
        <div className="mx-auto flex max-w-6xl flex-wrap justify-center gap-x-5 gap-y-1 text-xs font-medium">
          <Link href="/products" className={cn(navLink, "text-xs")}>
            Kits
          </Link>
          <Link href="/categories" className={cn(navLink, "text-xs")}>
            Categories
          </Link>
          <Link href="/search" className={cn(navLink, "text-xs")}>
            Search
          </Link>
          {user && (
            <>
              <CartNavLink className={cn(navLink, "text-xs")} />
              <Link href="/orders" className={cn(navLink, "text-xs")}>
                Orders
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
