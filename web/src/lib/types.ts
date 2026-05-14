export type User = {
  id: string;
  name: string;
  email: string;
  role: string;
  addresses?: Address[];
  segments?: string[];
  createdAt: string;
};

export type Address = {
  id?: string;
  fullName: string;
  phone: string;
  country: string;
  state: string;
  city: string;
  postalCode: string;
  addressLine: string;
  isDefault?: boolean;
};

export type ProductVariant = {
  id: string;
  sku: string;
  size: string;
  color: string;
  price: number;
  stock: number;
  weightKg?: number;
};

export type Product = {
  id: string;
  name: string;
  slug: string;
  description: string;
  image: string;
  categoryId: string;
  tags: string[];
  tagIds: string[];
  variants: ProductVariant[];
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
};

export type Category = {
  id: string;
  name: string;
  slug: string;
  description: string;
  isActive: boolean;
  createdAt: string;
};

export type Tag = {
  id: string;
  name: string;
  slug: string;
};

export type CartItem = {
  id: string;
  productId: string;
  variantId: string;
  name: string;
  sku: string;
  price: number;
  quantity: number;
  image: string;
};

export type Cart = {
  id: string;
  userId: string;
  items: CartItem[];
  createdAt: string;
  updatedAt: string;
};

export type OrderItem = {
  productId: string;
  variantId: string;
  name: string;
  sku: string;
  price: number;
  quantity: number;
};

export type Order = {
  id: string;
  userId: string;
  items: OrderItem[];
  shippingAddress: Address;
  subtotal: number;
  discount: number;
  shipping: number;
  shippingCarrier?: string;
  total: number;
  status: string;
  paymentStatus: string;
  invoiceNumber?: string;
  createdAt: string;
  updatedAt: string;
  paidAt?: string;
  fulfilledAt?: string;
  shippedAt?: string;
  cancelledAt?: string;
};

export type RMA = {
  id: string;
  userId: string;
  orderId: string;
  items: {
    productId: string;
    variantId: string;
    sku: string;
    name: string;
    quantity: number;
    price: number;
  }[];
  reason: string;
  status: string;
  adminNote?: string;
  refundAmount?: number;
  createdAt: string;
  updatedAt: string;
};

export type Banner = {
  id: string;
  slot: string;
  title: string;
  body: string;
  imageUrl: string;
  linkUrl: string;
  sortOrder: number;
  isActive: boolean;
  startsAt?: string;
  endsAt?: string;
  createdAt: string;
  updatedAt: string;
};

export type Discount = {
  id: string;
  code: string;
  type: string;
  value: number;
  minimumAmount: number;
  isActive: boolean;
  expiresAt: string;
  buyQty: number;
  getQty: number;
  productId?: string;
  categoryId?: string;
};

export type FeaturedHome = {
  newArrivals: Product[];
  bestSellers: Product[];
  featuredCategories: (Category & { productCount: number })[];
};

export type StorefrontStats = {
  productCount: number;
  categoryCount: number;
  freeShippingAtUsd: number;
};

/** Single response from GET /discovery/storefront */
export type StorefrontFeed = {
  featured: FeaturedHome;
  salePicks: Product[];
  announcementBanners: Banner[];
  carouselBanners: Banner[];
  stats: StorefrontStats;
  saleTagId?: string;
};

export type DashboardStats = {
  windowDays: number;
  totalRevenue: number;
  ordersPlaced: number;
  paidOrders: number;
  newCustomers: number;
  previous: {
    totalRevenue: number;
    ordersPlaced: number;
    paidOrders: number;
    newCustomers: number;
  };
};

export type TimelineEvent = { at: string; label: string; kind: string };

export type WishlistItem = {
  productId: string;
  variantId: string;
  sku: string;
  name: string;
  price: number;
  image: string;
  createdAt: string;
};

export type Wishlist = {
  userId: string;
  items: WishlistItem[];
  updatedAt: string;
};

export type SaveLaterItem = WishlistItem;

export type SaveLaterList = {
  userId: string;
  items: SaveLaterItem[];
  updatedAt: string;
};
