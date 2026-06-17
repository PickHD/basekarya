import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import type {
  BPJSRateConfig,
  BPJSRateConfigPayload,
  BPJSComponentType,
} from "../types";
import { BPJS_COMPONENT_LABELS, RISK_LEVEL_LABELS } from "../types";
import type { UseMutationResult } from "@tanstack/react-query";

const formSchema = z.object({
  employee_rate: z.coerce.number().min(0).max(100),
  employer_rate: z.coerce.number().min(0).max(100),
  max_salary_cap: z.coerce.number().min(0).optional().nullable(),
  industry_risk_level: z.string().optional().nullable(),
  effective_from: z.string().min(1, "Effective date is required"),
});

type FormValues = z.infer<typeof formSchema>;

interface BpjsConfigCardProps {
  type: BPJSComponentType;
  config?: BPJSRateConfig;
  createMutation: UseMutationResult<any, Error, BPJSRateConfigPayload>;
  updateMutation: UseMutationResult<any, Error, BPJSRateConfigPayload & { id: number }>;
  deleteMutation: UseMutationResult<any, Error, number>;
}

export function BpjsConfigCard({
  type,
  config,
  createMutation,
  updateMutation,
  deleteMutation,
}: BpjsConfigCardProps) {
  const [isEditing, setIsEditing] = useState(false);
  const isJKK = type === "JKK";

  const displayConfig = config?.is_active ? config : undefined;
  const isCompanyOverride = displayConfig && displayConfig.company_id !== null;

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      employee_rate: displayConfig?.employee_rate ? displayConfig.employee_rate * 100 : 0,
      employer_rate: displayConfig?.employer_rate ? displayConfig.employer_rate * 100 : 0,
      max_salary_cap: displayConfig?.max_salary_cap ?? null,
      industry_risk_level: displayConfig?.industry_risk_level ?? null,
      effective_from: displayConfig?.effective_from
        ? displayConfig.effective_from.split("T")[0]
        : new Date().toISOString().split("T")[0],
    },
  });

  const onSubmit = async (values: FormValues) => {
    const payload: BPJSRateConfigPayload = {
      type,
      employee_rate: values.employee_rate / 100,
      employer_rate: values.employer_rate / 100,
      max_salary_cap: values.max_salary_cap ?? null,
      industry_risk_level: isJKK ? (values.industry_risk_level as any) ?? null : null,
      effective_from: values.effective_from,
      is_active: true,
    };

    if (isCompanyOverride) {
      await updateMutation.mutateAsync({ ...payload, id: displayConfig!.id });
    } else {
      await createMutation.mutateAsync(payload);
    }
    setIsEditing(false);
  };

  const handleDelete = async () => {
    if (isCompanyOverride && displayConfig) {
      await deleteMutation.mutateAsync(displayConfig.id);
    }
  };

  const isPending = createMutation.isPending || updateMutation.isPending || deleteMutation.isPending;

  return (
    <Card>
      <CardHeader className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-2">
        <div>
          <CardTitle className="text-lg">{BPJS_COMPONENT_LABELS[type]}</CardTitle>
          <CardDescription>
            {displayConfig ? (
              isCompanyOverride ? (
                <Badge variant="secondary">Configured</Badge>
              ) : (
                <Badge variant="outline">Default</Badge>
              )
            ) : (
              <Badge variant="outline">Not Set</Badge>
            )}
          </CardDescription>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setIsEditing(!isEditing)}
        >
          {isEditing ? "Cancel" : isCompanyOverride ? "Edit" : "Configure"}
        </Button>
      </CardHeader>
      <CardContent>
        {!isEditing && displayConfig ? (
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Employee Rate</span>
              <span className="font-medium">{(displayConfig.employee_rate * 100).toFixed(2)}%</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Employer Rate</span>
              <span className="font-medium">{(displayConfig.employer_rate * 100).toFixed(2)}%</span>
            </div>
            {displayConfig.max_salary_cap && (
              <div className="flex justify-between">
                <span className="text-muted-foreground">Max Salary Cap</span>
                <span className="font-medium">
                  Rp {displayConfig.max_salary_cap.toLocaleString("id-ID")}
                </span>
              </div>
            )}
            {isJKK && displayConfig.industry_risk_level && (
              <div className="flex justify-between">
                <span className="text-muted-foreground">Risk Level</span>
                <span className="font-medium">
                  {RISK_LEVEL_LABELS[displayConfig.industry_risk_level as keyof typeof RISK_LEVEL_LABELS]}
                </span>
              </div>
            )}
            <div className="flex justify-between">
              <span className="text-muted-foreground">Effective From</span>
              <span className="font-medium">
                {new Date(displayConfig.effective_from).toLocaleDateString("id-ID")}
              </span>
            </div>
            {isCompanyOverride && (
              <Button
                variant="destructive"
                size="sm"
                className="mt-4 w-full"
                onClick={handleDelete}
                disabled={isPending}
              >
                Reset to Default
              </Button>
            )}
          </div>
        ) : isEditing ? (
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="employee_rate"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Employee Rate (%)</FormLabel>
                    <FormControl>
                      <Input type="number" step="0.01" min="0" max="100" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="employer_rate"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Employer Rate (%)</FormLabel>
                    <FormControl>
                      <Input type="number" step="0.01" min="0" max="100" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {(type === "KESEHATAN" || type === "JP") && (
                <FormField
                  control={form.control}
                  name="max_salary_cap"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Max Salary Cap (Rp)</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          step="1"
                          min="0"
                          value={field.value ?? ""}
                          onChange={(e) => field.onChange(e.target.value ? Number(e.target.value) : null)}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}

              {isJKK && (
                <FormField
                  control={form.control}
                  name="industry_risk_level"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Industry Risk Level</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        value={field.value ?? ""}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Select risk level" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {Object.entries(RISK_LEVEL_LABELS).map(([key, label]) => (
                            <SelectItem key={key} value={key}>
                              {label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}

              <FormField
                control={form.control}
                name="effective_from"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Effective From</FormLabel>
                    <FormControl>
                      <Input type="date" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" disabled={isPending} className="w-full">
                {isPending ? "Saving..." : isCompanyOverride ? "Update Override" : "Save Override"}
              </Button>
            </form>
          </Form>
        ) : (
          <p className="text-sm text-muted-foreground">No config found. Click Configure to create one.</p>
        )}
      </CardContent>
    </Card>
  );
}
