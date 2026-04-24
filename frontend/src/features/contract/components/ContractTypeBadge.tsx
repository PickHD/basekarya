import { Badge } from "@/components/ui/badge";
import type { ContractType } from "../types";

interface Props {
  type: ContractType;
}

export function ContractTypeBadge({ type }: Props) {
  if (type === "PKWT") {
    return (
      <Badge variant="outline" className="bg-blue-50 text-blue-700 border-blue-200">
        PKWT
      </Badge>
    );
  }

  if (type === "PKWTT") {
    return (
      <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
        PKWTT
      </Badge>
    );
  }

  return <Badge variant="outline">{type}</Badge>;
}
