"use client";
import dynamic from "next/dynamic";
import Link from "next/link";
import { useEffect, useState } from "react";

const Mermaid = dynamic(() => import("@/components/Mermaid"), { ssr: false });

const diagramInitial = `
graph TD
    A --> B
    B --> C
    B --> D
`;

const toAddInitial = ["D --> F", "F --> G", "C --> H", "H --> G"];
export default function Dummy() {
  const [diagram, setDiagram] = useState(diagramInitial);
  const [toAdd, setToAdd] = useState(toAddInitial);

  useEffect(() => {
    if (toAdd.length === 0) return;

    const timeout = setTimeout(() => {
      const [next, ...rest] = toAdd;
      setDiagram((prev) => prev + "\n" + next);
      setToAdd(rest);
    }, 2000);

    return () => clearTimeout(timeout);
  }, [toAdd]);

  return (
    <div>
      <Link
        className=" bg-blue-500 text-white font-bold py-2 px-4 rounded hover:bg-blue-700"
        href={"/"}
      >
        Home
      </Link>
      <h1>Mermaid Diagram</h1>
      <Mermaid chart={diagram} />
      <div>aaa</div>
    </div>
  );
}
