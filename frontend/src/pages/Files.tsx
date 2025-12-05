import { useEffect, useState } from "react";
import api from "../api/axios";

interface FileEntry {
  id: number;
  filename: string;
  size_bytes: number;
  created_at: string;
}

export default function FileList() {
  const [files, setFiles] = useState<FileEntry[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchFiles = async () => {
    try {
      const res = await api.get<FileEntry[]>("/files/list/");
      setFiles(res.data);
    } catch (err) {
      console.error("Failed to fetch files:", err);
    } finally {
      setLoading(false);
    }
  };

  const downloadFile = async (id: number, filename: string) => {
    try {
      const res = await api.get(`/files//download/${id}`, {
        responseType: "blob",
      });

      const url = window.URL.createObjectURL(res.data);
      const a = document.createElement("a");
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      console.error("Download failed:", err);
    }
  };

  useEffect(() => {
    fetchFiles();
  }, []);

  if (loading) return <p>Loading...</p>;

  return (
    <div>
      <h2>Your Files</h2>
      {files.length === 0 && <p>No files uploaded yet.</p>}

      <ul>
        {files.map((file) => (
          <li key={file.id} style={{ marginBottom: "12px" }}>
            <strong>{file.filename}</strong>{" "}
            <small>
              ({(file.size_bytes / 1024).toFixed(1)} KB) —{" "}
              {new Date(file.created_at).toLocaleString()}
            </small>
            <button
              style={{ marginLeft: "12px" }}
              onClick={() => downloadFile(file.id, file.filename)}
            >
              Download
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
