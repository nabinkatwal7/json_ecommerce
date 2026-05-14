import type {
  Address,
  Banner,
  Cart,
  Category,
  DashboardStats,
  Discount,
  FeaturedHome,
  StorefrontFeed,
  Order,
  Product,
  RMA,
  SaveLaterItem,
  Tag,
  TimelineEvent,
  User,
  WishlistItem,
} from "@/lib/types";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "/backend";

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("token");
}

export function setToken(token: string | null) {
  if (typeof window === "undefined") return;
  if (token) localStorage.setItem("token", token);
  else localStorage.removeItem("token");
}

type FetchOpts = RequestInit & {
  json?: unknown;
  auth?: boolean;
  /** Extra headers (e.g. X-Cron-Secret) */
  extraHeaders?: Record<string, string>;
};

export async function apiFetch<T = unknown>(
  path: string,
  opts: FetchOpts = {},
): Promise<T> {
  const { json, auth = true, extraHeaders, ...init } = opts;
  const headers = new Headers(init.headers);
  if (json !== undefined) {
    headers.set("Content-Type", "application/json");
    init.body = JSON.stringify(json);
  }
  if (auth) {
    const t = getToken();
    if (t) headers.set("Authorization", `Bearer ${t}`);
  }
  if (extraHeaders) {
    for (const [k, v] of Object.entries(extraHeaders)) {
      if (v) headers.set(k, v);
    }
  }
  const res = await fetch(`${API_BASE}${path}`, { ...init, headers });
  if (!res.ok) {
    let msg = res.statusText;
    try {
      const body = (await res.json()) as { error?: string };
      if (body?.error) msg = body.error;
    } catch {
      /* ignore */
    }
    throw new ApiError(res.status, msg);
  }
  if (res.status === 204) return undefined as T;
  const ct = res.headers.get("content-type") ?? "";
  if (ct.includes("application/pdf")) {
    return (await res.blob()) as T;
  }
  if (ct.includes("application/json")) {
    return (await res.json()) as T;
  }
  return (await res.text()) as T;
}

/* ---------- Public ---------- */
export const publicApi = {
  health: () => apiFetch<{ ok: boolean }>("/health", { auth: false }),
  register: (body: { name: string; email: string; password: string }) =>
    apiFetch<{ token: string; user: User }>("/register", {
      auth: false,
      method: "POST",
      json: body,
    }),
  login: (body: { email: string; password: string }) =>
    apiFetch<{ token: string; user: User }>("/login", {
      auth: false,
      method: "POST",
      json: body,
    }),
  forgotPassword: (body: { email: string }) =>
    apiFetch<{ ok?: boolean }>("/forgot-password", {
      auth: false,
      method: "POST",
      json: body,
    }),
  resetPassword: (body: { token: string; newPassword: string }) =>
    apiFetch<unknown>("/reset-password", {
      auth: false,
      method: "POST",
      json: body,
    }),
  featured: () =>
    apiFetch<FeaturedHome>("/collections/featured", { auth: false }),
  storefrontFeed: () =>
    apiFetch<StorefrontFeed>("/discovery/storefront", { auth: false }),
  banners: (slot?: string) =>
    apiFetch<Banner[]>(
      slot ? `/banners?slot=${encodeURIComponent(slot)}` : "/banners",
      { auth: false },
    ),
  products: (categoryId?: string, tagId?: string) => {
    const p = new URLSearchParams();
    if (categoryId) p.set("categoryId", categoryId);
    if (tagId) p.set("tagId", tagId);
    const qs = p.toString();
    return apiFetch<Product[]>(qs ? `/products?${qs}` : "/products", { auth: false });
  },
  product: (id: string) => apiFetch<Product>(`/products/${id}`, { auth: false }),
  productBySlug: (slug: string) =>
    apiFetch<Product>(`/products/slug/${encodeURIComponent(slug)}`, {
      auth: false,
    }),
  related: (id: string, limit?: number) =>
    apiFetch<Product[]>(
      `/products/${encodeURIComponent(id)}/related${limit ? `?limit=${limit}` : ""}`,
      { auth: false },
    ),
  categories: () => apiFetch<Category[]>("/categories", { auth: false }),
  category: (id: string) =>
    apiFetch<Category>(`/categories/${encodeURIComponent(id)}`, { auth: false }),
  tags: () => apiFetch<Tag[]>("/tags", { auth: false }),
  search: (q: string, categoryId?: string, limit?: number) => {
    const p = new URLSearchParams({ q });
    if (categoryId) p.set("categoryId", categoryId);
    if (limit) p.set("limit", String(limit));
    return apiFetch<{ hits: Product[] }>(`/search?${p}`, { auth: false });
  },
  searchSuggest: (q: string, limit?: number) => {
    const p = new URLSearchParams({ q });
    if (limit) p.set("limit", String(limit));
    return apiFetch<{ suggestions: string[] }>(`/search/suggest?${p}`, {
      auth: false,
    });
  },
  abandonedCron: (cronSecret: string) =>
    apiFetch<{ sent: number }>("/internal/cron/abandoned-carts", {
      auth: false,
      method: "POST",
      extraHeaders: { "X-Cron-Secret": cronSecret },
    }),
};

/* ---------- Customer ---------- */
export const customerApi = {
  me: () => apiFetch<User>("/me"),
  patchMe: (body: { name?: string; email?: string }) =>
    apiFetch<User>("/me", { method: "PATCH", json: body }),
  insights: () =>
    apiFetch<{
      user: User;
      paidOrders: number;
      lifetimeSpend: number;
      bigSpenderThresholdUsd: number;
    }>("/me/insights"),
  addresses: () => apiFetch<Address[]>("/me/addresses"),
  postAddress: (body: Address) => apiFetch<Address[]>("/me/addresses", { method: "POST", json: body }),
  putAddress: (id: string, body: Address) =>
    apiFetch<Address[]>(`/me/addresses/${encodeURIComponent(id)}`, {
      method: "PUT",
      json: body,
    }),
  deleteAddress: (id: string) =>
    apiFetch<Address[]>(`/me/addresses/${encodeURIComponent(id)}`, {
      method: "DELETE",
    }),
  shippingQuote: (addr: Address) =>
    apiFetch<{ rates: { code: string; label: string; amount: number }[] }>(
      "/shipping/quote",
      { method: "POST", json: addr },
    ),
  couponValidate: (code: string) =>
    apiFetch<{
      valid: boolean;
      code?: string;
      message: string;
      discountType?: string;
      discountAmount: number;
      subtotal: number;
    }>("/coupons/validate", { method: "POST", json: { code } }),
  cart: () => apiFetch<Cart>("/cart"),
  cartValidate: () =>
    apiFetch<{
      ok: boolean;
      lines: {
        itemId: string;
        productId: string;
        variantId: string;
        name: string;
        sku: string;
        cartPrice: number;
        currentPrice: number;
        quantity: number;
        availableStock: number;
        issues: string[];
        ok: boolean;
      }[];
    }>("/cart/validate"),
  addCartItem: (body: { productId: string; variantId: string; quantity: number }) =>
    apiFetch<Cart>("/cart/items", { method: "POST", json: body }),
  patchCartItem: (itemId: string, quantity: number) =>
    apiFetch<Cart>(`/cart/items/${encodeURIComponent(itemId)}`, {
      method: "PATCH",
      json: { quantity },
    }),
  deleteCartItem: (itemId: string) =>
    apiFetch<Cart>(`/cart/items/${encodeURIComponent(itemId)}`, {
      method: "DELETE",
    }),
  wishlist: () => apiFetch<WishlistItem[]>("/wishlist"),
  postWishlist: (body: { productId: string; variantId: string }) =>
    apiFetch<WishlistItem[]>("/wishlist/items", { method: "POST", json: body }),
  deleteWishlist: (productId: string, variantId: string) =>
    apiFetch<WishlistItem[]>(
      `/wishlist/items?productId=${encodeURIComponent(productId)}&variantId=${encodeURIComponent(variantId)}`,
      { method: "DELETE" },
    ),
  wishlistMoveSave: (body: { productId: string; variantId: string }) =>
    apiFetch<{ wishlist: WishlistItem[]; saveLater: SaveLaterItem[] }>(
      "/wishlist/move-to-save-later",
      { method: "POST", json: body },
    ),
  saveLater: () => apiFetch<SaveLaterItem[]>("/save-later"),
  postSaveLater: (body: { productId: string; variantId: string }) =>
    apiFetch<SaveLaterItem[]>("/save-later/items", { method: "POST", json: body }),
  deleteSaveLater: (productId: string, variantId: string) =>
    apiFetch<SaveLaterItem[]>(
      `/save-later/items?productId=${encodeURIComponent(productId)}&variantId=${encodeURIComponent(variantId)}`,
      { method: "DELETE" },
    ),
  saveLaterMoveWish: (body: { productId: string; variantId: string }) =>
    apiFetch<{ wishlist: WishlistItem[]; saveLater: SaveLaterItem[] }>(
      "/save-later/move-to-wishlist",
      { method: "POST", json: body },
    ),
  checkout: (body: {
    shippingAddress: Address;
    discountCode?: string;
    shippingCarrier?: string;
  }) => apiFetch<Order>("/orders/checkout", { method: "POST", json: body }),
  orders: () => apiFetch<Order[]>("/orders"),
  order: (id: string) => apiFetch<Order>(`/orders/${encodeURIComponent(id)}`),
  orderInvoicePdf: async (id: string) => {
    const blob = await apiFetch<Blob>(`/orders/${encodeURIComponent(id)}/invoice.pdf`);
    return blob;
  },
  cancelOrder: (id: string) =>
    apiFetch<Order>(`/orders/${encodeURIComponent(id)}/cancel`, { method: "POST" }),
  stripeIntent: (id: string) =>
    apiFetch<{ clientSecret: string; paymentIntentId: string }>(
      `/orders/${encodeURIComponent(id)}/stripe-payment-intent`,
      { method: "POST" },
    ),
  payOrder: (id: string, body: { stripePaymentIntentId?: string; stub?: boolean }) =>
    apiFetch<{ order: Order; payment: unknown }>(
      `/orders/${encodeURIComponent(id)}/pay`,
      { method: "POST", json: body },
    ),
  rmas: () => apiFetch<RMA[]>("/rmas"),
  rma: (id: string) => apiFetch<RMA>(`/rmas/${encodeURIComponent(id)}`),
  postRma: (body: {
    orderId: string;
    reason: string;
    items: { productId: string; variantId: string; quantity: number }[];
  }) => apiFetch<RMA>("/rmas", { method: "POST", json: body }),
};

/* ---------- Admin ---------- */
export const adminApi = {
  dashboardStats: (days?: number) =>
    apiFetch<DashboardStats>(
      days ? `/admin/dashboard/stats?days=${days}` : "/admin/dashboard/stats",
    ),
  products: () => apiFetch<Product[]>("/admin/products"),
  postProduct: (body: Record<string, unknown>) =>
    apiFetch<Product>("/admin/products", { method: "POST", json: body }),
  putProduct: (id: string, body: Record<string, unknown>) =>
    apiFetch<Product>(`/admin/products/${encodeURIComponent(id)}`, {
      method: "PUT",
      json: body,
    }),
  deleteProduct: (id: string) =>
    apiFetch<undefined>(`/admin/products/${encodeURIComponent(id)}`, {
      method: "DELETE",
    }),
  postCategory: (body: Record<string, unknown>) =>
    apiFetch<Category>("/admin/categories", { method: "POST", json: body }),
  putCategory: (id: string, body: Record<string, unknown>) =>
    apiFetch<Category>(`/admin/categories/${encodeURIComponent(id)}`, {
      method: "PUT",
      json: body,
    }),
  deleteCategory: (id: string) =>
    apiFetch<undefined>(`/admin/categories/${encodeURIComponent(id)}`, {
      method: "DELETE",
    }),
  postDiscount: (body: Record<string, unknown>) =>
    apiFetch<Discount>("/admin/discounts", { method: "POST", json: body }),
  banners: () => apiFetch<Banner[]>("/admin/banners"),
  postBanner: (body: Record<string, unknown>) =>
    apiFetch<Banner>("/admin/banners", { method: "POST", json: body }),
  putBanner: (id: string, body: Record<string, unknown>) =>
    apiFetch<Banner>(`/admin/banners/${encodeURIComponent(id)}`, {
      method: "PUT",
      json: body,
    }),
  deleteBanner: (id: string) =>
    apiFetch<undefined>(`/admin/banners/${encodeURIComponent(id)}`, {
      method: "DELETE",
    }),
  adminTags: () => apiFetch<Tag[]>("/admin/tags"),
  postTag: (body: { name: string; slug: string }) =>
    apiFetch<Tag>("/admin/tags", { method: "POST", json: body }),
  putTag: (id: string, body: { name: string; slug: string }) =>
    apiFetch<Tag>(`/admin/tags/${encodeURIComponent(id)}`, {
      method: "PUT",
      json: body,
    }),
  deleteTag: (id: string) =>
    apiFetch<undefined>(`/admin/tags/${encodeURIComponent(id)}`, {
      method: "DELETE",
    }),
  orders: () => apiFetch<Order[]>("/admin/orders"),
  orderTimeline: (id: string) =>
    apiFetch<{ events: TimelineEvent[] }>(
      `/admin/orders/${encodeURIComponent(id)}/timeline`,
    ),
  cancelOrder: (id: string) =>
    apiFetch<Order>(`/admin/orders/${encodeURIComponent(id)}/cancel`, {
      method: "POST",
    }),
  fulfillOrder: (id: string) =>
    apiFetch<Order>(`/admin/orders/${encodeURIComponent(id)}/fulfill`, {
      method: "POST",
    }),
  shipOrder: (id: string) =>
    apiFetch<Order>(`/admin/orders/${encodeURIComponent(id)}/ship`, {
      method: "POST",
    }),
  lowStock: (threshold?: number) =>
    apiFetch<unknown>(
      threshold
        ? `/admin/inventory/low-stock?threshold=${threshold}`
        : "/admin/inventory/low-stock",
    ),
  rmas: () => apiFetch<RMA[]>("/admin/rmas"),
  rma: (id: string) => apiFetch<RMA>(`/admin/rmas/${encodeURIComponent(id)}`),
  rmaApprove: (id: string, note?: string) =>
    apiFetch<RMA>(`/admin/rmas/${encodeURIComponent(id)}/approve`, {
      method: "POST",
      json: note ? { note } : {},
    }),
  rmaReject: (id: string, note?: string) =>
    apiFetch<RMA>(`/admin/rmas/${encodeURIComponent(id)}/reject`, {
      method: "POST",
      json: note ? { note } : {},
    }),
  rmaReceive: (id: string, note?: string) =>
    apiFetch<RMA>(`/admin/rmas/${encodeURIComponent(id)}/receive`, {
      method: "POST",
      json: note ? { note } : {},
    }),
  rmaRefund: (id: string, note?: string) =>
    apiFetch<RMA>(`/admin/rmas/${encodeURIComponent(id)}/refund`, {
      method: "POST",
      json: note ? { note } : {},
    }),
  searchReindex: () =>
    apiFetch<unknown>("/admin/search/reindex", { method: "POST" }),
};
