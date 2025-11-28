import { Link } from "react-router-dom";

export default function Dashboard() {
  return (
    <div style={{ padding: 20 }}>

      <div style={{ marginTop: 20 }}>
        <Link to="/files">View My Files</Link><br/>
        <Link to="/upload">Upload New File</Link><br/>
        <Link to="/account">Account Settings</Link>
      </div>
    </div>
  );
}