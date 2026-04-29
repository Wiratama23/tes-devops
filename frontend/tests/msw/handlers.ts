import { http, HttpResponse } from "msw";

const API = "http://api.test/api";

export const sampleProducts = [
  {
    product_id: "SKU10001",
    product_name: "Premium Coffee",
    product_quantity: 12,
    product_prices: "29.99",
    product_type: "10",
    created_at: "2026-04-19T12:00:00Z",
    created_by: "00000000-0000-0000-0000-000000000001",
    image_path: "assets/coffee.jpg",
  },
  {
    product_id: "SKU10002",
    product_name: "Cold Brew Concentrate",
    product_quantity: 5,
    product_prices: "14.50",
    product_type: "10",
    created_at: "2026-04-19T12:30:00Z",
    created_by: "00000000-0000-0000-0000-000000000001",
    image_path: "assets/cold-brew.jpg",
  },
];

export const sampleArticles = [
  {
    articles_id: 1,
    uid: "00000000-0000-0000-0000-000000000001",
    title: "Getting Started",
    article_text: "<p>Welcome to the blog.</p>",
    date_created: "2026-04-19T10:30:00Z",
    updated_at: "2026-04-19T10:30:00Z",
  },
  {
    articles_id: 2,
    uid: "00000000-0000-0000-0000-000000000001",
    title: "Advanced Patterns",
    article_text: "<p>Some advanced patterns…</p>",
    date_created: "2026-04-20T10:30:00Z",
    updated_at: "2026-04-20T10:30:00Z",
  },
];

export const handlers = [
  http.get(`${API}/products`, () =>
    HttpResponse.json({ data: sampleProducts, limit: 10, offset: 0 })
  ),
  http.get(`${API}/articles`, () =>
    HttpResponse.json({
      data: sampleArticles,
      total_count: sampleArticles.length,
      limit: 10,
      offset: 0,
    })
  ),
  http.get(`${API}/users`, () => HttpResponse.json([])),
  http.get(`${API}/auth/me`, () =>
    HttpResponse.json({
      uid: "00000000-0000-0000-0000-000000000001",
      username: "admin",
      email: "admin@example.com",
      is_admin: true,
    })
  ),
  http.post(`${API}/auth/login`, async ({ request }) => {
    const body = (await request.json()) as { username: string; password: string };
    if (body.username === "admin" && body.password === "secret123") {
      return HttpResponse.json({
        token: "test-token",
        user: {
          uid: "00000000-0000-0000-0000-000000000001",
          username: "admin",
          email: "admin@example.com",
          is_admin: true,
        },
        expires: Math.floor(Date.now() / 1000) + 3600,
      });
    }
    return new HttpResponse("invalid credentials", { status: 401 });
  }),
  http.post(`${API}/uploads/images`, () =>
    HttpResponse.json(
      {
        path: "assets/uploaded.png",
        url: `${API}/assets/uploaded.png`,
        filename: "uploaded.png",
        size: 100,
      },
      { status: 201 }
    )
  ),
  http.post(`${API}/logs`, () => new HttpResponse(null, { status: 204 })),
  http.delete(`${API}/articles/:id`, () => new HttpResponse(null, { status: 204 })),
  http.delete(`${API}/products/:id`, () => new HttpResponse(null, { status: 204 })),
];
