"use client";
import dynamic from "next/dynamic";
import { useEffect, useReducer, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { Checkbox } from "@/components/ui/checkbox";
import { AlertCircle } from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

const RecipeTree = dynamic(() => import("@/components/RecipeTree"), {
  ssr: false,
  loading: () => (
    <div className="flex justify-center items-center h-64">
      <div className="animate-pulse">Loading visualization...</div>
    </div>
  ),
});

type RecipeNode = {
  name: string;
  recipes?: [RecipeNode, RecipeNode][];
};

type State = {
  data: RecipeNode | undefined;
  errors: string[];
  time: number;
  nodeCount: number;
  recipeFound: number;
  isLoading: boolean;
  formError?: string;
};

type Action =
  | { type: "FETCH_START" }
  | { type: "FETCH_SUCCESS"; payload: Omit<State, "isLoading" | "formError"> }
  | { type: "FETCH_ERROR"; error: string }
  | { type: "FORM_ERROR"; error: string };

const initialState: State = {
  data: undefined,
  errors: [],
  time: 0,
  nodeCount: 0,
  recipeFound: 0,
  isLoading: false,
  formError: undefined,
};

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case "FETCH_START":
      return { ...state, isLoading: true, errors: [], formError: undefined };
    case "FETCH_SUCCESS":
      return { ...action.payload, isLoading: false };
    case "FETCH_ERROR":
      return {
        ...state,
        errors: [...state.errors, action.error],
        isLoading: false,
      };
    case "FORM_ERROR":
      return {
        ...state,
        formError: action.error,
        isLoading: false,
      };
    default:
      return state;
  }
}

type RecipeOption = {
  tier: number;
  name: string;
};

export default function Page() {
  const [state, dispatch] = useReducer(reducer, initialState);
  const [target, setTarget] = useState<string>("");
  const [algorithm, setAlgorithm] = useState("BFS");
  const [shortestPath, setShortestPath] = useState(false);
  const [maxRecipes, setMaxRecipes] = useState(1);
  const [open, setOpen] = useState(false);
  const [searchValue, setSearchValue] = useState("");
  const [recipeOptions, setRecipeOptions] = useState<RecipeOption[]>();

  
  useEffect(() => {
    
    const fetchRecipeOptions = async () => {
      try {
        const base_url = process.env.BACKEND_URL || "http://localhost:8081";
        const response = await fetch(base_url + "/elements", );
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const jsonData = await response.json();
        setRecipeOptions(jsonData);
      } catch (error) {
        console.error(
          error instanceof Error ? error.message : "Unknown error occurred"
        );
      }
    };
    fetchRecipeOptions();
  }, []);

  const filteredOptions = recipeOptions?.filter((option) =>
    option.name.toLowerCase().includes(searchValue.toLowerCase())
  );

  const fetchData = async (
    target: string,
    algorithm: string,
    shortest: boolean
  ) => {
    if (!target) {
      dispatch({ type: "FORM_ERROR", error: "Please select a recipe" });
      return;
    }

    dispatch({ type: "FETCH_START" });

    try {
      const base_url = process.env.BACKEND_URL || "http://localhost:8081";
      const url = base_url + (shortest
        ? `/search?target=${encodeURIComponent(
            target
          )}&algo=${algorithm}&shortest=true`
        : `/search?target=${encodeURIComponent(
            target
          )}&algo=${algorithm}&max=${maxRecipes}`);

      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const jsonData = await response.json();
      dispatch({
        type: "FETCH_SUCCESS",
        payload: {
          data: jsonData.data,
          errors: jsonData.errors || [],
          time: jsonData.time,
          nodeCount: jsonData.nodeCount,
          recipeFound: jsonData.recipefound,
        },
      });
    } catch (error) {
      dispatch({
        type: "FETCH_ERROR",
        error:
          error instanceof Error ? error.message : "Unknown error occurred",
      });
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    fetchData(target, algorithm, shortestPath);
  };

  const targetRecipe = recipeOptions?.find(
    (option) => option.name === target
  )?.name;

  return (
    <main className="min-h-screen p-4 max-w-6xl mx-auto space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Recipe Explorer</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="target">Target Recipe *</Label>
              <Popover open={open} onOpenChange={setOpen}>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={open}
                    className="w-full justify-between"
                  >
                    {targetRecipe || "Select recipe..."}
                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-[var(--radix-popover-trigger-width)] p-0">
                  <Command className="w-full">
                    <CommandInput
                      placeholder="Search recipes..."
                      value={searchValue}
                      onValueChange={setSearchValue}
                    />
                    <CommandEmpty>No recipe found.</CommandEmpty>
                    <CommandGroup className="max-h-60 overflow-y-auto w-full">
                      {filteredOptions?.map((option) => (
                        <CommandItem
                          key={option.name}
                          value={option.name}
                          onSelect={(currentValue) => {
                            setTarget(
                              currentValue === target ? "" : currentValue
                            );
                            setOpen(false);
                          }}
                          className="w-full"
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              target === option.name
                                ? "opacity-100"
                                : "opacity-0"
                            )}
                          />
                          {option.name} - Tier {option.tier}
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  </Command>
                </PopoverContent>
              </Popover>
              {state.formError && (
                <p className="text-sm font-medium text-destructive flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {state.formError}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="algorithm">Search Algorithm</Label>
              <select
                id="algorithm"
                value={algorithm}
                onChange={(e) => setAlgorithm(e.target.value)}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
              >
                <option value="BFS">Breadth-First Search (BFS)</option>
                <option value="DFS">Depth-First Search (DFS)</option>
              </select>
            </div>

            <div className="flex items-center space-x-2">
              <Checkbox
                id="shortestPath"
                checked={shortestPath}
                onCheckedChange={(checked) => setShortestPath(checked === true)}
              />
              <Label htmlFor="shortestPath">
                Find shortest path only (1 recipe)
              </Label>
            </div>

            {!shortestPath && (
              <div className="space-y-2">
                <Label htmlFor="maxRecipes">Maximum Recipes</Label>
                <Input
                  id="maxRecipes"
                  type="number"
                  min="1"
                  max="50"
                  value={maxRecipes}
                  onChange={(e) => setMaxRecipes(Number(e.target.value))}
                  placeholder="Default: 10"
                />
              </div>
            )}

            <Tooltip>
              <TooltipTrigger asChild>
                <div className="w-full">
                  <Button
                    type="submit"
                    disabled={state.isLoading || !target}
                    className="w-full"
                  >
                    {state.isLoading ? "Searching..." : "Search Recipes"}
                  </Button>
                </div>
              </TooltipTrigger>
              {!target && (
                <TooltipContent
                  side="top"
                  className="bg-transparent text-destructive-foreground text-4xl"
                >
                  <p>Please select a target recipe first</p>
                </TooltipContent>
              )}
            </Tooltip>
          </form>

          {state.data && (
            <div className="mt-4 text-sm text-muted-foreground space-y-1">
              <p>
                Found {state.recipeFound} recipes in {state.time}ms
              </p>
              <p>Visited {state.nodeCount} nodes</p>
              <p>Using {algorithm} algorithm</p>
              {shortestPath ? (
                <p>Showing shortest path only</p>
              ) : (
                <p>Showing up to {maxRecipes} recipes</p>
              )}
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardContent className="pt-6">
          {state.isLoading ? (
            <div className="flex justify-center items-center h-64">
              <div className="animate-pulse">Loading recipe data...</div>
            </div>
          ) : state.errors.length > 0 ? (
            <div className="text-destructive space-y-2">
              {state.errors.map((error, i) => (
                <p key={i}>{error}</p>
              ))}
              <Button
                onClick={() => fetchData(target, algorithm, shortestPath)}
                variant="outline"
              >
                Retry
              </Button>
            </div>
          ) : state.data ? (
            <div className="h-[600px]">
              <RecipeTree data={state.data} />
            </div>
          ) : (
            <div className="text-center py-12 text-muted-foreground">
              <p>No recipe data available</p>
              <p className="text-sm mt-2">{`Select a recipe and click "Search Recipes"`}</p>
            </div>
          )}
        </CardContent>
      </Card>
    </main>
  );
}
