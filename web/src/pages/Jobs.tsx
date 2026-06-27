import { useState, useEffect } from "react";
import { supabase } from "../lib/supabase";

interface Job {
  id: string;
  type: string;
  status: string;
  attempts: number;
  created_at: string;
}

function Jobs() {
  const [jobs, setJobs] = useState<Job[]>([]);

useEffect(() => {
  // 1. carga inicial (lo que ya tenés)
  supabase
    .from("jobs")
    .select("id, type, status, attempts, created_at")
    .order("created_at", { ascending: false })
    .then(({ data }) => {
      if (data) setJobs(data as Job[]);
    });

  // 2. suscripción a cambios en vivo
  const channel = supabase
    .channel("jobs-changes")
    .on(
      "postgres_changes",
      { event: "*", schema: "public", table: "jobs" },
      (payload) => {
        console.log("cambio en jobs:", payload);
        // recargar la lista cuando algo cambia
        supabase
          .from("jobs")
          .select("id, type, status, attempts, created_at")
          .order("created_at", { ascending: false })
          .then(({ data }) => {
            if (data) setJobs(data as Job[]);
          });
      }
    )
    .subscribe();

  // 3. limpieza: cerrar el canal al desmontar
  return () => {
    supabase.removeChannel(channel);
  };
}, []);

  return (
    <div>
      <h1>Jobs</h1>
      <ul>
        {jobs.map((job) => (
          <li key={job.id}>
            {job.type} — <strong>{job.status}</strong> (intentos: {job.attempts})
          </li>
        ))}
      </ul>
    </div>
  );
}

export default Jobs;