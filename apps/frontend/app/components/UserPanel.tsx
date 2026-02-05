import { useState } from "react";
import { Upload, X, MoreVertical, Download, Trash2 } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "./ui/card";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "./ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";

interface FileItem {
  id: string;
  name: string;
  size: string;
  timestamp: string;
}

// TODO: DEPLOYMENT - REMOVE THIS SAMPLE DATA
// Replace with actual API call to fetch user's files from backend
// Example: const [files, setFiles] = useState<FileItem[]>([]);
// useEffect(() => { fetchUserFiles().then(setFiles); }, []);
const sampleFiles: FileItem[] = [
  { id: "1", name: "Presentation.pptx", size: "10MB", timestamp: "15:00" },
  { id: "2", name: "File.txt", size: "160MB", timestamp: "14:10" },
  { id: "3", name: "Image.png", size: "2MB", timestamp: "08:45" },
  { id: "4", name: "Document.docx", size: "1.2MB", timestamp: "Yesterday 12:30" },
  { id: "5", name: "Spreadsheet.xlsx", size: "500KB", timestamp: "2026/02/02 09:15" },
];

export default function UserPanel() {
  const [showUploadDialog, setShowUploadDialog] = useState(false);
  const [showFeatureDialog, setShowFeatureDialog] = useState(false);
  const [dialogContent, setDialogContent] = useState({ title: "", description: "" });
  const [searchQuery, setSearchQuery] = useState("");

  const handleUploadClick = () => {
    setShowUploadDialog(true);
  };

  const handleFeatureClick = (feature: string) => {
    setDialogContent({
      title: `${feature} - Coming Soon`,
      description: `The ${feature.toLowerCase()} feature is not yet available. The backend API endpoint for this operation needs to be implemented.`,
    });
    setShowFeatureDialog(true);
    setShowUploadDialog(false);
  };

  // TODO: ENDPOINTS - Implement actual file download
  // Replace with API call: await fileApi.download(file.id)
  const handleDownload = (file: FileItem) => {
    handleFeatureClick("Download File");
  };

  // TODO: ENDPOINTS - Implement actual file deletion
  // Replace with API call: await fileApi.delete(file.id)
  const handleDelete = (file: FileItem) => {
    handleFeatureClick("Delete File");
  };

  const filteredFiles = sampleFiles.filter(file =>
    file.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <>
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex-1 max-w-md">
              <div className="relative">
                <Input
                  type="text"
                  placeholder="Search uploaded files"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
                <svg
                  className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  />
                </svg>
              </div>
            </div>
            <Button onClick={handleUploadClick} className="ml-4">
              <Upload className="h-4 w-4 mr-2" />
              Upload
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-1">
            {filteredFiles.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                No files found
              </div>
            ) : (
              filteredFiles.map((file) => (
                <div
                  key={file.id}
                  className="flex items-center justify-between py-3 px-4 hover:bg-gray-50 rounded-md transition-colors"
                >
                  <div className="flex-1">
                    <span className="text-gray-900 font-medium">{file.name}</span>
                  </div>
                  <div className="flex items-center gap-6 text-sm text-gray-600">
                    <span className="w-20 text-right">{file.size}</span>
                    <span className="w-40 text-right">{file.timestamp}</span>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                          <span className="sr-only">Open menu</span>
                          <MoreVertical className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onSelect={() => handleDownload(file)}>
                          <Download className="h-4 w-4 mr-2" />
                          Download
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          onSelect={() => handleDelete(file)}
                          className="text-red-600 hover:text-red-700"
                        >
                          <Trash2 className="h-4 w-4 mr-2" />
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </div>
              ))
            )}
          </div>
        </CardContent>
      </Card>

      <Dialog open={showUploadDialog} onOpenChange={setShowUploadDialog}>
        <DialogContent className="sm:max-w-md">
          <button
            onClick={() => setShowUploadDialog(false)}
            className="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100"
          >
            <X className="h-4 w-4" />
            <span className="sr-only">Close</span>
          </button>
          <div className="flex flex-col items-center space-y-4 pt-4">
            <div className="w-16 h-16 rounded-full border-2 border-gray-300 flex items-center justify-center">
              <Upload className="h-8 w-8 text-gray-600" />
            </div>
            <DialogHeader>
              <DialogTitle>Drag and drop a file to upload</DialogTitle>
            </DialogHeader>
          </div>
          <div className="flex flex-col items-center space-y-4 pb-4">
            <Button onClick={() => handleFeatureClick("Upload File")}>
              <Upload className="h-4 w-4 mr-2" />
              Upload
            </Button>
            <p className="text-xs text-gray-500">Max size: 2GB</p>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={showFeatureDialog} onOpenChange={setShowFeatureDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{dialogContent.title}</DialogTitle>
            <DialogDescription>{dialogContent.description}</DialogDescription>
          </DialogHeader>
          <div className="flex justify-end">
            <Button onClick={() => setShowFeatureDialog(false)}>
              OK
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}