"use client";

import { useEffect, useRef, useId } from "react";
import mermaid from "mermaid";

interface MermaidProps {
  chart: string;
}

export default function Mermaid({ chart }: MermaidProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const id = useId();

  useEffect(() => {
    mermaid.initialize({ startOnLoad: true });
    if (containerRef.current) {
      containerRef.current.innerHTML = `<div class="mermaid" id="${id}">${chart}</div>`;
      mermaid.contentLoaded();
    }
  }, [chart, id]);

  return <div ref={containerRef} />;
}
