import { describe, expect, it } from "vitest";
import { authorInitial, formatDate, formatDateTime } from "./format";

describe("format helpers", () => {
  it("formats Indonesian dates and date-times", () => {
    expect(formatDate("2026-11-14T08:00:00.000Z")).toContain("2026");
    expect(formatDate("2026-11-14T08:00:00.000Z")).toContain("November");
    expect(formatDateTime("2026-11-14T08:00:00.000Z")).toMatch(/\d{2}[.:]\d{2}/);
  });

  it("returns an empty string for invalid dates", () => {
    expect(formatDate("bukan tanggal")).toBe("");
    expect(formatDateTime("")).toBe("");
  });

  it("derives an author initial", () => {
    expect(authorInitial("rina")).toBe("R");
    expect(authorInitial("  bima")).toBe("B");
    expect(authorInitial("   ")).toBe("M");
  });
});
