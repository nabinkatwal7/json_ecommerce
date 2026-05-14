"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { customerApi, publicApi, setToken } from "@/lib/api";
import type { User } from "@/lib/types";

type AuthContextValue = {
  user: User | null;
  loading: boolean;
  refresh: () => Promise<void>;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  setUserFromAuth: (token: string, user: User) => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    try {
      const u = await customerApi.me();
      setUser(u);
    } catch {
      setUser(null);
      setToken(null);
    }
  }, []);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      if (typeof window === "undefined" || !localStorage.getItem("token")) {
        if (!cancelled) setLoading(false);
        return;
      }
      try {
        const u = await customerApi.me();
        if (!cancelled) setUser(u);
      } catch {
        if (!cancelled) {
          setUser(null);
          setToken(null);
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const res = await publicApi.login({ email, password });
    setToken(res.token);
    setUser(res.user);
  }, []);

  const register = useCallback(
    async (name: string, email: string, password: string) => {
      const res = await publicApi.register({ name, email, password });
      setToken(res.token);
      setUser(res.user);
    },
    [],
  );

  const logout = useCallback(() => {
    setToken(null);
    setUser(null);
  }, []);

  const setUserFromAuth = useCallback((token: string, u: User) => {
    setToken(token);
    setUser(u);
  }, []);

  const value = useMemo(
    () => ({
      user,
      loading,
      refresh,
      login,
      register,
      logout,
      setUserFromAuth,
    }),
    [user, loading, refresh, login, register, logout, setUserFromAuth],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
