export class ApiError extends Error {
  status: number;
  body: string;
  constructor(message: string, status: number, body: string) {
    super(message);
    this.status = status;
    this.body = body;
    this.name = "ApiError";
  }
}

// Sanitize error messages to hide internal API URLs
export function sanitizeErrorMessage(message: string | null | undefined): string {
  if (!message || typeof message !== "string") {
    return "API request failed";
  }
  // Hide localhost URLs and internal API paths
  return message
    .replace(/https?:\/\/[^\s/]+\/api/gi, "[API]")
    .replace(/\/api\/\w+/gi, "[endpoint]");
}
