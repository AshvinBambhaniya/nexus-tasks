"use client";

import Link from "next/link";
import { useUser } from "@/hooks/use-user";
import { Button } from "@/components/ui/button";
import { CheckSquare, Layout, Users, Zap, Shield, ArrowRight, Github } from "lucide-react";

export default function Home() {
  const { user, isLoading } = useUser();

  return (
    <div className="flex min-h-screen flex-col bg-white">
      {/* Navigation */}
      <header className="sticky top-0 z-50 w-full border-b border-gray-100 bg-white/80 backdrop-blur-md">
        <div className="container mx-auto flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600 text-white">
              <CheckSquare className="h-5 w-5" />
            </div>
            <span className="text-xl font-bold text-gray-900">Nexus Tasks</span>
          </div>
          
          <nav className="hidden md:flex items-center gap-6 text-sm font-medium text-gray-600">
            <Link href="#features" className="hover:text-blue-600 transition-colors">Features</Link>
            <Link href="#about" className="hover:text-blue-600 transition-colors">About</Link>
            <Link href="https://github.com/AshvinBambhaniya/nexus-tasks" target="_blank" className="hover:text-blue-600 transition-colors">GitHub</Link>
          </nav>

          <div className="flex items-center gap-4">
            {isLoading ? (
              <Button disabled variant="ghost" size="sm">Loading...</Button>
            ) : user ? (
              <Link href="/dashboard">
                <Button>Go to Dashboard</Button>
              </Link>
            ) : (
              <>
                <Link href="/login" className="hidden sm:block">
                  <Button variant="ghost" className="text-gray-700 hover:text-gray-900">Log in</Button>
                </Link>
                <Link href="/register">
                  <Button>Get Started</Button>
                </Link>
              </>
            )}
          </div>
        </div>
      </header>

      <main className="flex-1">
        {/* Hero Section */}
        <section className="relative overflow-hidden pt-16 pb-24 lg:pt-32 lg:pb-40">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8 relative z-10 text-center">
            <div className="mx-auto max-w-3xl">
              <div className="inline-flex items-center rounded-full border border-blue-100 bg-blue-50 px-3 py-1 text-sm font-medium text-blue-600 mb-8">
                <span className="flex h-2 w-2 rounded-full bg-blue-600 mr-2 animate-pulse"></span>
                v1.0 is now live
              </div>
              <h1 className="text-4xl font-extrabold tracking-tight text-gray-900 sm:text-6xl mb-6">
                Bridging Personal Focus <br className="hidden sm:block" />
                <span className="text-blue-600">and Team Collaboration</span>
              </h1>
              <p className="text-lg text-gray-600 mb-10 max-w-2xl mx-auto leading-relaxed">
                Nexus Tasks is the unified workspace where individual productivity meets team synergy. 
                Manage your personal tasks and collaborate on team projects in one seamless interface.
              </p>
              <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
                {user ? (
                   <Link href="/dashboard">
                    <Button size="lg" className="h-12 px-8 text-base">
                      Launch Workspace <ArrowRight className="ml-2 h-4 w-4" />
                    </Button>
                  </Link>
                ) : (
                  <Link href="/register">
                    <Button size="lg" className="h-12 px-8 text-base">
                      Start for Free <ArrowRight className="ml-2 h-4 w-4" />
                    </Button>
                  </Link>
                )}
                <Link href="https://github.com/AshvinBambhaniya/nexus-tasks" target="_blank">
                  <Button variant="outline" size="lg" className="h-12 px-8 text-base bg-white text-slate-700 hover:text-slate-900 border-slate-200">
                    <Github className="mr-2 h-4 w-4" /> Star on GitHub
                  </Button>
                </Link>
              </div>
            </div>
          </div>
          
          {/* Decorative Background Elements */}
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-full -z-10 pointer-events-none">
            <div className="absolute top-0 left-1/4 w-96 h-96 bg-blue-100 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob"></div>
            <div className="absolute top-0 right-1/4 w-96 h-96 bg-purple-100 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob animation-delay-2000"></div>
            <div className="absolute -bottom-32 left-1/3 w-96 h-96 bg-pink-100 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob animation-delay-4000"></div>
          </div>
        </section>

        {/* Features Section */}
        <section id="features" className="py-24 bg-gray-50">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center max-w-2xl mx-auto mb-16">
              <h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl mb-4">
                Everything you need to ship faster
              </h2>
              <p className="text-lg text-gray-600">
                Designed for modern software teams who want to maintain flow state while staying synced.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <FeatureCard 
                icon={Layout}
                title="Hybrid Workspaces"
                description="Seamlessly switch between your private 'Personal' workspace and shared 'Team' environments without losing context."
              />
              <FeatureCard 
                icon={Zap}
                title="Real-time Sync"
                description="Powered by WebSockets, see task updates, comments, and project changes instantly as they happen."
              />
              <FeatureCard 
                icon={Shield}
                title="Role-Based Security"
                description="Granular permissions ensure that admins retain control while members have the autonomy they need to contribute."
              />
            </div>
          </div>
        </section>

        {/* Bottom CTA */}
        <section className="py-24 bg-white border-t border-gray-100">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8 text-center">
            <h2 className="text-3xl font-bold tracking-tight text-gray-900 mb-6">
              Ready to streamline your workflow?
            </h2>
            <p className="text-lg text-gray-600 mb-10 max-w-xl mx-auto">
              Join thousands of developers who are regaining their focus with Nexus Tasks.
            </p>
            <Link href="/register">
              <Button size="lg" className="h-12 px-8">
                Create your account
              </Button>
            </Link>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="bg-gray-50 border-t border-gray-200 py-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 flex flex-col md:flex-row justify-between items-center gap-6">
          <div className="flex items-center gap-2">
            <div className="flex h-6 w-6 items-center justify-center rounded bg-gray-900 text-white">
              <CheckSquare className="h-3 w-3" />
            </div>
            <span className="text-sm font-semibold text-gray-900">Nexus Tasks</span>
          </div>
          <p className="text-sm text-gray-500">
            © 2026 Nexus Tasks. Open Source Project.
          </p>
          <div className="flex gap-4">
            <Link href="#" className="text-gray-400 hover:text-gray-900 transition-colors">
              <span className="sr-only">GitHub</span>
              <Github className="h-5 w-5" />
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
}

function FeatureCard({ icon: Icon, title, description }: { icon: any, title: string, description: string }) {
  return (
    <div className="flex flex-col items-center text-center p-8 bg-white rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
      <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600 mb-6">
        <Icon className="h-6 w-6" />
      </div>
      <h3 className="text-xl font-semibold text-gray-900 mb-3">{title}</h3>
      <p className="text-gray-500 leading-relaxed">
        {description}
      </p>
    </div>
  )
}