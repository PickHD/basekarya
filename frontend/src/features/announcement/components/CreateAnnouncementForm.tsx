"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";

import { usePublishAnnouncement } from "@/features/announcement/hooks/useAnnouncement";

const formSchema = z.object({
  title: z.string().min(1, "Title is required"),
  body: z.string().min(1, "Body is required"),
});

type FormValues = z.infer<typeof formSchema>;

export function CreateAnnouncementForm() {
  const { mutate: publish, isPending } = usePublishAnnouncement();

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      title: "",
      body: "",
    },
  });

  const onSubmit = (data: FormValues) => {
    publish(data, {
      onSuccess: () => {
        form.reset();
      },
    });
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="font-semibold text-slate-900">Announcement Title</FormLabel>
              <FormControl>
                <Input placeholder="Enter announcement title..." {...field} className="border-slate-300 focus-visible:ring-blue-600" />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="body"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="font-semibold text-slate-900">Announcement Body</FormLabel>
              <FormControl>
                <Textarea 
                  placeholder="Enter announcement body..." 
                  className="min-h-[200px] resize-y border-slate-300 focus-visible:ring-blue-600 leading-relaxed" 
                  {...field} 
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <div className="flex justify-end pt-2">
          <Button type="submit" disabled={isPending} className="w-full sm:w-auto bg-blue-700 hover:bg-blue-800 text-white shadow-sm transition-all duration-200">
            {isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Publishing...
              </>
            ) : (
              "Publish Announcement"
            )}
          </Button>
        </div>
      </form>
    </Form>
  );
}
