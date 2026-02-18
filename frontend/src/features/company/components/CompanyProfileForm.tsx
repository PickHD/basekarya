import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2, Upload, Building2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { useUpdateCompanyProfile } from "../hooks/useCompany";
import type { CompanyProfile } from "../types";

const companySchema = z.object({
  name: z
    .string()
    .min(3, "Company name must be at least 3 characters")
    .max(100, "Company name is too long")
    .trim(),

  email: z.string().email("Invalid email address").min(5, "Email is required"),

  phone_number: z
    .string()
    .min(6, "Phone number is too short")
    .max(20, "Phone number is too long")
    .regex(/^\d+$/, "Phone number must contain only numbers"),

  website: z.string().optional().or(z.literal("")), // Allow empty string

  tax_number: z
    .string()
    .min(10, "Tax number (NPWP) usually has at least 10 digits")
    .regex(/^\d+$/, "Tax number must contain only numbers")
    .optional()
    .or(z.literal("")),

  address: z
    .string()
    .min(5, "Address is too short")
    .max(500, "Address is too long"),
});

type CompanyFormValues = z.infer<typeof companySchema>;

interface CompanyProfileFormProps {
  initialData?: CompanyProfile;
}

export function CompanyProfileForm({ initialData }: CompanyProfileFormProps) {
  const { mutate: updateCompany, isPending } = useUpdateCompanyProfile();

  const [preview, setPreview] = useState<string | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const form = useForm<CompanyFormValues>({
    resolver: zodResolver(companySchema),
    defaultValues: {
      name: "",
      email: "",
      phone_number: "",
      website: "",
      tax_number: "",
      address: "",
    },
    values: initialData
      ? {
          name: initialData.name || "",
          email: initialData.email || "",
          phone_number: initialData.phone_number || "",
          website: initialData.website || "",
          tax_number: initialData.tax_number || "",
          address: initialData.address || "",
        }
      : undefined,
  });

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.size > 2 * 1024 * 1024) {
        alert("File size too big (max 2MB)");
        return;
      }
      setSelectedFile(file);
      const objectUrl = URL.createObjectURL(file);
      setPreview(objectUrl);
    }
  };

  const onSubmit = (values: CompanyFormValues) => {
    const formData = new FormData();

    formData.append("name", values.name);
    formData.append("email", values.email);
    formData.append("phone_number", values.phone_number);
    formData.append("address", values.address);
    if (values.website) formData.append("website", values.website);
    if (values.tax_number) formData.append("tax_number", values.tax_number);

    if (selectedFile) {
      formData.append("logo_url", selectedFile);
    }

    updateCompany(formData, {
      onSuccess: () => {
        if (preview) {
          URL.revokeObjectURL(preview);
        }
        setPreview(null);
        setSelectedFile(null);
      },
    });
  };

  const avatarSrc = preview || initialData?.logo_url || "";

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <div className="flex flex-col sm:flex-row items-center gap-6 p-6 border rounded-lg bg-white shadow-sm">
          <Avatar className="w-28 h-28 rounded-lg border-2 border-gray-100">
            <AvatarImage
              src={avatarSrc}
              className="object-contain bg-white"
              alt="Company Logo"
            />
            <AvatarFallback className="rounded-lg bg-slate-100 text-slate-400">
              <Building2 className="w-10 h-10" />
            </AvatarFallback>
          </Avatar>

          <div className="flex-1 space-y-2 text-center sm:text-left">
            <h3 className="text-lg font-medium">Company Logo</h3>
            <div className="flex items-center justify-center sm:justify-start gap-3 mt-2">
              <Input
                id="logo-upload"
                type="file"
                accept="image/*"
                className="hidden"
                onChange={handleFileChange}
              />
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => document.getElementById("logo-upload")?.click()}
              >
                <Upload className="w-4 h-4 mr-2" />
                Upload Logo
              </Button>
            </div>
          </div>
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem className="col-span-2 md:col-span-1">
                <FormLabel>
                  Company Name <span className="text-red-500">*</span>
                </FormLabel>
                <FormControl>
                  <Input placeholder="Acme Corp, Inc." {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem className="col-span-2 md:col-span-1">
                <FormLabel>
                  Official Email <span className="text-red-500">*</span>
                </FormLabel>
                <FormControl>
                  <Input placeholder="admin@acme.com" {...field} />
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
                <FormLabel>
                  Phone Number <span className="text-red-500">*</span>
                </FormLabel>
                <FormControl>
                  <Input placeholder="021-555-0199" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="website"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Website</FormLabel>
                <FormControl>
                  <Input placeholder="https://acme.com" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="tax_number"
            render={({ field }) => (
              <FormItem className="col-span-2">
                <FormLabel>Tax Number (NPWP)</FormLabel>
                <FormControl>
                  <Input placeholder="01.234.567.8-901.000" {...field} />
                </FormControl>
                <FormDescription>
                  Used for invoicing and payroll tax reports.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="address"
            render={({ field }) => (
              <FormItem className="col-span-2">
                <FormLabel>
                  Office Address <span className="text-red-500">*</span>
                </FormLabel>
                <FormControl>
                  <Textarea
                    placeholder="Jl. Sudirman No. 1..."
                    className="min-h-[100px]"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <div className="flex justify-end pt-4">
          <Button
            type="submit"
            disabled={isPending}
            className="w-full sm:w-auto min-w-[150px]"
          >
            {isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Saving Changes...
              </>
            ) : (
              "Save Company Profile"
            )}
          </Button>
        </div>
      </form>
    </Form>
  );
}
