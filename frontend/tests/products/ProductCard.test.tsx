import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { ProductCard } from "@/components/products/ProductCard";
import { sampleProducts } from "../msw/handlers";

describe("ProductCard", () => {
  it("renders name, price, and type label", () => {
    const product = sampleProducts[0];
    render(<ProductCard product={product} />);
    expect(screen.getByText(product.product_name)).toBeInTheDocument();
    expect(screen.getByText("Drinks")).toBeInTheDocument();
    // Locale-independent: just confirm "29" and "99" are both shown. The
    // currency separator depends on the runtime's default Intl locale.
    expect(screen.getByText(/29[.,]99/)).toBeInTheDocument();
    expect(
      screen.getByText(`${product.product_quantity} in stock`)
    ).toBeInTheDocument();
  });

  it("links to the product detail page", () => {
    const product = sampleProducts[0];
    render(<ProductCard product={product} />);
    const link = screen.getByRole("link");
    expect(link).toHaveAttribute(
      "href",
      `/products/${encodeURIComponent(product.product_id)}`
    );
  });
});
