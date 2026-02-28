import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Loader2 } from "lucide-react";
import { useCreateOvertime } from "@/features/overtime/hooks/useOvertime";
import type {
  CreateOvertimePayload,
  OvertimeFormDialogProps,
} from "@/features/overtime/types";

const overtimeSchema = z.object({
  date: z.string().min(1, "Tanggal harus diisi"),
  start_time: z.string().regex(/^([01]\d|2[0-3]):?([0-5]\d)$/, "Format waktu: HH:MM"),
  end_time: z.string().regex(/^([01]\d|2[0-3]):?([0-5]\d)$/, "Format waktu: HH:MM"),
  reason: z.string().min(3, "Alasan harus diisi"),
});

export function OvertimeFormDialog({ open, onOpenChange }: OvertimeFormDialogProps) {
  const { mutate, isPending } = useCreateOvertime();

  const form = useForm<z.infer<typeof overtimeSchema>>({
    resolver: zodResolver(overtimeSchema),
    defaultValues: {
      date: "",
      start_time: "",
      end_time: "",
      reason: "",
    },
  });

  useEffect(() => {
    if (open) {
      form.reset({
        date: "",
        start_time: "",
        end_time: "",
        reason: "",
      });
    }
  }, [open, form]);

  const onSubmit = (data: z.infer<typeof overtimeSchema>) => {
    const payload: CreateOvertimePayload = {
      date: data.date,
      start_time: data.start_time,
      end_time: data.end_time,
      reason: data.reason,
    };

    mutate(payload, {
      onSuccess: () => {
        onOpenChange(false);
      },
      onError: (err) => {
        console.error(err);
      }
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Ajukan Lembur (Overtime)</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="date"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tanggal Lembur (YYYY-MM-DD)</FormLabel>
                  <FormControl>
                    <Input placeholder="Misal: 2026-03-01" type="date" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="start_time"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Waktu Mulai</FormLabel>
                    <FormControl>
                      <Input placeholder="HH:MM (Contoh: 18:00)" type="time" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="end_time"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Waktu Selesai</FormLabel>
                    <FormControl>
                      <Input placeholder="HH:MM (Contoh: 21:00)" type="time" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="reason"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Pekerjaan yang Dilakukan</FormLabel>
                  <FormControl>
                    <Textarea className="resize-none" rows={3} placeholder="Mengerjakan tugas A..." {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="submit"
                disabled={isPending}
                className="w-full sm:w-auto"
              >
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isPending ? "Mengirim..." : "Kirim Pengajuan"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
