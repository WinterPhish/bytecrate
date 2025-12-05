import { Link } from "react-router-dom";
import { useState, useEffect } from "react";
import axios from "axios";
import api from "../api/axios";

type StatusResponse = {
  status: string;
};

export default function Dashboard() {
  const [data, setData] = useState<StatusResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

useEffect(() => {
  let isMounted = true;

  api.get("/dev/status")
    .then((res) => {
      if (isMounted) {
        setData(res.data);
        setLoading(false);
      }
    })
    .catch((err: unknown) => {
      if (isMounted) {
        if (axios.isAxiosError(err)) {
          setError(err.response?.data || err.message);
        } else if (err instanceof Error) {
          setError(err);
        } else {
          setError(new Error("Unknown error occurred"));
        }
        setLoading(false);
      }
    });

  return () => {
    isMounted = false;
  };
}, []);

  if (error) {
  return <div>Error: {error.message}</div>;
  }

  if (loading) {
	return <div>Loading</div>;
  }

  return (
    <div style={{ padding: 20 }}>
  
  <h1>Dashboard</h1>
    <p>Server Status: {data?.status}</p>
  <div style={{ marginTop: 20 }}>
        <Link to="/files">View My Files</Link><br/>
        <Link to="/upload">Upload New File</Link><br/>
        <Link to="/account">Account Settings</Link>
      </div>
    </div>
  );
}