import { Link } from "react-router-dom";
import { authenticatedFetch } from "../api/auth";
import { useState, useEffect } from "react";

type StatusResponse = {
  status: string;
};

export default function Dashboard() {
  const [data, setData] = useState<StatusResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    authenticatedFetch('http://localhost:8080/api/dev/status')
      .then((response) => {
        if (!response.ok) {
          throw new Error('Error response');
        }
        return response.json();
      })
      .then((data) => {
        setData(data);
        setLoading(false);
      })
      .catch((error) => {
        setError(error);
        setLoading(false);
      });
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