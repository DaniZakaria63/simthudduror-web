import { render } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { useReveal } from "./useReveal";

function Harness({ routeKey = "/" }: { routeKey?: string }) {
  useReveal(routeKey);
  return (
    <div>
      <p data-reveal>Pertama</p>
      <p data-reveal>Kedua</p>
    </div>
  );
}

class FakeIntersectionObserver {
  static instances: FakeIntersectionObserver[] = [];
  callback: IntersectionObserverCallback;
  observed: Element[] = [];
  disconnected = false;

  constructor(callback: IntersectionObserverCallback) {
    this.callback = callback;
    FakeIntersectionObserver.instances.push(this);
  }

  observe(element: Element) {
    this.observed.push(element);
  }

  unobserve(element: Element) {
    this.observed = this.observed.filter((item) => item !== element);
  }

  disconnect() {
    this.disconnected = true;
  }

  trigger(element: Element) {
    this.callback(
      [{ target: element, isIntersecting: true } as IntersectionObserverEntry],
      this as unknown as IntersectionObserver,
    );
  }
}

describe("useReveal", () => {
  it("observes marked nodes and reveals them once", () => {
    vi.stubGlobal("IntersectionObserver", FakeIntersectionObserver);
    const { container, unmount } = render(<Harness />);
    const observer = FakeIntersectionObserver.instances.at(-1);
    expect(observer).toBeDefined();
    expect(observer?.observed).toHaveLength(2);
    expect(document.documentElement.classList.contains("reveal-ready")).toBe(true);

    const first = container.querySelector("[data-reveal]") as HTMLElement;
    observer?.trigger(first);
    expect(first.classList.contains("is-in")).toBe(true);
    expect(observer?.observed).toHaveLength(1);

    unmount();
    expect(observer?.disconnected).toBe(true);
    expect(document.documentElement.classList.contains("reveal-ready")).toBe(false);
  });

  it("reveals immediately for reduced motion", () => {
    vi.stubGlobal("matchMedia", (query: string) => ({
      matches: true,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    }));
    const { container } = render(<Harness />);
    expect(document.documentElement.classList.contains("reveal-ready")).toBe(true);
    expect(container.querySelectorAll(".is-in")).toHaveLength(2);
  });

  it("reveals immediately when IntersectionObserver is unavailable", () => {
    vi.stubGlobal("IntersectionObserver", undefined);
    const { container } = render(<Harness />);
    expect(container.querySelectorAll(".is-in")).toHaveLength(2);
  });
});
