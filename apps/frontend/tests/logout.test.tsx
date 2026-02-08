import { render } from "@testing-library/react"
import { describe, it, expect, beforeEach, vi } from "vitest"
import { BrowserRouter } from "react-router"
import Logout from "../app/routes/logout"
import { AuthProvider } from "../app/contexts/AuthContext"

describe("Logout", () => {
  beforeEach(() => {
    localStorage.setItem("token", "fake-token")
    localStorage.setItem("user", JSON.stringify({
      id: "1",
      username: "testuser",
      role: "user"
    }))
  })

  it("clears localStorage on logout", () => {
    render(
      <BrowserRouter>
        <AuthProvider>
          <Logout />
        </AuthProvider>
      </BrowserRouter>
    )

    expect(localStorage.getItem("token")).toBeNull()
    expect(localStorage.getItem("user")).toBeNull()
  })
})
