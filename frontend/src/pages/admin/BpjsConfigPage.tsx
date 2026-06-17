import { useBpjsConfigs, useBpjsConfigMutations } from "@/features/bpjs-config/hooks/useBpjsConfig";
import { BpjsConfigCard } from "@/features/bpjs-config/components/BpjsConfigCard";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { BPJSComponentType, BPJSRateConfig } from "@/features/bpjs-config/types";
import { Loader2, HeartPulse } from "lucide-react";

const BPJS_TYPES: BPJSComponentType[] = ["KESEHATAN", "JHT", "JKK", "JKM", "JP"];

export default function BpjsConfigPage() {
  const { data: configs, isLoading } = useBpjsConfigs();
  const { createMutation, updateMutation, deleteMutation } = useBpjsConfigMutations();

  const findConfig = (type: BPJSComponentType): BPJSRateConfig | undefined =>
    configs?.find((c) => c.type === type);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">BPJS Config</h2>
          <p className="text-slate-500">
            Manage BPJS Kesehatan and Ketenagakerjaan contribution rates per component.
          </p>
        </div>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg font-semibold">Rate Configuration</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : !configs || configs.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
              <HeartPulse className="h-12 w-12 mb-2" />
              <p>No BPJS config found</p>
              <p className="text-sm">Default rates will be used until configured.</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              {BPJS_TYPES.map((type) => (
                <BpjsConfigCard
                  key={type}
                  type={type}
                  config={findConfig(type)}
                  createMutation={createMutation}
                  updateMutation={updateMutation}
                  deleteMutation={deleteMutation}
                />
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
