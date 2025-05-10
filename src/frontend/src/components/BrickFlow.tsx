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

const nodeDefaults = {
  sourcePosition: Position.Bottom,
  targetPosition: Position.Top,
};

let nodeIdCounter = 1;
let comboCounter = 1;
let recipeCounter = 1;
const generatedNodes: Node[] = [];
const generatedEdges: Edge[] = [];

const comboStyle = {
  width: 10,
  height: 10,
  background: "transparent",
  border: "none",
};

function traverse(node: recipeNode, parentId?: string, isRoot = false): string {
  const nodeId = `node${nodeIdCounter++}`;

  // @ts-expect-error ga kasih location disini, layouting dihandle @/utils/graphLayout
  generatedNodes.push({
    id: nodeId,
    data: { label: node.name },
    ...nodeDefaults,
  });

  if (node.recipes && node.recipes.length >= 2) {
    for (let i = 0; i + 1 < node.recipes.length; i += 2) {
      const left = node.recipes[i];
      const right = node.recipes[i + 1];

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
}

const recipe1: recipeNode = {
  name: "Brick",
  recipes: [
    {
      name: "Mud",
      recipes: [
        {
          name: "Water",
        },
        {
          name: "Earth",
        },
      ],
    },
    {
      name: "Fire",
    },
    {
      name: "Clay",
      recipes: [
        {
          name: "Mud",
          recipes: [
            {
              name: "Water",
            },
            {
              name: "Earth",
            },
          ],
        },
        {
          name: "Sand",
          recipes: [
            {
              name: "Stone",
              recipes: [
                {
                  name: "Lava",
                  recipes: [
                    {
                      name: "Earth",
                    },
                    {
                      name: "Fire",
                    },
                  ],
                },
                {
                  name: "Air",
                },
                {
                  name: "Earth",
                },
                {
                  name: "Pressure",
                  recipes: [{ name: "Air" }, { name: "Air" }],
                },
              ],
            },
            {
              name: "Air",
            },
          ],
        },
      ],
    },
    {
      name: "Stone",
      recipes: [
        {
          name: "Lava",
          recipes: [
            {
              name: "Earth",
            },
            {
              name: "Fire",
            },
          ],
        },
        {
          name: "Air",
        },
        {
          name: "Earth",
        },
        {
          name: "Pressure",
          recipes: [{ name: "Air" }, { name: "Air" }],
        },
      ],
    },
  ],
};

traverse(recipe1, undefined, true);

const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
  generatedNodes,
  generatedEdges,
  "TB"
);

export default function BrickFlow() {
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

type recipeNode = {
  name: string;
  recipes?: recipeNode[];
};
