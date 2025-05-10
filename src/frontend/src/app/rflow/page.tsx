"use client";
import dynamic from "next/dynamic";

const BrickFlow = dynamic(() => import("@/components/BrickFlow"), {
  ssr: false,
});

export default function Page() {
  return <BrickFlow />;
}
