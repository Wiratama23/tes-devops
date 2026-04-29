import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { ArticleCard } from "@/components/articles/ArticleCard";
import { sampleArticles } from "../msw/handlers";

describe("ArticleCard", () => {
  it("renders title and stripped preview", () => {
    const article = sampleArticles[0];
    render(<ArticleCard article={article} />);
    expect(screen.getByText(article.title)).toBeInTheDocument();
    expect(screen.getByText(/Welcome to the blog/)).toBeInTheDocument();
  });

  it("links to the detail page", () => {
    const article = sampleArticles[0];
    render(<ArticleCard article={article} />);
    expect(screen.getByRole("link")).toHaveAttribute(
      "href",
      `/articles/${article.articles_id}`
    );
  });
});
