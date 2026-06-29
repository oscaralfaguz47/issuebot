import { supabase } from "@/lib/supabase";
import { API_URL } from "@/lib/config"
import type { Membership } from "../types";


export async function getMemberships(): Promise<Membership[]> {
  const { data } = await supabase.auth.getSession();
  const token = data.session?.access_token;

  const res = await fetch(`${API_URL}/me`, {
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!res.ok) throw new Error(`API error: ${res.status}`);
  return res.json();
}