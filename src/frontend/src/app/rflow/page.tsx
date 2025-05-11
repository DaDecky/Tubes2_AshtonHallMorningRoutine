"use client";
import dynamic from "next/dynamic";

const RecipeTree = dynamic(() => import("@/components/RecipeTree"), {
  ssr: false,
});
// import { data } from "@/testing/sheet_music_recipes";
import { data } from "@/testing/brick";

export default function Page() {
  return <RecipeTree data={data} />;
}
