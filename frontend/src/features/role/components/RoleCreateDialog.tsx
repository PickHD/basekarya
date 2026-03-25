import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Loader2 } from "lucide-react";
import { useCreateRole } from "@/features/role/hooks/useRole";

interface RoleCreateDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function RoleCreateDialog({ open, onOpenChange }: RoleCreateDialogProps) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const { mutate: createRole, isPending } = useCreateRole();

  const handleOpenChangeWrapper = (isOpen: boolean) => {
    onOpenChange(isOpen);
    if (!isOpen) {
      setTimeout(() => {
        setName("");
        setDescription("");
      }, 300);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    createRole(
      { name: name.trim().toUpperCase(), description: description.trim() },
      {
        onSuccess: () => {
          handleOpenChangeWrapper(false);
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChangeWrapper}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create New Role</DialogTitle>
          <DialogDescription>
            Add a new role to the system. You can assign permissions to it later.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="name">Role Name (e.g., HR_MANAGER)</Label>
            <Input
              id="name"
              placeholder="Enter role name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className="uppercase"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">Description <span className="text-slate-400 font-normal">(Optional)</span></Label>
            <Textarea
              id="description"
              placeholder="Briefly describe the purpose of this role"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
            />
          </div>
          <DialogFooter className="pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => handleOpenChangeWrapper(false)}
              disabled={isPending}
            >
              Cancel
            </Button>
            <Button type="submit" className="bg-blue-600 hover:bg-blue-700" disabled={isPending || !name.trim()}>
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Create Role
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
