import { useEffect, useState } from "react";
import api from "../api/axios";
import axios from "axios";

interface FileEntry {
  id: number;
  filename: string;
  size_bytes: number;
  created_at: string;
  content_type: string;
}

export default function FileList() {
  const [files, setFiles] = useState<FileEntry[]>([]);
  const [loading, setLoading] = useState(true);

  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [previewType, setPreviewType] = useState<string | null>(null);

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

  const deleteFile = async (id: number) => {
    try {
      const res = await api.delete(`/files/${id}`);
      return res.data;
    } catch (err) {
      if (axios.isAxiosError(err)) {
        throw new Error(err.response?.data?.error || "Delete failed");
      }
      throw err;
    }
  }

  const renameFile = async (id: number, newName: string) => {
    try {
      const res = await api.put(`/files/${id}/rename`, {
        filename: newName,
      });
      return res.data;
    } catch (err) {
      if (axios.isAxiosError(err)) {
        throw new Error(err.response?.data?.error || "Rename failed");
    }
    throw err;
    }
  }

  // Possible enhancement
  const downloadFile = async (id: number, filename: string) => {
    try {
      const res = await api.get(`/files/${id}/download`, {
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

  const previewFile = async (id: number): Promise<string> => {
    const res = await api.get(`/files/${id}/preview`, {
      responseType: "blob",
    });

    const blob = new Blob([res.data], { type: res.headers["content-type"] });
    return URL.createObjectURL(blob);
  }

  const handlePreview = async (file: FileEntry) => {
    try {
      const url = await previewFile(file.id);
      setPreviewUrl(url);
      setPreviewType(file.content_type);
    } catch (err) {
      console.error(err);
      alert("Preview not available.");
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
  onClick={() => handlePreview(file)}
>
  Preview
</button>
            <button
              style={{ marginLeft: "12px" }}
              onClick={() => downloadFile(file.id, file.filename)}
            >
              Download
            </button>

            {/* Rename */}
            <button
              style={{ marginLeft: "8px" }}
              onClick={async () => {
                const newName = prompt("New filename:", file.filename);
                if (!newName) return;

                try {
                  await renameFile(file.id, newName);
                  await fetchFiles(); // auto refresh
                } catch (err) {
                  alert("Rename failed: " + err);
                }
              }}
            >
              Rename
            </button>

            {/* Delete */}
            <button
              style={{ marginLeft: "8px", color: "red" }}
              onClick={async () => {
                if (!confirm("Delete this file?")) return;

                try {
                  await deleteFile(file.id);
                  await fetchFiles(); // auto refresh
                } catch (err) {
                  alert("Delete failed: " + err);
                }
              }}
            >
              Delete
            </button>
          </li>
        ))}
      </ul>
         {previewUrl && (
        <div
          onClick={() => setPreviewUrl(null)}
          style={{
            position: "fixed",
            inset: 0,
            background: "rgba(0,0,0,0.65)",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            zIndex: 999,
          }}
        >
          <div
            onClick={(e) => e.stopPropagation()}
            style={{
              background: "white",
              padding: "20px",
              borderRadius: "10px",
              maxWidth: "80%",
              maxHeight: "80%",
              overflow: "auto",
            }}
          >
            <h3>Preview</h3>

            {/* IMAGES */}
            {previewType?.startsWith("image") && (
              <img
                src={previewUrl}
                style={{ maxWidth: "100%", maxHeight: "70vh" }}
              />
            )}

            {/* PDF */}
            {previewType === "application/pdf" && (
              <iframe
                src={previewUrl}
                width="800"
                height="600"
                style={{ border: "none" }}
              ></iframe>
            )}

            {/* TEXT */}
            {previewType?.startsWith("text") && (
              <iframe
                src={previewUrl}
                style={{ width: "600px", height: "500px" }}
              ></iframe>
            )}

            {/* VIDEO */}
            {previewType?.startsWith("video") && (
              <video
                controls
                src={previewUrl}
                style={{ maxWidth: "100%", maxHeight: "70vh" }}
              ></video>
            )}

            {/* AUDIO */}
            {previewType?.startsWith("audio") && (
              <audio controls src={previewUrl}></audio>
            )}

            {/* FALLBACK */}
            {!previewType?.match(/image|pdf|text|video|audio/) && (
              <p>No preview available for this file type.</p>
            )}

            <button
              style={{ marginTop: "10px" }}
              onClick={() => setPreviewUrl(null)}
            >
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
