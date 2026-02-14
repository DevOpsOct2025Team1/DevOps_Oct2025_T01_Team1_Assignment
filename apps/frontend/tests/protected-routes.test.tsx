import { describe, it, expect, beforeEach, vi } from "vitest"
import { BrowserRouter } from "react-router"
import Dashboard from "../app/routes/dashboard"
import { renderWithProviders } from "./test-utils"
import { clearAuthCache } from "../app/utils/auth"

const createUserMock = vi.fn(() => ({
  mutateAsync: vi.fn(),
  isPending: false,
}))

const deleteUserMock = vi.fn(() => ({
  mutateAsync: vi.fn(),
  isPending: false,
}))

const getUsersMock = vi.fn(() => ({
  data: [],
  isLoading: false,
}))

const getFilesMock = vi.fn(() => ({
  status: 200,
  data: { files: [] },
}))

vi.mock("../app/api/generated", () => ({
  usePostApiAdminCreateUser: () => createUserMock(),
  useDeleteApiAdminDeleteUser: () => deleteUserMock(),
  useGetApiAdminListUsers: () => getUsersMock(),
  getGetApiAdminListUsersQueryKey: () => ["users"],
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
  usePostApiFiles: () => ({ mutate: vi.fn(), isPending: false }),
  useDeleteApiFilesId: () => ({ mutate: vi.fn(), isPending: false }),
  getGetApiFilesIdDownloadUrl: (id: string) => `/api/files/${id}/download`,
}))

vi.mock("../app/utils/chunkedUpload", () => ({
  uploadFileInChunks: vi.fn(),
}))

describe("Protected Routes", () => {
  beforeEach(() => {
    vi.resetModules()
    clearAuthCache()
    localStorage.clear()
    createUserMock.mockClear()
    deleteUserMock.mockClear()
    getUsersMock.mockClear()
    getFilesMock.mockClear()
  })

  it("shows admin dashboard for admin users", async () => {
    localStorage.setItem("token", "fake-token")
    localStorage.setItem("user", JSON.stringify({
      id: "1",
      username: "admin",
      role: "admin"
    }))

    const { findByText } = renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    )

    expect(await findByText(/Admin/i)).toBeTruthy()
    expect(await findByText("What would you like to manage today?")).toBeTruthy()
  })

  it("shows user dashboard for regular users", async () => {
    localStorage.setItem("token", "fake-token")
    localStorage.setItem("user", JSON.stringify({
      id: "1",
      username: "user",
      role: "user"
    }))

    const { findByText } = renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    )

    expect(await findByText(/Good .* user/i)).toBeTruthy()
    expect(await findByText("What brings you to your files today?")).toBeTruthy()
  })
})
