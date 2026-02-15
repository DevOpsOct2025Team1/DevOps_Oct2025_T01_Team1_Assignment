import { describe, it, expect, vi, beforeEach } from "vitest"
import userEvent from "@testing-library/user-event"
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

    const { getByPlaceholderText, findByText, queryByText } = renderWithProviders(
      <BrowserRouter>
        <UserPanel />
      </BrowserRouter>
    )

    expect(await findByText("apple.txt")).toBeTruthy()
    expect(await findByText("banana.jpg")).toBeTruthy()

    const searchInput = getByPlaceholderText("Search uploaded files")
    await userEvent.type(searchInput, "apple")

    expect(queryByText("apple.txt")).toBeTruthy()
    expect(queryByText("banana.jpg")).toBeNull()
  })

  it("shows empty state when no files", async () => {
    getFilesMock.mockReturnValue({ status: 200, data: { files: [] } })
    const { findByText } = renderWithProviders(
      <BrowserRouter><UserPanel /></BrowserRouter>
    )
    expect(await findByText("No files uploaded yet")).toBeTruthy()
  })

  it("opens delete confirmation dialog when delete is triggered", async () => {
    const mockFiles = [
      { id: "1", filename: "file1.txt", size: 1024, created_at: 1700000000 },
    ]
    getFilesMock.mockReturnValue({ status: 200, data: { files: mockFiles } })
    const { findByText, getByRole } = renderWithProviders(
      <BrowserRouter><UserPanel /></BrowserRouter>
    )
    expect(await findByText("file1.txt")).toBeTruthy()

    const menuButton = getByRole("button", { name: /menu/i })
    await userEvent.click(menuButton)
    await userEvent.click(await findByText("Delete"))

    expect(await findByText("Are you sure?")).toBeTruthy()
    expect(await findByText(/permanently delete your file/)).toBeTruthy()
  })

  it("confirms file deletion", async () => {
    const mockFiles = [
      { id: "1", filename: "file1.txt", size: 1024, created_at: 1700000000 },
    ]
    getFilesMock.mockReturnValue({ status: 200, data: { files: mockFiles } })
    const { findByText, getByRole } = renderWithProviders(
      <BrowserRouter><UserPanel /></BrowserRouter>
    )
    expect(await findByText("file1.txt")).toBeTruthy()

    const menuButton = getByRole("button", { name: /menu/i })
    await userEvent.click(menuButton)
    await userEvent.click(await findByText("Delete"))
    await userEvent.click(getByRole("button", { name: /^delete$/i }))

    expect(deleteFileMock).toHaveBeenCalledWith({ id: "1" })
  })

  it("cancels file deletion", async () => {
    const mockFiles = [
      { id: "1", filename: "file1.txt", size: 1024, created_at: 1700000000 },
    ]
    getFilesMock.mockReturnValue({ status: 200, data: { files: mockFiles } })
    const { findByText, getByRole, queryByText } = renderWithProviders(
      <BrowserRouter><UserPanel /></BrowserRouter>
    )
    expect(await findByText("file1.txt")).toBeTruthy()

    const menuButton = getByRole("button", { name: /menu/i })
    await userEvent.click(menuButton)
    await userEvent.click(await findByText("Delete"))
    expect(await findByText("Are you sure?")).toBeTruthy()

    await userEvent.click(getByRole("button", { name: /cancel/i }))
    expect(queryByText("Are you sure?")).toBeNull()
  })

  it("handles non-200 status", async () => {
    getFilesMock.mockReturnValue({ status: 500, data: { files: [] } })
    const { findByText } = renderWithProviders(
      <BrowserRouter><UserPanel /></BrowserRouter>
    )
    expect(await findByText("No files uploaded yet")).toBeTruthy()
  })

  it("renders upload button", async () => {
    getFilesMock.mockReturnValue({ status: 200, data: { files: [] } })
    const { getByRole } = renderWithProviders(
      <BrowserRouter><UserPanel /></BrowserRouter>
    )
    expect(getByRole("button", { name: /upload/i })).toBeDefined()
  })

  it("displays greeting with username from auth context", async () => {
    getFilesMock.mockReturnValue({ status: 200, data: { files: [] } })
    const { findByText } = renderWithProviders(
      <BrowserRouter><UserPanel /></BrowserRouter>
    )
    expect(await findByText(/testuser/)).toBeTruthy()
  })

  it("displays file description text", async () => {
    getFilesMock.mockReturnValue({ status: 200, data: { files: [] } })
    const { findByText } = renderWithProviders(
      <BrowserRouter><UserPanel /></BrowserRouter>
    )
    expect(await findByText("What brings you to your files today?")).toBeTruthy()
  })
})
