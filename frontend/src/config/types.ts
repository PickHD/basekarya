import { type LucideIcon } from "lucide-react";

export interface MenuItem {
  title: string;
  href: string;
  icon: LucideIcon;
  permission?: string | string[];
  group?: string;
  hideForPlatformAdmin?: boolean;
  platformAdminOnly?: boolean;
  requiredModule?: string;
}
