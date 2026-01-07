import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Label } from "@/components/ui/label";
import { useUpdateProfile } from "../hooks/useProfile";

const generalSchema = z.object({
  phone_number: z.string().min(10, "Phone number minimum at least 10 digit"),
});

interface GeneralFormProps {
  user: any;
}

export function GeneralForm({ user }: GeneralFormProps) {
  const { mutate: updateProfile, isPending } = useUpdateProfile();
  const [preview, setPreview] = useState<string | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const form = useForm({
    resolver: zodResolver(generalSchema),
    defaultValues: { phone_number: user.phone_number || "" },
  });

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      const objectUrl = URL.createObjectURL(file);
      setPreview(objectUrl);
    }
  };

  const onSubmit = (values: any) => {
    const formData = new FormData();
    formData.append("phone_number", values.phone_number);
    if (selectedFile) {
      formData.append("photo", selectedFile);
    }
    updateProfile(formData, {
      onSuccess: () => {
        if (preview) {
          URL.revokeObjectURL(preview);
        }
        setPreview(null);
        setSelectedFile(null);
      },
    });
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <div className="flex items-center gap-6">
          <div className="shrink-0">
            <Avatar className="w-20 h-20 border">
              <AvatarImage
                src={preview || user.profile_picture_url}
                className="object-cover"
              />
              <AvatarFallback>Pic</AvatarFallback>
            </Avatar>
          </div>
          <div className="flex-1">
            <Label>Profile Photo</Label>
            <Input
              type="file"
              accept="image/*"
              className="mt-2 cursor-pointer"
              onChange={handleFileChange}
            />
            <p className="text-xs text-slate-500 mt-1">JPG, PNG up to 2MB.</p>
          </div>
        </div>

        <FormField
          control={form.control}
          name="phone_number"
          render={({ field }) => (
            <FormItem>
              <FormLabel>WhatsApp Number</FormLabel>
              <FormControl>
                <Input placeholder="0812..." {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="flex justify-end">
          <Button
            type="submit"
            disabled={isPending}
            className="bg-blue-600 hover:bg-blue-700"
          >
            {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Save Changes
          </Button>
        </div>
      </form>
    </Form>
  );
}
