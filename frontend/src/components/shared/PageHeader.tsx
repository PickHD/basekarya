import type { LucideIcon } from "lucide-react";

interface PageHeaderProps {
  title: string;
  description?: string;
  icon?: LucideIcon;
}

export function PageHeader({ title, description, icon: Icon }: PageHeaderProps) {
  return (
    <div className="flex flex-col gap-1 pb-4">
      <div className="flex items-center gap-2">
        {Icon && <Icon className="h-6 w-6 text-slate-700" />}
        <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-slate-900">
          {title}
        </h1>
      </div>
      {description && (
        <p className="text-sm sm:text-base text-slate-500 max-w-2xl">
          {description}
        </p>
      )}
    </div>
  );
}
