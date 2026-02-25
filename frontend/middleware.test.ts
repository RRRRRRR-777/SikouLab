import { describe, expect, it } from "vitest";

import { __internal__ } from "./middleware";

describe("middleware path matching", () => {
  it("treats '/' as exact match only", () => {
    expect(__internal__.matchesPath("/", "/")).toBe(true);
    expect(__internal__.matchesPath("/login", "/")).toBe(false);
    expect(__internal__.matchesPath("/articles", "/")).toBe(false);
  });

  it("matches exact path and sub paths for non-root routes", () => {
    expect(__internal__.matchesPath("/articles", "/articles")).toBe(true);
    expect(__internal__.matchesPath("/articles/123", "/articles")).toBe(true);
    expect(__internal__.matchesPath("/articles-legacy", "/articles")).toBe(false);
  });

  it("keeps /login as public and protected routes as expected", () => {
    expect(__internal__.isPublicPath("/login")).toBe(true);
    expect(__internal__.isPublicPath("/subscription")).toBe(true);
    expect(__internal__.isProtectedPath("/login")).toBe(false);
    expect(__internal__.isProtectedPath("/")).toBe(true);
    expect(__internal__.isProtectedPath("/articles/2026")).toBe(true);
  });
});
