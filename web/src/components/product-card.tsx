import { fmtMoney } from "@/lib/format";
import type { Product } from "@/lib/types";
import { Card, CardContent } from "@/components/ui/card";
import Link from "next/link";
import { Badge } from "@/components/ui/badge";

export function ProductCard({ product }: { product: Product }) {
  const minPrice = Math.min(...product.variants.map((v) => v.price));
  return (
    <Card className="group overflow-hidden border border-border/80 bg-card shadow-sm transition-all duration-200 hover:border-primary/25 hover:shadow-md">
      <Link href={`/products/${product.slug}`} className="block outline-none ring-offset-background focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
        <div className="relative aspect-[4/5] w-full overflow-hidden bg-muted">
          {product.image ? (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              src={product.image}
              alt={product.name}
              className="h-full w-full object-cover transition-transform duration-300 ease-out group-hover:scale-[1.03]"
            />
          ) : (
            <div className="flex h-full items-center justify-center text-xs text-muted-foreground">
              No image
            </div>
          )}
        </div>
        <CardContent className="space-y-1.5 p-4">
          <p className="line-clamp-2 text-[15px] font-medium leading-snug tracking-tight text-foreground">
            {product.name}
          </p>
          <p className="text-sm tabular-nums tracking-tight text-muted-foreground">
            From <span className="font-medium text-foreground">{fmtMoney(minPrice)}</span>
          </p>
          {!product.isActive && (
            <Badge variant="secondary" className="text-[10px]">
              Inactive
            </Badge>
          )}
        </CardContent>
      </Link>
    </Card>
  );
}
