import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { BrowserRouter } from "react-router"
import { renderWithProviders } from "./test-utils"
import { CreateUserDialog } from "~/components/CreateUserDialog"

const mutateAsyncMock = vi.fn()

vi.mock("../app/api/generated", () => ({
  usePostApiAdminCreateUser: () => ({
    mutateAsync: mutateAsyncMock, isPending: false,
  }),
  getGetApiAdminListUsersQueryKey: () => ["admin-list-users"],
}))

describe("CreateUserDialog", () => {
  beforeEach(() => { vi.clearAllMocks() })

  it("renders the Add User trigger button", () => {
    renderWithProviders(<BrowserRouter><CreateUserDialog /></BrowserRouter>)
    expect(screen.getByRole("button", { name: /add user/i })).toBeDefined()
  })

  it("opens dialog when trigger is clicked", async () => {
    const user = userEvent.setup()
    renderWithProviders(<BrowserRouter><CreateUserDialog /></BrowserRouter>)
    await user.click(screen.getByRole("button", { name: /add user/i }))
    expect(screen.getByText("Create New User")).toBeDefined()
    expect(screen.getByLabelText(/username/i)).toBeDefined()
    expect(screen.getByLabelText(/password/i)).toBeDefined()
  })

  it("shows validation error for empty fields", async () => {
    const user = userEvent.setup()
    renderWithProviders(<BrowserRouter><CreateUserDialog /></BrowserRouter>)
    await user.click(screen.getByRole("button", { name: /add user/i }))
    await user.click(screen.getByRole("button", { name: /create user/i }))
    expect(screen.getByText("Username and password are required")).toBeDefined()
  })

  it("submits form with valid data", async () => {
    mutateAsyncMock.mockResolvedValueOnce({})
    const user = userEvent.setup()
    renderWithProviders(<BrowserRouter><CreateUserDialog /></BrowserRouter>)
    await user.click(screen.getByRole("button", { name: /add user/i }))
    await user.type(screen.getByLabelText(/username/i), "newuser")
    await user.type(screen.getByLabelText(/password/i), "password123")
    await user.click(screen.getByRole("button", { name: /create user/i }))
    expect(mutateAsyncMock).toHaveBeenCalledWith({
      data: { username: "newuser", password: "password123" },
    })
  })

  it("shows error on failed submission", async () => {
    mutateAsyncMock.mockRejectedValueOnce(new Error("Server error"))
    const user = userEvent.setup()
    renderWithProviders(<BrowserRouter><CreateUserDialog /></BrowserRouter>)
    await user.click(screen.getByRole("button", { name: /add user/i }))
    await user.type(screen.getByLabelText(/username/i), "newuser")
    await user.type(screen.getByLabelText(/password/i), "password123")
    await user.click(screen.getByRole("button", { name: /create user/i }))
    expect(await screen.findByText("Server error")).toBeDefined()
  })
})
