import { resolveUrl } from "~/api/orval-client"
import { getStoredToken } from "~/utils/auth"

interface UploadResult {
  file: {
    id: string
    filename: string
    size: number
    content_type: string
    created_at: number
  }
}

interface PartInfo {
  part_number: number
  etag: string
}

function getAuthHeaders(): Record<string, string> {
  const token = getStoredToken()
  if (token) {
    return { Authorization: `Bearer ${token}` }
  }
  return {}
}

async function safeJsonError(response: Response): Promise<string> {
  try {
    const data = await response.json()
    return data.error || data.message || "Request failed"
  } catch {
    return `Request failed with status ${response.status}`
  }
}

const MAX_FILE_SIZE = 2 * 1024 * 1024 * 1024

export async function uploadFileInChunks(
  file: File,
  onProgress?: (percentage: number) => void
): Promise<UploadResult> {
  if (file.size > MAX_FILE_SIZE) {
    throw new Error("File size exceeds maximum allowed size of 2GB")
  }

  const initiateResponse = await fetch(resolveUrl("/api/files/multipart/initiate"), {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify({
      filename: file.name,
      content_type: file.type || "application/octet-stream",
      total_size: file.size,
    }),
  })

  if (!initiateResponse.ok) {
    const message = await safeJsonError(initiateResponse)
    throw new Error(message)
  }

  const initiateData = await initiateResponse.json()
  const { upload_id, chunk_size, total_parts } = initiateData

  const parts: PartInfo[] = []
  let uploadedParts = 0

  try {
    for (let partNumber = 1; partNumber <= total_parts; partNumber++) {
      const start = (partNumber - 1) * chunk_size
      const end = Math.min(start + chunk_size, file.size)
      const chunk = file.slice(start, end)

      const formData = new FormData()
      formData.append("chunk", chunk)

      const response = await fetch(
        resolveUrl(`/api/files/multipart/${upload_id}/part/${partNumber}`),
        {
          method: "POST",
          headers: getAuthHeaders(),
          body: formData,
        }
      )

      if (!response.ok) {
        const message = await safeJsonError(response)
        throw new Error(message)
      }

      const partResponse = await response.json()
      parts.push({
        part_number: partNumber,
        etag: partResponse.etag,
      })

      uploadedParts++
      if (onProgress) {
        onProgress(Math.round((uploadedParts / total_parts) * 100))
      }
    }

    const completeResponse = await fetch(
      resolveUrl(`/api/files/multipart/${upload_id}/complete`),
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...getAuthHeaders(),
        },
        body: JSON.stringify({ parts }),
      }
    )

    if (!completeResponse.ok) {
      const message = await safeJsonError(completeResponse)
      throw new Error(message)
    }

    const completeData = await completeResponse.json()
    return completeData
  } catch (error) {
    try {
      await fetch(resolveUrl(`/api/files/multipart/${upload_id}`), {
        method: "DELETE",
        headers: getAuthHeaders(),
      })
    } catch (abortError) {
      console.error("Failed to abort upload:", abortError)
    }
    throw error
  }
}
