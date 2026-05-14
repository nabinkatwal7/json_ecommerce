"use client";

import Link from "next/link";
import { useAuth } from "@/contexts/auth-context";
import { AdminNav } from "./admin-nav";

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();

  if (loading) {
    return <p className="p-8 text-sm">Loading…</p>;
  }
  if (!user) {
    return (
      <div className="p-8 text-center text-sm">
        <Link href="/login" className="underline">
          Log in
        </Link>{" "}
        as an admin user.
      </div>
    );
  }
  if (user.role !== "admin") {
    return <p className="p-8 text-center text-sm">Admin access only.</p>;
  }

  return (
    <div className="flex min-h-screen">
      <AdminNav />
      <div className="flex-1 overflow-auto p-6">{children}</div>
    </div>
  );
}
