import {MoreVertical, Trash} from "lucide-react"
import { Button } from "~/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu"
import { formatBytes, formatDate } from "~/lib/utils"
import type { InternalHandlersFileMetadata } from "~/api/generated/model"

interface FileItemProps {
  file: InternalHandlersFileMetadata
  onDelete: (id: string) => void
  onDownload: (id: string, name: string) => void
}

export const FileItem = ({ file, onDelete, onDownload }: FileItemProps) => {
  return (
    <div
      className="flex items-center justify-between py-4 border-b border-gray-100 last:border-0 hover:bg-gray-50 px-4"
    >
      <div>
        <button
            onClick={() => file.id && file.filename && onDownload(file.id, file.filename)}
            className="flex-1 font-semibold text-gray-900 truncate pr-4 text-left hover:text-blue-600 transition-colors"
        >
          {file.filename}
        </button>
      </div>
      <div className="flex items-center gap-2 sm:gap-8 text-sm text-gray-900 font-medium shrink-0">
        <span className="w-20 text-right hidden sm:block">
          {formatBytes(file.size || 0)}
        </span>
        <span className="w-40 text-right hidden md:block">
          {file.created_at ? formatDate(file.created_at * 1000) : "-"}
        </span>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <MoreVertical className="h-4 w-4" />
              <span className="sr-only">Menu</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem
              className="text-red-500"
              onClick={() => file.id && onDelete(file.id)}
            >
              <Trash className="h-4 w-4" />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  )
}
