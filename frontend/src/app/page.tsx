"use client";

import Link from "next/link";
import { useUser } from "@/hooks/use-user";
import { Button } from "@/components/ui/button";
import {
  CheckSquare,
  Layout,
  Users,
  Zap,
  Shield,
  ArrowRight,
  Github,
  Code2,
  Terminal,
  Cpu,
  GitBranch,
} from "lucide-react";

export default function Home() {
  const { user, isLoading } = useUser();

  return (
    <div className="flex min-h-screen flex-col bg-white selection:bg-blue-100">
      {/* Navigation */}
      <header className="sticky top-0 z-50 w-full border-b border-gray-100 bg-white/80 backdrop-blur-md">
        <div className="container mx-auto flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600 text-white shadow-md shadow-blue-200">
              <CheckSquare className="h-5 w-5" />
            </div>
            <span className="text-xl font-bold tracking-tight text-gray-900">
              Nexus Tasks
            </span>
          </div>

          <nav className="hidden items-center gap-8 text-sm font-medium text-gray-600 md:flex">
            <Link
              href="#features"
              className="transition-colors hover:text-blue-600"
            >
              Features
            </Link>
            <Link
              href="#tech-stack"
              className="transition-colors hover:text-blue-600"
            >
              Tech Stack
            </Link>
            <Link
              href="https://github.com/AshvinBambhaniya/nexus-tasks"
              target="_blank"
              className="transition-colors hover:text-blue-600"
            >
              GitHub
            </Link>
          </nav>

          <div className="flex items-center gap-4">
            {isLoading ? (
              <Button disabled variant="ghost" size="sm">
                Loading...
              </Button>
            ) : user ? (
              <Link href="/dashboard">
                <Button className="font-semibold shadow-md shadow-blue-100">
                  Go to Dashboard
                </Button>
              </Link>
            ) : (
              <>
                <Link href="/login" className="hidden sm:block">
                  <Button
                    variant="ghost"
                    className="font-medium text-gray-600 hover:text-gray-900"
                  >
                    Log in
                  </Button>
                </Link>
                <Link href="/register">
                  <Button className="font-semibold shadow-md shadow-blue-100">
                    Get Started
                  </Button>
                </Link>
              </>
            )}
          </div>
        </div>
      </header>

      <main className="flex-1">
        {/* Hero Section */}
        <section className="relative overflow-hidden pt-20 pb-28 lg:pt-32 lg:pb-40">
          <div className="relative z-10 container mx-auto px-4 text-center sm:px-6 lg:px-8">
            <div className="mx-auto max-w-4xl">
              <div className="mb-8 inline-flex items-center rounded-full border border-blue-100 bg-blue-50/50 px-3 py-1 text-sm font-medium text-blue-600 backdrop-blur-sm">
                <span className="mr-2 flex h-2 w-2 animate-pulse rounded-full bg-blue-600"></span>
                v1.0 is Live • Open Source
              </div>
              <h1 className="mb-6 text-4xl font-extrabold tracking-tight text-gray-900 sm:text-6xl">
                The Task Manager for <br className="hidden sm:block" />
                <span className="text-blue-600">Developers in Flow</span>
              </h1>
              <p className="mx-auto mb-10 max-w-2xl text-lg leading-relaxed text-gray-600">
                Stop fighting complex workflows. Nexus Tasks bridges your{" "}
                <strong className="font-semibold text-gray-900">
                  personal focus
                </strong>{" "}
                with{" "}
                <strong className="font-semibold text-gray-900">
                  team collaboration
                </strong>{" "}
                in a clean, high-performance interface built for speed.
              </p>

              <div className="mb-16 flex flex-col items-center justify-center gap-4 sm:flex-row">
                {user ? (
                  <Link href="/dashboard">
                    <Button size="lg" className="h-12 px-8 text-base">
                      Launch Workspace <ArrowRight className="ml-2 h-4 w-4" />
                    </Button>
                  </Link>
                ) : (
                  <Link href="/register">
                    <Button size="lg" className="h-12 px-8 text-base">
                      Start Building Free{" "}
                      <ArrowRight className="ml-2 h-4 w-4" />
                    </Button>
                  </Link>
                )}
                <Link
                  href="https://github.com/AshvinBambhaniya/nexus-tasks"
                  target="_blank"
                >
                  <Button
                    variant="outline"
                    size="lg"
                    className="h-12 border-slate-200 bg-white px-8 text-base text-slate-700 hover:text-slate-900"
                  >
                    <Github className="mr-2 h-4 w-4" /> Star on GitHub
                  </Button>
                </Link>
              </div>

              {/* Tech Stack Pills */}
              <div
                className="flex flex-wrap justify-center gap-3 opacity-80"
                id="tech-stack"
              >
                <TechBadge icon={Code2} label="Next.js 16" />
                <TechBadge icon={Terminal} label="FastAPI" />
                <TechBadge icon={GitBranch} label="TypeScript" />
                <TechBadge icon={Cpu} label="Real-time WebSockets" />
              </div>
            </div>
          </div>

          {/* Decorative Background Elements */}
          <div className="pointer-events-none absolute top-0 left-1/2 -z-10 h-full w-full -translate-x-1/2 overflow-hidden">
            <div className="animate-blob absolute top-[-10%] left-[10%] h-[500px] w-[500px] rounded-full bg-blue-100/50 opacity-60 mix-blend-multiply blur-[100px] filter"></div>
            <div className="animate-blob animation-delay-2000 absolute top-[-10%] right-[10%] h-[500px] w-[500px] rounded-full bg-indigo-100/50 opacity-60 mix-blend-multiply blur-[100px] filter"></div>
            <div className="animate-blob animation-delay-4000 absolute bottom-[-20%] left-[30%] h-[600px] w-[600px] rounded-full bg-sky-100/50 opacity-60 mix-blend-multiply blur-[100px] filter"></div>
          </div>
        </section>

        {/* Features Grid */}
        <section
          id="features"
          className="border-t border-gray-100 bg-gray-50/50 py-24"
        >
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="mx-auto mb-16 max-w-3xl text-center">
              <h2 className="mb-4 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
                Designed for the way you code
              </h2>
              <p className="text-lg text-gray-600">
                Nexus Tasks cuts through the enterprise bloat to give you
                exactly what you need: clarity, speed, and structure.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
              <FeatureCard
                icon={Layout}
                title="Hybrid Workspaces"
                description="Keep your personal side-projects and team sprints in one place. Switch contexts instantly without logging out."
              />
              <FeatureCard
                icon={Code2}
                title="Markdown Native"
                description="Write task descriptions and comments using standard Markdown. Code blocks, lists, and formatting just work."
              />
              <FeatureCard
                icon={Zap}
                title="Real-time Sync"
                description="Collaborate without refreshing. See comments, status changes, and new tasks appear instantly via WebSockets."
              />
              <FeatureCard
                icon={GitBranch}
                title="GitHub-style Workflow"
                description="Manage tasks with a familiar flow: Open, In Progress, Done. Track history and discussions linearly."
              />
              <FeatureCard
                icon={Users}
                title="Team Alignment"
                description="Organize members into Teams. Assign projects to specific groups and manage access with minimal friction."
              />
              <FeatureCard
                icon={Shield}
                title="Production Ready"
                description="Built with a robust FastAPI backend and Type-Safe Next.js frontend. Secure, scalable, and self-hostable."
              />
            </div>
          </div>
        </section>

        {/* Call to Action */}
        <section className="bg-white py-24">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="relative overflow-hidden rounded-3xl bg-gray-900 px-6 py-16 text-center shadow-2xl sm:px-12 sm:py-20">
              <div className="relative z-10 mx-auto max-w-2xl">
                <h2 className="mb-6 text-3xl font-bold tracking-tight text-white sm:text-4xl">
                  Ready to regain your focus?
                </h2>
                <p className="mb-10 text-lg text-gray-300">
                  Join the community of developers building better software with
                  less friction. Open source and free to get started.
                </p>
                <Link href="/register">
                  <Button
                    size="lg"
                    className="h-14 border-none bg-white px-8 text-base font-semibold text-gray-900 hover:bg-gray-100 hover:text-gray-900"
                  >
                    Get Started for Free
                  </Button>
                </Link>
              </div>

              {/* Abstract shapes */}
              <svg
                viewBox="0 0 1024 1024"
                className="absolute top-1/2 left-1/2 -z-10 h-[64rem] w-[64rem] -translate-x-1/2 [mask-image:radial-gradient(closest-side,white,transparent)]"
                aria-hidden="true"
              >
                <circle
                  cx="512"
                  cy="512"
                  r="512"
                  fill="url(#gradient)"
                  fillOpacity="0.7"
                />
                <defs>
                  <radialGradient id="gradient">
                    <stop stopColor="#4F46E5" />
                    <stop offset="1" stopColor="#80CAFF" />
                  </radialGradient>
                </defs>
              </svg>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="border-t border-gray-100 bg-white py-12">
        <div className="container mx-auto flex flex-col items-center justify-between gap-6 px-4 sm:px-6 md:flex-row lg:px-8">
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-900 text-white">
              <CheckSquare className="h-4 w-4" />
            </div>
            <div>
              <span className="block text-sm font-bold text-gray-900">
                Nexus Tasks
              </span>
              <span className="text-xs text-gray-500">Built for Builders</span>
            </div>
          </div>
          <p className="text-sm text-gray-500">
            © {new Date().getFullYear()} Nexus Tasks. Distributed under MIT
            License.
          </p>
          <div className="flex gap-6">
            <Link
              href="https://github.com/AshvinBambhaniya/nexus-tasks"
              target="_blank"
              className="flex items-center gap-2 text-sm text-gray-400 transition-colors hover:text-gray-900"
            >
              <Github className="h-5 w-5" />
              <span>Source Code</span>
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
}

function FeatureCard({
  icon: Icon,
  title,
  description,
}: {
  icon: React.ElementType;
  title: string;
  description: string;
}) {
  return (
    <div className="group flex flex-col items-start rounded-2xl border border-gray-100 bg-white p-8 text-left shadow-sm transition-all duration-200 hover:border-blue-100 hover:shadow-md">
      <div className="mb-6 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600 transition-colors group-hover:bg-blue-600 group-hover:text-white">
        <Icon className="h-6 w-6" />
      </div>
      <h3 className="mb-3 text-xl font-bold text-gray-900">{title}</h3>
      <p className="leading-relaxed text-gray-500">{description}</p>
    </div>
  );
}

function TechBadge({
  icon: Icon,
  label,
}: {
  icon: React.ElementType;
  label: string;
}) {
  return (
    <div className="flex items-center gap-2 rounded-full border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-600 shadow-sm">
      <Icon className="h-4 w-4 text-gray-400" />
      {label}
    </div>
  );
}
