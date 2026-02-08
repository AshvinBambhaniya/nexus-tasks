"use client";

import Link from "next/link";
import { useUser } from "@/hooks/use-user";
import { Button } from "@/components/ui/button";

export default function Home() {
  const { user, isLoading } = useUser();

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-24">
      <h1 className="text-4xl font-bold mb-4">Nexus Tasks</h1>
      <p className="text-lg text-gray-600 mb-8">
        Bridging personal productivity and team collaboration.
      </p>

      <div className="flex gap-4">
        {isLoading ? (
          <Button disabled isLoading>
            Loading...
          </Button>
        ) : user ? (
          <Link href="/dashboard">
            <Button size="lg">Go to Dashboard</Button>
          </Link>
        ) : (
          <>
            <Link href="/login">
              <Button variant="outline" size="lg">
                Login
              </Button>
            </Link>
            <Link href="/register">
              <Button size="lg">Register</Button>
            </Link>
          </>
        )}
      </div>
    </div>
  );
}
