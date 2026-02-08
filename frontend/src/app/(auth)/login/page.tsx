"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/hooks/use-auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { CheckSquare } from "lucide-react";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const { login, isLoading, error } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await login(email, password);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md space-y-8 bg-white p-8 rounded-xl shadow-sm border border-gray-200">
        <div className="flex flex-col items-center text-center">
          <div className="flex items-center justify-center h-12 w-12 rounded-full bg-blue-100 mb-4">
             <CheckSquare className="h-6 w-6 text-blue-600" />
          </div>
          <h2 className="text-2xl font-bold tracking-tight text-gray-900">Welcome back</h2>
          <p className="text-sm text-gray-500 mt-2">
            Enter your email to sign in to your account
          </p>
        </div>

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                placeholder="name@example.com"
                type="email"
                autoComplete="email"
                disabled={isLoading}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="grid gap-2">
              <div className="flex items-center justify-between">
                <Label htmlFor="password">Password</Label>
              </div>
              <Input
                id="password"
                placeholder="••••••••"
                type="password"
                autoComplete="current-password"
                disabled={isLoading}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
          </div>

          {error && (
            <div className="text-sm font-medium text-red-600 bg-red-50 p-3 rounded-md border border-red-200">
              {error}
            </div>
          )}

          <Button disabled={isLoading} isLoading={isLoading} className="w-full">
            Sign In
          </Button>
        </form>

        <div className="text-center text-sm">
          <span className="text-gray-500">Don&apos;t have an account? </span>
          <Link href="/register" className="font-medium text-blue-600 hover:text-blue-500 hover:underline">
            Sign up
          </Link>
        </div>
      </div>
    </div>
  );
}