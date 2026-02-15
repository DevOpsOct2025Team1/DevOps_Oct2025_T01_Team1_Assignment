import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { renderWithProviders } from "./test-utils"
import { FileItem } from "~/components/FileItem"

describe("FileItem", () => {
  const mockFile = {
    id: "file-1",
    filename: "document.pdf",
    size: 2048,
    created_at: 1700000000,
  }
  let onDelete: (id: string) => void
  let onDownload: (id: string, name: string) => void

  beforeEach(() => {
    onDelete = vi.fn<(id: string) => void>()
    onDownload = vi.fn<(id: string, name: string) => void>()
  })

  it("renders file name", () => {
    renderWithProviders(
      <FileItem file={mockFile} onDelete={onDelete} onDownload={onDownload} />
    )
    expect(screen.getByText("document.pdf")).toBeDefined()
  })

  it("renders formatted file size", () => {
    renderWithProviders(
      <FileItem file={mockFile} onDelete={onDelete} onDownload={onDownload} />
    )
    expect(screen.getByText("2KB")).toBeDefined()
  })

  it("renders formatted date", () => {
    renderWithProviders(
      <FileItem file={mockFile} onDelete={onDelete} onDownload={onDownload} />
    )
    // created_at is multiplied by 1000 in the component
    const expectedDate = new Date(1700000000 * 1000).toLocaleString()
    expect(screen.getByText(expectedDate)).toBeDefined()
  })

  it("renders dash when created_at is missing", () => {
    const fileWithoutDate = { id: "file-2", filename: "test.txt", size: 100 }
    renderWithProviders(
      <FileItem file={fileWithoutDate} onDelete={onDelete} onDownload={onDownload} />
    )
    expect(screen.getByText("-")).toBeDefined()
  })

  it("calls onDownload when filename is clicked", async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <FileItem file={mockFile} onDelete={onDelete} onDownload={onDownload} />
    )
    await user.click(screen.getByText("document.pdf"))
    expect(onDownload).toHaveBeenCalledWith("file-1", "document.pdf")
  })

  it("does not call onDownload when file has no id", async () => {
    const user = userEvent.setup()
    const fileNoId = { filename: "test.txt", size: 100 }
    renderWithProviders(
      <FileItem file={fileNoId} onDelete={onDelete} onDownload={onDownload} />
    )
    await user.click(screen.getByText("test.txt"))
    expect(onDownload).not.toHaveBeenCalled()
  })

  it("calls onDelete from dropdown menu", async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <FileItem file={mockFile} onDelete={onDelete} onDownload={onDownload} />
    )
    await user.click(screen.getByRole("button", { name: /menu/i }))
    await user.click(screen.getByText("Delete"))
    expect(onDelete).toHaveBeenCalledWith("file-1")
  })

  it("renders zero bytes correctly", () => {
    const emptyFile = { id: "file-3", filename: "empty.txt", size: 0 }
    renderWithProviders(
      <FileItem file={emptyFile} onDelete={onDelete} onDownload={onDownload} />
    )
    expect(screen.getByText("0 Bytes")).toBeDefined()
  })
})
