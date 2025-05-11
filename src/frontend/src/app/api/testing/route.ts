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
  const algo = params.algo == undefined ? "BFS" : params.algo;
  const max = params.max == undefined ? -1 : parseInt(params.max);
  const target = params.target == undefined ? "Brick" : params.target;
  // const mode = params.mode == undefined ? "Shortest" : params.mode;
  const isShortest =
    params.shortest == undefined
      ? false
      : params.shortest == "true"
      ? true
      : false;

  const result: JSONResponse = {
    data: target === "Sheet Music" ? data1 : data2,
    errors: [],
    time: 43,
    nodeCount: 100,
    recipefound: 2,
  };

  return NextResponse.json(result);
}
