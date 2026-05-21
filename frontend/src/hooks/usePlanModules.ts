import { useMemo } from "react";
import { useCompanyProfile } from "@/features/company/hooks/useCompany";

export function usePlanModules() {
  const { data: company } = useCompanyProfile();

  const modules = useMemo(() => {
    if (!company?.plan_modules) return null;
    try {
      const parsed = JSON.parse(company.plan_modules);
      return parsed.modules as string[];
    } catch {
      return null;
    }
  }, [company?.plan_modules]);

  const hasModule = (moduleName: string) => {
    if (!modules) return true;
    return modules.includes(moduleName);
  };

  return { modules, hasModule, company };
}
