"use client";
import dynamic from "next/dynamic";

const RecipeTree = dynamic(() => import("@/components/RecipeTree"), {
  ssr: false,
  loading: () => (
    <div className="flex justify-center items-center h-screen">
      <div className="animate-pulse">Loading recipe visualization...</div>
    </div>
  ),
});
// import { data as data1 } from "@/testing/sheet_music_recipes";
import { data as data2 } from "@/testing/brick";
import { useState } from "react";

type RecipeNode = {
  name: string;
  recipes?: [RecipeNode, RecipeNode][]; // each recipe requires 2 ingredients
};

export default function Page() {
  const [data] = useState<RecipeNode>(data2);
  // toggle between use image or not
  // const [useImage, setUseImage] = useState(false);

  return (
    <main className="">
      <RecipeTree data={data} />
    </main>
  );
}
