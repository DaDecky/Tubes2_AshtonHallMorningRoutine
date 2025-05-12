// this is a testing route, will be removed in the future
// in the future, api is served from the backend
import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { data as data1 } from "@/testing/sheet_music_recipes";
import { data as data2 } from "@/testing/brick";

type RecipeNode = {
  name: string;
  recipes?: [RecipeNode, RecipeNode][]; // each recipe requires 2 ingredients
};

type JSONResponse = {
  data: RecipeNode;
  errors: []; // list of error string
  time: number; // waktu nyari in ms
  nodeCount: number; //banyak node dikunjungi
  recipefound: number; //banyak resep yang didapet -> mungkin kurang dari max recipe yang di set
};

/* eslint-disable @typescript-eslint/no-unused-vars */
export async function GET(request: NextRequest) {
  const { searchParams } = new URL(request.url);
  const params = Object.fromEntries(searchParams.entries());
  const target = params.target == undefined ? "Brick" : params.target;
  const algo = params.algo == undefined ? "BFS" : params.algo;


  let url = `http://localhost:8081/search?target=${target}&algo=${algo}`
  if (params.shortest != undefined) {
    url += "&shortest=true"
  }
  if (params.max != undefined) {
    url += "&max=" + params.max
  }

  const response = await fetch(url);

  let jsonData = await response.json();

  return NextResponse.json(jsonData);
}
