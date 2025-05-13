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
import Image from "next/image";
import { useEffect, useState } from "react";

type RecipeNode = {
  name: string;
  recipes?: [RecipeNode, RecipeNode][];
};

type recipeTreeProps = {
  data: RecipeNode | undefined;
  direction?: "TB" | "LR";
  useImage?: boolean;
  algorithm?: "BFS" | "DFS";
  speed?: number
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
  direction: "TB" | "LR" = "TB",
  algorithm: "BFS" | "DFS" = "BFS"
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

    generatedNodes.push({
      id: nodeId,
      data: {
        label: (
          <div className="flex items-center">
            <div className="w-12 h-12 relative">
              <Image
                src={`/images/${node.name}.svg`}
                alt={node.name}
                fill
              />
            </div>
            <span className="text-xl">{node.name}</span>
          </div>
        ),
      },
      ...nodeDefaults,
    } as Node);

    if (node.recipes && node.recipes.length > 0) {
      for (const [left, right] of node.recipes) {
        const comboId = `combo${comboCounter++}`;

        generatedNodes.push({
          id: comboId,
          data: { label: "" },
          style: comboStyle,
          ...nodeDefaults,
        } as Node);

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

const generateGraphElementsAnimated = (
  data: RecipeNode,
  direction: "TB" | "LR" = "TB",
  algorithm: "BFS" | "DFS" = "BFS"
) => {
  const allNodes: Node[] = [];
  const allEdges: Edge[] = [];
  let nodeIdCounter = 1;
  let comboCounter = 1;
  let recipeCounter = 1;

  if (algorithm === "BFS") {
    // BFS implementation
    const queue: { node: RecipeNode; parentId?: string; isRoot?: boolean }[] = [];
    queue.push({ node: data, isRoot: true });

    while (queue.length > 0) {
      const current = queue.shift()!;
      const nodeId = `node${nodeIdCounter++}`;

      // Add current node
      allNodes.push({
        id: nodeId,
        data: {
          label: (
            <div className="flex items-center">
              <div className="w-12 h-12 relative">
                <Image
                  src={`/images/${current.node.name}.svg`}
                  alt={current.node.name}
                  fill
                />
              </div>
              <span className="text-xl">{current.node.name}</span>
            </div>
          ),
        },
        ...nodeDefaults,
      } as Node);

      if (current.node.recipes && current.node.recipes.length > 0) {
        for (const [left, right] of current.node.recipes) {
          const comboId = `combo${comboCounter++}`;

          // Add combo node
          allNodes.push({
            id: comboId,
            data: { label: "" },
            style: comboStyle,
            ...nodeDefaults,
          } as Node);

          // Add edges
          allEdges.push({
            id: `e-${comboId}-${nodeId}`,
            source: comboId,
            target: nodeId,
            label: current.isRoot ? `Recipe ${recipeCounter++}` : undefined,
            style: { stroke: "#4caf50", strokeWidth: 2 },
            markerEnd: { type: MarkerType.ArrowClosed },
          });

          // Process children (BFS)
          queue.push({ node: left, parentId: comboId });
          queue.push({ node: right, parentId: comboId });
        }
      }

      // If this node has a parent, connect it
      if (current.parentId) {
        allEdges.push({
          id: `e-${nodeId}-${current.parentId}`,
          source: nodeId,
          target: current.parentId,
          label: "combine",
          style: { strokeDasharray: "5 5" },
        });
      }
    }
  } else {
    // DFS implementation
    const stack: { node: RecipeNode; parentId?: string; isRoot?: boolean }[] = [];
    stack.push({ node: data, isRoot: true });

    while (stack.length > 0) {
      const current = stack.pop()!;
      const nodeId = `node${nodeIdCounter++}`;

      // Add current node
      allNodes.push({
        id: nodeId,
        data: {
          label: (
            <div className="flex items-center">
              <div className="w-12 h-12 relative">
                <Image
                  src={`/images/${current.node.name}.svg`}
                  alt={current.node.name}
                  fill
                />
              </div>
              <span className="text-xl">{current.node.name}</span>
            </div>
          ),
        },
        ...nodeDefaults,
      } as Node);

      if (current.node.recipes && current.node.recipes.length > 0) {
        // For DFS, we need to process children in reverse order
        const recipes = [...current.node.recipes].reverse();
        
        for (const [left, right] of recipes) {
          const comboId = `combo${comboCounter++}`;

          // Add combo node
          allNodes.push({
            id: comboId,
            data: { label: "" },
            style: comboStyle,
            ...nodeDefaults,
          } as Node);

          // Add edges
          allEdges.push({
            id: `e-${comboId}-${nodeId}`,
            source: comboId,
            target: nodeId,
            label: current.isRoot ? `Recipe ${recipeCounter++}` : undefined,
            style: { stroke: "#4caf50", strokeWidth: 2 },
            markerEnd: { type: MarkerType.ArrowClosed },
          });

          // Process children (DFS)
          stack.push({ node: right, parentId: comboId });
          stack.push({ node: left, parentId: comboId });
        }
      }

      // If this node has a parent, connect it
      if (current.parentId) {
        allEdges.push({
          id: `e-${nodeId}-${current.parentId}`,
          source: nodeId,
          target: current.parentId,
          label: "combine",
          style: { strokeDasharray: "5 5" },
        });
      }
    }
  }

  return getLayoutedElements(allNodes, allEdges, direction);
};

export default function RecipeTree({
  data,
  direction = "TB",
  algorithm = "BFS",
  speed = 500,
}: recipeTreeProps) {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  useEffect(() => {
    if (!data) return;
    const { nodes: allNodes, edges: allEdges } = generateGraphElementsAnimated(
      data,
      direction,
      algorithm
    );

    // Clear previous nodes and edges
    setNodes([]);
    setEdges([]);
    
    const nodeInterval = setInterval(() => {
      setNodes((prevNodes) => {
        if (allNodes.length > prevNodes.length) {
          const newNodes = [...prevNodes, allNodes[prevNodes.length]];

          setEdges((prevEdges) => {
            if (allEdges.length > prevEdges.length) {
              return [...prevEdges, allEdges[prevEdges.length]];
            }
            return prevEdges;
          });

          return newNodes;
        } else {
          clearInterval(nodeInterval);
          return prevNodes;
        }
      });
    }, speed);

    
    return () => clearInterval(nodeInterval);
  }, [data, direction, algorithm, speed]);

  
  return (
    <div style={{ width: "100%", height: "100%" }}>
      
      <ReactFlow
        nodes={nodes}
        edges={edges}
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