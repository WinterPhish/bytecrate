import api from "./axios";

export async function uploadFile(file: File, onProgress?: (pct: number) => void) {
  const formData = new FormData();
  formData.append("file", file);

  try {
  const res = await api.post("/files/upload", formData, {
    headers: { "Content-Type": "multipart/form-data" },
    onUploadProgress: (event) => {
      if (!event.total) return;
      const pct = Math.round((event.loaded / event.total) * 100);
      if (onProgress) onProgress(pct);
    }
  });

    return res.data; // automatically parses JSON
  } catch (err: unknown) {
    console.error("Upload failed:", err);
    if (err instanceof Error) {
    throw new Error(err.message || "Upload failed");
    }
  }
}