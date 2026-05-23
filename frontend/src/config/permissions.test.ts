import { describe, it, expect } from "vitest";
import { PERMISSIONS } from "@/config/permissions";

describe("PERMISSIONS", () => {
  it("should have permission keys matching their values", () => {
    Object.entries(PERMISSIONS).forEach(([key, value]) => {
      expect(key).toBe(value);
    });
  });

  it("should contain employee permissions", () => {
    expect(PERMISSIONS.VIEW_EMPLOYEE).toBe("VIEW_EMPLOYEE");
    expect(PERMISSIONS.CREATE_EMPLOYEE).toBe("CREATE_EMPLOYEE");
    expect(PERMISSIONS.UPDATE_EMPLOYEE).toBe("UPDATE_EMPLOYEE");
    expect(PERMISSIONS.DELETE_EMPLOYEE).toBe("DELETE_EMPLOYEE");
    expect(PERMISSIONS.EXPORT_EMPLOYEE).toBe("EXPORT_EMPLOYEE");
  });

  it("should contain payroll permissions", () => {
    expect(PERMISSIONS.VIEW_PAYROLL).toBe("VIEW_PAYROLL");
    expect(PERMISSIONS.GENERATE_PAYROLL).toBe("GENERATE_PAYROLL");
    expect(PERMISSIONS.DOWNLOAD_PAYSLIP).toBe("DOWNLOAD_PAYSLIP");
    expect(PERMISSIONS.MARK_AS_PAID).toBe("MARK_AS_PAID");
    expect(PERMISSIONS.SEND_PAYSLIP).toBe("SEND_PAYSLIP");
  });

  it("should contain leave permissions", () => {
    expect(PERMISSIONS.VIEW_LEAVE).toBe("VIEW_LEAVE");
    expect(PERMISSIONS.VIEW_SELF_LEAVE).toBe("VIEW_SELF_LEAVE");
    expect(PERMISSIONS.CREATE_LEAVE).toBe("CREATE_LEAVE");
    expect(PERMISSIONS.APPROVAL_LEAVE).toBe("APPROVAL_LEAVE");
    expect(PERMISSIONS.EXPORT_LEAVE).toBe("EXPORT_LEAVE");
  });

  it("should have all expected permission categories", () => {
    const keys = Object.keys(PERMISSIONS);
    const categories = [
      "PERMISSION",
      "ROLE",
      "MASTER",
      "EMPLOYEE",
      "ATTENDANCE",
      "PAYROLL",
      "LEAVE",
      "LOAN",
      "OVERTIME",
      "REIMBURSEMENT",
      "COMPANY",
      "ANNOUNCEMENT",
      "CONTRACT",
      "REQUISITION",
      "APPLICANT",
      "ONBOARDING",
      "FINANCE",
    ];
    categories.forEach((cat) => {
      expect(keys.some((k) => k.includes(cat))).toBe(true);
    });
  });
});
