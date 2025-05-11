"use client";

import ReactFlow, {
  Background,
  Controls,
  MiniMap,
  Position,
  Node,
  Edge,
  MarkerType,
} from "reactflow";
import "reactflow/dist/style.css";
import { getLayoutedElements } from "@/utils/graphLayout";

type RecipeNode = {
  name: string;
  recipes?: [RecipeNode, RecipeNode][]; // each recipe requires 2 ingredients
};

type recipeTreeProps = {
  data: RecipeNode;
  direction?: "TB" | "LR"; // Top-to-Bottom or Left-to-Right layout
};

const nodeDefaults = {
  sourcePosition: Position.Bottom,
  targetPosition: Position.Top,
};

const comboStyle = {
  width: 10,
  height: 10,
  background: "transparent",
  border: "none",
};

const generateGraphElements = (
  data: RecipeNode,
  direction: "TB" | "LR" = "TB"
) => {
  let nodeIdCounter = 1;
  let comboCounter = 1;
  let recipeCounter = 1;
  const generatedNodes: Node[] = [];
  const generatedEdges: Edge[] = [];

  const traverse = (
    node: RecipeNode,
    parentId?: string,
    isRoot = false
  ): string => {
    const nodeId = `node${nodeIdCounter++}`;

    // @ts-expect-error ga kasih location disini, layouting dihandle @/utils/graphLayout
    generatedNodes.push({
      id: nodeId,
      data: { label: node.name },
      ...nodeDefaults,
    });

    if (node.recipes && node.recipes.length > 0) {
      for (const [left, right] of node.recipes) {
        const comboId = `combo${comboCounter++}`;

        // @ts-expect-error ga kasih location disini, layouting dihandle @/utils/graphLayout
        generatedNodes.push({
          id: comboId,
          data: { label: "" },
          style: comboStyle,
          ...nodeDefaults,
        });

        const leftId = traverse(left);
        const rightId = traverse(right);

        generatedEdges.push({
          id: `e-${leftId}-${comboId}`,
          source: leftId,
          target: comboId,
          label: "combine",
          style: { strokeDasharray: "5 5" },
        });

        generatedEdges.push({
          id: `e-${rightId}-${comboId}`,
          source: rightId,
          target: comboId,
          label: "combine",
          style: { strokeDasharray: "5 5" },
        });

        generatedEdges.push({
          id: `e-${comboId}-${nodeId}`,
          source: comboId,
          target: nodeId,
          label: isRoot ? `Recipe ${recipeCounter++}` : undefined,
          style: { stroke: "#4caf50", strokeWidth: 2 },
          markerEnd: { type: MarkerType.ArrowClosed },
        });
      }
    }

    return nodeId;
  };

  traverse(data, undefined, true);

  return getLayoutedElements(generatedNodes, generatedEdges, direction);
};

export default function recipeTree({
  data,
  direction = "TB",
}: recipeTreeProps) {
  const { nodes: layoutedNodes, edges: layoutedEdges } = generateGraphElements(
    data,
    direction
  );

  return (
    <div style={{ width: "100%", height: "100vh" }}>
      <ReactFlow
        nodes={layoutedNodes}
        edges={layoutedEdges}
        fitView
        defaultEdgeOptions={{ style: { strokeWidth: 2 } }}
      >
        <Background />
        <Controls />
        <MiniMap />
      </ReactFlow>
    </div>
  );
}
