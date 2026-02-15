import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { BrowserRouter } from "react-router"
import { renderWithProviders } from "./test-utils"
import { UserRowActions } from "~/components/UserRowActions"

const mutateAsyncMock = vi.fn()

vi.mock("../app/api/generated", () => ({
  useDeleteApiAdminDeleteUser: () => ({
    mutateAsync: mutateAsyncMock,
    isPending: false,
  }),
  getGetApiAdminListUsersQueryKey: () => ["admin-list-users"],
}))

describe("UserRowActions", () => {
  const mockUser = { id: "user-1", username: "testuser", role: "ROLE_USER" }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("renders the menu trigger button", () => {
    renderWithProviders(
      <BrowserRouter><UserRowActions user={mockUser} /></BrowserRouter>
    )
    expect(screen.getByRole("button", { name: /open menu/i })).toBeDefined()
  })

  it("opens dropdown menu on click", async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <BrowserRouter><UserRowActions user={mockUser} /></BrowserRouter>
    )
    await user.click(screen.getByRole("button", { name: /open menu/i }))
    expect(screen.getByText("Delete User")).toBeDefined()
  })

  it("shows delete confirmation dialog", async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <BrowserRouter><UserRowActions user={mockUser} /></BrowserRouter>
    )
    await user.click(screen.getByRole("button", { name: /open menu/i }))
    await user.click(screen.getByText("Delete User"))
    expect(screen.getByText("Are you absolutely sure?")).toBeDefined()
    expect(screen.getByText(/testuser/)).toBeDefined()
  })

  it("calls delete mutation when confirmed", async () => {
    mutateAsyncMock.mockResolvedValueOnce({})
    const user = userEvent.setup()
    renderWithProviders(
      <BrowserRouter><UserRowActions user={mockUser} /></BrowserRouter>
    )
    await user.click(screen.getByRole("button", { name: /open menu/i }))
    await user.click(screen.getByText("Delete User"))
    await user.click(screen.getByRole("button", { name: /^delete$/i }))
    expect(mutateAsyncMock).toHaveBeenCalledWith({
      data: { id: "user-1" },
    })
  })

  it("closes dialog when cancel is clicked", async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <BrowserRouter><UserRowActions user={mockUser} /></BrowserRouter>
    )
    await user.click(screen.getByRole("button", { name: /open menu/i }))
    await user.click(screen.getByText("Delete User"))
    expect(screen.getByText("Are you absolutely sure?")).toBeDefined()
    await user.click(screen.getByRole("button", { name: /cancel/i }))
    expect(screen.queryByText("Are you absolutely sure?")).toBeNull()
  })

  it("handles delete error gracefully", async () => {
    const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {})
    mutateAsyncMock.mockRejectedValueOnce(new Error("Delete failed"))
    const user = userEvent.setup()
    renderWithProviders(
      <BrowserRouter><UserRowActions user={mockUser} /></BrowserRouter>
    )
    await user.click(screen.getByRole("button", { name: /open menu/i }))
    await user.click(screen.getByText("Delete User"))
    await user.click(screen.getByRole("button", { name: /^delete$/i }))
    expect(mutateAsyncMock).toHaveBeenCalled()
    consoleSpy.mockRestore()
  })
})
