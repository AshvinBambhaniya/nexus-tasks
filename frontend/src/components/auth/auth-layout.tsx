import Link from "next/link";
import { CheckSquare } from "lucide-react";

interface AuthLayoutProps {
  children: React.ReactNode;
  title: string;
  description: string;
  linkText: string;
  linkHref: string;
}

export function AuthLayout({ children, title, description, linkText, linkHref }: AuthLayoutProps) {
  return (
    <div className="container relative min-h-screen flex-col items-center justify-center grid lg:max-w-none lg:grid-cols-2 lg:px-0">
      <div className="relative hidden h-full flex-col bg-muted p-10 text-white dark:border-r lg:flex">
        <div className="absolute inset-0 bg-zinc-900" />
        <div className="relative z-20 flex items-center text-lg font-medium">
          <CheckSquare className="mr-2 h-6 w-6" />
          Nexus Tasks
        </div>
        <div className="relative z-20 mt-auto">
          <blockquote className="space-y-2">
            <p className="text-lg">
              &ldquo;Nexus Tasks bridges the gap between my personal focus and 
              team collaboration, making it the ultimate tool for engineers 
              who need to stay organized.&rdquo;
            </p>
            <footer className="text-sm">Ashvin Bambhaniya, Software Engineer</footer>
          </blockquote>
        </div>
      </div>
      <div className="lg:p-8 bg-gray-50 h-full flex items-center justify-center">
        <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
          <div className="flex flex-col space-y-2 text-center">
            <h1 className="text-2xl font-semibold tracking-tight text-gray-900">{title}</h1>
            <p className="text-sm text-muted-foreground text-gray-500">
              {description}
            </p>
          </div>
          {children}
          <p className="px-8 text-center text-sm text-muted-foreground text-gray-500">
            {linkText === "Create an account" ? "Don't have an account? " : "Already have an account? "}
            <Link
              href={linkHref}
              className="underline underline-offset-4 hover:text-primary text-blue-600"
            >
              {linkText}
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
