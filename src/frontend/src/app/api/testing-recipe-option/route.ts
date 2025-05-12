// this is a testing route, will be removed in the future
// in the future, api is served from the backend
import { NextResponse } from "next/server";

type RecipeOption = {
  tier: number;
  name: string;
};

const response = await fetch("http://localhost:8081/elements");
let jsonData = await response.json();


const RECIPE_OPTIONS: RecipeOption[] = jsonData

export async function GET() {
  return NextResponse.json(RECIPE_OPTIONS);
}
