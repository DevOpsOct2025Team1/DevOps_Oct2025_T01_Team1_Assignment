import { useState, useCallback } from "react"
import { Search, Upload, X } from "lucide-react"
import { getHours } from "date-fns"
import { useAuth } from "~/contexts/AuthContext"
import { Button } from "~/components/ui/button"
import { Input } from "~/components/ui/input"
import { Dialog, DialogContent, DialogClose } from "~/components/ui/dialog"
import { useGetApiFiles, usePostApiFiles, useDeleteApiFilesId } from "~/api/generated"
import { resolveUrl, getAuthHeaders } from "~/api/orval-client"
import { getGetApiFilesIdDownloadUrl } from "~/api/generated"
import { uploadFileInChunks } from "~/utils/chunkedUpload"
import { FileList } from "./FileList"
import type { InternalHandlersFileMetadata } from "~/api/generated/model"

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "~/components/ui/alert-dialog"

const UserPanel = () => {
  const { user } = useAuth()
  const [isUploadOpen, setIsUploadOpen] = useState(false)
  const [fileToDelete, setFileToDelete] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState("")
  const [uploadProgress, setUploadProgress] = useState<number>(0)
  const [isChunkedUploading, setIsChunkedUploading] = useState(false)

  const getGreeting = () => {
    const hours = getHours(new Date())
    if (hours < 12) return "Good morning"
    if (hours < 18) return "Good afternoon"
    return "Good evening"
  }

  const { data: files = [], isLoading, refetch } = useGetApiFiles(
    {
      query: {
        select: (data) => {
          if (data.status === 200 && data.data) {
             return (data.data.files || []) as InternalHandlersFileMetadata[]
          }
          return []
        }
      }
    }
  )

  const { mutate: uploadFile, isPending: isUploading } = usePostApiFiles({
    mutation: {
      onSuccess: () => {
        setIsUploadOpen(false)
        refetch().then()
      },
    },
  })
  const { mutate: deleteFile } = useDeleteApiFilesId({
    mutation: {
      onSuccess: () => {
        setFileToDelete(null)
        refetch().then()
      },
    },
  })

  const CHUNK_THRESHOLD = 50 * 1024 * 1024

  const handleFileUpload = async (file: File) => {
    if (file.size > CHUNK_THRESHOLD) {
      setIsChunkedUploading(true)
      setUploadProgress(0)

      try {
        await uploadFileInChunks(file, (percentage) => {
          setUploadProgress(percentage)
        })

        setIsUploadOpen(false)
        setUploadProgress(0)
        await refetch()
      } catch (error) {
        console.error("Upload failed:", error)
        alert("Upload failed. Please try again.")
      } finally {
        setIsChunkedUploading(false)
      }
    } else {
      uploadFile({ data: { file } })
    }
  }

  const handleDeleteConfirm = () => {
    if (fileToDelete) {
      deleteFile({ id: fileToDelete })
    }
  }

  const onDrop = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      handleFileUpload(e.dataTransfer.files[0])
      e.dataTransfer.clearData()
    }
  }, [])

  const onDragOver = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
  }, [])

  const onFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      handleFileUpload(e.target.files[0])
    }
  }
  
  const triggerDownload = async (id: string, name: string) => {
      try {
          const downloadUrl = resolveUrl(getGetApiFilesIdDownloadUrl(id));
          const headers = new Headers(getAuthHeaders());
          headers.delete("Content-Type");

          const response = await fetch(downloadUrl, {
            headers,
          });
          if (!response.ok) throw new Error("Download failed");
          const blob = await response.blob();
          const url = window.URL.createObjectURL(blob);
          const a = document.createElement("a");
          a.href = url;
          a.download = name;
          document.body.appendChild(a);
          a.click();
          window.URL.revokeObjectURL(url);
          document.body.removeChild(a);
      } catch (error) {
          console.error("Download error:", error);
      }
  }

  const filteredFiles = files.filter(file => 
    file.filename?.toLowerCase().includes(searchQuery.toLowerCase())
  )

  return (
    <div className="space-y-8">
      <div className="space-y-2">
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-gray-900">
          {getGreeting()} {user?.username || "User"}
        </h1>
        <p className="text-lg md:text-xl text-gray-600">
          What brings you to your files today?
        </p>
      </div>

      <div className="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4">
        <div className="relative flex-1 w-full md:max-w-md">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
          <Input
            placeholder="Search uploaded files"
            className="pl-10 bg-white border-gray-200 w-full"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        <Button
          onClick={() => setIsUploadOpen(true)}
        >
          <Upload className="mr-2 h-4 w-4" />
          Upload
        </Button>
      </div>

      <FileList 
        files={filteredFiles}
        isLoading={isLoading}
        onDelete={(id) => setFileToDelete(id)}
        onDownload={triggerDownload}
      />

      <Dialog open={isUploadOpen} onOpenChange={setIsUploadOpen}>
        <DialogContent className="sm:max-w-xl p-0 overflow-hidden bg-white gap-0">
          <DialogClose className="absolute right-4 top-4 rounded-sm opacity-70">
            <X className="h-6 w-6" />
            <span className="sr-only">Close</span>
          </DialogClose>
          
          <div 
            className="flex flex-col items-center justify-center p-12 min-h-[400px]"
            onDrop={onDrop}
            onDragOver={onDragOver}
          >
            <div className="w-full h-full flex flex-col items-center justify-center">
               <div className="mb-6">
                <Upload className="h-24 w-24 text-gray-900" strokeWidth={1} />
               </div>
               <h3 className="text-2xl font-medium text-gray-900 mb-2">
                 Drag and drop a file to upload
               </h3>
               <p className="text-sm text-gray-500 mb-6">or</p>
               <div className="relative">
                 <input
                   type="file"
                   className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                   onChange={onFileInputChange}
                   disabled={isUploading || isChunkedUploading}
                 />
                 <Button
                     size="lg"
                     disabled={isUploading || isChunkedUploading}
                 >
                   <Upload className="mr-2 h-4 w-4" />
                   {isUploading || isChunkedUploading ? "Uploading..." : "Upload"}
                 </Button>
               </div>
               {isChunkedUploading && (
                 <div className="mt-4 w-full max-w-xs">
                   <div className="flex items-center justify-between mb-2">
                     <span className="text-sm text-gray-600">
                       Uploading... {uploadProgress}%
                     </span>
                   </div>
                   <div className="w-full bg-gray-200 rounded-full h-2">
                     <div
                       className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                       style={{ width: `${uploadProgress}%` }}
                     />
                   </div>
                 </div>
               )}
               <p className="mt-4 text-sm text-gray-500">Max size: 2GB</p>
            </div>
          </div>
        </DialogContent>
      </Dialog>

      <AlertDialog open={!!fileToDelete} onOpenChange={(open) => !open && setFileToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete your file.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteConfirm} className="bg-red-600 hover:bg-red-700">
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}

export default UserPanel
