/**
 * Go API wraps JSON as `{ "data": <payload> }` (backend/pkg/httpresponse).
 * Orval assigns the full parsed body to `response.data`, so list endpoints often
 * need one more `.data` unwrap than the generated types suggest.
 */
export function normalizeListPayload<T>(parsedBody: unknown): T[] {
	if (parsedBody == null) {
		return []
	}
	if (Array.isArray(parsedBody)) {
		return parsedBody as T[]
	}
	if (
		typeof parsedBody === "object" &&
		parsedBody !== null &&
		"data" in parsedBody
	) {
		const inner = (parsedBody as { data: unknown }).data
		if (inner == null || !Array.isArray(inner)) {
			return []
		}
		return inner as T[]
	}
	return []
}
