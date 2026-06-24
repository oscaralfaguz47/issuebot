import { useState, useEffect } from "react";
import { getProjects } from "../lib/api";
import type { Project } from "../types";

function Dashboard() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getProjects()
      .then((data) => setProjects(data))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p>Loading projects...</p>;
  if (error) return <p>Error: {error}</p>;

  return (
    <div>
      <h1>Dashboard</h1>
      <h2>Your projects ({projects.length})</h2>
      <ul>
        {projects.map((p) => (
          <li key={p.ID}>{p.Name}</li>
        ))}
      </ul>
    </div>
  );
}

export default Dashboard;