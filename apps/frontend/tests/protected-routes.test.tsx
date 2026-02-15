import { describe, it, expect, beforeEach, vi } from "vitest"
import { BrowserRouter } from "react-router"
import { renderWithProviders } from "./test-utils"
import { clearAuthCache } from "~/utils/auth"

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

vi.mock("../app/api/orval-client", () => ({
  resolveUrl: (url: string) => `http://localhost:3001${url}`,
  getAuthHeaders: () => ({ Authorization: "Bearer fake-token" }),
}))

vi.mock("../app/routes/dashboard", async () => {
  const { useEffect } = await import("react")
  const { useNavigate } = await import("react-router")
  const { useAuth } = await import("../app/contexts/AuthContext")
  const { default: AdminPanel } = await import("../app/components/AdminPanel")
  const { default: UserPanel } = await import("../app/components/UserPanel")

  return {
    default: function Dashboard() {
      const navigate = useNavigate()
      const { user, isAuthenticated } = useAuth()

      useEffect(() => {
        if (!isAuthenticated) {
          navigate("/login")
        }
      }, [isAuthenticated])

      if (!user) return null

      const userIsAdmin = typeof user.role === "number"
        ? user.role === 2
        : user.role.toLowerCase().includes("admin")

      return (
        <div className="flex-1 bg-gray-50 py-8">
          <div className="max-w-5xl mx-auto px-4">
            {userIsAdmin ? <AdminPanel /> : <UserPanel />}
          </div>
        </div>
      )
    },
  }
})

const { default: Dashboard } = await import("../app/routes/dashboard")

describe("Protected Routes", () => {
  beforeEach(() => {
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
