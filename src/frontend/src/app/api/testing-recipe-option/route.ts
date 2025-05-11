// this is a testing route, will be removed in the future
// in the future, api is served from the backend
import { NextResponse } from "next/server";

type RecipeOption = {
  tier: number;
  name: string;
};

const RECIPE_OPTIONS: RecipeOption[] = [
  { tier: 0, name: "Brick" },
  { tier: 0, name: "Wood" },
  { tier: 0, name: "Steel" },
  { tier: 0, name: "Glass" },
  { tier: 0, name: "Plastic" },
  { tier: 0, name: "Concrete" },
  { tier: 0, name: "Ceramic" },
  { tier: 0, name: "Aluminum" },
];

export async function GET() {
  return NextResponse.json(RECIPE_OPTIONS);
}
