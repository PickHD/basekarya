import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Loader2, Plus, Pencil, Trash2, X } from "lucide-react";
import { useTerBrackets, useTerBracketMutations } from "../hooks/usePph21Config";
import { TER_CATEGORIES } from "../types";
import type { TERBracket } from "../types";
import type { UseMutationResult } from "@tanstack/react-query";

interface EditableRowProps {
  bracket?: TERBracket;
  category: string;
  createMutation: UseMutationResult<any, Error, any>;
  updateMutation: UseMutationResult<any, Error, any>;
  deleteMutation: UseMutationResult<any, Error, number>;
  onCancel?: () => void;
}

function EditableRow({
  bracket,
  category,
  createMutation,
  updateMutation,
  deleteMutation,
  onCancel,
}: EditableRowProps) {
  const [bracketNumber, setBracketNumber] = useState(bracket?.bracket_number ?? 0);
  const [minSalary, setMinSalary] = useState(bracket?.min_monthly_salary ?? 0);
  const [rate, setRate] = useState(bracket ? bracket.rate * 100 : 0);
  const [effectiveFrom, setEffectiveFrom] = useState(
    bracket?.effective_from?.split("T")[0] ?? new Date().toISOString().split("T")[0]
  );

  const isPending = createMutation.isPending || updateMutation.isPending;

  const handleSave = async () => {
    const payload = {
      category,
      bracket_number: bracketNumber,
      min_monthly_salary: minSalary,
      rate: rate / 100,
      effective_from: effectiveFrom,
    };

    if (bracket) {
      await updateMutation.mutateAsync({ ...payload, id: bracket.id });
    } else {
      await createMutation.mutateAsync(payload);
    }
    onCancel?.();
  };

  return (
    <TableRow>
      <TableCell>
        <Input
          type="number"
          value={bracketNumber}
          onChange={(e) => setBracketNumber(Number(e.target.value))}
          className="w-20"
        />
      </TableCell>
      <TableCell>
        <Input
          type="number"
          value={minSalary}
          onChange={(e) => setMinSalary(Number(e.target.value))}
          className="w-36"
        />
      </TableCell>
      <TableCell>
        <Input
          type="number"
          step="0.01"
          value={rate}
          onChange={(e) => setRate(Number(e.target.value))}
          className="w-24"
        />
      </TableCell>
      <TableCell>
        <Input
          type="date"
          value={effectiveFrom}
          onChange={(e) => setEffectiveFrom(e.target.value)}
          className="w-40"
        />
      </TableCell>
      <TableCell>
        <div className="flex gap-1">
          <Button size="sm" onClick={handleSave} disabled={isPending}>
            {isPending && <Loader2 className="mr-1 h-3 w-3 animate-spin" />}
            Save
          </Button>
          <Button size="sm" variant="ghost" onClick={onCancel}>
            <X className="h-4 w-4" />
          </Button>
        </div>
      </TableCell>
    </TableRow>
  );
}

function StaticRow({
  bracket,
  onEdit,
  deleteMutation,
}: {
  bracket: TERBracket;
  onEdit: () => void;
  deleteMutation: UseMutationResult<any, Error, number>;
}) {
  return (
    <TableRow>
      <TableCell className="font-medium">{bracket.bracket_number}</TableCell>
      <TableCell>Rp {bracket.min_monthly_salary.toLocaleString("id-ID")}</TableCell>
      <TableCell>{(bracket.rate * 100).toFixed(2)}%</TableCell>
      <TableCell>{new Date(bracket.effective_from).toLocaleDateString("id-ID")}</TableCell>
      <TableCell className="text-right">
        <div className="flex justify-end gap-1">
          <Button
            variant="ghost"
            size="icon"
            onClick={onEdit}
            className="h-8 w-8"
          >
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => deleteMutation.mutate(bracket.id)}
            disabled={deleteMutation.isPending}
            className="h-8 w-8 text-destructive hover:text-destructive"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </TableCell>
    </TableRow>
  );
}

export function TerBracketsTable() {
  const [activeCategory, setActiveCategory] = useState("A");
  const [editingId, setEditingId] = useState<number | "new" | null>(null);
  const { data: brackets, isLoading } = useTerBrackets(activeCategory);
  const { createMutation, updateMutation, deleteMutation } = useTerBracketMutations();

  const sortedBrackets = (brackets ?? [])
    .filter((b) => b.category === activeCategory)
    .sort((a, b) => a.bracket_number - b.bracket_number);

  return (
    <div className="space-y-4">
      <Tabs value={activeCategory} onValueChange={setActiveCategory}>
        <TabsList>
          {TER_CATEGORIES.map((cat) => (
            <TabsTrigger key={cat} value={cat}>
              Category {cat}
            </TabsTrigger>
          ))}
        </TabsList>
      </Tabs>

      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Salary brackets for Category {activeCategory}.
        </p>
        <Button
          size="sm"
          onClick={() => setEditingId("new")}
          disabled={editingId === "new"}
        >
          <Plus className="mr-2 h-4 w-4" /> Add Bracket
        </Button>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : (
        <>
          <div className="hidden md:block">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>#</TableHead>
                  <TableHead>Min Monthly Salary</TableHead>
                  <TableHead>Rate</TableHead>
                  <TableHead>Effective From</TableHead>
                  <TableHead className="w-24 text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {editingId === "new" && (
                  <EditableRow
                    category={activeCategory}
                    createMutation={createMutation}
                    updateMutation={updateMutation}
                    deleteMutation={deleteMutation}
                    onCancel={() => setEditingId(null)}
                  />
                )}
                {sortedBrackets.map((bracket) =>
                  editingId === bracket.id ? (
                    <EditableRow
                      key={bracket.id}
                      bracket={bracket}
                      category={activeCategory}
                      createMutation={createMutation}
                      updateMutation={updateMutation}
                      deleteMutation={deleteMutation}
                      onCancel={() => setEditingId(null)}
                    />
                  ) : (
                    <StaticRow
                      key={bracket.id}
                      bracket={bracket}
                      onEdit={() => setEditingId(bracket.id)}
                      deleteMutation={deleteMutation}
                    />
                  )
                )}
                {sortedBrackets.length === 0 && editingId !== "new" && (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                      No brackets configured for Category {activeCategory}.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>

          <div className="md:hidden space-y-3">
            {editingId === "new" && (
              <Card className="p-4 border-2 border-blue-200">
                <p className="text-sm font-medium mb-2">New Bracket</p>
                <div className="space-y-2">
                  <div>
                    <span className="text-xs text-muted-foreground">Bracket #</span>
                    <Input
                      type="number"
                      value={sortedBrackets.length + 1}
                      className="h-8 mt-1"
                      placeholder="Bracket number"
                    />
                  </div>
                  <div>
                    <span className="text-xs text-muted-foreground">Min Monthly Salary (Rp)</span>
                    <Input type="number" className="h-8 mt-1" placeholder="Min salary" />
                  </div>
                  <div>
                    <span className="text-xs text-muted-foreground">Rate (%)</span>
                    <Input type="number" step="0.01" className="h-8 mt-1" placeholder="Rate" />
                  </div>
                  <div>
                    <span className="text-xs text-muted-foreground">Effective From</span>
                    <Input type="date" className="h-8 mt-1" />
                  </div>
                  <div className="flex gap-2 pt-1">
                    <Button size="sm" className="flex-1">Save</Button>
                    <Button size="sm" variant="ghost" onClick={() => setEditingId(null)}>
                      <X className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </Card>
            )}
            {sortedBrackets.map((bracket) =>
              editingId === bracket.id ? (
                <Card key={bracket.id} className="p-4 border-2 border-blue-200">
                  <p className="text-sm font-medium mb-2">Edit Bracket #{bracket.bracket_number}</p>
                  <div className="space-y-2">
                    <div>
                      <span className="text-xs text-muted-foreground">Bracket #</span>
                      <Input
                        type="number"
                        defaultValue={bracket.bracket_number}
                        className="h-8 mt-1"
                      />
                    </div>
                    <div>
                      <span className="text-xs text-muted-foreground">Min Monthly Salary (Rp)</span>
                      <Input
                        type="number"
                        defaultValue={bracket.min_monthly_salary}
                        className="h-8 mt-1"
                      />
                    </div>
                    <div>
                      <span className="text-xs text-muted-foreground">Rate (%)</span>
                      <Input
                        type="number"
                        step="0.01"
                        defaultValue={bracket.rate * 100}
                        className="h-8 mt-1"
                      />
                    </div>
                    <div>
                      <span className="text-xs text-muted-foreground">Effective From</span>
                      <Input
                        type="date"
                        defaultValue={bracket.effective_from.split("T")[0]}
                        className="h-8 mt-1"
                      />
                    </div>
                    <div className="flex gap-2 pt-1">
                      <Button size="sm" className="flex-1">Save</Button>
                      <Button size="sm" variant="ghost" onClick={() => setEditingId(null)}>
                        <X className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </Card>
              ) : (
                <Card key={bracket.id} className="p-4">
                  <div className="flex items-start justify-between">
                    <div>
                      <p className="font-medium">
                        <span className="text-muted-foreground font-mono text-sm">#{bracket.bracket_number}</span>
                        {" · "}
                        {(bracket.rate * 100).toFixed(2)}%
                      </p>
                      <p className="text-sm text-muted-foreground">
                        Min: Rp {bracket.min_monthly_salary.toLocaleString("id-ID")}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        Since {new Date(bracket.effective_from).toLocaleDateString("id-ID")}
                      </p>
                    </div>
                    <div className="flex gap-1">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => setEditingId(bracket.id)}
                        className="h-8 w-8"
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => deleteMutation.mutate(bracket.id)}
                        disabled={deleteMutation.isPending}
                        className="h-8 w-8 text-destructive hover:text-destructive"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </Card>
              )
            )}
            {sortedBrackets.length === 0 && editingId !== "new" && (
              <div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
                <p>No brackets configured for Category {activeCategory}.</p>
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}
