import { useState } from "react";
import { uploadFile } from "../api/files";

export default function UploadPage() {
  const [file, setFile] = useState<File | null>(null);
  const [message, setMessage] = useState("");
  const [progress, setProgress] = useState<number>(0);

  const handleUpload = async () => {
    if (!file) return;

    setProgress(0);
    setMessage("");

    try {
      const res = await uploadFile(file, (pct) => {
        setProgress(pct);
      });

      setMessage(`Uploaded successfully: ${res.filename}`);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setMessage(err.message);
      } else {
        setMessage(String(err));
      }
    }
  };

  return (
    <div style={{ padding: "2rem", maxWidth: "400px" }}>
      <h1>Upload File</h1>

      <input
        type="file"
        onChange={e => setFile(e.target.files?.[0] || null)}
      />

      <button onClick={handleUpload} disabled={!file}>
        Upload
      </button>

      {/* Progress Bar */}
      {progress > 0 && progress < 100 && (
        <div style={{ marginTop: "1rem", width: "100%" }}>
          <div
            style={{
              height: "20px",
              background: "#ddd",
              borderRadius: "4px",
              overflow: "hidden"
            }}
          >
            <div
              style={{
                width: `${progress}%`,
                height: "100%",
                background: "#4caf50",
                transition: "width 0.2s"
              }}
            />
          </div>
          <p>{progress}%</p>
        </div>
      )}

      {progress === 100 && (
        <p style={{ color: "green" }}>Upload complete!</p>
      )}

      <p>{message}</p>
    </div>
  );
}
