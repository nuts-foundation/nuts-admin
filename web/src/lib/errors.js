/**
 * Parses API error responses into a formatted error message.
 * Handles both string and object error responses.
 * Supports RFC 7807 Problem Details format (title, status, detail).
 *
 * @param {string|object|Error} error - The error from the API call
 * @returns {string} Formatted error message in the format "status - title - detail" (detail is optional)
 * @example
 * // Example error response:
 * // { "title": "RevokeVC failed", "status": 400, "detail": "credential contains no (relevant) status" }
 * // Returns: "400 - RevokeVC failed - credential contains no (relevant) status"
 */
export function parseApiError(error) {
  // Parse string errors as JSON
  let errorObj = error
  if (typeof error === 'string') {
    try {
      errorObj = JSON.parse(error)
    } catch {
      return error
    }
  }

  if (!errorObj) return 'An error occurred'

  // Format: status - title - detail
  const parts = [errorObj.status, errorObj.title, errorObj.detail].filter(Boolean)
  return parts.length > 0 ? parts.join(' - ') : (errorObj.message || errorObj.toString() || 'An error occurred')
}


