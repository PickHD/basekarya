import { OnboardingTemplateManager } from "@/features/onboarding/components/OnboardingTemplateManager";

export default function OnboardingTemplatePage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-xl font-bold text-slate-800">Onboarding Templates</h1>
        <p className="text-sm text-slate-500 mt-1">
          Configure checklist templates that are automatically applied when a new hire starts onboarding.
        </p>
      </div>

      <OnboardingTemplateManager />
    </div>
  );
}
