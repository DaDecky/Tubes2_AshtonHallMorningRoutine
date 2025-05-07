"use client";

import { useEffect, useRef, useId, useState } from "react";
import mermaid from "mermaid";

interface MermaidProps {
  chart: string;
}

export default function Mermaid({ chart }: MermaidProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const renderId = useId();
  const [svg, setSvg] = useState<string>("");

  useEffect(() => {
    mermaid.initialize({ startOnLoad: false });
    mermaid
      .render(renderId, chart)
      .then(({ svg }) => {
        setSvg(svg);
      })
      .catch((err) => {
        console.error("Mermaid render error", err);
      });
  }, [chart, renderId]);

  return <div ref={containerRef} dangerouslySetInnerHTML={{ __html: svg }} />;
}
