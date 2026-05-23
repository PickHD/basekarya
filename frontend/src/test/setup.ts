import "@testing-library/jest-dom/vitest";

Object.defineProperty(globalThis, "IntersectionObserver", {
  writable: true,
  value: class IntersectionObserver {
    private callback: IntersectionObserverCallback;
    private options: IntersectionObserverInit;
    private elements: Element[] = [];

    constructor(callback: IntersectionObserverCallback, options?: IntersectionObserverInit) {
      this.callback = callback;
      this.options = options || {};
    }

    observe(element: Element) {
      this.elements.push(element);
    }

    unobserve(_element: Element) {
      void _element;
    }

    disconnect() {
      this.elements = [];
    }

    triggerIntersection(isIntersecting: boolean) {
      this.callback(
        isIntersecting
          ? ([{ isIntersecting, target: this.elements[0] }] as unknown as IntersectionObserverEntry[])
          : ([{ isIntersecting: false, target: this.elements[0] }] as unknown as IntersectionObserverEntry[]),
        this
      );
    }
  },
});
