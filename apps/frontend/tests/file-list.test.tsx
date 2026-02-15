import { describe, it, expect, vi } from "vitest"
import { screen } from "@testing-library/react"
import { renderWithProviders } from "./test-utils"
import { FileList } from "~/components/FileList"

describe("FileList", () => {
  const onDelete: (id: string) => void = vi.fn()
  const onDownload: (id: string, name: string) => void = vi.fn()

  it("shows loading state", () => {
    renderWithProviders(
      <FileList files={[]} isLoading={true} onDelete={onDelete} onDownload={onDownload} />
    )
    expect(screen.getByText("Loading files...")).toBeDefined()
  })

  it("shows empty state when no files", () => {
    renderWithProviders(
      <FileList files={[]} isLoading={false} onDelete={onDelete} onDownload={onDownload} />
    )
    expect(screen.getByText("No files uploaded yet")).toBeDefined()
  })

  it("renders file items when files are provided", () => {
    const files = [
      { id: "1", filename: "file1.txt", size: 1024, created_at: 1700000000 },
      { id: "2", filename: "file2.jpg", size: 2048, created_at: 1700000100 },
    ]
    renderWithProviders(
      <FileList files={files} isLoading={false} onDelete={onDelete} onDownload={onDownload} />
    )
    expect(screen.getByText("file1.txt")).toBeDefined()
    expect(screen.getByText("file2.jpg")).toBeDefined()
  })

  it("does not show loading or empty state when files exist", () => {
    const files = [
      { id: "1", filename: "file1.txt", size: 1024, created_at: 1700000000 },
    ]
    renderWithProviders(
      <FileList files={files} isLoading={false} onDelete={onDelete} onDownload={onDownload} />
    )
    expect(screen.queryByText("Loading files...")).toBeNull()
    expect(screen.queryByText("No files uploaded yet")).toBeNull()
  })

  it("does not show files while loading", () => {
    const files = [
      { id: "1", filename: "file1.txt", size: 1024, created_at: 1700000000 },
    ]
    renderWithProviders(
      <FileList files={files} isLoading={true} onDelete={onDelete} onDownload={onDownload} />
    )
    expect(screen.getByText("Loading files...")).toBeDefined()
    expect(screen.queryByText("file1.txt")).toBeNull()
  })
})
