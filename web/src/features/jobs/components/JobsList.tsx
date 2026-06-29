import { useEffect } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { supabase } from "@/lib/supabase";
import LoadingScreen from "@/components/LoadingScreen";
import type { Job } from "../types";

// Query function: fetches jobs from Supabase
async function fetchJobs(): Promise<Job[]> {
  const { data, error } = await supabase
    .from("jobs")
    .select("id, type, status, attempts, created_at")
    .order("created_at", { ascending: false });

  if (error) throw error;
  return data as Job[];
}

function JobsList() {
  const queryClient = useQueryClient();

  const { data: jobs, isPending, error } = useQuery({
    queryKey: ["jobs"],
    queryFn: fetchJobs,
  });

  // Realtime subscription: on any change, invalidate the cache so RQ refetches
  useEffect(() => {
    const channel = supabase
      .channel("jobs-changes")
      .on(
        "postgres_changes",
        { event: "*", schema: "public", table: "jobs" },
        () => {
          queryClient.invalidateQueries({ queryKey: ["jobs"] });
        }
      )
      .subscribe();

    return () => {
      supabase.removeChannel(channel);
    };
  }, [queryClient]);

  if (isPending) return <LoadingScreen />;
  if (error) return <p>Error loading jobs</p>;

  return (
    <div>
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

export default JobsList;