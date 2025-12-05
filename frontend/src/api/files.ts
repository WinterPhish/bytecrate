import api from "./axios";

export async function uploadFile(file: File) {
  const formData = new FormData();
  formData.append("file", file);

  try {
    const res = await api.post("/files/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      onUploadProgress(progressEvent) {
        const percentCompleted = Math.round(
          (progressEvent.loaded * 100) / progressEvent.total!
        );
        console.log(`Upload progress: ${percentCompleted}%`);
      },
    });

    return res.data; // automatically parses JSON
  } catch (err: unknown) {
    console.error("Upload failed:", err);
    if (err instanceof Error) {
    throw new Error(err.message || "Upload failed");
    }
  }
}