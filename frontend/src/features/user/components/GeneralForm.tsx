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
  full_name: z
    .string()
    .min(3, "Fullname minimum at least 3 characters")
    .max(100, "Fullname too long")
    .transform((val) => val.trim()),

  phone_number: z
    .string()
    .min(10, "Phone number invalid")
    .max(15, "Phone number too long")
    .regex(/^[0-9]+$/, "Phone number only numbers can be included"),

  bank_name: z.string().min(3, "Bank name minimum at least 3 characters"),

  bank_account_number: z
    .string()
    .min(5, "Bank account number minimum at least 5 characters")
    .max(20, "Bank account number too long")
    .regex(/^\d+$/, "Bank account number only numbers"),

  bank_account_holder: z
    .string()
    .min(3, "Invalid bank account holder")
    .transform((val) => val.toUpperCase().trim()),

  npwp: z
    .string()
    .refine((val) => /^\d{15,16}$/.test(val.replace(/[.-]/g, "")), {
      message: "NPWP at least 15-16 number digit",
    })
    .transform((val) => val.replace(/[.-]/g, "")),
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
    defaultValues: {
      full_name: user.full_name || "",
      phone_number: user.phone_number || "",
      bank_name: user.bank_name || "",
      bank_account_number: user.bank_account_number || "",
      bank_account_holder: user.bank_account_holder || "",
      npwp: user.npwp || "",
    },
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
    formData.append("full_name", values.full_name);
    formData.append("phone_number", values.phone_number);
    formData.append("bank_name", values.bank_name);
    formData.append("bank_account_number", values.bank_account_number);
    formData.append("bank_account_holder", values.bank_account_holder);
    formData.append("npwp", values.npwp);

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
          name="full_name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Fullname</FormLabel>
              <FormControl>
                <Input placeholder="Your name..." {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

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

        <FormField
          control={form.control}
          name="bank_name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Bank Name</FormLabel>
              <FormControl>
                <Input placeholder="BCA.." {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="bank_account_number"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Bank Account Number</FormLabel>
              <FormControl>
                <Input placeholder="79123..." {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="bank_account_holder"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Bank Account Holder</FormLabel>
              <FormControl>
                <Input placeholder="SENDY.." {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="npwp"
          render={({ field }) => (
            <FormItem>
              <FormLabel>NPWP Number</FormLabel>
              <FormControl>
                <Input placeholder="319287391.." {...field} />
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
