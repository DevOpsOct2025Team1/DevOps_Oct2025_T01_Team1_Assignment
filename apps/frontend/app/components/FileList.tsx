import { FileText } from "lucide-react"
import type { InternalHandlersFileMetadata } from "~/api/generated/model"
import { FileItem } from "./FileItem"

interface FileListProps {
  files: InternalHandlersFileMetadata[]
  isLoading: boolean
  onDelete: (id: string) => void
  onDownload: (id: string, name: string) => void
}

export const FileList = ({ files, isLoading, onDelete, onDownload }: FileListProps) => {
  if (isLoading) {
    return (
      <div className="bg-white rounded-lg p-8 text-center text-gray-500">
        Loading files...
      </div>
    )
  }

  if (files.length === 0) {
    return (
      <div className="bg-white rounded-lg p-8 text-center text-gray-500 flex flex-col items-center">
        <FileText className="h-12 w-12 text-gray-300 mb-4" />
        <p>No files uploaded yet</p>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg">
      {files.map((file) => (
        <FileItem
          key={file.id}
          file={file}
          onDelete={onDelete}
          onDownload={onDownload}
        />
      ))}
    </div>
  )
}
