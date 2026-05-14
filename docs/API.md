# JSON E‑Commerce HTTP API

Go + Gin JSON API. Default base URL: `http://localhost:8080` (override with env `ADDR`).

## Conventions

- **Content-Type**: `application/json` for bodies (except `GET /orders/:id/invoice.pdf`, which returns PDF).
- **Errors**: JSON object `{"error":"<message>"}` with appropriate HTTP status.
- **Auth (customer / admin JWT)**: `Authorization: Bearer <token>`
  Obtain a token from `POST /register` or `POST /login` (`token` field in response).
- **Admin**: Same JWT; user must have role `admin`. All routes below **`/admin/*`** require admin.
- **Rate limiting**: Global limiter on all routes; stricter limiter on register, login, forgot/reset password, and the abandoned-cart cron route.

---

## Public (no JWT)

| Method | Path                             | Description                                                                                                                      |
| ------ | -------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `GET`  | `/health`                        | Liveness check. Returns `{"ok":true}`.                                                                                           |
| `POST` | `/register`                      | Create account. Body: `name`, `email`, `password` (min 8 chars). Returns `token` + `user`.                                       |
| `POST` | `/login`                         | Body: `email`, `password`. Returns `token` + `user`.                                                                             |
| `POST` | `/forgot-password`               | Body: `email`. Sends reset email when SMTP is configured (anti-enumeration: always `200` with `{"ok":true}` on success path).    |
| `POST` | `/reset-password`                | Body: `token`, `newPassword`.                                                                                                    |
| `GET`  | `/products`                      | Active products. Query: optional `categoryId`.                                                                                   |
| `GET`  | `/products/:id`                  | Active product by ID.                                                                                                            |
| `GET`  | `/products/slug/:slug`           | Active product by slug.                                                                                                          |
| `GET`  | `/categories`                    | Active categories.                                                                                                               |
| `GET`  | `/categories/:id`                | Active category by ID.                                                                                                           |
| `GET`  | `/tags`                          | All tags.                                                                                                                        |
| `GET`  | `/search`                        | Search products. Query: `q`, optional `categoryId`, `limit` (default 20). Uses Algolia when configured, else local fuzzy search. |
| `POST` | `/internal/cron/abandoned-carts` | Abandoned-cart email job. Header **`X-Cron-Secret`** must match env `CRON_SECRET`. Stricter rate limit. No JWT.                  |

---

## Authenticated customer (`Authorization: Bearer …`)

| Method   | Path                                | Description                                                                                                                                                                    |
| -------- | ----------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `GET`    | `/me`                               | Current user profile (no password hash).                                                                                                                                       |
| `GET`    | `/me/insights`                      | Refreshes marketing `segments` on user; returns user, `paidOrders`, `lifetimeSpend`, `bigSpenderThresholdUsd`.                                                                 |
| `GET`    | `/me/addresses`                     | List saved addresses.                                                                                                                                                          |
| `POST`   | `/me/addresses`                     | Add address. Body: `Address` fields (`fullName`, `phone`, `country`, `state`, `city`, `postalCode`, `addressLine`, optional `isDefault`, `id`).                                |
| `PUT`    | `/me/addresses/:id`                 | Update address by `id`.                                                                                                                                                        |
| `DELETE` | `/me/addresses/:id`                 | Remove address.                                                                                                                                                                |
| `POST`   | `/shipping/quote`                   | Carrier rate quotes from current cart weight. Body: shipping `Address` (e.g. `country`; defaults to `US`). Returns `rates[]` (FedEx/UPS/DHL-style stub codes).                 |
| `GET`    | `/wishlist`                         | Wishlist line items.                                                                                                                                                           |
| `POST`   | `/wishlist/items`                   | Body: `productId`, `variantId`.                                                                                                                                                |
| `DELETE` | `/wishlist/items`                   | Query: `productId`, `variantId`.                                                                                                                                               |
| `POST`   | `/wishlist/move-to-save-later`      | Body: `productId`, `variantId`. Moves item to save-for-later.                                                                                                                  |
| `GET`    | `/save-later`                       | Save-for-later items.                                                                                                                                                          |
| `POST`   | `/save-later/items`                 | Body: `productId`, `variantId`.                                                                                                                                                |
| `DELETE` | `/save-later/items`                 | Query: `productId`, `variantId`.                                                                                                                                               |
| `POST`   | `/save-later/move-to-wishlist`      | Body: `productId`, `variantId`.                                                                                                                                                |
| `GET`    | `/cart`                             | Current user cart (created if missing).                                                                                                                                        |
| `POST`   | `/cart/items`                       | Add line. Body: `productId`, `variantId`, `quantity`.                                                                                                                          |
| `PATCH`  | `/cart/items/:itemId`               | Body: `quantity`.                                                                                                                                                              |
| `DELETE` | `/cart/items/:itemId`               | Remove line.                                                                                                                                                                   |
| `POST`   | `/rmas`                             | Create return request. Body: `orderId`, `reason`, `items[]` with `productId`, `variantId`, `quantity` (order must be paid + fulfilled/shipped).                                |
| `GET`    | `/rmas`                             | List caller’s RMAs.                                                                                                                                                            |
| `GET`    | `/rmas/:id`                         | RMA detail (owner only).                                                                                                                                                       |
| `POST`   | `/orders/checkout`                  | Create order, deduct stock, clear cart. Body: `shippingAddress`, optional `discountCode`, optional `shippingCarrier` (`flat` / `fedex_ground` / `ups_ground` / `dhl_express`). |
| `GET`    | `/orders`                           | List caller’s orders.                                                                                                                                                          |
| `GET`    | `/orders/:id/invoice.pdf`           | Download PDF invoice (assigns invoice number on first use).                                                                                                                    |
| `POST`   | `/orders/:id/cancel`                | Cancel unpaid `created` order; restores stock.                                                                                                                                 |
| `POST`   | `/orders/:id/stripe-payment-intent` | Create Stripe PaymentIntent; returns `clientSecret`, `paymentIntentId`.                                                                                                        |
| `POST`   | `/orders/:id/pay`                   | Confirm payment. Body: `stripePaymentIntentId` **or** (dev) `stub: true` when `DEV_PAYMENT_STUB` is enabled.                                                                   |
| `GET`    | `/orders/:id`                       | Order detail (owner only).                                                                                                                                                     |

---

## Admin (`Authorization: Bearer …` + role `admin`)

All paths are prefixed with **`/admin`**.

| Method   | Path                         | Description                                                                                                                                                                                                     |
| -------- | ---------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `GET`    | `/admin/products`            | All products (including inactive).                                                                                                                                                                              |
| `POST`   | `/admin/products`            | Create product. Body: PIM fields (`name`, `slug`, `description`, `image`, `categoryId`, `tags`, `tagIds`, `variants[]` with SKU, optional `weightKg`, etc.).                                                    |
| `PUT`    | `/admin/products/:id`        | Update product.                                                                                                                                                                                                 |
| `DELETE` | `/admin/products/:id`        | Delete product.                                                                                                                                                                                                 |
| `POST`   | `/admin/categories`          | Create category.                                                                                                                                                                                                |
| `PUT`    | `/admin/categories/:id`      | Update category.                                                                                                                                                                                                |
| `DELETE` | `/admin/categories/:id`      | Delete category.                                                                                                                                                                                                |
| `POST`   | `/admin/discounts`           | Create discount. Body: `code`, `type` (`percent` \| `fixed` \| `bogo`), `value`, `minimumAmount`, `isActive`, optional `expiresAt` (RFC3339), for BOGO: `buyQty`, `getQty`, optional `productId`, `categoryId`. |
| `GET`    | `/admin/tags`                | List tags.                                                                                                                                                                                                      |
| `POST`   | `/admin/tags`                | Create tag. Body: `name`, `slug`.                                                                                                                                                                               |
| `PUT`    | `/admin/tags/:id`            | Update tag.                                                                                                                                                                                                     |
| `DELETE` | `/admin/tags/:id`            | Delete tag.                                                                                                                                                                                                     |
| `GET`    | `/admin/orders`              | All orders.                                                                                                                                                                                                     |
| `POST`   | `/admin/orders/:id/cancel`   | Cancel order; restores stock.                                                                                                                                                                                   |
| `POST`   | `/admin/orders/:id/fulfill`  | Fulfill paid order.                                                                                                                                                                                             |
| `POST`   | `/admin/orders/:id/ship`     | Mark fulfilled order shipped.                                                                                                                                                                                   |
| `GET`    | `/admin/inventory/low-stock` | Low-stock report. Query: optional `threshold` (default from `LOW_STOCK_THRESHOLD`).                                                                                                                             |
| `GET`    | `/admin/rmas`                | All RMAs.                                                                                                                                                                                                       |
| `GET`    | `/admin/rmas/:id`            | RMA by ID.                                                                                                                                                                                                      |
| `POST`   | `/admin/rmas/:id/approve`    | Optional body: `note`.                                                                                                                                                                                          |
| `POST`   | `/admin/rmas/:id/reject`     | Optional body: `note`.                                                                                                                                                                                          |
| `POST`   | `/admin/rmas/:id/receive`    | Optional body: `note`.                                                                                                                                                                                          |
| `POST`   | `/admin/rmas/:id/refund`     | Refund + restore inventory. Optional body: `note`.                                                                                                                                                              |
| `POST`   | `/admin/search/reindex`      | Reindex active products to Algolia (requires Algolia env vars).                                                                                                                                                 |

---

## Run the server

```bash
go run ./cmd/server
```

See `internal/config/config.go` and environment variables used across the app (e.g. `DATA_DIR`, `JWT_SECRET`, `STRIPE_SECRET_KEY`, `SMTP_*`, `REDIS_ADDR`, `ALGOLIA_*`, `CRON_SECRET`, `ADMIN_EMAIL`, `TLS_CERT_FILE` / `TLS_KEY_FILE`).

---

## Route count summary

| Area                |  Count |
| ------------------- | -----: |
| Public              |     13 |
| Customer (JWT)      |     29 |
| Admin (JWT + admin) |     24 |
| **Total**           | **66** |

_(One row per HTTP method + path as registered in `internal/api/router.go`.)_
