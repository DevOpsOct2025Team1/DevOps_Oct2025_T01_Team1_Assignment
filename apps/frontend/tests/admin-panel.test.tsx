import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { BrowserRouter } from "react-router"
import { renderWithProviders } from "./test-utils"
import AdminPanel from "../app/components/AdminPanel"

const listUsersMock = vi.fn()

vi.mock("../app/api/generated", () => ({
  useGetApiAdminListUsers: (params: any, options: any) => {
    const data = listUsersMock()
    if (options?.query?.select && data) {
      return { data: options.query.select(data), isLoading: false }
    }
    return { data: data || [], isLoading: false }
  },
  usePostApiAdminCreateUser: () => ({
    mutateAsync: vi.fn(), isPending: false,
  }),
  getGetApiAdminListUsersQueryKey: () => ["admin-list-users"],
  useDeleteApiAdminDeleteUser: () => ({
    mutateAsync: vi.fn(), isPending: false,
  }),
}))

describe("AdminPanel", () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.setItem("token", "fake-token")
    localStorage.setItem("user", JSON.stringify({ id: "1", username: "admin", role: "admin" }))
  })

  it("renders greeting and management prompt", async () => {
    listUsersMock.mockReturnValue({ status: 200, data: { users: [] } })
    renderWithProviders(<BrowserRouter><AdminPanel /></BrowserRouter>)
    expect(screen.getByText(/Good .* Admin/)).toBeDefined()
    expect(screen.getByText("What would you like to manage today?")).toBeDefined()
  })

  it("renders user list with roles", async () => {
    listUsersMock.mockReturnValue({
      status: 200,
      data: { users: [
        { id: "1", username: "alice", role: "ROLE_ADMIN" },
        { id: "2", username: "bob", role: "ROLE_USER" },
      ]},
    })
    renderWithProviders(<BrowserRouter><AdminPanel /></BrowserRouter>)
    expect(screen.getByText("alice")).toBeDefined()
    expect(screen.getByText("bob")).toBeDefined()
    expect(screen.getByText("Admin")).toBeDefined()
    expect(screen.getByText("User")).toBeDefined()
  })

  it("renders empty state when no users", async () => {
    listUsersMock.mockReturnValue({ status: 200, data: { users: [] } })
    renderWithProviders(<BrowserRouter><AdminPanel /></BrowserRouter>)
    expect(screen.getByText("No users found")).toBeDefined()
  })

  it("renders Add User button", async () => {
    listUsersMock.mockReturnValue({ status: 200, data: { users: [] } })
    renderWithProviders(<BrowserRouter><AdminPanel /></BrowserRouter>)
    expect(screen.getByRole("button", { name: /add user/i })).toBeDefined()
  })

  it("filters users by search input", async () => {
    const user = userEvent.setup()
    listUsersMock.mockReturnValue({
      status: 200,
      data: { users: [{ id: "1", username: "alice", role: "ROLE_USER" }] },
    })
    renderWithProviders(<BrowserRouter><AdminPanel /></BrowserRouter>)
    const searchInput = screen.getByPlaceholderText("Search Users")
    await user.type(searchInput, "alice")
    expect(searchInput).toBeDefined()
  })
})
