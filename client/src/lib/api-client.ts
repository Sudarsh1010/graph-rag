const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080"

export class ApiError extends Error {
  constructor(
    public status: number,
    public statusText: string,
    public body?: unknown
  ) {
    super(`API Error: ${status} ${statusText}`)
    this.name = "ApiError"
  }
}

export async function apiFetch<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const url = `${API_BASE_URL}${path}`

  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  })

  if (!response.ok) {
    let body: unknown
    try {
      body = await response.json()
    } catch {
      body = await response.text()
    }
    throw new ApiError(response.status, response.statusText, body)
  }

  // Handle empty responses
  const contentType = response.headers.get("content-type")
  if (!contentType?.includes("application/json")) {
    return undefined as T
  }

  return response.json()
}
