import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
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
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Loader2, Pencil, Percent } from "lucide-react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { usePtkpConfigs, usePtkpConfigMutations } from "../hooks/usePph21Config";
import { PTKP_CODES } from "../types";
import type { PTKPConfig } from "../types";

const formSchema = z.object({
  code: z.string().min(1, "Code is required"),
  annual_amount: z.coerce.number().min(0, "Amount must be positive"),
  effective_year: z.coerce.number().int().min(2000).max(2100),
});

type FormValues = z.infer<typeof formSchema>;

function PtkpEditDialog({
  open,
  onOpenChange,
  config,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  config?: PTKPConfig;
}) {
  const { updateMutation } = usePtkpConfigMutations();

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      code: config?.code ?? "",
      annual_amount: config?.annual_amount ?? 0,
      effective_year: config?.effective_year ?? new Date().getFullYear(),
    },
  });

  const onSubmit = async (values: FormValues) => {
    if (config) {
      await updateMutation.mutateAsync({ ...values, id: config.id });
    }
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit PTKP: {config?.code}</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="code"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>PTKP Code</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger disabled={updateMutation.isPending}>
                        <SelectValue placeholder="Select code" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {PTKP_CODES.map((code) => (
                        <SelectItem key={code} value={code}>
                          {code}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="annual_amount"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Annual Amount (Rp)</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      step="1"
                      min="0"
                      disabled={updateMutation.isPending}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="effective_year"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Effective Year</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      disabled={updateMutation.isPending}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={updateMutation.isPending}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={updateMutation.isPending}>
                {updateMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Save
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export function PtkpTable() {
  const { data: configs, isLoading } = usePtkpConfigs();
  const [editingConfig, setEditingConfig] = useState<PTKPConfig | null>(null);

  const sortedConfigs = (configs ?? []).sort((a, b) => a.code.localeCompare(b.code));

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (sortedConfigs.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
        <Percent className="h-12 w-12 mb-2" />
        <p>No PTKP configs found</p>
        <p className="text-sm">Seed data to populate PTKP thresholds.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <p className="text-sm text-muted-foreground">
        PTKP (Penghasilan Tidak Kena Pajak) annual thresholds used for PPh 21 calculations.
      </p>

      <div className="hidden md:block">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Code</TableHead>
              <TableHead>Annual Amount</TableHead>
              <TableHead>Effective Year</TableHead>
              <TableHead className="w-24 text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sortedConfigs.map((config) => (
              <TableRow key={config.id}>
                <TableCell className="font-medium">{config.code}</TableCell>
                <TableCell>Rp {config.annual_amount.toLocaleString("id-ID")}</TableCell>
                <TableCell>{config.effective_year}</TableCell>
                <TableCell className="text-right">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => setEditingConfig(config)}
                    className="h-8 w-8"
                  >
                    <Pencil className="h-4 w-4" />
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <div className="md:hidden space-y-3">
        {sortedConfigs.map((config) => (
          <Card key={config.id} className="p-4">
            <div className="flex items-start justify-between">
              <div>
                <p className="font-medium font-mono text-sm">{config.code}</p>
                <p className="text-base">Rp {config.annual_amount.toLocaleString("id-ID")}</p>
                <p className="text-sm text-muted-foreground">Year: {config.effective_year}</p>
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setEditingConfig(config)}
                className="h-8 w-8 flex-shrink-0"
              >
                <Pencil className="h-4 w-4" />
              </Button>
            </div>
          </Card>
        ))}
      </div>

      {editingConfig && (
        <PtkpEditDialog
          open={!!editingConfig}
          onOpenChange={(open) => {
            if (!open) setEditingConfig(null);
          }}
          config={editingConfig}
        />
      )}
    </div>
  );
}
