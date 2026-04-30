import { cn } from "@/tools/utils";

// Renders the HTML produced by the Tiptap editor on the admin side. The
// content is currently stored as raw HTML in `article_text`. Sanitization is
// the backend's responsibility (the WAF + the admin-only write path); we
// intentionally avoid adding a heavyweight sanitizer to the bundle.
interface RichTextRendererProps {
  html: string;
  className?: string;
}

export function RichTextRenderer({ html, className }: RichTextRendererProps) {
  return (
    <div
      className={cn(
        "prose prose-sm prose-zinc max-w-none dark:prose-invert",
        "prose-headings:font-semibold prose-a:text-primary prose-img:rounded-lg",
        className
      )}
      dangerouslySetInnerHTML={{ __html: html }}
    />
  );
}
