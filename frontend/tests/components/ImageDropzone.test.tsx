import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { ImageDropzone } from "@/components/products/ImageDropzone";

describe("ImageDropzone", () => {
  it("calls onChange with the dropped file", async () => {
    const onChange = vi.fn();
    render(<ImageDropzone value={null} onChange={onChange} />);

    const file = new File(["hello"], "cat.png", { type: "image/png" });
    const input = screen.getByTestId("dropzone-input") as HTMLInputElement;

    Object.defineProperty(input, "files", {
      value: [file],
      configurable: true,
    });
    fireEvent.change(input);

    await waitFor(() => {
      expect(onChange).toHaveBeenCalledWith(file);
    });
  });

  it("rejects unsupported file types", async () => {
    const onChange = vi.fn();
    render(<ImageDropzone value={null} onChange={onChange} />);
    const file = new File(["MZ"], "evil.exe", {
      type: "application/octet-stream",
    });
    const input = screen.getByTestId("dropzone-input") as HTMLInputElement;
    Object.defineProperty(input, "files", { value: [file], configurable: true });
    fireEvent.change(input);

    await waitFor(() => {
      expect(onChange).toHaveBeenCalledWith(null);
    });
  });
});
