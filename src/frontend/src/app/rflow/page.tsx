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
import { data as data1 } from "@/testing/sheet_music_recipes";
import { data as data2 } from "@/testing/brick";

export default function Page() {
  // toggle between use image or not
  // const [useImage, setUseImage] = useState(false);

  return (
    <main className="">
      <RecipeTree data={data1} />
    </main>
  );
}
