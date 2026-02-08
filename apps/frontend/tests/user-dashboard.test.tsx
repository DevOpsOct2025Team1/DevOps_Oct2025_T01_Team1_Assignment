import { describe, it, expect, vi, beforeEach } from "vitest"
import { BrowserRouter } from "react-router"
import { renderWithProviders } from "./test-utils"
import UserPanel from "../app/components/UserPanel"

const getFilesMock = vi.fn()
const uploadFileMock = vi.fn()
const deleteFileMock = vi.fn()

vi.mock("../app/api/generated", () => ({
  useGetApiFiles: (options: any) => {
    const data = getFilesMock()
    if (options?.query?.select && data?.data) {
      return {
        data: options.query.select(data),
        isLoading: false,
        refetch: vi.fn(),
      }
    }
    return {
      data,
      isLoading: false,
      refetch: vi.fn(),
    }
  },
  usePostApiFiles: () => ({
    mutate: uploadFileMock,
    isPending: false,
  }),
  useDeleteApiFilesId: () => ({
    mutate: deleteFileMock,
    isPending: false,
  }),
  getGetApiFilesIdDownloadUrl: (id: string) => `/api/files/${id}/download`,
}))

vi.mock("../app/api/orval-client", () => ({
  resolveUrl: (url: string) => `http://localhost:3001${url}`,
  getAuthHeaders: () => ({ Authorization: "Bearer fake-token" }),
}))

describe("UserPanel", () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.setItem("token", "fake-token")
    localStorage.setItem("user", JSON.stringify({ id: "1", username: "testuser", role: "user" }))
  })

  it("renders the greeting and welcome message", async () => {
    getFilesMock.mockReturnValue({ status: 200, data: { files: [] } })
    const { findByText } = renderWithProviders(
      <BrowserRouter>
        <UserPanel />
      </BrowserRouter>
    )

    expect(await findByText(/Good .* testuser/)).toBeTruthy()
    expect(await findByText("What brings you to your files today?")).toBeTruthy()
  })

  it("renders a list of files", async () => {
    const mockFiles = [
      { id: "1", filename: "file1.txt", size: 1024, created_at: 1700000000 },
      { id: "2", filename: "file2.jpg", size: 2048, created_at: 1700000100 },
    ]
    getFilesMock.mockReturnValue({ status: 200, data: { files: mockFiles } })

    const { findByText } = renderWithProviders(
      <BrowserRouter>
        <UserPanel />
      </BrowserRouter>
    )

    expect(await findByText("file1.txt")).toBeTruthy()
    expect(await findByText("file2.jpg")).toBeTruthy()
    expect(await findByText("1KB")).toBeTruthy()
    expect(await findByText("2KB")).toBeTruthy()
  })

  it("opens the upload dialog when clicking upload button", async () => {
    getFilesMock.mockReturnValue({ status: 200, data: { files: [] } })
    const { findByText, getByRole } = renderWithProviders(
      <BrowserRouter>
        <UserPanel />
      </BrowserRouter>
    )

    const uploadButton = getByRole("button", { name: /upload/i })
    uploadButton.click()

    expect(await findByText("Drag and drop a file to upload")).toBeTruthy()
  })

  it("filters files based on search input", async () => {
    const mockFiles = [
      { id: "1", filename: "apple.txt", size: 1024, created_at: 1700000000 },
      { id: "2", filename: "banana.jpg", size: 2048, created_at: 1700000100 },
    ]
    getFilesMock.mockReturnValue({ status: 200, data: { files: mockFiles } })

    const { getByPlaceholderText, queryByText, findByText } = renderWithProviders(
      <BrowserRouter>
        <UserPanel />
      </BrowserRouter>
    )

    expect(await findByText("apple.txt")).toBeTruthy()
    expect(await findByText("banana.jpg")).toBeTruthy()

    const searchInput = getByPlaceholderText("Search uploaded files")
    searchInput.focus()
    const event = new Event("input", { bubbles: true })
    Object.defineProperty(event, "target", { value: { value: "apple" } })
    searchInput.dispatchEvent(new Event("change", { bubbles: true }))
  })
})
