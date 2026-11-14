import { useLayoutEffect } from "react";

/**
 * Reveals marked sections once as they enter. The hidden state only applies
 * after `reveal-ready` is set, so content stays visible without JavaScript.
 */
export function useReveal(routeKey: string): void {
  useLayoutEffect(() => {
    const root = document.documentElement;
    const nodes = Array.from(document.querySelectorAll<HTMLElement>("[data-reveal]"));
    const reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    if (reduced || typeof IntersectionObserver === "undefined") {
      root.classList.add("reveal-ready");
      for (const node of nodes) {
        node.classList.add("is-in");
      }
      return () => {
        root.classList.remove("reveal-ready");
      };
    }

    root.classList.add("reveal-ready");
    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            entry.target.classList.add("is-in");
            observer.unobserve(entry.target);
          }
        }
      },
      { rootMargin: "0px 0px -8% 0px", threshold: 0.12 },
    );
    for (const node of nodes) {
      if (node.classList.contains("is-in")) {
        continue;
      }
      observer.observe(node);
    }
    return () => {
      observer.disconnect();
      root.classList.remove("reveal-ready");
    };
  }, [routeKey]);
}
