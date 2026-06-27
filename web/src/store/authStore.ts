import { create } from "zustand";
import type { Session } from "@supabase/supabase-js";
import type { Membership } from "../types";

interface AuthState {
  session: Session | null;
  loading: boolean;
  memberships: Membership[];
  setSession: (session: Session | null) => void;
  setLoading: (loading: boolean) => void;
   setMemberships: (memberships: Membership[]) => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  session: null,
  loading: true,
  memberships: [],
  setSession: (session) => set({ session }),
  setLoading: (loading) => set({ loading }),
  setMemberships: (memberships) => set({ memberships }),
}));