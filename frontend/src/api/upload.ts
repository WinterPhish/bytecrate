import { authenticatedFetch } from "./auth";

export async function uploadFile(file: File) {

  const formData = new FormData();
  formData.append("file", file);

  const res = await authenticatedFetch("http://localhost:8080/api/files/upload", {
    method: "POST",
    body: formData,
  });

  if (!res.ok) {
    throw new Error(await res.text());
  }

  return await res.json();
}