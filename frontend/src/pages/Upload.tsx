import { useState } from "react";
import { uploadFile } from "../api/files";

export default function UploadPage() {
  const [file, setFile] = useState<File | null>(null);
  const [message, setMessage] = useState("");

  const handleUpload = async () => {
    if (!file) return;

    try {
      const res = await uploadFile(file);
      setMessage(`Uploaded: ${res.file.name}`);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setMessage(err.message);
      } else {
        setMessage(String(err));
      }
    }
  };

  return (
    <div style={{ padding: "2rem" }}>
      <h1>Upload File</h1>

      <input type="file" onChange={e => setFile(e.target.files?.[0] || null)} />

      <button onClick={handleUpload}>Upload</button>

      <p>{message}</p>
    </div>
  );
}