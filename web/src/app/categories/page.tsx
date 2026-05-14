"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { publicApi } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function CategoriesPage() {
  const q = useQuery({
    queryKey: ["categories"],
    queryFn: () => publicApi.categories(),
  });

  return (
    <div className="mx-auto max-w-6xl space-y-6 px-4 py-8">
      <h1 className="text-2xl font-semibold">Categories</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {(q.data ?? []).map((c) => (
          <Card key={c.id} className="border">
            <CardHeader>
              <CardTitle className="text-base">
                <Link href={`/categories/${c.id}`} className="hover:underline">
                  {c.name}
                </Link>
              </CardTitle>
            </CardHeader>
            <CardContent className="text-sm text-muted-foreground">
              {c.description}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
